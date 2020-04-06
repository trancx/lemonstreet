package vrfapi

import (
	"context"
	"fmt"
	"github.com/bilibili/kratos/pkg/ecode"
	"net/http"
	"strings"
	"sync"
	"github.com/bilibili/kratos/pkg/log"

	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"
	"google.golang.org/grpc"
)

// AppID .
const (
	_defaultDomain         = "localhost"
	_defaultCookieName     = "login_cookie"
	_defaultCookieLifeTime = 2592000
	AppID = "verify.service"
)


type Verify struct {
	lock   sync.RWMutex
	keys   map[int64]string
	client VerifyClient
}

func New() *Verify {
	c, err := newRPCVerifyClient(nil)
	if err != nil {
		panic(err)
	}
	v := &Verify{
		keys:   make(map[int64]string),
		client: c,
	}
	return v
}

// 1. Any cookies? Yes to_2, ELSE to_3
// 2. Cookies outdated? Yse to_3, ELSE to_4
// 3. Required Login
// 4. GetCookies Data
// 5. Verify Valid? (Save caches, etc... ) Yes go on,  ELSE refuse connection
// FIXME: LOG and Change on the token
func (v *Verify) verify(ctx *bm.Context) error {
	var (
		uid     int64
		ctoken 	string
		req     *TokenReq
		err     error
		cookies *http.Cookie
	)

	if cookies, err = ctx.Request.Cookie(_defaultCookieName); err != nil {
		log.Error("Cookie Invalid")
		return ecode.AccessDenied
	}
	if _, err = fmt.Sscanf(cookies.Value, "uid=%d&token=%s", &uid, &ctoken); err != nil {
		log.Error("Cookie Value Invalid")
		return ecode.AccessDenied
	}
	log.Info("uid = %d, token=%s", uid, ctoken)
	req = &TokenReq{
		Tk: &Token{
			Id:  uid,
			Key: ctoken,
		},
	}

	v.lock.RLock()
	token, ok := v.keys[uid]
	v.lock.RUnlock()

	if ok {
		if strings.EqualFold(ctoken, token) {
			ctx.Set("uid", uid)
			return nil
		}
	}

	// not found or not equal
	{
		reply, err := v.client.VrfKey(ctx, req)
		if err != nil {
			// FIXME: RPC error
			log.Error("RPC error")
			return ecode.ServerErr
		}
		if reply.IsValid {
			// means token is right one and we can cached it!
			v.lock.Lock()
			v.keys[uid] = ctoken // RPC will return nil, when invalid, see gRPC handler
			v.lock.Unlock()
			ctx.Set("uid", uid)
			return nil
		} else {
			// if reply.IsUpdated, not the right one, but we can cached it
			// for logout 这里会有一个循环，不正确就会确认一次。。
			v.lock.Lock()
			v.keys[uid] = reply.Tk.Key
			v.lock.Unlock()
			log.Error("Token Invalid")
			return ecode.AccessDenied
		}
		// unreachable!!!
	}
}

func (v *Verify) Verify(ctx *bm.Context) {
	if err := v.verify(ctx); err != nil {
		ctx.JSON(nil, err)
		ctx.Abort()
	}
}

func (v *Verify) GenToken(c context.Context, id int64) (token string, err error) {
	var (
		req *TokenReq
		rpl *TokenReply
	)
	req = &TokenReq{
		Tk:                   &Token{
			Id:                   id,
			Key:                  "",
		},
	}
	rpl, err = v.client.GenKey(c, req)
	if err != nil {
		err = ecode.ServerErr
		return
	}
	token = rpl.Tk.Key
	// update key cache
	v.lock.Lock()
	v.keys[id] = token
	v.lock.Unlock()
	return
}

// NewClient new grpc client
func newRPCVerifyClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (VerifyClient, error) {
	client := warden.NewClient(cfg, opts...)
	cc, err := client.Dial(context.Background(), fmt.Sprintf("discovery://default/%s", AppID))
	if err != nil {
		return nil, err
	}
	return NewVerifyClient(cc), nil
}



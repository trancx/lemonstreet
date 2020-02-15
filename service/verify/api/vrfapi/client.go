package vrfapi

import (
	"context"
	"fmt"
	"github.com/bilibili/kratos/pkg/ecode"
	"net/http"
	"strings"
	"sync"

	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"
	"google.golang.org/grpc"
)

// AppID .
const AppID = "verify.service"

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
// FIXME: LOG!!!!
func (v *Verify) verify(ctx *bm.Context) error {
	var (
		uid     int64
		req     *TokenReq
		err     error
		cookies *http.Cookie
	)

	if cookies, err = ctx.Request.Cookie("uid"); err != nil {
		err = ecode.AccessDenied
		return err
	}
	if _, err = fmt.Sscanf(cookies.Value, "%d", &uid); err != nil {
		return ecode.AccessDenied
	}
	if cookies, err = ctx.Request.Cookie("token"); err != nil {
		err = ecode.AccessDenied
		return err
	}
	req = &TokenReq{
		Tk: &Token{
			Id:  uid,
			Key: cookies.Value,
		},
	}
	v.lock.RLock()
	token, ok := v.keys[uid]
	v.lock.RUnlock()

	if !ok {
		reply, err := v.client.VrfKey(ctx, req)
		if err != nil {
			// FIXME: RPC error
			return ecode.AccessDenied
		}
		if reply.IsValid {
			// means token is right one and we can cached it!
			v.lock.Lock()
			v.keys[uid] = cookies.Value // RPC will return nil, when invalid, see gRPC handler
			v.lock.Unlock()
			return nil
		} else {
			// if reply.IsUpdated, not the right one, but we can cached it
			v.lock.Lock()
			v.keys[uid] = reply.Tk.Key
			v.lock.Unlock()
			return ecode.AccessDenied
		}
		// unreachable!!!
	}
	// cached hit
	if !strings.EqualFold(cookies.Value, token) {
		return ecode.AccessDenied
	}
	return nil
}

func (v *Verify) Verify(ctx *bm.Context) {
	if err := v.verify(ctx); err != nil {
		ctx.JSON(nil, err)
		ctx.Abort()
	}
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

// 生成 gRPC 代码
//go:generate kratos tool protoc --grpc --bm api.proto

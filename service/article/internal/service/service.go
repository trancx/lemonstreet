package service

import (
	acc "article/api/accapi"
	artapi "article/api/artapi"
	cmt "article/api/cmtapi"
	"article/internal/dao"
	"article/internal/model"
	"context"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/ecode"
	"github.com/bilibili/kratos/pkg/naming/discovery"
	"github.com/bilibili/kratos/pkg/net/rpc/warden/resolver"
	"github.com/prometheus/common/log"
)

// Service service.
type Service struct {
	ac     *paladin.Map
	dao    *dao.Dao
	accRPC acc.AccountClient
	cmtRPC cmt.CommentsClient
}

/* gRPC implementation, normal path, db -> cache, or specific amount(e.g. 100) comments
	cache together? but we have both uid-key and aid-key, cache is foolish here.
*/
func (s *Service) LatestArticles(ctx context.Context, in *artapi.TimeReq) (*artapi.ArticleBaseInfosReply, error) {
	var (
		err error
		infos []*artapi.ArticleBaseInfo
		res		[]artapi.ArticleBaseInfo
	)
	res, err = s.dao.RawArticleBaseInfoByDate(ctx, in.Beg, in.End)
	if err != nil {
		err = ecode.NothingFound
		return nil, err
	}
	// Ugly code!!!
	for _, temp := range res {
		infos = append(infos, &temp)
	}
	return &artapi.ArticleBaseInfosReply{
		Infos:                infos,
	}, err
}

// NEED_TESTED Cached
func (s *Service) SearchArticlesByUID(c context.Context, req *artapi.IDReq) (*artapi.ArticleBaseInfosReply, error) {
	var (
		err error
		uid int64
		infos []*artapi.ArticleBaseInfo
		res		[]artapi.ArticleBaseInfo
	)
	uid = req.Id

	res, err = s.dao.ArticleBaseInfosByUId(c, uid)
	if err != nil {
		err = ecode.NothingFound
		return nil, err
	}
	// Ugly code!!!
	for _, temp := range res {
		infos = append(infos, &temp)
	}

	return &artapi.ArticleBaseInfosReply{
		Infos:                infos,
	}, err
}

// NEED_TESTED
func (s *Service) SearchArticlesByTitle(c context.Context, req *artapi.NameReq) (*artapi.ArticleBaseInfosReply, error) {
	var (
		err error
		title string
		infos []*artapi.ArticleBaseInfo
		res		[]artapi.ArticleBaseInfo
	)
	title = req.Name
	res, err = s.dao.ArticleBaseInfosByTitle(c, title)
	if err != nil {
		err = ecode.NothingFound
		return nil, err
	}
	// Ugly code!!!
	for _, temp := range res {
		infos = append(infos, &temp)
	}

	return &artapi.ArticleBaseInfosReply{
		Infos:                infos,
	}, err
}

// username && title
// content and userinfo inside http body!
func (s *Service) PostArticle(c context.Context, info *artapi.ArticleBaseInfo, content string) (err error) {
	reply, err := s.accRPC.BaseInfo(c, &acc.UidReq{
		Uid:                  info.Uid,
		RealIp:               "",
	})
	if err != nil {
		log.Errorf("Account RPC error (%v)", err)
		return
	}
	info.Author = reply.Info.Name
	// save to DB and cache it!
	err = s.dao.PostArticle(c, info, content)
	// handle error
	return
}

// generate comments
func (s *Service) GetArticleAnnms(c context.Context, aid int64) (*model.ArticleInfo, error) {
	var (
		reply *acc.BaseInfoReply
		abi	*artapi.ArticleBaseInfo
		art *model.Article
		info *acc.UserBaseInfo
		cmtRpl *cmt.CommentsReply = nil
		cmts []*cmt.Comment
		err	error
	)

	abi, err = s.dao.ArticleBaseInfoByAId(c, aid) // cache it!
	if err != nil {
		err = ecode.NothingFound
		return nil, err
	}

	req := acc.UidReq{
		Uid:				  abi.Uid,
		RealIp:               "",
	}
	reply, err = s.accRPC.BaseInfo(c, &req)
	if err != nil {
		log.Errorf("wrong uid-%d in aid-%d", abi.Uid, abi.Aid)
		err = ecode.ServerErr
		return nil, err
	}
	info = reply.Info

	art, err = s.dao.Article(c, abi.Aid)
	if err != nil {
		err = ecode.ServerErr
		return nil, err
	}

	cmtRpl, _ = s.cmtRPC.CommentsOfAID(c, &cmt.IDReq{ Id: abi.Aid})
	// FIXME: ignore err, when not found means no comments
	// but inner RPC error should not be ignore
	cmts = cmtRpl.Comments
	return &model.ArticleInfo{
		UInfo:   info,
		ABI:     abi,
		Content: art,
		Comments: cmts,
	}, err
}

func (s *Service) GetArticle(c context.Context, user *acc.UserBaseInfo, article *model.ArticleBaseInfo) (info *model.ArticleInfo, err error) {
	var (

	)

	// generate article_id and fill the article baseinfo

	// save to DB and cache it!

	return
}

func init(){
	// NOTE: 注意这段代码，表示要使用discovery进行服务发现
	// NOTE: 还需注意的是，resolver.Register是全局生效的，所以建议该代码放在进程初始化的时候执行
	// NOTE: ！！！切记不要在一个进程内进行多个不同中间件的Register！！！
	// NOTE: 在启动应用时，可以通过flag(-discovery.nodes) 或者 环境配置(DISCOVERY_NODES)指定discovery节点
	resolver.Register(discovery.Builder())
}


// New new a service and return.
func New(d *dao.Dao) (s *Service, err error) {
	var (
		accRPC acc.AccountClient
		cmtRPC cmt.CommentsClient
	)

	if accRPC, err = acc.NewRPCAccountClient(nil); err != nil {
		panic(err)
	}
	cmtRPC, err = cmt.NewRPCCommentClient(nil)
	s = &Service{
		ac:     &paladin.TOML{},
		dao:    d,
		accRPC: accRPC,
		cmtRPC: cmtRPC,
	}
	err = paladin.Watch("application.toml", s.ac)
	return
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}

package artapi
import (
	"context"
	"fmt"

	"github.com/bilibili/kratos/pkg/net/rpc/warden"

	"google.golang.org/grpc"
)

// AppID .
const ArticleAppID = "article.service"

// NewClient new grpc client
func NewRPCArticleClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (ArticleClient, error) {
	client := warden.NewClient(cfg, opts...)
	cc, err := client.Dial(context.Background(), fmt.Sprintf("discovery://default/%s", ArticleAppID))
	if err != nil {
		return nil, err
	}
	return NewArticleClient(cc), nil
}

// 生成 gRPC 代码
//go:generate kratos tool protoc --grpc  api.proto

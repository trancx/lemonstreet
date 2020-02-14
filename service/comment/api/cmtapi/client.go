package cmtapi
import (
	"context"
	"fmt"

	"github.com/bilibili/kratos/pkg/net/rpc/warden"

	"google.golang.org/grpc"
)

// AppID .
const AppID = "comment.service"

// NewClient new grpc client NewRPCArticleClient
func NewRPCCommentClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (CommentsClient, error) {
	client := warden.NewClient(cfg, opts...)
	cc, err := client.Dial(context.Background(), fmt.Sprintf("discovery://default/%s", AppID))
	if err != nil {
		return nil, err
	}
	return NewCommentsClient(cc), nil
}

// 生成 gRPC 代码
//go:generate kratos tool protoc --grpc api.proto

package accapi

import (
	"context"
	"fmt"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"

	"google.golang.org/grpc"
)

// AppID .
const AccountAppID = "account.service"

// NewClient new grpc client
func NewRPCAccountClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (AccountClient, error) {
	client := warden.NewClient(cfg, opts...)
	cc, err := client.Dial(context.Background(), fmt.Sprintf("discovery://default/%s", AccountAppID))
	if err != nil {
		return nil, err
	}
	return NewAccountClient(cc), nil
}

// 生成 gRPC 代码  dont use !!! --bm
//go:generate kratos tool protoc --grpc api.proto
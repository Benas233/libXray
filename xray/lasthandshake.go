
package xray

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

func QueryLasthandshake(server string) (string, error) {

	conn, err := grpc.Dial(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()
	req := &emptypb.Empty{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp := &LastHandshakeResponse{}
	method := "/xray.app.lasthandshake.command.LasthandshakeService/GetLastHandshake"
	if err := conn.Invoke(ctx, method, req, resp); err != nil {
		return "", fmt.Errorf("RPC call failed: %w", err)
	}

	jsonResp, err := protojson.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal respone: %w", err)
	}

	return string(jsonResp), nil
}

package xray

import (
	context "context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func QueryLasthandshake(serverBase64 string) string {
	serverBytes, err := base64.StdEncoding.DecodeString(serverBase64)
	if err != nil {
		return encodeResponse("", fmt.Errorf("failed to decode server address: %w", err))
	}
	server := string(serverBytes)

	conn, err := grpc.Dial(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return encodeResponse("", fmt.Errorf("failed to connect: %w", err))
	}
	defer conn.Close()

	req := &emptypb.Empty{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp := &LastHandshakeResponse{}
	method := "/xray.app.lasthandshake.command.LasthandshakeService/GetLastHandshake"
	if err := conn.Invoke(ctx, method, req, resp); err != nil {
		return encodeResponse("", fmt.Errorf("RPC call failed: %w", err))
	}

	jsonResp, err := protojson.Marshal(resp)
	if err != nil {
		return encodeResponse("", fmt.Errorf("failed to marshal response: %w", err))
	}

	return encodeResponse(string(jsonResp), nil)
}

func encodeResponse(data string, err error) string {
	type wrapper struct {
		Data  string `json:"data"`
		Error string `json:"error,omitempty"`
	}

	w := wrapper{
		Data: data,
	}
	if err != nil {
		w.Error = err.Error()
	}

	jsonBytes, _ := json.Marshal(w)
	return base64.StdEncoding.EncodeToString(jsonBytes)
}

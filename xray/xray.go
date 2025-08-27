package xray

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/xtls/libxray/memory"
	"github.com/xtls/xray-core/common/cmdarg"
	"github.com/xtls/xray-core/core"
	_ "github.com/xtls/xray-core/main/distro/all"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Constants for environment variables
const (
	coreAsset = "xray.location.asset"
	coreCert  = "xray.location.cert"
)

var (
	coreServer *core.Instance
)

func StartXray(configPath string) (*core.Instance, error) {
	file := cmdarg.Arg{configPath}
	config, err := core.LoadConfig("json", file)
	if err != nil {
		return nil, err
	}

	server, err := core.New(config)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func StartXrayFromJSON(configJSON string) (*core.Instance, error) {
	// Convert JSON string to bytes
	configBytes := []byte(configJSON)

	// Use core.StartInstance which can load configuration directly from bytes
	server, err := core.StartInstance("json", configBytes)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func InitEnv(datDir string) {
	os.Setenv(coreAsset, datDir)
	os.Setenv(coreCert, datDir)
}

// Run Xray instance.
// datDir means the dir which geosite.dat and geoip.dat are in.
// configPath means the config.json file path.
func RunXray(datDir string, configPath string) (err error) {
	InitEnv(datDir)
	memory.InitForceFree()
	coreServer, err = StartXray(configPath)
	if err != nil {
		return
	}

	if err = coreServer.Start(); err != nil {
		return
	}

	debug.FreeOSMemory()
	return nil
}

// Run Xray instance with JSON configuration string.
// datDir means the dir which geosite.dat and geoip.dat are in.
// configJSON means the JSON configuration string.
func RunXrayFromJSON(datDir string, configJSON string) (err error) {
	InitEnv(datDir)
	memory.InitForceFree()
	coreServer, err = StartXrayFromJSON(configJSON)
	if err != nil {
		return
	}

	if err = coreServer.Start(); err != nil {
		return
	}

	debug.FreeOSMemory()
	return nil
}

// Get Xray State
func GetXrayState() bool {
	return coreServer.IsRunning()
}

// Stop Xray instance.
func StopXray() error {
	if coreServer != nil {
		err := coreServer.Close()
		coreServer = nil
		if err != nil {
			return err
		}
	}
	return nil
}

// Xray's version
func XrayVersion() string {
	return core.Version()
}

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

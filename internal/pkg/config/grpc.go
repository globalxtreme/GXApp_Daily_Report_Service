package config

import (
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"os"
	"time"
)

var (
	TestingRPC string
)

func InitRPC() {
	xtremepkg.RPCDialTimeout = 5 * time.Second
	TestingRPC = os.Getenv("GRPC_TESTING_HOST")
}

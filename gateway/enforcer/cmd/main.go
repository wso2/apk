package main

import(
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/xds"
)

func main() {
	cfg := config.GetConfig()
	
	xds.CreateXDSClients(cfg)

	// Wait forever
	select {}
}
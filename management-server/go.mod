module APKManagementServer

go 1.18

// todo(amaliMatharaarachchi) remove replace once adapter side get merged
replace github.com/wso2/product-microgateway/adapter => github.com/AmaliMatharaarachchi/product-microgateway/adapter v0.0.0-20221027040248-0052cbae398e

require (
	github.com/envoyproxy/go-control-plane v0.10.2-0.20211124143408-6141aee35516
	github.com/sirupsen/logrus v1.9.0
	github.com/wso2/product-microgateway/adapter v0.0.0
	google.golang.org/grpc v1.36.0
)

require (
	github.com/census-instrumentation/opencensus-proto v0.2.1 // indirect
	github.com/cncf/xds/go v0.0.0-20211001041855-01bcc9b48dfe // indirect
	github.com/envoyproxy/protoc-gen-validate v0.4.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	golang.org/x/net v0.0.0-20220909164309-bea034e7d591 // indirect
	golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20201201144952-b05cb90ed32e // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

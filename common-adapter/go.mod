module github.com/wso2/apk/common-adapter

go 1.19

require (
	github.com/envoyproxy/go-control-plane v0.11.0
	github.com/pelletier/go-toml v1.8.1
	github.com/sirupsen/logrus v1.9.0
	github.com/wso2/apk/adapter v0.0.0-20230628072301-3ed4c08eae55
	google.golang.org/grpc v1.52.0
)

replace github.com/wso2/apk/adapter => ../adapter

require (
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cncf/xds/go v0.0.0-20220314180256-7f1daf1720fc // indirect
	github.com/envoyproxy/protoc-gen-validate v0.9.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

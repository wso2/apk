module github.com/wso2/apk/management-server

go 1.19

require (
	github.com/envoyproxy/go-control-plane v0.11.2-0.20230802074621-eea0b3bd0f81
	github.com/pelletier/go-toml v1.9.5
	github.com/sirupsen/logrus v1.9.0
	github.com/wso2/apk/adapter v0.0.0-20231214082511-af2c8b8a19f1
	google.golang.org/grpc v1.58.3
	google.golang.org/protobuf v1.31.0
)

replace github.com/wso2/apk/adapter => ../adapter

require (
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cncf/xds/go v0.0.0-20230607035331-e9ce68804cb4 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto v0.0.0-20230731193218-e0aa005b6bdf // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230726155614-23370e0ffb3e // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230731190214-cbb8c96f2d6d // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

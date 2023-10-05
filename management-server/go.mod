module github.com/wso2/apk/management-server

go 1.19

require (
	github.com/envoyproxy/go-control-plane v0.11.2-0.20230802074621-eea0b3bd0f81
	github.com/jackc/pgx/v5 v5.3.1
	github.com/pelletier/go-toml v1.9.5
	github.com/sirupsen/logrus v1.9.0
	github.com/wso2/apk/adapter v0.0.0-20230313062104-25216c8acbc5
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.31.0
)

replace github.com/wso2/apk/adapter => ../adapter

require (
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cncf/xds/go v0.0.0-20230607035331-e9ce68804cb4 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.0 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230731193218-e0aa005b6bdf // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230726155614-23370e0ffb3e // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230731190214-cbb8c96f2d6d // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

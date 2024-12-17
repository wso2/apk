package ext_proc

import (
	"io"
	"fmt"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"

  "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExternalProcessingServer struct {
	log logging.Logger
}

func StartExternalProcessingServer(cfg *config.Server) {
	server, err := util.CreateGRPCServer(cfg.EnforcerPublicKeyPath, cfg.EnforcerPrivateKeyPath)
	if err != nil {
		panic(err)
	}
	envoy_service_proc_v3.RegisterExternalProcessorServer(server, &ExternalProcessingServer{cfg.Logger})
}

func (s *ExternalProcessingServer) Process(srv envoy_service_proc_v3.ExternalProcessor_ProcessServer) error {
	ctx := srv.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive stream request: %v", err)
		}

		resp := &envoy_service_proc_v3.ProcessingResponse{}
		switch v := req.Request.(type) {
		case *envoy_service_proc_v3.ProcessingRequest_RequestHeaders:
			if v.RequestHeaders != nil {
				hdrs := v.RequestHeaders.Headers.GetHeaders()
				for _, hdr := range hdrs {
					s.log.Info(fmt.Sprintf("Header: %v\n", hdr))
				}
			}

			rhq := &envoy_service_proc_v3.HeadersResponse{
				Response: &envoy_service_proc_v3.CommonResponse{},
			}
			resp = &envoy_service_proc_v3.ProcessingResponse{
				Response: &envoy_service_proc_v3.ProcessingResponse_RequestHeaders{
					RequestHeaders: rhq,
				},
			}
			break
		case *envoy_service_proc_v3.ProcessingRequest_ResponseHeaders:
			s.log.Info(fmt.Sprintf("response header"))
			rhq := &envoy_service_proc_v3.HeadersResponse{
				Response: &envoy_service_proc_v3.CommonResponse{
				},
			}
			resp = &envoy_service_proc_v3.ProcessingResponse{
				Response: &envoy_service_proc_v3.ProcessingResponse_ResponseHeaders{
					ResponseHeaders: rhq,
				},
			}
			break
		case *envoy_service_proc_v3.ProcessingRequest_ResponseBody:
			httpBody := req.GetResponseBody()
			s.log.Info(fmt.Sprintf("request body %v\n", httpBody))
			rbq := &envoy_service_proc_v3.BodyResponse{
				Response: &envoy_service_proc_v3.CommonResponse{
				},
			}
			resp = &envoy_service_proc_v3.ProcessingResponse{
					Response: &envoy_service_proc_v3.ProcessingResponse_ResponseBody{
							ResponseBody: rbq,
					},
			}
			
		case *envoy_service_proc_v3.ProcessingRequest_RequestBody:
			httpBody := req.GetRequestBody()
			s.log.Info(fmt.Sprintf("request body %v\n", httpBody))
			rbq := &envoy_service_proc_v3.BodyResponse{
				Response: &envoy_service_proc_v3.CommonResponse{
				},
			}
			resp = &envoy_service_proc_v3.ProcessingResponse{
					Response: &envoy_service_proc_v3.ProcessingResponse_RequestBody{
							RequestBody: rbq,
					},
			}
		default:
			s.log.Info(fmt.Sprintf("Unknown Request type %v\n", v))
		}
		if err := srv.Send(resp); err != nil {
			s.log.Info(fmt.Sprintf("send error %v", err))
		}
	}
}

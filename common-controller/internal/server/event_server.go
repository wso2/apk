package server

import (
	"github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/utils"
	apkmgt "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"google.golang.org/grpc/metadata"
)

// EventServer struct use to hold event server
type EventServer struct {
	apkmgt.UnimplementedEventStreamServiceServer
}

// StreamEvents streams events to the enforcer
func (s EventServer) StreamEvents(req *apkmgt.Request, srv apkmgt.EventStreamService_StreamEventsServer) error {
	// Read metadata from the request context
	md, ok := metadata.FromIncomingContext(srv.Context())
	if !ok {
		loggers.LoggerAPKOperator.Errorf("error : %v", "Failed to get metadata from the request context")
		return nil
		// Handle the case where metadata is not present
	}
	enforcerID := md.Get("enforcer-uuid")
	loggers.LoggerAPKOperator.Debugf("Enforcer ID : %v", enforcerID[0])
	utils.AddClientConnection(enforcerID[0], srv)
	utils.SendResetEvent()
	<-srv.Context().Done()
	loggers.LoggerAPKOperator.Infof("Connection closed by the client : %v", enforcerID[0])
	utils.DeleteClientConnection(enforcerID[0])
	return nil // Client closed the connection
}

package utils

import (
	apkmgt "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/service/apkmgt"
)

var clientConnections = make(map[string]apkmgt.EventStreamService_StreamEventsServer)

// AddClientConnection adds a client connection to the map
func AddClientConnection(clientID string, stream apkmgt.EventStreamService_StreamEventsServer) {
	clientConnections[clientID] = stream
}

// DeleteClientConnection deletes a client connection from the map
func DeleteClientConnection(clientID string) {
	delete(clientConnections, clientID)
}

// GetAllClientConnections returns all client connections
func GetAllClientConnections() map[string]apkmgt.EventStreamService_StreamEventsServer {
	return clientConnections
}

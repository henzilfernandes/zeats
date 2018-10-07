package wire

import (
	"zeats/types"
	"zeats/config"
	"zeats/cassandra"
	"zeats/gateway"
)

var (
	component types.Component
)

func Start() {

	// Cassandra
	cSession := cassandra.NewCassandraSession(config.Vals().CassandraNodes)

	// Gateway...
	gatewayHandle := gateway.NewGateway(cSession)
	go gatewayHandle.Start()
	component = gatewayHandle
}

func Stop() {
	component.Stop()
}


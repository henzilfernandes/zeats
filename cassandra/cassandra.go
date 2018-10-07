package cassandra

import (
	"zeats/config"
	"github.com/gocql/gocql"
	"github.com/golang/glog"
	"fmt"
)

func NewCassandraSession(cassandraNodes []*config.Node) *gocql.Session {
	cluster := gocql.NewCluster(getCassandraHosts(cassandraNodes)...)
	cSession, err := cluster.CreateSession()
	if err != nil {
		glog.Fatalf("Error while connecting cassandra $s\n", err)
	}
	fmt.Println("[✔]\tCassandra connected")
	return cSession
}

func getCassandraHosts(cassandraHosts []*config.Node) []string {
	var hosts []string
	for _, cassandraHost := range cassandraHosts {
		hosts = append(hosts, cassandraHost.Host)
	}
	return hosts
}

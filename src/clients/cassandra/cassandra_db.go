package cassandra

import (
	"fmt"
	"github.com/gocql/gocql"
)

var(
	cluster *gocql.ClusterConfig
)

func init()  {
	// Connect to Cassandra cluster:
	cluster = gocql.NewCluster("192.168.56.4")
	fmt.Println(cluster)
	cluster.Keyspace = "oauth"
	cluster.Consistency = gocql.Quorum
}

func GetSession() (*gocql.Session, error) {
	return cluster.CreateSession()
}
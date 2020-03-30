package cassandra

import (
	"fmt"
	"github.com/gocql/gocql"
)

var (
	session *gocql.Session
)

func init() {
	// Connect to Cassandra cluster:
	cluster := gocql.NewCluster("192.168.56.4")
	fmt.Println(cluster)
	cluster.Keyspace = "oauth"
	cluster.Consistency = gocql.Quorum

	var err error
	if session, err = cluster.CreateSession(); err != nil {
		panic(err)
	}
}

func GetSession() *gocql.Session {
	return session
}

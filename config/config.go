package config

// Node ..
type Node struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Config struct {
	Port string `yaml:"port"`
	CassandraNodes []*Node `yaml:"cassandra_nodes"`
}



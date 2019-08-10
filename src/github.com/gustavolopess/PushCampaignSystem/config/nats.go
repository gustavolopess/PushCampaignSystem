package config

// struct with NATS configurations
type NatsConfig struct {
	Host		string	`json:"host"`
	Port		int		`json:"port"`
	ClusterID	string	`json:"cluster_id"`
	Subject		string	`json:"subject"`
}

func (c *NatsConfig) LoadConfig(configPath string) {
	return
}


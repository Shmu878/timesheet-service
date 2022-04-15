package config

import (
	"os"
	"path/filepath"
)

// Here are defined all types for your configuration
// You can remove not needed types or add your own

type Storages struct {
	Redis    *redis.Config
	Database *db.DbClusterConfig
}

type Adapter struct {
	Grpc *grpc.ClientConfig
}

type Config struct {
	Grpc       *grpc.ServerConfig
	Storages   *Storages
	Nats       *queue.Config
	Log        *log.Config
	Cluster    *service.Config
	Adapters   map[string]*Adapter
	Monitoring *monitoring.Config
}

func Load() (*Config, error) {

	// get root folder from env
	rootPath := os.Getenv("FOCROOT")
	if rootPath == "" {
		return nil, kitConfig.ErrConfigPaErrConfigPathIsEmpty()
	}

	// config path
	configPath := filepath.Join(rootPath, meta.Meta.ServiceCode(), "config.yml")

	// load config
	config := &Config{}
	err := kitConfig.NewConfigLoader(logger.LF()).
		WithConfigPath(configPath).
		Load(config)

	if err != nil {
		return nil, err
	}
	return config, nil
}

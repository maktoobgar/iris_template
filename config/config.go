package config

import db "service/pkg/database"

type (
	Config struct {
		CurrentMicroservice   Microservice
		Translator            Translator              `yaml:"translator"`
		Logging               Logging                 `yaml:"logging"`
		Gateway               Microservice            `yaml:"gateway"`
		Microservices         map[string]Microservice `yaml:"microservices"`
		Debug                 bool                    `yaml:"debug"`
		Domain                string                  `yaml:"domain"`
		PWD                   string                  `yaml:"pwd"`
		AllowOrigins          string                  `yaml:"allow_origins"`
		AllowHeaders          string                  `yaml:"allow_headers"`
		MaxAge                int                     `yaml:"max_age"`
		Timeout               int64                   `yaml:"timeout"`
		MaxConcurrentRequests int                     `yaml:"max_concurrent_requests"`
		SecretKey             string                  `yaml:"secret_key"`
	}

	Translator struct {
		Path string `yaml:"path"`
	}

	Logging struct {
		Path         string `yaml:"path"`
		Pattern      string `yaml:"pattern"`
		MaxAge       string `yaml:"max_age"`
		RotationTime string `yaml:"rotation_time"`
		RotationSize string `yaml:"rotation_size"`
	}

	Microservice struct {
		Databases        map[string]db.Database `yaml:"databases"`
		MongodbDatabases map[string]db.Database `yaml:"mongodb_databases"`
		IP               string                 `yaml:"ip"`
		Port             string                 `yaml:"port"`
	}
)

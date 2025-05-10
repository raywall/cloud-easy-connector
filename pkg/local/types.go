package local

import (
	lc "github.com/raywall/cloud-easy-connector/internal/local"
)

func New() LocalResource {
	return &lc.LocalResource{}
}

type LocalResource interface {
	GetEnvOrDefault(key, defaultValue string) string
}

package local

import "os"

type LocalResource struct{}

// GetEnvOrDefault recupera uma variável de ambiente ou retorna um valor padrão caso a variável não exista
func (l *LocalResource) GetEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

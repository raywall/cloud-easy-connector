package auth

// AuthRequest representa os dados enviados para a API de autenticação
type AuthRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

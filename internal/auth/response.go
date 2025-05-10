package auth

// TokenResponse representa a resposta da API de autenticação
type TokenResponse struct {
	Token        string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresAt    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	Active       bool   `json:"active"`
}

package oauth

type TokenExchangeRequest struct {
	Code         string `json:"code" binding:"required"`
	RedirectUri  string `json:"redirect_uri" binding:"required"`
	ClientId     string `json:"client_id" binding:"required"`
	CodeVerifier string `json:"code_verifier" binding:"required"`
}

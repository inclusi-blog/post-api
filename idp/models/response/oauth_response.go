package response

import "post-api/idp/models/db"

type AcceptResponse struct {
	RedirectTo string `json:"redirect_to" example:"some redirect url"`
}

type GenericAuthError struct {
	Error            string `json:"error"`
	StatusCode       int    `json:"status_code"`
	ErrorDescription string `json:"error_description"`
	Debug            string `json:"debug"`
}

type ConsentResponse struct {
	Userprofile                  db.UserProfile `json:"Context"`
	RequestedAccessTokenAudience []string       `json:"requested_access_token_audience"`
	RequestedScope               []string       `json:"requested_scope"`
}

type TokenExchangeResponse struct {
	AccessToken      string `json:"access_token"`
	IdToken          string `json:"id_token,omitempty"`
	EncryptedIdToken string `json:"enc_id_token"`
	ExpiresIn        int    `json:"expires_in"`
	ExpiresAt        string `json:"expires_at"`
}

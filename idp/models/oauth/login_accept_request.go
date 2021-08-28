package oauth

import "post-api/idp/models/db"

type LoginAcceptRequest struct {
	Subject     string         `json:"subject"`
	Remember    bool           `json:"remember"`
	RememberFor int            `json:"remember_for"`
	Acr         string         `json:"acr"`
	Userprofile db.UserProfile `json:"Context"`
}

type ConsentAcceptRequest struct {
	GrantAccessTokenAudience []string `json:"grant_access_token_audience"`
	GrantScope               []string `json:"grant_scope"`
	Remember                 bool     `json:"remember"`
	RememberFor              int      `json:"remember_for"`
	Session                  Session  `json:"session"`
}

type Session struct {
	Userprofile db.UserProfile `json:"id_token"`
}

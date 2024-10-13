package oauthlogin

import (
	"net/http"
	"net/url"

	"userstyles.world/modules/config"
	"userstyles.world/modules/util"
)

// This is just an empty stub.
// However we will link all necessary functions to this stub.
type github struct{}

func (github) oauthMakeURL() string {
	// Base URL.
	oauthURL := "https://github.com/login/oauth/authorize"
	// Add our app client ID.
	oauthURL += "?client_id=" + config.OpenAuth.GitHubID
	// Add email scope.
	oauthURL += "&scope=" + url.QueryEscape("user:email")

	return oauthURL
}

func (github) enableState() bool {
	return true
}

func (github) appendToRedirect(state any) string {
	// Trying to follow our stateless design we encrypt the
	// Nonsense state so we later can re-use by decrypting it.
	// And than have the actual value. Also we use this to specify
	// From which site the callback was from.
	return util.EncryptText(state.(string), util.AEADOAuth, config.Secrets) + "/"
}

func (github) getAuthTokenURL(state any) string {
	authURL := "https://github.com/login/oauth/access_token"
	authURL += "?client_id=" + config.OpenAuth.GitHubID
	authURL += "&client_secret=" + config.OpenAuth.GitHubSecret
	// Add the nonsense state we uses earlier.
	authURL += "&state=" + state.(string)

	return authURL
}

func (github) isAuthTokenPost() bool {
	return false
}

func (github) getAuthTokenPostBody(any) authURLPostBody {
	return authURLPostBody{}
}

func (github) beforeRequest(authURLPostBody, *http.Request) error {
	return nil
}

func (github) getUserEndpoint() string {
	return "https://api.github.com/user"
}

func (github) getServiceType() Service {
	return GithubService
}

func (github) getEmailEndpoint() string {
	return "https://api.github.com/user/emails"
}

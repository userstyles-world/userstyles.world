package oauthlogin

import (
	"net/http"
	"net/url"

	"userstyles.world/modules/config"
	"userstyles.world/utils"
)

// This is just an empty stub.
// However we will link all neccessary functions to this stub.
type github struct {
}

func (github) oauthMakeURL() string {
	// Base URL.
	oauthURL := "https://github.com/login/oauth/authorize"
	// Add our app client ID.
	oauthURL += "?client_id=" + config.GITHUB_CLIENT_ID
	// Add email scope.
	oauthURL += "&scope=" + url.QueryEscape("user:email")

	return oauthURL
}

func (github) enableState() bool {
	return true
}

func (github) appendToRedirect(state interface{}) string {
	// Trying to follow our stateless design we encrypt the
	// Nonsense state so we later can re-use by decrypting it.
	// And than have the actual value. Also we use this to specify
	// From which site the callback was from.
	return utils.EncryptText(state.(string), utils.AEAD_OAUTH, config.ScrambleConfig) + "/"
}

func (github) getAuthTokenURL(state interface{}) string {
	authURL := "https://github.com/login/oauth/access_token"
	authURL += "?client_id=" + config.GITHUB_CLIENT_ID
	authURL += "&client_secret=" + config.GITHUB_CLIENT_SECRET
	// Add the nonsense state we uses earlier.
	authURL += "&state=" + state.(string)

	return authURL
}

func (github) isAuthTokenPost() bool {
	return false
}

func (github) getAuthTokenPostBody(data interface{}) authURLPostBody {
	return authURLPostBody{}
}

func (github) beforeRequest(body authURLPostBody, req *http.Request) error {
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

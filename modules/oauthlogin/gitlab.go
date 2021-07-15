package oauthlogin

import (
	"net/http"
	"net/url"

	"userstyles.world/modules/config"
)

// This is just an empty stub.
// However we will link all necessary functions to this stub.
type gitlab struct {
}

const gitlabStr = "gitlab"

func (gitlab) oauthMakeURL() string {
	oauthURL := "https://gitlab.com/oauth/authorize"
	// Add our app client ID.
	oauthURL += "?client_id=" + config.GITLAB_CLIENT_ID
	// Define we want a code back
	oauthURL += "&response_type=code"
	// Add read_user scope.
	oauthURL += "&scope=read_user"
	return oauthURL
}

func (gitlab) enableState() bool {
	return false
}

func (gitlab) appendToRedirect(data interface{}) string {
	return gitlabStr + "/"
}

func (gitlab) getAuthTokenURL(data interface{}) string {
	authURL := "https://gitlab.com/oauth/token"
	authURL += "?client_id=" + config.GITLAB_CLIENT_ID
	authURL += "&client_secret=" + config.GITLAB_CLIENT_SECRET
	// Define we log in trough the temp code.
	authURL += "&grant_type=authorization_code"
	// Specify the the redirect uri, because it is required
	authURL += "&redirect_uri=" + url.PathEscape(config.OAuthURL()+gitlabStr+"/")

	return authURL
}

func (gitlab) isAuthTokenPost() bool {
	return false
}

func (gitlab) getAuthTokenPostBody(data interface{}) authURLPostBody {
	return authURLPostBody{}
}

func (gitlab) beforeRequest(body authURLPostBody, req *http.Request) error {
	return nil
}

func (gitlab) getUserEndpoint() string {
	return "https://gitlab.com/api/v4/user"
}

func (gitlab) getServiceType() Service {
	return GitlabService
}

// We don't return anything, because when it's used, gitlab will never reach that function.
func (gitlab) getEmailEndpoint() string {
	return ""
}

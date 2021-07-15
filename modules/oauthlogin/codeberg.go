package oauthlogin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"userstyles.world/modules/config"
	"userstyles.world/utils"
)

// This is just an empty stub.
// However we will link all necessary functions to this stub.
type codeberg struct {
}

const codebergStr = "codeberg"

func (codeberg) oauthMakeURL() string {
	// Base URL.
	oauthURL := "https://codeberg.org/login/oauth/authorize"
	// Add our app client ID.
	oauthURL += "?client_id=" + config.CODEBERG_CLIENT_ID
	// Define we want a code back
	oauthURL += "&response_type=code"

	return oauthURL
}

func (codeberg) enableState() bool {
	return false
}

func (codeberg) appendToRedirect(data interface{}) string {
	return codebergStr + "/"
}

func (codeberg) getAuthTokenURL(data interface{}) string {
	return "https://codeberg.org/login/oauth/access_token"
}

func (codeberg) isAuthTokenPost() bool {
	return true
}

func (codeberg) getAuthTokenPostBody(data interface{}) authURLPostBody {
	return authURLPostBody{
		ClientID:     config.CODEBERG_CLIENT_ID,
		ClientSecret: config.CODEBERG_CLIENT_SECRET,
		Code:         data.(string),
		GrantType:    "authorization_code",
		RedirectURI:  config.OAuthURL() + "codeberg/",
	}
}

func (codeberg) beforeRequest(body authURLPostBody, req *http.Request) error {
	bodyString, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(strings.NewReader(utils.UnsafeString(bodyString)))
	req.ContentLength = int64(len(bodyString))
	req.Header.Set("Content-Type", "application/json")
	return nil
}

func (codeberg) getUserEndpoint() string {
	return "https://codeberg.org/api/v1/user"
}

func (codeberg) getServiceType() Service {
	return CodebergService
}

func (codeberg) getEmailEndpoint() string {
	return "https://codeberg.org/api/v1/user/emails"
}

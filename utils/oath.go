package utils

import (
	"encoding/json"
	"net/http"
	"net/url"

	"userstyles.world/config"
)

var (
	githubClientID = "f3f2d378d88794062895"
	gitlabClientID = "f906330603e645d54184fd13337ca51b762ea5da2188e7f248e109c686940897"
)

// https://docs.github.com/en/developers/apps/authorizing-oauth-apps#response
type OAuthTokenResponse struct {
	AccesToken string `json:"access_token"`
	TokenType  string `json:"token_type"`
}

type OAuthResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func OauthMakeURL(baseURL, service string) string {
	if service == "" {
		return ""
	}
	var oauthURL string
	nonsenseState := B2s(RandStringBytesMaskImprSrcUnsafe(16))
	if service == "github" {
		// Base URL.
		oauthURL = "https://github.com/login/oauth/authorize"
		// Add our app client ID.
		oauthURL += "?client_id=" + githubClientID
		// Add email scope.
		oauthURL += "&scope=" + url.QueryEscape("read:user")
		// Our non-guessable state of 16 characters.

		// Adding the nonsene state within the query.
		oauthURL += "&state=" + nonsenseState
	} else if service == "gitlab" {
		// Base URL.
		oauthURL = "https://gitlab.com/oauth/authorize"
		// Add our app client ID.
		oauthURL += "?client_id=" + gitlabClientID
		// Define we want a code back
		oauthURL += "&response_type=code"
		// Add read_user scope.
		oauthURL += "&scope=read_user"
	}

	// Trying to follow our stateless design we encrypt the
	// Nonsense state so we later can re-use by decrypting it.
	// And than have the actual value. Also we use this to specify
	// From which site the callback was from.
	redirectURL := baseURL + "/api/callback/"
	if service == "github" {
		redirectURL += PrepareText(service+"+"+nonsenseState, AEAD_OAUTH) + "/"
	} else if service == "gitlab" {
		redirectURL += "gitlab/"
	}
	oauthURL += "&redirect_uri=" + redirectURL

	return oauthURL
}

func CallbackOAuth(tempCode, state, service string) OAuthResponse {
	if service == "" {
		return OAuthResponse{}
	}
	// Now the hard part D:
	// With our temp code and orignial state, we need to request the auth code.
	// With that auth code we need to ask nicely for the user's email.
	// Which is then passed back.

	// Base URL
	// Add our app client ID.
	// Add our client secret.
	var authURL string
	if service == "github" {
		authURL = "https://github.com/login/oauth/access_token"
		authURL += "?client_id=" + githubClientID
		authURL += "&client_secret=" + config.GITHUB_CLIENT_SECRET
		// Add the nonsense state we uses earlier.
		authURL += "&state=" + state
	} else if service == "gitlab" {
		authURL = "https://gitlab.com/oauth/token"
		authURL += "?client_id=" + gitlabClientID
		authURL += "&client_secret=" + config.GITLAB_CLIENT_SECRET
		// Define we log in trough the temp code.
		authURL += "&grant_type=authorization_code"
		// Specify the the redirect uri? It is required
		authURL += "&redirect_uri=" + url.PathEscape("http://localhost:3000/api/callback/gitlab/")
	}
	// Add the temp code.
	authURL += "&code=" + tempCode
	client := &http.Client{}
	req, err := http.NewRequest("POST", authURL, nil)
	if err != nil {
		return OAuthResponse{}
	}
	// Ensure we get a json response.
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return OAuthResponse{}
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return OAuthResponse{}
	}
	var responseJson OAuthTokenResponse

	err = json.NewDecoder(res.Body).Decode(&responseJson)
	if err != nil {
		return OAuthResponse{}
	}
	var userEndpoint string
	if service == "github" {
		userEndpoint = "https://api.github.com/user"
	} else if service == "gitlab" {
		userEndpoint = "https://gitlab.com/api/v4/user"
	}

	reqEmail, err := http.NewRequest("GET", userEndpoint, nil)
	if err != nil {
		return OAuthResponse{}
	}
	if service == "github" {
		// Recommended
		reqEmail.Header.Set("Accept", "application/vnd.github.v3+json")
	}

	reqEmail.Header.Set("Authorization", responseJson.TokenType+" "+responseJson.AccesToken)

	resEmail, err := client.Do(reqEmail)
	if err != nil {
		return OAuthResponse{}
	}
	defer resEmail.Body.Close()
	if resEmail.StatusCode != 200 {
		return OAuthResponse{}
	}

	var oauthResponse OAuthResponse
	err = json.NewDecoder(resEmail.Body).Decode(&oauthResponse)
	if err != nil {
		return OAuthResponse{}
	}

	return oauthResponse
}

package utils

import (
	"encoding/json"
	"net/http"
	"net/url"

	"userstyles.world/config"
)

var (
	githubClientID = "f3f2d378d88794062895"
)

// https://docs.github.com/en/developers/apps/authorizing-oauth-apps#response
type GithubResponse struct {
	AccesToken string `json:"access_token"`
	TokenType  string `json:"token_type"`
}

// https://docs.github.com/en/rest/reference/users#list-email-addresses-for-the-authenticated-user
type GithubEmailResponse struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

type OAuthResponse struct {
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
}

func GithubMakeURL(baseURL string) string {
	// Base URL.
	oathURL := "https://github.com/login/oauth/authorize"
	// Add our app client ID.
	oathURL += "?client_id=" + githubClientID
	// Add email scope.
	oathURL += "&scope=" + url.QueryEscape("user:email")

	// Our non-guessable state of 16 characters.
	nonsenseState := B2s(RandStringBytesMaskImprSrcUnsafe(16))

	// Adding the nonsene state within the query.
	oathURL += "&state=" + nonsenseState

	// Trying to follow our stateless design we encrypt the
	// Nonsense state so we later can re-use by decrypting it.
	// And than have the actual value. Also we use this to specify
	// From which site the callback was from.
	redirectURL := baseURL + "/api/callback"
	redirectURL += "/" + PrepareText("github+"+nonsenseState, AEAD_OAUTH)
	oathURL += "&redirect_uri=" + redirectURL

	return oathURL
}

func GithubCallbackOAuth(tempCode, state string) OAuthResponse {
	// Now the hard part D:
	// With our temp code and orignial state, we need to request the auth code.
	// With that auth code we need to ask nicely for the user's email.
	// Which is then passed back.

	// Base URL
	authURL := "https://github.com/login/oauth/access_token"
	// Add our app client ID.
	authURL += "?client_id=" + githubClientID
	// Add our client secret.
	authURL += "&client_secret=" + config.GITHUB_CLIENT_SECRET
	// Add the temp code.
	authURL += "&code=" + tempCode
	// Add the nonsense state we uses earlier.
	authURL += "&state=" + state

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
	var responseJson GithubResponse

	err = json.NewDecoder(res.Body).Decode(&responseJson)
	if err != nil {
		return OAuthResponse{}
	}
	reqEmail, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return OAuthResponse{}
	}
	// Recommended
	reqEmail.Header.Set("Accept", "application/vnd.github.v3+json")
	reqEmail.Header.Set("Authorization", responseJson.TokenType+" "+responseJson.AccesToken)

	resEmail, err := client.Do(reqEmail)
	if err != nil {
		return OAuthResponse{}
	}
	defer resEmail.Body.Close()
	if resEmail.StatusCode != 200 {
		return OAuthResponse{}
	}

	var responseEmailJson []GithubEmailResponse
	err = json.NewDecoder(resEmail.Body).Decode(&responseEmailJson)
	if err != nil {
		return OAuthResponse{}
	}
	var email GithubEmailResponse

	for _, a := range responseEmailJson {
		if a.Primary {
			email = a
			break
		}
	}

	return OAuthResponse{
		Email:    email.Email,
		Verified: email.Verified,
	}
}

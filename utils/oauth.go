package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"userstyles.world/config"
)

type OAuthTokenResponse struct {
	AccesToken string `json:"access_token"`
	TokenType  string `json:"token_type"`
}

type OAuthResponse struct {
	// Gitlab returns "username" for the username
	UserName string `json:"username"`
	// Github/Gitea-based returns "login" for the username
	LoginName string `json:"login"`
}

type GiteaLikeAccessJSON struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
	RedirectURI  string `json:"redirect_uri"`
}

func OauthMakeURL(service string) string {
	if service == "" {
		return ""
	}
	oauthURL := ""
	var nonsenseState string
	switch service {
	case "github":
		nonsenseState = B2s(RandStringBytesMaskImprSrcUnsafe(16))
		// Base URL.
		oauthURL = "https://github.com/login/oauth/authorize"
		// Add our app client ID.
		oauthURL += "?client_id=" + config.GITHUB_CLIENT_ID
		// Add email scope.
		oauthURL += "&scope=" + url.QueryEscape("read:user")
		// Our non-guessable state of 16 characters.
		oauthURL += "&state=" + nonsenseState
	case "gitlab":
		// Base URL.
		oauthURL = "https://gitlab.com/oauth/authorize"
		// Add our app client ID.
		oauthURL += "?client_id=" + config.GITLAB_CLIENT_ID
		// Define we want a code back
		oauthURL += "&response_type=code"
		// Add read_user scope.
		oauthURL += "&scope=read_user"
	case "codeberg":
		// Base URL.
		oauthURL = "https://codeberg.org/login/oauth/authorize"
		// Add our app client ID.
		oauthURL += "?client_id=" + config.CODEBERG_CLIENT_ID
		// Define we want a code back
		oauthURL += "&response_type=code"
	}
	if oauthURL == "" {
		return ""
	}

	// Trying to follow our stateless design we encrypt the
	// Nonsense state so we later can re-use by decrypting it.
	// And than have the actual value. Also we use this to specify
	// From which site the callback was from.
	redirectURL := config.OAuthURL()
	if service == "github" {
		redirectURL += PrepareText(nonsenseState, AEAD_OAUTH) + "/"
	} else {
		redirectURL += service + "/"
	}
	oauthURL += "&redirect_uri=" + redirectURL

	return oauthURL
}

func CallbackOAuth(tempCode, state, service string) (OAuthResponse, error) {
	if service == "" {
		return OAuthResponse{}, errors.New("no service detected")
	}
	// Now the hard part D:
	// With our temp code and orignial state, we need to request the auth code.
	// With that auth code we need to ask nicely for the user's email.
	// Which is then passed back.

	// Base URL
	// Add our app client ID.
	// Add our client secret.
	var authURL string
	var body GiteaLikeAccessJSON
	switch service {
	case "github":
		authURL = "https://github.com/login/oauth/access_token"
		authURL += "?client_id=" + config.GITHUB_CLIENT_ID
		authURL += "&client_secret=" + config.GITHUB_CLIENT_SECRET
		// Add the nonsense state we uses earlier.
		authURL += "&state=" + state
	case "gitlab":
		authURL = "https://gitlab.com/oauth/token"
		authURL += "?client_id=" + config.GITLAB_CLIENT_ID
		authURL += "&client_secret=" + config.GITLAB_CLIENT_SECRET
		// Define we log in trough the temp code.
		authURL += "&grant_type=authorization_code"
		// Specify the the redirect uri? It is required
		authURL += "&redirect_uri=" + url.PathEscape(config.OAuthURL()+"gitlab/")
	case "codeberg":
		authURL = "https://codeberg.org/login/oauth/access_token"
		body = GiteaLikeAccessJSON{
			ClientID:     config.CODEBERG_CLIENT_ID,
			ClientSecret: config.CODEBERG_CLIENT_SECRET,
			Code:         tempCode,
			GrantType:    "authorization_code",
			RedirectURI:  config.OAuthURL() + "codeberg/",
		}
	}
	if authURL == "" {
		return OAuthResponse{}, errors.New("no authURL was set")
	}
	if body.ClientID == "" {
		// Add the temp code.
		authURL += "&code=" + tempCode
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", authURL, nil)
	if err != nil {
		return OAuthResponse{}, err
	}
	// Ensure we get a json response.
	req.Header.Set("Accept", "application/json")
	if body.ClientID != "" {
		bodyString, err := json.Marshal(body)
		if err != nil {
			return OAuthResponse{}, err
		}
		req.Body = ioutil.NopCloser(strings.NewReader(B2s(bodyString)))
		req.ContentLength = int64(len(bodyString))
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := client.Do(req)
	if err != nil {
		return OAuthResponse{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return OAuthResponse{}, errors.New("status code returned wasn't OK 200")
	}
	var responseJSON OAuthTokenResponse

	err = json.NewDecoder(res.Body).Decode(&responseJSON)
	if err != nil {
		return OAuthResponse{}, err
	}
	var userEndpoint string
	switch service {
	case "github":
		userEndpoint = "https://api.github.com/user"
	case "gitlab":
		userEndpoint = "https://gitlab.com/api/v4/user"
	case "codeberg":
		userEndpoint = "https://codeberg.org/api/v1/user"
	}
	userInformationReq, err := http.NewRequest("GET", userEndpoint, nil)
	if err != nil {
		return OAuthResponse{}, err
	}
	if service == "github" {
		// Recommended
		userInformationReq.Header.Set("Accept", "application/vnd.github.v3+json")
	}

	userInformationReq.Header.Set("Authorization", responseJSON.TokenType+" "+responseJSON.AccesToken)

	resEmail, err := client.Do(userInformationReq)
	if err != nil {
		return OAuthResponse{}, err
	}
	defer resEmail.Body.Close()
	if resEmail.StatusCode != 200 {
		return OAuthResponse{}, errors.New("status code returned for userInformationReq wasn't OK 200")
	}

	var oauthResponse OAuthResponse
	err = json.NewDecoder(resEmail.Body).Decode(&oauthResponse)
	if err != nil {
		return OAuthResponse{}, err
	}

	if oauthResponse.LoginName != "" {
		oauthResponse.UserName = oauthResponse.LoginName
	}

	return oauthResponse, nil
}

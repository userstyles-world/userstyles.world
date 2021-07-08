package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"userstyles.world/modules/config"
	"userstyles.world/modules/errors"
)

const (
	codeberg = "codeberg"
	gitlab   = "gitlab"
	github   = "github"
)

type OAuthTokenResponse struct {
	AccesToken string `json:"access_token"`
	TokenType  string `json:"token_type"`
}

type userResponse struct {
	// Gitlab returns "username" for the username
	UserName string `json:"username"`
	// Github/Gitea-based returns "login" for the username
	LoginName string `json:"login"`

	// Gitlab has this bug with the email :)
	// BUt does include within the /user endpoint.
	Email string `json:"email"`
}

type OAuthResponse struct {
	Email    string
	Username string
}

type emailResponseStruct struct {
	Email string `json:"email"`

	// Github & Gitea
	Verified bool `json:"verified"`
	Primary  bool `json:"primary"`

	// Gitlab
	ConfirmedAt string `json:"confirmed_at"`
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
	case github:
		nonsenseState = UnsafeString(RandStringBytesMaskImprSrcUnsafe(16))
		// Base URL.
		oauthURL = "https://github.com/login/oauth/authorize"
		// Add our app client ID.
		oauthURL += "?client_id=" + config.GITHUB_CLIENT_ID
		// Add email scope.
		oauthURL += "&scope=" + url.QueryEscape("user:email")
		// Our non-guessable state of 16 characters.
		oauthURL += "&state=" + nonsenseState
	case gitlab:
		// Base URL.
		oauthURL = "https://gitlab.com/oauth/authorize"
		// Add our app client ID.
		oauthURL += "?client_id=" + config.GITLAB_CLIENT_ID
		// Define we want a code back
		oauthURL += "&response_type=code"
		// Add read_user scope.
		oauthURL += "&scope=read_user"
	case codeberg:
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
	if service == github {
		redirectURL += PrepareText(nonsenseState, AEAD_OAUTH) + "/"
	} else {
		redirectURL += service + "/"
	}
	oauthURL += "&redirect_uri=" + redirectURL

	return oauthURL
}

func CallbackOAuth(tempCode, state, service string) (OAuthResponse, error) {
	if service == "" {
		return OAuthResponse{}, errors.ErrNoServiceDetected
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
	case github:
		authURL = "https://github.com/login/oauth/access_token"
		authURL += "?client_id=" + config.GITHUB_CLIENT_ID
		authURL += "&client_secret=" + config.GITHUB_CLIENT_SECRET
		// Add the nonsense state we uses earlier.
		authURL += "&state=" + state
	case gitlab:
		authURL = "https://gitlab.com/oauth/token"
		authURL += "?client_id=" + config.GITLAB_CLIENT_ID
		authURL += "&client_secret=" + config.GITLAB_CLIENT_SECRET
		// Define we log in trough the temp code.
		authURL += "&grant_type=authorization_code"
		// Specify the the redirect uri? It is required
		authURL += "&redirect_uri=" + url.PathEscape(config.OAuthURL()+"gitlab/")
	case codeberg:
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
		return OAuthResponse{}, errors.ErrNoAuthURL
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
		req.Body = ioutil.NopCloser(strings.NewReader(UnsafeString(bodyString)))
		req.ContentLength = int64(len(bodyString))
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := client.Do(req)
	if err != nil {
		return OAuthResponse{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return OAuthResponse{}, errors.ErrNot200Ok
	}
	var responseJSON OAuthTokenResponse

	err = json.NewDecoder(res.Body).Decode(&responseJSON)
	if err != nil {
		return OAuthResponse{}, err
	}
	// Move the collecting of information.
	return getUserInformation(service, responseJSON)
}

func getUserInformation(service string, responseJSON OAuthTokenResponse) (OAuthResponse, error) {
	client := &http.Client{}
	var userEndpoint string
	switch service {
	case github:
		userEndpoint = "https://api.github.com/user"
	case gitlab:
		userEndpoint = "https://gitlab.com/api/v4/user"
	case codeberg:
		userEndpoint = "https://codeberg.org/api/v1/user"
	}
	userInformationReq, err := http.NewRequest("GET", userEndpoint, nil)
	if err != nil {
		return OAuthResponse{}, err
	}
	if service == github {
		// Recommended
		userInformationReq.Header.Set("Accept", "application/vnd.github.v3+json")
	}

	userInformationReq.Header.Set("Authorization", responseJSON.TokenType+" "+responseJSON.AccesToken)

	resUserInformation, err := client.Do(userInformationReq)
	if err != nil {
		return OAuthResponse{}, err
	}
	defer resUserInformation.Body.Close()
	if resUserInformation.StatusCode != 200 {
		return OAuthResponse{}, errors.ErrNot200Ok
	}

	var UserResponse userResponse
	var oauthResponse OAuthResponse
	err = json.NewDecoder(resUserInformation.Body).Decode(&UserResponse)
	if err != nil {
		return OAuthResponse{}, err
	}

	oauthResponse.Username = UserResponse.UserName
	if UserResponse.LoginName != "" {
		oauthResponse.Username = UserResponse.LoginName
	}
	oauthResponse.Username = strings.ToLower(oauthResponse.Username)

	if service == "gitlab" {
		if UserResponse.Email == "" {
			return OAuthResponse{}, errors.ErrPrimaryEmailNotVerified
		}
		oauthResponse.Email = UserResponse.Email
		return oauthResponse, nil
	}

	return getUserEmail(service, responseJSON, oauthResponse)
}

func getUserEmail(service string, responseJSON OAuthTokenResponse, oauthResponse OAuthResponse) (OAuthResponse, error) {
	client := &http.Client{}
	var emailEndpoint string
	switch service {
	case github:
		emailEndpoint = "https://api.github.com/user/emails"
	case codeberg:
		emailEndpoint = "https://codeberg.org/api/v1/user/emails"
	}
	emailInformationReq, err := http.NewRequest("GET", emailEndpoint, nil)
	if err != nil {
		return OAuthResponse{}, err
	}
	if service == github {
		// Recommended
		emailInformationReq.Header.Set("Accept", "application/vnd.github.v3+json")
	}

	emailInformationReq.Header.Set("Authorization", responseJSON.TokenType+" "+responseJSON.AccesToken)

	resEmailInformation, err := client.Do(emailInformationReq)
	if err != nil {
		return OAuthResponse{}, err
	}
	defer resEmailInformation.Body.Close()
	if resEmailInformation.StatusCode != 200 {
		return OAuthResponse{}, errors.ErrNot200Ok
	}
	var emailResponse []emailResponseStruct
	err = json.NewDecoder(resEmailInformation.Body).Decode(&emailResponse)
	if err != nil {
		return OAuthResponse{}, err
	}

	// Check if primary email is verified
	var email emailResponseStruct
	for i := 0; i < len(emailResponse); i++ {
		email = emailResponse[i]
		if email.Verified && email.Primary {
			oauthResponse.Email = email.Email
			break
		}
	}
	if oauthResponse.Email == "" {
		return OAuthResponse{}, errors.ErrPrimaryEmailNotVerified
	}

	return oauthResponse, nil
}

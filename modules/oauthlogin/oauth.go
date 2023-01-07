package oauthlogin

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"userstyles.world/modules/config"
	"userstyles.world/modules/errors"
	"userstyles.world/utils"
)

type Service string

const (
	CodebergService Service = "codeberg"
	GitlabService   Service = "gitlab"
	GithubService   Service = "github"
)

type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type userResponse struct {
	ID int `json:"id"`

	// Gitlab returns "username" for the username
	UserName string `json:"username"`
	// Github/Gitea-based returns "login" for the username
	LoginName string `json:"login"`

	// Gitlab has this bug with the email :)
	// But does include within the /user endpoint.
	Email string `json:"email"`
}

type OAuthResponse struct {
	Provider    Service
	ExternalID  int
	AccessToken string
	Email       string
	Username    string
	RawData     string
}

func (o *OAuthResponse) normalize(username string) {
	if username != "" {
		o.Username = username
	}
	o.Username = strings.ToLower(o.Username)
}

func (o *OAuthResponse) ProfileURL() string {
	switch o.Provider {
	case GithubService:
		return "https://github.com/" + o.Username
	case GitlabService:
		return "https://gitlab.com/" + o.Username
	case CodebergService:
		return "https://codeberg.org/" + o.Username
	default:
		return ""
	}
}

type emailResponseStruct struct {
	Email string `json:"email"`

	// Github & Gitea
	Verified bool `json:"verified"`
	Primary  bool `json:"primary"`

	// Gitlab
	ConfirmedAt string `json:"confirmed_at"`
}

type authURLPostBody struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
	RedirectURI  string `json:"redirect_uri"`
}

type ProviderFunctions interface {
	// Get the provider specific URL to be logged in with.
	oauthMakeURL() string

	// Let's the implementation append to the redirect URL.
	appendToRedirect(data any) string

	// See if the current implementation allows state.
	enableState() bool

	// Get the providers specic URL to get the auth token.
	getAuthTokenURL(data any) string

	// Check if the provider needs to set a POST body for the request.
	isAuthTokenPost() bool

	// And if the provider needs to set such body, we get it via the special function.
	getAuthTokenPostBody(data any) authURLPostBody

	// Let the provider's implementation do some work before sending the request.
	beforeRequest(body authURLPostBody, req *http.Request) error

	// Get the `/user` endpoint of the provider.
	getUserEndpoint() string

	// Return the type `Service` of the provider.
	getServiceType() Service

	// Get the '/email' endpoint of the provider.
	getEmailEndpoint() string
}

var (
	githubFunc   = github{}
	gitlabFunc   = gitlab{}
	codebergFunc = codeberg{}
)

func GetInterfaceForService(service string) (ProviderFunctions, error) {
	switch Service(service) {
	case GithubService:
		return githubFunc, nil
	case GitlabService:
		return gitlabFunc, nil
	case CodebergService:
		return codebergFunc, nil
	}
	return nil, errors.ErrNoServiceDetected
}

func OauthMakeURL(serviceType string) string {
	service, err := GetInterfaceForService(serviceType)
	if err != nil {
		return ""
	}

	oauthURL := service.oauthMakeURL()
	var state string
	if service.enableState() {
		state = utils.RandomString(16)
		oauthURL += "&state=" + state
	}
	if oauthURL == "" {
		return ""
	}

	oauthURL += "&redirect_uri=" + config.OAuthURL() + service.appendToRedirect(state)
	return oauthURL
}

func CallbackOAuth(tempCode, state, serviceType string) (OAuthResponse, error) {
	service, err := GetInterfaceForService(serviceType)
	if err != nil {
		return OAuthResponse{}, err
	}
	// Now the hard part D:
	// With our temp code and orignial state, we need to request the auth code.
	// With that auth code we need to ask nicely for the user's email.
	// Which is then passed back.

	var body authURLPostBody
	authURL := service.getAuthTokenURL(state)
	if authURL == "" {
		return OAuthResponse{}, errors.ErrNoAuthURL
	}

	if service.isAuthTokenPost() {
		body = service.getAuthTokenPostBody(tempCode)
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
	if err = service.beforeRequest(body, req); err != nil {
		return OAuthResponse{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return OAuthResponse{}, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return OAuthResponse{}, err
	}

	if res.StatusCode != 200 {
		return OAuthResponse{}, errors.ErrNot200Ok
	}

	var responseJSON OAuthTokenResponse
	err = json.Unmarshal(resBody, &responseJSON)
	if err != nil {
		return OAuthResponse{}, err
	}

	// Move the collecting of information.
	return getUserInformation(service, responseJSON)
}

func getUserInformation(service ProviderFunctions, responseJSON OAuthTokenResponse) (OAuthResponse, error) {
	client := &http.Client{}
	userEndpoint := service.getUserEndpoint()

	userInformationReq, err := http.NewRequest("GET", userEndpoint, nil)
	if err != nil {
		return OAuthResponse{}, err
	}
	if service.getServiceType() == GithubService {
		// Recommended
		userInformationReq.Header.Set("Accept", "application/vnd.github.v3+json")
	}

	userInformationReq.Header.Set("Authorization", responseJSON.TokenType+" "+responseJSON.AccessToken)

	resUserInformation, err := client.Do(userInformationReq)
	if err != nil {
		return OAuthResponse{}, err
	}
	defer resUserInformation.Body.Close()
	if resUserInformation.StatusCode != 200 {
		return OAuthResponse{}, errors.ErrNot200Ok
	}

	var userResponseJSON userResponse
	resBody, err := io.ReadAll(resUserInformation.Body)
	if err != nil {
		return OAuthResponse{}, err
	}

	err = json.Unmarshal(resBody, &userResponseJSON)
	if err != nil {
		return OAuthResponse{}, err
	}

	oauthResponse := OAuthResponse{
		Provider:    service.getServiceType(),
		ExternalID:  userResponseJSON.ID,
		Email:       userResponseJSON.Email,
		Username:    userResponseJSON.UserName,
		AccessToken: responseJSON.AccessToken,
		RawData:     string(resBody),
	}
	oauthResponse.normalize(userResponseJSON.LoginName)

	// GitLab returns email address early, so we can return here.
	if service.getServiceType() == GitlabService {
		if userResponseJSON.Email == "" {
			return OAuthResponse{}, errors.ErrPrimaryEmailNotVerified
		}
		return oauthResponse, nil
	}

	email, err := getUserEmail(service, responseJSON)
	if err != nil {
		return OAuthResponse{}, err
	}
	oauthResponse.Email = email

	return oauthResponse, nil
}

func getUserEmail(service ProviderFunctions, responseJSON OAuthTokenResponse) (string, error) {
	client := &http.Client{}
	emailEndpoint := service.getEmailEndpoint()

	emailInformationReq, err := http.NewRequest("GET", emailEndpoint, nil)
	if err != nil {
		return "", err
	}
	if service.getServiceType() == GithubService {
		// Recommended
		emailInformationReq.Header.Set("Accept", "application/vnd.github.v3+json")
	}

	emailInformationReq.Header.Set("Authorization", responseJSON.TokenType+" "+responseJSON.AccessToken)

	resEmailInformation, err := client.Do(emailInformationReq)
	if err != nil {
		return "", err
	}
	defer resEmailInformation.Body.Close()
	if resEmailInformation.StatusCode != 200 {
		return "", errors.ErrNot200Ok
	}

	var emailResponse []emailResponseStruct
	err = json.NewDecoder(resEmailInformation.Body).Decode(&emailResponse)
	if err != nil {
		return "", err
	}

	// Check if primary email is verified
	var email emailResponseStruct
	primaryEmail := ""
	for i := 0; i < len(emailResponse); i++ {
		email = emailResponse[i]
		if email.Verified && email.Primary {
			primaryEmail = email.Email
			break
		}
	}
	if primaryEmail == "" {
		return "", errors.ErrPrimaryEmailNotVerified
	}

	return primaryEmail, nil
}

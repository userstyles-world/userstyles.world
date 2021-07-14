package oauthlogin

import (
	"encoding/json"
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
	AccesToken string `json:"access_token"`
	TokenType  string `json:"token_type"`
}

type userResponse struct {
	// Gitlab returns "username" for the username
	UserName string `json:"username"`
	// Github/Gitea-based returns "login" for the username
	LoginName string `json:"login"`

	// Gitlab has this bug with the email :)
	// But does include within the /user endpoint.
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
	appendToRedirect(data interface{}) string

	// See if the current implementation allows state.
	enableState() bool

	// Get the providers specic URL to get the auth token.
	getAuthTokenURL(data interface{}) string

	// Check if the provider needs to set a POST body for the request.
	isAuthTokenPost() bool

	// And if the provider needs to set such body, we get it via the special function.
	getAuthTokenPostBody(data interface{}) authURLPostBody

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
		state = utils.UnsafeString(utils.RandStringBytesMaskImprSrcUnsafe(16))
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
		return OAuthResponse{}, errors.ErrNoServiceDetected
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

	// Because of gitlab oauth's implementation we already receive the email at this point.
	// Meaning we don't need to do another request.
	if service.getServiceType() == GitlabService {
		if UserResponse.Email == "" {
			return OAuthResponse{}, errors.ErrPrimaryEmailNotVerified
		}
		oauthResponse.Email = UserResponse.Email
		return oauthResponse, nil
	}

	return getUserEmail(service, responseJSON, oauthResponse)
}

func getUserEmail(service ProviderFunctions, responseJSON OAuthTokenResponse, oauthResponse OAuthResponse) (OAuthResponse, error) {
	client := &http.Client{}
	emailEndpoint := service.getEmailEndpoint()

	emailInformationReq, err := http.NewRequest("GET", emailEndpoint, nil)
	if err != nil {
		return OAuthResponse{}, err
	}
	if service.getServiceType() == GithubService {
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
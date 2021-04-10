package utils

import (
	"net/url"
)

var (
	githubClientID = "f3f2d378d88794062895"
)

func OATHGithub(baseURL string) string {
	// Base URL
	oathURL := "https://github.com/login/oauth/authorize"
	// Add our app client ID
	oathURL += "?client_id=" + githubClientID
	// Add email scope
	oathURL += "?scope=" + url.QueryEscape("user:email")

	// Our non-guessable state of 16 characters.
	nonsenseState := B2s(RandStringBytesMaskImprSrcUnsafe(16))

	// Adding the nonsene state within the query
	oathURL += "?state=" + nonsenseState

	// Trying to follow our stateless design we encrypt the
	// Nonsense state so we later can re-use by decrypting it.
	// And than have the actual value. Also we use this to specify
	// From which site the callback was from.
	redirectURL := baseURL + "api/callback"
	redirectURL += "/" + PrepareText("github+"+nonsenseState, AEAD_OAUTH)
	oathURL += "?redirect_uri=" + redirectURL

	return oathURL
}

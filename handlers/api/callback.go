package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"userstyles.world/utils"
)

func CallbackGet(c *fiber.Ctx) error {
	// Get the necessary information.
	redirectCode, tempCode, state := c.Params("rcode"), c.Query("code"), c.Query("state")

	if redirectCode == "" || tempCode == "" || state == "" {
		// Give them the bad enpoint error.
		return c.Next()
	}
	// Decode the string so we get our actual information back.
	code, err := utils.DecodePreparedText(redirectCode, utils.AEAD_OAUTH)
	if err != nil {
		return c.Next()
	}
	// We added the service within the the information and use the '+'
	// As seperator so now unseperate them.
	var service string
	var rState string
	if splitted := strings.Split(code, "+"); len(splitted) == 2 {
		service, rState = splitted[0], splitted[1]
	} else {
		return c.Next()
	}

	if rState != state {
		return c.Next()
	}

	var email utils.OAuthResponse
	switch service {
	case "github":
		email = utils.GithubCallbackOAuth(tempCode, rState)
	}
	if email == (utils.OAuthResponse{}) {
		return c.Next()
	}

	var verified string
	if email.Verified {
		verified = "verified"
	} else {
		verified = "unverified"
	}
	stateText := utils.PrepareText(verified+"+"+email.Email, utils.AEAD_OAUTH)

	return c.Render("more_info", fiber.Map{
		"Email": email.Email,
		"State": stateText,
	})
}

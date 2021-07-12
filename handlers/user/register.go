package user

import (
	"log"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func RegisterGet(c *fiber.Ctx) error {
	if u, ok := jwt.User(c); ok {
		log.Printf("User %d has set session, redirecting.", u.ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	}

	return c.Render("user/register", fiber.Map{
		"Title": "Register",
	})
}

func RegisterPost(c *fiber.Ctx) error {
	u := models.User{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
		Email:    c.FormValue("email"),
	}

	err := utils.Validate().StructPartial(u, "Username", "Email", "Password")
	if err != nil {
		errors := err.(validator.ValidationErrors)
		log.Println("Validation errors:", errors)

		return c.Status(fiber.StatusInternalServerError).
			Render("user/register", fiber.Map{
				"Title": "Register failed",
				"Error": "Failed to register. Make sure your input is correct.",
			})
	}

	token, err := utils.NewJWTToken().
		SetClaim("username", strings.ToLower(u.Username)).
		SetClaim("password", u.Password).
		SetClaim("email", u.Email).
		SetExpiration(time.Now().Add(time.Hour * 2)).
		GetSignedString(utils.VerifySigningKey)
	if err != nil {
		log.Println("Couldn't create a JWT Token, due to", err)
		return c.Status(fiber.StatusInternalServerError).
			Render("err", fiber.Map{
				"Title": "Internal server error",
			})
	}

	link := c.BaseURL() + "/verify/" + utils.PrepareText(token, utils.AEAD_CRYPTO)

	partPlain := utils.NewPart().
		SetBody("Verify your UserStyles.world account by clicking the link below.\n" +
			"The link will expire in 2 hours\n\n" +
			link + "\n\n" +
			"You can safely ignore this e-mail if you never made an account for UserStyles.world.")
	partHTML := utils.NewPart().
		SetBody("<p>Verify your UserStyles.world account by clicking the link below.</p>\n" +
			"<b>The link will expire in 2 hours</b>\n" +
			"<br>\n" +
			"<a target=\"_blank\" clicktracking=\"off\" href=\"" + link + "\">Verifcation link</a>\n" +
			"<br><br>\n" +
			"<p>You can safely ignore this e-mail if you never made an account for UserStyles.world.</p>").
		SetContentType("text/html")

	err = utils.NewEmail().
		SetTo(u.Email).
		SetSubject("Verify your email address").
		AddPart(*partPlain).
		AddPart(*partHTML).
		SendEmail()

	if err != nil {
		log.Println("Couldn't send a email, due to", err)
		return c.Status(fiber.StatusInternalServerError).
			Render("err", fiber.Map{
				"Title": "Internal server error",
				"Error": "Failed to send e-mail.",
			})
	}

	return c.Render("user/email-sent", fiber.Map{
		"Title":  "Email verifcation",
		"Reason": "Verification link has been sent to your e-mail address.",
	})
}

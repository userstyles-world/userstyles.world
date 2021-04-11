package user

import (
	"log"
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

	return c.Render("register", fiber.Map{
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

		c.SendStatus(fiber.StatusInternalServerError)
		if c.Locals("Email") != nil {
			return c.Render("more_info", fiber.Map{
				"Title": "Register failed",
				"Error": "Failed to register. Make sure your input is correct.",
				"Email": c.Locals("email"),
			})
		} else {
			return c.Render("register", fiber.Map{
				"Title": "Register failed",
				"Error": "Failed to register. Make sure your input is correct.",
			})
		}
	}

	jwt, err := utils.NewJWTToken().
		SetClaim("username", u.Username).
		SetClaim("password", u.Password).
		SetClaim("email", u.Email).
		SetExpiration(time.Hour * 2).
		GetSignedString(utils.VerifySigningKey)

	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
		})
	}

	link := c.BaseURL() + "/verify/" + utils.PrepareText(jwt, utils.AEAD_CRYPTO)

	PlainPart := utils.NewPart().
		SetBody("Verify your UserStyles.world account by clicking the link below.\n" +
			"The link will expire in 2 hours\n\n" +
			link + "\n\n" +
			"You can safely ignore this e-mail if you never made an account for UserStyles.world.")
	HTMLPart := utils.NewPart().
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
		AddPart(*PlainPart).
		AddPart(*HTMLPart).
		SendEmail()

	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"Error": "Failed to send e-mail.",
		})
	}

	return c.Render("send_email", fiber.Map{
		"Title":  "Email verifcation",
		"Reason": "Verification link has been sent to your e-mail address.",
	})
}

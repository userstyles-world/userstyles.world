package sessions

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"userstyles.world/models"
)

var store = session.New()

func GetStore() *session.Store {
	return store
}

func State(c *fiber.Ctx) *session.Session {
	s, err := store.Get(c)
	if err != nil {
		log.Println(err)
	}

	return s
}

func User(c *fiber.Ctx) *models.APIUser {
	s := State(c)

	// Type assertion will convert interface{} to other types.
	u := &models.APIUser{
		Username: s.Get("name").(string),
		Email:    s.Get("email").(string),
		ID:       s.Get("id").(uint),
	}

	return u
}

package sessions

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
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

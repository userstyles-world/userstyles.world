package sessions

import (
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store = session.New()

func GetStore() *session.Store {
	return store
}

package routes

import (
	"sync"

	controller "github.com/Improwised/xkcd-api/controllers/api/v1"
	"github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App) error {
	mu.Lock()

	app.Static("/assets/", "./assets")
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Render("./assets/index.html", fiber.Map{})
	})
	router := app.Group("/api")

	v1 := router.Group("/v1")

	err := setupItemController(v1)
	if err != nil {
		return err
	}

	mu.Unlock()
	return nil
}

func setupItemController(v1 fiber.Router) error {
	svcController, err := controller.NewItemController()
	if err != nil {
		return err
	}

	v1.Get("/getdata", svcController.XkcdGet)

	return nil
}

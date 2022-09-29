package v1

import (
	"net/http"

	"github.com/Improwised/xkcd-api/services"
	"github.com/Improwised/xkcd-api/utils"

	"github.com/gofiber/fiber/v2"
)

// ItemController for user controllers
type ItemController struct {
}

// NewItemController returns a user
func NewItemController() (*ItemController, error) {
	return &ItemController{}, nil
}

func (ctrl *ItemController) UserGet(c *fiber.Ctx) error {
	// fmt.Println("here comes")
	items, err := services.GetData()

	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}
	return utils.JSONWrite(c, http.StatusOK, items)
}

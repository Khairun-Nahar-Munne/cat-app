package controllers

import (
	"cat-app/models"

	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}




// Get handles the main page
func (c *MainController) Get() {
	// Get initial cat image
	resultChannel := make(chan models.VoteResponse)
	go func() {
		resultChannel <- fetchNewImage()
	}()
	result := <-resultChannel

	if result.Status == "success" {
		c.Data["ImageID"] = result.ImageID
		c.Data["ImageURL"] = result.ImageURL
	}

	c.TplName = "index.tpl"
}

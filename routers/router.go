package routers

import (
    "cat-app/controllers"
    "github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/api/vote", &controllers.MainController{}, "post:SubmitVote")
    beego.Router("/api/cat/fetch", &controllers.MainController{}, "get:FetchNewImage")
	beego.Router("/api/vote", &controllers.MainController{}, "get:GetVotes")

    beego.Router("/api/favorite", &controllers.MainController{}, "post:SubmitFavorite")
    beego.Router("/api/favourite", &controllers.MainController{}, "get:GetFavorites")
    beego.Router("/api/favourite/:id", &controllers.MainController{}, "delete:DeleteFavorite") 

    beego.Router("/api/breeds", &controllers.MainController{}, "get:GetBreeds")
    beego.Router("/api/breed/:id", &controllers.MainController{}, "get:GetBreedDetails")
    
}
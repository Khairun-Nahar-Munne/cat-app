package routers

import (
    "cat-app/controllers"
    "github.com/beego/beego/v2/server/web"
)

func init() {
    web.Router("/", &controllers.MainController{})
    web.Router("/api/vote", &controllers.MainController{}, "post:SubmitVote")
    web.Router("/api/cat/fetch", &controllers.MainController{}, "get:FetchNewImage")
	web.Router("/api/vote", &controllers.MainController{}, "get:GetVotes")

    web.Router("/api/favorite", &controllers.MainController{}, "post:SubmitFavorite")
    web.Router("/api/favourite", &controllers.MainController{}, "get:GetFavorites")
    web.Router("/api/favourite/:id", &controllers.MainController{}, "delete:DeleteFavorite") 

    web.Router("/api/breeds", &controllers.MainController{}, "get:GetBreeds")
    web.Router("/api/breed/:id", &controllers.MainController{}, "get:GetBreedDetails")
    
}
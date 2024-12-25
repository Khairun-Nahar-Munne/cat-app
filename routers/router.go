package routers

import (
	"cat-app/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/vote", &controllers.VoteController{}, "post:SubmitVote")
	beego.Router("/api/cat/fetch", &controllers.VoteController{}, "get:FetchNewImage")
	beego.Router("/api/vote", &controllers.VoteController{}, "get:GetVotes")

	beego.Router("/api/favorite", &controllers.FavController{}, "post:SubmitFavorite")
	beego.Router("/api/favourite", &controllers.FavController{}, "get:GetFavorites")
	beego.Router("/api/favourite/:id", &controllers.FavController{}, "delete:DeleteFavorite")

	beego.Router("/api/breeds", &controllers.BreedController{}, "get:GetBreeds")
	beego.Router("/api/breed/:id", &controllers.BreedController{}, "get:GetBreedDetails")

}

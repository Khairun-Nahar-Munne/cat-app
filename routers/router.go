package routers

import (
	"cat-app/controllers"
	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/", &controllers.VotingController{})
	web.Router("/voting", &controllers.VotingController{})
}

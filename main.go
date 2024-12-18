package main

import (
	_ "cat-app/routers"
	"github.com/beego/beego/v2/server/web"
)

func main() {
	web.Run()
}

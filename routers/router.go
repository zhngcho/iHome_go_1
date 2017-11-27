package routers

import (
	"github.com/astaxie/beego"
	"iHome_go_1/controllers"
)

func init() {

	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/v1.0/areas", &controllers.AreaController{}, "get:GetAreas")
	beego.Router("/api/v1.0/session", &controllers.SessionController{}, "get:Get")
	beego.Router("api/v1.0/users", &controllers.UserController{}, "post:Reg")
	//登陆
	beego.Router("api/v1.0/sessions", &controllers.UserController{}, "post:Login")

	//更新用户名
	beego.Router("api/v1.0/user/name", &controllers.UserController{}, "put:UpdateUsername")

	//上传文件头像
	beego.Router("api/v1.0/user/avatar", &controllers.UserController{}, "post:GetAvatar")

	//实名认证post
	beego.Router("api/v1.0/user/auth", &controllers.AuthController{}, "post:UpdateAuth")

	//实名认证校验get
	beego.Router("api/v1.0/user/auth", &controllers.AuthController{}, "get:GetAuth")

	//user
	beego.Router("api/v1.0/user", &controllers.UserController{}, "get:GetUser")

}

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

	//get用户发布的房源
	beego.Router("api/v1.0/user/houses", &controllers.UserHousesController{}, "get:GetUserHouses")

	// 发布房源路由
	// 发布房源路由
	beego.Router("api/v1.0/user/houses", &controllers.HouseController{}, "get:Gethouse")
	//请求查看房东/租客订单信息
	//查看租客信息
	//beego.Router("api/v1.0/user/orders?role=custom", &controllers.HouseController{}, "get:GetRenetrInfo")

	// 发布房源信息
	beego.Router("api/v1.0/houses", &controllers.HouseController{}, "post:NewHouse")

	//退出
	beego.Router("/api/v1.0/session", &controllers.HouseController{}, "delete:Delete")

	// 搜索房源
	beego.Router("/api/v1.0/houses", &controllers.HouseController{}, "get:GetHouseInfo")

	// 上传房源图片
	beego.Router("/api/v1.0/houses/:id([0-9])+/images", &controllers.HouseController{}, "post:Uplodpicture")

	//获取房源详细信息
	beego.Router("api/v1.0/houses/:id:int", &controllers.HouseDetailController{}, "get:GetHouseDetail")

}

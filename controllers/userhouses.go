package controllers

import (
	//"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	//"path"
)

type HouseInfoResp struct {
	Address     string `json:"address"`
	Areaname    string `json:"area_name"`
	Ctime       string `json:"ctime"`
	House_id    int    `json:"house_id"`
	Img_url     string `json:"img_url"`
	Order_count int    `json:"order_count"`
	Price       int    `json:"price"`
	Room_count  int    `json:"room_count"`
	Title       string `json:"title"`
	User_avatar string `json:"user_avatar"`
}

type HousesInfoResp struct {
	Houses []HouseInfoResp `json:"houses"`
}
type UserHousesResp struct {
	Errno  string         `json:"errno"`
	Errmsg string         `json:"errmsg"`
	Data   HousesInfoResp `json:"data"`
}

type UserHousesController struct {
	beego.Controller
}

func (this *UserHousesController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

// /api/1.0/user/houses [get]
func (this *UserHousesController) GetUserHouses() {
	resp := UserHousesResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)

	//从当前Session中得到当前客户端的user_id
	user_id := this.GetSession("user_id")

	//查询数据库
	var houses []models.House

	o := orm.NewOrm()
	//select * from user where mobile = request_data.mobile
	if _, err := o.QueryTable("house").Filter("user_id", user_id).RelatedSel().All(&houses); err == orm.ErrNoRows {
		//表示没有任何数据
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(resp.Errno)
		return

	}

	var houses_info []HouseInfoResp

	for _, value := range houses {
		var house_info HouseInfoResp
		house_info.Address = value.Address

		/////////////////////
		o.LoadRelated(&value, "Area")
		house_info.Areaname = value.Area.Name
		house_info.Ctime = value.Ctime.Format("2006-01-02 15:04:05")
		house_info.House_id = value.Id

		/////////////////////
		house_info.Img_url = value.Index_image_url
		house_info.Order_count = value.Order_count
		house_info.Price = value.Price
		house_info.Room_count = value.Room_count
		house_info.Title = value.Title
		house_info.User_avatar = value.User.Avatar_url

		houses_info = append(houses_info, house_info)
	}

	resp.Data = HousesInfoResp{Houses: houses_info}

	//更新Session
	this.SetSession("user_id", user_id)

	return

}

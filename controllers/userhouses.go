package controllers

import (
	//"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	//"path"
)

type UserHousesResp struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
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
	if _, err := o.QueryTable("house").Filter("user_id", user_id).All(&houses); err == orm.ErrNoRows {
		//表示没有任何数据
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(resp.Errno)
		return

	}

	resp.Data = houses

	//更新Session
	this.SetSession("user_id", user_id)

	return

}

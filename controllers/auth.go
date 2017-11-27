package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	_ "time"
)

type AuthController struct {
	beego.Controller
}

//update auth 业务请求
type UpdateAuthRequest struct {
	Real_name string `json:"real_name"`
	Id_card   string `json:"id_card"`
}

//update auth业务回复
type UpdateAuthResp struct {
	Errno  string `json:"errno"`
	Errmsg string `json:"errmsg"`
}

//实名信息校验业务回复
type AuthResp struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func (this *AuthController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

//校验实名信息
func (this *AuthController) GetAuth() {
	beego.Debug("get /api/v1.0/user/auth....")

	resp := AuthResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)

	//1. 从当前Session中得到当前客户端的user_id
	user_id := this.GetSession("user_id")

	//2.根据user_id 查找信息
	o := orm.NewOrm()
	user := models.User{Id: user_id.(int)}

	err := o.Read(&user)

	if err == orm.ErrNoRows {
		beego.Debug("查询不到")
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(models.RECODE_NODATA)
	} else if err == orm.ErrMissPK {
		beego.Debug("找不到主键")
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(models.RECODE_NODATA)
	} else {
		resp.Data = user
	}

	this.SetSession("user_id", user_id)

	return
}

//更新实名认证信息 /api/v1.0/user/auth [post]
func (this *AuthController) UpdateAuth() {
	beego.Debug("post /api/v1.0/user/auth....")

	resp := UpdateAuthResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)

	//1. 从当前Session中得到当前客户端的user_id
	user_id := this.GetSession("user_id")

	//得到用户的post请求的数
	//request
	var request_data UpdateAuthRequest
	json.Unmarshal(this.Ctx.Input.RequestBody, &request_data)

	fmt.Printf("request data : %+v\n", request_data)

	//校验信息

	if request_data.Real_name == "" || request_data.Id_card == "" {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//更新数据库 User 的 name字段
	o := orm.NewOrm()
	user := models.User{Id: user_id.(int), Real_name: request_data.Real_name, Id_card: request_data.Id_card}

	if _, err := o.Update(&user, "real_name", "id_card"); err != nil {
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Debug(err)
		return
	}

	//更新Session
	this.SetSession("user_id", user_id)

	//response data
	return
}

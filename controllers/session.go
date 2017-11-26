package controllers

import (
	"github.com/astaxie/beego"
	"iHome_go_1/models"
)

type SessionController struct {
	beego.Controller
}

type SessionResp struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

type Name struct {
	Name string `json:"name"`
}

func (this *SessionController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

func (this *SessionController) Get() {
	beego.Debug("get /api/v1.0/session....")

	resp := SessionResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	name := this.GetSession("name")
	if name == nil {
		resp.Errno = models.RECODE_SESSIONERR
		resp.Errno = models.RecodeText(resp.Errno)
		return
	}

	NameData := Name{Name: name.(string)}
	resp.Data = NameData

	return
}

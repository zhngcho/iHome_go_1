package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	"strconv"
	_ "time"
)

type AcceptOrderController struct {
	beego.Controller
}

//accept order业务请求
type AcceptOrderRequest struct {
	Action string `json:"action"`
	Reason string `json:"reason"`
}

//accept order业务回复
type AcceptOrderResp struct {
	Errno  string `json:"errno"`
	Errmsg string `json:"errmsg"`
}

func (this *AcceptOrderController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

//房东提交确认信息 /api/v1.0/orders/4/status [put]
func (this *AcceptOrderController) PutAcceptOrder() {
	beego.Debug("/api/v1.0/orders/4/status [put]....")

	resp := AcceptOrderResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)

	//1. 从当前Session中得到当前客户端的user_id
	user_id := this.GetSession("user_id")

	order_id_str := this.Ctx.Input.Param(":id")
	order_id, _ := strconv.Atoi(order_id_str)

	//得到用户的post请求的数据
	var request_data AcceptOrderRequest
	json.Unmarshal(this.Ctx.Input.RequestBody, &request_data)

	fmt.Printf("request data : %+v\n", request_data)

	//校验信息

	if request_data.Action == "" {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//查询订单，找到订单，确认订单状态为wait_accept
	o := orm.NewOrm()
	order := models.OrderHouse{Id: order_id}

	if err := o.Read(&order); err != nil {
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Debug(err)
		return
	}

	o.LoadRelated(&order, "House")
	//	o.LoadRelated(&order.House, "User")
	if order.Status != models.ORDER_STATUS_WAIT_ACCEPT || order.House.User.Id != user_id {
		fmt.Println("+++++77hang+++++\n")
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return

	}

	//接受订单业务
	var temp_order models.OrderHouse
	if request_data.Action == "accept" {
		temp_order = models.OrderHouse{Id: order_id, Status: models.ORDER_STATUS_WAIT_COMMENT}
		if _, err := o.Update(&temp_order, "status"); err != nil {
			resp.Errno = models.RECODE_DATAERR
			resp.Errmsg = models.RecodeText(resp.Errno)
			beego.Debug(err)
			return
		}
	} else {
		reason := request_data.Reason
		temp_order = models.OrderHouse{Id: order_id, Status: models.ORDER_STATUS_REJECTED, Comment: reason}
		if _, err := o.Update(&temp_order, "status", "comment"); err != nil {
			resp.Errno = models.RECODE_DATAERR
			resp.Errmsg = models.RecodeText(resp.Errno)
			beego.Debug(err)
			return
		}

	}

	//更新Session
	this.SetSession("user_id", user_id)

	//response data
	return
}

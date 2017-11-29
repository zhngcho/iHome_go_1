package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	"strconv"
	"time"
)

type PostOrderController struct {
	beego.Controller
}

//post order 业务请求
type PostOrderRequest struct {
	HouseId   string `json:"house_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

//post order业务回复的data结构体
type PostOrderRespData struct {
	OrderId string `json:"order_id"`
}

//post order业务回复
type PostOrderResp struct {
	Errno  string            `json:"errno"`
	Errmsg string            `json:"errmsg"`
	Data   PostOrderRespData `json:"data"`
}

func (this *PostOrderController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

//提交订单信息
func (this *PostOrderController) PostOrder() {
	beego.Debug("get /api/v1.0/orders....")

	resp := AuthResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)

	//从当前Session中得到当前客户端的user_id
	user_id := this.GetSession("user_id")

	//得到用户的post请求的数
	var request_data PostOrderRequest
	json.Unmarshal(this.Ctx.Input.RequestBody, &request_data)

	beego.Debug("request data : %+v\n", request_data)

	//校验信息

	if request_data.HouseId == "" || request_data.StartDate == "" || request_data.EndDate == "" || request_data.StartDate > request_data.EndDate {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//得到一共入住的天数
	startTime, _ := time.Parse("2006-01-02 15:04:05", request_data.StartDate+" 00:00:00")
	endTime, _ := time.Parse("2006-01-02 15:04:05", request_data.EndDate+" 00:00:00")
	durationDays := int(endTime.Sub(startTime).Hours()/24 + 1)

	//2.根据house_id查找相关房源信息
	o := orm.NewOrm()
	house_id, _ := strconv.Atoi(request_data.HouseId)
	house := models.House{Id: house_id}

	if err := o.Read(&house); err == orm.ErrNoRows {
		beego.Debug("查询不到")
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(models.RECODE_NODATA)
		return
	} else if err == orm.ErrMissPK {
		beego.Debug("找不到主键")
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(models.RECODE_NODATA)
		return
	} else if err != nil {
		beego.Debug("数据库查询错误")
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(models.RECODE_NODATA)
		return
	}

	if house.User.Id == user_id {
		beego.Debug("房东不能预订自己的房子")
		resp.Errno = models.RECODE_ROLEERR
		resp.Errmsg = models.RecodeText(models.RECODE_ROLEERR)
		return
	}

	//////////////////
	//判断时间的冲突//
	//////////////////

	var houseorder models.OrderHouse
	user := models.User{Id: user_id.(int)}
	houseorder.User = &user
	houseorder.House = &house
	houseorder.Begin_date = startTime
	houseorder.End_date = endTime
	houseorder.Days = durationDays
	houseorder.House_price = house.Price
	houseorder.Amount = house.Price * durationDays
	houseorder.Status = models.ORDER_STATUS_WAIT_ACCEPT
	//////houseorder.Comment = //////////////////////////////////nil
	fmt.Printf("houseorder to be inserted %+v\n", houseorder)

	id, err := o.Insert(&houseorder)
	if err != nil {
		beego.Info("inert error = ", err)
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	beego.Info("reg insert succ id = ", id)

	id_str := strconv.Itoa(int(id))
	resp.Data = PostOrderRespData{OrderId: id_str}

	this.SetSession("user_id", user_id)

	return
}

package controllers

import (
	_ "context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	"strconv"
)

type HouseController struct {
	beego.Controller
}

type HouseInfo struct {
	Title      string   `json:"title"`
	Price      string   `json:"price"`
	Area_id    string   `json:"area_id"`
	Address    string   `json:"address"`
	Room_count string   `json:"room_count"`
	Acreage    string   `json:"acreage"`
	Unit       string   `json:"unit"`
	Capacity   string   `json:"capacity"`
	Beds       string   `json:"beds"`
	Deposit    string   `json:"deposit"`
	Min_days   string   `json:"min_days"`
	Max_days   string   `json:"max_days"`
	Facility   []string `json:"facility"`
}

/*
type Facility_id struct {
	Facility []string `json:"facility"`
}
*/
type HouseResp struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

type House_id struct {
	Id int64 `json:"house_id"`
}

func (this *HouseController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

func (this *HouseController) Gethouse() {

	resp := HouseResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	resp.Data = House_id{}

	defer this.RetData(&resp)

}
func (this *HouseController) NewHouse() {
	resp := HouseResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	beego.Info(this.Ctx.Input.RequestBody)

	//	var house_info models.House
	var house_info HouseInfo

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &house_info)
	if err != nil {
		beego.Info("解码失败!")
		beego.Info(err)
		return
	}
	beego.Info(house_info)

	var data_base models.House

	// 获取用户Id
	user_id := this.GetSession("user_id")
	user := models.User{Id: user_id.(int)}
	data_base.User = &user

	temp_id, _ := strconv.Atoi(house_info.Area_id)
	area := models.Area{Id: temp_id}
	data_base.Area = &area

	data_base.Address = house_info.Address
	data_base.Acreage, _ = strconv.Atoi(house_info.Acreage)
	data_base.Beds = house_info.Beds
	data_base.Capacity, _ = strconv.Atoi(house_info.Capacity)
	data_base.Deposit, _ = strconv.Atoi(house_info.Deposit)
	data_base.Max_days, _ = strconv.Atoi(house_info.Max_days)
	data_base.Min_days, _ = strconv.Atoi(house_info.Min_days)
	data_base.Price, _ = strconv.Atoi(house_info.Price)
	data_base.Room_count, _ = strconv.Atoi(house_info.Room_count)
	data_base.Title = house_info.Title
	data_base.Unit = house_info.Unit

	beego.Info(data_base)
	o := orm.NewOrm()

	// 插入数据库
	house_id, err1 := o.Insert(&data_base)
	if err1 != nil {
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	beego.Info(house_id)

	facilities := []*models.Facility{}
	for _, fid := range house_info.Facility {

		id, _ := strconv.Atoi(fid)
		faci := &models.Facility{Id: id}
		facilities = append(facilities, faci)
	}

	// 创建多对多映射关系
	m2m := o.QueryM2M(&data_base, "Facilities")

	num, err := m2m.Add(facilities)
	if err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	beego.Info("num的value: ")
	beego.Info(num)

	resp.Data = House_id{Id: house_id}

}

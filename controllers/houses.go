package controllers

import (
	_ "context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	"path"
	"strconv"
)

type A1vatarUrl struct {
	Url string `json:"avatar_url"`
}

// 上传头像的返回结构
type A1vatarResp struct {
	Errno  string     `json:"errno"`
	Errmsg string     `json:"errmsg"`
	Data   A1vatarUrl `json:"data"`
}

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

func (this *HouseController) Uplodpicture() {
	resp := A1vatarResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)
	//获取房子的ID
	HouseId := this.Ctx.Input.Param(":id")
	//获取文件数据
	file, header, err := this.GetFile("house_image")

	if err != nil {
		resp.Errno = models.RECODE_SERVERERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Info("get file error")
		return
	}
	defer file.Close()

	//创建一个文件的缓冲
	fileBuffer := make([]byte, header.Size)

	_, err = file.Read(fileBuffer)
	if err != nil {
		resp.Errno = models.RECODE_IOERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Info("read file error")
		return
	}

	//home1.jpg
	suffix := path.Ext(header.Filename) // suffix = ".jpg"
	groupName, fileId, err1 := models.FDFSUploadByBuffer(fileBuffer, suffix[1:])
	if err1 != nil {
		resp.Errno = models.RECODE_IOERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Info("fdfs upload  file error")
		return
	}
	beego.Debug("groupName:", groupName, " fileId:", fileId)
	//添加Avatar_url字段到数据库中
	o := orm.NewOrm()

	House_id, _ := strconv.Atoi(HouseId)
	OneHouse := models.House{Id: House_id}

	houseImage := models.HouseImage{Url: fileId, House: &OneHouse}

	if _, err := o.Insert(&houseImage); err != nil {

		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//////////////////////房子的图片添加进  图片链接数据库
	if err := o.Read(&OneHouse); err != nil {
		resp.Errno = models.RECODE_DBERR
		fmt.Println("111111111111111111111")
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	house_image := models.HouseImage{House: &OneHouse, Url: fileId}
	OneHouse.Images = append(OneHouse.Images, &house_image)
	//根据house_id 查询house_image 是否为空
	if OneHouse.Index_image_url == "" {
		//如果为空 那么就用当前image_url为house的主image_url

		OneHouse.Index_image_url = fileId
		beego.Debug("set index_image_url ", fileId)
	}

	//将house_image入库
	if _, err := o.Insert(&house_image); err != nil {
		resp.Errno = models.RECODE_DBERR
		fmt.Println("211111111111111111111")
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Debug("insert house image error")
		return

	}

	//添加Avatar_url字段到数据库中

	if _, err := o.Update(&OneHouse); err != nil {
		resp.Errno = models.RECODE_DBERR
		fmt.Println("311111111111111111111")
		resp.Errmsg = models.RecodeText(resp.Errno)

		return
	}

	//拼接一个完整的路径
	avatar_url := models.AddDomain2Url(fileId)

	resp.Data.Url = avatar_url
	return
}

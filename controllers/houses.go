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

// 返回搜索房源结构体
type Query_data struct {
	Address     string `json:"address"`
	Area_Name   string `json:"area_name"`
	Ctime       string `json:"ctime"`
	House_Id    int    `json:"house_id"`
	Img_Url     string `json:"img_url"`
	Order_Count int    `json:"order_count"`
	Price       int    `json:"price"`
	Room_Count  int    `json:"room_count"`
	Title       string `json:"title"`
	User_Avatar string `json:"user_avatar"`
}

type Base_info struct {
	Current_page int         `json:"current_page"`
	Houses       interface{} `json:"houses"`
	Total_page   int         `json:"total_page"`
}

// 搜索房源
func (this *HouseController) GetHouseInfo() {

	resp := HouseResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	// api/v1.0/houses?aid=3&sd=2017-11-30&ed=2017-11-30&sk=new&p=1
	var (
		aid int
		sd  string
		ed  string
		sk  string
		p   int
	)
	this.Ctx.Input.Bind(&aid, "aid")
	this.Ctx.Input.Bind(&sd, "sd")
	this.Ctx.Input.Bind(&ed, "ed")
	this.Ctx.Input.Bind(&sk, "sk")
	this.Ctx.Input.Bind(&p, "p")

	beego.Info("测试!!!")
	fmt.Println(aid, sd, ed, sk, p)

	// 1、判断开始时间一定要小于结束时间
	// 2、检测p的值，一定要大于零
	if p <= 0 {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	// 3、尝试从redis数据库中获取数据,有的话返回数据

	// 4、没有从mysql数据库中查询数据，获取
	o := orm.NewOrm()
	qs := o.QueryTable("house")

	base := []models.House{}
	num, err := qs.Filter("area_id", aid).All(&base)
	if err != nil {

		beego.Info(err)
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	beego.Info(base)
	beego.Info(num)

	// 5、返回结果
	/*
		ret_data := make([]Query_data, num)

		for index, value := range base {
			ret_data[index].Address = value.Address
			ret_data[index].Area_Name = value.Area.Name
			ret_data[index].Ctime = value.Ctime.String()
			ret_data[index].House_Id = value.Id
			ret_data[index].Img_Url = value.Index_image_url
			ret_data[index].Order_Count = value.Order_count
			ret_data[index].Price = value.Price
			ret_data[index].Room_Count = value.Room_count
			ret_data[index].Title = value.Title
			ret_data[index].User_Avatar = value.User.Name
		}
	*/ /*
		var base_data Base_info
		base_data.Houses = &ret_data
		base_data.Current_page = 1
		base_data.Total_page = 1
		resp.Data = &base_data

	*/
	total_page := int(num)/models.HOUSE_LIST_PAGE_CAPACITY + 1
	house_page := 1

	house_list := []interface{}{}
	for _, house := range base {
		o.LoadRelated(&house, "Area")
		o.LoadRelated(&house, "User")
		o.LoadRelated(&house, "Images")
		o.LoadRelated(&house, "Facilities")
		house_list = append(house_list, house.To_house_info())
	}

	data := map[string]interface{}{}
	data["houses"] = house_list
	data["total_page"] = total_page
	data["current_page"] = house_page

	resp.Data = data

	// 将数据保存到redis中

	return

}

// 上传头像的返回结构
type A1vatarResp struct {
	Errno  string     `json:"errno"`
	Errmsg string     `json:"errmsg"`
	Data   A1vatarUrl `json:"data"`
}
type A1vatarUrl struct {
	Url string `json:"avatar_url"`
}

// 上传图片
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

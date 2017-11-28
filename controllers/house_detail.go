package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	"time"
	//"github.com/astaxie/beego/config"
	"strconv"
)

type CommentT struct {
	Comment   string `json:"comment"`
	Ctime     string `json:"ctime"`
	User_name string `"user_name"`
}

type HouseDetail struct {
	Acreage     int        `json:"acreage"`
	Address     string     `json:"address"`
	Beds        string     `json:"beds"`
	Capacity    int        `json:"capacity"`
	Comments    []CommentT `json:"comments"`
	Deposit     int        `json:"deposit"`
	Facilities  []int      `json:"facilities"`
	Hid         int        `json:"hid"`
	Img_urls    []string   `json:"img_urls"`
	Max_days    int        `json:"max_days"`
	Min_days    int        `json:"min_days"`
	Price       int        `json:"price"`
	Room_count  int        `json:"room_count"`
	Title       string     `json:"title"`
	Unit        string     `json:"unit"`
	User_avatar string     `json:"user_avatar"`
	User_id     int        `json:"user_id"`
	User_name   string     `json:"user_name"`
}
type HouseDetailDataResp struct {
	House   HouseDetail `json:"house"`
	User_id int         `json:"user_id"`
}

type HouseDetailResp struct {
	Errno  string              `json:"errno"`
	Errmsg string              `json:"errmsg"`
	Data   HouseDetailDataResp `json:"data"`
}

type HouseDetailController struct {
	beego.Controller
}

func (this *HouseDetailController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

// /api/1.0/houses/1 [get]
func (this *HouseDetailController) GetHouseDetail() {
	resp := HouseDetailResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)

	//从当前Session中得到当前客户端的user_id
	user_id := this.GetSession("user_id")

	//从url中得到房屋id
	house_id_str := this.Ctx.Input.Param(":id")
	house_id, _ := strconv.Atoi(house_id_str)

	//1 从redis查询是否有有area数据的缓存  如有有直接返回
	cache_conn, err := cache.NewCache("redis", `{"key":"ihome_go_1","conn":"127.0.0.1:6400","dbNum":"0"}`)
	if err != nil {
		beego.Debug("connect redis server error")
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//如果连接数据库成功 直接从redis中去“house_info” 将这个里面的value 直接返回给前端
	var house_info HouseDetailDataResp
	house_info_value := cache_conn.Get("house_detail_" + house_id_str)
	if house_info_value != nil {
		//代表缓存有数据， 直接将数据返回
		beego.Info("==== get area_info from cache =======")

		//将areas_info_value字符串变成 go的结构体
		json.Unmarshal(house_info_value.([]byte), &house_info)

		resp.Data = house_info
		return
	}

	//如果缓存中没有房屋详细信息，则查询数据库
	beego.Info("正在查询数据库.....\n")

	o := orm.NewOrm()
	house := models.House{Id: house_id}

	err = o.Read(&house)
	//表示没有任何数据
	if err == orm.ErrNoRows {
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	o.LoadRelated(&house, "Area")
	o.LoadRelated(&house, "User")
	o.LoadRelated(&house, "Facilities")
	o.LoadRelated(&house, "Images")
	o.LoadRelated(&house, "Orders")

	fmt.Printf("house: %+v\n", house)

	house_info.User_id = user_id.(int)
	house_info.House.Acreage = house.Acreage
	house_info.House.Address = house.Address
	house_info.House.Beds = house.Beds
	house_info.House.Capacity = house.Capacity

	var comments_value []CommentT
	for _, value := range house.Orders {
		var comment_value CommentT
		comment_value.Comment = value.Comment
		comment_value.Ctime = value.Ctime.Format("2006-01-02 15:04:05")

		o.LoadRelated(&value, "User")
		comment_value.User_name = value.User.Name
		comments_value = append(comments_value, comment_value)
	}

	house_info.House.Comments = comments_value

	house_info.House.Deposit = house.Deposit

	var facility_value []int
	for _, value := range house.Facilities {
		facility_value = append(facility_value, value.Id)
	}
	house_info.House.Facilities = facility_value

	house_info.House.Hid = house.Id

	var img_url_value []string
	for _, value := range house.Images {
		img_url_value = append(img_url_value, value.Url)
	}
	house_info.House.Img_urls = img_url_value

	house_info.House.Max_days = house.Max_days
	house_info.House.Min_days = house.Min_days
	house_info.House.Price = house.Price
	house_info.House.Room_count = house.Room_count
	house_info.House.Title = house.Title
	house_info.House.Unit = house.Unit
	house_info.House.User_avatar = house.User.Avatar_url
	house_info.House.User_id = house.User.Id
	house_info.House.User_name = house.User.Name

	resp.Data = house_info

	//将 house_info存储到缓存数据库中
	//将house_info转换成json字符串再存
	house_info_str, _ := json.Marshal(house_info)
	if err := cache_conn.Put("house_detail_"+house_id_str, house_info_str, 3600*time.Second); err != nil {
		beego.Debug("set house_detail_digit to cache error, err = ", err)
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//更新Session
	this.SetSession("user_id", user_id)

	return

}

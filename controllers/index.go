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
)

type HouseInfoIndex struct {
	Address     string `json:"address"`
	House_id    int    `json:"house_id"`
	Img_url     string `json:"img_url"`
	Order_count int    `json:"order_count"`
	Price       int    `json:"price"`
	Room_count  int    `json:"room_count"`
	Title       string `json:"title"`
	User_avatar string `json:"user_avatar"`
	Area_name   string `json:"area_name"`
	Ctime       string `json:"ctime"`
}

type HouseInfoIndexResp struct {
	Errno  string           `json:"errno"`
	Errmsg string           `json:"errmsg"`
	Data   []HouseInfoIndex `json:"data"`
}

type IndexController struct {
	beego.Controller
}

func (this *IndexController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

// /api/1.0/houses/index [get]
func (this *IndexController) GetHouseInfoIndex() {
	resp := HouseInfoIndexResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)

	//从当前Session中得到当前客户端的user_id
	user_id := this.GetSession("user_id")

	fmt.Printf("user_id:%+v\n", user_id)

	//1 从redis查询是否有有area数据的缓存  如有有直接返回
	cache_conn, err := cache.NewCache("redis", `{"key":"ihome_go_1","conn":"127.0.0.1:6400","dbNum":"0"}`)
	if err != nil {
		beego.Debug("connect redis server error")
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//如果连接数据库成功 直接从redis中去“houses_info_index” 将这个里面的value 直接返回给前端
	houses_info_value := cache_conn.Get("houses_info_index")
	if houses_info_value != nil {
		//代表缓存有数据， 直接将数据返回
		beego.Info("==== get area_info from cache =======")

		//将areas_info_value字符串变成 go的结构体
		json.Unmarshal(houses_info_value.([]byte), &resp.Data)

		return
	}

	//如果缓存中没有房屋详细信息，则查询数据库
	beego.Info("正在查询数据库.....\n")

	houses_info := make([]models.House, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("house").RelatedSel("User", "Area")

	image_num, _ := qs.Count()
	//表示没有任何数据
	if image_num == 0 {
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	if image_num <= 5 {
		qs.All(&houses_info)
	} else {
		for i := 0; i < 5; {
			qs.One(&houses_info[i])
			i++
		}
	}

	for _, value := range houses_info {
		var temp_house_info HouseInfoIndex
		temp_house_info.Address = value.Address
		temp_house_info.Area_name = value.Area.Name
		temp_house_info.Ctime = value.Ctime.Format("2006-01-02 15:04:05")
		temp_house_info.House_id = value.Id
		temp_house_info.Img_url = models.AddDomain2Url(value.Index_image_url)
		temp_house_info.Order_count = value.Order_count
		temp_house_info.Price = value.Price
		temp_house_info.Room_count = value.Room_count
		temp_house_info.Title = value.Title
		temp_house_info.User_avatar = models.AddDomain2Url(value.User.Avatar_url)
		resp.Data = append(resp.Data, temp_house_info)
	}

	fmt.Printf("resp.Data: %+v\n", resp.Data)

	//将 house_info存储到缓存数据库中
	//将house_info转换成json字符串再存
	houses_info_str, _ := json.Marshal(resp.Data)
	if err := cache_conn.Put("houses_info_index", houses_info_str, 3600*time.Second); err != nil {
		beego.Debug("set house_detail_digit to cache error, err = ", err)
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	fmt.Printf("houses_info_index: %+v\n", resp.Data)
	//更新Session
	this.SetSession("user_id", user_id)

	return

}

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
)

type AreaController struct {
	beego.Controller
}

type AreaResp struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func (this *AreaController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

func (this *AreaController) GetAreas() {
	beego.Debug("get /api/v1.0/areas....")

	resp := AreaResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)

	//1 从redis查询是否有有area数据的缓存  如有有直接返回
	cache_conn, err := cache.NewCache("redis", `{"key":"ihome_go_1","conn":"127.0.0.1:6400","dbNum":"0"}`)
	if err != nil {
		beego.Debug("connect redis server error")
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//如果连接数据库成功 直接从redis中去“area_info” 将这个里面的value 直接返回给前端
	areas_info_value := cache_conn.Get("area_info")
	if areas_info_value != nil {
		//代表缓存有数据， 直接将数据返回
		beego.Info("==== get area_info from cache =======")

		//将areas_info_value字符串变成 go的结构体
		var areas_info interface{}
		json.Unmarshal(areas_info_value.([]byte), &areas_info)

		resp.Data = areas_info
		return
	}

	//2 如果没有应该从数据库中查询
	o := orm.NewOrm()
	var areas []models.Area

	qs := o.QueryTable("area")
	num, err := qs.All(&areas)
	//select * from area
	if err != nil {
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	if num == 0 {
		//没有数据
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	fmt.Printf("areas = %+v\n", areas)
	resp.Data = areas

	//将 areas存储到缓存数据库中
	//将areas转换成json字符串再存
	areas_info_str, _ := json.Marshal(areas)
	if err := cache_conn.Put("area_info", areas_info_str, 3600*time.Second); err != nil {
		beego.Debug("set area_info to cache error, err = ", err)
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	return
}

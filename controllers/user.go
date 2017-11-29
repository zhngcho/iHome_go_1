package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/cache"
	//_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"iHome_go_1/models"
	//	"time"
	//"github.com/astaxie/beego/config"
	"path"
)

//reg客户端请求的数据
type RegRequest struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Sms_code string `json:"sms_code"`
}

//reg业务回复
type RegResp struct {
	Errno  string `json:"errno"`
	Errmsg string `json:"errmsg"`
}

//login 客户端请求的数据
type LoginResquest struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}
type LoginResp struct {
	Errno  string `json:"errno"`
	Errmsg string `json:"errmsg"`
}

type AvatarUrl struct {
	Url string `json:"avatar_url"`
}

// 上传头像的返回结构
type AvatarResp struct {
	Errno  string    `json:"errno"`
	Errmsg string    `json:"errmsg"`
	Data   AvatarUrl `json:"data"`
}

type NameResp struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

type UserResp struct {
	Errno  string      `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

type UserController struct {
	beego.Controller
}

// 返回查看订单信息结构体
type Orders_info struct {
	Amount     int    `json : "amount"`
	Comment    string `json : "comment"`
	Ctime      string `json : "ctime"`
	Days       int    `json : "days"`
	End_Date   string `json : "end_date"`
	Img_Url    string `json : "img_url"`
	Order_Id   int    `json : "order_id"`
	Start_Date string `json : "start_date"`
	Status     string `json : "status"`
	Title      string `json : "title"`
}

func (this *UserController) RetData(resp interface{}) {
	//给客户端返回json数据
	this.Data["json"] = resp
	//将json写回客户端
	this.ServeJSON()
}

//  /api/v1.0/users [post]
func (this *UserController) Reg() {
	resp := RegResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)

	//得到用户的post请求的数
	//request
	var request_data RegRequest
	json.Unmarshal(this.Ctx.Input.RequestBody, &request_data)

	fmt.Printf("request data : %+v\n", request_data)

	//校验信息

	if request_data.Mobile == "" || request_data.Password == "" || request_data.Sms_code == "" {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//对短信进行校验

	//将用户信息入库

	user := models.User{}

	user.Mobile = request_data.Mobile
	user.Password_hash = request_data.Password
	user.Name = request_data.Mobile

	o := orm.NewOrm()
	id, err := o.Insert(&user)
	if err != nil {
		fmt.Println("inert error = ", err)
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	beego.Info("reg insert succ id = ", id)

	//将用户名存入session中，
	//将用户的user_id存入session中
	//将用户的name存入session中
	this.SetSession("user_id", user.Id)
	this.SetSession("name", user.Name)
	this.SetSession("mobile", user.Mobile)

	return
}

// /api/v1.0/sessions [post] 登陆
func (this *UserController) Login() {
	resp := LoginResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)

	//得到用户的post请求的数
	//request
	var request_data LoginResquest
	json.Unmarshal(this.Ctx.Input.RequestBody, &request_data)

	fmt.Printf("request data : %+v\n", request_data)

	//校验参数合法性
	if request_data.Mobile == "" || request_data.Password == "" {
		resp.Errno = models.RECODE_PARAMERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//查询数据库
	var user models.User

	o := orm.NewOrm()
	//select * from user where mobile = request_data.mobile
	if err := o.QueryTable("user").Filter("mobile", request_data.Mobile).One(&user); err == orm.ErrNoRows {
		//表示没有任何数据
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(resp.Errno)
		return

	}

	//比对密码
	if user.Password_hash != request_data.Password {
		//密码错误
		resp.Errno = models.RECODE_PWDERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//存储session
	this.SetSession("user_id", user.Id)
	this.SetSession("name", user.Mobile)
	this.SetSession("mobile", user.Mobile)

	return
}

// api/v1.0/user/name [put]
func (this *UserController) UpdateUsername() {
	resp := NameResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	//从session得到user_id
	user_id := this.GetSession("user_id")

	type Name struct {
		Name string `json:"name"`
	}
	//request post data
	var req_name Name
	//得到客户端请求数据
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req_name); err != nil {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	if req_name.Name == "" {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = "name is Empty!"
		return
	}

	//更新数据库 User 的 name字段
	o := orm.NewOrm()
	user := models.User{Id: user_id.(int), Name: req_name.Name}

	if _, err := o.Update(&user, "name"); err != nil {
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Debug(err)
		return
	}

	//更新Session
	this.SetSession("user_id", user_id)
	this.SetSession("name", req_name.Name)

	//response data
	resp.Data = req_name
	return
}

// /api/1.0/user/avatar [post]
func (this *UserController) GetAvatar() {
	resp := AvatarResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)

	//获取文件数据
	file, header, err := this.GetFile("avatar")

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

	beego.Info("groupname,", groupName, " file id ", fileId)

	//通过session得到当前用户
	user_id := this.GetSession("user_id")

	//添加Avatar_url字段到数据库中
	o := orm.NewOrm()
	user := models.User{Id: user_id.(int), Avatar_url: fileId}

	if _, err := o.Update(&user, "avatar_url"); err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//拼接一个完整的路径
	avatar_url := models.AddDomain2Url(fileId)

	resp.Data.Url = avatar_url
	return
}

func (this *UserController) GetUser() {
	beego.Debug("get /api/v1.0/user....")

	resp := UserResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

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

// api/v1.0/user/orders?role=custom
// 查看我的订单
func (this *UserController) GetOrders() {
	resp := UserResp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	// 根据session或区user_id
	user_id := this.GetSession("user_id")

	beego.Info(user_id)

	// 绑定URL参数，判断是房东还是租客， 如果为空则返回错误信息
	var role string
	this.Ctx.Input.Bind(&role, "role")

	// 返回信息变量
	var data []Orders_info
	var order_info []models.OrderHouse
	var row int64
	var err error
	o := orm.NewOrm()

	if role == "" {
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	} else if role == "landlord" {
		// 如果是房东返回自己发布的房屋订单表，查询mysql
		beego.Info("我是房东")

		var house_info []models.House

		// 查询house表，自己发布的房源, 获取house_id
		row, err = o.QueryTable("house").Filter("user_id", user_id).All(&house_info, "id")
		if err != nil {
			resp.Errno = models.RECODE_DBERR
			resp.Errmsg = models.RecodeText(resp.Errno)
			beego.Info(err)
			return
		}

		beego.Info(house_info)

		if row > 0 {
			// 根据house_id查询房屋订单表
			for _, value := range house_info {
				row, err = o.QueryTable("order_house").Filter("house_id", value.Id).All(&order_info)
			}
		}
		fmt.Println("房东查询到的行数: ", row)

	} else if role == "custom" {
		// 租客
		beego.Info("我是租客")

		row, err = o.QueryTable("order_house").Filter("user_id", user_id).All(&order_info)
		if err != nil {
			resp.Errno = models.RECODE_DBERR
			resp.Errmsg = models.RecodeText(resp.Errno)
			return
		}
		fmt.Println("租客查询到的行数: ", row)
	}

	if row > 0 {
		// 根据house_id查询房屋订单表
		for index, _ := range order_info {
			data[index].Amount = order_info[index].Amount
			data[index].Comment = order_info[index].Comment
			data[index].Ctime = order_info[index].Ctime.String()
			data[index].Days = order_info[index].Days
			data[index].End_Date = order_info[index].End_date.String()
			data[index].Img_Url = order_info[index].House.Index_image_url
			data[index].Order_Id = order_info[index].Id
			data[index].Start_Date = order_info[index].Begin_date.String()
			data[index].Status = order_info[index].Status
			data[index].Title = order_info[index].House.Title
		}
	}

	// 返回json
	resp.Data = &data
	return
}

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

type UserController struct {
	beego.Controller
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
	avatar_url := "http://101.200.170.171:8080/" + fileId

	resp.Data.Url = avatar_url
	return
}

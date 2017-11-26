package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"iHome_go_1/models"
	_ "iHome_go_1/routers"
	"net/http"
	"strings"
)

func init() {
	// set default database
	//绑定orm此时用的是哪个数据库的驱动
	//第四个参数表示 数据连接最大的空闲个数(可选)
	//第五个参数表示 数据库最大的链接个数(可选)
	orm.RegisterDataBase("default", "mysql", "root:mysql@tcp(127.0.0.1:3306)/ihome_go_1?charset=utf8", 30)

	// register model
	//注册orm都有哪些模块， 目前orm需用同步哪些表
	orm.RegisterModel(new(models.User), new(models.Area), new(models.Facility), new(models.House), new(models.HouseImage), new(models.OrderHouse))

	// create table
	//第二个参数表示是否强制替换
	//第三个表示 如果没有是否创建
	orm.RunSyncdb("default", false, true)
}

func main() {
	//设置一个fastdfs 请求的静态路径
	//http://101.200.170.171:8080/group1/M00/00/00/Zciqq1oaGW-ABnxDAAAHFIcthTk%207176.go
	beego.SetStaticPath("/group1/M00", "fastdfs/storage_data/data")

	//测试fastdfs接口
	//	models.FDFSUploadByFileName("home01.jpg")

	ignoreStaticPath()
	beego.Run()
}

//重定向static静态路径
func ignoreStaticPath() {

	//透明static

	beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)
	beego.InsertFilter("/*", beego.BeforeRouter, TransparentStatic)
}

func TransparentStatic(ctx *context.Context) {
	orpath := ctx.Request.URL.Path
	beego.Debug("request url: ", orpath)
	//如果请求uri还有api字段,说明是指令应该取消静态资源路径重定向
	if strings.Index(orpath, "api") >= 0 {
		return
	}
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/html/"+ctx.Request.URL.Path)

	//将全部的静态资源重定向 加上/static/html路径
	//http://ip:port:8080/index.html----> http://ip:port:8080/static/html/index.html
	//如果restFUL api  那么就取消冲定向
	//http://ip:port:8080/api/v1.0/areas ---> http://ip:port:8080/static/html/api/v1.0/areas
}

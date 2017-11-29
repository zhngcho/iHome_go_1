package controllers

import (
	_ "context"
	_ "encoding/json"
	_ "fmt"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/orm"
	_ "iHome_go_1/models"
	_ "path"
	_ "strconv"
)

type Comment struct {
	Order_id string `json:"order_id"`
	Comment  string `json:"comment"`
}

type Retcomment struct {
	Errno  string `json:"errno"`
	Errmsg string `json:"errmsg"`
}
type Housecomment struct {
	beego.Controller
}

func (this *HouseController) Comment() {

}

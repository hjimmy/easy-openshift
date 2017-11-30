/*
* @Author: haodaquan
* @Date:   2017-08-16 10:27:40
* @Last Modified by:   haodaquan
* @Last Modified time: 2017-08-16 09:17:22
 */

package controllers

import (
	"github.com/astaxie/beego"
	"github.com/hjimmy/easy-openshift/libs"
	"github.com/hjimmy/easy-openshift/models"
        "github.com/hjimmy/easy-openshift/openshift"
	"strconv"
	"strings"
	"time"
	"fmt"
)

type AppController struct {
	BaseController
}

func (this *AppController) List() {
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}

	result, count := models.TaskAppGetList(page, this.pageSize)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
<<<<<<< HEAD
		row["type"] = v.Type
=======
		apptype := v.Type
                switch apptype {
                    case 0:
			row["type"] = "owncloud"
                    case 1:
			row["type"] = "mysql"
		}
>>>>>>> 9170e3490fd9e3343c696cb5ee73c67accb698fd
		row["name"] = v.Name
		row["ip"] = v.Ip
                row["port"] = v.Port
		row["replica"] = v.Replica
                row["route"] = v.Route
                row["size"] = v.Size
                row["type"] = v.Type
		row["detail"] = v.Detail
		row["create_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		list[k] = row
	}
	this.Data["pageTitle"] = "应用列表"
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("AppController.List"), true).ToString()
	this.display()
}

func (this *AppController) Add() {
	if this.isPost() {
		app := new(models.TaskApp)
		app.Name = strings.TrimSpace(this.GetString("name"))
                fmt.Println(app.Name)
		app.Ip = strings.TrimSpace(this.GetString("ip"))
		app.Port,_ = strconv.Atoi(this.GetString("port"))
                app.Replica,_ = strconv.Atoi(this.GetString("replica"))
<<<<<<< HEAD
                app.Size,_ = strconv.Atoi(this.GetString("size"))
                app.Route = strings.TrimSpace(this.GetString("route"))
		app.Type = strings.TrimSpace(this.GetString("type"))
=======
                app.Size = strings.TrimSpace(this.GetString("size"))
                app.Route = strings.TrimSpace(this.GetString("route"))
		app.Type,_ = strconv.Atoi(this.GetString("type"))
>>>>>>> 9170e3490fd9e3343c696cb5ee73c67accb698fd
		app.Detail = strings.TrimSpace(this.GetString("detail"))
		app.CreateTime = time.Now().Unix()
		app.UpdateTime = time.Now().Unix()
		app.Status = 0
<<<<<<< HEAD
                fmt.Println(app.Type)
                if app.Type == "owncloud"{
		    fmt.Println(openshift.Serveraddr)
                    openshift.Create_owncloud(app.Name, app.Port, app.Size, app.Type, app.Replica)
=======
                fmt.Println(app_Type)
                if app.Type == 0{
                    openshift.create_owncloud(app.Name, app.Port, app.Size, owncloud, app.Replica)
>>>>>>> 9170e3490fd9e3343c696cb5ee73c67accb698fd
                }
		_, err := models.TaskAppAdd(app)
                
		if err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}
		this.ajaxMsg("", MSG_OK)
	}
	this.Data["pageTitle"] = "添加服务器"
	this.display()
}

func (this *AppController) Edit() {
	id, _ := this.GetInt("id")
<<<<<<< HEAD
        fmt.Println(id)
=======
>>>>>>> 9170e3490fd9e3343c696cb5ee73c67accb698fd
	app, err := models.TaskAppGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	if this.isPost() {
		app.Name = strings.TrimSpace(this.GetString("name"))
		app.Ip = strings.TrimSpace(this.GetString("ip"))
<<<<<<< HEAD
                app.Port,_ = strconv.Atoi(this.GetString("port")) 
		app.Replica,_ = strconv.Atoi(this.GetString("replica"))
		app.Route = strings.TrimSpace(this.GetString("route"))
		app.Size,_ = strconv.Atoi(this.GetString("size"))
		app.Type = strings.TrimSpace(this.GetString("type"))
		app.Detail = strings.TrimSpace(this.GetString("detail"))
		app.UpdateTime = time.Now().Unix()
		app.Status = 0
		
		if app.Type == "owncloud"{
		    openshift.Update_owncloud(app.Name, app.Type, app.Port, app.Replica)
                }
=======
		app.Type,_ = strconv.Atoi(this.GetString("type"))
		app.Detail = strings.TrimSpace(this.GetString("detail"))
		//app.UpdateTime = time.Now().Unix()
		app.Status = 0
>>>>>>> 9170e3490fd9e3343c696cb5ee73c67accb698fd
		err := app.Update()
		if err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}
		this.ajaxMsg("", MSG_OK)
	}

	this.Data["pageTitle"] = "编辑服务器"
	this.Data["app"] = app
	this.display()
}

//TODO删除更新
func (this *AppController) Batch() {
	action := this.GetString("action")
	ids := this.GetStrings("ids")
	if len(ids) < 1 {
		this.ajaxMsg("请选择要操作的项目", MSG_ERR)
	}

	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}
		switch action {
		case "delete":
<<<<<<< HEAD
                     fmt.Println(id)
		     app, err := models.TaskAppGetById(id)
                     if err != nil {
                       this.showMsg(err.Error())
                     }
                     if app.Type == "owncloud"{
                    	    fmt.Println(openshift.Serveraddr)
                    	    openshift.Delete_owncloud(app.Name, app.Type)
			    //删除数据库中内容
                            models.TaskAppDelById(id)
		     }
=======
			//查询服务器是否被占用
		        models.TaskAppDelById(id)
>>>>>>> 9170e3490fd9e3343c696cb5ee73c67accb698fd
		}
	}

	this.ajaxMsg("", MSG_OK)
}

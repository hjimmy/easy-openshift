package main

import (
	"github.com/astaxie/beego"
	"github.com/hjimmy/easy-openshift/models"
	_ "github.com/hjimmy/easy-openshift/routers"
	"github.com/hjimmy/easy-openshift/jobs"
)

const (
	VERSION = "1.0.0"
)

func init() {
	//初始化数据模型
	models.Init()
	jobs.InitJobs()
}

func main() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run()
}

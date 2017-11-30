/*
* @Author: haodaquan
* @Date:   2017-08-16 12:22:37
* @Last Modified by:   haodaquan
* @Last Modified time: 2017-08-16 12:22:55
 */

package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type TaskApp struct {
	Id            int
	Name    string
	Ip      string
        Port    int
	Replica  int
	Route    string
<<<<<<< HEAD
	Size     int
	Type          string
=======
	Size     string
	Type          int
>>>>>>> 9170e3490fd9e3343c696cb5ee73c67accb698fd
	Detail        string
	CreateTime    int64
	UpdateTime    int64
	Status        int
}

func (t *TaskApp) TableName() string {
	return TableName("task_app")
}

func (t *TaskApp) Update(fields ...string) error {
	if t.Name == "" {
		return fmt.Errorf("App名不能为空")
	}
<<<<<<< HEAD
        fmt.Println(t.Ip)
=======
>>>>>>> 9170e3490fd9e3343c696cb5ee73c67accb698fd
	if t.Ip == "" {
		return fmt.Errorf("App IP 地址不能为空")
	}

        fmt.Println(t.Port)
        if t.Port == 0 {
                return fmt.Errorf("App端口不能为空")
        }

        if t.Route == "" {
                return fmt.Errorf("路由不能为空")
        }



	if t.Replica <= 0 {
		return fmt.Errorf("副本数必须大于1")
	}

<<<<<<< HEAD
	if t.Type == ""  {
=======
	if t.Type == 0  {
>>>>>>> 9170e3490fd9e3343c696cb5ee73c67accb698fd
		return fmt.Errorf("App 类型不能为空")
	}


	if _, err := orm.NewOrm().Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func TaskAppAdd(obj *TaskApp) (int64, error) {
        fmt.Println(obj)	
	fmt.Println(obj.Type)
	if obj.Name == "" {
                return 0, fmt.Errorf("App名不能为空")
        }
        if obj.Ip == "" {
                return 0, fmt.Errorf("App地址不能为空")
        }

        if obj.Replica == 0 {
                return 0, fmt.Errorf("登录账户不能为空")
        }

<<<<<<< HEAD
        if obj.Type == ""  {
=======
        if obj.Type == 0  {
>>>>>>> 9170e3490fd9e3343c696cb5ee73c67accb698fd
                return 0, fmt.Errorf("App 类型不能为空")
        }

	return orm.NewOrm().Insert(obj)
}

func TaskAppGetById(id int) (*TaskApp, error) {
	obj := &TaskApp{
		Id: id,
	}
	err := orm.NewOrm().Read(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func TaskAppDelById(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("task_app")).Filter("id", id).Delete()
	return err
}

func TaskAppGetList(page, pageSize int) ([]*TaskApp, int64) {
	offset := (page - 1) * pageSize
	list := make([]*TaskApp, 0)
	query := orm.NewOrm().QueryTable(TableName("task_app"))
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)

	return list, total
}

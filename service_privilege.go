// Title：权限相关的服务
//
// Description:
//
// Author:Samurai
//
// Createtime:2013-09-23 10:06
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
package lessgo

import (
	"github.com/hjqhezgh/commonlib"
	"net/http"
)

const (
	SESSION_USER      = "SESSION_USER"      //用户登录后信息存储
	KEY_USER_ID       = "KEY_USER_ID"       //用户ID
	KEY_USER_NAME     = "KEY_USER_NAME"     //用户名
	KEY_REALLY_NAME   = "KEY_REALLY_NAME"   //真实姓名
	KEY_DEPARTMENT_ID = "KEY_DEPARTMENT_ID" //部门ID
	KEY_ROLE_ID       = "KEY_ROLE_ID"       //角色ID
	KEY_ROLE_CODE     = "KEY_ROLE_CODE"     //角色代号
	KEY_ROLE_LEVEL    = "KEY_ROLE_LEVEL"    //角色等级
)

type Menu struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Url      string `json:"url"`
	Children []Menu `json:"children"`
}

//用于存储登陆员工的信息
type Employee struct {
	UserId       string `json:"userId"`
	UserName     string `json:"userName"`
	ReallyName   string `json:"reallyName"`
	DepartmentId string `json:"departmentId"`
	RoleId       string `json:"roleId"`
	RoleCode     string `json:"roleCode"`
	RoleLevel    string `json:"roleLevel"`
}

func GetMenus(username string) []Menu {

	var menus []Menu
	var parent_id string
	db := GetMySQL()
	defer db.Close()
	sql := `select a.action_id,a.action_name,a.icon,a.parent_id,a.url from action a where
				a.action_id in (select distinct(ra.action_id) from role_action ra where
					ra.role_id in (select er.role_id from employee_role er where
						er.user_id=(select e.user_id from employee e where e.username=?))) order by a.parent_id,a.action_id`

	rows, err := db.Query(sql, username)
	if err != nil {
		Log.Error(err.Error())
		return nil
	}
	for rows.Next() {
		menu := new(Menu)
		err := rows.Scan(&menu.Id, &menu.Name, &menu.Icon, &parent_id, &menu.Url)
		if err != nil {
			Log.Error(err.Error())
			return nil
		}

		if parent_id > "0" {
			for i := 0; i < len(menus); i++ {
				if menus[i].Id == parent_id {
					menus[i].Children = append(menus[i].Children, *menu)
				}
			}
		} else {
			menus = append(menus, *menu)
		}
	}
	return menus
}

func QueryMenusAction(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})

	username := r.FormValue("username")
	if username == "" {
		Log.Error("username is NULL!")
		return
	}

	data["menus"] = GetMenus(username)
	Log.Debug(data)
	commonlib.OutputJson(w, data, "")
}

//获取当前登陆用户
func GetCurrentEmployee(r *http.Request) Employee {

	session, err := Store.Get(r, SESSION_USER)

	if err != nil {
		Log.Error("前台用户获取session发生错误，信息如下：", err.Error())
		return Employee{}
	}

	user_id, ok := session.Values[KEY_USER_ID].(string)
	if !ok {
		Log.Error("get session value error!", ok)
		return Employee{}
	}

	user_name, ok := session.Values[KEY_USER_NAME].(string)
	if !ok {
		Log.Error("get session value error!", ok)
		return Employee{}
	}

	really_name, ok := session.Values[KEY_REALLY_NAME].(string)
	if !ok {
		Log.Error("get session value error!", ok)
		return Employee{}
	}

	department_id, ok := session.Values[KEY_DEPARTMENT_ID].(string)
	if !ok {
		Log.Error("get session value error!", ok)
		return Employee{}
	}

	role_id, ok := session.Values[KEY_ROLE_ID].(string)
	if !ok {
		Log.Error("get session value error!", ok)
		return Employee{}
	}

	role_code, ok := session.Values[KEY_ROLE_CODE].(string)
	if !ok {
		Log.Error("get session value error!", ok)
		return Employee{}
	}

	role_level, ok := session.Values[KEY_ROLE_LEVEL].(string)
	if !ok {
		Log.Error("get session value error!", ok)
		return Employee{}
	}

	return Employee{
		UserId:       user_id,
		UserName:     user_name,
		ReallyName:   really_name,
		DepartmentId: department_id,
		RoleId:       role_id,
		RoleCode:     role_code,
		RoleLevel:    role_level,
	}
}

//设置当前用户信息
func SetCurrentEmployee(employee Employee, w http.ResponseWriter, r *http.Request) {

	session, err := Store.Get(r, SESSION_USER)

	if err != nil {
		Log.Error(err)
	}

	session.Values[KEY_USER_ID] = employee.UserId
	session.Values[KEY_USER_NAME] = employee.UserName
	session.Values[KEY_REALLY_NAME] = employee.ReallyName
	session.Values[KEY_DEPARTMENT_ID] = employee.DepartmentId
	session.Values[KEY_ROLE_ID] = employee.RoleId
	session.Values[KEY_ROLE_CODE] = employee.RoleCode
	session.Values[KEY_ROLE_LEVEL] = employee.RoleLevel

	//session.Options.MaxAge = 10 * 24 * 60 * 60
	session.Save(r, w)
}

package admin

import (
	"encoding/json"
	"fmt"
	"github.com/alanwhen/education-mini/enums"
	"github.com/alanwhen/education-mini/models"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

type StudentController struct {
	BaseController
}

func (this *StudentController) Prepare() {
	this.BaseController.Prepare()

	this.checkLogin()
}

func (this *StudentController) Index() {
	this.Data["pageTitle"] = "学籍列表"
	this.Data["showMoreQuery"] = true
	this.Data["activeSidebarUrl"] = this.URLFor(this.controllerName + "." + this.actionName)

	this.setTpl("admin/student/index.html")
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["header"] = "admin/student/index_header.html"
	this.LayoutSections["footer"] = "admin/student/index_footer.html"

	this.Data["canEdit"] = this.checkActionAuthor(this.controllerName, "Edit")
	this.Data["canDelete"] = this.checkActionAuthor(this.controllerName, "Delete")
}

func (this *StudentController) DataGrid() {
	var params models.StudentQueryParam
	json.Unmarshal(this.Ctx.Input.RequestBody, &params)

	data, total := models.StudentPageList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data

	this.Data["json"] = result
	this.ServeJSON()
}

func (this *StudentController) Edit() {
	if this.Ctx.Request.Method == "POST" {
		this.Save()
	}

	Id, _ := this.GetInt(":id", 0)
	m := models.Student{Id: Id}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			this.pageError("数据无效，请刷新后重试")
		}
	} else {
		//设置m的默认值
	}

	this.Data["m"] = m
	this.setTpl("admin/student/edit.html", "shared/layout_pull_box.html")
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["footer"] = "admin/student/edit_footer.html"
}

func (this *StudentController) Save() {
	var err error
	m := models.Student{}

	if err = this.ParseForm(&m); err != nil {
		this.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}

	o := orm.NewOrm()
	if m.Id == 0 {
		if err = o.Begin(); err != nil {
			this.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
			return
		}

		if _, err = o.Insert(&m); err == nil {
			if err = o.Commit(); err != nil {
				this.jsonResult(enums.JRCodeFailed, "添加提交失败", m.Id)
				o.Rollback()
			} else {
				this.jsonResult(enums.JRCodeSuccess, "添加成功", m.Id)
			}
		} else {
			if err = o.Rollback(); err != nil {
				this.jsonResult(enums.JRCodeFailed, "添加回滚失败", m.Id)
			} else {
				this.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
			}
		}
	} else {
		if _, err = o.Update(&m); err == nil {
			this.jsonResult(enums.JRCodeSuccess, "编辑成功", m.Id)
		} else {
			this.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}
}

func (this *StudentController) Delete() {
	strS := this.GetString("ids")
	ids := make([]int, 0, len(strS))
	for _, str := range strings.Split(strS, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}

	if num, err := models.StudentBatchDelete(ids); err == nil {
		this.jsonResult(enums.JRCodeSuccess, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		this.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}

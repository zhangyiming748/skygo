package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/log"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mysql_model"
)

type KnowledgeTestCaseModuleController struct{}

func (this KnowledgeTestCaseModuleController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	s := mysql.GetSession()

	// 查询组键
	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	all := widget.PaginatorFind(s, &[]mysql_model.TestCaseModule{})
	response.RenderSuccess(ctx, all)
}

func (this KnowledgeTestCaseModuleController) GetOne(ctx *gin.Context) {
	id := request.ParamString(ctx, "id")
	s := mysql.GetSession()
	s.Where("id=?", id)

	w := orm.PWidget{}
	result, err := w.One(s, &mysql_model.TestCaseModule{})

	if err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this KnowledgeTestCaseModuleController) Create(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)

	module := new(mysql_model.TestCaseModule)
	module.ModuleName = req.MustString("module_name")
	module.ModuleNameCode = req.MustString("module_name_code")
	module.ModuleType = req.MustString("module_type")
	module.ModuleTypeCode = req.MustString("module_type_code")
	if _, err := module.Create(); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, orm.StructToMap(*module))
	}
}

func (this KnowledgeTestCaseModuleController) Update(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	id := request.ParamString(ctx, "id")
	module := new(mysql_model.TestCaseModule)

	if request.IsExist(ctx, "module_name") {
		module.ModuleName = req.MustString("module_name")
	}
	if request.IsExist(ctx, "module_name_code") {
		module.ModuleNameCode = req.MustString("module_name_code")
	}
	if request.IsExist(ctx, "module_type") {
		module.ModuleType = req.MustString("module_type")
	}
	if request.IsExist(ctx, "module_type_code") {
		module.ModuleTypeCode = req.MustString("module_type_code")
	}
	if _, err := module.Update(id); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, orm.StructToMap(*module))
	}
}

func (this KnowledgeTestCaseModuleController) BulkDelete(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	successNum := 0
	if _, has := req.TrySlice("ids"); has {
		ids := req.SliceInt("ids")
		for _, id := range ids {
			_, err := new(mysql_model.TestCaseModule).RemoveById(id)
			if err != nil {
				log.GetHttpLogLogger().Error(fmt.Sprintf("%v", err))
				break
			} else {
				successNum++
			}
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"number": successNum})
}

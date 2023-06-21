package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mysql_model"
)

type Screen struct{}

func (this Screen) GetScreenInfo(ctx *gin.Context) {
	session := mysql.GetSession()
	widget := orm.PWidget{}
	widget.SetQueryStr(ctx.Request.URL.RawQuery)
	all, err := widget.All(session, &[]mysql_model.ScreenInfo{})
	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, all)
	}
}

func (this Screen) CreateScreenInfo(ctx *gin.Context) {
	body := new(mysql_model.ScreenInfo)
	if content, err := json.Marshal(request.GetRequestBody(ctx)); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		if err := json.Unmarshal(content, body); err != nil {
			response.RenderFailure(ctx, err)
		} else {
			body.Create()
			response.RenderSuccess(ctx, body)
		}
	}
}

func (this Screen) DeleteScreenInfo(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	ids := req.SliceInt("ids")
	var successNum int
	var failNum int
	for _, id := range ids {
		err := new(mysql_model.ScreenInfo).Remove(id)
		if err != nil {
			failNum++
		} else {
			successNum++
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"success": successNum, "fail": failNum})
}

func (this Screen) GetScreenPieceTestProgress(ctx *gin.Context) {
	session := mysql.GetSession()
	widget := orm.PWidget{}
	widget.SetQueryStr(ctx.Request.URL.RawQuery)
	all, err := widget.All(session, &[]mysql_model.ScreenPieceTestProgress{})
	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, all)
	}
}

func (this Screen) CreateScreenPieceTestProgress(ctx *gin.Context) {
	body := new(mysql_model.ScreenPieceTestProgress)
	if content, err := json.Marshal(request.GetRequestBody(ctx)); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		if err := json.Unmarshal(content, body); err != nil {
			response.RenderFailure(ctx, err)
		} else {
			body.Create()
			response.RenderSuccess(ctx, body)
		}
	}
}

func (this Screen) DeleteScreenPieceTestProgress(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	ids := req.SliceInt("ids")
	var successNum int
	var failNum int
	for _, id := range ids {
		err := new(mysql_model.ScreenPieceTestProgress).Remove(id)
		if err != nil {
			failNum++
		} else {
			successNum++
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"success": successNum, "fail": failNum})
}

func (this Screen) GetScreenTaskInfo(ctx *gin.Context) {
	session := mysql.GetSession()
	widget := orm.PWidget{}
	widget.SetQueryStr(ctx.Request.URL.RawQuery)
	all, err := widget.All(session, &[]mysql_model.ScreenTaskInfo{})
	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, all)
	}
}

func (this Screen) CreateScreenTaskInfo(ctx *gin.Context) {
	body := new(mysql_model.ScreenTaskInfo)
	if content, err := json.Marshal(request.GetRequestBody(ctx)); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		if err := json.Unmarshal(content, body); err != nil {
			response.RenderFailure(ctx, err)
		} else {
			body.Create()
			response.RenderSuccess(ctx, body)
		}
	}
}

func (this Screen) DeleteScreenTaskInfo(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	ids := req.SliceInt("ids")
	var successNum int
	var failNum int
	for _, id := range ids {
		err := new(mysql_model.ScreenTaskInfo).Remove(id)
		if err != nil {
			failNum++
		} else {
			successNum++
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"success": successNum, "fail": failNum})
}

func (this Screen) GetScreenTestCase(ctx *gin.Context) {
	session := mysql.GetSession()
	widget := orm.PWidget{}
	widget.SetQueryStr(ctx.Request.URL.RawQuery)
	all, err := widget.All(session, &[]mysql_model.ScreenTestCase{})
	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, all)
	}
}

func (this Screen) CreateScreenTestCase(ctx *gin.Context) {
	body := new(mysql_model.ScreenTestCase)
	if content, err := json.Marshal(request.GetRequestBody(ctx)); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		if err := json.Unmarshal(content, body); err != nil {
			response.RenderFailure(ctx, err)
		} else {
			body.Create()
			response.RenderSuccess(ctx, body)
		}
	}
}

func (this Screen) DeleteScreenTestCase(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	ids := req.SliceInt("ids")
	var successNum int
	var failNum int
	for _, id := range ids {
		err := new(mysql_model.ScreenTestCase).Remove(id)
		if err != nil {
			failNum++
		} else {
			successNum++
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"success": successNum, "fail": failNum})
}

func (this Screen) GetScreenVehicleInfo(ctx *gin.Context) {
	session := mysql.GetSession()
	widget := orm.PWidget{}
	widget.SetQueryStr(ctx.Request.URL.RawQuery)
	all, err := widget.All(session, &[]mysql_model.ScreenVehicleInfo{})
	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, all)
	}
}

func (this Screen) CreateScreenVehicleInfo(ctx *gin.Context) {
	body := new(mysql_model.ScreenVehicleInfo)
	if content, err := json.Marshal(request.GetRequestBody(ctx)); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		if err := json.Unmarshal(content, body); err != nil {
			response.RenderFailure(ctx, err)
		} else {
			body.Create()
			response.RenderSuccess(ctx, body)
		}
	}
}

func (this Screen) DeleteScreenVehicleInfo(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	ids := req.SliceInt("ids")
	var successNum int
	var failNum int
	for _, id := range ids {
		err := new(mysql_model.ScreenVehicleInfo).Remove(id)
		if err != nil {
			failNum++
		} else {
			successNum++
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"success": successNum, "fail": failNum})
}

// 整车测试进展
func (this Screen) GetScreenVehicleTestProgress(ctx *gin.Context) {
	session := mysql.GetSession()
	widget := orm.PWidget{}
	widget.SetQueryStr(ctx.Request.URL.RawQuery)
	all, err := widget.All(session, &[]mysql_model.ScreenVehicleTestProgress{})
	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, all)
	}
}

func (this Screen) CreateScreenVehicleTestProgress(ctx *gin.Context) {
	body := new(mysql_model.ScreenVehicleTestProgress)
	if content, err := json.Marshal(request.GetRequestBody(ctx)); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		if err := json.Unmarshal(content, body); err != nil {
			response.RenderFailure(ctx, err)
		} else {
			body.Create()
			response.RenderSuccess(ctx, body)
		}
	}
}

func (this Screen) DeleteScreenVehicleTestProgress(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	ids := req.SliceInt("ids")
	var successNum int
	var failNum int
	for _, id := range ids {
		err := new(mysql_model.ScreenVehicleTestProgress).Remove(id)
		if err != nil {
			failNum++
		} else {
			successNum++
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"success": successNum, "fail": failNum})
}

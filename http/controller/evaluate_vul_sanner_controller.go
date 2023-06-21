package controller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/license"
	"skygo_detection/mongo_model"
	"skygo_detection/mysql_model"
	"skygo_detection/service"
)

type EvaluateVulScannerController struct{}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_vul_scanners  下载漏洞检测用例
 * @apiVersion 0.1.0
 * @apiName GetAll
 * @apiGroup EvaluateVulScanner
 *
 * @apiDescription 下载漏洞检测用例
 *
 * @apiUse authHeader
 *
 * @apiParam {string}  name    漏洞检测任务名称
 *
 * @apiParamExample {json}  请求参数示例:
 *  {"name":"漏洞检测任务名称"}
 *
 * @apiSuccessExample {json} 请求成功示例:
 *{
 *    "code": 0,
 *    "data": {
 *        "url": "http://pub-zzdt.s3.360.cn/skaderb8c6"
 *    }
 *}
 */
func (this EvaluateVulScannerController) GetAll(ctx *gin.Context) {
	// 根据参数下载对应的测试用例包 (platform)_(cpu_version)_(sys_sdk_ver).zip
	platform := strings.ToLower(request.QueryString(ctx, "platform"))
	cpuVersion := strings.ToLower(request.QueryString(ctx, "cpu_version"))
	sysSdkVer := strings.ToLower(request.QueryString(ctx, "sys_sdk_ver"))
	fileName := fmt.Sprintf("%s_%s_%s.zip", platform, cpuVersion, sysSdkVer)
	packagePath := fmt.Sprintf("%s/vul_scanner/%s", service.LoadConfig().Download.DownloadPath, fileName)
	if content, err := ioutil.ReadFile(packagePath); err == nil {
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		ctx.Header("Content-Type", "*")
		ctx.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
		ctx.Writer.Write(content)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_vul_scanners  漏洞检测结果上传
 * @apiVersion 0.1.0
 * @apiName Create
 * @apiGroup EvaluateVulScanner
 *
 * @apiDescription 漏洞检测结果上传
 *
 * @apiUse authHeader
 *
 * @apiParam {string}  task_id    任务id
 *
 * @apiParamExample {json}  请求参数示例:
 *
 *{
 *    "task_id": "123456",
 *    "deviceInfos": {
 *        "company": "LGE",
 *        "sysVersion": "7.1.1",
 *        "cpuMode": "AArch64 Processor rev 3 (aarch64)",
 *        "cpuVersion":64,
 *        "platform":"Android",
 *        "sysSdkVer":25,
 *        "brand":"Android",
 *        "carMode":"AOSP on BullHead"
 *    },
 *    "scanResult": [
 *        {
 *            "cveId": "CVE-2014-9322",
 *            "googleSeverityLevel": 1,
 *            "dateExposure": "2014-04-30 07:46:09",
 *            "dateBulletin": "2016-04-02",
 *            "sketch": "内核中的提权漏洞",
 *            "description": "Linux kernel是美国Linux基金会发布的开源操作系统Linux所使用的内核。NFSv4 implementation是其中的一个分布式文件系统协议。Linux kernel 3.17.5之前版本的arch\/x86\/kernel\/entry_64.S文件中存在安全漏洞，该漏洞源于程序没有正确处理与Stack Segment(SS)段寄存器相关的错误。本地攻击者可借助IRET指令利用该漏洞获取权限。"
 *        }
 *    ]
 *}
 */
func (this EvaluateVulScannerController) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	if result, err := new(mongo_model.EvaluateVulScanner).Create(*req); err == nil {
		taskId := req.MustString("task_id")
		body := qmap.QM{
			"status": common.VUL_PRELIMINARY_END,
		}
		new(mysql_model.VulTask).UpdateByTaskId(taskId, body)
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_vul_scanners/:id  获取单条漏洞检测结果
 * @apiVersion 0.1.0
 * @apiName GetOne
 * @apiGroup EvaluateVulScanner
 *
 * @apiDescription 获取单条漏洞检测结果
 *
 * @apiUse authHeader
 *
 * @apiParam {string}  id    漏洞ID
 *
 * @apiParamExample   请求参数示例:
 *  curl localhost:3001/api/v1/evaluate_vul_scanners/5fdc7894aee3d1849a612bc4
 *
 * @apiSuccessExample {json} 请求成功示例:
 *{
 *    "code": 0,
 *    "data": {
 *        "_id": "5fdc7894aee3d1849a612bc4",
 *        "deviceInfos": {
 *            "brand": "Android",
 *            "carMode": "AOSP on BullHead",
 *            "company": "LGE",
 *            "cpuMode": "AArch64 Processor rev 3 (aarch64)",
 *            "cpuVersion": 64,
 *            "platform": "Android",
 *            "sysSdkVer": 25,
 *            "sysVersion": "7.1.1"
 *        },
 *        "scanResult": [
 *            {
 *                "cveId": "CVE-2014-9322",
 *                "dateBulletin": "2016-04-02",
 *                "dateExposure": "2014-04-30 07:46:09",
 *                "description": "Linux kernel是美国Linux基金会发布的开源操作系统Linux所使用的内核。NFSv4 implementation是其中的一个分布式文件系统协议。Linux kernel 3.17.5之前版本的arch/x86/kernel/entry_64.S文件中存在安全漏洞，该漏洞源于程序没有正确处理与Stack Segment(SS)段寄存器相关的错误。本地攻击者可借助IRET指令利用该漏洞获取权限。",
 *                "googleSeverityLevel": 1,
 *                "sketch": "内核中的提权漏洞"
 *            }
 *        ],
 *        "taskId": "123456"
 *    }
 *}
 */
func (this EvaluateVulScannerController) GetOne(ctx *gin.Context) {
	id := ctx.Param("id")

	params := qmap.QM{
		"e__id": bson.ObjectIdHex(id),
	}
	ormSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_SCANNER, params)
	data, _ := ormSession.GetOne()
	response.RenderSuccess(ctx, data)
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_vul_scanner/distribution/:id  获取漏洞分布
 * @apiVersion 0.1.0
 * @apiName Distribution
 * @apiGroup EvaluateVulScanner
 *
 * @apiDescription 获取漏洞分布
 *
 * @apiUse authHeader
 *
 * @apiParam {string}  id    任务ID
 *
 * @apiParamExample   请求参数示例:
 *  curl localhost:3001/api/v1/evaluate_vul_scanner/distribution/123456
 *
 * @apiSuccessExample {json} 请求成功示例:
 *{
 *    "code": 0,
 *    "data": {
 *        "level_distribution": [
 *            {
 *                "name": "低危",
 *                "value": 10
 *            },
 *            ...
 *        ],
 *       "module_distribution": [
 *            {
 *                "name": "锁定屏幕",
 *                "value": 10
 *            },
 *            ...
 *        ],
 *        "type_distribution": [
 *            {
 *                "name": "越权",
 *                "value": 10
 *            },
 *            ...
 *        ]
 *    }
 *}
 */
func (this EvaluateVulScannerController) Distribution(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	// db.getCollection('evaluate_vul_info').aggregate([{$match:{"task_id":"123456"}},{$group:{_id:"$cve_type", num_tutorial:{$sum:  1}}}])
	taskId := ctx.Param("id")
	distribution := qmap.QM{}
	// 拼接聚合语句
	match := bson.M{"task_id": taskId}
	session := mongo.NewMgoSession(common.MC_EVALUATE_VUL_INFO).Session
	total, _ := session.Find(match).Count()
	// 获取影响模块分布 module_distribution
	{
		moduleDistribution := &[]qmap.QM{}
		group := bson.M{"_id": "$involve_module", "num": bson.M{"$sum": 1}}
		agg := []bson.M{{"$match": match}, {"$group": group}, {"$sort": bson.M{"num": -1}}}
		result := []qmap.QM{}
		pipe := session.Pipe(agg)
		pipe.Iter().All(&result)
		moduleDistribution = this.distribution(result, total)
		distribution["module_distribution"] = moduleDistribution
	}
	// 获取漏洞等级分布 level_distribution
	{
		levelDistribution := &[]qmap.QM{}
		group := bson.M{"_id": "$google_severity_level", "num": bson.M{"$sum": 1}}
		agg := []bson.M{{"$match": match}, {"$group": group}, {"$sort": bson.M{"num": -1}}}
		result := []qmap.QM{}
		pipe := session.Pipe(agg)
		pipe.Iter().All(&result)
		levelDistribution = this.distributionLevel(result, total)
		distribution["level_distribution"] = levelDistribution
	}
	// 获取漏洞漏洞分布 type_distribution
	{
		typeDistribution := &[]qmap.QM{}
		group := bson.M{"_id": "$cve_type", "num": bson.M{"$sum": 1}}
		agg := []bson.M{{"$match": match}, {"$group": group}}
		result := []qmap.QM{}
		pipe := session.Pipe(agg)
		pipe.Iter().All(&result)
		typeDistribution = this.distribution(result, total)
		distribution["type_distribution"] = typeDistribution
	}
	response.RenderSuccess(ctx, distribution)
}

func (this EvaluateVulScannerController) distribution(input []qmap.QM, total int) *[]qmap.QM {
	result := []qmap.QM{}
	// 只去前八个数据做展示
	if len(input) > 8 {
		sum := 0
		countSum := 0
		for i := 0; i < 7; i++ {
			tmp := input[i]
			name := tmp.MustString("_id")
			count := tmp.MustInt("num")
			value := count * 100 / total
			result = append(result, qmap.QM{"name": name, "value": value, "count": count})
			sum += value
			countSum += count
		}
		result = append(result, qmap.QM{"name": "others", "value": 100 - sum, "count": countSum * (100 - sum) / 100})
	} else {
		for _, tmp := range input {
			name := tmp.MustString("_id")
			value := tmp.MustInt("num") * 100 / total
			count := tmp.MustInt("num")
			result = append(result, qmap.QM{"name": name, "value": value, "count": count})
		}
	}
	return &result
}

func (this EvaluateVulScannerController) distributionLevel(input []qmap.QM, total int) *[]qmap.QM {
	result := []qmap.QM{
		qmap.QM{"name": "严重", "value": 0},
		qmap.QM{"name": "高危", "value": 0},
		qmap.QM{"name": "中危", "value": 0},
		// qmap.QM{"name": "低危", "value": 0},
	}
	for _, tmp := range input {
		name := tmp.MustInt("_id")
		count := tmp.MustInt("num")
		value := count * 100 / total
		level := qmap.QM{}
		switch name {
		case 3:
			level = qmap.QM{"name": "严重", "value": value, "count": count}
			result[0] = level
		case 2:
			level = qmap.QM{"name": "高危", "value": value, "count": count}
			result[1] = level
		case 1:
			level = qmap.QM{"name": "中危", "value": value, "count": count}
			result[2] = level
		}
	}
	return &result
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_vul_scanner/vul_numbers/:id  获取漏洞数量
 * @apiVersion 0.1.0
 * @apiName VulNumbers
 * @apiGroup EvaluateVulScanner
 *
 * @apiDescription 获取漏洞数量
 *
 * @apiUse authHeader
 *
 * @apiParam {string}  id    漏洞ID
 *
 * @apiParamExample   请求参数示例:
 *  curl localhost:3001/api/v1/evaluate_vul_scanner/vul_numbers/5fdc7894aee3d1849a612bc4
 *
 * @apiSuccessExample {json} 请求成功示例:
 *{
 *    "code": 0,
 *    "data": {
 *        "vul_all": {
 *            "level_number": [
 *                {
 *                    "name": "全部",
 *                    "number": 100
 *                },
 *                ...
 *            ],
 *            "number": 100
 *        },
 *        "vul_other": {
 *            "level_number": [
 *                {
 *                    "name": "全部",
 *                    "number": 100
 *                },
 *                ...
 *            ],
 *            "number": 100
 *        },
 *       "vul_repaired": {
 *            "level_number": [
 *                {
 *                    "name": "全部",
 *                    "number": 100
 *                },
 *                ...
 *            ],
 *            "number": 100
 *        },
 *        "vul_unrepair": {
 *            "level_number": [
 *                {
 *                    "name": "全部",
 *                    "number": 100
 *                },
 *                ...
 *            ],
 *            "number": 100
 *        }
 *    }
 *}
 */
func (this EvaluateVulScannerController) VulNumbers(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")

	taskId := req.MustString("id")
	vulNumbers := qmap.QM{}
	// 拼接聚合语句
	session := mongo.NewMgoSession(common.MC_EVALUATE_VUL_INFO).Session
	// 全部漏洞
	{
		match := bson.M{"task_id": taskId}
		total, _ := session.Find(match).Count()
		match["google_severity_level"] = common.VUL_GOOGLE_LEVEL_SERIOUS
		serious, _ := session.Find(match).Count()
		match["google_severity_level"] = common.VUL_GOOGLE_LEVEL_HIGHT
		hight, _ := session.Find(match).Count()
		match["google_severity_level"] = common.VUL_GOOGLE_LEVEL_MIDDLE
		middlle, _ := session.Find(match).Count()
		levelResult := []qmap.QM{
			qmap.QM{"name": "全部", "number": total},
			qmap.QM{"name": "严重", "number": serious},
			qmap.QM{"name": "高危", "number": hight},
			qmap.QM{"name": "中危", "number": middlle},
		}
		result := qmap.QM{}
		result["number"] = total
		result["level_number"] = levelResult
		vulNumbers["vul_all"] = result
	}
	// 未修复的漏洞
	{
		match := bson.M{"task_id": taskId, "fix_status": common.VUL_FIX_UNREPAIR}
		total, _ := session.Find(match).Count()
		match["google_severity_level"] = common.VUL_GOOGLE_LEVEL_SERIOUS
		serious, _ := session.Find(match).Count()
		match["google_severity_level"] = common.VUL_GOOGLE_LEVEL_HIGHT
		hight, _ := session.Find(match).Count()
		match["google_severity_level"] = common.VUL_GOOGLE_LEVEL_MIDDLE
		middlle, _ := session.Find(match).Count()
		levelResult := []qmap.QM{
			qmap.QM{"name": "全部", "number": total},
			qmap.QM{"name": "严重", "number": serious},
			qmap.QM{"name": "高危", "number": hight},
			qmap.QM{"name": "中危", "number": middlle},
		}
		result := qmap.QM{}
		result["number"] = total
		result["level_number"] = levelResult
		vulNumbers["vul_unrepair"] = result
	}
	// 已修复的漏洞
	{
		match := bson.M{"task_id": taskId, "fix_status": common.VUL_FIX_REPAIR}
		total, _ := session.Find(match).Count()
		match["google_severity_level"] = common.VUL_GOOGLE_LEVEL_SERIOUS
		serious, _ := session.Find(match).Count()
		match["google_severity_level"] = common.VUL_GOOGLE_LEVEL_HIGHT
		hight, _ := session.Find(match).Count()
		match["google_severity_level"] = common.VUL_GOOGLE_LEVEL_MIDDLE
		middlle, _ := session.Find(match).Count()
		levelResult := []qmap.QM{
			qmap.QM{"name": "全部", "number": total},
			qmap.QM{"name": "严重", "number": serious},
			qmap.QM{"name": "高危", "number": hight},
			qmap.QM{"name": "中危", "number": middlle},
		}
		result := qmap.QM{}
		result["number"] = total
		result["level_number"] = levelResult
		vulNumbers["vul_repaired"] = result
	}
	// 未涉及漏洞
	{
		match := bson.M{"task_id": taskId, "fix_status": common.VUL_FIX_OTHER}
		total, _ := session.Find(match).Count()
		match["google_severity_level"] = common.VUL_GOOGLE_LEVEL_SERIOUS
		serious, _ := session.Find(match).Count()
		match["google_severity_level"] = common.VUL_GOOGLE_LEVEL_HIGHT
		hight, _ := session.Find(match).Count()
		match["google_severity_level"] = common.VUL_GOOGLE_LEVEL_MIDDLE
		middlle, _ := session.Find(match).Count()
		levelResult := []qmap.QM{
			qmap.QM{"name": "全部", "number": total},
			qmap.QM{"name": "严重", "number": serious},
			qmap.QM{"name": "高危", "number": hight},
			qmap.QM{"name": "中危", "number": middlle},
		}
		result := qmap.QM{}
		result["number"] = total
		result["level_number"] = levelResult
		vulNumbers["vul_other"] = result
	}
	response.RenderSuccess(ctx, vulNumbers)
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_vul_scanner/sys_info/:id  获取系统信息
 * @apiVersion 0.1.0
 * @apiName GetSysInfo
 * @apiGroup EvaluateVulScanner
 *
 * @apiDescription 获取系统信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}  id    漏洞ID
 *
 * @apiParamExample   请求参数示例:
 *  curl localhost:3001/api/v1/evaluate_vul_scanner/sys_info/5fdc7894aee3d1849a612bc4
 *
 * @apiSuccessExample {json} 请求成功示例:
 *{
 *    "code": 0,
 *    "data": {
 *        "brand": "Android",
 *        "carMode": "AOSP on BullHead",
 *        "company": "LGE",
 *        "cpuMode": "AArch64 Processor rev 3 (aarch64)",
 *        "cpuVersion": 64,
 *        "platform": "Android",
 *        "sysSdkVer": 25,
 *        "sysVersion": "7.1.1"
 *    }
 *}
 */
func (this EvaluateVulScannerController) GetSysInfo(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")
	(*req)["query_params"] = ctx.Request.URL.RawQuery

	params := qmap.QM{
		"e_task_id": req.MustString("id"),
	}
	ormSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_DEVICE_INFO, params)
	if res, err := ormSession.GetOne(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderSuccess(ctx, qmap.QM{})
	}
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_vul_scanner/vul_info/:id  获取漏洞信息
 * @apiVersion 0.1.0
 * @apiName GetVulInfo
 * @apiGroup EvaluateVulScanner
 *
 * @apiDescription 获取漏洞信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}  id    任务ID
 *
 * @apiParamExample   请求参数示例:
 *  curl localhost:3001/api/v1/evaluate_vul_scanner/vul_info/123456
 *
 * @apiSuccessExample {json} 请求成功示例:
 *{
 *    "code": 0,
 *    "data": {
 *        "list": [
 *            {
 *                "_id": "5fe057d7aee3d1849a6a025b",
 *                "cve_id": "CVE-2016-0844",
 *                "cve_type": "Eop",
 *                "date_bulletin": "2016-04-02",
 *                "date_exposure": "2016-01-06 05:16:34",
 *                "description": "Android是美国谷歌（Google）公司和开放手持设备联盟（简称OHA）共同开发的一套以Linux为基础的开源操作系统。Qualcomm RF是使用在其中的一个美国高通（Qualcomm）公司开发的前端解决方案（包含功率放大器、天线开关、天线调谐器以及包络跟踪器等一系列芯片技术）组件。Android的Qualcomm RF组件中存在提权漏洞，该漏洞源于程序没有正确限制使用套接字ioctl调用。本地攻击者可借助特制的应用程序利用该漏洞获取权限，在内核上下文中执行任意代码。以下版本受到影响：Android 4.4.4之前版本，5.0.2之前版本，5.1.1之前版本，6.0之前版本和6.0.1之前版本。",
 *                "fix_status": 4,
 *                "google_severity_level": 1,
 *                "involve_module": "Qualcomm RF 组件",
 *                "sketch": "Qualcomm RF 组件中的提权漏洞",
 *                "task_id": "123456"
 *            },
 *            ...
 *        ],
 *        "pagination": {
 *            "count": 2,
 *            "current_page": 1,
 *            "per_page": 20,
 *            "total": 2,
 *            "total_pages": 1
 *        }
 *    }
 *}
 */
func (this EvaluateVulScannerController) GetVulInfo(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")
	(*req)["query_params"] = ctx.Request.URL.RawQuery

	queryParams := req.String("query_params")
	params := qmap.QM{
		"e_task_id": req.MustString("id"),
	}
	ormSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_INFO, params).AddUrlQueryCondition(queryParams)
	ormSession.AddSorter("fix_status", 0)
	ormSession.AddSorter("google_severity_level", 1)
	ormSession.SetTransformFunc(vulScannerTransform)
	if res, err := ormSession.GetPage(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func vulScannerTransform(data qmap.QM) qmap.QM {
	id := data["_id"]
	searchContent := fmt.Sprintf("%s_%s", data.String("cve_id"), data.String("sketch"))
	// todo 临时脚本，刷search_content内容
	selector := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{"search_content": searchContent},
	}
	if err := mongo.NewMgoSession(common.MC_EVALUATE_VUL_INFO).Update(selector, update); err != nil {
		panic(err)
	}
	data["search_content"] = searchContent
	data["relate_task_case"] = new(mysql_model.VulnerabilityScannerVulRelation).GetRelationVulIds(id.(bson.ObjectId).Hex())
	return data
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_vul_scanner/check_auth  漏洞检测授权校验
 * @apiVersion 0.1.0
 * @apiName CheckAuth
 * @apiGroup EvaluateVulScanner
 *
 * @apiDescription 漏洞检测授权校验
 *
 * @apiUse authHeader
 *
 * @apiParam {string}  task_id    任务id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *    "task_id": "123456"
 * }
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *    "code": 0,
 *    "data": {},
 *	  "msg":""
 * }
 */
func (this EvaluateVulScannerController) CheckAuth(ctx *gin.Context) {
	if result := license.VerifyMenu(license.TEST_TASK); !result {
		response.RenderFailure(ctx, errors.New("授权证书无效，授权失败"))
		return
	}
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	params := qmap.QM{
		"e_task_id": req.String("task_id"),
	}
	vulTask := new(mysql_model.VulTask)
	if has, _ := sys_service.NewSessionWithCond(params).GetOne(vulTask); has {
		response.RenderSuccess(ctx, qmap.QM{})
	} else {
		response.RenderFailure(ctx, errors.New("请先在平台创建漏洞检测任务"))
	}
}

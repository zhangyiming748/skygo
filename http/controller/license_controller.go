package controller

import (
	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/app/http/request"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/license"
)

type LicenseController struct{}

/**
 * apiType http
 * @api {get} /api/v1/license/info 查询授权信息
 * @apiVersion 1.0.0
 * @apiName GetLicense
 * @apiGroup License
 *
 * @apiDescription 查询授权信息
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "active_time": 1635931733,
 *         "create_time": 1635931932,
 *         "disk_uuid": "aeb6fc55-7fb2-4a6b-aed8-3dff04c2766e",
 *         "expired_time": 1636018133,
 *		   "release_version": "1.0",
 *		   "build_version":   "202111",
 *         "menu": [
 *             {
 *                 "expired_time": 1636018133,
 *                 "id": 1,
 *                 "name": "测试任务"
 *             },
 *             {
 *                 "expired_time": -1,
 *                 "id": 2,
 *                 "name": "报告管理"
 *             },
 *             {
 *                 "expired_time": -1,
 *                 "id": 3,
 *                 "name": "漏洞管理"
 *             },
 *             {
 *                 "expired_time": -1,
 *                 "id": 0,
 *                 "name": "资产管理"
 *             }
 *         ]
 *     },
 *     "msg": ""
 * }
 */
func (this *LicenseController) GetLicense(ctx *gin.Context) {
	response.RenderSuccess(ctx, license.GetLicense())
}

/**
 * apiType http
 * @api {post} /api/v1/license/generate 生成授权码
 * @apiVersion 1.0.0
 * @apiName GenerateLicense
 * @apiGroup License
 *
 * @apiDescription 生成授权码
 *
 * @apiUse authHeader
 *
 * @apiParam {int}           active_time                许可证激活时间(秒时间戳)
 * @apiParam {int}           authorize_time             授权时长(例如:授权1天，该值为86400秒)
 * @apiParam {string}        disk_uuid	              	设备磁盘id
 *
 * @apiParamExample {json} 请求参数示例:
 * {
 *     "license":"YUpUd2E3TGFqUmZMS0JnLURjdkY4T0h4Y1FjRWtVQ0c4Q05DMm1xZUN3RVBsN1FaZUFaWkZTTW9ZaEM1UHlqSXlJZzBrM0sxdnluRjNFODZjS1NBTVNiWjA5VEFtZ3RzTU8zU2p4X2VIVHlxQ1Z2aXlTTjlZSDhsaEg1WExmQU9YVjRva3pKME50TW93SkVsZExnQ2wyS1p5TDIyeHA3WUswSDBXQmUycXJKM01ELVVjRGh5MFdrYmtHN1hMM3dFT0dRYm5fZXl5LWR6cjY2TVR4ekVHV2IzeURiaWtEQl9FNll3YTRUUkV2RDk4ckp0N1NzTFRLNE5tTE5tYUVCTDJ3ZTVZdV9saEUwM0NPRW5fZUxvcjBTWG5LSDNWYXdaekNka05FbHVxZi1mRHB5OEtIZmtFdGJxSTdwSldLVnpnLU1HSGpTNGhDX2p5RUJQV2lmazY1U1pmRGQ3ZWg0SEh2ZEF2dk5mMzg3UTk1anFxUEgyalpHYzZCMzhrY2dIdUk2NEJNSUc5d3dNcFZ1T2UzTEtGOHJwMWY0ejdNREpSSnNkVklENm5QOWwwYU40dnhVWmhZcnl0dzhvbGhtUDE5TUx2ckI3dF9sOUFwNm9HSkJQWXcudFBoZEpWZXpvZm1vNTRucVpqQlpBV3J6clV6c1plakU1Q29WMnJ3NkNFVQ"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "msg": ""
 * }
 */
func (this *LicenseController) GenerateLicense(ctx *gin.Context) {
	body := request.GetRequestBody(ctx)
	activeTime := body.MustInt("active_time")
	authorizeTime := body.MustInt("authorize_time")
	diskUUID := body.MustString("disk_uuid")
	// diskUUID, err := license.GetDiskUUID()
	// if err != nil {
	// 	response.RenderFailure(ctx, err)
	// }
	if licenseStr, err := license.GenerateLicence(activeTime, activeTime+authorizeTime, diskUUID); err == nil {
		response.RenderSuccess(ctx, qmap.QM{"license": licenseStr})
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/license/import 导入授权码
 * @apiVersion 1.0.0
 * @apiName ImportLicense
 * @apiGroup License
 *
 * @apiDescription 导入授权码
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           license 		许可证
 *
 * @apiParamExample {json} 请求参数示例:
 * {
 *     "license":YUpUd2E3TGFqUmZMS0JnLURjdkY4T0h4Y1FjRWtVQ0c4Q05DMm1xZUN3RVBsN1FaZUFaWkZTTW9ZaEM1UHlqSXlJZzBrM0sxdnluRjNFODZjS1NBTVNiWjA5VEFtZ3RzTU8zU2p4X2VIVHlxQ1Z2aXlTTjlZSDhsaEg1WExmQU9YVjRva3pKME50TW93SkVsZExnQ2wyS1p5TDIyeHA3WUswSDBXQmUycXJKM01ELVVjRGh5MFdrYmtHN1hMM3dFT0dRYm5fZXl5LWR6cjY2TVR4ekVHV2IzeURiaWtEQl9FNll3YTRUUkV2RDk4ckp0N1NzTFRLNE5tTE5tYUVCTDJ3ZTVZdV9saEUwM0NPRW5fZUxvcjBTWG5LSDNWYXdaekNka05FbHVxZi1mRHB5OEtIZmtFdGJxSTdwSldLVnpnLU1HSGpTNGhDX2p5RUJQV2lmazY1U1pmRGQ3ZWg0SEh2ZEF2dk5mMzg3UTk1anFxUEgyalpHYzZCMzhrY2dIdUk2NEJNSUc5d3dNcFZ1T2UzTEtGOHJwMWY0ejdNREpSSnNkVklENm5QOWwwYU40dnhVWmhZcnl0dzhvbGhtUDE5TUx2ckI3dF9sOUFwNm9HSkJQWXcudFBoZEpWZXpvZm1vNTRucVpqQlpBV3J6clV6c1plakU1Q29WMnJ3NkNFVQ
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "active_time": 1635931733,
 *         "create_time": 1635931932,
 *         "disk_uuid": "aeb6fc55-7fb2-4a6b-aed8-3dff04c2766e",
 *         "expired_time": 1636018133,
 *		   "release_version": "1.0",
 *		   "build_version":   "202111",
 *         "menu": [
 *             {
 *                 "expired_time": 1636018133,
 *                 "id": 1,
 *                 "name": "测试任务"
 *             },
 *             {
 *                 "expired_time": -1,
 *                 "id": 2,
 *                 "name": "报告管理"
 *             },
 *             {
 *                 "expired_time": -1,
 *                 "id": 3,
 *                 "name": "漏洞管理"
 *             },
 *             {
 *                 "expired_time": -1,
 *                 "id": 0,
 *                 "name": "资产管理"
 *             }
 *         ]
 *     },
 *     "msg": ""
 * }
 */
func (this *LicenseController) ImportLicense(ctx *gin.Context) {
	body := request.GetRequestBody(ctx)
	licenseStr := body.MustString("license")
	if err := license.ImportLicense(licenseStr); err == nil {
		response.RenderSuccess(ctx, license.GetLicense())
	} else {
		response.RenderFailure(ctx, err)
	}
}

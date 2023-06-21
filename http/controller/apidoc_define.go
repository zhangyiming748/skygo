package controller

//通用的api注释定义
//其他代码注释中可能引入的通用注释模块
//在文档导出的时候会用到(如果缺失,会影响文档的正常导出)

/**
 * @apiDefine authHeader
 *
 * @apiHeader {string}              Authorization       认证口令，值为: Bearer {{authToken}}, {{authToken}}为用户认证成功返回的Token
 *
 */

/**
 * @apiDefine urlQueryParams
 *
 * @apiParam {int}      [page=1]        页码
 * @apiParam {int}      [limit=100]     页面大小，限制返回数据个数
 * @apiParam {int}      [offset=0]      分页时查询数据的偏移量
 * @apiParam {string}   [fields]        指定返回数据包含哪些域，用于减少数据返回量
 * @apiParam {json}     [sort]          排序方式，如sort={"id":-1}，默认按资源ID逆序排序
 * @apiParam {json}     [where]         where条件查询参数，如where={"vehicles":{"os_version":5.1.1}}
 * @apiParam {json}     [or]            或条件查询参数, 如or=[{"author":{"e":"Jake"}},{author:{"e":"Alex"}}]
 * @apiParam {json}     [in]            匹配条件查询参数, 如in={"author":["Jake","Billy"]}
 * @apiParam {json}     [having]        having条件查询参数，如having={"status":0,"is_root":1}
 *
 */

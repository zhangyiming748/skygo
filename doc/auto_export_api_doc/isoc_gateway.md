# 工业SOC服务接口 v1.0.0


# Role

## 添加系统角色

<p>添加系统角色</p>

	POST /isoc/v1/roles

请求头参数列表:

| 名称    | 类型      |必填       | 说明                          |
|---------|-----------|-----------|--------------------------------------|
| Authorization			| string			| yes  | <p>认证口令，值为: Bearer {{authToken}}, {{authToken}}为用户认证成功返回的Token</p>							|

请求体参数列表:

| 名称    | 类型      |必填       | 说明                          |
|---------|-----------|-----------|--------------------------------------|
| channel_id			| string			| yes  | <p>渠道号</p>  						|
| name			| string			| yes  | <p>角色名称</p>  						|
| parent_id			| int			| yes  | <p>父角色id</p>  						|


请求参数示例:
```json
{
    "channel_id":"O68468",
    "name":"超级管理员",
    "parent_id":0
}
```

请求成功示例:
```json
{
    "code": 0,
    "data": {
        "channel_id": "O68468",
        "id": 5,
        "name": "超级管理员",
        "parent_id": 1
    }
}
```
## 批量删除系统角色

<p>批量删除系统角色</p>

	DELETE /isoc/v1/roles

请求头参数列表:

| 名称    | 类型      |必填       | 说明                          |
|---------|-----------|-----------|--------------------------------------|
| Authorization			| string			| yes  | <p>认证口令，值为: Bearer {{authToken}}, {{authToken}}为用户认证成功返回的Token</p>							|

请求体参数列表:

| 名称    | 类型      |必填       | 说明                          |
|---------|-----------|-----------|--------------------------------------|
| ids			| []int			| yes  | <p>用户id</p>  						|


请求参数示例:
```json
{
    "ids":[10,11]
}
```

请求成功示例:
```json
     {
          "code": 0
			 "data":{
				"number":2
			}
     }
```
## 文件下载

<p>文件下载</p>

	GET /api/v1/role/file/download

请求头参数列表:

| 名称    | 类型      |必填       | 说明                          |
|---------|-----------|-----------|--------------------------------------|
| Authorization			| string			| yes  | <p>认证口令，值为: Bearer {{authToken}}, {{authToken}}为用户认证成功返回的Token</p>							|

请求体参数列表:

| 名称    | 类型      |必填       | 说明                          |
|---------|-----------|-----------|--------------------------------------|
| file_id			| string			| yes  | <p>文件id</p>  						|
| page			| int			| no  | <p>页码</p><br/>默认值:`1`  						|
| limit			| int			| no  | <p>页面大小，限制返回数据个数</p><br/>默认值:`100`  						|
| offset			| int			| no  | <p>分页时查询数据的偏移量</p><br/>默认值:`0`  						|
| fields			| string			| no  | <p>指定返回数据包含哪些域，用于减少数据返回量</p>  						|
| sort			| json			| no  | <p>排序方式，如sort={&quot;id&quot;:-1}，默认按资源ID逆序排序</p>  						|
| where			| json			| no  | <p>where条件查询参数，如where={&quot;vehicles&quot;:{&quot;os_version&quot;:5.1.1}}</p>  						|
| or			| json			| no  | <p>或条件查询参数, 如or=[{&quot;author&quot;:{&quot;e&quot;:&quot;Jake&quot;}},{author:{&quot;e&quot;:&quot;Alex&quot;}}]</p>  						|
| in			| json			| no  | <p>匹配条件查询参数, 如in={&quot;author&quot;:[&quot;Jake&quot;,&quot;Billy&quot;]}</p>  						|
| having			| json			| no  | <p>having条件查询参数，如having={&quot;status&quot;:0,&quot;is_root&quot;:1}</p>  						|

请求示例:
```curl
curl -i http://localhost/api/v1/role/file/download
```


## 查询所有系统角色列表

<p>查询所有系统角色列表</p>

	GET /isoc/v1/roles

请求头参数列表:

| 名称    | 类型      |必填       | 说明                          |
|---------|-----------|-----------|--------------------------------------|
| Authorization			| string			| yes  | <p>认证口令，值为: Bearer {{authToken}}, {{authToken}}为用户认证成功返回的Token</p>							|



请求成功示例:
```json
{
    "code": 0,
    "data": [
        {
            "channel_id": "Q00001",
            "id": 1,
            "name": "超级管理员",
            "parent_id": 0
        },
        {
            "channel_id": "T56205",
            "id": 2,
            "name": "亿咖通管理员1",
            "parent_id": 1
        }
    ]
}
```
## 查询某一系统角色信息

<p>查询某一系统角色信息</p>

	GET /isoc/v1/roles/:id

请求头参数列表:

| 名称    | 类型      |必填       | 说明                          |
|---------|-----------|-----------|--------------------------------------|
| Authorization			| string			| yes  | <p>认证口令，值为: Bearer {{authToken}}, {{authToken}}为用户认证成功返回的Token</p>							|

请求体参数列表:

| 名称    | 类型      |必填       | 说明                          |
|---------|-----------|-----------|--------------------------------------|
| id			| int			| yes  | <p>系统角色id</p>  						|



请求成功示例:
```json
{
    "code": 0,
    "data": {
        "channel_id": "T56205",
        "id": 2,
        "name": "管理员",
        "parent_id": 1
    }
}
```
## 更新系统角色

<p>更新系统角色</p>

	PUT /isoc/v1/roles/:id

请求头参数列表:

| 名称    | 类型      |必填       | 说明                          |
|---------|-----------|-----------|--------------------------------------|
| Authorization			| string			| yes  | <p>认证口令，值为: Bearer {{authToken}}, {{authToken}}为用户认证成功返回的Token</p>							|

请求体参数列表:

| 名称    | 类型      |必填       | 说明                          |
|---------|-----------|-----------|--------------------------------------|
| channel_id			| string			| yes  | <p>渠道号</p>  						|
| name			| string			| yes  | <p>角色名称</p>  						|
| parent_id			| int			| yes  | <p>父角色id</p>  						|


请求参数示例:
```json
{
    "channel_id":"O68468",
    "name":"超级管理员1",
    "parent_id":0
}
```

请求成功示例:
```json
{
    "code": 0,
    "data": {
         "id":1,
         "channel_id":"O68468",
         "name":"超级管理员1",
         "parent_id":0
    }
}
```
## 文件上传

<p>文件上传</p>

	POST /api/v1/role/file/upload

请求头参数列表:

| 名称    | 类型      |必填       | 说明                          |
|---------|-----------|-----------|--------------------------------------|
| Authorization			| string			| yes  | <p>认证口令，值为: Bearer {{authToken}}, {{authToken}}为用户认证成功返回的Token</p>							|

请求体参数列表:

| 名称    | 类型      |必填       | 说明                          |
|---------|-----------|-----------|--------------------------------------|
| file_name			| string			| no  | <p>文件名称</p>  						|
| file			| file			| yes  | <p>文件</p>  						|

请求示例:
```curl
curl -i http://localhost/api/v1/role/file/upload
```


请求成功示例:
```json
 {
		"code":0,
		"msg":"",
		"data":{
			"file_id":"a834qafmcxvadfq1123"
		}
 }
```


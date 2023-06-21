define({ "api": [
  {
    "type": "post",
    "url": "/isoc/v1/roles",
    "title": "添加系统角色",
    "version": "1.0.0",
    "name": "Create",
    "group": "Role",
    "description": "<p>添加系统角色</p>",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "string",
            "optional": false,
            "field": "channel_id",
            "description": "<p>渠道号</p>"
          },
          {
            "group": "Parameter",
            "type": "string",
            "optional": false,
            "field": "name",
            "description": "<p>角色名称</p>"
          },
          {
            "group": "Parameter",
            "type": "int",
            "optional": false,
            "field": "parent_id",
            "description": "<p>父角色id</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "请求参数示例:",
          "content": "{\n    \"channel_id\":\"O68468\",\n    \"name\":\"超级管理员\",\n    \"parent_id\":0\n}",
          "type": "json"
        }
      ]
    },
    "success": {
      "examples": [
        {
          "title": "请求成功示例:",
          "content": "{\n    \"code\": 0,\n    \"data\": {\n        \"channel_id\": \"O68468\",\n        \"id\": 5,\n        \"name\": \"超级管理员\",\n        \"parent_id\": 1\n    }\n}",
          "type": "json"
        }
      ]
    },
    "filename": "../../http/controller/role_controller.go",
    "groupTitle": "Role",
    "sampleRequest": [
      {
        "url": "http://industry.isoc.360.cn/isoc/v1/roles"
      }
    ],
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "type": "string",
            "optional": false,
            "field": "Authorization",
            "description": "<p>认证口令，值为: Bearer {{authToken}}, {{authToken}}为用户认证成功返回的Token</p>"
          }
        ]
      }
    }
  },
  {
    "type": "delete",
    "url": "/isoc/v1/roles",
    "title": "批量删除系统角色",
    "version": "1.0.0",
    "name": "DeleteBulk",
    "group": "Role",
    "description": "<p>批量删除系统角色</p>",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "[]int",
            "optional": false,
            "field": "ids",
            "description": "<p>用户id</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "请求参数示例:",
          "content": "{\n    \"ids\":[10,11]\n}",
          "type": "json"
        }
      ]
    },
    "success": {
      "examples": [
        {
          "title": "请求成功示例:",
          "content": "     {\n          \"code\": 0\n\t\t\t \"data\":{\n\t\t\t\t\"number\":2\n\t\t\t}\n     }",
          "type": "json"
        }
      ]
    },
    "filename": "../../http/controller/role_controller.go",
    "groupTitle": "Role",
    "sampleRequest": [
      {
        "url": "http://industry.isoc.360.cn/isoc/v1/roles"
      }
    ],
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "type": "string",
            "optional": false,
            "field": "Authorization",
            "description": "<p>认证口令，值为: Bearer {{authToken}}, {{authToken}}为用户认证成功返回的Token</p>"
          }
        ]
      }
    }
  },
  {
    "type": "get",
    "url": "/api/v1/role/file/download",
    "title": "文件下载",
    "version": "1.0.0",
    "name": "Download",
    "group": "Role",
    "description": "<p>文件下载</p>",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "string",
            "optional": false,
            "field": "file_id",
            "description": "<p>文件id</p>"
          },
          {
            "group": "Parameter",
            "type": "int",
            "optional": true,
            "field": "page",
            "defaultValue": "1",
            "description": "<p>页码</p>"
          },
          {
            "group": "Parameter",
            "type": "int",
            "optional": true,
            "field": "limit",
            "defaultValue": "100",
            "description": "<p>页面大小，限制返回数据个数</p>"
          },
          {
            "group": "Parameter",
            "type": "int",
            "optional": true,
            "field": "offset",
            "defaultValue": "0",
            "description": "<p>分页时查询数据的偏移量</p>"
          },
          {
            "group": "Parameter",
            "type": "string",
            "optional": true,
            "field": "fields",
            "description": "<p>指定返回数据包含哪些域，用于减少数据返回量</p>"
          },
          {
            "group": "Parameter",
            "type": "json",
            "optional": true,
            "field": "sort",
            "description": "<p>排序方式，如sort={&quot;id&quot;:-1}，默认按资源ID逆序排序</p>"
          },
          {
            "group": "Parameter",
            "type": "json",
            "optional": true,
            "field": "where",
            "description": "<p>where条件查询参数，如where={&quot;vehicles&quot;:{&quot;os_version&quot;:5.1.1}}</p>"
          },
          {
            "group": "Parameter",
            "type": "json",
            "optional": true,
            "field": "or",
            "description": "<p>或条件查询参数, 如or=[{&quot;author&quot;:{&quot;e&quot;:&quot;Jake&quot;}},{author:{&quot;e&quot;:&quot;Alex&quot;}}]</p>"
          },
          {
            "group": "Parameter",
            "type": "json",
            "optional": true,
            "field": "in",
            "description": "<p>匹配条件查询参数, 如in={&quot;author&quot;:[&quot;Jake&quot;,&quot;Billy&quot;]}</p>"
          },
          {
            "group": "Parameter",
            "type": "json",
            "optional": true,
            "field": "having",
            "description": "<p>having条件查询参数，如having={&quot;status&quot;:0,&quot;is_root&quot;:1}</p>"
          }
        ]
      }
    },
    "examples": [
      {
        "title": "请求示例:",
        "content": "curl -i http://localhost/api/v1/role/file/download",
        "type": "curl"
      }
    ],
    "filename": "../../http/controller/role_controller.go",
    "groupTitle": "Role",
    "sampleRequest": [
      {
        "url": "http://industry.isoc.360.cn/api/v1/role/file/download"
      }
    ],
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "type": "string",
            "optional": false,
            "field": "Authorization",
            "description": "<p>认证口令，值为: Bearer {{authToken}}, {{authToken}}为用户认证成功返回的Token</p>"
          }
        ]
      }
    }
  },
  {
    "type": "get",
    "url": "/isoc/v1/roles",
    "title": "查询所有系统角色列表",
    "version": "1.0.0",
    "name": "GetAll",
    "group": "Role",
    "description": "<p>查询所有系统角色列表</p>",
    "success": {
      "examples": [
        {
          "title": "请求成功示例:",
          "content": "{\n    \"code\": 0,\n    \"data\": [\n        {\n            \"channel_id\": \"Q00001\",\n            \"id\": 1,\n            \"name\": \"超级管理员\",\n            \"parent_id\": 0\n        },\n        {\n            \"channel_id\": \"T56205\",\n            \"id\": 2,\n            \"name\": \"亿咖通管理员1\",\n            \"parent_id\": 1\n        }\n    ]\n}",
          "type": "json"
        }
      ]
    },
    "filename": "../../http/controller/role_controller.go",
    "groupTitle": "Role",
    "sampleRequest": [
      {
        "url": "http://industry.isoc.360.cn/isoc/v1/roles"
      }
    ],
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "type": "string",
            "optional": false,
            "field": "Authorization",
            "description": "<p>认证口令，值为: Bearer {{authToken}}, {{authToken}}为用户认证成功返回的Token</p>"
          }
        ]
      }
    }
  },
  {
    "type": "get",
    "url": "/isoc/v1/roles/:id",
    "title": "查询某一系统角色信息",
    "version": "1.0.0",
    "name": "GetOne",
    "group": "Role",
    "description": "<p>查询某一系统角色信息</p>",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "int",
            "optional": false,
            "field": "id",
            "description": "<p>系统角色id</p>"
          }
        ]
      }
    },
    "success": {
      "examples": [
        {
          "title": "请求成功示例:",
          "content": "{\n    \"code\": 0,\n    \"data\": {\n        \"channel_id\": \"T56205\",\n        \"id\": 2,\n        \"name\": \"管理员\",\n        \"parent_id\": 1\n    }\n}",
          "type": "json"
        }
      ]
    },
    "filename": "../../http/controller/role_controller.go",
    "groupTitle": "Role",
    "sampleRequest": [
      {
        "url": "http://industry.isoc.360.cn/isoc/v1/roles/:id"
      }
    ],
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "type": "string",
            "optional": false,
            "field": "Authorization",
            "description": "<p>认证口令，值为: Bearer {{authToken}}, {{authToken}}为用户认证成功返回的Token</p>"
          }
        ]
      }
    }
  },
  {
    "type": "put",
    "url": "/isoc/v1/roles/:id",
    "title": "更新系统角色",
    "version": "1.0.0",
    "name": "Update",
    "group": "Role",
    "description": "<p>更新系统角色</p>",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "string",
            "optional": false,
            "field": "channel_id",
            "description": "<p>渠道号</p>"
          },
          {
            "group": "Parameter",
            "type": "string",
            "optional": false,
            "field": "name",
            "description": "<p>角色名称</p>"
          },
          {
            "group": "Parameter",
            "type": "int",
            "optional": false,
            "field": "parent_id",
            "description": "<p>父角色id</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "请求参数示例:",
          "content": "{\n    \"channel_id\":\"O68468\",\n    \"name\":\"超级管理员1\",\n    \"parent_id\":0\n}",
          "type": "json"
        }
      ]
    },
    "success": {
      "examples": [
        {
          "title": "请求成功示例:",
          "content": "{\n    \"code\": 0,\n    \"data\": {\n         \"id\":1,\n         \"channel_id\":\"O68468\",\n         \"name\":\"超级管理员1\",\n         \"parent_id\":0\n    }\n}",
          "type": "json"
        }
      ]
    },
    "filename": "../../http/controller/role_controller.go",
    "groupTitle": "Role",
    "sampleRequest": [
      {
        "url": "http://industry.isoc.360.cn/isoc/v1/roles/:id"
      }
    ],
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "type": "string",
            "optional": false,
            "field": "Authorization",
            "description": "<p>认证口令，值为: Bearer {{authToken}}, {{authToken}}为用户认证成功返回的Token</p>"
          }
        ]
      }
    }
  },
  {
    "type": "post",
    "url": "/api/v1/role/file/upload",
    "title": "文件上传",
    "version": "1.0.0",
    "name": "Upload",
    "group": "Role",
    "description": "<p>文件上传</p>",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "string",
            "optional": true,
            "field": "file_name",
            "description": "<p>文件名称</p>"
          },
          {
            "group": "Parameter",
            "type": "file",
            "optional": false,
            "field": "file",
            "description": "<p>文件</p>"
          }
        ]
      }
    },
    "examples": [
      {
        "title": "请求示例:",
        "content": "curl -i http://localhost/api/v1/role/file/upload",
        "type": "curl"
      }
    ],
    "success": {
      "examples": [
        {
          "title": "请求成功示例:",
          "content": " {\n\t\t\"code\":0,\n\t\t\"msg\":\"\",\n\t\t\"data\":{\n\t\t\t\"file_id\":\"a834qafmcxvadfq1123\"\n\t\t}\n }",
          "type": "json"
        }
      ]
    },
    "filename": "../../http/controller/role_controller.go",
    "groupTitle": "Role",
    "sampleRequest": [
      {
        "url": "http://industry.isoc.360.cn/api/v1/role/file/upload"
      }
    ],
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "type": "string",
            "optional": false,
            "field": "Authorization",
            "description": "<p>认证口令，值为: Bearer {{authToken}}, {{authToken}}为用户认证成功返回的Token</p>"
          }
        ]
      }
    }
  }
] });

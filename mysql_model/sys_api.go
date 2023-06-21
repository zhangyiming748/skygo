package mysql_model

import (
	"errors"

	"skygo_detection/guardian/src/net/qmap"
	"xorm.io/xorm"

	"skygo_detection/lib/common_lib/mysql"
)

type SysApi struct {
	Id          int    `xorm:"not null pk autoincr comment('api编号') INT(11)" json:"id"`
	Method      string `xorm:"not null default 'GET' comment('请求方法') ENUM('DELETE','GET','POST','PUT')" json:"method"`
	Url         string `xorm:"not null default '' comment('资源定位符') VARCHAR(255)" json:"url"`
	Version     string `xorm:"not null default '' comment('api版本') VARCHAR(255)" json:"version"`
	ApiType     string `xorm:"not null default 'rpc' comment('接口类型') ENUM('http','rpc')" json:"api_type"`
	Resource    string `xorm:"not null default '' comment('接口组名称') VARCHAR(32)" json:"resource"`
	Name        string `xorm:"comment('接口名称') VARCHAR(1024)" json:"name"`
	Description string `xorm:"comment('接口描述') VARCHAR(1024)" json:"description"`
}

const (
	API_TYPE_HTTP = "http" //接口类型：http
	API_TYPE_RPC  = "rpc"  //接口类型:rpc
)

func ApiGetWhere(params map[string]interface{}) *xorm.Session {
	session := mysql.GetSession()
	if id, ok := params["id"]; ok {
		session.And("id=?", id)
	}
	if url, ok := params["url"]; ok {
		session.And("url=?", url)
	}

	return session
}

func ApiFetchIdByUrl(url, method string) int {
	api := new(SysApi)
	params := map[string]interface{}{
		"url":    url,
		"method": method,
	}
	if has, err := ApiGetWhere(params).Get(api); err == nil && has {
		return api.Id
	}
	panic(errors.New("AuthorizationDeny"))
}

func (this SysApi) GetApiTreeList() (*qmap.QM, error) {
	httpApi := qmap.QM{}
	rpcApi := qmap.QM{}
	var resourceList []struct {
		ApiType  string
		Resource string
	}
	if err := mysql.GetSession().Table(this).GroupBy("api_type, resource").Find(&resourceList); err == nil {
		for _, item := range resourceList {
			if apis, apiErr := this.GetApiList(item.ApiType, item.Resource); apiErr == nil {
				if item.ApiType == API_TYPE_HTTP {
					httpApi[item.Resource] = apis
				} else if item.ApiType == API_TYPE_RPC {
					rpcApi[item.Resource] = apis
				}
			} else {
				return nil, apiErr
			}
		}
	} else {
		return nil, err
	}
	result := qmap.QM{}
	if len(httpApi) > 0 {
		result["http"] = httpApi
	}
	if len(rpcApi) > 0 {
		result["rpc"] = rpcApi
	}
	return &result, nil
}

func (this SysApi) GetApiList(apiType, resource string) ([]SysApi, error) {
	models := make([]SysApi, 0)
	err := mysql.GetSession().Table(this).Where("api_type = ?", apiType).And("resource = ?", resource).Limit(10000, 0).Find(&models)
	return models, err
}

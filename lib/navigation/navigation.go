package navigation

import (
	"fmt"
	"net/url"
	"skygo_detection/guardian/src/net/qmap"
	"strings"
)

type Navigation struct {
	EvaluateType string
	ModuleName   string
	ModuleType   string
	Argument     []string
}

func (this *Navigation) ArgumentEncoding() {

}

func (this *Navigation) ArgumentDecoding() []Navigation {
	result := []Navigation{}
	if len(this.Argument) > 0 {
		for _, nn := range this.Argument {
			tmp := new(Navigation)
			menu := strings.Split(nn, "_")
			if len(menu) == 3 {
				tmp.EvaluateType = menu[0]
				tmp.ModuleName = menu[1]
				tmp.ModuleType = menu[2]
				result = append(result, *tmp)
			} else {
				continue
			}
		}
	}
	return result
}

func (this *Navigation) IsNavigation(queryParams string) bool {
	urlParams := map[string]string{}
	u := url.URL{RawQuery: queryParams}
	for k, v := range u.Query() {
		if len(v) != 1 {
			continue
		}
		urlParams[k] = v[0]
	}
	if navigation, has := urlParams["navigation"]; has {
		if navigation == "" {
			return false
		}
		t := []string{}
		t = strings.Split(navigation, ",")
		this.Argument = t
		return has
	}
	return false
}

func (this *Navigation) TransformToNavigation(qm qmap.QM) *qmap.QM {
	//"对象":{组件1:[A,B,C],组件2:[A,B,C]
	// id:对象，label:对象，children:[]interface{}
	result := qmap.QM{}
	for k, v := range qm {
		t := result.Slice("data")
		parent := qmap.QM{
			"id":    k,
			"label": k,
		}
		for kk, vv := range v.(qmap.QM) {
			tt := parent.Slice("children")
			children := qmap.QM{
				"id":    fmt.Sprintf("%s_%s", k, kk),
				"label": kk,
			}
			fmt.Println("================")
			for _, vvv := range vv.([]interface{}) {
				ttt := children.Slice("children")
				grandchild := qmap.QM{
					"id":    fmt.Sprintf("%s_%s_%s", k, kk, vvv),
					"label": vvv,
				}
				children["children"] = append(ttt, grandchild)
			}
			parent["children"] = append(tt, children)
		}
		result["data"] = append(t, parent)
	}
	return &result
}

package orm

import "net/url"

type UrlParamsParser interface {
	Parse(UrlParams) *Query
}

type UrlParams = map[string]string

func GetUrlParamsFromRawQuery(queryParamStr string) UrlParams {
	u := url.URL{RawQuery: queryParamStr}
	urlParams := UrlParams{}
	for k, v := range u.Query() {
		if len(v) != 1 {
			continue
		}
		urlParams[k] = v[0]
	}
	return urlParams
}

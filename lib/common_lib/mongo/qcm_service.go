package mongo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type QCMConn struct {
	Ip   string
	Port int
}

const (
	QCM_REQUEST_TYPE = "GET" //qcm配置拉取请求方法类型
)

func GetQCMConfig(url, key string) ([]*QCMConn, error) {
	qcmUrl := fmt.Sprintf("%s?key=%s", url, key)
	req, err := http.NewRequest(QCM_REQUEST_TYPE, qcmUrl, nil)
	if err != nil {
		return nil, err
	}
	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var qcmConn []*QCMConn
	if body, reqErr := ioutil.ReadAll(resp.Body); reqErr == nil {
		connSlices := strings.Split(string(body), ",")
		for _, conn := range connSlices {
			if connInfo := strings.Split(conn, ":"); len(connInfo) == 2 {
				if port, err := strconv.Atoi(connInfo[1]); err == nil {
					tempConn := &QCMConn{
						Ip:   connInfo[0],
						Port: port,
					}
					qcmConn = append(qcmConn, tempConn)
				} else {
					return nil, err
				}
			}
		}
	} else {
		return nil, reqErr
	}
	return qcmConn, nil
}

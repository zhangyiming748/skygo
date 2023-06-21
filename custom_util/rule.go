package custom_util

import (
	"errors"
	"regexp"
	"strings"
)

// 验证ip:port,ip:port,ip:port是否合法
//
//	str := "123.123.123.123:65535,123.123.123.123:65535,123.123.123.123:65535"
func CheckIpPortsStr(str string) (bool, error) {
	// ip匹配正则 `((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))`
	// port匹配正则 `([0-9]|[1-9]\d|[1-9]\d{2}|[1-9]\d{3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])`
	// 这个用来校验整个ip:port,ip:port
	pattern := `^(((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))):([0-9]|[1-9]\d|[1-9]\d{2}|[1-9]\d{3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5]))(,(((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))):([0-9]|[1-9]\d|[1-9]\d{2}|[1-9]\d{3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])))*$`
	matched, err := regexp.MatchString(pattern, str)
	return matched, err
}

// 验证ip:port是否合法
func CheckIpPortStr(str string) (bool, error) {
	pattern := `^(((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))):([0-9]|[1-9]\d|[1-9]\d{2}|[1-9]\d{3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5]))$`
	matched, err := regexp.MatchString(pattern, str)
	return matched, err
}

// 验证url:port是否合法
func CheckUrlPortStr(str string) (bool, error) {
	pattern := `^(\w|-)+:([0-9]|[1-9]\d|[1-9]\d{2}|[1-9]\d{3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$`
	matched, err := regexp.MatchString(pattern, str)
	return matched, err
}

func CheckHttpIpPorts(str string, splitNo string) (bool, error) {
	infos := strings.Split(str, splitNo)
	for _, v := range infos {
		vs := strings.SplitN(v, "://", 2)
		if len(vs) != 2 {
			return false, errors.New("无://")
		}
		if vs[0] != "http" && vs[0] != "https" {
			return false, errors.New("无http/https")
		}
		if ok, _ := CheckBrokers(vs[1], ","); ok == false {
			return false, errors.New("ip:port错误")
		}
	}
	return true, nil
}

// kafka brokers的访问,
// 例子1： iisop-m:9201
// 例子2：127.0.0.1:9201
func CheckBrokers(str string, splitNo string) (bool, error) {
	infos := strings.Split(str, splitNo)
	var ok bool
	var err error
	for _, v := range infos {
		ok, err = CheckUrlPortStr(v)
		if ok {
			continue
		}

		ok, err = CheckIpPortStr(v)
		if ok {
			continue
		}

		return ok, err
	}
	return true, nil
}

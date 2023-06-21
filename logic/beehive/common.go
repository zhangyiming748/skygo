package beehive

import (
	"errors"
	"skygo_detection/service"
)

const (
	GsmSniffer = "gsmSniffer"
	LteSystem  = "lteSystem"
	GsmSystem  = "gsmSystem"
)

func GetHost(srv string) (host string, err error) {
	beehiveConfig := service.LoadBeehiveConfig()
	switch srv {
	case GsmSniffer:
		return beehiveConfig.Host + beehiveConfig.GsmSnifferPort, nil
	case LteSystem:
		return beehiveConfig.Host + beehiveConfig.LteSystemPort, nil
	case GsmSystem:
		return beehiveConfig.Host + beehiveConfig.GsmSystemPort, nil
	default:
		return "", errors.New("请检查参数")
	}
}

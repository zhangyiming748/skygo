package mysql

import (
	"strconv"
	"time"
)

func FindById(id int, modelPtr interface{}) (bool, error) {
	has, err := GetSession().ID(id).Get(modelPtr)
	return has, err
}

func DeleteById(id int, modelPtr interface{}) (bool, error) {
	count, err := GetSession().ID(id).Delete(modelPtr)
	if count == 0 {
		return false, err
	} else {
		return true, err
	}
}

func GetTaskId() string {
	now := int(time.Now().UnixNano() / 1000)
	str := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	assetID := ""
	var remainder int
	var remainderStr string
	for now != 0 {
		remainder = now % 36
		if remainder < 36 && remainder > 9 {
			remainderStr = str[remainder]
		} else {
			remainderStr = strconv.Itoa(remainder)
		}
		assetID = remainderStr + assetID
		now = now / 36
	}
	if len(assetID) > 8 {
		rs := []rune(assetID)
		assetID = string(rs[:8])
	}

	return assetID
}

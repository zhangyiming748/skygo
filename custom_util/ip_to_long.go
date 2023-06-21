package custom_util

import (
	"errors"
	"github.com/thinkeridea/go-extend/exnet"
	"math"
	"strconv"
	"strings"
)

func IpPrint(key int, data string) (uint, error) {
	dataLen, ipValue, err := getLen(4, data)
	if err != nil {
		return 0, err
	} else {
		if key < 0 {
			if int(math.Abs(float64(key))) <= dataLen {
				key = dataLen - int(math.Abs(float64(key)))
			} else {
				return 0, errors.New("IndexError")
			}
		} else {
			if key >= dataLen {
				return 0, errors.New("IndexError")
			}
		}
		return ipValue + uint(key), nil
	}
}
func ipVersionToLen(version int) (int, error) {
	/*Return number of bits in address for a certain IP version.

	>>> _ipVersionToLen(4)
	32
	>>> _ipVersionToLen(6)
	128
	>>> _ipVersionToLen(5)
	Traceback (most recent call last):
	File "<stdin>", line 1, in ?
	File "IPy.py", line 1076, in _ipVersionToLen
	raise ValueError("only IPv4 and IPv6 supported")
	ValueError: only IPv4 and IPv6 supported
	*/

	if version == 4 {
		return 32, nil
	} else if version == 6 {
		return 128, nil
	} else {
		return 0, errors.New("only IPv4 and IPv6 supported")
	}
}
func getLen(version int, data string) (int, uint, error) {
	bits, err := ipVersionToLen(version)
	if err != nil {
		return 0, 0, err
	} else {
		prefixLen, ipValue, err := getPrefixLen(data)
		if err != nil {
			return 0, 0, err
		} else {
			localLen := bits - prefixLen
			return int(math.Pow(2, float64(localLen))), ipValue, nil
		}
	}

}
func getPrefixLen(data string) (prefixLen int, ipValue uint, err error) {
	ip := ""
	xx := strings.Split(data, "/")
	if len(xx) == 1 {
		ip = xx[0]
		prefixLen = -1
	} else if len(xx) > 2 {
		return 0, 0, errors.New("only one '/' allowed in IP Address")
	} else {
		ip = xx[0]
		prefixStr := xx[1]
		prefixLen, err = MaskToLen(prefixStr)
		if err != nil {
			return 0, 0, err
		}
	}
	ipLong, err := exnet.IPString2Long(ip)
	if err != nil {
		return 0, 0, errors.New("invalid ipv4 format")
	}
	return prefixLen, ipLong, err
}
func netmaskToPrefixLen(netmask int) (int, error) {
	/*Convert an Integer representing a netmask to a prefixlen.

	E.g. 0xffffff00 (255.255.255.0) returns 24
	*/
	netLen, err := count0Bits(netmask)
	if err != nil {
		return 0, err
	}
	maskLen := count1Bits(netmask)
	err = checkNetmask(netmask, maskLen)
	if err != nil {
		return 0, err
	}
	return maskLen - netLen, err
}
func count1Bits(num int) int {
	//Find the highest bit set to 1 in an integer.
	ret := 0
	for {
		if num <= 0 {
			break
		}
		num = num >> 1
		ret += 1
	}
	return ret
}
func count0Bits(num int) (int, error) {
	/*
		"Find the highest bit set to 0 in an integer."
	*/

	// this could be so easy if _count1Bits(~int(num)) would work as excepted
	if num < 0 {
		return 0, errors.New("Only positive Numbers please:" + string(num))
	}
	ret := 0
	for {
		if num <= 0 {
			break
		}
		if num&1 == 1 {
			break
		}
		num = num >> 1
		ret += 1
	}
	return ret, nil
}
func checkNetmask(netmask int, maskLen int) error {
	//Checks if a netmask is expressable as a prefixlen.
	num := netmask
	bits := maskLen

	// remove zero bits at the end
	for {
		if (num&1) != 0 || bits == 0 {
			break
		}
		num = num >> 1
		bits -= 1
		if bits == 0 {
			break
		}
	}

	// now check if the rest consists only of ones
	for {
		if bits <= 0 {
			return nil
		}
		if (num & 1) == 0 {
			return errors.New("Netmask 0x" + string(netmask) + "%x can't be expressed as an prefix.")
		}

		num = num >> 1
		bits -= 1
	}

}
func MaskToLen(mask interface{}) (int, error) {
	switch mask.(type) { //多选语句switch
	case string:
		//是字符时做的事情
		strMask := mask.(string)
		if strings.Index(strMask, ".") != -1 {
			// check if the user might have used a netmask like
			// a.b.c.d/255.255.255.0
			bytes := strings.Split(strMask, ".")
			var data [4]int
			if len(bytes) > 4 {
				return 0, errors.New("IPv4 Address with more than 4 bytes")
			} else {
				bytesLen := 4 - len(bytes)
				for i := 0; i < bytesLen; i++ {
					bytes = append(bytes, "0")
				}
				for index, value := range bytes {
					x, err := strconv.Atoi(value)
					if err == nil {
						if x > 255 || x < 0 {
							return 0, errors.New(strMask + "single byte must be 0 <= byte < 256")
						}
						data[index] = x
					}

				}
				netmask := (data[0] << 24) + (data[1] << 16) + (data[2] << 8) + data[3]
				maskLen, err := netmaskToPrefixLen(netmask)
				return maskLen, err
			}
		} else {
			maskLen, err := strconv.Atoi(mask.(string))
			if err != nil {
				return 0, errors.New("mask is error")
			} else {
				return maskLen, nil
			}
		}
	case int:
		//是整数时做的事情
		return mask.(int), nil
	}
	return 0, errors.New("mask type is error")
}

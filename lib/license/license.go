package license

import (
	"bufio"
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/service"
)

var (
	AES_KEY     = []byte("K*{rMpAaa|vaLvW8l(F3:*]M4BY~?+UL")
	HMAC256_KEY = []byte("2NMAY_TV3iTn51R=>;7D}8_R}ItFjnMxb9-)>>hoG,i~bZ@4")
)
var (
	ErrLicenceInvalid = errors.New("licence is invalid")
)
var (
	licenseInfo    = qmap.QM{}
	licenseMenuMap = map[int]int{}
	licenseRaw     = ""
	loadFlag       = false
	mu             sync.Mutex
)

const (
	RELEASE_VERSION = "1.1"
	BUILD_VERSION   = "202111"
)

const (
	LICENSE_TYPE_TRIAL    = 1 // 许可证类型:试用版
	LICENSE_TYPE_OFFICIAL = 2 // 许可证类型:正式版
)

func loadLicenseFromDisk() {
	// 初始化尝试加载许可证
	if licenseContent, err := ioutil.ReadFile(service.LoadLicenseConfig().Path); err == nil {
		if info, err := ValidLicence(string(licenseContent)); err == nil {
			if diskId, err := GetDiskUUID(); err == nil {
				if diskId != info.String("disk_uuid") {
					println("许可证校验失败，请导入正确的证书")
					return
				}
			} else {
				println(err.Error())
				return
			}
			mu.Lock()
			defer mu.Unlock()
			licenseRaw = string(licenseContent)
			licenseInfo = info
			licenseInfo["license"] = licenseRaw
			licenseMenuMap = map[int]int{}
			for _, item := range licenseInfo.Slice("menu") {
				itemQM := qmap.QM(item.(map[string]interface{}))
				licenseMenuMap[itemQM.Int("id")] = itemQM.Int("expired_time")
			}
		} else {
			println("系统许可证校验失败", err.Error())
		}
	} else {
		println("系统许可证加载失败", err.Error())
	}
}

func VerifyMenu(menuId int) bool {
	menuMap := GetLicenseMenuMap()
	if expireTime, has := menuMap[menuId]; has {
		if expireTime == -1 || expireTime > int(time.Now().Unix()) {
			return true
		}
	}
	return false
}

/*
	{
	     "active_time": 1,
	     "create_time": 1634111263,
	     "disk_uuid": "aeb6fc55-7fb2-4a6b-aed8-3dff04c2766e",
	     "expired_time": 2,
	     "menu": [
	          {
	               "expired_time": -1,
	               "id": 0,
	               "name": "资产管理"
	          },
	          {
	               "expired_time": 2,
	               "id": 1,
	               "name": "测试任务"
	          },
	          {
	               "expired_time": -1,
	               "id": 2,
	               "name": "报告管理"
	          },
	          {
	               "expired_time": -1,
	               "id": 3,
	               "name": "漏洞管理"
	          }
	     ]
	}
*/
func GenerateLicence(activeTime, expireTime int, disUuid string) (string, error) {
	// 设置license的菜单权限
	menuInfo := []qmap.QM{}
	for k, v := range MenuMap {
		item := qmap.QM{
			"id":   k,
			"name": v,
		}
		if k == TEST_TASK {
			item["expired_time"] = expireTime
		} else {
			item["expired_time"] = -1
		}
		menuInfo = append(menuInfo, item)
	}
	encryptStr, err := EncryptLicence(menuInfo, activeTime, expireTime, disUuid)
	if err != nil {
		return "", err
	}
	signStr, signErr := HMAC256Sign(encryptStr, HMAC256_KEY)
	if signErr != nil {
		return "", signErr
	}
	return EncodeSegment([]byte(fmt.Sprintf("%s.%s", encryptStr, signStr))), nil
}

func GetLicense() qmap.QM {
	if !loadFlag {
		loadLicenseFromDisk()
		loadFlag = true
	}
	return licenseInfo
}
func GetLicenseMenuMap() map[int]int {
	if !loadFlag {
		loadLicenseFromDisk()
		loadFlag = true
	}
	return licenseMenuMap
}
func ImportLicense(licenseStr string) error {
	if info, err := ValidLicence(licenseStr); err == nil {
		if diskId, err := GetDiskUUID(); err == nil {
			if diskId != info.String("disk_uuid") {
				return errors.New("许可证校验失败，请导入正确的证书")
			}
		} else {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		if err := ioutil.WriteFile(service.LoadLicenseConfig().Path, []byte(licenseStr), 0666); err == nil {
			licenseRaw = licenseStr
			licenseInfo = info
			licenseInfo["license"] = licenseRaw
			licenseMenuMap = map[int]int{}
			for _, item := range licenseInfo.Slice("menu") {
				itemQM := qmap.QM(item.(map[string]interface{}))
				licenseMenuMap[itemQM.Int("id")] = itemQM.Int("expired_time")
			}
			loadFlag = true
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

var uuidReg = regexp.MustCompile(`^\s*UUID=([\w\d-]*).*`)

// 获取系统磁盘id
func GetDiskUUID() (string, error) {
	c := exec.Command("bash", "-c", "cat  /etc/fstab")
	if output, err := c.CombinedOutput(); err == nil {
		uuid := ""
		r := strings.NewReader(string(output))
		rd := bufio.NewReader(r)
		for {
			line, err := rd.ReadString('\n')
			if err != nil || io.EOF == err {
				break
			}
			if params := uuidReg.FindStringSubmatch(line); len(params) > 0 && params[1] != "" {
				uuid = params[1]
				break
			}
		}
		return uuid, nil
	} else {
		return "", errors.New("disk uuid not found")
	}
}

func ValidLicence(licence string) (qmap.QM, error) {
	licenceBase64, err := DecodeSegment(licence)
	if err != nil {
		return nil, ErrLicenceInvalid
	}
	s := strings.Split(string(licenceBase64), ".")
	if len(s) != 2 {
		return nil, ErrLicenceInvalid
	}
	if err := HMAC256Verify(s[0], s[1], HMAC256_KEY); err != nil {
		return nil, err
	}
	return DecryptLicence(s[0])
}

func EncryptLicence(menuInfo []qmap.QM, activeTime, expiredTime int, disUuid string) (string, error) {
	licence := qmap.QM{
		"active_time":     activeTime,
		"expired_time":    expiredTime,
		"menu":            menuInfo,
		"create_time":     time.Now().Unix(),
		"disk_uuid":       disUuid,
		"release_version": RELEASE_VERSION,
		"build_version":   BUILD_VERSION,
		"license_type":    LICENSE_TYPE_OFFICIAL,
	}
	encodeByte, err := new(sys_service.CryptService).AesEncrypt([]byte(licence.ToString()), AES_KEY[0:16], AES_KEY[16:32])
	if err != nil {
		return "", err
	}
	return EncodeSegment(encodeByte), nil
}

func DecryptLicence(license string) (qmap.QM, error) {
	licenseBase64, err := DecodeSegment(license)
	if err != nil {
		return nil, ErrLicenceInvalid
	}
	decrypt, err := new(sys_service.CryptService).AesDecrypt(licenseBase64, AES_KEY[0:16], AES_KEY[16:32])
	if err != nil {
		return nil, err
	}
	if license, err := qmap.NewWithString(string(decrypt)); err == nil {
		return license, nil
	} else {
		return nil, err
	}
}

// Implements the HMAC-SHA family of signing methods signing methods
// Expects key type of []byte for both signing and validation
type SigningMethodHMAC256 struct{}

var (
	ErrInvalidKey       = errors.New("key is invalid")
	ErrInvalidKeyType   = errors.New("key is of invalid type")
	ErrHashUnavailable  = errors.New("the requested hash function is unavailable")
	ErrSignatureInvalid = errors.New("signature is invalid")
)

// Verify the signature of HSXXX tokens.  Returns nil if the signature is valid.
func HMAC256Verify(signingString, signature string, key interface{}) error {
	// Verify the key is the right type
	keyBytes, ok := key.([]byte)
	if !ok {
		return ErrInvalidKeyType
	}

	// Decode signature, for comparison
	sig, err := DecodeSegment(signature)
	if err != nil {
		return err
	}

	// This signing method is symmetric, so we validate the signature
	// by reproducing the signature from the signing string and key, then
	// comparing that against the provided signature.
	hasher := hmac.New(crypto.SHA256.New, keyBytes)
	hasher.Write([]byte(signingString))
	if !hmac.Equal(sig, hasher.Sum(nil)) {
		return ErrSignatureInvalid
	}

	// No validation errors.  Signature is good.
	return nil
}

// Implements the Sign method from SigningMethod for this signing method.
// Key must be []byte
func HMAC256Sign(signingString string, key interface{}) (string, error) {
	if keyBytes, ok := key.([]byte); ok {
		hasher := hmac.New(crypto.SHA256.New, keyBytes)
		hasher.Write([]byte(signingString))

		return EncodeSegment(hasher.Sum(nil)), nil
	}

	return "", ErrInvalidKeyType
}

// Encode JWT specific base64url encoding with padding stripped
func EncodeSegment(seg []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(seg), "=")
}

// Decode JWT specific base64url encoding with padding stripped
func DecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}

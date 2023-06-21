package common

const (
	// DefaultLimitNum 默认显示条数
	DefaultLimitNum = 10

	// MaxLimitNum 最大显示条数
	MaxLimitNum = 200

	// 通用名称最大长度
	CommonNameMaxLen = 30

	// 字符编码
	CharacterEncodingUtf8 = 1
)

var CharacterEncodingMap = map[int]string{
	CharacterEncodingUtf8: "UTF-8",
}

package ump

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"skygo_detection/custom_util"
	"skygo_detection/custom_util/blog"
	"skygo_detection/lib/aes256"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/mysql_model"
	ump_model "skygo_detection/mysql_model/ump"
	"skygo_detection/service"
	"time"
)

const (
	CodeSuccess = 200
	UrlVerify   = "umpapp/verify/"
)

type VerifyParams struct {
	ClientId    string `json:"client_id"`
	AccessToken string `json:"access_token"`
}

type VerifyResponse struct {
	Code   int64    `json:"code"`
	Msg    string   `json:"msg"`
	Data   UserInfo `json:"data"`
	Status string   `json:"status"`
}

type UserInfo struct {
	UserId      int    `json:"user_id"`
	UserName    string `json:"username"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	IsStaff     bool   `json:"is_staff"`
	IsSuperuser bool   `json:"is_superuser"`
	IsActive    bool   `json:"is_active"`
}

// 应用验证
func Authenticate(rawAccessToken string, ctx *gin.Context) (userInfo UserInfo, err error) {
	umpConfig := service.LoadUmpConfig()
	clientSecret := umpConfig.ClientSecret
	accessToken := aes256.Decrypt(rawAccessToken, clientSecret)
	args := VerifyParams{
		ClientId:    umpConfig.ClientID,
		AccessToken: accessToken,
	}
	authVerifyUrl := fmt.Sprintf("%s/%s", umpConfig.ServerSso, UrlVerify)
	resp, err := custom_util.HttpPostJson(nil, args, authVerifyUrl)
	if err != nil {
		blog.Error("Authenticate HttpPostJson", zap.Any("request:", args),
			zap.Any("authVerifyUrl:", authVerifyUrl),
			zap.Any("responds:", resp), zap.Any("err:", err))
	}
	var verifyResponse VerifyResponse
	err = json.Unmarshal(resp, &verifyResponse)
	if err != nil {
		blog.Debug("Authenticate Unmarshal", zap.Any("resp:", resp),
			zap.Any("authVerifyUrl:", authVerifyUrl),
			zap.Any("responds:", resp), zap.Any("err:", err))
		return
	}
	if verifyResponse.Code != CodeSuccess {
		blog.Error("Authenticate Unmarshal", zap.Any("request:", args),
			zap.Any("responds:", verifyResponse), zap.Any("err:", err))
		return userInfo, errors.New(verifyResponse.Msg)
	}
	data := verifyResponse.Data
	userInfo = UserInfo{
		UserName:    data.UserName,
		Email:       data.Email,
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		IsStaff:     data.IsStaff,
		IsSuperuser: data.IsSuperuser,
		IsActive:    data.IsActive,
	}
	userInfo.UserName = "skygo"
	blog.Debug("Authenticate HttpPostJson", zap.Any("request:", args), zap.Any("userInfo:", userInfo),
		zap.Any("authVerifyUrl:", authVerifyUrl),
		zap.Any("responds:", resp), zap.Any("err:", err))
	user := new(mysql_model.SysUser).GetUserFindByUsername(userInfo.UserName)
	userId := 0
	// 不存在该用户
	if user == nil {
		// 添加 sys 用户
		_, err := CreateUserByUmp(userInfo, accessToken)
		if err != nil {
			return userInfo, err
		}
		// 获取用户ID
		user = new(mysql_model.SysUser).GetUserFindByUsername(userInfo.UserName)
		userId = user.Id
		userInfo.UserId = userId
		// 添加 ump 用户
		_, err = createUmpUser(userInfo)
		if err != nil {
			return userInfo, err
		}
	} else {
		userId = user.Id
	}
	// 授权登录
	token, err := service.GenerateJWT(fmt.Sprintf("%d", userId), "", 3153600000)
	rURL := fmt.Sprintf("%s?token=%s", umpConfig.LoginUrl, token)
	if err != nil {
		blog.Error("Authenticate GenerateJWT", zap.Any("user_id:", userId),
			zap.Any("redirect url:", rURL),
			zap.Any("token:", token), zap.Any("err:", err))
		return userInfo, err
	}
	blog.Debug("Authenticate GenerateJWT",
		zap.Any("redirect url:", rURL), zap.Any("err:", err))
	ctx.Redirect(http.StatusMovedPermanently, rURL)
	return
}

func createUmpUser(useInfo UserInfo) (res int64, err error) {
	umpModel := new(ump_model.UmpUser)
	// 查询该用户是否存在
	has, err := umpModel.Get(useInfo.UserId)
	if err != nil {
		return
	}
	if !has {
		nowTime := time.Now()
		umpModel.IsStaff = ump_model.IsNotStaff
		umpModel.IsSuperuser = ump_model.IsNotSuperUser
		umpModel.IsActive = ump_model.IsNotActive
		if useInfo.IsStaff {
			umpModel.IsStaff = ump_model.IsStaff
		}
		if useInfo.IsSuperuser {
			umpModel.IsSuperuser = ump_model.IsSuperUser
		}
		if useInfo.IsActive {
			umpModel.IsActive = ump_model.IsActive
		}
		umpModel.UserId = useInfo.UserId
		umpModel.UmpUserId = useInfo.UserId
		umpModel.UserName = useInfo.UserName
		umpModel.Email = useInfo.Email
		umpModel.FirstName = useInfo.FirstName
		umpModel.LastName = useInfo.LastName
		umpModel.CreateTime = nowTime.Format("2006-01-02 15:04:05")
		umpModel.UpdateTime = nowTime.Format("2006-01-02 15:04:05")
		res, err = umpModel.Create()
		if err != nil {
			return res, err
		}
	}
	return
}

// 统一登录创建用户
func CreateUserByUmp(userInfo UserInfo, password string) (res int64, err error) {
	newUser := mysql_model.SysUser{
		Username:    userInfo.UserName,
		ChannelId:   "Q00001",
		RoleId:      1,
		Realname:    userInfo.UserName,
		Password:    service.HashPassword(password),
		AccountType: "user",
		Nickname:    userInfo.UserName,
		Email:       userInfo.Email,
		Status:      mysql_model.USER_STATUS_NORMAL,
	}
	s := mysql.GetSession().Table(mysql_model.SysUser{})
	res, err = s.InsertOne(newUser)
	if err != nil {
		return
	}
	return
}

package logic

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"
	"xorm.io/builder"

	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/session"
	"skygo_detection/mysql_model"
	"skygo_detection/service"
)

type AuthLogic struct {
}

// 检验token合法性
func (this *AuthLogic) VerifyToken(token string) (*VerifyTokenResponse, error) {
	if jwtClaim, err := service.TokenValid(token); err == nil {
		userModel := mysql_model.SysUser{}
		if has, err := mysql.GetSession().ID(jwtClaim.Id).Get(&userModel); has && err == nil {
			if userId, err := strconv.ParseInt(jwtClaim.Id, 10, 64); err == nil {
				req := &VerifyTokenResponse{
					UserId:       userId,
					Username:     userModel.Username,
					ChannelId:    userModel.ChannelId,
					Hmd5:         service.HMD5(fmt.Sprintf("%v-%s-%s", userModel.Id, userModel.Username, userModel.ChannelId), ""),
					RoleId:       int64(userModel.RoleId), // todo
					AccountType:  userModel.AccountType,
					CryptKey:     "", // todo
					GlobalRoleId: int64(userModel.RoleId),
				}
				return req, nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

type VerifyTokenResponse struct {
	UserId       int64  `json:"user_id,omitempty"`
	Username     string `json:"username,omitempty"`
	ChannelId    string `json:"channel_id,omitempty"`
	Hmd5         string `json:"hmd5,omitempty"`
	RoleId       int64  `json:"role_id,omitempty"`
	AccountType  string `json:"account_type,omitempty"`
	CryptKey     string `json:"crypt_key,omitempty"`
	GlobalRoleId int64  `json:"global_role_id,omitempty"`
}

// 查询指定服务用户列表
func (this *AuthLogic) GetSpecifiedServiceUsers(serviceName string, roleId int) ([]map[string]interface{}, error) {
	saasRoleIds := new(mysql_model.SysSaasRoleDetail).GetServiceSaasRoleIds(serviceName, roleId)

	session := mysql.GetSession().Where(builder.In("role_id", saasRoleIds))

	widget := orm.PWidget{}
	widget.SetTransformerFunc(this.userTransform)
	if list, err := widget.All(session, &[]mysql_model.SysUser{}); err != nil {
		return nil, err
	} else {
		return list, nil
	}
}

func (this *AuthLogic) userTransform(data qmap.QM) qmap.QM {
	data["channel_name"] = new(mysql_model.SysVehicleFactory).GetChannelName(data.String("channel_id"))
	data["role_name"] = new(mysql_model.SysSaasRole).GetSaasRoleName(data.Int("role_id"))
	return data
}

func (this *AuthLogic) GetCurrentUserInfo(ctx gin.Context) (qmap.QM, error) {
	if user, has := new(mysql_model.SysUser).FindById(int(session.GetUserId(http_ctx.GetHttpCtx(&ctx)))); has {
		res := &qmap.QM{
			"id":         user.Id,
			"role_id":    user.RoleId,
			"username":   user.Username,
			"realname":   user.Realname,
			"channel_id": user.ChannelId,
			"nickname":   user.Nickname,
			"email":      user.Email,
			"status":     user.Status,
		}
		return *res, nil
	} else {
		return nil, errors.New("AccountNotFound")
	}
}

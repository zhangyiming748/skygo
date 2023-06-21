package mysql_model

import (
	"errors"

	"skygo_detection/guardian/src/net/qmap"

	"xorm.io/xorm"

	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/service"
)

type SysUser struct {
	Id            int    `xorm:"not null pk autoincr INT(20)"`
	ChannelId     string `xorm:"not null default '' comment('渠道号') CHAR(6)"`
	Username      string `xorm:"comment('用户名') unique VARCHAR(64)"`
	Password      string `xorm:"comment('密码') VARCHAR(127)"`
	Nickname      string `xorm:"comment('昵称') VARCHAR(64)"`
	Realname      string `xorm:"comment('真实姓名') VARCHAR(30)"`
	RoleId        int    `xorm:"not null default 0 comment('角色id') INT(11)"`
	Email         string `xorm:"comment('邮箱') index VARCHAR(200)"`
	Mobile        int64  `xorm:"not null default 0 comment('手机号') BIGINT(11)"`
	HeadPic       string `xorm:"comment('头像地址') VARCHAR(100)"`
	Sex           int    `xorm:"not null default 0 comment('性别(0:女 1:男)') TINYINT(1)"`
	AccountType   string `xorm:"not null default 'user' comment('账户类型(用户、车机、tbox)') ENUM('tbox','user','vehicle')"`
	Status        int    `xorm:"not null default 2 comment('状态(0:待审核 1:禁用 2:正常 99:删除)') TINYINT(1)"`
	AuthorizeTime int    `xorm:"default 0 comment('token单次授权时长(单位:秒)，如果为0，则单次授权时长取决于配置文件') INT(10)"`
	CreateTime    int    `xorm:"created not null default 0 comment('创建时间') INT(10)"`
}

const (
	ACCOUNT_TYPE_USER    = "user"
	ACCOUNT_TYPE_VEHICLE = "vehicle"
	ACCOUNT_TYPE_TBOX    = "tbox"

	USER_STATUS_UNCHECK   = 0  // 用户状态:未审核
	USER_STATUS_FORBIDDEN = 1  // 用户状态:禁用
	USER_STATUS_NORMAL    = 2  // 用户状态:正常
	USER_STATUS_REMOVE    = 99 // 用户状态:删除
)

func UserGetWhere(params map[string]interface{}) *xorm.Session {
	session := mysql.GetSession()
	if id, ok := params["id"]; ok {
		session.And("id=?", id)
	}
	if channelId, ok := params["channelId"]; ok && channelId != "" {
		session.And("channel_id=?", channelId)
	}
	if username, ok := params["username"]; ok {
		session.And("username=?", username)
	}

	return session
}

func (this *SysUser) Create() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}

func (this *SysUser) FindById(id int) (*SysUser, bool) {
	if has, err := mysql.GetSession().ID(id).Get(this); err != nil {
		panic(err)
	} else {
		return this, has
	}
}

func (this *SysUser) GetUserRoleId(id int) int {
	if has, err := mysql.GetSession().ID(id).Get(this); err == nil && has {
		return this.RoleId
	} else {
		return 0
	}
}

func (this *SysUser) UpdateById(id int, data qmap.QM) (*SysUser, error) {
	if _, has := this.FindById(id); has {
		if _, err := mysql.GetSession().Table(this).ID(this.Id).Update(data); err != nil {
			return nil, err
		} else {
			newUser := new(SysUser)
			newUser.FindById(this.Id)
			return newUser, nil
		}
	} else {
		return nil, errors.New("AccountNotFound") // tODO
	}
}

func (this *SysUser) Delete(channelId string, id int) (int64, error) {
	params := qmap.QM{
		"channelId": channelId,
		"id":        id,
	}
	return UserGetWhere(params).Delete(this)
}

func (this *SysUser) GetUserFindByUsername(username string) *SysUser {
	user := new(SysUser)
	params := map[string]interface{}{"username": username}
	if has, err := UserGetWhere(params).Get(user); err != nil {
		panic(err)
	} else {
		if has {
			return user
		} else {
			return nil
		}
	}
}

func (this *SysUser) ChangePassword(channelId string, id int, oldPassword, newPassword string) error {
	if oldPassword == newPassword {
		return errors.New("SameNewPasswordError")
	}
	if len(newPassword) < 6 {
		return errors.New("NewPasswordTooShotError")
	}

	user := &SysUser{}
	if has, _ := mysql.GetSession().Where("id = ?", id).Get(user); has {
		if checkErr := service.CheckPassword(user.Password, oldPassword); checkErr == nil {
			updateData := qmap.QM{
				"password": service.HashPassword(newPassword),
			}

			if _, err := mysql.GetSession().Table(this).ID(id).Update(updateData); err == nil {
				return nil
			} else {
				return err
			}
		} else {
			return errors.New("OldPasswordIncorrectError")
		}
	} else {
		return errors.New("AccountNotFound")
	}
}

// 根据用户id查询用户信息
func (this *SysUser) GetUserInfo(userId int) (qmap.QM, error) {
	if user, has := new(SysUser).FindById(userId); has {
		res := qmap.QM{
			"id":         user.Id,
			"role_id":    user.RoleId,
			"username":   user.Username,
			"realname":   user.Realname,
			"channel_id": user.ChannelId,
			"nickname":   user.Nickname,
			"email":      user.Email,
			"status":     user.Status,
		}
		return res, nil
	} else {
		return nil, errors.New("This account was not found")
	}
}

// 查询某一个用户管辖范围内的所有用户id列表
// func (this *SysUser) GeRestrainUserIds(channelId, service string, userId int64) []int {
//	roleId := UserRoleGetRoleId(userId, service)
//	userParams := qmap.QM{
//		"e_service":   service,
//		"gte_role_id": roleId,
//	}
//	roleUserIds := []int{}
//	//先查询对应服务中小于等于给定用户权限的用户id列表
//	if err := sys_service.NewSessionWithCond(userParams).Table("sys_user_role").Cols("user_id").Find(&roleUserIds); err == nil {
//		//如果该用户不是超级管理员，并且提供了查询渠道号，则对用户列表按照渠道号进行再过滤
//		if channelId != "" && roleId != common.SUPER_ADMINISTRATE_ROLE_ID {
//			userIds := []int{}
//			params := qmap.QM{
//				"in_id":        userIds,
//				"e_channel_id": channelId,
//			}
//			if err := sys_service.NewSessionWithCond(params).Table("sys_user").Cols("id").Find(&userIds); err == nil {
//				return userIds
//			} else {
//				panic(err)
//			}
//		}
//		return roleUserIds
//	} else {
//		panic(err)
//	}
// }

func SysUserFindById(userId int) (*SysUser, error) {
	model := SysUser{}
	if has, err := mysql.GetSession().ID(userId).Get(&model); err != nil {
		return nil, err
	} else {
		if has {
			return &model, nil
		} else {
			return nil, errors.New("not found")
		}
	}
}

func (this *SysUser) FindByIds(ids []int) ([]SysUser, error) {
	model := []SysUser{}
	if len(ids) < 1 {
		return model, nil
	}

	if err := mysql.GetSession().In("id", ids).Find(&model); err != nil {
		return nil, err
	} else {
		return model, nil
	}
}

package user

import (
	"errors"
	"soul/apis/dao"
	"soul/apis/dto"
	"soul/apis/service/system/role"
	log "soul/internal/logger"
	"soul/model"
	"soul/model/common"
	"soul/utils"
)

type User struct{}

func (u *User) Register(user dto.SystemRegister) (string, bool) {
	if user.Username == "admin" {
		return "禁止使用admin注册", false
	}

	existUser := dao.SystemUser.GetUserByMobile(user.Mobile)
	if existUser != nil {
		return "手机号已存在", false
	}

	// 创建新用户
	newUser := &model.SystemUser{
		Username: user.Username,
		Mobile:   user.Mobile,
		Password: utils.PasswdMd5Digest(user.Password),
		Email:    user.Email,
	}

	err := dao.SystemUser.CreateUser(newUser)
	if err != nil {
		return err.Error(), false
	}

	return "注册成功", true
}

func (u *User) Add(user dto.SystemAdd) (string, bool) {
	if user.Username == "admin" {
		return "禁止添加admin注册", false
	}

	// todo 判断账号是否存在，手机号、用户名、邮箱
	existUser := dao.SystemUser.GetUserByAccount(user.Mobile)
	if existUser != nil {
		return "手机号已存在", false
	}
	existUser = dao.SystemUser.GetUserByAccount(user.Username)
	if existUser != nil {
		return "用户名已存在", false
	}

	existUser = dao.SystemUser.GetUserByAccount(user.Email)
	if existUser != nil {
		return "邮箱已存在", false
	}

	// 构造roles
	var roles []model.SystemRole
	roleService := role.Role{}
	for _, item := range user.Roles {
		_, ok := roleService.Info(item)
		if !ok {
			return "角色不存在", false
		}
		roles = append(roles, model.SystemRole{
			ID: common.ID{ID: item},
		})
	}

	// 创建新用户
	newUser := &model.SystemUser{
		Username: user.Username,
		Nickname: user.Nickname,
		Mobile:   user.Mobile,
		Avatar:   user.Avatar,
		Password: utils.PasswdMd5Digest(user.Password),
		Email:    user.Email,
		Roles:    roles,
	}

	err := dao.SystemUser.CreateUser(newUser)
	if err != nil {
		return err.Error(), false
	}

	return "注册成功", true
}

func (u *User) Login(user dto.SystemLogin) (map[string]any, error) {
	existUser := dao.SystemUser.GetUserByAccount(user.Account)
	if existUser == nil {
		return nil, errors.New("账号不存在")
	}
	if existUser.Password != utils.PasswdMd5Digest(user.Password) {
		return nil, errors.New("手机号或密码错误")
	}

	token, err := utils.CreateJwtToken(existUser.ID.ID, existUser.Username)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New("生成token发生错误")
	}

	jwtToken, err := utils.ParseJwtToken(token)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"accessToken": token,
		"expires":     jwtToken.ExpiresAt,
		"username":    existUser.Username,
		"roles":       existUser.RolesToList(),
	}, nil
}

func (u *User) Info(userId uint) (*dto.SystemUserInfo, bool) {
	user := dao.SystemUser.GetUserById(userId)

	if user == nil {
		return nil, false
	}

	userinfo := &dto.SystemUserInfo{}
	return userinfo.FromModel(user), true
}

func (u *User) AssignRole(roleId, userId uint) error {
	_, exists := u.Info(userId)
	if !exists {
		return errors.New("用户不存在")
	}

	r := role.Role{}
	_, exists = r.Info(roleId)
	if !exists {
		return errors.New("角色不存在")
	}

	err := dao.SystemUser.AssignRoleById(roleId, userId)
	if err != nil {
		return err
	}

	return nil
}

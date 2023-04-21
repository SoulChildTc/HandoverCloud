package user

import (
	"errors"
	"soul/apis/dao"
	"soul/apis/dto"
	log "soul/internal/logger"
	"soul/model"
	"soul/utils"
)

type User struct{}

func (s *User) Register(user dto.SystemRegister) (string, bool) {
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

func (s *User) Login(user dto.SystemLogin) (map[string]any, error) {
	existUser := dao.SystemUser.GetUserByMobile(user.Mobile)
	if existUser == nil {
		return nil, errors.New("账号不存在")
	}
	if existUser.Password != utils.PasswdMd5Digest(user.Password) {
		return nil, errors.New("手机号或密码错误")
	}

	token, err := utils.CreateJwtToken(int(existUser.ID.ID), existUser.Username)
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
		"roles":       []string{},
	}, nil
}

func (s *User) Info(userId uint) (*model.SystemUser, bool) {
	user := dao.SystemUser.GetUserById(userId)
	if user == nil {
		return nil, false
	}
	return user, true
}

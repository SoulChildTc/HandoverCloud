package dbInitializer

import (
	"errors"
	"soul/apis/dao"
	"soul/apis/service/system/user"
	"soul/model"
	"soul/utils"
)

type InitData struct {
}

func (d *InitData) IsInit() (error, bool) {
	err, has := dao.SystemInitData.HasLock()
	if err != nil {
		return err, false
	}

	if has {
		return nil, true
	}
	return nil, false
}

func (d *InitData) InitData() error {
	err, yes := d.IsInit()
	if err != nil {
		return err
	}

	if yes {
		return errors.New("请勿重复初始化")
	}

	// 创建管理员角色
	newRole := &model.SystemRole{
		RoleName: "admin",
	}
	err = dao.SystemRole.CreateRole(newRole)
	if err != nil {
		return err
	}

	// 创建管理员用户
	newUser := &model.SystemUser{
		Username: "admin",
		Nickname: "管理员",
		Password: utils.PasswdMd5Digest("admin"),
		Roles: []model.SystemRole{
			{RoleName: "admin"},
		},
	}
	err = dao.SystemUser.CreateUser(newUser)
	if err != nil {
		return err
	}

	// 分配角色
	userService := user.User{}
	err = userService.AssignRole(1, 1)
	if err != nil {
		return err
	}

	// 标记已初始化
	err = dao.SystemInitData.Locked()
	if err != nil {
		return err
	}

	return nil
}

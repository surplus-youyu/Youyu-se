package models

import (
	"github.com/jinzhu/gorm"
	"github.com/surplus-youyu/Youyu-se/utils"
)

type Group struct {
	ID       int    `gorm:"column:gid" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	Owner    int    `gorm:"column:owner" json:"owner"`
	Intro    string `gorm:"column:intro" json:"intro"`
	IsPublic int    `gorm:"column:is_public" json:"is_public"`
	Type     int    `gorm:"column:type" json:"type"`
}

type GroupUser struct {
	ID     int `gorm:"column:id"`
	GID    int `gorm:"column:gid"`
	UID    int `gorm:"column:uid"`
	Status int `gorm:"column:status"`
}

func CreateGroup(group Group) {
	tx := DB.Begin()

	var user User
	if err := tx.Find(&user, User{Uid: group.Owner}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(utils.Error{404, "未找到用户", err})
		}
		tx.Rollback()
		panic(err)
	}

	if err := tx.Create(&group); err != nil {
		tx.Rollback()
		panic(err)
	}

	if err := tx.Create(&GroupUser{UID: group.Owner, GID: group.ID}); err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()
}

func JoinedGroupList(uid int) []Group {
	var group []Group

	if err := DB.
		Table("group").
		Select("group.gid, group.name, group.owner, group.intro, intro.is_public task.type").
		Joins("LEFT JOIN group_user ON group_user.gid = group.gid").
		Where("group_user.uid = ?", uid).
		Find(&group).Error; err != nil {
		panic(err)
	}

	return group
}

func GetGroupList() []Group {
	var group []Group

	if err := DB.Find(&group, Group{IsPublic: 1}).Error; err != nil {
		panic(err)
	}

	return group
}

func JoinGroup(uid int, gid int) {
	var user User
	if err := DB.Find(&user, User{Uid: uid}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(utils.Error{404, "未找到用户", err})
		}
		panic(err)
	}

	var group Group
	if err := DB.Find(&group, Group{ID: gid}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(utils.Error{404, "未找到组", err})
		}
		panic(err)
	}

	DB.Create(&GroupUser{GID: gid, UID: uid, Status: 1})
}

func RemoveFromGroup(owner int, gid int, member int) {
	var group Group
	if err := DB.Find(&group, Group{ID: gid, Owner: owner}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(utils.Error{404, "未找到组", err})
		}
		panic(err)
	}

	if err := DB.Table("group_user").
		Where("uid = ? AND gid = ?", member, gid).
		Update(GroupUser{Status: 0}).Error; err != nil {
		panic(err)
	}
}

func DeleteGroup(owner int, gid int) {
	var group Group
	if err := DB.Find(&group, Group{ID: gid, Owner: owner}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(utils.Error{404, "未找到组", err})
		}
		panic(err)
	}

	if err := DB.Table("group_user").
		Where("gid = ?", gid).
		Delete(GroupUser{}).Error; err != nil {
		panic(err)
	}

	if err := DB.Delete(Group{ID: gid, Owner: owner}).Error; err != nil {
		panic(err)
	}
}

func GetGroupMembers(gid int) []User {
	var group Group
	var users []User
	if err := DB.Find(&group, Group{ID: gid}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(utils.Error{404, "未找到组", err})
		}
		panic(err)
	}

	if err := DB.
		Table("user").
		Joins("LEFT JOIN group_user ON group_user.uid = user.uid").
		Where("group_user.gid = ?", gid).
		Find(&users).Error; err != nil {
		panic(err)
	}

	return users
}

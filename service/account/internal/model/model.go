package model

import (
	pb "account/api/accapi"
	"fmt"
	"time"
)

// Kratos hello kratos.
type Kratos struct {
	Hello string
}

type Article struct {
	ID int64
	Content string
	Author string
}

type UserBaseInfo struct {
	UserID			int64
	Name			string
	Gender			string
	Avatar			string
	Description		string
	CreatedDate		string
}

type ParamUid struct {
	UId int64 	`json:"uid"`
} 

type UserInfo struct {
	 UserID 		int64
	 Name			string
	 Tel			string
	 Mail			string
	 Gender 		string
	 Avatar			string
	 Description	string
	 CreatedDate	string
}

func NewUser() *UserInfo{
	return &UserInfo{
		Name:        "tourist",
		CreatedDate: fmt.Sprintf("%d", time.Now().Unix()),
	}
}

func (info *UserInfo) ToBaseInfo() *pb.UserBaseInfo {
	return &pb.UserBaseInfo{
		Uid:      info.UserID,
		Name:        info.Name,
		Gender:      info.Gender,
		Avatar:      info.Avatar,
		Desc: info.Description,
		Created: info.CreatedDate,
	}
}
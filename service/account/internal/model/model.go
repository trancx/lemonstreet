package model

import (
	pb "account/api"
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
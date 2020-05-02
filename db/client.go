package db

import "github.com/n4wei/memo/model"

type Client interface {
	GetUser(string) (*model.User, bool)
	GetAllUsers() []*model.User
	AddUser(string, *model.User) bool
	RemoveUser(string) bool

	GetUserMemo(string, string) (*model.Memo, bool)
	GetAllUserMemos(string) ([]*model.Memo, bool)
	AddUserMemo(string, *model.Memo) bool
	RemoveUserMemo(string, string) bool

	Close() error
}

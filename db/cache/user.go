package cache

import "github.com/n4wei/memo/model"

func (this *cache) GetUser(key string) (*model.User, bool) {
	userData, exist := this.store.Users[key]
	if !exist {
		return nil, false
	}
	return userData.User, true
}

func (this *cache) GetAllUsers() []*model.User {
	users := make([]*model.User, 0, len(this.store.Users))
	for _, userData := range this.store.Users {
		users = append(users, userData.User)
	}
	return users
}

func (this *cache) AddUser(key string, user *model.User) bool {
	if _, exist := this.store.Users[key]; exist {
		return false // do not overwrite
	}

	this.store.Users[key] = &UserData{
		User:  user,
		Memos: map[string]*model.Memo{},
	}
	return true
}

func (this *cache) RemoveUser(key string) bool {
	if _, exist := this.store.Users[key]; !exist {
		return false
	}

	delete(this.store.Users, key)
	return true
}

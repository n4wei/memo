package cache

import (
	"testing"

	"github.com/n4wei/memo/lib/test_helper"
	"github.com/n4wei/memo/model"
)

func setupUserTest() (*cache, []*model.User) {
	cache := &cache{
		store: &Store{Users: map[string]*UserData{}},
	}
	users := []*model.User{
		{UserId: "user0"},
		{UserId: "user1"},
	}
	return cache, users
}

func TestCache_AddRemoveUser(t *testing.T) {
	cache, users := setupUserTest()
	userId := users[0].UserId
	var exist, added, removed bool
	var user *model.User

	_, exist = cache.GetUser(userId)
	test_helper.AssertEqual(t, exist, false)

	added = cache.AddUser(userId, users[0])
	test_helper.AssertEqual(t, added, true)
	user, exist = cache.GetUser(userId)
	test_helper.AssertEqual(t, exist, true)
	test_helper.AssertEqual(t, user, users[0])

	removed = cache.RemoveUser(userId)
	test_helper.AssertEqual(t, removed, true)
	_, exist = cache.GetUser(userId)
	test_helper.AssertEqual(t, exist, false)

	removed = cache.RemoveUser(userId)
	test_helper.AssertEqual(t, removed, false)
}

func TestCache_GetAllUsers(t *testing.T) {
	cache, users := setupUserTest()
	var added bool
	var actualUsers []*model.User

	added = cache.AddUser(users[0].UserId, users[0])
	test_helper.AssertEqual(t, added, true)
	added = cache.AddUser(users[1].UserId, users[1])
	test_helper.AssertEqual(t, added, true)

	actualUsers = cache.GetAllUsers()
	test_helper.AssertDeepEqual(t, actualUsers, users)
}

func TestCache_AddSameUser(t *testing.T) {
	cache, _ := setupUserTest()
	user1 := &model.User{UserId: "same-id", CreateTimestamp: "1"}
	user2 := &model.User{UserId: "same-id", CreateTimestamp: "2"}
	var added, exist bool
	var user *model.User

	added = cache.AddUser("same-id", user1)
	test_helper.AssertEqual(t, added, true)
	user, exist = cache.GetUser("same-id")
	test_helper.AssertEqual(t, exist, true)
	test_helper.AssertEqual(t, user, user1)

	added = cache.AddUser("same-id", user2)
	test_helper.AssertEqual(t, added, false)
	user, exist = cache.GetUser("same-id")
	test_helper.AssertEqual(t, exist, true)
	test_helper.AssertEqual(t, user, user1)
}

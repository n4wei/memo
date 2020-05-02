package cache

import (
	"testing"

	"github.com/n4wei/memo/lib/test_helper"
	"github.com/n4wei/memo/model"
)

func setupMemoTest() (*cache, *model.User, []*model.Memo) {
	cache := &cache{
		store: &Store{Users: map[string]*UserData{}},
	}
	user := &model.User{UserId: "user0"}
	memos := []*model.Memo{
		{MemoId: "memo0"},
		{MemoId: "memo1"},
	}
	return cache, user, memos
}

func TestCache_AddRemoveMemo(t *testing.T) {
	cache, user, memos := setupMemoTest()
	userId := user.UserId
	memoId := memos[0].MemoId
	var exist, added, removed bool
	var memo *model.Memo

	_, exist = cache.GetUserMemo(userId, memoId)
	test_helper.AssertEqual(t, exist, false)

	added = cache.AddUserMemo(userId, memos[0])
	test_helper.AssertEqual(t, added, false)

	_, exist = cache.GetUserMemo(userId, memoId)
	test_helper.AssertEqual(t, exist, false)

	added = cache.AddUser(userId, user)
	test_helper.AssertEqual(t, added, true)

	added = cache.AddUserMemo(userId, memos[0])
	test_helper.AssertEqual(t, added, true)

	memo, exist = cache.GetUserMemo(userId, memoId)
	test_helper.AssertEqual(t, exist, true)
	test_helper.AssertDeepEqual(t, memo, memos[0])

	removed = cache.RemoveUserMemo(userId, memoId)
	test_helper.AssertEqual(t, removed, true)

	_, exist = cache.GetUserMemo(userId, memoId)
	test_helper.AssertEqual(t, exist, false)

	removed = cache.RemoveUserMemo(userId, memoId)
	test_helper.AssertEqual(t, removed, false)
}

func TestCache_GetAllMemos(t *testing.T) {
	cache, user, memos := setupMemoTest()
	userId := user.UserId
	var exist, added bool
	var actualMemos []*model.Memo

	_, exist = cache.GetAllUserMemos(userId)
	test_helper.AssertEqual(t, exist, false)

	added = cache.AddUser(userId, user)
	test_helper.AssertEqual(t, added, true)
	added = cache.AddUserMemo(userId, memos[0])
	test_helper.AssertEqual(t, added, true)
	added = cache.AddUserMemo(userId, memos[1])
	test_helper.AssertEqual(t, added, true)

	actualMemos, exist = cache.GetAllUserMemos(userId)
	test_helper.AssertDeepEqual(t, actualMemos, memos)
}

func TestCache_AddSameMemo(t *testing.T) {
	cache, user, _ := setupMemoTest()
	memo1 := &model.Memo{MemoId: "same-id", Title: "title1"}
	memo2 := &model.Memo{MemoId: "same-id", Title: "title2"}
	userId := user.UserId
	var exist, added bool
	var actualMemo *model.Memo

	added = cache.AddUser(userId, user)
	test_helper.AssertEqual(t, added, true)

	added = cache.AddUserMemo(userId, memo1)
	test_helper.AssertEqual(t, added, true)
	actualMemo, exist = cache.GetUserMemo(userId, "same-id")
	test_helper.AssertEqual(t, exist, true)
	test_helper.AssertDeepEqual(t, actualMemo, memo1)

	added = cache.AddUserMemo(userId, memo2)
	test_helper.AssertEqual(t, added, true)
	actualMemo, exist = cache.GetUserMemo(userId, "same-id")
	test_helper.AssertEqual(t, exist, true)
	test_helper.AssertDeepEqual(t, actualMemo, memo2)
}

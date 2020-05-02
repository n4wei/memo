package cache

import "github.com/n4wei/memo/model"

func (this *cache) GetUserMemo(userKey string, memoKey string) (*model.Memo, bool) {
	userData, exist := this.store.Users[userKey]
	if !exist {
		return nil, false
	}

	memo, exist := userData.Memos[memoKey]
	if !exist {
		return nil, false
	}

	return memo, true
}

func (this *cache) GetAllUserMemos(userKey string) ([]*model.Memo, bool) {
	userData, exist := this.store.Users[userKey]
	if !exist {
		return nil, false
	}

	memos := make([]*model.Memo, 0, len(userData.Memos))
	for _, memo := range userData.Memos {
		memos = append(memos, memo)
	}
	return memos, true
}

func (this *cache) AddUserMemo(userKey string, memo *model.Memo) bool {
	userData, exist := this.store.Users[userKey]
	if !exist {
		return false
	}

	userData.Memos[memo.MemoId] = memo
	return true
}

func (this *cache) RemoveUserMemo(userKey string, memoKey string) bool {
	userData, exist := this.store.Users[userKey]
	if !exist {
		return false
	}

	_, exist = userData.Memos[memoKey]
	if !exist {
		return false
	}

	delete(userData.Memos, memoKey)
	return true
}

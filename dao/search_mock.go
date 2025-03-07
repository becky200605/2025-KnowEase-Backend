package dao

import (
	"KnowEase/models"

	"github.com/stretchr/testify/mock"
)

type MockSearchDao struct {
	mock.Mock
}

// 存入搜索记录
func (ms *MockSearchDao) SyncSearchRecord(Record models.SearchRecord) error {
	args := ms.Called(Record)
	return args.Error(0)
}

// 查询用户搜索记录
func (ms *MockSearchDao) SearchUserRecord(UserID string) ([]models.SearchRecord, error) {
	args := ms.Called(UserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.SearchRecord), args.Error(1)
}

// 模糊搜索帖子内容
func (ms *MockSearchDao) SearchPost(Message string) ([]models.PostMessage, error) {
	args := ms.Called(Message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PostMessage), args.Error(1)
}

// 模糊搜索问答内容
func (ms *MockSearchDao) SearchQA(Message string) ([]models.QAs, error) {
	args := ms.Called(Message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.QAs), args.Error(1)
}

// 更新搜索库
func (ms *MockSearchDao) UpdateRecord(UserID string, Record []string) error {
	args := ms.Called(UserID, Record)
	return args.Error(0)
}

// 删除搜索记录
func (ms *MockSearchDao) DeleteRecord(UserID, Message string) error {
	args := ms.Called(UserID, Message)
	return args.Error(0)
}

// 清空搜索记录
func (ms *MockSearchDao) DeleteAllRecord(UserID string) error {
	args := ms.Called(UserID)
	return args.Error(0)
}

// 查询用户搜索记录
func (ms *MockSearchDao) SearchMessage(Message, UserID string) error {
	args := ms.Called(Message, UserID)
	return  args.Error(0)
}

// 更新搜索库
func (ms *MockSearchDao) UpdateCreatedTime(Record models.SearchRecord) error {
	args := ms.Called(Record)
	return args.Error(0)
}

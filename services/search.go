package services

import (
	"KnowEase/dao"
	"KnowEase/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type SearchService struct {
	SearchDao dao.SearchDaoInterface
}

func NewSearchService(SearchDao dao.SearchDaoInterface) *SearchService {
	return &SearchService{SearchDao: SearchDao}
}

// 将搜索记录存入数据库
func (ss *SearchService) SyncSearchRecord(Record models.SearchRecord) error {
	err := ss.SearchDao.SearchMessage(Record.SearchMessage, Record.UserID)
	if err != gorm.ErrRecordNotFound {
		Record.CreatedAt = time.Now()
		err := ss.SearchDao.UpdateCreatedTime(Record)
		if err != nil {
			//fmt.Printf("%v",err)
			return err
		}
		return nil
	}
	return ss.SearchDao.SyncSearchRecord(Record)
}

// 获取生活帖
func (ss *SearchService) SearchPostByMessage(Message string) ([]models.PostMessage, error) {
	if Message == "" {
		return nil, nil
	}
	return ss.SearchDao.SearchPost(Message)
}

// 获取问答
func (ss *SearchService) SearchQAByMessage(Message string) ([]models.QAs, error) {
	if Message == "" {
		return nil, nil
	}
	return ss.SearchDao.SearchQA(Message)
}

// 获取用户搜索记录
func (ss *SearchService) GetUserSearchRecord(UserID string) ([]models.SearchRecord, error) {

	return ss.SearchDao.SearchUserRecord(UserID)
}

// 更新用户搜索记录
func (ss *SearchService) UpdateSearchRecord(UserID string) {
	Record, err := ss.SearchDao.SearchUserRecord(UserID)
	if err != nil {
		log.Printf("failed to update user %s search record:%v", UserID, err)
	}
	var Message []string
	for _, record := range Record {
		Message = append(Message, record.SearchMessage)
	}
	err = ss.SearchDao.UpdateRecord(UserID, Message)
	if err != nil {
		log.Printf("failed to update user %s search record:%v", UserID, err)
	}
}

// 删除用户搜索记录
func (ss *SearchService) DeleteUserSearchRecord(UserID, Message string) error {
	return ss.SearchDao.DeleteRecord(UserID, Message)
}

// 删除用户全部搜索记录
func (ss *SearchService) DeleteUserAllSearchRecord(UserID string) error {
	return ss.SearchDao.DeleteAllRecord(UserID)
}

// 获取热搜榜
func (ss *SearchService) GetHotRank() ([]models.Hot, error) {
	url := "https://weibo.com/ajax/side/hotSearch"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get hot rank %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response:%v", err)
	}
	var hotrank models.Outer
	err = json.Unmarshal(body, &hotrank)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal:%v", err)
	}
	return hotrank.Data.Realtime[:8], nil

}

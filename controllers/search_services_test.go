package controllers

import (
	"KnowEase/dao"
	"KnowEase/models"
	"KnowEase/services"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)
type Record struct{
	UserID string `json:"userid"`
}
// 忽略时间比对
func equalSearchRecordWithoutCreatedAt(expected, actual interface{}) bool {
	e, ok := expected.(models.SearchRecord)
	if !ok {
		return false
	}
	a, ok := actual.(models.SearchRecord)
	if !ok {
		return false
	}
	e.CreatedAt = time.Time{}
	a.CreatedAt = time.Time{}
	return reflect.DeepEqual(e, a)
}
func TestSearchRelatedService(t *testing.T) {
	//依赖注入
	mockSearchDao := &dao.MockSearchDao{}
	SearchService := services.NewSearchService(mockSearchDao)
	SearchControllers := NewSearchController(SearchService)
	//创建Gin引擎
	r := gin.Default()
	search := r.Group("/api")
	{
		search.POST("/search", SearchControllers.Search)
		search.GET("/search/:userid/getrecord", SearchControllers.GetUserSearchRecord)
		search.POST("/search/:userid/recommend", SearchControllers.GetSearchRecommend)
		search.DELETE("/search/delete", SearchControllers.DeleteUserRecord)
		search.DELETE("/search/:userid/deleteall", SearchControllers.DeleteUserALLRecord)
		search.GET("/search/gethotrank", SearchControllers.GetHotRank)
		search.POST("/search/syncRecord", SearchControllers.SyncSearchRecord)
	}

	t.Run("SearchTest1", func(t *testing.T) {
		//搜索
		TestRecord := models.SearchRecord{
			SearchMessage: "啦啦啦啦",
			UserID:        "111",
			CreatedAt:     time.Now(),
		}
		TestRecordJson, err := json.Marshal(TestRecord)
		if err != nil {
			t.Fatalf("failed to marshal test SearchRecord")
		}
		req, err := http.NewRequest("POST", "/api/search", bytes.NewBuffer(TestRecordJson))
		if err != nil {
			t.Fatalf("failed to publish search message :%v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		re := httptest.NewRecorder()
		mockSearchDao.On("SearchPost", TestRecord.SearchMessage).Return([]models.PostMessage{}, nil)
		mockSearchDao.On("SearchQA", TestRecord.SearchMessage).Return([]models.QAs{}, nil)
		//mockSearchDao.On("SyncSearchRecord",TestRecord).Return(nil)
		r.ServeHTTP(re, req)
		//mockSearchDao.AssertCalled(t,"SyncSearchRecord",TestRecord)
		assert.Equal(t, http.StatusOK, re.Code, "unexpected status code for SyncSearchRecord")
	})

	//模拟搜索帖子时，输入无效的情况
	t.Run("SearchTest2", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/search", nil)
		if err != nil {
			t.Fatalf("failed to create request:%v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		re := httptest.NewRecorder()
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusBadRequest, re.Code, "unexpected status code for testing SyncSearchRecord with invaild json:%v", err)
	})
	//搜索失败
	t.Run("SearchTest3", func(t *testing.T) {
		//搜索
		TestRecord := models.SearchRecord{
			SearchMessage: "啦啦啦啦1",
			UserID:        "1112",
			CreatedAt:     time.Now(),
		}
		TestRecordJson, err := json.Marshal(TestRecord)
		if err != nil {
			t.Fatalf("failed to marshal test SearchRecord")
		}
		req, err := http.NewRequest("POST", "/api/search", bytes.NewBuffer(TestRecordJson))
		if err != nil {
			t.Fatalf("failed to publish search message :%v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		re := httptest.NewRecorder()
		mockSearchDao.On("SearchPost", TestRecord.SearchMessage).Return(nil, assert.AnError)
		mockSearchDao.On("SearchQA", TestRecord.SearchMessage).Return([]models.QAs{}, nil)
		//mockSearchDao.On("SyncSearchRecord",TestRecord).Return(nil)
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusInternalServerError, re.Code, "unexpected status code for SyncSearchRecord")
	})


	//获取用户搜索记录
	//正常情况
	t.Run("GetUserSearchRecord", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/search/111/getrecord", nil)
		if err != nil {
			t.Fatalf("failed to create request:%v", err)
		}
		re := httptest.NewRecorder()
		mockSearchDao.On("SearchUserRecord", "111").Return([]models.SearchRecord{
			models.SearchRecord{
				UserID:        "111",
				SearchMessage: "啦啦啦",
				CreatedAt:     time.Now(),
			}}, nil)
		mockSearchDao.On("UpdateRecord", "111", []string{"啦啦啦"}).Return(nil)
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusOK, re.Code, "unexpected status code for GetUserSearchRecord")
	})

	//输入无效的情况
	t.Run("GetUserRecord2", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/search/invaild_id/getrecord", nil)
		if err != nil {
			t.Fatalf("failed to create request:%v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		re := httptest.NewRecorder()
		mockSearchDao.On("SearchUserRecord", "invaild_id").Return(nil, gorm.ErrRecordNotFound)
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusInternalServerError, re.Code, "unexpected status code for GetUserSearchRecord")
	})
	//输入无效的情况
	t.Run("GetUserRecord3", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/search/1/getrecord", nil)
		if err != nil {
			t.Fatalf("failed to create request:%v", err)
		}
		req.URL.Path = "/api/search//getrecord" 
		req.Header.Set("Content-Type", "application/json")
		re := httptest.NewRecorder()
		//mockSearchDao.On("SearchUserRecord", "invaild_id").Return(nil, gorm.ErrRecordNotFound)
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusBadRequest, re.Code, "unexpected status code for GetUserSearchRecord")
	})
	//搜索失败
	t.Run("GetUserSearchRecord2", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/search/1114/getrecord", nil)
		if err != nil {
			t.Fatalf("failed to create request:%v", err)
		}
		re := httptest.NewRecorder()
		mockSearchDao.On("SearchUserRecord", "1114").Return(nil, assert.AnError)
		//mockSearchDao.On("UpdateRecord", "111", []string{"啦啦啦"}).Return(nil)
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusInternalServerError, re.Code, "unexpected status code for GetUserSearchRecord")
	})

	//搜索推荐
	//正常
	t.Run("GetRecommend1", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/search/111/recommend", nil)
		if err != nil {
			t.Fatalf("failed to create request %v", err)
		}
		re := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		mockSearchDao.On("SearchUserRecord", "111").Return([]models.SearchRecord{
			models.SearchRecord{SearchMessage: "啦啦啦啦",
				UserID: "111",
			}}, nil)
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusOK, re.Code, "unexpected status code for GetSearchRecommend")
	})

	//搜索失败
	t.Run("GetRecommend4", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/search/1113/recommend", nil)
		if err != nil {
			t.Fatalf("failed to create request %v", err)
		}
		re := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		mockSearchDao.On("SearchUserRecord", "1113").Return(nil, assert.AnError)
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusInternalServerError, re.Code, "unexpected status code for GetSearchRecommend")
	})
	//输入无效的情况
	t.Run("GetRecommend2", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/search/invaild_id/recommend", nil)
		if err != nil {
			t.Fatalf("failed to create request %v", err)
		}
		re := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		mockSearchDao.On("SearchUserRecord", "invaild_id").Return(nil, gorm.ErrRecordNotFound)
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusInternalServerError, re.Code, "unexpected status code for GetSearchRecommend")
	})
	//输入无效的情况
	t.Run("GetRecommend3", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/search/1/recommend", nil)
		if err != nil {
			t.Fatalf("failed to create request %v", err)
		}
		req.URL.Path = "/api/search//recommend" 
		req.Header.Set("Content-Type", "application/json")
		re := httptest.NewRecorder()
		//mockSearchDao.On("SearchUserRecord", "invaild_id").Return(nil, gorm.ErrRecordNotFound)
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusBadRequest, re.Code, "unexpected status code for GetSearchRecommend")
	})


	//删除用户搜索记录
	//正常情况
	t.Run("DeleteUserRecord1", func(t *testing.T) {
		TestRecord := models.SearchRecord{
			SearchMessage: "啦啦啦啦",
			UserID:        "111",
			CreatedAt:     time.Now(),
		}
		TestRecordJson, err := json.Marshal(TestRecord)
		if err != nil {
			t.Fatalf("failed to marshal test SearchRecord")
		}
		req, err := http.NewRequest("DELETE", "/api/search/delete", bytes.NewBuffer(TestRecordJson))
		if err != nil {
			t.Fatalf("failed to create request %v", err)
		}
		re := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		mockSearchDao.On("DeleteRecord", "111", TestRecord.SearchMessage).Return(nil)
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusCreated, re.Code, "unexpected status code for deleteuserrecord")
	})

	//输入无效的情况
	t.Run("DeleteUserRecord2", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/api/search/delete", nil)
		if err != nil {
			t.Fatalf("failed to create request %v", err)
		}
		re := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		//mockSearchDao.On("DeleteRecord","111",TestRecord.SearchMessage).Return(nil)
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusBadRequest, re.Code, "unexpected status code for deleteuserrecord")
	})

	//没找到记录
	t.Run("DeleteUserRecord3", func(t *testing.T) {
		TestRecord1 := models.SearchRecord{
			SearchMessage: "啦啦啦啦",
			UserID:        "112",
			CreatedAt:     time.Now(),
		}
		TestRecordJson1, err := json.Marshal(TestRecord1)
		if err != nil {
			t.Fatalf("failed to marshal test SearchRecord")
		}
		req, err := http.NewRequest("DELETE", "/api/search/delete", bytes.NewBuffer(TestRecordJson1))
		if err != nil {
			t.Fatalf("failed to create request %v", err)
		}
		re := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		mockSearchDao.On("DeleteRecord", "112", TestRecord1.SearchMessage).Return(gorm.ErrRecordNotFound)
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusInternalServerError, re.Code, "unexpected status code for deleteuserrecord")
	})

	//清空用户搜索记录
	//正常
	t.Run("DeleteAlRecord", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/api/search/111/deleteall", nil)
		if err != nil {
			t.Fatalf("failed to create request %v", err)
		}
		re := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		mockSearchDao.On("DeleteAllRecord", "111").Return(nil)
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusCreated, re.Code, "unexpected status code for DeleteUserAllRecord")
	})

	//用户不存在
	t.Run("DeleteAllRecord2", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/api/search/1112/deleteall", nil)
		if err != nil {
			t.Fatalf("failed to create request %v", err)
		}
		re := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		mockSearchDao.On("DeleteAllRecord", "1112").Return(gorm.ErrRecordNotFound)
		r.ServeHTTP(re, req)
		mockSearchDao.AssertCalled(t, "DeleteAllRecord", "1112")
		assert.Equal(t, http.StatusInternalServerError, re.Code, "unexpected status code for DeleteUserAllRecord")
	})
	//输入无效
	t.Run("DeleteAllRecord3", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/api/search/1/deleteall", nil)
		if err != nil {
			t.Fatalf("failed to create request %v", err)
		}
		req.URL.Path = "/api/search//deleteall" 
		req.Header.Set("Content-Type", "application/json")
		re := httptest.NewRecorder()
		//mockSearchDao.On("DeleteAllRecord", "1112").Return(gorm.ErrRecordNotFound)
		r.ServeHTTP(re, req)
		//mockSearchDao.AssertCalled(t, "DeleteAllRecord", "1112")
		assert.Equal(t, http.StatusBadRequest, re.Code, "unexpected status code for DeleteUserAllRecord")
	})

	//获取热搜榜
	t.Run("GetHotRank", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/search/gethotrank", nil)
		if err != nil {
			t.Fatalf("failed to create request %v", err)
		}
		re := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusOK, re.Code, "unexpected status code for GetHotRank")
	})

	//搜索记录存入库
	//正常
	t.Run("SyncSearchRecord1", func(t *testing.T) {
		TestRecord := models.SearchRecord{
			SearchMessage: "嘻嘻",
			UserID:        "111",
			CreatedAt:     time.Now(),
		}
		TestRecordJson, err := json.Marshal(TestRecord)
		if err != nil {
			t.Fatalf("failed to marshal test SearchRecord")
		}
		fmt.Printf("Marshaled JSON: %s\n", string(TestRecordJson)) // 打印序列化后的 JSON

		req, err := http.NewRequest("POST", "/api/search/syncRecord", bytes.NewBuffer(TestRecordJson))
		if err != nil {
			t.Fatalf("failed to create request %v", err)
		}
		re := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		mockSearchDao.On("SearchMessage", TestRecord.SearchMessage, TestRecord.UserID).Return(gorm.ErrRecordNotFound)
		mockSearchDao.On("SyncSearchRecord", mock.MatchedBy(func(record models.SearchRecord) bool {
			return equalSearchRecordWithoutCreatedAt(TestRecord, record)
		})).Return(nil)

		// 打印请求信息
		fmt.Printf("Request URL: %s\n", req.URL)
		fmt.Printf("Request Method: %s\n", req.Method)
		fmt.Printf("Request Headers: %v\n", req.Header)

		// 验证 SyncSearchRecord 方法是否被调用
		//mockSearchDao.AssertCalled(t, "SyncSearchRecord", TestRecord)
		r.ServeHTTP(re, req)
		//mockSearchDao.AssertCalled(t, "SyncSearchRecord", mock.MatchedBy(func(record models.SearchRecord) bool {
		//    return equalSearchRecordWithoutCreatedAt(TestRecord, record)
		//}))

		fmt.Printf("Response Status Code: %d\n", re.Code)
		fmt.Printf("Response Body: %s\n", re.Body.String())

		// 验证 SyncSearchRecord 方法是否被调用
		//mockSearchDao.AssertCalled(t, "SyncSearchRecord", TestRecord)
		assert.Equal(t, http.StatusCreated, re.Code, "unexpected status code for SyncRecord")
	})

	//搜索记录已存在的情况
	t.Run("SyncSearchRecord2", func(t *testing.T) {
		TestRecord := models.SearchRecord{
			SearchMessage: "啦啦啦啦",
			UserID:        "111",
			CreatedAt:     time.Now(),
		}
		TestRecordJson, err := json.Marshal(TestRecord)
		if err != nil {
			t.Fatalf("failed to marshal test SearchRecord")
		}
		fmt.Printf("Marshaled JSON: %s\n", string(TestRecordJson))

		req, err := http.NewRequest("POST", "/api/search/syncRecord", bytes.NewBuffer(TestRecordJson))
		if err != nil {
			t.Fatalf("failed to create request %v", err)
		}
		re := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		mockSearchDao.On("SearchMessage", TestRecord.SearchMessage, TestRecord.UserID).Return(nil)
		//TestRecord.CreatedAt=time.Now()
		// 设置 UpdateCreatedTime 方法的调用预期
		mockSearchDao.On("UpdateCreatedTime", mock.MatchedBy(func(record models.SearchRecord) bool {
			return equalSearchRecordWithoutCreatedAt(TestRecord, record)
		})).Return(nil)
		mockSearchDao.On("SyncSearchRecord", mock.MatchedBy(func(record models.SearchRecord) bool {
			return equalSearchRecordWithoutCreatedAt(TestRecord, record)
		})).Return(nil)

		// 打印请求信息
		fmt.Printf("Request URL: %s\n", req.URL)
		fmt.Printf("Request Method: %s\n", req.Method)
		fmt.Printf("Request Headers: %v\n", req.Header)

		// 验证 SyncSearchRecord 方法是否被调用
		//mockSearchDao.AssertCalled(t, "SyncSearchRecord", TestRecord)
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusCreated, re.Code, "unexpected status code for SyncRecord")
	})

	//输入无效
	t.Run("SyncSearchRecord3", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/search/syncRecord", nil)
		if err != nil {
			t.Fatalf("failed to create request %v", err)
		}
		re := httptest.NewRecorder()
		//mockSearchDao.On("SyncSearchRecord",TestRecord).Return(nil)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusBadRequest, re.Code, "unexpected status code for SyncRecord")
	})
	//存入失败的情况
	t.Run("SyncSearchRecord4", func(t *testing.T) {
		TestRecord := models.SearchRecord{
			SearchMessage: "啦啦啦啦",
			UserID:        "1115",
			CreatedAt:     time.Now(),
		}
		TestRecordJson, err := json.Marshal(TestRecord)
		if err != nil {
			t.Fatalf("failed to marshal test SearchRecord")
		}
		fmt.Printf("Marshaled JSON: %s\n", string(TestRecordJson))
		req, err := http.NewRequest("POST", "/api/search/syncRecord", bytes.NewBuffer(TestRecordJson))
		if err != nil {
			t.Fatalf("failed to create request %v", err)
		}
		re := httptest.NewRecorder()
		mockSearchDao.On("SearchMessage", TestRecord.SearchMessage, TestRecord.UserID).Return(gorm.ErrRecordNotFound)
		mockSearchDao.On("SyncSearchRecord",TestRecord).Return(assert.AnError)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(re, req)
		assert.Equal(t, http.StatusInternalServerError, re.Code, "unexpected status code for SyncRecord")
	})
//存入失败的情况2
t.Run("SyncSearchRecord5", func(t *testing.T) {
	TestRecord := models.SearchRecord{
		SearchMessage: "啦啦啦啦1",
		UserID:        "11111",
		CreatedAt:     time.Now(),
	}
	TestRecordJson, err := json.Marshal(TestRecord)
	if err != nil {
		t.Fatalf("failed to marshal test SearchRecord")
	}
	fmt.Printf("Marshaled JSON: %s\n", string(TestRecordJson))
	req, err := http.NewRequest("POST", "/api/search/syncRecord", bytes.NewBuffer(TestRecordJson))
	if err != nil {
		t.Fatalf("failed to create request %v", err)
	}
	re := httptest.NewRecorder()
	mockSearchDao.On("SearchMessage", TestRecord.SearchMessage, TestRecord.UserID).Return(nil)
	//mockSearchDao.On("SyncSearchRecord",TestRecord).Return(nil)
	mockSearchDao.On("UpdateCreatedTime", TestRecord).Return(assert.AnError)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(re, req)
	assert.Equal(t, http.StatusInternalServerError, re.Code, "unexpected status code for SyncRecord")
})
}

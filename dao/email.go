package dao

import (
	"KnowEase/models"

	"gorm.io/gorm"
)

type EmailDao struct {
	db *gorm.DB
}

func NewEmailDao(db *gorm.DB) *EmailDao {
	return &EmailDao{db: db}
}
func ProvideEmailDao(db *gorm.DB) EmailDaoInterface {
	return NewEmailDao(db)
}

type EmailDaoInterface interface {
	WriteCode(email, code string) error
	SearchVerificationCode(code string) (*models.Emailverify, error)
}

// 将验证码信息写入数据库
func (ed *EmailDao) WriteCode(email, code string) error {
	ver := models.Emailverify{
		Email: email,
		Code:  code,
	}
	r := ed.db.Create(&ver)
	return r.Error
}

// 查找验证码
func (ed *EmailDao) SearchVerificationCode(code string) (*models.Emailverify, error) {
	var email models.Emailverify
	re := ed.db.Where("code = ?", code).Take(&email)
	if re.Error != nil {
		if re.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, re.Error
	}
	return &email, nil
}

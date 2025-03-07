package dao

import (
	"KnowEase/models"

	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}
func ProvideUserDao(db *gorm.DB) UserDaoInterface {
	return NewUserDao(db)
}

type UserDaoInterface interface {
	GetUserFromEmail(email string) (*models.User, error)
	CreateNewUser(user *models.User) error
	SearchUserByID(UserID string) (models.User, error)
	ChangeUserMessage(userID, field, value string) error
	SearchAllUser() ([]string, error)
}

// 通过邮箱查找该用户信息
func (ud *UserDao) GetUserFromEmail(email string) (*models.User, error) {
	var user models.User
	re := ud.db.Where("email = ?", email).First(&user)
	if re.Error != nil {
		if re.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, re.Error
	}
	return &user, nil
}

// 创建新用户
func (ud *UserDao) CreateNewUser(user *models.User) error {
	re := ud.db.Create(user)
	return re.Error
}

// 查询用户id
func (ud *UserDao) SearchUserByID(UserID string) (models.User, error) {
	var User models.User
	err := ud.db.Where("id = ?", UserID).First(&User).Error
	return User, err
}

// 修改用户信息
func (ud *UserDao) ChangeUserMessage(UserID, UpdateField, UpdateValue string) error {
	err := ud.db.Model(&models.User{}).Where("id = ?", UserID).Update(UpdateField, UpdateValue).Error
	return err
}

// 查找所有用户
func (ud *UserDao) SearchAllUser() ([]string, error) {
	var UserIDs []string
	err := ud.db.Model(&models.User{}).Select("id").Find(&UserIDs).Error
	if err != nil {
		return nil, err
	}
	return UserIDs, nil
}

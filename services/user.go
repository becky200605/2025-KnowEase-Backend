package services

import (
	"KnowEase/dao"
	"KnowEase/models"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserDao dao.UserDaoInterface
}

func NewUserService(UserDao dao.UserDaoInterface) *UserService {
	return &UserService{UserDao: UserDao}
}

// 用户注册
func (us *UserService) Register(user *models.User) error {
	encryptedpassword, err := EncryptPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to envrypt password:%w", err)
	}
	user.Password = encryptedpassword
	err = us.UserDao.CreateNewUser(user)
	if err != nil {
		return fmt.Errorf("failed to create user:%w", err)
	}
	return nil
}

// 用户密码登录
func (us *UserService) LoginByPassword(LoginMessqge models.Login) (*models.User, error) {
	User, err := us.UserDao.GetUserFromEmail(LoginMessqge.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find usermessage")
	} else if User == nil {
		return nil, fmt.Errorf("this user is not registered")
	}
	return User, nil
}

// 密码加密
func EncryptPassword(password string) (string, error) {
	encryptedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(encryptedpassword), nil
}

// 密码比对
func (us *UserService) ComparePassword(password1, password2 string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))
	return err
}

// 通过邮箱地址查找用户
func (us *UserService) GetUserFromEmail(Email string) (*models.User, error) {
	return us.UserDao.GetUserFromEmail(Email)
}

// 修改密码
func (us *UserService) ChangePassword(UserID, NewPassword string) error {
	EncryptedPassword, err := EncryptPassword(NewPassword)
	if err != nil {
		return fmt.Errorf("failed to encrypt password")
	}
	if err := us.UserDao.ChangeUserMessage(UserID, "password", EncryptedPassword); err != nil {
		return fmt.Errorf("failed to change password")
	}
	return nil
}

// 修改用户背景
func (us *UserService) ChangeUserBackground(UserID, UpdateValue string) error {
	return us.UserDao.ChangeUserMessage(UserID, "page_background_url", UpdateValue)
}

// 修改用户头像
func (us *UserService) ChangeUserPicture(UserID, UpdateValue string) error {
	return us.UserDao.ChangeUserMessage(UserID, "image_url", UpdateValue)
}

// 修改用户邮箱
func (us *UserService) ChangeUserEmail(UserID, UpdateValue string) error {
	return us.UserDao.ChangeUserMessage(UserID, "email", UpdateValue)
}

// 修改用户名
func (us *UserService) ChangeUsername(UserID, UpdateValue string) error {
	return us.UserDao.ChangeUserMessage(UserID, "username", UpdateValue)
}

// 查询用户个人信息
func (us *UserService) SearchUserByID(UserID string) (models.User, error) {
	return us.UserDao.SearchUserByID(UserID)
}
//查询发帖人信息
func (us *UserService) SearchPosterMessage(UserID string) (string, string, error) {
	PosterMessage, err := us.UserDao.SearchUserByID(UserID)
	if err != nil {
		return "", "", err
	}
	return PosterMessage.Username, PosterMessage.ImageURL, nil
}

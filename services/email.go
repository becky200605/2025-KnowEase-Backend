package services

import (
	"KnowEase/dao"
	"fmt"
	"math/rand"

	"github.com/spf13/viper"
	gomail "gopkg.in/gomail.v2"
)

type EmailService struct {
	EmailDao dao.EmailDaoInterface
	UserDao  dao.UserDaoInterface
}

func NewEmailService(EmailDao dao.EmailDaoInterface, UserDao dao.UserDaoInterface) *EmailService {
	return &EmailService{EmailDao: EmailDao, UserDao: UserDao}
}

// 邮箱验证码-注册
func (es *EmailService) Register(email string) error {
	//检验该用户是否已注册
	ExistingUser, err := es.UserDao.GetUserFromEmail(email)
	if err != nil {
		return fmt.Errorf("failed to find user:%w", err)
	}
	if ExistingUser != nil {
		return fmt.Errorf("this user already exists")
	}
	err = es.SendEmailVerification(email)
	if err != nil {
		return err
	}
	return nil
}

// 邮箱验证码-登录
func (es *EmailService) Login(email string) error {
	//检验该用户是否已注册
	ExistingUser, err := es.UserDao.GetUserFromEmail(email)
	if err != nil {
		return fmt.Errorf("failed to find user:%w", err)
	}
	if ExistingUser == nil {
		return fmt.Errorf("this user is not registered")
	}
	return nil
}

// 根据邮箱发送并检验验证码
func (es *EmailService) SendEmailVerification(email string) error {
	code := es.RandomCode(6)
	err := es.EmailDao.WriteCode(email, code)
	if err != nil {
		return fmt.Errorf("failed to insert:%w", err)
	}
	fmt.Println("insert successfully!")
	if err = es.SendEmail(email, code); err != nil {
		return fmt.Errorf("failed to send verification code:%w", err)
	}
	return nil
}

// 发送验证邮件
func (es *EmailService) SendEmail(email, code string) error {
	message := gomail.NewMessage()
	//设置邮件头
	message.SetHeader("From", "2061291860@qq.com")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "小知账号注册")
	//验证码邮件正文部分
	password := viper.GetString("EMAIL_PASSWORD")
	if password == "" {
		return fmt.Errorf("未设置邮件授权码，请检查环境变量")
	}
	message.SetBody("text/plain", "【小知创想家】您的验证码"+code+"该验证码5分钟内有效，请勿泄露于他人！")
	dialer := gomail.NewDialer("smtp.qq.com", 587, "2061291860@qq.com", password)
	if err := dialer.DialAndSend(message); err != nil {
		return err //邮件发送失败
	}

	return nil //发送成功

}

// 生成指定位数的验证码
func (e *EmailService) RandomCode(length int) string {
	const charset = "0123456789"
	var result []byte
	for i := 0; i < length; i++ {
		result = append(result, charset[rand.Intn(len(charset))])
	}
	return string(result)

}

// 检验验证码
func (e *EmailService) VerifyCode(code string) error {
	Email, err := e.EmailDao.SearchVerificationCode(code)
	if err != nil {
		return fmt.Errorf("failed to find code")
	}
	if Email == nil {
		return fmt.Errorf("this verification code was not found")
	}
	return nil
}

package controllers

import (
	"KnowEase/models"
	"KnowEase/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserControllers struct {
	UserService  *services.UserService
	EmailService *services.EmailService
	TokenService *services.TokenService
}

func NewUserControllers(UserService *services.UserService, EmailService *services.EmailService, TokenService *services.TokenService) *UserControllers {
	return &UserControllers{UserService: UserService, EmailService: EmailService, TokenService: TokenService}
}

// Sendemail 发送验证码邮件-注册
// @Summary 发送验证码
// @Description 接收用户的邮箱地址后判断该用户是否已注册并发送验证码邮件
// @Tags 用户注册
// @Accept application/json
// @Produce application/json
// @Param email body models.EmailAddress true "邮箱地址"
// @Success 200 {object} models.Response "验证码邮件发送成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 409 {object} models.Response "该用户已注册"
// @Failure 500 {object} models.Response "发送验证码失败"
// @Router /api/register/sendemail [post]
func (u *UserControllers) RegisterSendemail(c *gin.Context) {
	var email models.EmailAddress
	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := u.EmailService.Register(email.Email); err != nil {
		c.JSON(http.StatusConflict, models.Write(err.Error()))
		return
	}
	if err := u.EmailService.SendEmailVerification(email.Email); err != nil {
		response := models.Write(err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, models.Write("验证码邮件发送成功！"))
}

// Sendemail 发送验证码邮件-登录
// @Summary 发送验证码
// @Description 接收用户的邮箱地址后判断该用户是否注册并发送验证码邮件
// @Tags 用户登录
// @Accept application/json
// @Produce application/json
// @Param email body models.EmailAddress true "邮箱地址"
// @Success 200 {object} models.Response "验证码邮件发送成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 404 {object} models.Response "该用户未注册"
// @Failure 500 {object} models.Response "发送验证码失败"
// @Router /api/login/sendemail [post]
func (u *UserControllers) LoginSendemail(c *gin.Context) {
	var email models.EmailAddress
	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := u.EmailService.Login(email.Email); err != nil {
		c.JSON(http.StatusNotFound, models.Write(err.Error()))
		return
	}
	if err := u.EmailService.SendEmailVerification(email.Email); err != nil {
		response := models.Write(err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, models.Write("验证码邮件发送成功！"))
}

// VerifyEmail 验证验证码-注册
// @Summary 验证验证码
// @Description 接收用户提交的验证码并验证其有效性
// @Tags 用户注册
// @Accept application/json
// @Produce application/json
// @Param code body models.Code true "验证码信息"
// @Success 200 {object} models.Response "验证码验证成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "检验验证码失败"
// @Router /api/register/sendemail/verifycode [post]
func (u *UserControllers) ReigisterVerifyEmail(c *gin.Context) {
	var code models.Code
	if err := c.BindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := u.EmailService.VerifyCode(code.Code); err != nil {
		response := err.Error()
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, models.Write("验证码验证成功！"))

}

// VerifyEmail 验证验证码-登录
// @Summary 验证验证码
// @Description 接收用户提交的验证码并验证其有效性,并返回登录信息
// @Tags 用户登录
// @Accept application/json
// @Produce application/json
// @Param code body models.Emailverify true "验证码和邮箱信息"
// @Success 200 {object} map[string]interface{} "验证码验证成功并返回token和用户信息"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "检验验证码失败"
// @Router /api/login/sendemail/verifycode [post]
func (u *UserControllers) LoginVerifyEmail(c *gin.Context) {
	var code models.Emailverify
	if err := c.BindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := u.EmailService.VerifyCode(code.Code); err != nil {
		response := err.Error()
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	User, err := u.EmailService.UserDao.GetUserFromEmail(code.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
	}
	token := u.TokenService.GenerateToken(User)
	c.JSON(http.StatusOK, gin.H{"message": "登录成功!", "token": token, "usermessage": *User})

}

// Register 用户注册
// @Summary 用户注册
// @Description 接收用户的注册信息，生成默认用户名并完成注册，返回用户信息和认证 token
// @Tags 用户注册
// @Accept application/json
// @Produce application/json
// @Param user body models.User true "用户注册信息"
// @Success 201 {object} map[string]interface{} "注册成功的响应，返回token和用户信息"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "数据库错误或其他服务器错误"
// @Router /api/register/sendemail/verifycode/sendmessage [post]
func (u *UserControllers) Register(c *gin.Context) {
	var User models.User
	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
	}
	//设置身份默认为用户
	User.Role = "user"
	//随机用户uid
	User.ID = u.EmailService.RandomCode(9)
	//防止id重复
	for {
		_, err := u.UserService.SearchUserByID(User.ID)
		if err == gorm.ErrRecordNotFound {
			break
		}
	}
	//生成随机默认用户名
	rand := u.EmailService.RandomCode(5)
	User.Username = "小知用户" + rand
	if err := u.UserService.Register(&User); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	token := u.TokenService.GenerateToken(&User)
	//成功响应后返回token和用户信息
	c.JSON(http.StatusCreated, gin.H{"message": "注册成功", "token": token, "userMessage": User})
}

// LoginByPassword 用户登录
// @Summary 密码登录
// @Description 通过密码和邮箱进行登录，检验成功后返回用户信息和认证 token
// @Tags 用户登录
// @Accept application/json
// @Produce application/json
// @Param user body models.Login true "用户注册信息"
// @Success 201 {object} map[string]interface{} "注册成功的响应，返回token和用户信息"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 404 {object} models.Response "该用户未找到"
// @Failure 409 {object} models.Response "密码不正确"
// @Router /api/login/bypassword [post]
func (u *UserControllers) LoginByPassword(c *gin.Context) {
	var User models.Login
	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Usermessage, err := u.UserService.LoginByPassword(User)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Write(err.Error()))
		return
	}
	if err := u.UserService.ComparePassword((*Usermessage).Password, User.Password); err != nil {
		c.JSON(http.StatusConflict, models.Write("密码不正确"))
		return
	}
	token := u.TokenService.GenerateToken(Usermessage)
	c.JSON(http.StatusOK, gin.H{"message": "登陆成功！", "token": token, "usermessage": *Usermessage})
}

// UpdateSendemail 发送验证码邮件-忘记密码
// @Summary 发送验证码
// @Description 接收用户的邮箱地址后发送验证码邮件
// @Tags 忘记密码
// @Accept application/json
// @Produce application/json
// @Param email body models.EmailAddress true "邮箱地址"
// @Success 200 {object} models.Response "验证码邮件发送成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "发送验证码失败"
// @Router /api/forgetPassword/sendemail [post]
func (u *UserControllers) UpdateSendemail(c *gin.Context) {
	var email models.EmailAddress
	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := u.EmailService.SendEmailVerification(email.Email); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Write("验证码邮件发送成功！"))
}

// VerifyEmail 验证验证码-忘记密码
// @Summary 验证验证码
// @Description 接收用户提交的验证码并验证其有效性
// @Tags 忘记密码
// @Accept application/json
// @Produce application/json
// @Param code body models.Code true "验证码信息"
// @Success 200 {object} models.Response "验证码验证成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "检验验证码失败"
// @Router /api/forgetPassword/sendemail/verifycode [post]
func (u *UserControllers) UpdateVerifyEmail(c *gin.Context) {
	var code models.Emailverify
	if err := c.BindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := u.EmailService.VerifyCode(code.Code); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Write("验证码验证成功！"))
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 通过邮箱进行修改密码
// @Tags 忘记密码
// @Accept application/json
// @Produce application/json
// @Param user body models.Login true "邮箱和新密码"
// @Success 201 {object} models.Response "修改密码成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 404 {object} models.Response "未找到该用户"
// @Failure 500 {object} models.Response "修改密码失败"
// @Router /api/forgetPassword/sendemail/verifycode/update [post]
func (u *UserControllers) ChangePassword(c *gin.Context) {
	var User models.Login
	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	UserMessage, err := u.UserService.GetUserFromEmail(User.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Write("未找到该用户"))
		return
	}
	if err := u.UserService.ChangePassword(UserMessage.ID, User.Password); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("密码修改成功！"))
}

// @Summary 退出登录
// @Description 退出登录并将token强制过期
// @Tags 个人主页-退出登录
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Success 201 {object} models.Response "登出成功"
// @Failure 400 {object} models.Response "token获取失败"
// @Failure 500 {object} models.Response "服务器错误"
// @Router /api/logout [post]
func (uc *UserControllers) Logout(c *gin.Context) {
	// 从请求头中获取token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, models.Write("token获取失败！"))
		return
	}
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	err := uc.TokenService.AddToBlacklist(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("登出成功！"))
}

// @Summary 修改个人主页背景
// @Description 修改用户个人主页的背景图片
// @Tags 个人主页
// @Accept  json
// @Produce  json
// @Param userbackground body models.User true "新背景"
// @Param userid path string true "用户ID"
// @Success 201 {object} map[string]interface{} "响应成功信息以及背景图片的url"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "修改背景失败"
// @Router /api/{userid}/userpage/alterbackground [post]
func (uc *UserControllers) ChangeUserBackground(c *gin.Context) {
	var UserBackground models.User
	if err := c.BindJSON(&UserBackground); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	UserBackground.ID = c.Param("userid")
	if UserBackground.ID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := uc.UserService.ChangeUserBackground(UserBackground.ID, UserBackground.PageBackgroundURL); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("修改背景失败！"))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "修改成功！", "newbackground": UserBackground.PageBackgroundURL})
}

// @Summary 修改个人主页头像
// @Description 修改用户个人主页的头像
// @Tags 个人主页
// @Accept  json
// @Produce  json
// @Param user body models.User true "新头像"
// @Param userid path string true "用户ID"
// @Success 201 {object} map[string]interface{} "响应成功信息以及新头像的url"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "修改头像失败"
// @Router /api/{userid}/userpage/alterimage [post]
func (uc *UserControllers) ChangeUserPicture(c *gin.Context) {
	var User models.User
	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	User.ID = c.Param("userid")
	if User.ID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := uc.UserService.ChangeUserPicture(User.ID, User.ImageURL); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("更新用户头像失败"))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Message": "更新用户头像成功！", "NewImageURL": User.ImageURL})

}

// @Summary 修改个人主页密码
// @Description 修改用户个人主页的密码
// @Tags 个人主页
// @Accept  json
// @Produce  json
// @Param user body models.User true "新密码"
// @Param userid path string true "用户ID"
// @Success 201 {object} models.Response "修改密码成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "修改密码成功"
// @Router /api/{userid}/userpage/alterpassword [post]
func (uc *UserControllers) ChangeUserPassword(c *gin.Context) {
	var User models.User
	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	User.ID = c.Param("userid")
	if User.ID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := uc.UserService.ChangePassword(User.ID, User.Password); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("修改密码成功"))
}

// Sendemail 发送验证码邮件-修改邮箱
// @Summary 发送验证码
// @Description 绑定用户新邮箱地址
// @Tags 个人主页
// @Accept application/json
// @Produce application/json
// @Param email body models.EmailAddress true "邮箱地址"
// @Success 201 {object} models.Response "验证码邮件发送成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "发送验证码失败"
// @Router /api/{userid}/userpage/sendemail [post]
func (u *UserControllers) AlterSendemail(c *gin.Context) {
	var email models.EmailAddress
	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := u.EmailService.SendEmailVerification(email.Email); err != nil {
		response := models.Write(err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, models.Write("验证码邮件发送成功！"))
}

// VerifyEmail 验证验证码-修改邮箱
// @Summary 验证验证码
// @Description 验证新邮箱地址可用性，并修改用户信息
// @Tags 个人主页
// @Accept application/json
// @Produce application/json
// @Param code body models.Emailverify true "验证码和邮箱信息"
// @Param userid path string true "用户ID"
// @Success 201 {object} map[string]interface{} "验证码验证成功并返回新邮箱地址"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 409 {object} models.Response "检验验证码失败"
// @Failure 500 {object} models.Response "修改邮箱地址失败"
// @Router /api/{userid}/userpage/alteremail [post]
func (u *UserControllers) ChangeUserEmail(c *gin.Context) {
	var code models.Emailverify
	if err := c.BindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := u.EmailService.VerifyCode(code.Code); err != nil {
		response := err.Error()
		c.JSON(http.StatusConflict, response)
		return
	}
	if err := u.UserService.ChangeUserEmail(UserID, code.Email); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("更新邮箱地址失败"))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "修改密码成功!", "newEmail": code.Email})

}

// @Summary 修改用户名
// @Description 修改用户个人用户名
// @Tags 个人主页
// @Accept  json
// @Produce  json
// @Param user body models.User true "新名字"
// @Param userid path string true "用户ID"
// @Success 201 {object} map[string]interface{} "响应成功信息以及新名称"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "修改用户名失败"
// @Router /api/{userid}/userpage/alterusername [post]
func (uc *UserControllers) ChangeUsername(c *gin.Context) {
	var User models.User
	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	User.ID = c.Param("userid")
	if User.ID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := uc.UserService.ChangeUsername(User.ID, User.Username); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("更新用户名失败"))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Message": "更新用户名成功！", "NewUsername": User.Username})

}

// @Summary 获取用户基本信息
// @Description 在他人主页获取用户基本信息
// @Tags 个人主页
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "响应成功信息以及用户基本信息"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/userpage/{userid} [get]
func (uc *UserControllers) GetUserMessage(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	UserMessage, err := uc.UserService.SearchUserByID(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询用户信息失败！"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取用户信息成功！", "username": UserMessage.Username, "userUrl": UserMessage.ImageURL, "backgroundUrl": UserMessage.PageBackgroundURL, "userid": UserMessage.ID})
}

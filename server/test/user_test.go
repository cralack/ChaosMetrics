package test

import (
	"testing"

	"github.com/cralack/ChaosMetrics/server/app/provider/user"
	"github.com/cralack/ChaosMetrics/server/model/usermodel"
)

func Test_Register(t *testing.T) {
	serv := user.NewUserService()

	usr1 := &usermodel.User{
		UserName: "snoop dogg",
		Password: "123456",
		Email:    "cralack@foxmail.com",
	}
	usr1.ID = 10086

	// usr2 := &usermodel.User{
	// 	UserName: "tupac",
	// 	Password: "123456",
	// 	Email:    "cralack@qq.com",
	// }
	// register
	userToken, err := serv.PreRegister(usr1)
	if err != nil {
		t.Log(err)
	}
	t.Log(usr1.Token == userToken)

	ok, err := serv.VerifyRegister(usr1.Token)
	if err != nil {
		t.Log(err)
	}
	t.Log(ok)
	usrDB := &usermodel.User{}
	if err = db.Where("username=?", usr1.UserName).First(&usrDB).Error; err != nil {
		t.Log(err)
	}
	t.Log(usrDB.CreatedAt)

	// login
	token, err := serv.Login(usr1.UserName, usr1.Password)
	if err != nil || usrDB == nil {
		t.Log(err)
	}
	t.Log(token)
	usr1.Token = token

	// get user
	usrDB, err = serv.GetUser(usr1.ID)
	t.Log(usrDB.ID == usr1.ID)

	// verify
	sess, err := serv.VerifyLogin(usr1.Token)
	if err != nil || sess == nil {
		t.Log(err)
	}

	// logout
	err = serv.Logout(usr1)
	if err != nil {
		t.Log(err)
	}
	sess, err = serv.VerifyLogin(usr1.Token)
	if err == nil || sess != nil {
		t.Log("?")
	}
}

package test

import (
	"fmt"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_secret"
	"testing"
)

func TestGenTpl(t *testing.T) {
	password := "admin123"
	pwd, err := lv_secret.PasswordHash(password)
	lv_log.Error("------------" + pwd)
	//校验密码
	eq := lv_secret.PasswordVerify(password, pwd)
	fmt.Println(eq, err)
}

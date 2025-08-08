package test

import (
	"common/models"
	"common/myconf"
	"common/util"
	"fmt"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"testing"
)

func TestI18n(t *testing.T) {
	cfg := myconf.GetConfigInstance()
	lv_db.GetOrmDefault()
	fmt.Println(cfg)
	paramId := 1
	model := models.DevParam{}
	modelPtr, _ := model.FindById(paramId)
	//model.GroupCode = "home"
	util.TranslateI18nTagAll("zh", modelPtr)
	fmt.Println("==============modelPtr===========", modelPtr.Name)

	list := []models.DevParam{}
	list = append(list, model)
	util.TranslateI18nTagAll("zh", &list)
	str := lv_conv.ToJsonStr(list)
	fmt.Println("==============list===========", str)

	list = append(list, model)
	util.TranslateI18nTagAll("zh", list)
	str = lv_conv.ToJsonStr(list)
	fmt.Println("============== arr ===========", str)
	//
	//listPtr, _ := dao.GetDevParamDaoInstance().SelectByCodes("home", []string{"FLOAT_AC1_PHASEVOL_A_ID"})
	//util.TranslateI18nTagAll("en", listPtr)
	//str = lv_conv.ToJsonStr(list)
	fmt.Println("==============listPtr ===========" + str)

	fmt.Println("============== over ===========")
}

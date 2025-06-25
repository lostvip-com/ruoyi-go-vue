package test

import (
	"common/myconf"
	"fmt"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/utils/lv_file"
	"path/filepath"
	"system/service"
	"testing"
)

func TestGenTpl(t *testing.T) {
	cfg := myconf.GetConfigInstance()
	lv_db.GetMasterGorm()
	fmt.Println(cfg)
	tableId := int64(7)
	tableService := service.TableService{}
	entity, _ := tableService.FindGenTableById(tableId)
	target := filepath.Join(lv_file.GetCurrentPath(), "tmp")
	distName := "index.vue"
	tpl := service.TplInfo{
		PathSrc:  lv_file.GetCurrentPath() + "/../resources/tpl_gen/vue3/src/views/{{.ModuleName}}/{{.BusinessName}}",
		NameSrc:  distName + ".tpl",
		PathDist: target,
		NameDist: distName,
	}
	svc := service.CodeGenService{}
	src, err := svc.GenCodeByTpl(entity, &tpl)
	fmt.Println(src, err)
}

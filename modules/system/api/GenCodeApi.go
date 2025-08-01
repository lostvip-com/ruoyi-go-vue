package api

import (
	api2 "common/api"
	util2 "common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/lv_batis"
	"github.com/lostvip-com/lv_framework/lv_global"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/utils/lv_file"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"github.com/spf13/cast"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"system/model"
	"system/service"
	"system/vo"
)

type GenCodeApi struct {
	api2.BaseApi
}

// 表单构建
func (w *GenCodeApi) Build(c *gin.Context) {
	util2.BuildTpl(c, "tool/build").WriteTpl()
}

func (w *GenCodeApi) CreateMenu(c *gin.Context) {
	tableName := c.Query("tableName")
	if tableName == "" {
		util2.Fail(c, "param wrong")
		return
	}
	var tableService service.TableService
	entity, err := tableService.FindGenTableByName(tableName)
	lv_err.HasErrAndPanic(err)
	exist := service.GetMenuServiceInstance().IsMenuNameExist(entity.FunctionName)
	if exist {
		util2.Fail(c, "生成代码失败，生成菜单名称已存在")
		return
	}
	// Loads queries from file
	sqlFile := path.Join(lv_global.Config().GetTmpPath(), tableName+"_menu.sql")
	exist = lv_file.IsFileExist(sqlFile)
	if !exist {
		sqlFile = path.Join("../"+lv_global.Config().GetTmpPath(), tableName+"_menu.sql")
		exist = lv_file.IsFileExist(sqlFile)
		if !exist { //检查上组目录中是否存在
			util2.Fail(c, "未找到文件")
			return
		}
	}
	lv_log.Info(lv_file.GetCurrentPath()+"---------file:", sqlFile)
	batis, err := lv_batis.LoadFromFile(sqlFile)
	lv_err.HasErrAndPanic(err)
	// Run queries
	tb := lv_db.GetOrmDefault()
	//cfg := global.GetConfigInstance()
	_, err = batis.ExecMultiSqlInTransaction(tb, "create_menu")
	lv_err.HasErrAndPanic(err)
	util2.Success(c, nil)
}

// swagger文档
func (w *GenCodeApi) Swagger(c *gin.Context) {
	a := c.Query("a")
	if a == "r" {
		//重新生成文档
		curDir, err := os.Getwd()
		if err != nil {
			util2.Fail(c, err.Error())
			return
		}
		genPath := curDir + "/static/swagger"
		err = w.generateSwaggerFiles(genPath)
		if err != nil {
			util2.Fail(c, err.Error())
			c.Abort()
			return
		}
	}
	c.Redirect(http.StatusFound, "/static/swagger/index.html")
}

// 自动生成文档 swag init -o static/swagger
func (w *GenCodeApi) generateSwaggerFiles(output string) error {

	cmd := exec.Command("swag", "init -o "+output)
	// 保证关闭输出流
	if err := cmd.Start(); err != nil { // 运行命令
		return err
	}

	return nil
}
func (w *GenCodeApi) GenList(c *gin.Context) {
	var req *vo.GenTablePageReq
	tableService := service.TableService{}
	if err := c.ShouldBind(&req); err != nil {
		util2.Fail(c, err.Error())
		return
	}
	rows, total, err := tableService.FindPage(req)
	if err != nil {
		util2.Fail(c, err.Error())
		return
	}
	util2.SuccessPage(c, rows, total)
}

func (w *GenCodeApi) RemoveByTableId(c *gin.Context) {
	tableId := c.Param("tableId")
	tableService := service.TableService{}
	err := tableService.DeleteByIds(tableId)
	if err != nil {
		util2.Fail(c, err.Error())
		return
	}
	util2.Success(c, nil)
}

// EditSave 修改数据保存
func (w *GenCodeApi) EditSave(c *gin.Context) {
	var req = new(vo.EditGenTableVO)
	if err := c.ShouldBind(&req); err != nil {
		util2.Fail(c, err.Error())
		return
	}
	w.FillInUpdate(c, &req.BaseModel)
	tableService := service.TableService{}
	err := tableService.SaveEdit(req)
	if err != nil {
		util2.Fail(c, err.Error())
		return
	}
	util2.Success(c, nil)
}

// Preview 预览代码
func (w *GenCodeApi) Preview(c *gin.Context) {
	tableId := lv_conv.Int64(c.Param("tableId"))
	tableService := service.TableService{}
	entity, err := tableService.FindGenTableById(tableId)
	if err != nil {
		util2.Fail(c, err.Error())
	}
	tableService.SetPkColumn(entity, entity.Columns)
	var codeGenService service.CodeGenService
	dataMap := codeGenService.PreviewCode(entity)
	retMap := make(map[string]string)
	for _, mp := range dataMap {
		for k, v := range mp {
			retMap[k] = v
		}
	}
	util2.Success(c, retMap)
}
func (w *GenCodeApi) DataList(c *gin.Context) {
	var req *vo.GenTablePageReq

	err := c.ShouldBind(&req)
	lv_err.HasErrAndPanic(err)
	tableService := service.TableService{}
	rows := make([]model.GenTable, 0)
	result, total, err := tableService.SelectDbTableList(req)
	if err == nil && len(result) > 0 {
		rows = result
	}

	c.JSON(http.StatusOK, lv_dto.TableDataInfo{
		Code:  200,
		Msg:   "操作成功",
		Total: total,
		Rows:  rows,
	})
}
func (w *GenCodeApi) ColumnList(c *gin.Context) {
	tableId := lv_conv.Int64(c.Query("tableId"))
	rows := make([]model.GenTableColumn, 0)
	tableService := service.TableColumnService{}
	result, err := tableService.SelectGenTableColumnListByTableId(tableId)
	if err == nil && len(result) > 0 {
		rows = result
	}
	c.JSON(http.StatusOK, lv_dto.TableDataInfo{
		Code:  200,
		Msg:   "操作成功",
		Total: len(rows),
		Rows:  rows,
	})
}
func (w *GenCodeApi) ImportTableSave(c *gin.Context) {
	tables := c.Query("tables")
	if tables == "" {
		util2.Fail(c, "参数错误tables未选中")
	}
	user := w.GetCurrUser(c)
	operName := user.UserName
	tableService := service.TableService{}
	tableArr := strings.Split(tables, ",")
	tableList, err := tableService.SelectDbTableListByNames(tableArr)
	if err != nil {
		util2.Fail(c, err.Error())
		return
	}
	if tableList == nil {
		util2.Fail(c, "请选择需要导入的表")
		return
	}
	err = tableService.ImportGenTable(&tableList, operName)
	if err != nil {
		util2.Fail(c, err.Error())
		return
	}
	util2.Success(c, nil)
}

// 生成代码
func (w *GenCodeApi) GenCode(c *gin.Context) {
	tableName := c.Param("tableName")
	tableService := service.TableService{}
	entity, err := tableService.FindGenTableByName(tableName)
	lv_err.HasErrAndPanic(err)
	tableService.SetPkColumn(entity, entity.Columns)
	var codeGenService service.CodeGenService
	codeGenService.GenCode(entity, true)
	//(genService)
	util2.Success(c, nil)
}

func (w *GenCodeApi) GetGenTableInfo(c *gin.Context) {
	tableIdStr := c.Param("tableId")
	tableId := cast.ToInt64(tableIdStr)
	m := make(map[string]any)
	var svc service.TableService
	table, err := svc.FindGenTableById(tableId)
	lv_err.HasErrAndPanic(err)
	var svcCol service.TableColumnService
	columns, err := svcCol.SelectGenTableColumnListByTableId(tableId)
	lv_err.HasErrAndPanic(err)
	m["info"] = table
	m["rows"] = columns
	//m["tables"] = selectGenTableAllInfoById()
	util2.Success(c, m)
}

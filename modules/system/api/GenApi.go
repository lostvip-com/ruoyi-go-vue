package api

import (
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/lv_batis"
	"github.com/lostvip-com/lv_framework/lv_global"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"net/http"
	"os"
	"os/exec"
	"system/model"
	"system/service"
	"system/vo"
)

type GenApi struct {
	BaseApi
}

// 表单构建
func (w *GenApi) Build(c *gin.Context) {
	util.BuildTpl(c, "tool/build").WriteTpl()
}

func (w *GenApi) ExecSqlFile(c *gin.Context) {
	tableId := lv_conv.Int64(c.Query("tableId"))

	if tableId <= 0 {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).SetMsg("参数错误").Log("执行SQL文件错误", gin.H{"tableId": tableId})
	}
	genTable := model.GenTable{}
	po, err := genTable.FindById(tableId)
	if err != nil {
		panic(err.Error())
	}

	//err = lv_db.ExecSqlFile(sqlFile)
	// Loads queries from file
	batis, err := lv_batis.LoadFromFile(lv_global.Config().GetTmpPath() + "/" + po.TbName + "_menu.sql")
	// Run queries
	tb := lv_db.GetMasterGorm()
	//cfg := global.GetConfigInstance()
	batis.Exec(tb, "menu")
	menuName := po.FunctionName
	sysmenu := model.SysMenu{}
	sysmenu.MenuName = menuName
	err = sysmenu.FindLastOne()
	lv_err.HasErrAndPanic(err)
	pmenuId := sysmenu.MenuId
	_, err = batis.Exec(tb, "menu_button_create", pmenuId)
	_, err = batis.Exec(tb, "menu_button_retrieve", pmenuId)
	_, err = batis.Exec(tb, "menu_button_update", pmenuId)
	_, err = batis.Exec(tb, "menu_button_delete", pmenuId)
	_, err = batis.Exec(tb, "menu_button_export", pmenuId)
	if err != nil {
		panic(err)
	}
	util.Success(c, nil)
}

// swagger文档
func (w *GenApi) Swagger(c *gin.Context) {
	a := c.Query("a")
	if a == "r" {
		//重新生成文档
		curDir, err := os.Getwd()
		if err != nil {
			util.BuildTpl(c, lv_dto.ERROR_PAGE).WriteTpl(gin.H{
				"desc": "参数错误",
			})
			c.Abort()
			return
		}
		genPath := curDir + "/static/swagger"
		err = w.generateSwaggerFiles(genPath)
		if err != nil {
			util.BuildTpl(c, lv_dto.ERROR_PAGE).WriteTpl(gin.H{
				"desc": "参数错误",
			})
			c.Abort()
			return
		}
	}
	c.Redirect(http.StatusFound, "/static/swagger/index.html")
}

// 自动生成文档 swag init -o static/swagger
func (w *GenApi) generateSwaggerFiles(output string) error {

	cmd := exec.Command("swag", "init -o "+output)
	// 保证关闭输出流
	if err := cmd.Start(); err != nil { // 运行命令
		return err
	}

	return nil
}
func (w *GenApi) GenList(c *gin.Context) {
	var req *vo.GenTablePageReq
	tableService := service.TableService{}

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("生成代码", req).WriteJsonExit()
		return
	}
	rows := make([]model.GenTable, 0)
	result, total, err := tableService.FindPage(req)

	if err == nil && len(result) > 0 {
		rows = result
	}
	util.SuccessPage(c, rows, total)
}

// 删除数据
func (w *GenApi) Remove(c *gin.Context) {
	var req *lv_dto.IdsReq

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Del).Log("生成代码", req).WriteJsonExit()
		return
	}
	tableService := service.TableService{}
	err := tableService.DeleteByIds(req.Ids)

	if err == nil {
		util.SucessResp(c).SetBtype(lv_dto.Buniss_Del).Log("生成代码", req).WriteJsonExit()
	} else {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Del).Log("生成代码", req).WriteJsonExit()
	}
}

// EditSave 修改数据保存
func (w *GenApi) EditSave(c *gin.Context) {
	var req vo.GenTableEditReq

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).SetBtype(lv_dto.Buniss_Edit).Log("生成代码", gin.H{"tableName": req.TableName}).WriteJsonExit()
		return
	}
	tableService := service.TableService{}
	err := tableService.SaveEdit(&req, c)
	if err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).SetBtype(lv_dto.Buniss_Edit).Log("生成代码", gin.H{"tableName": req.TableName}).WriteJsonExit()
		return
	}
	util.SucessResp(c).SetBtype(lv_dto.Buniss_Edit).Log("生成代码", gin.H{"tableName": req.TableName}).WriteJsonExit()
}

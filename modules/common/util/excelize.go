package util

import (
	"fmt"
	"github.com/lostvip-com/lv_framework/lv_log"
	"math/rand"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/xuri/excelize/v2"
)

// 导出表格

var (
	defaultSheetName = "Sheet1" //默认Sheet名称
	defaultHeight    = 25.0     //默认行高度
)

type lzExcelExport struct {
	file      *excelize.File
	sheetName string //可定义默认sheet名称
}

func NewMyExcel() *lzExcelExport {
	return &lzExcelExport{file: createFile(), sheetName: defaultSheetName}
}

func (xls *lzExcelExport) ExportToPath(params []map[string]string, data []map[string]interface{}, path string) (string, error) {
	xls.export(params, data)
	name := createFileName()
	filePath := path + "/" + name
	err := xls.file.SaveAs(filePath)
	return filePath, err
}

func (xls *lzExcelExport) ExportToWeb(c *gin.Context, params []map[string]string, data []map[string]any) {
	xls.export(params, data)
	buffer, _ := xls.file.WriteToBuffer()
	//设置文件类型
	c.Header("Content-Type", "application/vnd.ms-excel;charset=utf8")
	//设置文件名称
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(createFileName()))
	_, _ = c.Writer.Write(buffer.Bytes())
}

// writeTop 设置首行
func (xls *lzExcelExport) writeHeader(params []map[string]string) {
	style := &excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	}
	topStyle, _ := xls.file.NewStyle(style)
	var word = 'A'
	//首行写入
	for _, conf := range params { //首行写入标题，A1 B1 C1
		title := conf["title"]
		width, _ := strconv.ParseFloat(conf["width"], 64)
		line := fmt.Sprintf("%c1", word)
		//设置标题
		_ = xls.file.SetCellValue(xls.sheetName, line, title)
		//列宽
		_ = xls.file.SetColWidth(xls.sheetName, fmt.Sprintf("%c", word), fmt.Sprintf("%c", word), width)
		//设置样式
		_ = xls.file.SetCellStyle(xls.sheetName, line, line, topStyle)
		word++
	}
}

// writeData 写入数据行
func (xls *lzExcelExport) writeData(headers []map[string]string, data []map[string]any) {
	style := &excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	}
	lineStyle, _ := xls.file.NewStyle(style)
	//数据写入
	var j = 2 //数据开始行数
	for i, val := range data {
		//设置行高
		_ = xls.file.SetRowHeight(xls.sheetName, i+1, defaultHeight)
		//逐列写入
		var word = 'A'
		for _, conf := range headers {
			valKey := conf["key"]
			line := fmt.Sprintf("%c%v", word, j)
			//设置值
			_ = xls.file.SetCellValue(xls.sheetName, line, val[valKey])
			//设置样式
			_ = xls.file.SetCellStyle(xls.sheetName, line, line, lineStyle)
			word++
		}
		j++
	}
	//设置行高 尾行
	_ = xls.file.SetRowHeight(xls.sheetName, len(data)+1, defaultHeight)
}

func (xls *lzExcelExport) export(params []map[string]string, data []map[string]interface{}) {
	xls.writeHeader(params)
	xls.writeData(params, data)
}

func createFile() *excelize.File {
	f := excelize.NewFile()
	// 创建一个默认工作表
	sheetName := defaultSheetName
	index, err := f.NewSheet(sheetName)
	lv_err.HasErrAndPanic(err)
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	return f
}

func createFileName() string {
	name := time.Now().Format("2006-01-02-15-04-05")
	rand.NewSource(time.Now().UnixNano())
	return fmt.Sprintf("excle-%v-%v.xlsx", name, rand.Int63n(time.Now().Unix()))
}

func (xls *lzExcelExport) ExportExcelByStruct(titleList []string, data []interface{}, fileName string, sheetName string, c *gin.Context) error {
	xls.file.SetSheetName("Sheet1", sheetName)
	header := make([]string, 0)
	for _, v := range titleList {
		header = append(header, v)
	}
	style := excelize.Style{
		Font: &excelize.Font{
			Color:  "#666666",
			Size:   13,
			Family: "arial",
		},
		Alignment: &excelize.Alignment{
			Vertical:   "center",
			Horizontal: "center",
		},
	}
	rowStyleID, _ := xls.file.NewStyle(&style)
	_ = xls.file.SetSheetRow(sheetName, "A1", &header)
	_ = xls.file.SetRowHeight("Sheet1", 1, 30)
	length := len(titleList)
	headStyle := Letter(length)
	var lastRow string
	var widthRow string
	for k, v := range headStyle {

		if k == length-1 {

			lastRow = fmt.Sprintf("%s1", v)
			widthRow = v
		}
	}
	if err := xls.file.SetColWidth(sheetName, "A", widthRow, 30); err != nil {
		lv_log.Error("错误--", err.Error())
	}
	rowNum := 1
	for _, v := range data {

		t := reflect.TypeOf(v)
		value := reflect.ValueOf(v)
		row := make([]interface {
		}, 0)
		for l := 0; l < t.NumField(); l++ {

			val := value.Field(l).Interface()
			row = append(row, val)
		}
		rowNum++
		err := xls.file.SetSheetRow(sheetName, "A"+strconv.Itoa(rowNum), &row)
		_ = xls.file.SetCellStyle(sheetName, "A"+strconv.Itoa(rowNum), lastRow, rowStyleID)
		if err != nil {
			return err
		}
	}
	disposition := fmt.Sprintf("attachment; filename=%s.xlsx", url.QueryEscape(fileName))
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", disposition)
	c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	return xls.file.Write(c.Writer)
}

// Letter 遍历a-z
func Letter(length int) []string {
	var str []string
	for i := 0; i < length; i++ {
		str = append(str, string(rune('A'+i)))
	}
	return str
}

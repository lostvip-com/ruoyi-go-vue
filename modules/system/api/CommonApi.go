package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_global"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"io"
	"os"
)

type CommonApi struct{}

// DownloadTmp 从临时目录下载，如excell等动态生成的数据（默认下载）
func (w *CommonApi) DownloadTmp(c *gin.Context) {
	fileName := c.Query("fileName")
	filepath := lv_global.Config().GetTmpPath() + "/" + fileName
	file, err := os.Open(filepath)
	defer file.Close()
	lv_err.HasErrAndPanic(err)
	b, _ := io.ReadAll(file)
	c.Writer.Header().Add("Content-Disposition", "attachment")
	c.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetmxls.sheet")
	c.Writer.Write(b)
	c.Abort()
}

// DownloadUpload 从upload目录下载,下载 public/upload 文件头像之类
func (w *CommonApi) DownloadUpload(c *gin.Context) {
	fileName := c.Query("fileName")
	filepath := lv_global.Config().GetUploadPath() + "/" + fileName
	file, err := os.Open(filepath)
	defer file.Close()
	lv_err.HasErrAndPanic(err)
	b, _ := io.ReadAll(file)
	c.Writer.Header().Add("Content-Disposition", "attachment")
	c.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetmxls.sheet")
	c.Writer.Write(b)
	c.Abort()
}

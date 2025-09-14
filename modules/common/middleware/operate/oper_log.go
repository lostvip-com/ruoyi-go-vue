package operate

import (
	"bytes"
	"common/global"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_if"
	"github.com/spf13/cast"
	"io"
	"strings"
	"system/model"
	"system/service"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		//startTime := time.Now()
		// 记录请求参数
		method := c.Request.Method
		if method == "GET" || method == "OPTION" {
			c.Next()
			return //忽略查询类操作
		}
		path := c.Request.URL.Path
		params := c.Request.URL.Query()

		// 从context中获取用户信息
		u, ok := c.Get(global.KEY_GIN_USER_PTR)
		if !ok {
			util.Fail(c, "获取用户信息失败")
		}
		userPtr := u.(*model.SysUser)
		var bodyStr string
		// 读取请求体（需要特殊处理）
		if c.Request.Body != nil {
			body, _ := io.ReadAll(c.Request.Body)
			bodyStr = string(body)
			// 重新设置body，因为读取后会被清空
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		lv_log.Debug("Logger ----------> Request \nMethod:", method, "\n Path:", path, "\n Params:", params, "\n Body:", bodyStr)

		// 处理请求
		c.Next()

		// 记录响应结果
		//duration := time.Since(startTime)
		status := c.Writer.Status()
		var resBody string

		// 判断是否为文件下载请求
		fileHeader := c.Writer.Header().Get("Content-Disposition")
		isFileDownload := strings.Contains(fileHeader, "attachment")
		if isFileDownload {
			resBody = fileHeader
		} else {
			blw := &BodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw
			if blw != nil {
				resBody = blw.body.String()
			}
		}
		// 将用户信息传递给SaveLog方法（如果需要）
		service.GetOperLogServiceInstance().SaveLog(c, status, lv_if.IfEmpty(bodyStr, cast.ToString(params)), resBody, userPtr)
		//lv_log.Debug("Logger ----------> Response \nStatus:", status, "\n Body:", resBody, "\n Duration:", duration, "\n User:", userPtr)
	}
}

type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w BodyLogWriter) Write(b []byte) (int, error) {
	// 检查是否是文件下载响应，并尝试从Content-Disposition头中提取文件名
	//contentDisposition := w.Header().Get("Content-Disposition")
	//if contentDisposition != "" && strings.Contains(contentDisposition, "attachment") {
	//	// 使用正则表达式从Content-Disposition中提取filename参数
	//	re := regexp.MustCompile(`filename[^;=\n]*=((['](?P<fname>[^']+)['])|(["](?P<fname2>[^"]+)["])|(?P<fname3>[^;\n]*))`)
	//	matches := re.FindStringSubmatch(contentDisposition)
	//	if len(matches) > 1 {
	//		for i := 1; i < len(matches); i++ {
	//			if matches[i] != "" && !strings.HasPrefix(matches[i], `"`) && !strings.HasPrefix(matches[i], `'`) {
	//				// 解码URL编码的文件名
	//				if decodedName, err := url.QueryUnescape(matches[i]); err == nil {
	//					// 将文件名记录到日志中
	//					lv_log.Debug("检测到文件下载，文件名: ", decodedName)
	//				}
	//				break
	//			}
	//		}
	//	}
	//}

	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

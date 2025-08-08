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
	"system/model"
	"system/service"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()
		// 记录请求参数
		method := c.Request.Method
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
		// 包装ResponseWriter来捕获响应
		blw := &BodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 记录响应结果
		duration := time.Since(startTime)
		status := c.Writer.Status()
		resBody := blw.body.String()

		// 将用户信息传递给SaveLog方法（如果需要）
		service.GetOperLogServiceInstance().SaveLog(c, status, lv_if.IfEmpty(bodyStr, cast.ToString(params)), resBody, userPtr)
		lv_log.Debug("Logger ----------> Response \nStatus:", status, "\n Body:", resBody, "\n Duration:", duration, "\n User:", userPtr)
	}
}

// 自定义ResponseWriter用于捕获响应体
type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

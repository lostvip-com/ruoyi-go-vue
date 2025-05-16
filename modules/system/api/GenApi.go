package api

import (
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"net/http"
	"system/model"
	"system/service"
)

type GenApi struct {
	BaseApi
}

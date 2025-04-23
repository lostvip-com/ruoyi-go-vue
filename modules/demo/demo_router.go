package demo

import (
	"demo/api"
	"fmt"
	"github.com/lostvip-com/lv_framework/web/router"
)

func init() {
	fmt.Println("############## demo init ################")
	//g1 := router.New( "/demo/form",token.TokenCheck())
	demo := api.DemoController{}
	g0 := router.New("/demo/lv_db")
	//mybatis
	g0.GET("/mybatis1", "", demo.MybatisMap)
	g0.GET("/mybatis2", "", demo.MybatisStruct)
	g0.GET("/mybatis3", "", demo.MybatisStructPage)
	g0.GET("/redis", "", demo.TestRedis)
}

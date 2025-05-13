package api

import (
	"common/global"
	"common/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_cache/lv_redis"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_net"
	"strings"
	"system/service"
	"system/vo"
)

type MonitorApi struct {
}

func (m MonitorApi) CacheHandler(c *gin.Context) {
	var list []vo.CacheVO
	list = append(list, vo.CacheVO{CacheName: global.LoginCacheKey, Remark: "用户信息"})
	list = append(list, vo.CacheVO{CacheName: "sys_config:", Remark: "配置信息"})
	list = append(list, vo.CacheVO{CacheName: global.SysDictCacheKey, Remark: "数据字典"})
	list = append(list, vo.CacheVO{CacheName: "captcha_codes:", Remark: "验证码"})
	list = append(list, vo.CacheVO{CacheName: "repeat_submit:", Remark: "防重提交"})
	list = append(list, vo.CacheVO{CacheName: "rate_limit:", Remark: "限流处理"})
	list = append(list, vo.CacheVO{CacheName: "pwd_err_cnt:", Remark: "密码错误次数"})
	util.Success(c, list)
}
func (m MonitorApi) GetCacheKeysHandler(c *gin.Context) {
	cacheName := c.Param("cacheName")
	redisCache := lv_redis.GetInstance(0)
	keys, _, err := redisCache.Scan(0, cacheName+"*", global.ScanCountMax)
	if err != nil {
		util.Fail(c, err.Error())
	}
	util.Success(c, keys)
}

func (m MonitorApi) GetCacheValueHandler(c *gin.Context) {
	cacheName := c.Param("cacheName")
	cacheKey := c.Param("cacheKey")
	redisCache := lv_redis.GetInstance(0)
	value, err := redisCache.Get(cacheKey)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	var cache = vo.CacheVO{
		CacheName:  cacheName,
		CacheKey:   cacheKey,
		CacheValue: value,
	}
	util.Success(c, cache)
}

func (m MonitorApi) ClearCacheNameHandler(c *gin.Context) {
	cacheName := c.Param("cacheName")
	redisCache := lv_redis.GetInstance(0)
	keys, _, err := redisCache.Scan(0, cacheName+"*", global.ScanCountMax)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	for i := range keys {
		err := redisCache.Del(keys[i])
		if err != nil {
			lv_log.Error("ClearCacheNameHandler error:", err.Error())
		}
	}
	util.Success(c, nil)
}

func (m MonitorApi) ClearCacheKeyHandler(c *gin.Context) {
	cacheKey := c.Param("cacheKey")
	redisCache := lv_redis.GetInstance(0)
	err := redisCache.Del(cacheKey)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}

func (m MonitorApi) ClearCacheAllHandler(c *gin.Context) {
	redisCache := lv_redis.GetInstance(0)
	keys, _, err := redisCache.Scan(0, "*", global.ScanCountMax)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	for i := range keys {
		err := redisCache.Del(keys[i])
		if err != nil {
			return
		}
	}
	util.Success(c, nil)
}

func (m MonitorApi) ServerInfo(c *gin.Context) {
	ip := lv_net.GetRemoteClientIp(c.Request)
	var server service.MonitorService
	var info = vo.ServerInfo{}
	info.Cpu = server.GetCpu()
	info.GoInfo = server.GetGoInfo()
	info.Mem = server.GetMem()
	info.Sys = server.GetSys(ip)
	info.SysFile = server.GetSysFile()
	util.Success(c, info)
}

func (m MonitorApi) ListOnLine(c *gin.Context) {
	var ipaddr = c.DefaultQuery("ipaddr", "")
	var userName = c.DefaultQuery("userName", "")
	key := global.LoginCacheKey + "*"
	redisCache := lv_redis.GetInstance(0)
	keyList, _, _ := redisCache.Scan(0, key, 0)
	var rows []vo.LoginUserCache
	for i := range keyList {
		keyString := keyList[i]
		result, _ := redisCache.Get(keyString)
		var loginUser vo.LoginUserCache
		json.Unmarshal([]byte(result), &loginUser)
		rows = append(rows, loginUser)
	}
	if userName != "" || ipaddr != "" { //按条件搜索
		rows = *m.FindSearchTarget(rows, userName, ipaddr)
	}
	util.SuccessPage(c, rows, int64(len(rows)))
}

func (m MonitorApi) FindSearchTarget(rows []vo.LoginUserCache, userName string, ipaddr string) *[]vo.LoginUserCache {
	var search = make([]vo.LoginUserCache, 0)
	for i := range rows {
		row := &rows[i]
		if userName != "" || row.UserName == userName {
			if strings.Contains(userName, row.UserName) {
				search = append(search, *row)
			}
		}
		if ipaddr != "" || row.Ipaddr == ipaddr {
			if strings.Contains(ipaddr, row.Ipaddr) {
				search = append(search, *row)
			}
		}
	}
	return &search
}

func (m MonitorApi) DetectOnLine(context *gin.Context) {

}

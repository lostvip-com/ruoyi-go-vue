package system

import (
	auth2 "common/middleware/auth"
	"github.com/lostvip-com/lv_framework/web/router"
	"system/api"
)

func init() {

	common := router.New("/common", auth2.TokenCheck(), auth2.PermitCheck)
	commonApi := api.CommonApi{}
	common.GET("/download", "", commonApi.DownloadTmp)
	common.GET("/downloadUpload", "", commonApi.DownloadUpload)
	//系统配置
	system := router.New("/system", auth2.TokenCheck(), auth2.PermitCheck)
	config := api.ConfigApi{}
	system.GET("/config/:configId", "system:config:list", config.GetConfigInfo)
	system.GET("/config/list", "system:config:list", config.ListAjax)
	system.GET("/config/configKey/:configKey", "", config.GetConfigKey)
	system.POST("/config", "system:config:add", config.AddSave)
	system.PUT("/config", "system:config:edit", config.EditSave)
	system.DELETE("/config/:configIds", "system:config:remove", config.Remove)
	system.POST("/config/export", "system:config:export", config.Export)
	// 字典类型参数路由
	dictType := api.DictTypeApi{}
	system.GET("/dict/type/list", "system:dict:list", dictType.ListAjax)
	system.GET("/dict/type/:dictId", "", dictType.GetTypeDict)
	system.POST("/dict/type", "system:dict:add", dictType.AddSave)
	system.PUT("/dict/type", "system:dict:edit", dictType.EditSave)
	system.DELETE("/dict/type/:dictIds", "system:dict:remove", dictType.Remove)
	system.POST("/dict/type/export", "system:dict:export", dictType.Export)
	//system.POST("/dict/type/checkDictTypeUniqueAll", "system:dict:view", dictType.CheckDictTypeUniqueAll)
	//system.POST("/dict/type/checkDictTypeUnique", "system:dict:view", dictType.CheckDictTypeUnique)
	system.GET("/dict/type/optionselect", "system:dict:view", dictType.GetOptionSelect)
	// 字典内容参数路由
	dictData := api.DictDataApi{}
	system.GET("/dict/data/list", "", dictData.ListAjax)
	system.GET("/dict/data/type/:dictType", "", dictData.GetDictDataByDictType)
	system.POST("/dict/data", "system:dict:add", dictData.AddSave)
	system.PUT("/dict/data", "system:dict:edit", dictData.EditSave)
	system.DELETE("/dict/data/:dictCodes", "system:dict:remove", dictData.Remove)
	system.POST("/dict/data/export", "system:dict:export", dictData.Export)
	//dept
	deptApi := api.DeptApi{}
	system.GET("/dept/:deptId", "system:dept:list", deptApi.GetDept)
	system.GET("/dept/list", "system:dept:list", deptApi.ListAjax)
	system.GET("dept/list/exclude/:deptId", "system:dept:list", deptApi.ExcludeDept)
	system.POST("/dept", "system:dept:add", deptApi.AddSave)
	system.PUT("/dept", "system:dept:edit", deptApi.EditSave)
	system.DELETE("/dept/:deptId", "system:dept:remove", deptApi.Remove)
	// 用户管理路由
	userApi := api.UserApi{}
	system.GET("/user/:userId", "system:user:list", userApi.GetUserInfo)
	system.GET("/user/list", "system:user:list", userApi.ListAjax)
	system.POST("/user", "system:user:add", userApi.AddSave)
	system.PUT("/user", "system:user:edit", userApi.EditSave)
	system.DELETE("/user/:userIds", "system:user:remove", userApi.Remove)
	system.PUT("/user/resetPwd", "system:user:resetPwd", userApi.ResetPwdSave)
	system.PUT("/user/changeStatus", "system:user:edit", userApi.ChangeStatus)
	system.PUT("/user/authRole", "system:user:edit", userApi.PutAuthUserRoleIds)
	system.GET("/user/authRole/:userId", "", userApi.GetAuthUserRole)
	system.GET("/user/deptTree", "", userApi.GetUserDeptTree)
	system.POST("/user/importData", "system:user:add", userApi.ImportUserData)
	system.POST("/user/importTemplate", "system:user:add", userApi.ImportTemplate)
	system.POST("/user/export", "system:userApi:export", userApi.Export)
	// 个人中心路由
	profile := api.ProfileApi{}
	system.GET("userApi/profile", "", profile.Profile)
	system.POST("/profile/update", "", profile.Update)
	system.POST("/profile/resetSavePwd", "", profile.UpdatePassword)
	system.POST("/profile/checkPhoneOK", "", profile.CheckPhoneOK)
	system.POST("/profile/checkEmailOK", "", profile.CheckEmailOK)
	system.POST("/profile/checkUserNameOK", "", profile.CheckUserNameOK)
	system.POST("/profile/checkPassword", "", profile.CheckPassword)
	system.POST("/profile/updateAvatar", "", profile.UpdateAvatar)
	// 角色路由
	roleApi := api.RoleApi{}
	system.GET("/role/list", "system:role:list", roleApi.ListAjax)
	system.GET("/role/:roleId", "system:role:detail", roleApi.GetRoleInfo)
	system.GET("/role/optionselect", "system:role:view", roleApi.GetRoleOptionSelect)
	system.POST("/role", "system:role:add", roleApi.AddSave)
	system.PUT("/role", "system:role:edit", roleApi.EditSave)
	system.POST("/role/changeStatus", "system:role:edit", roleApi.ChangeStatus)
	system.PUT("/role/dataScope", "system:role:edit", roleApi.PutDataScope)
	system.DELETE("/role/:roleIds", "system:role:remove", roleApi.Remove)

	system.PUT("/role/authUser/cancel", "system:role:view", roleApi.Cancel)
	system.PUT("/role/authUser/cancelAll", "system:role:view", roleApi.CancelAll)
	system.PUT("/role/authUser/selectAll", "system:role:view", roleApi.AuthRoleToUsers)
	system.GET("/role/authUser/allocatedList", "system:role:view", roleApi.AllocatedList)
	system.GET("/role/authUser/unallocatedList", "system:role:view", roleApi.GetUnAllocatedList)
	system.GET("/role/deptTree/:roleId", "system:role:view", roleApi.GetDeptTreeRole)
	// 菜单路由
	menuApi := api.MenuApi{}
	system.GET("/menu/:menuId", "system:menu:detail", menuApi.GetMenuInfo)
	system.GET("/menu/list", "system:menu:list", menuApi.ListMenu)
	system.GET("/menu/treeselect", "system:menu:list", menuApi.GetTreeSelect)
	system.GET("/menu/roleMenuTreeselect/:roleId", "system:menu:list", menuApi.TreeSelectByRole)
	system.POST("/menu", "system:menu:add", menuApi.AddSave)
	system.PUT("/menu", "system:menu:edit", menuApi.EditSave)
	system.DELETE("/menu", "system:menu:remove", menuApi.Remove)
	//system.GET("/menu/treeData", "", menuApi.MenuTreeData)
	// 岗位路由
	postApi := api.PostApi{}
	system.GET("/post/:postId", "", postApi.GetPostInfo)
	system.GET("/post/list", "system:post:list", postApi.ListAjax)
	system.POST("/post", "system:post:add", postApi.AddSave)
	system.PUT("/post", "system:post:edit", postApi.EditSave)
	system.DELETE("/post/:postIds", "system:post:remove", postApi.Remove)
	system.POST("/post/export", "system:post:export", postApi.Export)
	system.GET("/optionselect", "system:post:list", postApi.GetPostOptionSelect)
}

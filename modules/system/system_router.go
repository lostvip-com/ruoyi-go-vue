package system

import (
	auth2 "common/middleware/auth"
	"github.com/lostvip-com/lv_framework/web/router"
	"system/api"
)

func init() {

	gcommon := router.New("/common", auth2.TokenCheck(), auth2.PermitCheck)
	common := api.CommonController{}
	gcommon.GET("/download", "", common.DownloadTmp)
	gcommon.GET("/downloadUpload", "", common.DownloadUpload)
	//系统配置
	system := router.New("/system", auth2.TokenCheck(), auth2.PermitCheck)
	config := api.ConfigController{}
	system.GET("/config", "system:config:view", config.List)
	system.POST("/config/list", "system:config:list", config.ListAjax)
	system.GET("/config/add", "system:config:add", config.Add)
	system.POST("/config/add", "system:config:add", config.AddSave)
	system.POST("/config/remove", "system:config:remove", config.Remove)
	system.GET("/config/edit", "system:config:edit", config.Edit)
	system.POST("/config/edit", "system:config:edit", config.EditSave)
	system.POST("/config/export", "system:config:export", config.Export)
	system.POST("/config/checkConfigKeyUnique", "system:config:view", config.CheckConfigKeyUnique)
	// 字典类型参数路由
	dictType := api.DictTypeController{}
	system.POST("/dict/type/list", "system:dict:list", dictType.ListAjax)
	system.POST("/dict/type/add", "system:dict:add", dictType.AddSave)
	system.POST("/dict/type/remove", "system:dict:remove", dictType.Remove)
	system.GET("/dict/type/remove", "system:dict:remove", dictType.Remove)
	system.POST("/dict/type/edit", "system:dict:edit", dictType.EditSave)
	system.POST("/dict/type/export", "system:dict:export", dictType.Export)
	system.POST("/dict/type/checkDictTypeUniqueAll", "system:dict:view", dictType.CheckDictTypeUniqueAll)
	system.POST("/dict/type/checkDictTypeUnique", "system:dict:view", dictType.CheckDictTypeUnique)
	system.GET("/dict/type/treeData", "system:dict:view", dictType.TreeData)
	// 字典内容参数路由
	dictData := api.DictDataController{}
	system.POST("/dict/data/list", "system:dict:view", dictData.ListAjax)
	system.GET("/dict/data/add", "system:dict:add", dictData.Add)
	system.POST("/dict/data/add", "system:dict:add", dictData.AddSave)
	system.POST("/dict/data/remove", "system:dict:remove", dictData.Remove)
	system.POST("/dict/data/edit", "system:dict:edit", dictData.EditSave)
	system.POST("/dict/data/export", "system:dict:export", dictData.Export)
	//dept
	deptContr := api.DeptController{}
	system.POST("/dept/list", "system:dept:list", deptContr.ListAjax)
	system.POST("/dept/add", "system:dept:add", deptContr.AddSave)
	system.POST("/dept/remove", "system:dept:remove", deptContr.Remove)
	system.POST("/dept/edit", "system:dept:edit", deptContr.EditSave)
	system.GET("/dept/treeData", "system:dept:view", deptContr.TreeData)
	system.GET("/dept/roleDeptTreeData", "system:dept:view", deptContr.RoleDeptTreeData)
	// 用户管理路由
	user := api.UserApi{}
	system.POST("/user/list", "system:user:list", user.ListAjax)
	system.POST("/user/add", "system:user:add", user.AddSave)
	system.POST("/user/remove", "system:user:remove", user.Remove)
	system.POST("/user/edit", "system:user:edit", user.EditSave)
	system.POST("/user/export", "system:user:export", user.Export)
	system.POST("/user/resetPwd", "system:user:resetPwd", user.ResetPwdSave)
	system.POST("/user/changeStatus", "system:user:edit", user.ChangeStatus)
	// 个人中心路由
	profile := api.ProfileController{}
	system.GET("user/profile", "", profile.Profile)
	system.POST("/profile/update", "", profile.Update)
	system.POST("/profile/resetSavePwd", "", profile.UpdatePassword)
	system.POST("/profile/checkPhoneOK", "", profile.CheckPhoneOK)
	system.POST("/profile/checkEmailOK", "", profile.CheckEmailOK)
	system.POST("/profile/checkUserNameOK", "", profile.CheckUserNameOK)
	system.POST("/profile/checkPassword", "", profile.CheckPassword)
	system.POST("/profile/updateAvatar", "", profile.UpdateAvatar)
	// 角色路由
	roleController := api.RoleController{}
	system.POST("/role/list", "system:role:list", roleController.ListAjax)
	system.POST("/role/add", "system:role:add", roleController.AddSave)
	system.POST("/role/remove", "system:role:remove", roleController.Remove)
	system.POST("/role/edit", "system:role:edit", roleController.EditSave)
	system.POST("/role/changeStatus", "system:role:edit", roleController.ChangeStatus)
	system.POST("/role/authDataScope", "system:role:view", roleController.AuthDataScopeSave)
	system.POST("/role/allocatedList", "system:role:view", roleController.AllocatedList)
	system.POST("/role/unallocatedList", "system:role:view", roleController.UnallocatedList)
	system.POST("/role/selectAll", "system:role:view", roleController.SelectAll)
	system.POST("/role/cancel", "system:role:view", roleController.Cancel)
	system.POST("/role/cancelAll", "system:role:view", roleController.CancelAll)

	// 菜单路由
	menuController := api.MenuApi{}
	system.POST("/menu/list", "system:menu:list", menuController.ListAjax)

	system.POST("/menu/add", "system:menu:add", menuController.AddSave)
	system.GET("/menu/remove", "system:menu:remove", menuController.Remove)
	system.POST("/menu/remove", "system:menu:remove", menuController.Remove)
	system.POST("/menu/edit", "system:menu:edit", menuController.EditSave)
	system.GET("/menu/roleMenuTreeData", "system:menu:view", menuController.RoleMenuTreeData)
	system.GET("/menu/treeData", "system:menu:view", menuController.MenuTreeData)
	// 岗位路由
	postController := api.PostController{}
	system.POST("/post/post/list", "system:post:list", postController.ListAjax)
	system.GET("/post/add", "system:post:add", postController.Add)
	system.POST("/post/add", "system:post:add", postController.AddSave)
	system.POST("/post/remove", "system:post:remove", postController.Remove)
	system.GET("/post/edit", "system:post:edit", postController.Edit)
	system.POST("/post/edit", "system:post:edit", postController.EditSave)
	system.POST("/post/export", "system:post:export", postController.Export)
	system.POST("/post/isPostCodeExist", "system:post:add", postController.IsPostCodeExist)
}

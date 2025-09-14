package global

const ErrTimes2Lock = 20

// KeyOssUrl 获取oss地址Key
const KeyOssUrl = "sys.resource.url"

// 配置文件目录搜索顺序
var BaseFilePathArr = []string{".", "./resources", "../", "../resources"}

const (
	TimeFormat = "2006-01-02 15:04:05"
	DateFormat = "2006-01-02"
)

// cache key
const (
	SYS_DICT_KEY         = "sys_dict:"
	CAPTCHA_CODE_KEY     = "captcha_codes:"
	REPEAT_SUBMIT_KEY    = "repeat_submit:" //防重提交 redis key
	RATE_LIMIT_KEY       = "rate_limit:"
	PWD_ERR_CNT_KEY      = "pwd_err_cnt:" // 登录账户密码错误次数 redis key
	LoginCacheKey        = "login_tokens:"
	CaptchaCodesKey      = "captcha_codes:"
	SysDictCacheKey      = "sys_dict:"
	SysConfigCacheKey    = "sys_config:"
	PwdErrCntCacheKey    = "pwd_err_cnt:"
	RepeatSubmitCacheKey = "repeat_submit:"
	RateLimitCacheKey    = "rate_limit:"
	ScanCountMax         = 1000
)

// yaml key
const (
	KEY_CACHE_TYPE = "application.cache-type"
)

// gin 协程内共享的变量 set get 使用
const (
	KEY_GIN_USER_PTR = "user"
	KEY_GIN_USERNAME = "userName"
	KEY_GIN_USER_ID  = "userId"
	KEY_GIN_DEPT_ID  = "deptId"
	KEY_GIN_BIZ_TYPE = "bizCode" //操作编码，用于业务日志记录
	KEY_GIN_IN_PARAM = "inParam"
)

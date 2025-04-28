package global

import "os"

const Layout = "2006-01-02 15:04:05" //时间常量
const ErrTimes2Lock = 20

// KeyOssUrl 获取oss地址Key
const KeyOssUrl = "sys.resource.url"

const JQ_BE_NO = "false"
const JQ_BE_OK = "true"

const DIR_DIST_CODE = ""
const DIR_TPL_CODE_GEN = "resources" + string(os.PathSeparator) + "tpl_gen"

// 登录
const SYS_DICT_KEY = "sys_dict:"
const LOGIN_TOKEN_KEY = "login_tokens:"
const CAPTCHA_CODE_KEY = "captcha_codes:"
const REPEAT_SUBMIT_KEY = "repeat_submit:" //防重提交 redis key
const RATE_LIMIT_KEY = "rate_limit:"
const PWD_ERR_CNT_KEY = "pwd_err_cnt:" // 登录账户密码错误次数 redis key

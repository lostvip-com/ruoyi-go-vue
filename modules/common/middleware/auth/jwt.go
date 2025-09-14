package auth

import (
	"common/global"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/spf13/cast"
	"net/http"
	"strings"
	"system/service"
	"time"
)

const Authorization = "Authorization"
const Bearer = "Bearer "
const Secret = "3Bde3BGEbYqtqyEUzW3ry8jKFcaPH17fRmTmqE7MDr05Lwj95uruRKrrkb44TJ4s"
const JwtTtl = 12 * 3600 //12h

func CreateToken(UserName string, UserId int, DeptId int, uuid string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name": UserName,
		"user_id":   UserId,
		"dept_id":   DeptId,
		"exp":       time.Now().Unix() + int64(JwtTtl),
		"iss":       "lostvip.com",
		"uuid":      uuid,
	})

	mySigningKey := []byte(Secret)
	tokenId, err := token.SignedString(mySigningKey)
	if err != nil {
		tokenId = "error_" + cast.ToString(time.Now().Unix())
	}
	return tokenId
}

func verifyToken(tokenStr string) (*jwt.Token, error) {
	mySigningKey := []byte(Secret)
	tokenStr = strings.ReplaceAll(tokenStr, Bearer, "")
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
}

func TokenCheck() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		//根据实际情况取TOKEN, 这里从request header取
		header := ctx.Request.Header
		tokenStr := header.Get(Authorization)
		if len(tokenStr) < 1 {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": "参数错误", "code": http.StatusInternalServerError})
			return
		}
		token, err := verifyToken(tokenStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": "认证失败", "code": http.StatusUnauthorized})
			return
		}
		uuid := token.Claims.(jwt.MapClaims)["uuid"].(string)
		yes := service.GetSessionServiceInstance().IsSignedIn(uuid)
		if !yes {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": "token timeout", "code": http.StatusUnauthorized})
			return
		}

		//续期会话（使用redis）
		service.GetSessionServiceInstance().Refresh(uuid)
		userId := token.Claims.(jwt.MapClaims)["user_id"]
		deptId := token.Claims.(jwt.MapClaims)["dept_id"]
		username := token.Claims.(jwt.MapClaims)["user_name"]
		userPtr, err := service.GetUserServiceInstance().GetProfile(uuid)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": "USER NOT FOUND!", "code": http.StatusUnauthorized})
			return
		}
		// 此处已经通过了, 可以把Claims中的有效信息拿出来放入上下文使用
		ctx.Set(global.KEY_GIN_USER_ID, userId)    //供api使用
		ctx.Set(global.KEY_GIN_USERNAME, username) //供api使用
		ctx.Set(global.KEY_GIN_USER_ID, userId)
		ctx.Set(global.KEY_GIN_DEPT_ID, deptId)
		ctx.Set(global.KEY_GIN_USERNAME, username)
		ctx.Set(global.KEY_GIN_USER_PTR, userPtr)
		tokenId, err0 := GetJwtUuid(ctx)
		lv_err.HasErrAndPanic(err0)
		ctx.Set("tokenId", tokenId) //给后面的接口使用，避免重复
		//v, ok := ctx.Get(global.KEY_GIN_USER_PTR)
		//lv_log.Debug("token check ", v, ok)
		ctx.Next()
	}
}

func GetJwtUuid(ctx *gin.Context) (string, error) {
	tokenStr := ctx.Request.Header.Get(Authorization)
	if len(tokenStr) <= 0 {
		return "", nil
	}
	token, err := verifyToken(tokenStr)
	if err != nil {
		return "", err
	}
	uuid := token.Claims.(jwt.MapClaims)["uuid"]
	return uuid.(string), nil

}

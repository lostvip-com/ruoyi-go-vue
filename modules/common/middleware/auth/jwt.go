package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/cast"
	"net/http"
	"strings"
	"system/service"
	"time"
)

const Authorization = "Authorization"
const Bearer = "Bearer "
const Secret = "3Bde3BGEbYqtqyEUzW3ry8jKFcaPH17fRmTmqE7MDr05Lwj95uruRKrrkb44TJ4s"
const JwtTtl = 43200

func CreateToken(UserName string, UserId int64, DeptId int64, uuid string) string {
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
		// 根据实际情况取TOKEN, 这里从request header取
		header := ctx.Request.Header
		tokenStr := header.Get(Authorization)
		if len(tokenStr) < 1 {
			ctx.JSON(http.StatusOK, gin.H{
				"msg":  "参数错误",
				"code": http.StatusInternalServerError,
			})
			ctx.Abort()
			return
		}
		token, err := verifyToken(tokenStr)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"msg":  "认证失败",
				"code": http.StatusUnauthorized,
			})
			ctx.Abort()
			return
		}
		uuid := token.Claims.(jwt.MapClaims)["uuid"].(string)
		yes := service.GetSessionServiceInstance().IsSignedIn(uuid)
		if !yes {
			ctx.JSON(http.StatusOK, gin.H{
				"msg":  "token timeout",
				"code": http.StatusUnauthorized,
			})
			ctx.Abort()
			return
		}
		//续期会话（使用redis）
		service.GetSessionServiceInstance().Refresh(uuid)

		userId := token.Claims.(jwt.MapClaims)["user_id"]
		deptId := token.Claims.(jwt.MapClaims)["dept_id"]
		userName := token.Claims.(jwt.MapClaims)["user_name"]
		// 此处已经通过了, 可以把Claims中的有效信息拿出来放入上下文使用
		ctx.Set("userId", userId)
		ctx.Set("deptId", deptId)
		ctx.Set("userName", userName)
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

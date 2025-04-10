package utils

import (
	"takeout/common/constant"
	"takeout/common/errs"
	"time"

	"takeout/common/global"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateToken 生成JWT令牌 - 保留原有函数以兼容现有代码
func GenerateToken(claimName, claimData string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(global.Config.JWT.AdminTTL) * time.Second).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		claimName: claimData,
		"exp":     expirationTime,
	})
	// 获得签名后的完整token
	signedToken, err := token.SignedString([]byte(global.Config.JWT.AdminSecretKey))
	return signedToken, err
}

// ParseToken 解析JWT令牌
func ParseToken(tokenStr, claimName string) (string, error) {
	// 解析 token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		// jwt.SigningMethodHS256 是 jwt.SigningMethodHMAC 的一个具体实现
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.New(constant.CodeJWTParseError, constant.MsgJWTUnKnownSigningMethod)
		}
		return []byte(global.Config.JWT.AdminSecretKey), nil
	})
	// 错误处理
	if err != nil {
		return "", errs.Wrap(err, constant.CodeJWTParseError, constant.MsgJWTParseFail)
	}
	// 获取携带信息
	if claims, isMap := token.Claims.(jwt.MapClaims); isMap && token.Valid {
		// ☆ JWT 数字解析的默认类型是 float64，不能直接用 .(int)
		// 所以 id 还是转化成 string 再写入 token 比较好
		claimsData, ok := claims[claimName].(string)
		// fmt.Println(claims)
		if !ok {
			return "", errs.New(constant.CodeJWTParseError, constant.MsgGetAccountInfoFail)
		}
		return claimsData, nil
	}
	return "", errs.New(constant.CodeJWTParseError, constant.MsgJWTParseFail)
}

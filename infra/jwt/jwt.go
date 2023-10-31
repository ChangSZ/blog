package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/ChangSZ/blog/infra/conf"
	"github.com/ChangSZ/blog/infra/log"
	"github.com/ChangSZ/blog/infra/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
)

type JwtParam struct {
	DefaultIss      string
	DefaultAudience string
	DefaultJti      string
	SecretKey       string
	TokenKey        string
	TokenLife       time.Duration
	RedisCache      *redis.Client
}

func (jp *JwtParam) SetTokenKey(tk string) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.TokenKey
		jp.TokenKey = tk
		return i
	}
}

func (jp *JwtParam) SetTokenLife(tl time.Duration) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.TokenLife
		jp.TokenLife = tl
		return i
	}
}

func (jp *JwtParam) SetDefaultIss(iss string) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.DefaultIss
		jp.DefaultIss = iss
		return i
	}
}

func (jp *JwtParam) SetDefaultAudience(ad string) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.DefaultAudience
		jp.DefaultAudience = ad
		return i
	}
}

func (jp *JwtParam) SetDefaultJti(jti string) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.DefaultJti
		jp.DefaultJti = jti
		return i
	}
}

func (jp *JwtParam) SetDefaultSecretKey(sk string) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.SecretKey
		jp.SecretKey = sk
		return i
	}
}

func (jp *JwtParam) SetRedisCache(rc *redis.Client) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.RedisCache
		jp.RedisCache = rc
		return i
	}
}

var jwtParam *JwtParam

func (jp *JwtParam) JwtInit(options ...func(jp *JwtParam) interface{}) error {
	q := &JwtParam{
		DefaultJti:      conf.JWTJTI,
		DefaultAudience: conf.JWTAUDIENCE,
		DefaultIss:      conf.JWTISS,
		SecretKey:       conf.JWTSECRETKEY,
		TokenLife:       conf.JWTTOKENLIFE,
		TokenKey:        conf.JWTTOKENKEY,
	}
	for _, option := range options {
		option(q)
	}
	jwtParam = q
	return nil
}

func CreateToken(ctx context.Context, userIdString string) (token string, err error) {
	//	iss: jwt签发者
	//	sub: jwt所面向的用户
	//	aud: 接收jwt的一方
	//	exp: jwt的过期时间，这个过期时间必须要大于签发时间
	//	nbf: 定义在什么时间之前，该jwt都是不可用的.
	//	iat: jwt的签发时间
	//	jti: jwt的唯一身份标识，主要用来作为一次性token,从而回避重放攻击。

	tk := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	//claims["exp"] = time.Now().Add(time.Hour * time.Duration(72)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["iss"] = jwtParam.DefaultIss
	claims["sub"] = userIdString
	claims["aud"] = jwtParam.DefaultAudience
	claims["jti"] = utils.Md5(jwtParam.DefaultJti + jwtParam.DefaultIss)
	tk.Claims = claims

	SecretKey := jwtParam.SecretKey
	tokenString, err := tk.SignedString([]byte(SecretKey))
	if err != nil {
		log.WithTrace(ctx).Errorf("签名出错: %v", err)
		return "", err
	}

	err = jwtParam.RedisCache.Set(jwtParam.TokenKey+userIdString, tokenString, jwtParam.TokenLife).Err()
	if err != nil {
		log.WithTrace(ctx).Errorf("创建token出错: %v", err)
		return "", err
	}

	return tokenString, nil
}

func ParseToken(ctx context.Context, myToken string) (userId string, err error) {

	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtParam.SecretKey), nil
	})
	if err != nil {
		log.WithTrace(ctx).Errorf("解析token出错: %v", err)
		return "", err
	}

	if !token.Valid {
		log.WithTrace(ctx).Error("token无效")
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)

	sub, ok := claims["sub"].(string)
	if !ok {
		log.WithTrace(ctx).Errorf("claims断言失败: %v", err)
		return "", errors.New("claims断言失败")
	}

	res, err := jwtParam.RedisCache.Get(jwtParam.TokenKey + sub).Result()
	if err != nil {
		log.WithTrace(ctx).Errorf("从redis获取token失败: %v", err)
		return "", err
	}

	if res == "" || res != myToken {
		log.WithTrace(ctx).Error("token无效")
		return "", errors.New("token无效")
	}

	//refresh the token life time
	err = jwtParam.RedisCache.Set(jwtParam.TokenKey+sub, myToken, jwtParam.TokenLife).Err()
	if err != nil {
		log.WithTrace(ctx).Errorf("创建token出错: %v", err)
		return "", err
	}

	return sub, nil
}

func UnsetToken(ctx context.Context, myToken string) (bool, error) {
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtParam.SecretKey), nil
	})
	if err != nil {
		log.WithTrace(ctx).Errorf("解析token出错: %v", err)
		return false, err
	}
	claims := token.Claims.(jwt.MapClaims)

	sub, ok := claims["sub"].(string)
	if !ok {
		log.WithTrace(ctx).Errorf("claims断言失败: %v", err)
		return false, errors.New("claims断言失败")
	}
	err = jwtParam.RedisCache.Del(jwtParam.TokenKey + sub).Err()
	if err != nil {
		log.WithTrace(ctx).Errorf("unset token出错: %v", err)
		return false, err
	}
	return true, nil
}

package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/conf"
	"github.com/ChangSZ/blog/infra/gin/api"
	"github.com/ChangSZ/blog/infra/jwt"
	"github.com/ChangSZ/blog/infra/log"
	"github.com/ChangSZ/blog/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/mojocn/base64Captcha"
	"golang.org/x/crypto/bcrypt"
)

type ConsoleAuth interface {
	Register(*gin.Context)
	AuthRegister(*gin.Context)
	Login(*gin.Context)
	AuthLogin(*gin.Context)
	Logout(*gin.Context)
	DelCache(*gin.Context)
}
type Auth struct {
}

func NewAuth() ConsoleAuth {
	return &Auth{}
}

// customizeRdsStore An object implementing Store interface
type customizeRdsStore struct {
	redisClient *redis.Client
}

// customizeRdsStore implementing Set method of  Store interface
func (s *customizeRdsStore) Set(id string, value string) error {
	return s.redisClient.Set(id, value, time.Minute*10).Err()
}

// customizeRdsStore implementing Get method of  Store interface
func (s *customizeRdsStore) Get(id string, clear bool) string {
	val, err := s.redisClient.Get(id).Result()
	if err != nil {
		log.Error(err)
		return val
	}
	if clear {
		err := s.redisClient.Del(id).Err()
		if err != nil {
			log.Error(err)
			return val
		}
	}
	return val
}

func (s *customizeRdsStore) Verify(id, answer string, clear bool) bool {
	v := s.Get(id, clear)
	return v == answer
}

func (c *Auth) Register(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	cnt, err := service.GetUserCnt(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 400001004, nil)
		return
	}
	if cnt >= int64(conf.Cnf.UserCnt) {
		log.WithTrace(ctx).Infof("User cnt(%v) beyond expectation(%v)", cnt, conf.Cnf.UserCnt)
		appG.Response(http.StatusOK, 407000015, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (c *Auth) AuthRegister(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	requestJson, exists := ctx.Get("json")
	if !exists {
		log.WithTrace(ctx).Error("get request_params from context fail")
		appG.Response(http.StatusOK, 401000004, nil)
		return
	}
	ar, ok := requestJson.(common.AuthRegister)
	if !ok {
		log.WithTrace(ctx).Error("request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	cnt, err := service.GetUserCnt(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 400001004, nil)
		return
	}
	if cnt >= int64(conf.Cnf.UserCnt) {
		log.WithTrace(ctx).Infof("User cnt(%v) beyond expectation(%v)", cnt, conf.Cnf.UserCnt)
		appG.Response(http.StatusOK, 400001004, nil)
		return
	}
	_, err = service.UserStore(ctx, ar)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 407000016, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (c *Auth) Login(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	customStore := &customizeRdsStore{conf.CacheClient}
	driverDigit := &base64Captcha.DriverDigit{
		Height:   80,
		Width:    240,
		Length:   5,
		MaxSkew:  0.7,
		DotCount: 80,
	}

	captcha := base64Captcha.NewCaptcha(driverDigit, customStore)
	// 生成验证码
	id, b64s, err := captcha.Generate()
	if err != nil {
		appG.Response(http.StatusOK, 500, err)
		return
	}

	data := make(map[string]interface{})
	data["key"] = id
	data["png"] = b64s
	appG.Response(http.StatusOK, 0, data)
	return
}

func (c *Auth) AuthLogin(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	requestJson, exists := ctx.Get("json")
	if !exists {
		log.WithTrace(ctx).Error("get request_params from context fail")
		appG.Response(http.StatusOK, 401000004, nil)
		return
	}
	al, ok := requestJson.(common.AuthLogin)
	if !ok {
		log.WithTrace(ctx).Error("request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	customStore := &customizeRdsStore{conf.CacheClient}
	verifyResult := customStore.Verify(al.CaptchaKey, al.Captcha, true)
	if !verifyResult {
		log.WithTrace(ctx).Error("captcha is error")
		appG.Response(http.StatusOK, 407000008, nil)
		return
	}

	user, err := service.GetUserByEmail(ctx, al.Email)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 407000010, nil)
		return
	}
	if user.ID <= 0 {
		log.WithTrace(ctx).Error("Can get user")
		appG.Response(http.StatusOK, 407000010, nil)
		return
	}

	password := []byte(al.Password)
	hashedPassword := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 407000010, nil)
		return
	}

	userIdStr := strconv.Itoa(user.ID)
	token, err := jwt.CreateToken(ctx, userIdStr)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 407000011, nil)
		return
	}
	appG.Response(http.StatusOK, 0, token)
	return
}

func (c *Auth) Logout(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	token, exist := ctx.Get("token")
	if !exist || token == "" {
		log.WithTrace(ctx).Error("Can not get token")
		appG.Response(http.StatusOK, 400001004, nil)
		return
	}
	_, err := jwt.UnsetToken(ctx, token.(string))
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 407000014, nil)
		return
	}
	appG.Response(http.StatusOK, 0, token)
	return
}

func (c *Auth) DelCache(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	service.DelAllCache()
	appG.Response(http.StatusOK, 0, nil)
	return
}

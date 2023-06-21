package response

import (
	"net/http"
	"websocket-pool/pkg/common/xerr"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
	ReqId string      `json:"req_id"`
}

const (
	AUTH_ERROR                       = -1 // token验证异常
	ERROR                            = 7
	SUCCESS                          = 0
	ERROR_SMS_NUM_PEER10MIN          = 8  //10分钟内，超过发送次数
	ERROR_SMS_NUM_TODAY              = 9  //今日获取验证码已达上限
	ERROR_VERIFY_INVALID             = 10 //验证码失效
	ERROR_SMS_GET_FAILED             = 11 //验证码获取失败
	ERROR_SMS_GET_FREQUEN            = 12 //验证码获取频繁，请稍后重试
	ERROR_USER_REGISTERED            = 13 //该用户已经注册
	ERROR_VERIFY_CODE                = 14 //验证码错误
	ERROR_REQ_PARAM                  = 15 //请求参数错误
	ERROR_TOKEN_INVALID              = 17 //token失效
	ERROR_ABNORMAL_LOGIN             = 18 //您的账户被异地登录
	ERROR_UNAUTH_ACCESS              = 19 //未登录或者非法访问
	ERROR_USER_FORBIDDEN             = 20 //用户被禁止登录
	ERROR_VIDEO_CALL_APPROVAL_FAILED = 21 // 开启摄像头审批失败
	ERROR_VIDEO_CALL_OPEN_THRESHOLD  = 22 // 加入视频人数已达上限

	CookieToken   = "MINDTOKEN"
	CookieUser    = "MINDUSERID"
	CookieTeam    = "MINDTEAMID"
	CookieImToken = "MINDIMTOKEN"
	CookieImUser  = "MINDIMUSERID"

	BiUserToken = "user_mind_name"
	BiSessionId = "sessionid"
	BiHomeToken = "home_mind_token"
)

func SetCookie(c *gin.Context, name, value string) {
	// c.SetCookie(name, value, global.GVA_CONFIG.OAuth.MaxAge, global.GVA_CONFIG.OAuth.Path, global.GVA_CONFIG.OAuth.Domain, global.GVA_CONFIG.OAuth.Secure, global.GVA_CONFIG.OAuth.HttpOnly)
}

func UnSetCookie(c *gin.Context, name string) {
	// c.SetCookie(name, "", -1, global.GVA_CONFIG.OAuth.Path, global.GVA_CONFIG.OAuth.Domain, global.GVA_CONFIG.OAuth.Secure, global.GVA_CONFIG.OAuth.HttpOnly)
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
		c.GetHeader("X-Request-Id"),
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "操作成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}
func FailWithMessageByCode(message string, errCode int, c *gin.Context) {
	Result(errCode, map[string]interface{}{}, message, c)
}

func FailAuthWithMessage(message string, c *gin.Context) {
	Result(AUTH_ERROR, map[string]interface{}{}, message, c)
}

func ErrorWithMessage(err error, c *gin.Context) {
	//错误返回
	errCode := xerr.ServerCommonError
	errMsg := "服务器开小差啦，稍后再来试一试"
	causeErr := errors.Cause(err)                // err类型
	if e, ok := causeErr.(*xerr.CodeError); ok { //自定义错误类型
		//自定义CodeError
		errCode = e.GetErrCode()
		errMsg = e.GetErrMsg()
	}

	Result(int(errCode), map[string]interface{}{}, errMsg, c)
}

func ErrorWithData(err error, data interface{}, c *gin.Context) {
	//错误返回
	errCode := xerr.ServerCommonError
	errMsg := "服务器开小差啦，稍后再来试一试"
	causeErr := errors.Cause(err)                // err类型
	if e, ok := causeErr.(*xerr.CodeError); ok { //自定义错误类型
		//自定义CodeError
		errCode = e.GetErrCode()
		errMsg = e.GetErrMsg()
	}

	Result(int(errCode), data, errMsg, c)
}

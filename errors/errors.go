package errors

import (
	"github.com/zm-dev/gerrors"
)

// Swagger API documents need this structure.
type GlobalError struct {
	Code        int    `json:"code" example:"10001"`
	ServiceName string `json:"service_name" example:"cloud_disk"`
	Message     string `json:"message" example:"error message"`
	InnerErr    error  `json:"inner_err"`
	StatusCode  int    `json:"status_code" example:"500"`
}

// 参数绑定出错
func BindError(err error) error {
	return gerrors.BadRequest(10001, err.Error(), err)
}

func BadRequest(msg string, err ...error) error {
	return gerrors.BadRequest(10002, msg, err...)
}

func InternalServerError(msg string, err ...error) error {
	return gerrors.InternalServerError(10003, msg, err...)
}

func Unauthorized(message ...string) error {
	var msg string
	if len(message) == 0 {
		msg = "请先登录"
	} else {
		msg = message[0]
	}
	return gerrors.Unauthorized(10004, msg, nil)
}

// NotFound generates a 404 error.
func NotFound(message string, err ...error) error {
	return gerrors.NotFound(10005, message, err...)
}

// 记录不存在
func RecordNotFound(message string) error {
	return gerrors.NotFound(10006, message, nil)
}

// 文件已存在
func FileAlreadyExist(message ...string) error {
	var msg string
	if len(message) == 0 {
		msg = "文件已存在"
	} else {
		msg = message[0]
	}
	return gerrors.New(10007, 400, msg, nil)
}

// 没有权限
func Forbidden(msg string, err ...error) error {
	return gerrors.Forbidden(10008, msg, err...)
}

func ErrAccountAlreadyExisted() error {
	return gerrors.BadRequest(10009, "account already existed", nil)
}

func ErrPassword() error {
	return gerrors.BadRequest(10010, "密码错误", nil)
}

func ErrAccountNotFound() error {
	return gerrors.NotFound(10011, "账号不存在", nil)
}

func UserIsBanned(message ...string) error {
	var msg string
	if len(message) == 0 {
		msg = "此用户已禁用"
	} else {
		msg = message[0]
	}
	return gerrors.Forbidden(10012, msg, nil)
}

func UserNotAllowBeBan(message ...string) error {
	var msg string
	if len(message) == 0 {
		msg = "不允许 ban 该用户"
	} else {
		msg = message[0]
	}
	return gerrors.Forbidden(10013, msg, nil)
}

func GroupNotAllowBeDelete(message ...string) error {
	var msg string
	if len(message) == 0 {
		msg = "不允许删除该组"
	} else {
		msg = message[0]
	}
	return gerrors.Forbidden(10013, msg, nil)
}

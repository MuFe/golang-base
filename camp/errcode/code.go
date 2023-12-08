package errcode

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error 错误
type Error struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// SetData 设置data，枚举是全局的，必须返回一个新的对象
func (t *Error) SetData(data interface{}) *Error {
	return &Error{
		Code: t.Code,
		Msg:  t.Msg,
		Data: data,
	}
}

// AppendMsg 追加错误信息，枚举是全局的，必须返回一个新的对象
func (t *Error) AppendMsg(msg string) *Error {
	return &Error{
		Code: t.Code,
		Msg:  fmt.Sprintf("%s，%s", t.Msg, msg),
	}
}

// SetMsg 追加错误信息，枚举是全局的，必须返回一个新的对象
func (t *Error) SetMsg(msg string) *Error {
	return &Error{
		Code: t.Code,
		Msg: msg,
	}
}

// RPCError 转rpc错误类型
func (t *Error) RPCError() error {
	return status.Error(codes.Code(t.Code), t.Msg)
}

// ParseError 组装错误
func ParseError(err error) Error {
	derr := Error{Code: http.StatusInternalServerError}
	r, ok := status.FromError(err)
	derr.Msg = r.Message()
	codeStr := int(r.Code())
	if ok && r.Code() != codes.Unknown {
		derr.Code = codeStr
	}
	return derr
}

// ParseMsgf 组装错误
func ParseMsgf(f string, params ...interface{}) *Error {
	return ParseMsg(fmt.Sprintf(f, params...))
}


// ParseMsg 组装错误
func ParseMsg(s string) *Error {
	return &Error{Code: http.StatusInternalServerError, Msg: s}
}

// ParseOK 组装成功返回
func ParseOK(data interface{}) *Error {
	return &Error{Code: http.StatusOK, Msg: "OK", Data: data}
}

var commError = 200 // 公共错误开头


// http错误
var (
	HttpErrorOK         = Error{Code: http.StatusOK, Msg: "OK"}
	HttpErrorNotFound   = Error{Code: http.StatusNotFound, Msg: "Not Found"}
	HttpErrorWringParam = Error{Code: http.StatusBadRequest, Msg: "参数错误"}
)

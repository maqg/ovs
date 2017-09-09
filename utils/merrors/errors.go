package merrors

const (
	// ErrSuccess cmd successfully
	ErrSuccess = iota

	// ErrDbErr error for database
	ErrDbErr

	//ErrNotEnoughParas error for not enough paras
	ErrNotEnoughParas

	// ErrTooManyParas error for too many paras
	ErrTooManyParas

	// ErrBadParas error for bad paras
	ErrBadParas

	// ErrCmdErr error for cmd error
	ErrCmdErr

	// ErrCommonErr error for common error
	ErrCommonErr

	// ErrSegmentNotExist error for segment not exist
	ErrSegmentNotExist

	// ErrSegmentAlreadyExist error for segment not exist
	ErrSegmentAlreadyExist

	// ErrTimeout error for timeout error
	ErrTimeout

	// ErrSyscallErr error for system calling error
	ErrSyscallErr

	// ErrSystemErr error for system error
	ErrSystemErr

	// ErrNoSuchAPI error for no such api
	ErrNoSuchAPI

	// ErrNotImplemented error for not implemented error
	ErrNotImplemented

	// User

	// ErrUserNotExist error for user not exist
	ErrUserNotExist

	// ErrUserAlreadyExist error for user already exist
	ErrUserAlreadyExist

	// ErrPasswordDontMatch error for password not match
	ErrPasswordDontMatch

	// ErrUserNotLogin error for user not login
	ErrUserNotLogin
)

// GErrors for global errors mapping
var GErrors = map[int]string{
	ErrSuccess:             "Command Success",
	ErrDbErr:               "Database Error",
	ErrNotEnoughParas:      "No Enough Paras",
	ErrTooManyParas:        "Too Many Paras",
	ErrBadParas:            "Unaccept Paras",
	ErrCmdErr:              "Command Error",
	ErrCommonErr:           "Common Error",
	ErrSegmentNotExist:     "Segment Not Exist",
	ErrSegmentAlreadyExist: "Segment Already Exist",
	ErrTimeout:             "Timeout Error",
	ErrSyscallErr:          "System Call Error",
	ErrSystemErr:           "System Error",
	ErrNoSuchAPI:           "No Such API",
	ErrNotImplemented:      "Function not Implemented",

	// User
	ErrUserNotExist:      "User Not Exist",
	ErrUserAlreadyExist:  "User Already Exist",
	ErrPasswordDontMatch: "User And Password Not Match",
	ErrUserNotLogin:      "User Not Login",
}

// GErrorsCN Global error for Chinese
var GErrorsCN = map[int]string{
	ErrSuccess:             "操作成功",
	ErrDbErr:               "数据库错误",
	ErrNotEnoughParas:      "参数不足",
	ErrTooManyParas:        "太多参数",
	ErrBadParas:            "参数不合法",
	ErrCmdErr:              "命令执行错误",
	ErrCommonErr:           "通用错误",
	ErrSegmentNotExist:     "对象不存在",
	ErrSegmentAlreadyExist: "对象已存在",
	ErrTimeout:             "超时错误",
	ErrSyscallErr:          "系统调用错误",
	ErrSystemErr:           "系统错误",
	ErrNoSuchAPI:           "无此API",
	ErrNotImplemented:      "功能未实现",

	// User
	ErrUserNotExist:      "用户不存在",
	ErrUserAlreadyExist:  "用户已经存在",
	ErrPasswordDontMatch: "用户和密码不匹配",
	ErrUserNotLogin:      "用户未登录",
}

// MError base error structure
type MError struct {
	ErrorNo  int    `json:"no"`
	ErrorMsg string `json:"msg"`
}

// NewError to new an error
func NewError(code int, message string) *MError {
	return &MError{
		ErrorNo:  code,
		ErrorMsg: message,
	}
}

// GetMsg from errorNo
func GetMsg(errorNo int) string {
	return GErrors[errorNo]
}

// GetMsgCN from errorNo
func GetMsgCN(errorNo int) string {
	return GErrorsCN[errorNo]
}

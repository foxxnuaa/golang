package lib

import (
	"time"
)

//载荷发生器的状态
type GenStatus int

const (
	STATUS_ORIGINAL GenStatus = 0
	STATUS_STARTED  GenStatus = 1
	STATUS_STOPPED  GenStatus = 2
)

type ResultCode int

// 保留 1 ~ 1000 给载荷承受者使用。
const (
	RESULT_CODE_SUCCESS                         = 0    // 成功。
	RESULT_CODE_WARNING_CALL_TIMEOUT ResultCode = 1001 // 调用超时警告。
	RESULT_CODE_ERROR_CALL           ResultCode = 2001 // 调用错误。
	RESULT_CODE_ERROR_RESPONSE       ResultCode = 2002 // 响应内容错误。
	RESULT_CODE_ERROR_CALEE          ResultCode = 2003 // 被调用方（被测软件）的内部错误。
	RESULT_CODE_FATAL_CALL           ResultCode = 3001 // 调用过程中发生了致命错误！
)

func GetResultCodePlain(code ResultCode) string {
	var codePlain string
	switch code {
	case RESULT_CODE_SUCCESS:
		codePlain = "Success"
	case RESULT_CODE_WARNING_CALL_TIMEOUT:
		codePlain = "Call Timeout Warning"
	case RESULT_CODE_ERROR_CALL:
		codePlain = "Call Error"
	case RESULT_CODE_ERROR_RESPONSE:
		codePlain = "Response Error"
	case RESULT_CODE_ERROR_CALEE:
		codePlain = "Callee Error"
	case RESULT_CODE_FATAL_CALL:
		codePlain = "Call Fatal Error"
	default:
		codePlain = "Unknown result code"
	}
	return codePlain
}

//载荷发生器的接口
type Generator interface {
	//启动载荷发生器
	Start()
	//停止载荷发生器
	Stop() (uint64, bool)
	//获取状态
	Status() GenStatus
}

// 调用结果的结构
type CallResult struct {
	Id     int64         //ID
	Req    RawReq        //原生请求
	Resp   RawResp       //原生响应
	Code   ResultCode    //响应代码
	Msg    string        //结果成因的简述
	Elapse time.Duration //耗时
}

//原生请求
type RawReq struct {
	Id  int64
	Req []byte
}

//原生响应
type RawResp struct {
	Id     int64
	Resp   []byte
	Err    error
	Elapse time.Duration
}

//调用器
type Caller interface {
	//构建请求
	BuildReq() RawReq
	//调用
	Call(req []byte, timeoutNs time.Duration) ([]byte, error)
	//检查响应
	CheckResp(rawReq RawReq, rawResp RawResp) *CallResult
}

/*goroutine池接口*/
type Pooler interface {
	//获取
	Take()
	//归还
	Return()
	//是否被激活
	Active() bool
	//池内总数
	Total() uint32
	//剩余
	Remainder() uint32
}

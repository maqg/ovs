package vyos

import (
	"net/http"
	"octlink/ovs/utils"
)

type commandHandlerWrap struct {
	path    string
	handler http.HandlerFunc
	async   bool
}

type Options struct {
	Ip           string
	Port         uint
	ReadTimeout  uint
	WriteTimeout uint
	LogFile      string
}

type CommandResponseHeader struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type CommandContext struct {
	responseWriter http.ResponseWriter
	request        *http.Request
}

func (ctx *CommandContext) GetCommand(cmd interface{}) {
	if err := utils.JsonDecodeHttpRequest(ctx.request, cmd); err != nil {
		panic(err)
	}
}

type CommandHandler func(ctx *CommandContext) interface{}

type HttpInterceptor func(http.HandlerFunc) http.HandlerFunc

var (
	commandHandlers     map[string]*commandHandlerWrap = make(map[string]*commandHandlerWrap)
	commandOptions      Options
	CALLBACK_IP         = ""
	CURRENT_CALLBACK_IP = ""
)

const (
	CALLBACK_URL = "callbackurl"
	TASK_UUID    = "taskuuid"
)

func SetOptions(o Options) {
	commandOptions = o
}

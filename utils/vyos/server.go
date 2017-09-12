package vyos

import (
	"net/http"
)

// CommandContext for command descriptor
type CommandContext struct {
	responseWriter http.ResponseWriter
	request        *http.Request
}

// CommandHandler for command manager
type CommandHandler func(ctx *CommandContext) interface{}

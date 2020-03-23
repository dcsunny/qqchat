package safe

import (
	"github.com/dcsunny/qqchat/context"
)

type WxSafe struct {
	*context.Context
}

func NewWxSafe(context *context.Context) *WxSafe {
	tpl := new(WxSafe)
	tpl.Context = context
	return tpl
}

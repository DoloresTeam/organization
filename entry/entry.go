package entry

import (
	"github.com/rs/xid"
)

var subffix string

func SetSubffix(s string) {
	subffix = s
}

func generatorID() string {
	return xid.New().String()
}

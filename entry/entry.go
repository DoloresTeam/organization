package entry

import (
	"github.com/rs/xid"
)

var subfix string

func SetSubfix(s string) {
	subfix = s
}

func generatorID() string {
	return xid.New().String()
}

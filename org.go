package organization

import (
	"errors"

	"github.com/DoloresTeam/organization/entry"
	ldap "gopkg.in/ldap.v2"
)

type Organization struct {
	l      *ldap.Conn
	subfix string
}

func NewOrganization(subfix string, ldapBindConn *ldap.Conn) (*Organization, error) {

	if len(subfix) == 0 || ldapBindConn == nil {
		return nil, errors.New(`subfix and ldapBindConn must not be nil`)
	}

	// 方便构造各种 Request
	entry.SetSubfix(subfix)

	// 验证ldap 的目录结构
	return &Organization{l: ldapBindConn, subfix: subfix}, nil
}

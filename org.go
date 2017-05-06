package organization

import (
	"errors"
	"fmt"

	"github.com/DoloresTeam/organization/gorbacx"

	ldap "gopkg.in/ldap.v2"
)

var UnimplementError = errors.New(`THIS FEATURE UNIMPLEMENT !!`)

type Organization struct {
	l       *ldap.Conn
	rbacx   *gorbacx.RBACX
	subffix string
}

func NewOrganization(subffix string, ldapBindConn *ldap.Conn) (*Organization, error) {

	if len(subffix) == 0 || ldapBindConn == nil {
		return nil, errors.New(`subfix and ldapBindConn must not be nil`)
	}

	// TODO 验证ldap 的目录结构
	org := &Organization{ldapBindConn, gorbacx.New(), subffix}
	org.initial()

	return org, nil
}

func NewOrganizationWithSimpleBind(subffix, host, rootDN, rootPWD string, port int) (*Organization, error) {

	l, err := ldap.Dial(`tcp`, fmt.Sprintf(`%s:%d`, host, port))
	if err != nil {
		return nil, errors.New(`dial ldap server failed`)
	}

	err = l.Bind(rootDN, rootPWD)
	if err != nil {
		return nil, err
	}

	return NewOrganization(subffix, l)
}

func (org *Organization) initial() {

}

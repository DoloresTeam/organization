package organization

import (
	"errors"
	"fmt"

	"github.com/DoloresTeam/organization/gorbacx"

	ldap "gopkg.in/ldap.v2"
)

// Organization ldap operation handler
type Organization struct {
	l       *ldap.Conn
	rbacx   *gorbacx.RBACX
	subffix string
}

// NewOrganization ...
func NewOrganization(subffix string, ldapBindConn *ldap.Conn) (*Organization, error) {

	if len(subffix) == 0 || ldapBindConn == nil {
		return nil, errors.New(`subfix and ldapBindConn must not be nil`)
	}

	// TODO 验证ldap 的目录结构
	org := &Organization{ldapBindConn, gorbacx.New(), subffix}
	err := org.initial()
	if err != nil {
		return nil, err
	}

	return org, nil
}

// NewOrganizationWithSimpleBind ...
func NewOrganizationWithSimpleBind(subffix, host, rootDN, rootPWD string, port int) (*Organization, error) {

	l, err := ldap.Dial(`tcp`, fmt.Sprintf(`%s:%d`, host, port))
	if err != nil {
		return nil, errors.New(`dial ldap server failed`)
	}
	l.Debug = true
	err = l.Bind(rootDN, rootPWD)
	if err != nil {
		return nil, err
	}

	return NewOrganization(subffix, l)
}

func (org *Organization) initial() error {

	rs, err := org.AllRoles()
	if err != nil {
		return err
	}

	var roles []*gorbacx.Role
	for _, v := range rs {

		urs, err := org.PermissionByIDs(v[`upid`].([]string))
		if err != nil {
			return err
		}

		mrs, err := org.PermissionByIDs(v[`ppid`].([]string))
		if err != nil {
			return err
		}

		var ups []*gorbacx.Permission
		for _, info := range urs.Data {
			ups = append(ups, gorbacx.NewPermission(info[`id`].(string), info[`rbacType`].([]string)))
		}

		var mps []*gorbacx.Permission
		for _, info := range mrs.Data {
			mps = append(mps, gorbacx.NewPermission(info[`id`].(string), info[`rbacType`].([]string)))
		}

		roles = append(roles, gorbacx.NewRole(v[`id`].(string), ups, mps))
	}

	org.rbacx.Add(roles)

	// org.rbacx.PrettyPrint()

	return nil
}

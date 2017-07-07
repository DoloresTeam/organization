package organization

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/casbin/casbin"
	"github.com/doloresteam/organization/ldap-pool"
	ldap "gopkg.in/ldap.v2"
)

// Organization ldap operation handler
type Organization struct {
	pool                  ldappool.Pool
	enforcer              *casbin.Enforcer
	subffix               string
	latestResetVersion    string
	organizationViewEvent chan []string
}

// NewOrganization ...
// func NewOrganization(subffix string, ldapBindConn *ldap.Conn, orgViewChangeEvent chan []string) (*Organization, error) {
//
// 	if len(subffix) == 0 || ldapBindConn == nil {
// 		return nil, errors.New(`subfix and ldapBindConn must not be nil`)
// 	}
//
// 	pool := ldappool.NewChannelPool(1, 5, `ldap-default-pool`, func(string) {
//
// 	}, []uint8{ldap.LDAPResultTimeLimitExceeded, ldap.ErrorNetwork})
//
// 	// TODO 验证ldap 的目录结构
// 	org := &Organization{ldapBindConn, gorbacx.New(), subffix, ``, orgViewChangeEvent}
//
// 	return org, org.RefreshRBAC()
// }

// NewOrganizationWithSimpleBind ...
func NewOrganizationWithSimpleBind(subffix, host, rootDN, rootPWD string, port int, orgViewChangeEvent chan []string) (*Organization, error) {

	pool, err := ldappool.NewChannelPool(1, 5, `ldap-default-pool`, func(string) (ldap.Client, error) {
		c, err := ldap.Dial(`tcp`, fmt.Sprintf(`%s:%d`, host, port))
		if err != nil {
			return nil, errors.New(`dial ldap server failed`)
		}
		return c, c.Bind(rootDN, rootPWD)
	}, []uint8{ldap.LDAPResultTimeLimitExceeded, ldap.ErrorNetwork})
	if err != nil {
		return nil, err
	}

	// TODO 验证ldap 的目录结构
	_org := &Organization{pool, &casbin.Enforcer{}, subffix, ``, orgViewChangeEvent}

	m := casbin.NewModel()

	m.AddDef(`r`, `r`, `sub, obj, act`)
	m.AddDef(`p`, `p`, `sub, obj, act`)
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "r.obj == p.obj && r.act == p.act")

	_org.enforcer.InitWithModelAndAdapter(m, _org)

	return _org, nil
}

// OrganizationView get this member's visible departments 、members and version of this `organization view`
func (org *Organization) OrganizationView(id string) ([]map[string]interface{}, []map[string]interface{}, string, error) { // departments, members, version, error
	// 通过id 拿到所有的 角色
	roleIDs, err := org.RoleIDsByMemberID(id)
	if err != nil {
		return nil, nil, ``, err
	}
	// 类型
	// types := org.rbacx.MatchedTypes(roleIDs) // 这个Type包含了当前角色下所有的部门和员工Type， 所有的Type ID 都是全局唯一的
	types := org.fetchAllowedTypesInRoles(roleIDs)

	log.Printf(`%s allowed types %v`, id, types)

	filter, err := sqConvertArraysToFilter(`rbacType`, types)
	if err != nil {
		return nil, nil, ``, err
	}

	unitIDs, err := org.UnitIDsByTypeIDs(types)
	if err != nil {
		return nil, nil, ``, err
	}
	uFilter, err := sqConvertArraysToFilter(`unitID`, unitIDs)
	if err != nil {
		return nil, nil, ``, err
	}

	f := fmt.Sprintf(`(&(%s)(%s))`, filter, uFilter) // 添加部门过滤条件 确定有访问此部门的权限啊
	msq := &searchRequest{org.parentDN(member), f, MemberSignleAttrs[:], MemberMultipleAttrs[:], 0, nil}
	msr, err := org.search(msq)
	if err != nil {
		return nil, nil, ``, err
	}

	departments, err := org.searchUnit(filter, false)
	if err != nil {
		return nil, nil, ``, err
	}

	return departments, msr.Data, newTimeStampVersion(), err
}

func newTimeStampVersion() string {
	return time.Now().UTC().Format(`20060102150405Z`)
}

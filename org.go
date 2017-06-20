package organization

import (
	"errors"
	"fmt"
	"time"

	"github.com/DoloresTeam/organization/gorbacx"
	"github.com/doloresteam/organization/ldap-pool"
	ldap "gopkg.in/ldap.v2"
)

// Organization ldap operation handler
type Organization struct {
	pool                  ldappool.Pool
	rbacx                 *gorbacx.RBACX
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
	_org := &Organization{pool, gorbacx.New(), subffix, ``, orgViewChangeEvent}

	return _org, _org.RefreshRBAC()
}

// RefreshRBAC sync type permission and role with ldap server.
func (org *Organization) RefreshRBAC() error {

	err := func() error {
		org.rbacx.Clear()

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

		return nil
	}()
	return err
}

// OrganizationView get this member's visible departments 、members and version of this `organization view`
func (org *Organization) OrganizationView(id string) ([]map[string]interface{}, []map[string]interface{}, string, error) { // departments, members, version, error
	// 通过id 拿到所有的 角色
	roleIDs, err := org.RoleIDsByMemberID(id)
	if err != nil {
		return nil, nil, ``, err
	}
	// 类型
	types := org.rbacx.MatchedTypes(roleIDs) // 这个Type包含了当前角色下所有的部门和员工Type， 所有的Type ID 都是全局唯一的

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

	f := fmt.Sprintf(`(&(%s)(%s))`, filter, uFilter) // 添加部门过滤条件 确定有访问此部门的全新啊
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

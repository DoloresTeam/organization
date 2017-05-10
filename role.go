package organization

import (
	"fmt"

	"github.com/DoloresTeam/organization/gorbacx"
	ldap "gopkg.in/ldap.v2"
)

// AddRole to ldap server, this method will automatically update org's rbacx
func (org *Organization) AddRole(name, description string, ups, pps []string) error {

	upObjects, err := org.convertIDToObject(ups, true)
	if err != nil {
		return err
	}
	ppObjects, err := org.convertIDToObject(pps, false)
	if err != nil {
		return err
	}

	id := generatorID()
	dn := org.dn(id, role)
	aq := ldap.NewAddRequest(dn)

	aq.Attribute(`objectClass`, []string{`role`, `top`})

	aq.Attribute(`cn`, []string{name})
	aq.Attribute(`description`, []string{description})
	aq.Attribute(`upid`, ups)
	aq.Attribute(`ppid`, pps)

	err = org.l.Add(aq) // 先写数据库
	if err != nil {
		return err
	}

	role := gorbacx.NewRole(id, upObjects, ppObjects)
	org.rbacx.Add([]*gorbacx.Role{role})

	return nil
}

// RemoveRole from ldap server, automatically update org's rbacx
func (org *Organization) RemoveRole(id string) error {

	dn := org.dn(id, role)
	dq := ldap.NewDelRequest(dn, nil)

	err := org.l.Del(dq)
	if err != nil {
		return err
	}

	org.rbacx.Remove([]string{id})

	return nil
}

// ModifyRole in ldap server, automatically update org's rbacx
func (org *Organization) ModifyRole(id, name, description string, ups, pps []string) error {

	upObjects, err := org.convertIDToObject(ups, true)
	if err != nil {
		return err
	}
	ppObjects, err := org.convertIDToObject(pps, true)
	if err != nil {
		return err
	}

	dn := org.dn(id, role)
	mq := ldap.NewModifyRequest(dn)

	if len(name) > 0 {
		mq.Replace(`name`, []string{name})
	}
	if len(description) > 0 {
		mq.Replace(`description`, []string{description})
	}
	if len(ups) > 0 {
		mq.Replace(`upid`, ups)
	}
	if len(pps) > 0 {
		mq.Replace(`ppid`, pps)
	}

	err = org.l.Modify(mq)
	if err != nil {
		return err
	}

	role, err := org.rbacx.RoleByID(id)
	if err != nil {
		return err
	}

	if len(upObjects) > 0 {
		role.Replace(upObjects, true)
	}
	if len(ppObjects) > 0 {
		role.Replace(ppObjects, false)
	}

	return nil
}

// AllRoles in ldap
func (org *Organization) AllRoles() ([]map[string]interface{}, error) {
	return org.search(org.roleSC(``))
}

// RoleByIDs in ldap
func (org *Organization) RoleByIDs(ids []string) ([]map[string]interface{}, error) {

	filter, err := scConvertIDsToFilter(ids)
	if err != nil {
		return nil, err
	}

	return org.search(org.roleSC(filter))
}

// RoleByPermission which role contain this permission
func (org *Organization) RoleByPermission(id string, isUnit bool) ([]string, error) {

	var filter string
	if isUnit {
		filter = fmt.Sprintf(`(upid=%s)`, id)
	} else {
		filter = fmt.Sprintf(`(ppid=%s)`, id)
	}

	sc := &SearchCondition{
		DN:         org.parentDN(role),
		Filter:     filter,
		Attributes: []string{`id`},
		Convertor: func(sr *ldap.SearchResult) []map[string]interface{} {
			var roles []map[string]interface{}
			for _, e := range sr.Entries {
				roles = append(roles, map[string]interface{}{
					`id`: e.GetAttributeValue(`id`),
				})
			}
			return roles
		},
	}

	roles, err := org.search(sc)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, v := range roles {
		ids = append(ids, v[`id`].(string))
	}
	return ids, nil
}

func (org *Organization) convertIDToObject(ids []string, isUnit bool) ([]*gorbacx.Permission, error) {

	var objects []*gorbacx.Permission // 权限有效性判断

	for _, id := range ids {
		p, _ := org.rbacx.PermissionByID(id, isUnit)
		if p != nil {
			objects = append(objects, p)
		} else {
			infos, _ := org.PermissionByIDs([]string{id}, isUnit)
			if len(infos) != 1 {
				return nil, fmt.Errorf(`convert failed no this permission info id: %s`, id)
			}
			objects = append(objects, permissionWithLDAP(infos[0]))
		}
	}

	return objects, nil
}

func permissionWithLDAP(info map[string]interface{}) *gorbacx.Permission {
	return gorbacx.NewPermission(info[`id`].(string), info[`types`].([]string))
}

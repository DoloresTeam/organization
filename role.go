package organization

import (
	"fmt"

	"github.com/DoloresTeam/organization/gorbacx"
	ldap "gopkg.in/ldap.v2"
)

func (org *Organization) AddRole(name, description string, ups, pps []string) error {

	upObjects, err := org.convertIDToObject(ups, true)
	if err != nil {
		return err
	}
	ppObjects, err := org.convertIDToObject(pps, true)
	if err != nil {
		return err
	}

	oid := generatorOID()
	dn := org.dn(oid, ROLE)
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

	role := gorbacx.NewRole(oid, upObjects, ppObjects)
	org.rbacx.Add([]*gorbacx.Role{role})

	return nil
}

func (org *Organization) RemoveRole(oid string) error {

	dn := org.dn(oid, ROLE)
	dq := ldap.NewDelRequest(dn, nil)

	err := org.l.Del(dq)
	if err != nil {
		return err
	}

	org.rbacx.Remove([]string{oid})

	return nil
}

func (org *Organization) ModifyRole(oid, name, description string, ups, pps []string) error {

	upObjects, err := org.convertIDToObject(ups, true)
	if err != nil {
		return err
	}
	ppObjects, err := org.convertIDToObject(pps, true)
	if err != nil {
		return err
	}

	dn := org.dn(oid, ROLE)
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

	role, err := org.rbacx.RoleByID(oid)
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

func (org *Organization) AllRoles() ([]map[string]interface{}, error) {
	return org.filterRole(`(objectClass=role)`)
}

func (org *Organization) RoleByIDs(ids []string) ([]map[string]interface{}, error) {

	filter := `(|(objectClass=role)`

	for _, id := range ids {
		filter = filter + fmt.Sprintf(`(oid=%s)`, id)
	}

	filter = filter + `)`

	return org.filterRole(filter)
}

func (org *Organization) filterRole(filter string) ([]map[string]interface{}, error) {

	sq := ldap.NewSearchRequest(org.parentDN(ROLE),
		ldap.ScopeSingleLevel,
		ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{`oid`, `cn`, `description`, `upid`, `ppid`}, nil)

	sr, err := org.l.Search(sq)
	if err != nil {
		return nil, err
	}

	return convertRoleSearchResult(sr), nil
}

func convertRoleSearchResult(sr *ldap.SearchResult) []map[string]interface{} {

	var roles []map[string]interface{}
	for _, e := range sr.Entries {
		roles = append(roles, map[string]interface{}{
			`id`:                  e.GetAttributeValue(`objectIdentifier`),
			`name`:                e.GetAttributeValue(`cn`),
			`description`:         e.GetAttributeValue(`description`),
			`unitPermissionIDs`:   e.GetAttributeValues(`unitpermissionIdentifier`),
			`personPermissionIDs`: e.GetAttributeValues(`personpermissionIdentifier`),
		})
	}

	return roles
}

func (org *Organization) RoleByPermission(oid string, isUnit bool) ([]string, error) {
	var filter string
	if isUnit {
		filter = fmt.Sprintf(`(upid=%s)`, oid)
	} else {
		filter = fmt.Sprintf(`(ppid=%s)`, oid)
	}
	sq := ldap.NewSearchRequest(org.parentDN(ROLE),
		ldap.ScopeSingleLevel,
		ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{`oid`}, nil)

	sr, err := org.l.Search(sq)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, entry := range sr.Entries {
		ids = append(ids, entry.GetAttributeValue(`oid`))
	}

	return ids, nil
}

func (org *Organization) convertIDToObject(ids []string, isUnit bool) ([]*gorbacx.Permission, error) {

	var objects []*gorbacx.Permission // 权限有效性判断
	for _, id := range ids {
		p, _ := org.rbacx.PermissionByID(id, true)
		if p != nil {
			objects = append(objects, p)
		} else {
			p, err := org.PermissionByID(id, true) // 从服务器获取
			if p == nil {
				return nil, err
			}
			objects = append(objects, p)
		}
	}

	return objects, nil
}

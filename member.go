package organization

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

// MemberSignleAttrs default signle attributes
var MemberSignleAttrs = [...]string{`id`, `name`, `cn`, `telephoneNumber`, `labeledURI`, `gender`, `title`, `priority`}

// MemberSignleACLAttrs default signle acl attributes
var MemberSignleACLAttrs = [...]string{`rbacType`}

// MemberMultipleAttrs default multiple attributes
var MemberMultipleAttrs = [...]string{`email`, `unitID`}

// MemberMultipleACLAttrs default multiple acl attributes
var MemberMultipleACLAttrs = [...]string{`rbacRole`}

// AddMember to ldap server
func (org *Organization) AddMember(info map[string][]string) (string, error) {

	id := generateNewID()
	aq := ldap.NewAddRequest(org.dn(id, member))

	aq.Attribute(`objectClass`, []string{`member`, `top`})
	aq.Attribute(`id`, []string{id})

	for k, v := range info {
		aq.Attribute(k, v)
	}

	err := org.Add(aq)
	if err != nil {
		return ``, err
	}
	nm, err := org.MemberByID(id, true, false)
	if err != nil {
		return ``, err
	}

	return id, org.logAddMember(nm)
}

// ModifyMember ...
func (org *Organization) ModifyMember(id string, info map[string][]string) error {

	oMember, err := org.MemberByID(id, true, false)
	if err != nil {
		return err
	}
	dn := org.dn(id, member)
	mq := ldap.NewModifyRequest(dn)

	for k, v := range info {
		mq.Replace(k, v)
	}

	err = org.Modify(mq)
	if err == nil {
		go func() {
			nMember, _ := org.MemberByID(id, true, false)
			org.logModifyMember(oMember, nMember)
		}()
	}

	return err
}

// DelMember by id
func (org *Organization) DelMember(id string) error {
	m, err := org.MemberByID(id, true, false)
	if err != nil {
		return err
	}
	mids, err := org.relatedMIDs(m[`rbacType`].(string))
	if err != nil {
		return err
	}
	dq := ldap.NewDelRequest(fmt.Sprintf(`id=%s,%s`, id, org.parentDN(member)), nil)
	err = org.Del(dq)
	if err != nil {
		return err
	}
	return org.logDelMember(id, mids)
}

// AuthMember ...
func (org *Organization) AuthMember(telephoneNumber, pwd string) (string, error) {
	filter := fmt.Sprintf(`(telephoneNumber=%s)`, telephoneNumber)
	sq := ldap.NewSearchRequest(org.parentDN(member), ldap.ScopeSingleLevel, ldap.DerefAlways, 0, 0, false, filter, []string{`id`}, nil)
	r, err := org.Search(sq)
	if err != nil {
		return ``, err
	}
	if len(r.Entries) != 1 {
		return ``, fmt.Errorf(`can't find this member by tel: [%s]`, telephoneNumber)
	}

	success, err := org.Compare(r.Entries[0].DN, `userPassword`, pwd)
	if err != nil {
		return ``, err
	}
	if !success {
		return ``, errors.New(`password incorrect`)
	}
	return r.Entries[0].GetAttributeValue(`id`), nil
}

// ModifyPassword ...
func (org *Organization) ModifyPassword(id, originalPassword, newPassword string) error {
	dn := org.dn(id, member)
	success, err := org.Compare(dn, `userPassword`, originalPassword)
	if err != nil {
		return err
	}
	if !success {
		return errors.New(`password incorrect`)
	}

	mq := ldap.NewModifyRequest(dn)
	mq.Replace(`userPassword`, []string{newPassword})
	return org.Modify(mq) // 这样写可以避免产生审计日志
}

// Members return all members
func (org *Organization) Members(pageSize uint32, cookie []byte) (*SearchResult, error) {

	sq := &searchRequest{
		org.parentDN(member), `(objectClass=member)`,
		append(MemberSignleAttrs[:], MemberSignleACLAttrs[:]...),
		append(MemberMultipleAttrs[:], MemberMultipleACLAttrs[:]...), pageSize, cookie}

	return org.search(sq)
}

// MemberIDsByTypeIDs ...
func (org *Organization) MemberIDsByTypeIDs(tids []string) ([]string, error) {
	filter, err := sqConvertArraysToFilter(`rbacType`, tids)
	if err != nil {
		return nil, err
	}
	return org.memberIDsByFilter(filter)
}

// MemberIDsByRoleIDs ...
func (org *Organization) MemberIDsByRoleIDs(rids []string) ([]string, error) {
	filter, err := sqConvertArraysToFilter(`rbacRole`, rids)
	if err != nil {
		return nil, err
	}
	return org.memberIDsByFilter(filter)
}

// MemberIDsByDepartmentIDs ...
func (org *Organization) MemberIDsByDepartmentIDs(ids []string) ([]string, error) {
	filter, err := sqConvertArraysToFilter(`unitID`, ids)
	if err != nil {
		return nil, err
	}
	return org.memberIDsByFilter(filter)
}

func (org *Organization) memberIDsByFilter(filter string) ([]string, error) {
	dn := org.parentDN(member)
	sq := ldap.NewSearchRequest(dn, ldap.ScopeSingleLevel, ldap.DerefAlways,
		0, 0, false, filter, []string{`id`}, nil)
	sr, e := org.Search(sq)
	if e != nil {
		return nil, e
	}
	var ids []string
	for _, entry := range sr.Entries {
		ids = append(ids, entry.GetAttributeValue(`id`))
	}
	return ids, nil
}

// MemberByIDs ...
func (org *Organization) MemberByIDs(ids []string, containACL bool, containPwd bool) ([]map[string]interface{}, error) {
	sa := MemberSignleAttrs[:]
	ma := MemberMultipleAttrs[:]
	if containACL {
		sa = append(sa, MemberSignleACLAttrs[:]...)
		ma = append(ma, MemberMultipleACLAttrs[:]...)
	}
	if containPwd {
		sa = append(sa, `thirdPassword`)
	}

	dn := org.parentDN(member)
	filter, err := sqConvertIDsToFilter(ids)
	if err != nil {
		return nil, err
	}
	sq := &searchRequest{dn, filter, sa, ma, 0, nil}

	r, e := org.search(sq)
	if e != nil {
		return nil, e
	}
	return r.Data, nil
}

// MemberByID search member by id
func (org *Organization) MemberByID(id string, containACL bool, containPwd bool) (map[string]interface{}, error) {
	members, err := org.MemberByIDs([]string{id}, containACL, containPwd)
	if err != nil {
		return nil, err
	}
	if len(members) != 1 {
		return nil, errors.New(`not found`)
	}
	return members[0], nil
}

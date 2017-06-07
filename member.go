package organization

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

var memberSignleAttrs = [...]string{`id`, `name`, `cn`, `telephoneNumber`, `labeledURI`, `gender`, `title`}
var memberSignleACLAttrs = [...]string{`rbacType`}

var memberMutAttrs = [...]string{`email`, `unitID`}
var memberMutACLAttrs = [...]string{`rbacRole`}

// AddMember to ldap server
func (org *Organization) AddMember(info map[string][]string) (string, error) {

	id := generatorID()
	aq := ldap.NewAddRequest(org.dn(id, member))

	aq.Attribute(`objectClass`, []string{`member`, `top`})
	aq.Attribute(`id`, []string{id})

	for k, v := range info {
		aq.Attribute(k, v)
	}

	return id, org.l.Add(aq)
}

// ModifyMember ...
func (org *Organization) ModifyMember(id string, info map[string][]string) error {

	dn := org.dn(id, member)

	mq := ldap.NewModifyRequest(dn)

	for k, v := range info {
		mq.Replace(k, v)
	}

	return org.l.Modify(mq)
}

// DelMember by id
func (org *Organization) DelMember(id string) error {
	dq := ldap.NewDelRequest(fmt.Sprintf(`id=%s,%s`, id, org.parentDN(member)), nil)
	return org.l.Del(dq)
}

// AuthMember ...
func (org *Organization) AuthMember(telephoneNumber, pwd string) (string, error) {

	dn := org.parentDN(member)
	filter := fmt.Sprintf(`(&(telephoneNumber=%s)(userPassword=%s))`, telephoneNumber, fmt.Sprintf(`{MD5}%s`, pwd))
	sq := &searchRequest{dn, filter, []string{`id`}, nil, 0, nil}

	r, e := org.search(sq)
	if e != nil {
		return ``, e
	}
	if len(r.Data) != 1 {
		return ``, errors.New(`not found`)
	}
	return r.Data[0][`id`].(string), nil
}

// Members return all members
func (org *Organization) Members(pageSize uint32, cookie []byte) (*SearchResult, error) {

	sq := &searchRequest{
		org.parentDN(member), `(objectClass=member)`,
		append(memberSignleAttrs[:], memberSignleACLAttrs[:]...),
		append(memberMutAttrs[:], memberMutACLAttrs[:]...), 0, nil}

	return org.search(sq)
}

// MemberIDsByTypeIDs
func (org *Organization) MemberIDsByTypeIDs(tids []string) ([]string, error) {
	filter, err := sqConvertArraysToFilter(`rbacType`, tids)
	if err != nil {
		return nil, err
	}
	return org.memberIDsByFilter(filter)
}

// MemeberIDsByRoleID
func (org *Organization) MemberIDsByRoleIDs(rids []string) ([]string, error) {
	filter, err := sqConvertArraysToFilter(`rbacRole`, rids)
	if err != nil {
		return nil, err
	}
	return org.memberIDsByFilter(filter)
}

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
	sr, e := org.l.Search(sq)
	if e != nil {
		return nil, e
	}
	var ids []string
	for _, entry := range sr.Entries {
		ids = append(ids, entry.GetAttributeValue(`id`))
	}
	return ids, nil
}

// MemberByID search member by id
func (org *Organization) MemberByID(id string, containACL bool, containPwd bool) (map[string]interface{}, error) {
	sa := memberSignleAttrs[:]
	ma := memberMutAttrs[:]
	if containACL {
		sa = append(sa, memberSignleACLAttrs[:]...)
		ma = append(ma, memberMutACLAttrs[:]...)
	}
	if containPwd {
		sa = append(sa, `thirdPassword`)
	}

	dn := org.parentDN(member)
	filter := fmt.Sprintf(`(id=%s)`, id)
	sq := &searchRequest{dn, filter, sa, ma, 0, nil}

	r, e := org.search(sq)
	if e != nil {
		return nil, e
	}
	if len(r.Data) != 1 {
		return nil, errors.New(`not found`)
	}
	return r.Data[0], nil
}

// OrganizationMemberByMemberID ...
func (org *Organization) OrganizationMemberByMemberID(id string) ([]map[string]interface{}, error) {
	dn := org.parentDN(member)
	filter, err := org.filterConditionByMemberID(id, false)
	if err != nil {
		return nil, err
	}

	sq := &searchRequest{dn, filter, memberSignleAttrs[:], memberMutAttrs[:], 0, nil}

	r, e := org.search(sq)
	if e != nil {
		return nil, e
	}

	return r.Data, nil
}

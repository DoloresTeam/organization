package organization

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

var memberSignleAttrs = [...]string{`id`, `name`, `cn`, `telephoneNumber`, `labeledURI`, `gender`}
var memberSignleACLAttrs = [...]string{`thirdAccount`, `thirdPassword`}

var memberMutAttrs = [...]string{`email`, `title`, `unitID`}
var memberMutACLAttrs = [...]string{`rbacType`, `rbacRole`}

// AddMember to ldap server
func (org *Organization) AddMember(info map[string][]string) error {

	aq := ldap.NewAddRequest(org.dn(generatorID(), member))

	aq.Attribute(`objectClass`, []string{`member`, `top`})

	for k, v := range info {
		aq.Attribute(k, v)
	}

	return org.l.Add(aq)
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
	if len(id) == 0 {
		return errors.New(`member id is empty`)
	}

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
		return ``, errors.New(`404 Not Found`)
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
	dn := org.parentDN(member)
	filter, err := sqConvertArraysToFilter(`rbacType`, tids)
	if err != nil {
		return nil, err
	}
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
func (org *Organization) MemberByID(id string, containACL bool) (map[string]interface{}, error) {
	sa := memberSignleAttrs[:]
	ma := memberMutAttrs[:]
	if containACL {
		sa = append(sa, memberSignleACLAttrs[:]...)
		ma = append(ma, memberMutACLAttrs[:]...)
	}

	dn := org.parentDN(member)
	filter := fmt.Sprintf(`(id=%s)`, id)
	sq := &searchRequest{dn, filter, sa, ma, 0, nil}

	r, e := org.search(sq)
	if e != nil {
		return nil, e
	}
	if len(r.Data) != 1 {
		return nil, errors.New(`404 Not Found`)
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

	sq := &searchRequest{dn, filter, append(memberSignleAttrs[:], `thirdAccount`), memberMutAttrs[:], 0, nil}

	r, e := org.search(sq)
	if e != nil {
		return nil, e
	}

	return r.Data, nil
}

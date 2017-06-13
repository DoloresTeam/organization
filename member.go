package organization

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

var MemberSignleAttrs = [...]string{`id`, `name`, `cn`, `telephoneNumber`, `labeledURI`, `gender`, `title`, `priority`}
var MemberSignleACLAttrs = [...]string{`rbacType`}

var MemberMutAttrs = [...]string{`email`, `unitID`}
var MemberMutACLAttrs = [...]string{`rbacRole`}

// AddMember to ldap server
func (org *Organization) AddMember(info map[string][]string) (string, error) {

	id := generateNewID()
	aq := ldap.NewAddRequest(org.dn(id, member))

	aq.Attribute(`objectClass`, []string{`member`, `top`})
	aq.Attribute(`id`, []string{id})

	for k, v := range info {
		aq.Attribute(k, v)
	}

	err := org.l.Add(aq)
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

	err = org.l.Modify(mq)
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
	mids, err := org.MemberIDsByTypeIDs([]string{m[`rbacType`].(string)})
	if err != nil {
		return err
	}
	dq := ldap.NewDelRequest(fmt.Sprintf(`id=%s,%s`, id, org.parentDN(member)), nil)
	err = org.l.Del(dq)
	if err != nil {
		return err
	}
	return org.logDelMember(id, mids)
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
		append(MemberSignleAttrs[:], MemberSignleACLAttrs[:]...),
		append(MemberMutAttrs[:], MemberMutACLAttrs[:]...), pageSize, cookie}

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

func (org *Organization) MemberByIDs(ids []string, containACL bool, containPwd bool) ([]map[string]interface{}, error) {
	sa := MemberSignleAttrs[:]
	ma := MemberMutAttrs[:]
	if containACL {
		sa = append(sa, MemberSignleACLAttrs[:]...)
		ma = append(ma, MemberMutACLAttrs[:]...)
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

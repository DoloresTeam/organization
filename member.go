package organization

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

// AuthMember ...
func (org *Organization) AuthMember(telephoneNumber, pwd string) (string, error) {

	sq := searchRequest(org.parentDN(member), fmt.Sprintf(`(&(telephoneNumber=%s)(userPassword=%s))`, telephoneNumber, pwd), []string{`id`})

	sr, err := org.l.Search(sq)
	if err != nil {
		return ``, err
	}

	if len(sr.Entries) != 1 {
		return ``, errors.New(`Auth failed !`)
	}
	return sr.Entries[0].GetAttributeValue(`id`), nil
}

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

// RoleIDsByMemberID ...
func (org *Organization) RoleIDsByMemberID(id string) ([]string, error) {

	if len(id) == 0 {
		return nil, errors.New(`id must not be empty`)
	}

	filter := fmt.Sprintf(`(id=%s)`, id)

	sq := ldap.NewSearchRequest(org.parentDN(member),
		ldap.ScopeSingleLevel,
		ldap.DerefAlways, 0, 0, false, filter, []string{`rbacRole`}, nil)
	sr, err := org.l.Search(sq)
	if err != nil {
		return nil, err
	}
	if len(sr.Entries) != 1 {
		return nil, errors.New(`can't find this member`)
	}
	return sr.Entries[0].GetAttributeValues(`rbacRole`), nil
}

// MemberByID search member by id
func (org *Organization) MemberByID(id string, containACL bool) (map[string]interface{}, error) {

	sc := org.memberSC(fmt.Sprintf(`(id=%s)`, id), containACL)
	rs, err := org.search(sc)
	if err != nil {
		return nil, err
	}

	if len(rs) != 1 {
		return nil, errors.New(`Found member than 1`)
	}
	return rs[0], nil
}

// OrganizationMemberByMemberID ...
func (org *Organization) OrganizationMemberByMemberID(id string) ([]map[string]interface{}, error) {

	filter, err := org.filterByMemberID(id, false)
	if err != nil {
		return nil, err
	}

	return org.search(org.memberSC(filter, false))
}

package entry

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

const (
	UnitType   = `unit`
	PersonType = `person`
)

func AddType(name, description string, which string) (*ldap.AddRequest, error) {

	aq := ldap.NewAddRequest(dn(generatorID(), which))

	aq.Attribute(`objectClass`, []string{`doloresType`, `top`})

	if len(name) == 0 {
		return nil, errors.New(`name must be not nil`)
	}
	aq.Attribute(`cn`, []string{name})

	if len(description) != 0 {
		aq.Attribute(`description`, []string{description})
	}

	return aq, nil
}

func ModifyType(oid, name, description string, which string) (*ldap.ModifyRequest, error) {

	if len(oid) == 0 {
		return nil, errors.New(`type identifer must not be nil`)
	}

	mq := ldap.NewModifyRequest(dn(oid, which))

	if len(name) != 0 {
		mq.Add(`cn`, []string{name})
	}
	if len(description) != 0 {
		mq.Add(`description`, []string{description})
	}
	if len(mq.AddAttributes) == 0 {
		return nil, errors.New(`non attribute will be modify`)
	}
	return mq, nil
}

func DelType(oid string, which string) *ldap.DelRequest {

	dq := ldap.NewDelRequest(dn(oid, which), nil)

	return dq
}

func FetchAllTypes(which string) *ldap.SearchRequest {

	searchBaseDN := fmt.Sprintf(`ou=%s,%s`, which, baseDN())

	return ldap.NewSearchRequest(searchBaseDN,
		ldap.ScopeSingleLevel,
		ldap.NeverDerefAliases, 0, 0, false,
		`(objectClass=doloresType)`,
		[]string{`oid`, `cn`, `description`}, nil)
}

func dn(oid, which string) string {

	return fmt.Sprintf(`oid=%s,ou=%s,%s`, oid, which, baseDN())
}

func baseDN() string {
	return fmt.Sprintf(`ou=type,%s`, subffix)
}

package organization

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

func (org *Organization) search(sc *SearchCondition) ([]map[string]interface{}, error) {

	sr, err := org.l.Search(sc.sq())
	if err != nil {
		return nil, err
	}

	return sc.Convertor(sr), nil
}

type SearchCondition struct {
	DN         string
	Filter     string
	Attributes []string
	Convertor  func(*ldap.SearchResult) []map[string]interface{}
}

func (sc *SearchCondition) sq() *ldap.SearchRequest {
	return ldap.NewSearchRequest(sc.DN,
		ldap.ScopeSingleLevel,
		ldap.NeverDerefAliases, 0, 0, false,
		sc.Filter, sc.Attributes, nil)
}

func (org *Organization) typeSC(filter string, isUnit bool) *SearchCondition {
	if len(filter) > 0 {
		filter = fmt.Sprintf(`(&(objectClass=doloresType)(%s))`, filter)
	} else {
		filter = `(objectClass=doloresType)`
	}
	return &SearchCondition{
		DN:         org.parentDN(typeCategory(isUnit)),
		Filter:     filter,
		Attributes: []string{`oid`, `cn`, `description`},
		Convertor: func(sr *ldap.SearchResult) []map[string]interface{} {
			var types []map[string]interface{}
			for _, e := range sr.Entries {
				types = append(types, map[string]interface{}{
					`id`:          e.GetAttributeValue(`objectIdentifier`),
					`name`:        e.GetAttributeValue(`cn`),
					`description`: e.GetAttributeValue(`description`),
				})
			}
			return types
		},
	}
}

func (org *Organization) roleSC(filter string) *SearchCondition {
	return &SearchCondition{
		DN:         org.parentDN(ROLE),
		Filter:     fmt.Sprintf(`(&(objectClass=role)%s)`, filter),
		Attributes: []string{`oid`, `cn`, `description`, `upid`, `ppid`},
		Convertor: func(sr *ldap.SearchResult) []map[string]interface{} {

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
		},
	}
}

func (org *Organization) permissionSC(filter string, isUnit bool) *SearchCondition {
	if len(filter) > 0 {
		filter = fmt.Sprintf(`(&(objectClass=permission)(%s))`, filter)
	} else {
		filter = `(objectClass=permission)`
	}
	return &SearchCondition{
		DN:         org.parentDN(permissionCategory(isUnit)),
		Filter:     filter,
		Attributes: []string{`oid`, `cn`, `description`, `rbacType`},
		Convertor: func(sr *ldap.SearchResult) []map[string]interface{} {
			var types []map[string]interface{}
			for _, e := range sr.Entries {
				types = append(types, map[string]interface{}{
					`id`:          e.GetAttributeValue(`objectIdentifier`),
					`name`:        e.GetAttributeValue(`cn`),
					`description`: e.GetAttributeValue(`description`),
					`types`:       e.GetAttributeValues(`rbacType`),
				})
			}
			return types
		},
	}
}

func scConvertIDsToFilter(ids []string) (string, error) {
	if len(ids) == 0 {
		return ``, errors.New(`At least one id`)
	}

	filter := `(|`
	for _, id := range ids {
		filter += fmt.Sprintf(`(oid=%s)`, id)
	}
	filter += `)`

	return filter, nil
}

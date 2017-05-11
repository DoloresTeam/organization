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

// SearchCondition desgin to constrctor search request
type SearchCondition struct {
	DN         string
	Filter     string
	Attributes []string
	Convertor  func(*ldap.SearchResult) []map[string]interface{}
}

func (sc *SearchCondition) sq() *ldap.SearchRequest {
	return ldap.NewSearchRequest(sc.DN,
		ldap.ScopeWholeSubtree,
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
		Attributes: []string{`id`, `cn`, `description`},
		Convertor: func(sr *ldap.SearchResult) []map[string]interface{} {
			var types []map[string]interface{}
			for _, e := range sr.Entries {
				types = append(types, map[string]interface{}{
					`id`:          e.GetAttributeValue(`id`),
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
		DN:         org.parentDN(role),
		Filter:     fmt.Sprintf(`(&(objectClass=role)%s)`, filter),
		Attributes: []string{`id`, `cn`, `description`, `upid`, `ppid`},
		Convertor: func(sr *ldap.SearchResult) []map[string]interface{} {

			var roles []map[string]interface{}
			for _, e := range sr.Entries {
				roles = append(roles, map[string]interface{}{
					`id`:                  e.GetAttributeValue(`id`),
					`name`:                e.GetAttributeValue(`cn`),
					`description`:         e.GetAttributeValue(`description`),
					`unitPermissionIDs`:   e.GetAttributeValues(`unitpermissionIdentifier`),
					`memberPermissionIDs`: e.GetAttributeValues(`memberpermissionIdentifier`),
				})
			}

			return roles
		},
	}
}

func (org *Organization) permissionSC(filter string, isUnit bool) *SearchCondition {
	if len(filter) > 0 {
		filter = fmt.Sprintf(`(&(objectClass=permission)%s)`, filter)
	} else {
		filter = `(objectClass=permission)`
	}
	return &SearchCondition{
		DN:         org.parentDN(permissionCategory(isUnit)),
		Filter:     filter,
		Attributes: []string{`id`, `cn`, `description`, `rbacType`},
		Convertor: func(sr *ldap.SearchResult) []map[string]interface{} {
			var types []map[string]interface{}
			for _, e := range sr.Entries {
				types = append(types, map[string]interface{}{
					`id`:          e.GetAttributeValue(`id`),
					`name`:        e.GetAttributeValue(`cn`),
					`description`: e.GetAttributeValue(`description`),
					`types`:       e.GetAttributeValues(`rbacType`),
				})
			}
			return types
		},
	}
}

func (org *Organization) unitSC(filter string, containACL bool) *SearchCondition {

	attributes := []string{`id`, `ou`, `description`}
	if containACL {
		attributes = append(attributes, `rbacType`)
	}

	if len(filter) > 0 {
		filter = fmt.Sprintf(`(&(objectClass=organizationalUnit)%s)`, filter)
	} else {
		filter = `(objectClass=organizationalUnit)`
	}

	return &SearchCondition{
		DN:         org.parentDN(unit),
		Filter:     filter,
		Attributes: attributes,
		Convertor: func(sr *ldap.SearchResult) []map[string]interface{} {
			var units []map[string]interface{}
			for _, e := range sr.Entries {
				unit := make(map[string]interface{}, 0)
				dn, _ := ldap.ParseDN(e.DN)
				if len(dn.RDNs) > 5 {
					unit[`pid`] = dn.RDNs[1].Attributes[0].Value
				}
				unit[`id`] = e.GetAttributeValue(`id`)
				unit[`name`] = e.GetAttributeValue(`ou`)
				unit[`description`] = e.GetAttributeValue(`description`)
				if containACL {
					unit[`types`] = e.GetAttributeValue(`rbacType`)
				}
				units = append(units, unit)
			}
			return units
		},
	}
}

func scConvertIDsToFilter(ids []string) (string, error) {
	if len(ids) == 0 {
		return ``, errors.New(`At least one id`)
	}

	filter := `(|`
	for _, id := range ids {
		filter += fmt.Sprintf(`(id=%s)`, id)
	}
	filter += `)`

	return filter, nil
}

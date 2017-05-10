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

type mapperItem struct {
	key           string
	value         string
	isSignleValue bool
}

type ldapMapper struct {
	items []*mapperItem
	keys  []string
}

type mapper map[string]interface{}

func newLDAPMapper(m mapper) *ldapMapper {
	var items []*mapperItem
	var keys []string

	delete(m, `id`)

	for key, value := range m {

		if s, ok := value.(string); ok {
			items = append(items, &mapperItem{key, s, true})
			keys = append(keys, key)
			continue
		}
		if ss, ok := value.([]string); ok && len(ss) == 1 {
			items = append(items, &mapperItem{key, ss[0], false})
			keys = append(keys, key)
		}
	}
	return &ldapMapper{items, keys}
}

func (lm *ldapMapper) attributes() []string {
	return append(lm.keys, `id`)
}

func (lm *ldapMapper) mapEntries(sr *ldap.SearchResult) []map[string]interface{} {

	var values []map[string]interface{}

	for _, e := range sr.Entries {
		value := make(map[string]interface{}, len(lm.keys))
		for _, item := range lm.items {
			if item.isSignleValue {
				value[item.value] = e.GetAttributeValue(item.key)
			} else {
				value[item.value] = e.GetAttributeValues(item.key)
			}
		}
		value[`id`] = e.GetAttributeValue(`objectIdentifier`)
		values = append(values, value)
	}

	return values
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

func (org *Organization) unitSC(filter string, m mapper) *SearchCondition {

	lm := newLDAPMapper(m)

	if len(filter) > 0 {
		filter = fmt.Sprintf(`(&(objectClass=organizationalUnit)%s)`, filter)
	} else {
		filter = `(objectClass=organizationalUnit)`
	}

	return &SearchCondition{
		DN:         org.parentDN(unit),
		Filter:     filter,
		Attributes: lm.attributes(),
		Convertor:  lm.mapEntries,
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

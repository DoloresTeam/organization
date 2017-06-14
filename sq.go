package organization

import (
	"fmt"
	"strings"

	ldap "gopkg.in/ldap.v2"
)

type searchRequest struct {
	dn     string
	filter string
	sAttrs []string
	mAttrs []string
	size   uint32
	cookie []byte
}

// SearchResult ...
type SearchResult struct {
	Size   uint32
	Cookie []byte
	Data   []map[string]interface{}
}

func (org *Organization) search(sq *searchRequest) (*SearchResult, error) {

	var controls []ldap.Control

	if sq.size > 0 {
		control := ldap.NewControlPaging(sq.size)
		control.SetCookie(sq.cookie)
		controls = append(controls, control)
	}

	lsq := ldap.NewSearchRequest(sq.dn,
		ldap.ScopeWholeSubtree,
		ldap.DerefAlways, 0, 0, false,
		sq.filter,
		append(sq.sAttrs, sq.mAttrs...),
		controls)

	var lsr *ldap.SearchResult
	var err error

	if sq.size == 0 && sq.cookie == nil {
		lsr, err = org.SearchWithPaging(lsq, 200) // 会循环拿完所有对象
	} else {
		lsr, err = org.Search(lsq)
	}

	if err != nil {
		return nil, err
	}

	data := make([]map[string]interface{}, 0)

	// 所有的人不返回dn属性
	shouldAppendDN := !strings.Contains(sq.dn, `ou=member`)

	for _, e := range lsr.Entries {
		v := make(map[string]interface{}, 0)
		for _, s := range sq.sAttrs {
			v[s] = e.GetAttributeValue(s)
		}
		for _, m := range sq.mAttrs {
			v[m] = e.GetAttributeValues(m)
		}
		if shouldAppendDN {
			v[`dn`] = e.DN
		}
		data = append(data, v)
	}

	var cookie []byte
	var size uint32

	pagingResult := ldap.FindControl(lsr.Controls, ldap.ControlTypePaging)
	if pagingResult != nil {
		cookie = pagingResult.(*ldap.ControlPaging).Cookie
		size = pagingResult.(*ldap.ControlPaging).PagingSize
	}

	return &SearchResult{size, cookie, data}, nil
}

// func (org *Organization) searchType(filter, dn string, pageSize uint32, cookie []byte) (*SearchResult, error) {
// 	if len(filter) == 0 {
// 		filter = `(objectClass=doloresType)`
// 	} else {
// 		filter = fmt.Sprintf(`(&(objectClass=doloresType)%s)`, filter)
// 	}
// 	return org.search(&searchRequest{(dn),
// 		filter,
// 		[]string{`id`, `cn`, `description`, `modifyTimestamp`}, nil,
// 		pageSize,
// 		cookie,
// 	})
// }

func (org *Organization) searchPermission(filter, dn string, pageSize uint32, cookie []byte) (*SearchResult, error) {
	if len(filter) == 0 {
		filter = `(objectClass=permission)`
	} else {
		filter = fmt.Sprintf(`(&(objectClass=permission)%s)`, filter)
	}
	return org.search(&searchRequest{dn,
		filter,
		[]string{`id`, `cn`, `description`},
		[]string{`rbacType`},
		pageSize,
		cookie,
	})
}

func (org *Organization) searchRole(filter string, pageSize uint32, cookie []byte) (*SearchResult, error) {
	if len(filter) == 0 {
		filter = `(objectClass=role)`
	} else {
		filter = fmt.Sprintf(`(&(objectClass=role)%s)`, filter)
	}
	return org.search(&searchRequest{org.parentDN(role),
		filter,
		[]string{`id`, `cn`, `description`},
		[]string{`upid`, `ppid`},
		pageSize,
		cookie,
	})
}

func (org *Organization) searchUnit(filter string, containACL bool) ([]map[string]interface{}, error) {
	attrs := UnitAttributes[0:]
	if containACL {
		attrs = append(attrs, `rbacType`)
	}

	if len(filter) == 0 {
		filter = `(objectClass=organizationalUnit)`
	}
	sq := ldap.NewSearchRequest(org.parentDN(unit),
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		filter, attrs, nil)

	sr, e := org.Search(sq)
	if e != nil {
		return nil, e
	}
	units := make([]map[string]interface{}, 0)
	for _, e := range sr.Entries {
		unit := make(map[string]interface{}, 0)
		for _, attr := range attrs {
			unit[attr] = e.GetAttributeValue(attr)
		}
		// unit[`id`] = e.GetAttributeValue(`id`)
		// unit[`cn`] = e.GetAttributeValue(`ou`)
		// unit[`description`] = e.GetAttributeValue(`description`)
		if containACL {
			unit[`dn`] = e.DN
		}
		dn, _ := ldap.ParseDN(e.DN)
		if len(dn.RDNs) > 5 {
			unit[`parentID`] = dn.RDNs[1].Attributes[0].Value
		}

		units = append(units, unit)
	}

	return units, nil
}

func sqConvertIDsToFilter(ids []string) (string, error) {
	return sqConvertArraysToFilter(`id`, ids)
}

func sqConvertArraysToFilter(label string, datas []string) (string, error) {
	if len(datas) == 0 {
		return ``, fmt.Errorf(`At least one %s`, label)
	}
	if len(datas) == 1 {
		return fmt.Sprintf(`(%s=%s)`, label, datas[0]), nil
	}

	filter := `(|`
	for _, id := range datas {
		filter += fmt.Sprintf(`(%s=%s)`, label, id)
	}
	filter += `)`

	return filter, nil
}

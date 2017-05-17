package organization

import (
	"fmt"

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
		lsr, err = org.l.SearchWithPaging(lsq, 200) // 会循环拿完所有对象
	} else {
		lsr, err = org.l.Search(lsq)
	}

	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	for _, e := range lsr.Entries {
		v := make(map[string]interface{}, 0)
		for _, s := range sq.sAttrs {
			v[s] = e.GetAttributeValue(s)
		}
		for _, m := range sq.mAttrs {
			v[m] = e.GetAttributeValues(m)
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

func (org *Organization) searchType(filter string, isUnit bool, pageSize uint32, cookie []byte) (*SearchResult, error) {
	if len(filter) == 0 {
		filter = `(objectClass=doloresType)`
	} else {
		filter = fmt.Sprintf(`(&(objectClass=doloresType)%s)`, filter)
	}
	return org.search(&searchRequest{org.parentDN(typeCategory(isUnit)),
		filter,
		[]string{`id`, `cn`, `description`}, nil,
		pageSize,
		cookie,
	})
}

func (org *Organization) searchPermission(filter string, isUnit bool, pageSize uint32, cookie []byte) (*SearchResult, error) {
	if len(filter) == 0 {
		filter = `(objectClass=permission)`
	} else {
		filter = fmt.Sprintf(`(&(objectClass=permission)%s)`, filter)
	}
	return org.search(&searchRequest{org.parentDN(permissionCategory(isUnit)),
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

func (org *Organization) searchUnit(filter string, containACL bool, pageSize uint32, cookie []byte) (*SearchResult, error) {
	var mAttrs []string
	if containACL {
		mAttrs = []string{`rbacType`}
	}

	if len(filter) == 0 {
		filter = `(objectClass=organizationalUnit)`
	} else {
		filter = fmt.Sprintf(`(&(objectClass=organizationalUnit)%s)`, filter)
	}

	return org.search(&searchRequest{org.parentDN(unit),
		filter,
		[]string{`id`, `cn`, `description`},
		mAttrs,
		pageSize,
		cookie,
	})
}

func sqConvertIDsToFilter(ids []string) (string, error) {
	return sqConvertArraysToFilter(`id`, ids)
}

func sqConvertArraysToFilter(label string, datas []string) (string, error) {
	if len(datas) == 0 {
		return ``, fmt.Errorf(`At least one %s`, label)
	}

	filter := `(|`
	for _, id := range datas {
		filter += fmt.Sprintf(`(%s=%s)`, label, id)
	}
	filter += `)`

	return filter, nil
}

package gorbacx

import (
	"fmt"
	"sync"

	"github.com/deckarep/golang-set"
)

type RBACX struct {
	mutex sync.Mutex
	roles map[string]*Role
}

func New() *RBACX {
	return &RBACX{
		roles: make(map[string]*Role, 0),
	}
}

func (rbacx *RBACX) Add(roles []*Role) {
	rbacx.mutex.Lock()
	defer rbacx.mutex.Unlock()

	for _, r := range roles {
		rbacx.roles[r.ID] = r
	}
}

func (rbacx *RBACX) Remove(ids []string) {
	rbacx.mutex.Lock()
	rbacx.mutex.Unlock()

	for _, id := range ids {
		delete(rbacx.roles, id)
	}
}

func (rbacx *RBACX) RoleByID(id string) (*Role, error) {
	rbacx.mutex.Lock()
	rbacx.mutex.Unlock()

	if role, ok := rbacx.roles[id]; ok {
		return role, nil
	}
	return nil, fmt.Errorf(`not found role id: %s`, id)
}

func (rbacx *RBACX) PermissionByID(id string, isUnit bool) (*Permission, error) {
	rbacx.mutex.Lock()
	rbacx.mutex.Unlock()

	for _, role := range rbacx.roles {
		if isUnit {
			for _, p := range role.unitPermissions {
				if p.ID == id {
					return p, nil
				}
			}
		} else {
			for _, p := range role.personPermissions {
				if p.ID == id {
					return p, nil
				}
			}
		}
	}

	return nil, fmt.Errorf(`not found permission id: %s`, id)
}

func (rbacx *RBACX) MatchedTypes(roleIDs []string, isUnit bool) []string {
	rbacx.mutex.Lock()
	rbacx.mutex.Unlock()

	set := mapset.NewSet()

	for _, id := range roleIDs {
		if role, ok := rbacx.roles[id]; ok {
			set = set.Union(role.matchedTypes(isUnit))
		} else {
			fmt.Printf(`[Warning-rbacx] cant't find role: %s, please add this role`, id)
		}
	}

	var types []string
	it := set.Iterator()
	for t := range it.C {
		types = append(types, t.(string))
	}

	return types
}

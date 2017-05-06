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

func (rbacx *RBACX) Remove(roles []*Role) {
	rbacx.mutex.Lock()
	rbacx.mutex.Unlock()

	for _, r := range roles {
		delete(rbacx.roles, r.ID)
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

func (rbacx *RBACX) PermissionByID(id string) (*Permission, error) {
	rbacx.mutex.Lock()
	rbacx.mutex.Unlock()

	for _, role := range rbacx.roles {
		for _, p := range role.permissions {
			if p.ID == id {
				return p, nil
			}
		}
	}
	return nil, fmt.Errorf(`not found permission id: %s`, id)
}

func (rbacx *RBACX) HasRefrenceType(id string, isUnit bool) string {
	for _, role := range rbacx.roles {
		for _, p := range role.permissions {
			if isUnit {
				if p.unitTypes.Contains(id) {
					return p.ID
				}
			} else {
				if p.personTypes.Contains(id) {
					return p.ID
				}
			}
		}
	}
	return ``
}

func (rbacx *RBACX) MatchedTypes(roleIDs []string) (unitTypes []string, personTypes []string) {
	rbacx.mutex.Lock()
	rbacx.mutex.Unlock()

	unit := mapset.NewSet()
	person := mapset.NewSet()

	for _, id := range roleIDs {
		if role, ok := rbacx.roles[id]; ok {
			u, p := role.matchedTypes()
			unit = unit.Union(u)
			person = person.Union(p)
		} else {
			fmt.Printf(`[Warning-rbacx] cant't find role: %s, please add this role`, id)
		}
	}

	uit := unit.Iterator()
	for t := range uit.C {
		unitTypes = append(unitTypes, t.(string))
	}

	pit := person.Iterator()
	for t := range pit.C {
		personTypes = append(personTypes, t.(string))
	}

	return unitTypes, personTypes
}

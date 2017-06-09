package gorbacx

import (
	"fmt"
	"sync"

	"github.com/deckarep/golang-set"
)

// RBACX ...
type RBACX struct {
	mutex sync.Mutex
	roles map[string]*Role
}

// New ...
func New() *RBACX {
	return &RBACX{
		roles: make(map[string]*Role, 0),
	}
}

func (rbacx *RBACX) Clear() {
	rbacx.mutex.Lock()
	defer rbacx.mutex.Unlock()

	rbacx.roles = make(map[string]*Role, 0)
}

// Add ...
func (rbacx *RBACX) Add(roles []*Role) {
	rbacx.mutex.Lock()
	defer rbacx.mutex.Unlock()

	for _, r := range roles {
		rbacx.roles[r.ID] = r
	}
}

// Remove ...
func (rbacx *RBACX) Remove(ids []string) {
	rbacx.mutex.Lock()
	rbacx.mutex.Unlock()

	for _, id := range ids {
		delete(rbacx.roles, id)
	}
}

// RoleByID ...
func (rbacx *RBACX) RoleByID(id string) (*Role, error) {
	rbacx.mutex.Lock()
	rbacx.mutex.Unlock()

	if role, ok := rbacx.roles[id]; ok {
		return role, nil
	}
	return nil, fmt.Errorf(`not found role id: %s`, id)
}

// RoleIDsByTypeID 有哪些Role包含此Type
func (rbacx *RBACX) RoleIDsByTypeID(id string, isUnit bool) []string {
	var ids []string
	for _, r := range rbacx.roles {
		if r.matchedTypes(isUnit).Contains(id) {
			ids = append(ids, id)
		}
	}
	return ids
}

// PermissionByID ...
func (rbacx *RBACX) PermissionByID(id string) (*Permission, error) {
	rbacx.mutex.Lock()
	rbacx.mutex.Unlock()

	for _, role := range rbacx.roles {
		// id 具有全局唯一性
		for _, p := range role.unitPermissions {
			if p.ID == id {
				return p, nil
			}
		}
		for _, p := range role.memberPermissions {
			if p.ID == id {
				return p, nil
			}
		}
	}

	return nil, fmt.Errorf(`not found permission id: %s`, id)
}

// MatchedTypes ...
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

// PrettyPrint ...
func (rbacx *RBACX) PrettyPrint() {
	fmt.Println(`----------------------rbacx PrettyPrint Begin----------------------`)
	fmt.Printf(`Role Count: %v`, len(rbacx.roles))
	fmt.Println()
	for _, role := range rbacx.roles {
		fmt.Printf(`Role: %s`, role.ID)
		fmt.Println()
		for _, v := range role.unitPermissions {
			fmt.Printf(`---up[%s]: %v`, v.ID, v.types.ToSlice())
			fmt.Println()
		}
		for _, v := range role.memberPermissions {
			fmt.Printf(`---mp[%s]: %v`, v.ID, v.types.ToSlice())
			fmt.Println()
		}
	}
	fmt.Println(`----------------------rbacx PrettyPrint End----------------------`)
}

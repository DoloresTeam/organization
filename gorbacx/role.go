package gorbacx

import "github.com/deckarep/golang-set"

type Role struct {
	ID          string
	permissions map[string]*Permission
}

func NewRole(id string, ps []*Permission) *Role {

	permissions := make(map[string]*Permission, len(ps))
	for _, p := range ps {
		permissions[p.ID] = p
	}

	return &Role{id, permissions}
}

func (r *Role) Add(ps []*Permission) {
	for _, p := range ps {
		r.permissions[p.ID] = p
	}
}

func (r *Role) Remove(ps []*Permission) {
	for _, p := range ps {
		delete(r.permissions, p.ID)
	}
}

func (r *Role) matchedTypes() (mapset.Set, mapset.Set) {

	unit := mapset.NewSet()
	person := mapset.NewSet()

	for _, v := range r.permissions {
		unit = unit.Union(v.unitTypes)
		person = person.Union(v.personTypes)
	}

	return unit, person
}

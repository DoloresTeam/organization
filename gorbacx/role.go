package gorbacx

import (
	"sync"

	"github.com/deckarep/golang-set"
)

// Role ...
type Role struct {
	mutex             sync.Mutex
	ID                string
	unitPermissions   map[string]*Permission
	memberPermissions map[string]*Permission
}

// NewRole ...
func NewRole(id string, ups, pps []*Permission) *Role {

	upMap := make(map[string]*Permission, len(ups))
	for _, p := range ups {
		upMap[p.ID] = p
	}

	ppMap := make(map[string]*Permission, len(pps))
	for _, p := range pps {
		ppMap[p.ID] = p
	}

	return &Role{ID: id,
		unitPermissions: upMap, memberPermissions: ppMap}
}

// Add ...
func (r *Role) Add(ps []*Permission, isUnit bool) {
	if isUnit {
		for _, p := range ps {
			r.unitPermissions[p.ID] = p
		}
	} else {
		for _, p := range ps {
			r.memberPermissions[p.ID] = p
		}
	}
}

// Remove ...
func (r *Role) Remove(ps []*Permission, isUnit bool) {
	if isUnit {
		for _, p := range ps {
			delete(r.unitPermissions, p.ID)
		}
	} else {
		for _, p := range ps {
			delete(r.memberPermissions, p.ID)
		}
	}
}

// Replace ...
func (r *Role) Replace(ps []*Permission, isUnit bool) {
	if isUnit {
		r.unitPermissions = make(map[string]*Permission, 0)
	} else {
		r.memberPermissions = make(map[string]*Permission, 0)
	}
	r.Add(ps, isUnit)
}

func (r *Role) matchedTypes() mapset.Set {
	set := mapset.NewSet()
	for _, v := range r.unitPermissions {
		set = set.Union(v.types)
	}
	for _, v := range r.memberPermissions {
		set = set.Union(v.types)
	}
	return set
}

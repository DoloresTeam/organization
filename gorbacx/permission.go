package gorbacx

import "github.com/deckarep/golang-set"

// Permission ...
type Permission struct {
	ID    string
	types mapset.Set
}

// NewPermission ...
func NewPermission(id string, types []string) *Permission {

	p := &Permission{id, nil}

	p.Replace(types)

	return p
}

// Add ...
func (p *Permission) Add(ids []string) {
	for _, id := range ids {
		p.types.Add(id)
	}
}

// Remove ...
func (p *Permission) Remove(ids []string) {
	for _, id := range ids {
		p.types.Remove(id)
	}
}

// Replace ...
func (p *Permission) Replace(ids []string) {
	p.types = mapset.NewSet()
	p.Add(ids)
}

// TypeIDs ...
func (p *Permission) TypeIDs() []string {

	it := p.types.Iter()

	var ids []string
	for id := range it {
		ids = append(ids, id.(string))
	}

	return ids
}

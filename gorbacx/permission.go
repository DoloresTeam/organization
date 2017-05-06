package gorbacx

import "github.com/deckarep/golang-set"

type Permission struct {
	ID    string
	types mapset.Set
}

func NewPermission(id string, types []string) *Permission {

	p := &Permission{id, nil}

	p.Replace(types)

	return p
}

func (p *Permission) Add(ids []string) {
	for _, id := range ids {
		p.types.Add(id)
	}
}

func (p *Permission) Remove(ids []string) {
	for _, id := range ids {
		p.types.Remove(id)
	}
}

func (p *Permission) Replace(ids []string) {
	p.types = mapset.NewSet()
	p.Add(ids)
}

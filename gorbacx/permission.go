package gorbacx

import "github.com/deckarep/golang-set"

type Permission struct {
	ID          string
	unitTypes   mapset.Set
	personTypes mapset.Set
}

func NewPermission(id string, unit, person []string) *Permission {

	p := &Permission{id, mapset.NewSet(), mapset.NewSet()}

	p.Add(unit, true)
	p.Add(person, false)

	return p
}

func (p *Permission) Add(ids []string, isUnit bool) {
	for _, id := range ids {
		if isUnit {
			p.unitTypes.Add(id)
		} else {
			p.personTypes.Add(id)
		}
	}
}

func (p *Permission) Remove(ids []string, isUnit bool) {
	for _, id := range ids {
		if isUnit {
			p.unitTypes.Remove(id)
		} else {
			p.personTypes.Remove(id)
		}
	}
}

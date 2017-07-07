package organization

import (
	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/util"
)

// adapter interface

// LoadPolicy ...
func (org *Organization) LoadPolicy(model model.Model) error {
	rs, err := org.AllRoles() // 获取所有的Role
	if err != nil {
		return err
	}
	for _, r := range rs {
		org.insertNewPolicyByRole(r)
	}

	return nil
}

// SavePolicy ...
func (org *Organization) SavePolicy(model model.Model) error {
	panic(`Organization doesn't support save policy`)
}

//
func (org *Organization) fetchAllowedTypesInRoles(rids []string) []string {
	var allowedTypes []string
	for _, rid := range rids {
		var types []string
		for _, policy := range org.enforcer.GetFilteredPolicy(0, rid) {
			types = append(types, policy[1])
		}
		allowedTypes = append(allowedTypes, types...)
	}
	util.ArrayRemoveDuplicates(&allowedTypes) // 去重
	return allowedTypes
}

func (org *Organization) fetchAllRolesByTypeID(tid string) []string {
	var rids []string
	for _, policy := range org.enforcer.GetFilteredPolicy(1, tid) {
		rids = append(rids, policy[1])
	}
	return rids
}

func (org *Organization) insertNewPolicyByRole(r map[string]interface{}) {
	permissionIDs := append(r[`upid`].([]string), r[`ppid`].([]string)...)
	searchResult, err := org.PermissionByIDs(permissionIDs)
	if err == nil {
		id := r[`id`].(string)
		for _, p := range searchResult.Data {
			types := p[`rbacType`].([]string)
			for _, t := range types {
				org.enforcer.AddPolicy([]string{id, t, `read`})
			}
		}
	}
}

func (org *Organization) removePolicyByRoleID(rid string) {
	org.enforcer.RemoveFilteredPolicy(0, rid)
}

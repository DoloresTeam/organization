package organization

import (
	"github.com/deckarep/golang-set"
)

// 更新权限和角色时会影响整个数据视图
func (org *Organization) refreshRBACIfNeeded(o, n []string) {

	oSet := set(o)
	nSet := set(n)
	if oSet.Equal(nSet) { // 没有变化
		return
	}
	org.RefreshRBAC() // 更新RBAC缓存
}

// 更新个人 类型，角色时影响数据视图

// 更新部门 类型，角色时影响数据视图

// utils
func set(v []string) mapset.Set {
	set := mapset.NewSet()
	for _, s := range v {
		set.Add(s)
	}
	return set
}

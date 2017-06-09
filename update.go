package organization

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/deckarep/golang-set"
	ldap "gopkg.in/ldap.v2"
)

const (
	AuditActionAdd    = 0
	AuditActionDel    = 1
	AuditActionUpdate = 2
)

const (
	AuditCategoryUnit   = 0
	AuditCategoryMember = 1
)

type AuditContent []map[string]interface{}

// 更新权限和角色时会影响整个数据视图
func (org *Organization) refreshRBACIfNeeded(o, n []string) {
	if set(o).Equal(set(n)) { // 没有变化
		return
	}
	org.RefreshRBAC() // 更新RBAC缓存
}

// utils
func set(v []string) mapset.Set {
	set := mapset.NewSet()
	for _, s := range v {
		set.Add(s)
	}
	return set
}

func vset(ss []interface{}) []string {
	var strs []string
	for _, s := range ss {
		strs = append(strs, s.(string))
	}
	return strs
}

func (org *Organization) relatedMIDs(tid string) ([]string, error) {
	return org.MemberIDsByTypeIDs([]string{tid})
}

func (org *Organization) logAddMember(member map[string]interface{}) error {
	return org.typeAdded(AuditCategoryMember, member)
}

func (org *Organization) logModifyMember(oMember, nMember map[string]interface{}) error {

	// 1. 判断用户角色有没有变化
	oRole := set(oMember[`rbacRole`].([]string))
	nRole := set(nMember[`rbacRole`].([]string))
	if !oRole.Equal(nRole) { // 用户的角色发生了变化
		// 计算发生变化前后的`Type`变化y
		oType := set(org.rbacx.MatchedTypes(oMember[`rbacRole`].([]string), false))
		nType := set(org.rbacx.MatchedTypes(nMember[`rbacRole`].([]string), false))

		addTypes := oType.Difference(nType)

		// 下面这2块代码，需要重构

		// 获取增加的用户ID
		mids, err := org.MemberIDsByTypeIDs(vset(addTypes.ToSlice()))
		if err != nil {
			return err
		}
		members, err := org.MemberByIDs(mids, false, false)
		if err != nil {
			return err
		}
		err = org.addAuditLog(AuditActionAdd, AuditCategoryMember, []string{nMember[`id`].(string)}, members)
		if err != nil {
			return err
		}

		// 获取删除的用户ID
		delTypes := nType.Difference(oType)

		mids, err = org.MemberIDsByTypeIDs(vset(delTypes.ToSlice()))
		if err != nil {
			return err
		}
		members, err = org.MemberByIDs(mids, false, false)
		if err != nil {
			return err
		}
		err = org.addAuditLog(AuditActionDel, AuditCategoryMember, []string{nMember[`id`].(string)}, members)
		if err != nil {
			return err
		}

		// 写两条日志
		// 1. 这个member 通讯录新增加了那几条访问权限
		// ADD|member[`id`].(string)|-----
		// 2. 这个member 删除了那几条访问权限
		// DEL|member[`id`].(string)|------
	}

	oType := oMember[`rbacType`].(string)
	nType := nMember[`rbacType`].(string)
	if oType != nType { // 修改了用户的类型
		org.typeChange(oType, nType, AuditCategoryMember, nMember)
		// 写日志
		// ADD|add|member
		// DEL|del|member
	}

	return org.addAuditLog(AuditActionUpdate, AuditCategoryMember, []string{nMember[`id`].(string)}, AuditContent{nMember})
}

func (org *Organization) logDelMember(id string, relatedMIDs []string) error {
	return org.addAuditLog(AuditActionDel, AuditCategoryMember, relatedMIDs, AuditContent{
		map[string]interface{}{
			`id`: id,
		},
	})
}

func (org *Organization) logAddUnit(unit map[string]interface{}) error {
	return org.typeAdded(AuditCategoryUnit, unit)
}

func (org *Organization) logModifyUnit(oUnit, nUnit map[string]interface{}) error {
	oType := oUnit[`rbacType`].(string)
	nType := nUnit[`rbacType`].(string)
	if oType != nType {
		org.typeChange(oType, nType, AuditCategoryUnit, nUnit)
	}
	mids, err := org.relatedMIDs(nType)
	if err != nil {
		return err
	}
	return org.addAuditLog(AuditActionUpdate, AuditCategoryUnit, mids, AuditContent{nUnit})
}

func (org *Organization) logDelUnit(id string, tid string) error {
	mids, err := org.relatedMIDs(tid)
	if err != nil {
		return err
	}
	return org.addAuditLog(AuditActionDel, AuditCategoryUnit, mids, AuditContent{
		map[string]interface{}{
			`id`: id,
		},
	})
}

func (org *Organization) typeAdded(category int, entry map[string]interface{}) error {
	mids, err := org.relatedMIDs(entry[`rbacType`].(string))
	if err != nil {
		return err
	}
	return org.addAuditLog(AuditActionAdd, category, mids, AuditContent{entry})
}

func (org *Organization) typeChange(oType, nType string, category int, entry map[string]interface{}) error {
	oMIDs, err := org.MemberIDsByTypeIDs([]string{oType})
	if err != nil {
		return err
	}
	nMIDs, err := org.MemberIDsByTypeIDs([]string{nType})
	if err != nil {
		return err
	}

	// 哪些用户失去了对当前用户的访问
	delMIDs := vset(set(oMIDs).Difference(set(nMIDs)).ToSlice())
	addMIDs := vset(set(nMIDs).Difference(set(oMIDs)).ToSlice())

	err = org.addAuditLog(AuditActionAdd, category, addMIDs, AuditContent{entry})
	if err != nil {
		return err
	}
	return org.addAuditLog(AuditActionDel, category, delMIDs, AuditContent{entry})
}

func (org *Organization) addAuditLog(action, category int, mids []string, content []map[string]interface{}) error {

	if len(mids) == 0 || len(content) == 0 {
		return nil
	}

	json, err := json.Marshal(content)
	if err != nil {
		return err
	}

	fmt.Println(mids)
	fmt.Println(string(json))

	aq := ldap.NewAddRequest(org.dn(generatorID(), audit))

	aq.Attribute(`objectClass`, []string{`audit`, `top`})
	aq.Attribute(`action`, []string{strconv.Itoa(action)})
	aq.Attribute(`category`, []string{strconv.Itoa(category)})
	aq.Attribute(`mid`, mids)
	aq.Attribute(`auditContent`, []string{string(json)})

	return org.l.Add(aq)
}

func (org *Organization) fetchAuditLog(memberID, lastedLogID string) ([]map[string]interface{}, error) {
	if len(memberID) == 0 || len(lastedLogID) == 0 {
		return nil, errors.New(`memberID && lastedLogID must not be empty`)
	}
	filter := fmt.Sprintf(`(&(mid=%s)(id>=%s))`, memberID, lastedLogID)
	sq := ldap.NewSearchRequest(org.parentDN(audit),
		ldap.ScopeSingleLevel, ldap.DerefAlways, 0, 0, false, filter,
		[]string{`id`, `action`, `auditContent`, `category`}, nil)
	sr, err := org.l.Search(sq)
	if err != nil {
		return nil, err
	}

	for _, entry := range sr.Entries {
		entry.PrettyPrint(2)
	}

	return nil, nil
}

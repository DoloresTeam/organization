package organization

// update 增量更新的逻辑说明
// 在权限变更， 或者角色变更的时候会引起每个用户的通讯录视图发生变化
// 增加角色，或者删除角色 不会更新版本号
import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/deckarep/golang-set"
	ldap "gopkg.in/ldap.v2"
)

const (
	AuditActionAdd    = `add`
	AuditActionDel    = `delete`
	AuditActionUpdate = `update`
)

const (
	AuditCategoryUnit   = `department`
	AuditCategoryMember = `member`
)

type AuditContent []map[string]interface{}

// 更新权限和角色时会影响整个数据视图
func (org *Organization) refreshRBACIfNeeded(o, n []string) {
	if set(o).Equal(set(n)) { // 没有变化
		return
	}
	err := org.RefreshRBAC() // 更新RBAC缓存
	if err == nil {
		func(dn string) {
			sq := ldap.NewSearchRequest(dn, ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, true, `(objectClass=audit)`, nil, nil)
			sr, _ := org.l.Search(sq)
			for _, r := range sr.Entries {
				org.l.Del(ldap.NewDelRequest(r.DN, nil))
			}
		}(org.parentDN(audit))
		org.latestResetVersion = newTimeStampVersion()
	} else {
		fmt.Print(`err: %s`, err.Error())
	}
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
	rids := org.rbacx.RoleIDsByTypeID(tid)
	if len(rids) == 0 {
		return nil, errors.New(`no role refrence this type`)
	}
	return org.MemberIDsByRoleIDs(rids)
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
		oType := set(org.rbacx.MatchedTypes(oMember[`rbacRole`].([]string)))
		nType := set(org.rbacx.MatchedTypes(nMember[`rbacRole`].([]string)))

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

	mids, err := org.relatedMIDs(nMember[`rbacType`].(string))
	if err != nil {
		return err
	}

	return org.addAuditLog(AuditActionUpdate, AuditCategoryMember, mids, AuditContent{nMember})
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

func (org *Organization) typeAdded(category string, entry map[string]interface{}) error {
	mids, err := org.relatedMIDs(entry[`rbacType`].(string))
	if err != nil {
		return err
	}
	return org.addAuditLog(AuditActionAdd, category, mids, AuditContent{entry})
}

func (org *Organization) typeChange(oType, nType string, category string, entry map[string]interface{}) error {
	oMIDs, err := org.relatedMIDs(oType)
	if err != nil {
		return err
	}
	nMIDs, err := org.relatedMIDs(nType)
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

func (org *Organization) addAuditLog(action, category string, mids []string, content []map[string]interface{}) error {

	if len(mids) == 0 || len(content) == 0 {
		return nil
	}

	for _, c := range content {
		delete(c, `rbacType`)
		delete(c, `rbacRole`)
		delete(c, `dn`)
	}

	json, err := json.Marshal(content)
	if err != nil {
		return err
	}

	id := generateNewID()
	aq := ldap.NewAddRequest(org.dn(id, audit))

	aq.Attribute(`objectClass`, []string{`audit`, `top`})
	aq.Attribute(`action`, []string{action})
	aq.Attribute(`category`, []string{category})
	aq.Attribute(`mid`, mids)
	aq.Attribute(`auditContent`, []string{string(json)})

	err = org.l.Add(aq)
	if err != nil {
		return err
	}

	go func() {
		org.organizationViewEvent <- mids
	}()

	return nil
}

func (org *Organization) fetchAuditLog(memberID, lastedLogID string) ([]map[string]interface{}, error) {
	filter := fmt.Sprintf(`(&(mid=%s)(createTimestamp>=%s))`, memberID, lastedLogID)
	sq := ldap.NewSearchRequest(org.parentDN(audit),
		ldap.ScopeSingleLevel, ldap.DerefAlways, 0, 0, false, filter,
		[]string{`createTimestamp`, `action`, `auditContent`, `category`}, nil)
	sr, err := org.l.Search(sq)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	for _, e := range sr.Entries {
		var content []map[string]interface{}
		err := json.Unmarshal(e.GetRawAttributeValue(`auditContent`), &content)
		if err != nil {
			continue
		}
		log := make(map[string]interface{}, 0)
		log[`content`] = content
		log[`createTimestamp`] = e.GetAttributeValue(`createTimestamp`)
		log[`action`] = e.GetAttributeValue(`action`)
		log[`category`] = e.GetAttributeValue(`category`)
		result = append(result, log)
	}

	return result, nil
}

func (org *Organization) IsValidVersion(version string) bool {
	return org.latestResetVersion == `` || version > org.latestResetVersion
}

func (org *Organization) GenerateChangeLogFromVersion(version string, mid string) (string, []map[string]interface{}, error) {
	if org.IsValidVersion(version) {
		logs, err := org.fetchAuditLog(mid, version)
		return newTimeStampVersion(), logs, err
	} else {
		return ``, nil, fmt.Errorf(`version invalid. latest reset version %s`, org.latestResetVersion)
	}
}

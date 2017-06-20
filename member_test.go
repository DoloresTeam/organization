package organization

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"testing"
)

func TestAddDelMember(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	m := md5.New()
	m.Write([]byte(`123456`))
	pwd := m.Sum(nil)

	roles, err := org.AllRoles()
	if err != nil {
		t.Fatal(err)
	}
	var rids []string
	for _, r := range roles {
		rids = append(rids, r[`id`].(string))
	}

	types, err := org.Types(true, 10, nil)
	if err != nil {
		t.Fatal(err)
	}
	var tids []string
	for _, t := range types.Data {
		tids = append(tids, t[`id`].(string))
	}

	units, err := org.AllUnit()
	if err != nil {
		t.Fatal(err)
	}
	var uids []string
	for _, u := range units {
		uids = append(uids, u[`id`].(string))
	}
	// 添加一个用户
	id, err := org.AddMember(map[string][]string{
		`name`:            []string{`JustForTTTTTesting`},
		`telephoneNumber`: []string{`13134564321`},
		`cn`:              []string{`王聪灵`},
		`email`:           []string{`heath.wang@dolores.store`},
		`title`:           []string{`Developer`},
		`thirdAccount`:    []string{`heath`},
		`thirdPassword`:   []string{`111`},
		`rbacRole`:        rids,
		`rbacType`:        tids,
		`unitID`:          uids,
		`userPassword`:    []string{hex.EncodeToString(pwd)},
	})

	if err != nil {
		t.Fatal(err)
	}

	// 使用手机号和密码登录 然后返回用户ID
	aid, err := org.AuthMember(`13134564321`, hex.EncodeToString(pwd))
	if err != nil {
		t.Fatal(err)
	}

	if aid != id {
		t.Fatal(errors.New(`Auth method occured error !!!`))
	}

	// 删除这个用户
	err = org.DelMember(id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFetchAllMembers(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	_, err := org.Members(50, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestModifyPassword(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	id := `b4vb0vh1scghuujqilo0`

	m := md5.New()
	m.Write([]byte(`123456`))
	op := hex.EncodeToString(m.Sum(nil))

	m.Reset()
	m.Write([]byte(`123456`))
	np := hex.EncodeToString(m.Sum(nil))

	err := org.ModifyPassword(id, op, np)
	if err != nil {
		t.Fatal(err)
	}

	member, err := org.MemberByID(id, false, false)
	if err != nil {
		t.Fatal(err)
	}
	tel := member[`telephoneNumber`].(string)

	_, err = org.AuthMember(tel, np)
	if err != nil {
		t.Fatal(err)
	}

	err = org.ModifyPassword(id, np, op) // 还原密码
	if err != nil {
		t.Fatal(err)
	}
}

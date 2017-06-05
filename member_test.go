package organization

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
)

func TestAddDelMember(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	m := md5.New()
	m.Write([]byte(`123456`))
	pwd := m.Sum(nil)

	// 添加一个用户
	id, err := org.AddMember(map[string][]string{
		`name`:            []string{`JustForTTTTTesting`},
		`telephoneNumber`: []string{`1888888888`},
		`cn`:              []string{`王聪灵`},
		`email`:           []string{`heath.wang@dolores.store`},
		`title`:           []string{`Developer`},
		`thirdAccount`:    []string{`heath`},
		`thirdPassword`:   []string{`111`},
		`rbacRole`:        []string{`b4ds07lhfpcr37ut14a0`},
		`rbacType`:        []string{`b4drqelhfpcqn7f7du5g`},
		`unitID`:          []string{`b4ds0t5hfpcr4h3thtd0`},
		`userPassword`:    []string{fmt.Sprintf(`{MD5}%s`, hex.EncodeToString(pwd))},
	})

	if err != nil {
		t.Fatal(err)
	}

	// 使用手机号和密码登录 然后返回用户ID
	aid, err := org.AuthMember(`1888888888`, hex.EncodeToString(pwd))
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

	sr, err := org.Members(30, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(sr.Data)
}

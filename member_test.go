package organization

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestAddMember(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	m := md5.New()
	m.Write([]byte(`123456`))
	pwd := m.Sum(nil)

	err = org.AddMember(map[string][]string{
		`name`:            []string{`巩祥`},
		`telephoneNumber`: []string{`13918839402`},
		`cn`:              []string{`Kevin.Gong`},
		`email`:           []string{`aoxianglele@icloud.com`},
		`title`:           []string{`Developer`},
		`rbacRole`:        []string{`b49jug06h301nm494sd0`},
		`rbacType`:        []string{`b49jtn06h301mgko5jo0`},
		`unitID`:          []string{`b49kdrg6h302hrpggg8g`},
		`userPassword`:    []string{fmt.Sprintf(`{MD5}%s`, hex.EncodeToString(pwd))},
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthMember(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	m := md5.New()
	m.Write([]byte(`123456`))
	pwd := m.Sum(nil)

	id, err := org.AuthMember(`13918839401`, fmt.Sprintf(`{MD5}%s`, hex.EncodeToString(pwd)))
	if err != nil {
		t.Fatal(err)
	}

	member, err := org.MemberByID(id, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(member)
}

func TestFetchMemberRoles(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	ids, err := org.RoleIDsByMemberID(`b49kehg6h302jg98oi70`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ids)
}

func TestFetchOrgMemberByID(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	ids, err := org.OrganizationMemberByMemberID(`b49kehg6h302jg98oi70`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ids)
}

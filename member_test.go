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

	t.Log(hex.EncodeToString(pwd))

	err = org.AddMember(map[string][]string{
		`name`:            []string{`Heath.Wang`},
		`telephoneNumber`: []string{`18627800585`},
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

	id, err := org.AuthMember(`18627800585`, hex.EncodeToString(pwd))
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

	ids, err := org.MemberByID(`b4dsitlhfpcs1aerd0l0`, true) // org.OrganizationMemberByMemberID(`b4dsitlhfpcs1aerd0l0`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ids)
}

func TestFetchAllMembers(t *testing.T) {

	if testing.Short() {
		t.SkipNow()
	}

	org, err := neworg()
	if err != nil {
		t.Fatal(err)
	}

	sr, err := org.Members(30, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(sr.Data)
}

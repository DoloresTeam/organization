package organization

import "testing"

func TestRole(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	upes, err := org.Permissions(true, 2, nil)
	if err != nil {
		t.Fatal(err)
	}
	var ups []string
	for _, u := range upes.Data {
		ups = append(ups, u[`id`].(string))
	}

	ppes, err := org.Permissions(false, 2, nil)
	if err != nil {
		t.Fatal(err)
	}
	var pps []string
	for _, u := range ppes.Data {
		pps = append(pps, u[`id`].(string))
	}

	id, err := org.AddRole(`Test`, `Test and Role`, ups, pps)
	if err != nil {
		t.Fatal(err)
	}

	err = org.ModifyRole(id, ``, ``, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = org.DelRole(id)
	if err != nil {
		t.Fatal(err)
	}
}

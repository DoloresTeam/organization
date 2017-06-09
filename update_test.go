package organization

import "testing"

func TestUpdate(t *testing.T) {

	_, err := org.fetchAuditLog(`b4shdfthfpcg8d3knb5g`, `b4sj7edhfpci55od4mag`)
	if err != nil {
		t.Fatal(err)
	}
}

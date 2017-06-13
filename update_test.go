package organization

import (
	"testing"
)

func TestUpdate(t *testing.T) {
	// b4vb7p91scghuujqim3g
	mids, err := org.relatedMIDs(`b4va01h1scghr2vkka4g`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(mids)
}

package train

import (
	"testing"
)

func TestListArchives(t *testing.T) {
	err := ListArchives()
	if err != nil {
		t.Errorf("unable to list buckets: %v", err)
	}
}

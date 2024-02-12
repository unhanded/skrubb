package pivot_test

import (
	"testing"

	"github.com/unhanded/skrubb/internal/app/libskrubb/pivot"
)

func TestNewSkrubbRoot(t *testing.T) {
	fp, err := pivot.NewSkrubbRoot("banana")
	if err != nil {
		t.Error(err)
	}
	if len(fp) < 5 {
		t.Error("Filepath too short, something must've gone wrong.")
	} else {
		t.Logf("New dir filepath: %s", fp)
	}
}

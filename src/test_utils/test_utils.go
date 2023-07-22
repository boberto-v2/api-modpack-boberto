// TODO: When we init the CI/D at github actions,
// we need to decide the best practice to us.
// If we uses test with tags or a function to determine if the program is running at CI.
// check package test_utils
// at this time i will use this approach to not break the CI/CD rules
// https://stackoverflow.com/questions/24030059/skip-some-tests-with-go-test

package test_utils

import (
	"os"
	"testing"
)

func SkipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
	// and we put the PG_URI env here if we attack this form.
	t.Setenv("PG_URI", "postgres://root:test@127.0.0.1:5555/test?sslmode=disable")
}

package testutils

import (
	"fmt"
	"net/http"
	"testing"
)

// Given a request that is expected to return a single id,
// reads the request's body and returns the id
func IdOf(t *testing.T, resp *http.Response) int {
	var id int
	if _, err := fmt.Fscan(resp.Body, &id); err != nil {
		t.Fatal(err)
	}

	return id
}

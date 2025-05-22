package session

import (
	"net/http"
)

func IsAuthenticated(r *http.Request) bool {
	username, err := GetSessionValue(r, "username")
	if err != nil || username == nil {
		return false
	}
	return true
}

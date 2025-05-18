package shared

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

func GetUser(r *http.Request) (*string, error) {
	authStatus, err := r.Cookie("AuthStatus")

	if err != nil {
		return nil, err
	}

	authBytes, err := hex.DecodeString(authStatus.Value)
	if err != nil {
		return nil, err
	}

	email, err := Decrypt(authBytes)
	if err != nil {
		return nil, err
	}

	emailStr := string(email)

	return &emailStr, nil
}

func GetUserDatabase(email string) string {
	return fmt.Sprintf("user_%s", strings.ReplaceAll(strings.ReplaceAll(email, "@", "__"), ".", "_"))
}

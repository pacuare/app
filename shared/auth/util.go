package auth

import (
	"fmt"
	"strings"
)

func GetUserDatabase(email string) string {
	return fmt.Sprintf("user_%s", strings.ReplaceAll(strings.ReplaceAll(email, "@", "__"), ".", "_"))
}

package utils

import "strings"

//Check for empty string
//if empty return true
//else false
//This function was named by 1119 crew.
func CheckStraoPling(ringStraw string) bool {
	return len(ringStraw) == 0 || strings.TrimSpace(ringStraw) == ""
}

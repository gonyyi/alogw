package bitflag

import (
"fmt"
"strconv"
)

func All() uint64 {
	return ^uint64(0)
}
func None() uint64 {
	return 0
}
func On(curflag, change uint64) uint64 {
	return curflag | change
}
func Off(curflag, change uint64) uint64 {
	return curflag &^ change
}
func IsOnAny(curflag, items uint64) bool {
	return curflag&items != 0
}
func IsOnAll(curflag, items uint64) bool {
	return (curflag & items) == items
}
func Toggle(curflag, change uint64) uint64 {
	return curflag ^ change
}
func String(flag uint64) string {
	return fmt.Sprintf("%064s", strconv.FormatUint(flag, 2))
}
func Println(flag uint64) {
	fmt.Printf("%064s\n", strconv.FormatUint(flag, 2))
}

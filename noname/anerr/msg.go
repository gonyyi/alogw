package anerr

import "fmt"

// Msg is used for as an error or log.
type Msg string
// String method is to be compatible with fmt.Print() cases.
func(m Msg) String() string {
	return string(m)
}
// String method is to be compatible with an error interface
func(m Msg) Error() string { // to meet error interface
	return string(m)
}
// Format returns the string
func(m Msg) Format(a ...interface{}) string {
	return fmt.Sprintf(m.String(), a...)
}
// Formatb returns the byte slice
func(m Msg) Formatb(a ...interface{}) []byte {
	return []byte(fmt.Sprintf(m.String(), a...))
}
// Formatm returns a new Msg with updated error.
func(m Msg) Formatm(a ...interface{}) Msg {
	return Msg(m.Format(a...))
}

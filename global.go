package alogw

// CONST: bufwriter size
const (
	_ = 1 << (10 * iota)
	KB
	MB
	GB
	TB
	PB
	EB
)

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	ErrCannotCreateDir Err = `cannot create directory`
)

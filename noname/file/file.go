package file

const(
	OX = 1 << iota
	OW
	OR
	GX
	GW
	GR
	UX
	UW
	UR
	URW = UR|UW
	GRW = GR|GW
	ORW = OR|OW
	URWX = UR|UW|UX
	GRWX = GR|GW|GX
	ORWX = OR|OW|OX
)


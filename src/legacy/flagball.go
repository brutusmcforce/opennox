package legacy

var (
	Sub_4E8290 func(a1 byte, a2 uint16) int
)

func sub_4E8290(a1 int8, a2 int16) int {
	return Sub_4E8290(byte(a1), uint16(a2))
}

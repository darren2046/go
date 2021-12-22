package golanglibs

import (
	"net"

	"github.com/mitchellh/go-vnc"
)

type VNCCfg struct {
	Password string
	Delay    interface{} // seconds to wait after actions
}

type vncStruct struct {
	nc    net.Conn
	VC    *vnc.ClientConn
	lastX int
	lastY int
	delay interface{} // wait after actions
}

type vncNonAsciiKeyMapStruct struct {
	Shift   uint32
	Enter   uint32
	Windows uint32
	Left    uint32
	Up      uint32
	Down    uint32
	Right   uint32
	Delete  uint32
	Escape  uint32
	Tab     uint32
	Control uint32
	Alt     uint32
	F1      uint32
	F2      uint32
	F3      uint32
	F4      uint32
	F5      uint32
	F6      uint32
	F7      uint32
	F8      uint32
	F9      uint32
	F10     uint32
	F11     uint32
	F12     uint32
}

var vncNonAsciiKeyMap = &vncNonAsciiKeyMapStruct{
	Shift:   0xffe1,
	Enter:   0xff0d,
	Windows: 0xffe9,
	Left:    0xff51,
	Up:      0xff52,
	Down:    0xff54,
	Right:   0xff53,
	Delete:  0xff08,
	Escape:  0xff1b,
	Tab:     0xff09,
	Control: 0xffe4,
	Alt:     0xffe9,
	F1:      0xffbe,
	F2:      0xffbf,
	F3:      0xffc0,
	F4:      0xffc1,
	F5:      0xffc2,
	F6:      0xffc3,
	F7:      0xffc4,
	F8:      0xffc5,
	F9:      0xffc6,
	F10:     0xffc7,
	F11:     0xffc8,
	F12:     0xffc9,
}

var vncAsciiKeyMap = map[string][]uint32{
	"a": {0x0061},
	"b": {0x0062},
	"c": {0x0063},
	"d": {0x0064},
	"e": {0x0065},
	"f": {0x0066},
	"g": {0x0067},
	"h": {0x0068},
	"i": {0x0069},
	"j": {0x006a},
	"k": {0x006b},
	"l": {0x006c},
	"m": {0x006d},
	"n": {0x006e},
	"o": {0x006f},
	"p": {0x0070},
	"q": {0x0071},
	"r": {0x0072},
	"s": {0x0073},
	"t": {0x0074},
	"u": {0x0075},
	"v": {0x0076},
	"w": {0x0077},
	"x": {0x0078},
	"y": {0x0079},
	"z": {0x007a},

	"A": {vncNonAsciiKeyMap.Shift, 0x0061},
	"B": {vncNonAsciiKeyMap.Shift, 0x0062},
	"C": {vncNonAsciiKeyMap.Shift, 0x0063},
	"D": {vncNonAsciiKeyMap.Shift, 0x0064},
	"E": {vncNonAsciiKeyMap.Shift, 0x0065},
	"F": {vncNonAsciiKeyMap.Shift, 0x0066},
	"G": {vncNonAsciiKeyMap.Shift, 0x0067},
	"H": {vncNonAsciiKeyMap.Shift, 0x0068},
	"I": {vncNonAsciiKeyMap.Shift, 0x0069},
	"J": {vncNonAsciiKeyMap.Shift, 0x006a},
	"K": {vncNonAsciiKeyMap.Shift, 0x006b},
	"L": {vncNonAsciiKeyMap.Shift, 0x006c},
	"M": {vncNonAsciiKeyMap.Shift, 0x006d},
	"N": {vncNonAsciiKeyMap.Shift, 0x006e},
	"O": {vncNonAsciiKeyMap.Shift, 0x006f},
	"P": {vncNonAsciiKeyMap.Shift, 0x0070},
	"Q": {vncNonAsciiKeyMap.Shift, 0x0071},
	"R": {vncNonAsciiKeyMap.Shift, 0x0072},
	"S": {vncNonAsciiKeyMap.Shift, 0x0073},
	"T": {vncNonAsciiKeyMap.Shift, 0x0074},
	"U": {vncNonAsciiKeyMap.Shift, 0x0075},
	"V": {vncNonAsciiKeyMap.Shift, 0x0076},
	"W": {vncNonAsciiKeyMap.Shift, 0x0077},
	"X": {vncNonAsciiKeyMap.Shift, 0x0078},
	"Y": {vncNonAsciiKeyMap.Shift, 0x0079},
	"Z": {vncNonAsciiKeyMap.Shift, 0x007a},

	"1": {0x0031},
	"2": {0x0032},
	"3": {0x0033},
	"4": {0x0034},
	"5": {0x0035},
	"6": {0x0036},
	"7": {0x0037},
	"8": {0x0038},
	"9": {0x0039},
	"0": {0x0030},

	"!": {vncNonAsciiKeyMap.Shift, 0x0031},
	"@": {vncNonAsciiKeyMap.Shift, 0x0032},
	"#": {vncNonAsciiKeyMap.Shift, 0x0033},
	"$": {vncNonAsciiKeyMap.Shift, 0x0034},
	"%": {vncNonAsciiKeyMap.Shift, 0x0035},
	"^": {vncNonAsciiKeyMap.Shift, 0x0036},
	"&": {vncNonAsciiKeyMap.Shift, 0x0037},
	"*": {vncNonAsciiKeyMap.Shift, 0x0038},
	"(": {vncNonAsciiKeyMap.Shift, 0x0039},
	")": {vncNonAsciiKeyMap.Shift, 0x0030},

	" ": {0x0020},

	"`":  {0x0060},
	"-":  {0x002d},
	"=":  {0x003d},
	"[":  {0x005b},
	"]":  {0x005d},
	"\\": {0x005c},
	";":  {0x003b},
	"'":  {0x0027},
	",":  {0x002c},
	".":  {0x002e},
	"/":  {0x002f},

	"~":  {vncNonAsciiKeyMap.Shift, 0x0060},
	"_":  {vncNonAsciiKeyMap.Shift, 0x002d},
	"+":  {vncNonAsciiKeyMap.Shift, 0x003d},
	"{":  {vncNonAsciiKeyMap.Shift, 0x005b},
	"}":  {vncNonAsciiKeyMap.Shift, 0x005d},
	"|":  {vncNonAsciiKeyMap.Shift, 0x005c},
	":":  {vncNonAsciiKeyMap.Shift, 0x003b},
	"\"": {vncNonAsciiKeyMap.Shift, 0x0027},
	"<":  {vncNonAsciiKeyMap.Shift, 0x002c},
	">":  {vncNonAsciiKeyMap.Shift, 0x002e},
	"?":  {vncNonAsciiKeyMap.Shift, 0x002f},
}

func getVNC(server string, cfg ...VNCCfg) *vncStruct {
	nc, err := net.Dial("tcp", server)
	Panicerr(err)

	vcc := &vnc.ClientConfig{
		Exclusive: false,
	}
	var delay interface{} = 0
	if len(cfg) != 0 {
		if cfg[0].Password != "" {
			vcc.Auth = []vnc.ClientAuth{
				&vnc.PasswordAuth{
					Password: cfg[0].Password,
				},
			}
		}
		delay = cfg[0].Delay
	}

	vc, err := vnc.Client(nc, vcc)
	if err != nil {
		nc.Close()
		Panicerr(err)
	}

	return &vncStruct{
		nc:    nc,
		VC:    vc,
		delay: delay,
	}
}

func (m *vncStruct) Close() {
	m.VC.Close()
	m.nc.Close()
}

func (m *vncStruct) Move(x, y int) *vncStruct {
	m.lastX = x
	m.lastY = y
	err := m.VC.PointerEvent(0, Uint16(x), Uint16(y))
	Panicerr(err)
	sleep(m.delay)
	return m
}

func (m *vncStruct) Click() *vncStruct {
	err := m.VC.PointerEvent(1, Uint16(m.lastX), Uint16(m.lastY))
	Panicerr(err)
	err = m.VC.PointerEvent(0, Uint16(m.lastX), Uint16(m.lastY))
	Panicerr(err)
	sleep(m.delay)
	return m
}

func (m *vncStruct) RightClick() *vncStruct {
	err := m.VC.PointerEvent(4, Uint16(m.lastX), Uint16(m.lastY))
	Panicerr(err)
	err = m.VC.PointerEvent(0, Uint16(m.lastX), Uint16(m.lastY))
	Panicerr(err)
	sleep(m.delay)
	return m
}

func (m *vncStruct) Input(s string) *vncStruct {
	for i := 0; i < len(s); i++ {
		c := string([]byte{s[i]})
		// 顺序按那数组里面的按键, 然后反顺序放开
		if Map(vncAsciiKeyMap).Has(c) {
			for _, k := range vncAsciiKeyMap[c] {
				m.VC.KeyEvent(k, true)
			}
			for i := len(vncAsciiKeyMap[c]) - 1; i >= 0; i-- {
				m.VC.KeyEvent(vncAsciiKeyMap[c][i], false)
			}
		}
		if c == "\n" {
			m.Key().Enter()
		}
	}
	sleep(m.delay)
	return m
}

type vncKeyStruct struct {
	vc    *vnc.ClientConn
	delay interface{}
}

func (m *vncStruct) Key() *vncKeyStruct {
	return &vncKeyStruct{vc: m.VC, delay: m.delay}
}

func (m *vncKeyStruct) Enter() *vncKeyStruct {
	m.vc.KeyEvent(vncNonAsciiKeyMap.Enter, true)
	m.vc.KeyEvent(vncNonAsciiKeyMap.Enter, false)
	sleep(m.delay)
	return m
}

func (m *vncKeyStruct) CtrlA() *vncKeyStruct {
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, true)
	m.vc.KeyEvent(vncAsciiKeyMap["a"][0], true)
	m.vc.KeyEvent(vncAsciiKeyMap["a"][0], false)
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, false)
	sleep(m.delay)
	return m
}

func (m *vncKeyStruct) CtrlC() *vncKeyStruct {
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, true)
	m.vc.KeyEvent(vncAsciiKeyMap["c"][0], true)
	m.vc.KeyEvent(vncAsciiKeyMap["c"][0], false)
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, false)
	sleep(m.delay)
	return m
}

func (m *vncKeyStruct) CtrlV() *vncKeyStruct {
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, true)
	m.vc.KeyEvent(vncAsciiKeyMap["v"][0], true)
	m.vc.KeyEvent(vncAsciiKeyMap["v"][0], false)
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, false)
	sleep(m.delay)
	return m
}

func (m *vncKeyStruct) CtrlZ() *vncKeyStruct {
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, true)
	m.vc.KeyEvent(vncAsciiKeyMap["z"][0], true)
	m.vc.KeyEvent(vncAsciiKeyMap["z"][0], false)
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, false)
	sleep(m.delay)
	return m
}

func (m *vncKeyStruct) CtrlX() *vncKeyStruct {
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, true)
	m.vc.KeyEvent(vncAsciiKeyMap["x"][0], true)
	m.vc.KeyEvent(vncAsciiKeyMap["x"][0], false)
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, false)
	sleep(m.delay)
	return m
}

func (m *vncKeyStruct) CtrlF() *vncKeyStruct {
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, true)
	m.vc.KeyEvent(vncAsciiKeyMap["f"][0], true)
	m.vc.KeyEvent(vncAsciiKeyMap["f"][0], false)
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, false)
	sleep(m.delay)
	return m
}

func (m *vncKeyStruct) CtrlD() *vncKeyStruct {
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, true)
	m.vc.KeyEvent(vncAsciiKeyMap["d"][0], true)
	m.vc.KeyEvent(vncAsciiKeyMap["d"][0], false)
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, false)
	sleep(m.delay)
	return m
}

func (m *vncKeyStruct) CtrlS() *vncKeyStruct {
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, true)
	m.vc.KeyEvent(vncAsciiKeyMap["s"][0], true)
	m.vc.KeyEvent(vncAsciiKeyMap["s"][0], false)
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, false)
	sleep(m.delay)
	return m
}

func (m *vncKeyStruct) CtrlR() *vncKeyStruct {
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, true)
	m.vc.KeyEvent(vncAsciiKeyMap["r"][0], true)
	m.vc.KeyEvent(vncAsciiKeyMap["r"][0], false)
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, false)
	sleep(m.delay)
	return m
}

func (m *vncKeyStruct) CtrlE() *vncKeyStruct {
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, true)
	m.vc.KeyEvent(vncAsciiKeyMap["e"][0], true)
	m.vc.KeyEvent(vncAsciiKeyMap["e"][0], false)
	m.vc.KeyEvent(vncNonAsciiKeyMap.Control, false)
	sleep(m.delay)
	return m
}

func (m *vncKeyStruct) Delete() *vncKeyStruct {
	m.vc.KeyEvent(vncNonAsciiKeyMap.Delete, true)
	m.vc.KeyEvent(vncNonAsciiKeyMap.Delete, false)
	sleep(m.delay)
	return m
}

func (m *vncKeyStruct) Tab() *vncKeyStruct {
	m.vc.KeyEvent(vncNonAsciiKeyMap.Tab, true)
	m.vc.KeyEvent(vncNonAsciiKeyMap.Tab, false)
	sleep(m.delay)
	return m
}

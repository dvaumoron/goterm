// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package term

/*
Some simple functions to add colors and attributes to terminals.

The base colors are types implementing the Stringer interface, this makes
it very simple to give a color to arbitrary strings. Also handy to have the raw string still
available for comparisons and such.

	g := Green("Green world")
	fmt.Println("Hello",g)
	fmt.Println(Red("Warning!"))

	var col fmt.Stringer
	switch {
	case atk == 0:
		col = Blue("5 FADE OUT")
	case atk < 4:
		col = Green("4 DOUBLE TAKE")
	case atk <10:
		col = Yellow("3 ROUND HOUSE")
	case atk <50:
		col = Red("2 FAST PACE")
	case atk >= 50:
		col = Blinking("1 COCKED PISTOL")
	}
	fmt.Println("Defcon: ",col)
*/

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type empty = struct{}

type stringer interface {
	String() string
}

// colorEnable toggles colors on/off.
var colorEnable = true

// ColorEnable activates the terminal colors , this is the default.
func ColorEnable() {
	colorEnable = true
}

// ColorDisable disables the terminal colors.
func ColorDisable() {
	colorEnable = false
}

// Terminal Color and modifier codes
const (
	CSI    = "\033["
	SepCSI = ";"
	EndCSI = "m"

	FgBlack   = "30"
	FgRed     = "31"
	FgGreen   = "32"
	FgYellow  = "33"
	FgBlue    = "34"
	FgMagenta = "35"
	FgCyan    = "36"
	FgWhite   = "37"
	FgDefault = "39"
	F256      = "38"
	BgBlack   = "40"
	BgRed     = "41"
	BgGreen   = "42"
	BgYellow  = "43"
	BgBlue    = "44"
	BgMagenta = "45"
	BgCyan    = "46"
	BgWhite   = "47"
	BgDefault = "49"
	Bg256     = "48"
	Blink     = "5"
	Ital      = "3"
	Underln   = "4"
	Faint     = "2"
	Bld       = "1"
	NoMode    = "0"
)

const (
	sepCSI = ';'
	endCSI = 'm'

	csiLen  = len(CSI)
	baseLen = 2 * (csiLen + 1)

	eraser       = CSI + FgDefault + SepCSI + BgDefault + EndCSI
	colorBaseLen = csiLen + 2 + len(eraser)

	mod256          = ";5;"
	start256        = CSI + F256 + mod256
	sep256          = SepCSI + Bg256 + mod256
	eraser256       = CSI + FgDefault + mod256 + BgDefault + mod256 + EndCSI
	eraser256Len    = len(eraser256)
	color256BaseLen = len(start256) + len(sep256) + 1 + eraser256Len

	modRGB          = ";2;"
	startRGB        = CSI + F256 + modRGB
	colorRGBBaseLen = len(startRGB) + 3 + eraser256Len
)

// Standard colors
// Foreground

// Green implements the Stringer interface to print string foreground in Green color.
type Green string

// Blue implements the Stringer interface to print string foreground in Blue color.
type Blue string

// Red implements the Stringer interface to print string foreground in Red color.
type Red string

// Yellow implements the Stringer interface to print string foreground in Yellow color.
type Yellow string

// Magenta implements the Stringer interface to print string foreground in Magenta color.
type Magenta string

// Cyan implements the Stringer interface to print string foreground in Cyan color.
type Cyan string

// White implements the Stringer interface to print string foreground in White color.
type White string

// Black implements the Stringer interface to print string foreground in Black color.
type Black string

// Random implements the Stringer interface to print string foreground in Random color.
type Random string

// Background

// BGreen implements the Stringer interface to print string background in Green color.
type BGreen string

// BBlue implements the Stringer interface to print string background in Blue color.
type BBlue string

// BRed implements the Stringer interface to print string background in Red color.
type BRed string

// BYellow implements the Stringer interface to print string background in Yellow color.
type BYellow string

// BRandom implements the Stringer interface to print string background in Random color.
type BRandom string

// BMagenta implements the Stringer interface to print string background in Magenta color.
type BMagenta string

// BCyan implements the Stringer interface to print string background in Cyan color.
type BCyan string

// BWhite implements the Stringer interface to print string background in White color.
type BWhite string

// BBlack implements the Stringer interface to print string background in Black color.
type BBlack string

// Set color

// Color is the type returned by the colour setters to print any terminal colour.
type Color string

// ColorRandom implements the Stringer interface to print string Random color.
type ColorRandom string

// Color256Random implements the Stringer interface to print string random 256 color Term style.
type Color256Random string

// Some modifiers

// Blinking implements the Stringer interface to print string in Blinking mode.
type Blinking string

// Underline implements the Stringer interface to print string in Underline mode.
type Underline string

// Bold implements the Stringer interface to print string in Bold mode.
type Bold string

//type Bright string -- Doesn't seem to work well

// Italic implements the Stringer interface to print string foreground in Italic color.
type Italic string

// colConcat add beginning and ending color mode modifer
func colConcat(mode string, s string, nMode string) string {
	if !colorEnable {
		return s
	}
	var buffer strings.Builder
	buffer.Grow(baseLen + len(mode) + len(s) + len(nMode))
	buffer.WriteString(CSI)
	buffer.WriteString(mode)
	buffer.WriteByte(endCSI)
	buffer.WriteString(s)
	buffer.WriteString(CSI)
	buffer.WriteString(nMode)
	buffer.WriteByte(endCSI)
	return buffer.String()
}

// colConcatByte add beginning and ending color mode modifer (byte slice)
func colConcatByte(mode string, b []byte, nMode string) []byte {
	if !colorEnable {
		return b
	}
	buffer := make([]byte, 0, baseLen+len(mode)+len(b)+len(nMode))
	buffer = append(buffer, CSI...)
	buffer = append(buffer, mode...)
	buffer = append(buffer, endCSI)
	buffer = append(buffer, b...)
	buffer = append(buffer, CSI...)
	buffer = append(buffer, nMode...)
	buffer = append(buffer, endCSI)
	return buffer
}

// Stringers for all the base colors , just fill it in with something and print it
// Foreground

// String implements the Stringer interface for type Green.
func (c Green) String() string {
	return colConcat(FgGreen, string(c), FgDefault)
}

func GreenByte(b []byte) []byte {
	return colConcatByte(FgGreen, b, FgDefault)
}

// Greenf returns a Green formatted string.
func Greenf(format string, a ...any) string {
	return Green(fmt.Sprintf(format, a...)).String()
}

// String implements the Stringer interface for type Blue.
func (c Blue) String() string {
	return colConcat(FgBlue, string(c), FgDefault)
}

func BlueByte(b []byte) []byte {
	return colConcatByte(FgBlue, b, FgDefault)
}

// Bluef returns a Blue formatted string.
func Bluef(format string, a ...any) string {
	return Blue(fmt.Sprintf(format, a...)).String()
}

// String implements the Stringer interface for type Red.
func (c Red) String() string {
	return colConcat(FgRed, string(c), FgDefault)
}

func RedByte(b []byte) []byte {
	return colConcatByte(FgRed, b, FgDefault)
}

// Redf returns a Red formatted string.
func Redf(format string, a ...any) string {
	return Red(fmt.Sprintf(format, a...)).String()
}

// String implements the Stringer interface for type Yellow.
func (c Yellow) String() string {
	return colConcat(FgYellow, string(c), FgDefault)
}

func YellowByte(b []byte) []byte {
	return colConcatByte(FgYellow, b, FgDefault)
}

// Yellowf returns a Yellow formatted string.
func Yellowf(format string, a ...any) string {
	return Yellow(fmt.Sprintf(format, a...)).String()
}

// String implements the Stringer interface for type Magenta.
func (c Magenta) String() string {
	return colConcat(FgMagenta, string(c), FgDefault)
}

func MagentaByte(b []byte) []byte {
	return colConcatByte(FgMagenta, b, FgDefault)
}

// Magentaf returns a Magenta formatted string.
func Magentaf(format string, a ...any) string {
	return Magenta(fmt.Sprintf(format, a...)).String()
}

// String implements the Stringer interface for type White.
func (c White) String() string {
	return colConcat(FgWhite, string(c), FgDefault)
}

func WhiteByte(b []byte) []byte {
	return colConcatByte(FgWhite, b, FgDefault)
}

// Whitef returns a White formatted string.
func Whitef(format string, a ...any) string {
	return White(fmt.Sprintf(format, a...)).String()
}

// String implements the Stringer interface for type Black.
func (c Black) String() string {
	return colConcat(FgBlack, string(c), FgDefault)
}

func BlackByte(b []byte) []byte {
	return colConcatByte(FgBlack, b, FgDefault)
}

// Blackf returns a Black formatted string.
func Blackf(format string, a ...any) string {
	return Black(fmt.Sprintf(format, a...)).String()
}

// String implements the Stringer interface for type Cyan.
func (c Cyan) String() string {
	return colConcat(FgCyan, string(c), FgDefault)
}

func CyanByte(b []byte) []byte {
	return colConcatByte(FgCyan, b, FgDefault)
}

// Cyanf returns a Cyan formatted string.
func Cyanf(format string, a ...any) string {
	return Cyan(fmt.Sprintf(format, a...)).String()
}

// Background

// String implements the Stringer interface for type BGreen.
func (c BGreen) String() string {
	return colConcat(BgGreen, string(c), BgDefault)
}

func BGreenByte(b []byte) []byte {
	return colConcatByte(BgGreen, b, BgDefault)
}

// BGreenf returns a BGreen formatted string.
func BGreenf(format string, a ...any) string {
	return BGreen(fmt.Sprintf(format, a...)).String()
}

// String implements the Stringer interface for type BBlue.
func (c BBlue) String() string {
	return colConcat(BgBlue, string(c), BgDefault)
}

func BBlueByte(b []byte) []byte {
	return colConcatByte(BgBlue, b, BgDefault)
}

// BBluef returns a BBlue formatted string.
func BBluef(format string, a ...any) string {
	return BBlue(fmt.Sprintf(format, a...)).String()
}

// String implements the Stringer interface for type BRed.
func (c BRed) String() string {
	return colConcat(BgRed, string(c), BgDefault)
}

func BRedByte(b []byte) []byte {
	return colConcatByte(BgRed, b, BgDefault)
}

// BRedf returns a BRed formatted string.
func BRedf(format string, a ...any) string {
	return BRed(fmt.Sprintf(format, a...)).String()
}

// String implements the Stringer interface for type BYellow.
func (c BYellow) String() string {
	return colConcat(BgYellow, string(c), BgDefault)
}

func BYellowGreenByte(b []byte) []byte {
	return colConcatByte(BgYellow, b, BgDefault)
}

// BYellowf returns a BYellow formatted string.
func BYellowf(format string, a ...any) string {
	return BYellow(fmt.Sprintf(format, a...)).String()
}

// String implements the Stringer interface for type BMagenta.
func (c BMagenta) String() string {
	return colConcat(BgMagenta, string(c), BgDefault)
}

func BMagentaByte(b []byte) []byte {
	return colConcatByte(BgMagenta, b, BgDefault)
}

// BMagentaf returns a BMagenta formatted string.
func BMagentaf(format string, a ...any) string {
	return BMagenta(fmt.Sprintf(format, a...)).String()
}

// String implements the Stringer interface for type BWhite.
func (c BWhite) String() string {
	return colConcat(BgWhite, string(c), BgDefault)
}

func BWhiteByte(b []byte) []byte {
	return colConcatByte(BgWhite, b, BgDefault)
}

// BWhitef returns a BWhite formatted string.
func BWhitef(format string, a ...any) string {
	return BWhite(fmt.Sprintf(format, a...)).String()
}

// String implements the Stringer interface for type BBlack.
func (c BBlack) String() string {
	return colConcat(BgBlack, string(c), BgDefault)
}

func BBlackByte(b []byte) []byte {
	return colConcatByte(BgBlack, b, BgDefault)
}

// BBlackf returns a BBlack formatted string.
func BBlackf(format string, a ...any) string {
	return BBlack(fmt.Sprintf(format, a...)).String()
}

// String implements the Stringer interface for type BCyan.
func (c BCyan) String() string {
	return colConcat(BgCyan, string(c), BgDefault)
}

func BCyanByte(b []byte) []byte {
	return colConcatByte(BgCyan, b, BgDefault)
}

// BCyanf returns a BCyan formatted string.
func BCyanf(format string, a ...any) string {
	return BCyan(fmt.Sprintf(format, a...)).String()
}

// Modifier codes

// String implements the Stringer interface for type Blinking.
func (c Blinking) String() string {
	return colConcat(Blink, string(c), NoMode)
}

func BlinkingByte(b []byte) []byte {
	return colConcatByte(Blink, b, NoMode)
}

// String implements the Stringer interface for type Underline.
func (c Underline) String() string {
	return colConcat(Underln, string(c), NoMode)
}

func UnderlineByte(b []byte) []byte {
	return colConcatByte(Underln, b, NoMode)
}

// String implements the Stringer interface for type Bold.
func (c Bold) String() string {
	return colConcat(Bld, string(c), NoMode)
}

func BoldByte(b []byte) []byte {
	return colConcatByte(Bld, b, NoMode)
}

// String implements the Stringer interface for type Italic.
func (c Italic) String() string {
	return colConcat(Ital, string(c), NoMode)
}

func ItalicByte(b []byte) []byte {
	return colConcatByte(Ital, b, NoMode)
}

// NewColor gives a type Color back with specified fg/bg colors set that can
// be printed with anything using the Stringer iface.
func NewColor(str string, fg string, bg string) (Color, error) {
	if fg != "" {
		ifg, err := strconv.Atoi(fg)
		if err != nil {
			return Color(""), err
		}
		if ifg < 30 && ifg > 37 {
			return Color(""), errors.New("fg: " + fg + "not a valid color 30-37")
		}
	} else {
		fg = FgDefault
	}
	if bg != "" {
		ibg, err := strconv.Atoi(bg)
		if err != nil {
			return Color(""), err
		}
		if ibg < 40 && ibg > 47 {
			return Color(""), errors.New("Bg: " + bg + "not a valid color 40-47")
		}
	} else {
		bg = BgDefault
	}
	return Color(createColor(str, fg, bg)), nil
}

func createColor(str string, fg string, bg string) string {
	var buffer strings.Builder
	buffer.Grow(colorBaseLen + len(fg) + len(bg) + len(str))
	buffer.WriteString(CSI)
	buffer.WriteString(fg)
	buffer.WriteByte(sepCSI)
	buffer.WriteString(bg)
	buffer.WriteByte(endCSI)
	buffer.WriteString(str)
	buffer.WriteString(eraser)
	return buffer.String()
}

// String the stringer interface for all base color types.
func (c Color) String() string {
	if !colorEnable {
		clean := make([]byte, 0, len(c))
		src := []byte(c)
	L1:
		for i := 0; i < len(src); i++ {
			// Shortest possible mod.
			if len(src) < i+4 {
				clean = append(clean, src[i:]...)
				return string(clean)
			}
			if string(src[i:i+2]) == CSI {
				// Save current index incase this is not a term mod code.
				s := i
				// skip forward to end of mod
				for i += 2; i < len(src); i++ {
					switch src[i] {
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ';':
						// Legal characters in a term mod code.
						continue
					case 'm':
						// End of the term mod code.
						continue L1
					default:
						// Not a term mod code.
						i = s
					}
				}
			}
			clean = append(clean, src[i])
		}
		return string(clean)
	}
	return string(c)
}

// NewColor256 gives a type Color back using Term 256 color that can be printed with anything using the Stringer iface.
func NewColor256(str string, fg string, bg string) (Color, error) {
	if fg != "" {
		ifg, err := strconv.Atoi(fg)
		if err != nil {
			return Color(""), err
		}
		if ifg < 0 && ifg > 256 {
			return Color(""), errors.New("fg: " + fg + " not a valid color 0-256")
		}
	}
	if bg != "" {
		ibg, err := strconv.Atoi(bg)
		if err != nil {
			return Color(""), err
		}
		if ibg < 0 && ibg > 256 {
			return Color(""), errors.New("bg: " + bg + " not a valid color 0-256")
		}
	}
	var buffer strings.Builder
	buffer.Grow(color256BaseLen + len(bg) + len(fg) + len(str))
	buffer.WriteString(start256)
	buffer.WriteString(fg)
	buffer.WriteString(sep256)
	buffer.WriteString(bg)
	buffer.WriteByte(endCSI)
	buffer.WriteString(str)
	buffer.WriteString(eraser256)
	return Color(buffer.String()), nil
}

// NewColorRGB takes R G B and returns a ColorRGB type that can be printed by anything using the Stringer iface.
// Only Konsole to my knowledge that supports 24bit color
func NewColorRGB(str string, red uint8, green uint8, blue uint8) Color {
	ired := strconv.Itoa(int(red))
	igreen := strconv.Itoa(int(green))
	iblue := strconv.Itoa(int(blue))
	var buffer strings.Builder
	buffer.Grow(colorRGBBaseLen + len(ired) + len(igreen) + len(iblue) + len(str))
	buffer.WriteString(startRGB)
	buffer.WriteString(ired)
	buffer.WriteByte(sepCSI)
	buffer.WriteString(igreen)
	buffer.WriteByte(sepCSI)
	buffer.WriteString(iblue)
	buffer.WriteByte(endCSI)
	buffer.WriteString(str)
	buffer.WriteString(eraser256)
	return Color(buffer.String())
}

// String is a random color stringer.
func (c ColorRandom) String() string {
	if !colorEnable {
		return string(c)
	}
	ifg := rand.Int()%8 + 30
	ibg := rand.Int()%8 + 40
	return createColor(string(c), strconv.Itoa(ifg), strconv.Itoa(ibg))
}

// String gives a random fg color everytime it's printed.
func (c Random) String() string {
	if !colorEnable {
		return string(c)
	}
	ifg := int(rand.Int()%8 + 30)
	return colConcat(strconv.Itoa(ifg), string(c), FgDefault)
}

// String gives a random bg color everytime it's printed.
func (c BRandom) String() string {
	if !colorEnable {
		return string(c)
	}
	ibg := rand.Int()%8 + 40
	return colConcat(strconv.Itoa(ibg), string(c), BgDefault)
}

// NewCombo Takes a combination of modes and return a string with them all combined.
func NewCombo(s string, mods ...string) Color {
	var col, bcol, mod bool
	var buffer strings.Builder
	buffer.WriteString(CSI)
	first := true
	tracking := make(map[string]empty)
	for _, m := range mods {
		switch m {
		case FgBlack, FgRed, FgGreen, FgYellow, FgBlue, FgMagenta, FgCyan, FgWhite:
			if col {
				continue
			}
			col = true
		case BgBlack, BgRed, BgGreen, BgYellow, BgBlue, BgMagenta, BgCyan, BgWhite:
			if bcol {
				continue
			}
			bcol = true
		case Bld, Faint, Ital, Underln, Blink:
			if _, present := tracking[m]; present {
				continue
			}
			tracking[m] = empty{}
			mod = true
		default:
			continue
		}
		if first {
			first = false
		} else {
			buffer.WriteByte(sepCSI)
		}
		buffer.WriteString(m)
	}
	buffer.WriteByte(endCSI)
	buffer.WriteString(s)
	buffer.WriteString(CSI)
	if col {
		buffer.WriteString(FgDefault)
		if bcol {
			buffer.WriteByte(sepCSI)
			buffer.WriteString(BgDefault)
		}
	} else if bcol {
		buffer.WriteString(BgDefault)
	}
	if mod {
		if col || bcol {
			buffer.WriteByte(sepCSI)
		}
		buffer.WriteString(NoMode)
	}
	buffer.WriteByte(endCSI)
	return Color(buffer.String())
}

// TestTerm tries out most of the functions in this package and return
// a colourful string. Could be used to check what your terminal supports.
func TestTerm() string {
	res := "Standard 8:\n"
	res += "Fg:\t"
	for c := 30; c < 38; c++ {
		tres, _ := NewColor("#", strconv.Itoa(c), "")
		res += tres.String()
	}
	res += "\nBg:\t"
	for c := 40; c < 48; c++ {
		tres, _ := NewColor(" ", "", strconv.Itoa(c))
		res += tres.String()
	}
	res += "\nStandard 16:\t"
	for c := 0; c < 16; c++ {
		tcol, _ := NewColor256(" ", "", strconv.Itoa(c))
		res += tcol.String()
	}
	res += "\n"
	res += "256 Col:\n"
	// 6x6x6 cubes are trendy
	for row, base := 1, 0; row <= 6; row++ {
		base = (row * 6) + 9 // Step over the first 16 base colors
		for cubes := 1; cubes <= 6; cubes++ {
			for column := 1; column <= 6; column++ {
				tcol, _ := NewColor256(" ", "", strconv.Itoa(base+column))
				res += tcol.String()
			}
			base += 36 // 6 * 6
		}
		res += "\n"
	}
	// Grayscale left.
	res += "Grayscales:\n"
	for c := 232; c <= 255; c++ {
		tcol, _ := NewColor256(" ", "", strconv.Itoa(c))
		res += tcol.String()
	}
	return res
}

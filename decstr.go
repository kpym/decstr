// docstr is a package to detect and convert decimal strings.
package decstr

import (
	"strings"
)

// NoSeparator is a constant for no separator.
const NoSeparator rune = -1

type DecimalFormat struct {
	Point    rune // Decimal separator, if -1 then no decimal separator
	Group    rune // Grouping separator, if -1 then no grouping separator
	Standard bool // true if grouping is by 3, false if it's non-standard (e.g., 3 then by 2)
}

// String returns a string representation (like {`,`, ` `, standard}) for DecimalFormat.
func (df DecimalFormat) String() string {
	sep := func(r rune) string {
		if r == NoSeparator {
			return "<none>"
		}
		return string(r)
	}
	std := "non-standard"
	if df.Standard {
		std = "standard"
	}
	return "{`" + sep(df.Point) + "`, `" + sep(df.Group) + "`, " + std + "}"
}

// possibleGrouping is a map of possible grouping separators for each decimal separator.
var possibleGrouping = map[rune][]rune{
	',':  {' ', '.', '\''},
	'.':  {' ', ',', '\''},
	'·':  {','},
	'\'': {'.'},
}

// isPossible returns true if the grouping separator is possible for the decimal separator.
func isPossible(point, group rune) bool {
	groups, ok := possibleGrouping[point]
	if !ok {
		return false
	}
	for _, g := range groups {
		if g == group {
			return true
		}
	}
	return false
}

// getSign returns the sign and the absolute value of a decimal string.
// The sign is an empty string for positive numbers and a '-' for negative numbers.
func getSign(decimal string) (sign string, abs string) {
	abs = strings.TrimSpace(decimal)
	if len(abs) == 0 {
		return "", ""
	}
	switch abs[0] {
	case '-':
		return "-", strings.TrimLeft(abs[1:], " ")
	case '+':
		return "", strings.TrimLeft(abs[1:], " ")
	default:
		return "", abs
	}
}

// DetectFormat detects the decimal format of a string.
// It returns the DecimalFormat and a boolean indicating if the format was detected.
// ok is false because the string do not contain a decimal or because the format is ambiguous.
// If we can't detect if it is a standard or non-standard grouping, we assume it is standard.
func DetectFormat(decimal string) (df DecimalFormat, ok bool) {
	// temporary variables
	var (
		first        rune // first separator found
		point, group rune // decimal and grouping separators
		before       int  // number of digits before the separator
		mode         int  // 0: unknown, 2: non-standard grouping, 3: standard grouping
		hasDigit     bool // if we have at least one digit
	)
	_, decimal = getSign(decimal)
	// loop over the bytes of the string
	for i := 0; i < len(decimal); i++ {
		// is it a digit?
		if '0' <= decimal[i] && decimal[i] <= '9' {
			before++
			hasDigit = true
			continue
		}
		// is it the first non-digit character
		if first == 0 {
			switch decimal[i] {
			case ',', '.', '\'':
				first = rune(decimal[i])
				// is the rist separator a decimal separator necessarily?
				if before == 0 || before > 3 {
					point = first
				}
			case ' ':
				if before > 3 {
					return df, false
				}
				first = ' '
				group = ' '
			case 0xC2:
				if i+1 >= len(decimal) || decimal[i+1] != 0xB7 {
					return df, false
				}
				i++
				first = '·'
				point = '·'
			default:
				return df, false
			}
			before = 0
			continue
		}
		// are we after the decimal separator?
		if point != 0 {
			return df, false
		}
		// is it the next grouping separator?
		if first == rune(decimal[i]) {
			if (before != 2 && before != 3) || (mode > 0 && before != mode) {
				return df, false
			}
			group = first
			mode = before
			before = 0
			continue
		}
		// is it a midpoint?
		if decimal[i] == 0xC2 && i+1 < len(decimal) && decimal[i+1] == 0xB7 {
			i++
			point = '·'
		} else {
			point = rune(decimal[i])
		}
		group = first
		// is it the decimal separator?
		if before != 3 || !isPossible(point, group) {
			return df, false
		}
		before = 0
	}
	// if the string has no digit
	if !hasDigit {
		return df, false
	}
	// if no separator was found
	if first == 0 {
		df.Point = NoSeparator
		df.Group = NoSeparator
		df.Standard = true
		return df, true
	}
	// if we have a decimal separator
	if point != 0 {
		df.Point = point
		if group == 0 {
			df.Group = NoSeparator
		} else {
			df.Group = group
		}
		df.Standard = mode != 2
		return df, true
	}
	// if the only separator is a grouping separator
	if group != 0 {
		if before != 3 {
			return df, false
		}
		df.Point = NoSeparator
		df.Group = group
		df.Standard = mode != 2
		return df, true
	}
	// are we in the ambiguous case?
	if before == 3 {
		return df, false
	}
	// the only separator is necessarily a decimal separator
	df.Point = first
	df.Group = NoSeparator
	df.Standard = true
	return df, true
}

// Normalize returns the normalized decimal string.
// The normalized string is the decimal string:
//   - without leading or trailing spaces
//   - with a leading '-' if the number is negative
//   - with a decimal separator '.'
//   - without a grouping separator
//   - it is not starting with decimal separator : -.123 → -0.123
//   - it is not ending with decimal separator : 123. → 123.
//
// If the input string is not a valid decimal string,
// it returns the input string unchanged.
func Normalize(decimal string) (normalized string) {
	// get the sign and the absolute value
	sign, abs := getSign(decimal)
	// get the decimal format
	df, ok := DetectFormat(abs)
	if !ok {
		return decimal
	}
	sb := strings.Builder{}
	if sign == "-" {
		sb.WriteByte('-')
	}
	first := true
	hasDot := false
	for _, c := range abs {
		if c == df.Group {
			continue
		}
		if c == df.Point {
			if first {
				sb.WriteByte('0')
			}
			sb.WriteByte('.')
			hasDot = true
			continue
		}
		// skip leading '0' if any
		if first && c == '0' {
			continue
		}
		first = false
		sb.WriteRune(c)
	}
	normalized = sb.String()
	n := len(normalized)
	// trim trailing '0' from the decimal part
	if hasDot {
		for n > 0 && normalized[n-1] == '0' {
			n--
		}
	}
	// if the last character is the decimal separator remove it
	if normalized[n-1] == '.' {
		n--
	}
	return normalized[:n]
}

// IsNormalized returns true if the decimal string is normalized.
// A normalized decimal string is a string :
//   - may start with a '-'
//   - followed by a digit(s)
//   - followed by a '.' and a digit(s)
//   - can't start with '0' if the integer part is not 0
//   - can't have trailing zeros after the '.'
//   - can't have a trailing '.'
func IsNormalized(decimal string) bool {
	if len(decimal) == 0 {
		return false
	}
	if decimal == "0" {
		return true
	}
	var (
		first     bool
		after     bool
		c         rune
		expectDot bool
	)
	first = true
	for _, c = range decimal {
		// skip leading '-' if any
		if first && c == '-' {
			continue
		}
		if c == '.' {
			// can't start with '.' or have multiple '.'
			if first || after {
				return false
			}
			// we are after the '.'
			after = true
			expectDot = false
			continue
		}
		// if it is not a digit
		if c < '0' || c > '9' {
			return false
		}
		// if it is a digit but we expect a '.' (after first '0')
		if expectDot {
			return false
		}
		// if the integer part starts with '0'
		if first {
			expectDot = (c == '0')
		}
		first = false
	}
	// trailing '.' ?
	if c == '.' {
		return false
	}
	// trailing zeros after the '.'
	if c == '0' && after {
		return false
	}
	// '-0' case ?
	if expectDot {
		return false
	}
	return true
}

// Convert converts a decimal string to a formatted decimal string.
// If the input string is not a valid decimal string, it returns "0" and false.
// The input do not need to be normalized decimal string.
func (df DecimalFormat) Convert(decimal string) (new string, ok bool) {
	// try to normalize the decimal string
	if !IsNormalized(decimal) {
		decimal = Normalize(decimal)
		if !IsNormalized(decimal) {
			return "0", false
		}
	}
	var group int
	if df.Standard {
		group = 3
	} else {
		group = 2
	}
	sb := strings.Builder{}
	if decimal[0] == '-' {
		sb.WriteByte('-')
		decimal = decimal[1:]
	}
	parts := strings.Split(decimal, ".")
	n := len(parts[0])
	k, l := 0, (n-3)%group
	if l == 0 {
		l = group
	}
	for n > 3 {
		sb.WriteString(parts[0][k:l])
		sb.WriteRune(df.Group)
		k = l
		l += group
		n -= group
	}
	sb.WriteString(parts[0][k:])
	if len(parts) == 2 {
		sb.WriteRune(df.Point)
		sb.WriteString(parts[1])
	}
	return sb.String(), true
}

// decstr is a package for detecting and converting decimal strings.
// It provides utilities for identifying decimal formats and converting between them.
package decstr

import (
	"strings"
)

// NoSeparator represents the absence of a separator and is the 0 rune.
const NoSeparator = rune(0)

// DecimalFormat describes the format of a decimal string.
//   - Point: The decimal separator (or NoSeparator if absent).
//   - Group: The grouping separator (or NoSeparator if absent).
//   - Standard: True if grouping follows a standard pattern (e.g., groups of 3 digits),
//     False if it uses a non-standard pattern (e.g., 3 digits then 2 digits).
type DecimalFormat struct {
	Point    rune
	Group    rune
	Standard bool
}

// String returns a string representation of the DecimalFormat,
// formatted as {`<Point>`, `<Group>`, <standard|non-standard>}.
func (df DecimalFormat) String() string {
	// sep converts a rune to its string representation or "<none>" if NoSeparator.
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

// isPossible checks if the given grouping separator is valid for the specified decimal separator.
// Following https://en.wikipedia.org/wiki/Decimal_separator
// 1,234,567.89
// Australia,[51][52] Cambodia, Canada (English-speaking; unofficial), China,[53] Cyprus (currency numbers), Hong Kong, Iran, Ireland, Israel, Japan, Korea, Macau (in Chinese and English text), Malaysia, Mexico, Namibia, New Zealand, Pakistan, Peru (currency numbers), Philippines, Singapore, South Africa (English-speaking; unofficial), Taiwan, Thailand, United Kingdom and other Commonwealth states except Mozambique, United States.
// 1 234 567.89
// Canada (English-speaking; official), China,[53] Estonia (currency numbers), Hong Kong (in education), Mexico, Namibia, South Africa (English-speaking; unofficial), Sri Lanka, Switzerland (in federal texts for currency numbers only[54]), United Kingdom (in education), United States (in education)[citation needed]. SI style (English version) but SI doesn't include currency.
// 1_234_567.89
// Ada, C#, D, Eiffel, Fortran 90, Go, Haskell, Java, JavaScript, Julia, Kotlin , Perl, Python, Ruby, Rust, Swift... programming languages.
// 1 234 567,89
// Albania, Belgium (French), Brazil, Bulgaria, Canada (French-speaking), Costa Rica, Croatia, Czech Republic, Estonia, Finland,[55] France, Hungary, Italy (in education), Latin America, Latin Europe, Latvia, Lithuania, Macau (in Portuguese text), Mozambique, Norway, Peru, Poland, Portugal, Russia, Serbia (informal), Slovakia, Slovenia, South Africa (official[56]), Spain (official use since 2010, according to the RAE and CSIC), Sweden, Switzerland (in federal texts, except currency numbers[54]), Ukraine, Vietnam (in education). SI style (French version) but SI doesn't include currency.
// 1.234.567,89
// Austria, Belgium (Dutch), Bosnia and Herzegovina, Brazil (informal and in technology), Chile, Colombia, Croatia (in bookkeeping and technology),[57] Denmark, Germany, Greece, Indonesia, Italy, Latin America (informal), Netherlands, Romania, Slovenia, Serbia, Spain (used until 2010, inadvisable use according to the RAE and CSIC),[d][59] Turkey, Uruguay, Vietnam.
// 1,234,567·89
// Malaysia, Malta, Philippines (uncommon today), Singapore, Taiwan, United Kingdom (older, typically handwritten; in education)
// 12,34,567.89
// 12 34 567.89
// Bangladesh, India, Nepal, Pakistan (see Indian numbering system).
// 1'234'567.89
// Switzerland (computing), Liechtenstein.
// C++14, Rebol, and Red programming languages.
// 1'234'567,89
// Switzerland (handwriting), Italy (handwriting).
// 1.234.567'89
// Spain (handwriting, used until 1980s, inadvisable use according to the RAE and CSIC[citation needed]).
func isPossible(point, group rune) bool {
	switch point {
	case '.':
		return group == ' ' || group == '\u00A0' || group == ',' || group == '\'' || group == '_'
	case ',':
		return group == ' ' || group == '\u00A0' || group == '.' || group == '\''
	case '·': // Malaysia, Malta, Philippines (uncommon today), Singapore, Taiwan
		return group == ','
	case '\'':
		return group == '.'
	}
	return false
}

// bytestr is a type constraint for []byte and string, used for functions
// that operate generically on these types.
type bytestr interface {
	~[]byte | ~string
}

// trimLeft removes all leading occurrences of the specified character from the given byte slice or string.
func trimLeft[T bytestr](decimal T, c byte) T {
	var i int
	for i = 0; i < len(decimal); i++ {
		if decimal[i] != c {
			break
		}
	}
	return decimal[i:]
}

// trimRight removes all trailing occurrences of the specified character from the given byte slice or string.
func trimRight[T bytestr](decimal T, c byte) T {
	var i int
	for i = len(decimal) - 1; i >= 0; i-- {
		if decimal[i] != c {
			break
		}
	}
	return decimal[:i+1]
}

// trimSpace removes leading and trailing spaces from the given byte slice or string.
func trimSpace[T bytestr](decimal T) T {
	return trimRight(trimLeft(decimal, ' '), ' ')
}

// getSign extracts the sign and the absolute value of a decimal string.
// - decimal: The input decimal string or byte slice (may include leading/trailing spaces).
// - Returns:
//   - sign: An empty string for positive numbers, or a "-" for negative numbers.
//   - abs: The absolute value of the input (without the sign or leading spaces).
//
// If the input is empty or contains only spaces, both sign and abs are empty.
// Example:
//
//	getSign("-123") => "-", "123"
//	getSign("+123") => "", "123"
//	getSign("  123") => "", "123"
//	getSign("   ") => "", ""
func getSign[T bytestr](decimal T) (sign T, abs T) {
	abs = trimSpace(decimal)
	if len(abs) == 0 {
		return abs, abs
	}
	switch abs[0] {
	case '-': // Negative sign detected; trim it and return.
		return abs[:1], trimLeft(abs[1:], ' ')
	case '+': // Positive sign detected; trim it and return.
		return abs[:0], trimLeft(abs[1:], ' ')
	default: // No sign detected; return the absolute value.
		return abs[:0], abs
	}
}

// flushAtoB appends the contents of b to a and resets b to an empty slice.
func flushBtoA(a, b *[]byte) {
	if len(*b) > 0 {
		*a = append(*a, *b...)
		*b = (*b)[:0]
	}
}

// compose returns the normalized decimal string from the integer and decimal parts.
func compose(a, b []byte) []byte {
	a = trimLeft(a, '0')
	if len(a) == 0 {
		a = append(a, '0')
	}
	b = trimRight(b, '0')
	if len(b) == 0 {
		return a
	}
	a = append(a, '.')
	a = append(a, b...)
	return a
}

// detectAndNormalize detects the format of a decimal string and returns a normalized version of it.
// - decimal: The input decimal string or byte slice to process.
// - Returns:
//   - normalized: The normalized decimal string (with grouping separators removed and decimal part normalized).
//   - df: The detected decimal format (point, grouping, and whether grouping is standard or not).
//   - ok: A boolean indicating if the detection and normalization succeeded.
//
// The function supports various separators, such as ',', '.', '\”, and the midpoint '·'.
// Whitespace, non-standard grouping, and invalid formats are handled gracefully.
// Examples:
//
//	"1,234.56" -> "1234.56", {Point: '.', Group: ',', Standard: true}, true
//	"123.45"   -> "123.45", {Point: '.', Group: NoSeparator, Standard: true}, true
//	"123 45"   -> "", {}, false
//	""         -> "", {}, false
func detectAndNormalize[T bytestr](decimal T) (normalized T, df DecimalFormat, ok bool) {
	// temporary variables
	var (
		firstsep     rune // first separator found
		newsep       rune // new separator found
		point, group rune // decimal and grouping separators
		before       int  // number of digits before the separator
		mode         int  // 0: unknown, 2: non-standard grouping, 3: standard grouping
		hasDigit     bool // if we have at least one digit
	)
	a := make([]byte, 0, len(decimal)) // the integer part (before the decimal separator)
	b := make([]byte, 0, len(decimal)) // the decimal part (after the decimal separator)
	buf := &a                          // the current buffer (a or b)
	sign, abs := getSign(decimal)
	*buf = append(*buf, sign...)
	// loop over the bytes of the string
	for i := 0; i < len(abs); i++ {
		// handle digits
		if '0' <= abs[i] && abs[i] <= '9' {
			before++
			hasDigit = true
			*buf = append(*buf, abs[i])
			continue
		}

		// handle the first non-digit character
		if firstsep == 0 {
			// we never enter twice in this block
			switch abs[i] {
			case ',', '.', '\'':
				firstsep = rune(abs[i])
				// is the first separator a decimal separator?
				if before == 0 || before > 3 {
					point = firstsep
				}
				buf = &b // we start the possible decimal part (if not we will copy it back to a)
			case ' ', '_':
				if before > 3 {
					return decimal, df, false
				}
				firstsep, group = rune(abs[i]), rune(abs[i])
			case 0xC2:
				if i+1 >= len(abs) {
					// not a decimal number
					return decimal, df, false
				}
				switch abs[i+1] {
				case 0xB7: // center dot
					i++
					firstsep, point = '·', '·'
					buf = &b // we start the decimal part
				case 0xA0: // non-breaking space
					if before > 3 {
						// not a decimal number
						return decimal, df, false
					}
					i++
					firstsep, group = '\u00A0', '\u00A0'
				default:
					// not a decimal number
					return decimal, df, false
				}
			default:
				// not a decimal number
				return decimal, df, false
			}
			before = 0
			continue
		}

		// no more separator is allowed after the decimal separator
		if point != 0 {
			return decimal, df, false
		}

		if abs[i] == 0xC2 && i+1 < len(abs) {
			switch abs[i+1] {
			case 0xB7: // center dot
				newsep = '·'
			case 0xA0: // non-breaking space
				newsep = '\u00A0'
			default:
				// not a decimal number
				return decimal, df, false
			}
			i++
		} else {
			newsep = rune(abs[i])
		}

		// handle the grouping separator
		if firstsep == newsep {
			// grouping must match standard or non-standard rules (2 or 3 digits).
			if (before != 2 && before != 3) || (mode > 0 && before != mode) {
				return decimal, df, false
			}
			group, mode, before = firstsep, before, 0
			// if we were hesitating between a grouping and a decimal separator
			flushBtoA(&a, &b)
			buf = &a
			continue
		}
		// the new separator could be only a decimal separator
		// so the previous one is necessarily a grouping separator
		group = firstsep
		point = newsep
		// check if the decimal separator is valid
		if before != 3 || !isPossible(point, group) {
			return decimal, df, false
		}

		// handle ambiguity between grouping and decimal separator,
		// if we have collected some digits in the decimal part
		// transfer them to the integer part
		flushBtoA(&a, &b)
		// start collecting the decimal part
		buf = &b
		before = 0
	}

	// At this point df is zero, {NoSeparator, NoSeparator, false}.
	// We have to fill it with the detected values.

	// handle strings with no digits
	if !hasDigit {
		return decimal, df, false
	}

	// handle digits without any separator
	if firstsep == 0 {
		df.Standard = true
		return T(compose(a, b)), df, true
	}

	// handle digits with decimal separator
	if point != 0 {
		df.Point, df.Group, df.Standard = point, group, mode != 2
		return T(compose(a, b)), df, true
	}

	// handle digits only with grouping separator
	if group != 0 {
		if before != 3 {
			return decimal, df, false
		}
		df.Group, df.Standard = group, mode != 2
		return T(compose(a, b)), df, true
	}

	// handle digits with single unknown separator
	if before == 3 {
		// we are in the ambiguous case (3 digits before the separator)
		return decimal, df, false
	}
	// the only separator is necessarily a decimal separator
	df.Point, df.Standard = firstsep, true
	return T(compose(a, b)), df, true
}

// DetectFormat detects the decimal format of a string.
// It returns the detected DecimalFormat and a boolean indicating success.
// The boolean `ok` is false if the string does not contain a valid decimal format
// or if the format is ambiguous.
// If it is impossible to determine whether the grouping is standard or non-standard,
// it defaults to standard.
func DetectFormat[T bytestr](decimal T) (df DecimalFormat, ok bool) {
	_, df, ok = detectAndNormalize(decimal)
	return df, ok
}

// Normalize returns a normalized decimal string.
// A normalized decimal string adheres to the following rules:
//   - May start with a '-' (negative sign).
//   - Is followed by one or more digits.
//   - If a '.' is present, it is followed by one or more digits (e.g., "123." -> "123").
//   - Cannot start with '0' unless the integer part is exactly 0 (e.g., "0123.4" -> "123.4").
//   - Cannot have trailing zeros after the '.' (e.g., "123.000" -> "123").
//   - Cannot have a trailing '.' (e.g., "123." -> "123").
func Normalize[T bytestr](decimal T) (normalized T) {
	normalized, _, _ = detectAndNormalize(decimal)
	return normalized
}

// NormalizeCheck returns a normalized decimal string and a boolean.
// The boolean `ok` is true if the input string was successfully normalized;
// otherwise, it is false, indicating the input string is unchanged.
func NormalizeCheck[T bytestr](decimal T) (normalized T, ok bool) {
	normalized, _, ok = detectAndNormalize(decimal)
	return normalized, ok
}

// IsNormalized checks if a decimal string is normalized.
// A normalized decimal string adheres to the following rules:
//   - May start with a '-' (negative sign).
//   - Must be followed by one or more digits.
//   - If a '.' is present, it must be followed by one or more digits.
//   - Cannot start with '0' unless the integer part is exactly 0.
//   - Cannot have trailing zeros after the '.' (e.g., "123.000" -> false).
//   - Cannot have a trailing '.' (e.g., "123." -> false).
//   - The string cannot be empty.
func IsNormalized[T bytestr](decimal T) bool {
	if len(decimal) == 0 {
		return false
	}
	if len(decimal) == 1 && decimal[0] == '0' {
		return true
	}
	var (
		first     bool // whether we're processing the first character
		after     bool // whether we're after the '.'
		c         byte // current character
		expectDot bool // whether we expect a '.' after a leading '0'
	)
	first = true
	for i := 0; i < len(decimal); i++ {
		c = decimal[i]
		// skip leading '-' if any
		if first && c == '-' {
			continue
		}
		if c == '.' {
			// '.' cannot be the first character or appear multiple times.
			if first || after {
				return false
			}
			// we're now processing the decimal part (after the '.')
			after = true
			expectDot = false
			continue
		}
		// if we expect a '.' but encounter a digit, it's invalid
		if c < '0' || c > '9' {
			return false
		}
		// if we expect a '.' but encounter a digit, it's invalid
		if expectDot {
			return false
		}
		// check if the integer part starts with '0'
		if first {
			expectDot = (c == '0')
		}
		first = false
	}
	// ensure the last character is not '.' or '0' (if we're after '.')
	if c == '.' || (c == '0' && after) {
		return false
	}
	// special case for '-0'
	if expectDot {
		return false
	}
	return true
}

// Convert converts a decimal string to a formatted decimal string using the specified DecimalFormat.
// If the input string is not a valid decimal string, it returns "0" and false.
// The input string does not need to be a normalized decimal string.
// The output string is formatted based on the following rules:
//   - Grouping separators are inserted every 3 or 2 digits (depending on `df.Standard`).
//   - A custom decimal separator (`df.Point`) is used.
//   - Negative numbers retain their '-' sign. If + is present, it is removed.
func (df DecimalFormat) Convert(decimal string) (new string, ok bool) {
	// attempt to normalize the decimal string
	if !IsNormalized(decimal) {
		decimal = Normalize(decimal)
		// if normalization fails, return "0" and false
		if !IsNormalized(decimal) {
			return "0", false
		}
	}
	// determine the grouping size: 3 for standard formats, 2 for non-standard
	group := 3
	if !df.Standard {
		group = 2
	}

	// use a strings.Builder for efficient string construction
	sb := strings.Builder{}

	// handle negative numbers by writing the '-' sign and removing it from the input
	if decimal[0] == '-' {
		sb.WriteByte('-')
		decimal = decimal[1:]
	}

	// split the string into integer and fractional parts
	parts := strings.Split(decimal, ".")
	n := len(parts[0])

	// calculate initial grouping positions
	k, l := 0, (n-3)%group
	if l == 0 {
		l = group
	}

	// insert grouping separators for the integer part
	for n > 3 {
		sb.WriteString(parts[0][k:l])
		sb.WriteRune(df.Group)
		k = l
		l += group
		n -= group
	}
	sb.WriteString(parts[0][k:])

	// append the decimal separator and the fractional part if any
	if len(parts) == 2 {
		sb.WriteRune(df.Point)
		sb.WriteString(parts[1])
	}

	// return the formatted string and true, indicating success
	return sb.String(), true
}

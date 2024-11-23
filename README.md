# decstr

`decstr` is a Golang package to handle decimal strings. It provides utilities to normalize, check, detect formats, and convert decimal strings, without relying on external dependencies.

## Features

- Normalize decimal strings:
  - Removes grouping separators.
  - Standardizes the decimal separator to `.` (dot).
  - Removes leading zeros, trailing zeros, and trailing decimal separators for integers.
  
- Detect decimal format (grouping separator, decimal separator, and grouping style).
- Convert decimal strings to a specified format.

## Functions

All functions accept both `string` or `[]byte` as input and return results in the same type.
In what follows, only "string" is used for simplicity.

### `Normalize`
Normalizes a decimal string:
- No grouping separators.
- Decimal separator is `.` (dot).
- No leading or trailing zeros (and no trailing decimal for integers).
  
If the input string is not a valid decimal, it returns the string as-is.
### `NormalizeCheck`
Same as `Normalize`, but also returns a boolean indicating whether the string was normalized.

### `IsNormalized`
Checks if the decimal string is normalized.

### `DetectFormat`
Detects the decimal format:
- Returns the decimal separator.
- Returns the grouping separator (if any).
- Indicates whether the grouping is standard (3 digits per group) or non-standard (first 3 digits, then 2 per group).

### `Convert`
Converts a decimal string to the specified format.

## Example

```go
package main

import (
  "fmt"
  "github.com/kpym/decstr"
)

func main() {
    decimal := "1'234'567,89"

	// Normalize example
	normalized := decstr.Normalize(decimal)
	fmt.Println("Normalized:", normalized) // 1234567.89

	// Detect format example
	format, ok := decstr.DetectFormat(decimal)
	fmt.Println("Detected format:", format, "ok:", ok) // Detected format: {`,`, `'`, standard} ok: true
	// Convert example
	df := decstr.DecimalFormat{Point: '.', Group: ' ', Standard: false}
	converted, ok := df.Convert(decimal)
	fmt.Println("Converted:", converted, "ok:", ok) // Converted: 12 34 567.89 ok: true
}
```

## Documentation

The package documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/kpym/decstr).

## LICENSE

This package is released under the [MIT License](LICENSE).

## Decimal Format

Possible decimal writings from [Wikipedia](https://en.wikipedia.org/wiki/Decimal_separator):

- standard grouping (by 3)
  - `1,234,567.89` : Australia, Cambodia, Canada (English-speaking; unofficial), China, Cyprus (currency numbers), Hong Kong, Iran, Ireland, Israel, Japan, Korea, Macau (in Chinese and English text), Malaysia, Mexico, Namibia, New Zealand, Pakistan, Peru (currency numbers), Philippines, Singapore, South Africa (English-speaking; unofficial), Taiwan, Thailand, United Kingdom and other Commonwealth states except Mozambique, United States.
  - `1 234 567.89` : Canada (English-speaking; official), China, Estonia (currency numbers), Hong Kong (in education), Mexico, Namibia, South Africa (English-speaking; unofficial), Sri Lanka, Switzerland (in federal texts for currency numbers only), United Kingdom (in education), United States (in education). SI style (English version) but SI doesn't include currency.
  - `1 234 567,89` : Albania, Belgium (French), Brazil, Bulgaria, Canada (French-speaking), Costa Rica, Croatia, Czech Republic, Estonia, Finland, France, Hungary, Italy (in education), Latin America, Latin Europe, Latvia, Lithuania, Macau (in Portuguese text), Mozambique, Norway, Peru, Poland, Portugal, Russia, Serbia (informal), Slovakia, Slovenia, South Africa (official), Spain (official use since 2010, according to the RAE and CSIC), Sweden, Switzerland (in federal texts, except currency numbers[56]), Ukraine, Vietnam (in education). SI style (French version) but SI doesn't include currency.
  - `1.234.567,89` : Austria, Belgium (Dutch), Bosnia and Herzegovina, Brazil (informal and in technology), Chile, Colombia, Croatia (in bookkeeping and technology), Denmark, Germany, Greece, Indonesia, Italy, Latin America (informal), Netherlands, Romania, Slovenia, Serbia, Spain (used until 2010, inadvisable use according to the RAE and CSIC), Turkey, Uruguay, Vietnam.
  - `1,234,567·89` : Malaysia, Malta, Philippines (uncommon today), Singapore, Taiwan, United Kingdom (older, typically handwritten; in education)
  - `1'234'567.89` : Switzerland (computing), Liechtenstein.
  - `1'234'567,89` : Switzerland (handwriting), Italy (handwriting).
  - `1.234.567'89` : Spain (handwriting, used until 1980s).
- non standard grouping (3, then by 2)
  - `12,34,567.89` : Bangladesh, India, Nepal, Pakistan
  - `12 34 567.89` : Bangladesh, India, Nepal, Pakistan

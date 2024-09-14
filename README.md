# decstr

A golang package to work with decimal strings. No external dependencies.

## Functions

- `Normalize` get a decimal string and return a normalized string:
  - no grouping,
  - decimal separator is `.` (dot),
  - no leading zeros,
  - no trailing zeros and no trailing decimal separator (for integers).

  If the string is not a valid decimal string, it returns it as is.
- `IsNormalized` check if a decimal string is normalized.
- `DetectFormat` get a decimal string and return the decimal format:
  - decimal separator,
  - grouping separator,
  - standard grouping (by 3) or non standard grouping (3, then by 2).

  The only ambiguos decimal/integer is `##D<sep>DDD` (4.321 | 54.321 | 654.321).
- `Convert` get a decimal string and convert it to the specified format.

## Decimal Format

Possible decimal writings [Wikipedia](https://en.wikipedia.org/wiki/Decimal_separator).

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

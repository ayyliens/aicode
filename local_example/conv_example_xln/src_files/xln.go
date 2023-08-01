package xln

import "fmt"

var LANGS = []string{LANG_EN}

const (
	LANG_DEFAULT = LANG_EN
	LANG_EN      = `en`
	LANG_FR      = `fr`
)

var XLN_INVALID_LONGITUDE = Xln[func(lon float64) string]{
	LANG_EN: func(lon float64) string { return fmt.Sprintf(`invalid longitude %v`, lon) },
	LANG_FR: func(lon float64) string { return fmt.Sprintf(`longitude invalide %v`, lon) },
}.ValidGet()

var XLN_INVALID_LATITUDE = Xln[func(lat float64) string]{
	LANG_EN: func(lat float64) string { return fmt.Sprintf(`invalid latitude %v`, lat) },
	LANG_FR: func(lat float64) string { return fmt.Sprintf(`latitude invalide %v`, lat) },
}.ValidGet()

var XLN_MISSING_GEOGRAPHIC_COORDINATES = XlnS{
	LANG_EN: `missing geographic coordinates`,
	LANG_FR: `coordonnées géographiques manquantes`,
}.ValidGet()

var XLN_CURRENCY_NOT_EQUALS = Xln[func(exp, got string) string]{
	LANG_EN: func(exp, got string) string {
		return fmt.Sprintf(`expected currency %q, got %q`, exp, got)
	},
	LANG_FR: func(exp, got string) string {
		return fmt.Sprintf(`devise attendue %q, obtenue %q`, exp, got)
	},
}.ValidGet()

var XLN_MISSING_CURRENCY = XlnS{
	LANG_EN: `missing currency`,
	LANG_FR: `devise manquante`,
}.ValidGet()

var XLN_CURRENCY_INVALID = Xln[func(currency string) string]{
	LANG_EN: func(currency string) string {
		return fmt.Sprintf(`%v doesn't appear to be a valid currency`, currency)
	},
	LANG_FR: func(currency string) string {
		return fmt.Sprintf(`%v ne semble pas être une devise valide`, currency)
	},
}.ValidGet()

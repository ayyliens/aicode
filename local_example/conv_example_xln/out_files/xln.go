package xln

import "fmt"

var LANGS = []string{LANG_EN, LANG_ES, LANG_IT, LANG_PT}

const (
	LANG_DEFAULT = LANG_EN
	LANG_EN      = `en`
	LANG_FR      = `fr`
	LANG_ES      = `es`
	LANG_IT      = `it`
	LANG_PT      = `pt`
)

var XLN_INVALID_LONGITUDE = Xln[func(lon float64) string]{
	LANG_EN: func(lon float64) string { return fmt.Sprintf(`invalid longitude %v`, lon) },
	LANG_FR: func(lon float64) string { return fmt.Sprintf(`longitude invalide %v`, lon) },
	LANG_ES: func(lon float64) string { return fmt.Sprintf(`longitud inválida %v`, lon) },
	LANG_IT: func(lon float64) string { return fmt.Sprintf(`longitudine non valida %v`, lon) },
	LANG_PT: func(lon float64) string { return fmt.Sprintf(`longitude inválida %v`, lon) },
}.ValidGet()

var XLN_INVALID_LATITUDE = Xln[func(lat float64) string]{
	LANG_EN: func(lat float64) string { return fmt.Sprintf(`invalid latitude %v`, lat) },
	LANG_FR: func(lat float64) string { return fmt.Sprintf(`latitude invalide %v`, lat) },
	LANG_ES: func(lat float64) string { return fmt.Sprintf(`latitud inválida %v`, lat) },
	LANG_IT: func(lat float64) string { return fmt.Sprintf(`latitudine non valida %v`, lat) },
	LANG_PT: func(lat float64) string { return fmt.Sprintf(`latitude inválida %v`, lat) },
}.ValidGet()

var XLN_MISSING_GEOGRAPHIC_COORDINATES = XlnS{
	LANG_EN: `missing geographic coordinates`,
	LANG_FR: `coordonnées géographiques manquantes`,
	LANG_ES: `faltan coordenadas geográficas`,
	LANG_IT: `mancano coordinate geografiche`,
	LANG_PT: `faltam coordenadas geográficas`,
}.ValidGet()

var XLN_CURRENCY_NOT_EQUALS = Xln[func(exp, got string) string]{
	LANG_EN: func(exp, got string) string {
		return fmt.Sprintf(`expected currency %q, got %q`, exp, got)
	},
	LANG_FR: func(exp, got string) string {
		return fmt.Sprintf(`devise attendue %q, obtenue %q`, exp, got)
	},
	LANG_ES: func(exp, got string) string {
		return fmt.Sprintf(`se esperaba la moneda %q, pero se obtuvo %q`, exp, got)
	},
	LANG_IT: func(exp, got string) string {
		return fmt.Sprintf(`si aspettava la valuta %q, ma si è ottenuto %q`, exp, got)
	},
	LANG_PT: func(exp, got string) string {
		return fmt.Sprintf(`esperava-se a moeda %q, mas obteve-se %q`, exp, got)
	},
}.ValidGet()

var XLN_MISSING_CURRENCY = XlnS{
	LANG_EN: `missing currency`,
	LANG_FR: `devise manquante`,
	LANG_ES: `falta moneda`,
	LANG_IT: `mancanza di valuta`,
	LANG_PT: `falta de moeda`,
}.ValidGet()

var XLN_CURRENCY_INVALID = Xln[func(currency string) string]{
	LANG_EN: func(currency string) string {
		return fmt.Sprintf(`%v doesn't appear to be a valid currency`, currency)
	},
	LANG_FR: func(currency string) string {
		return fmt.Sprintf(`%v ne semble pas être une devise valide`, currency)
	},
	LANG_ES: func(currency string) string {
		return fmt.Sprintf(`%v no parece ser una moneda válida`, currency)
	},
	LANG_IT: func(currency string) string {
		return fmt.Sprintf(`%v non sembra essere una valuta valida`, currency)
	},
	LANG_PT: func(currency string) string {
		return fmt.Sprintf(`%v não parece ser uma moeda válida`, currency)
	},
}.ValidGet()

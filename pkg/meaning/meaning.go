package meaning

import (
	"fmt"
	"github.com/dark-enstein/crontable/pkg/reader"
	"golang.org/x/exp/slices"
	"io"
	"strconv"
	"strings"
)

const (
	Minute     = "minute"
	Hour       = "hour"
	DayOfMonth = "month"
	Month      = "month"
	DayOfWeek  = "day of the week"
)

const (
	TextWildCard = "every %v"
	TextEvery    = "every %v of the %v"
	TextNone     = "on the %v %v"
	TextComma    = "on the %v and %v %v"
	TextCommaPre = "on the "
	TextRange    = "between the %v and %v %v"
)

var (
	min        = "on the 5th and 9th minute"
	hour       = "between the 3rd and 19th hour"
	dayOfMonth = "every 5th of the month"
	month      = "on the 9th month"
	dayOfWeek  = "every day of the week"
)

type Writer interface {
}

type Written []byte

func Write(exit io.Writer, s []byte) (int, error) {
	return exit.Write(s)
}

func Explain(dec *reader.CronExpressionDecoded) []byte {
	// Get ordered keys so that I can range over the map in an orderly manner
	keys, mapDec := dec.FlattenToMap()
	chain := []string{}
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		v := *mapDec[k]
		switch k {
		case reader.Minute:
			suffix := Minute
			// resolve struct and ordinals
			chunkInterface := []interface{}{}
			lowOrd, highOrd, structure, _ := resolveAll(v)
			if v.DelimKind != reader.DelimWildcard {
				chunkInterface = append(chunkInterface, lowOrd)
			}
			chunk := ""
			if v.DelimKind != reader.DelimNone && v.DelimKind != reader.DelimEvery && v.DelimKind != reader.DelimWildcard {
				for i := 0; i < len(highOrd); i++ {
					chunkInterface = append(chunkInterface, highOrd[i])
				}
			}
			chunkInterface = append(chunkInterface, suffix)
			chunk += fmt.Sprintf(structure, chunkInterface...)
			chain = append(chain, chunk)
		case reader.Hour:
			suffix := Hour
			// resolve struct and ordinals
			chunkInterface := []interface{}{}
			lowOrd, highOrd, structure, _ := resolveAll(v)
			if v.DelimKind != reader.DelimWildcard {
				chunkInterface = append(chunkInterface, lowOrd)
			}
			chunk := ""
			if v.DelimKind != reader.DelimNone && v.DelimKind != reader.DelimEvery && v.DelimKind != reader.DelimWildcard {
				for i := 0; i < len(highOrd); i++ {
					chunkInterface = append(chunkInterface, highOrd[i])
				}
			}
			chunkInterface = append(chunkInterface, suffix)
			chunk += fmt.Sprintf(structure, chunkInterface...)
			chain = append(chain, chunk)
		case reader.DayOfTheMonth:
			suffix := DayOfMonth
			// resolve struct and ordinals
			chunkInterface := []interface{}{}
			lowOrd, highOrd, structure, _ := resolveAll(v)
			if v.DelimKind != reader.DelimWildcard {
				chunkInterface = append(chunkInterface, lowOrd)
			}
			chunk := ""
			if v.DelimKind != reader.DelimNone && v.DelimKind != reader.DelimEvery && v.DelimKind != reader.DelimWildcard {
				for i := 0; i < len(highOrd); i++ {
					chunkInterface = append(chunkInterface, highOrd[i])
				}
			}
			chunkInterface = append(chunkInterface, suffix)
			chunk += fmt.Sprintf(structure, chunkInterface...)
			chain = append(chain, chunk)
		case reader.Month:
			suffix := Month
			// resolve struct and ordinals
			chunkInterface := []interface{}{}
			lowOrd, highOrd, structure, _ := resolveAll(v)
			if v.DelimKind != reader.DelimWildcard {
				chunkInterface = append(chunkInterface, lowOrd)
			}
			chunk := ""
			if v.DelimKind != reader.DelimNone && v.DelimKind != reader.DelimEvery && v.DelimKind != reader.DelimWildcard {
				for i := 0; i < len(highOrd); i++ {
					chunkInterface = append(chunkInterface, highOrd[i])
				}
			}
			chunkInterface = append(chunkInterface, suffix)
			chunk += fmt.Sprintf(structure, chunkInterface...)
			chain = append(chain, chunk)
		case reader.DayOfTheWeek:
			suffix := DayOfWeek
			// resolve struct and ordinals
			chunkInterface := []interface{}{}
			lowOrd, highOrd, structure, _ := resolveAll(v)
			if v.DelimKind != reader.DelimWildcard {
				chunkInterface = append(chunkInterface, lowOrd)
			}
			chunk := ""
			if v.DelimKind != reader.DelimNone && v.DelimKind != reader.DelimEvery && v.DelimKind != reader.DelimWildcard {
				for i := 0; i < len(highOrd); i++ {
					chunkInterface = append(chunkInterface, highOrd[i])
				}
			}
			chunkInterface = append(chunkInterface, suffix)
			chunk += fmt.Sprintf(structure, chunkInterface...)
			chain = append(chain, chunk)
		}
	}
	title := titulate(strings.Join(chain, ", "))
	return []byte(title)
}

func titulate(s string) string {
	sRune := []rune(s)
	sRune[0] = []rune(strings.ToUpper(string(sRune[0])))[0]
	return string(sRune)
}

func NorminalToOrdinal(i int) string {
	iStr := []rune(strconv.Itoa(i))
	switch true {
	case slices.Equal(iStr, []rune("11")), slices.Equal(iStr, []rune("12")), slices.Equal(iStr, []rune("13")):
		return string(iStr) + "th"
	default:
		switch iStr[len(iStr)-1] {
		case '1':
			return string(iStr) + "st"
		case '2':
			return string(iStr) + "nd"
		case '3':
			return string(iStr) + "rd"
		default:
			return string(iStr) + "th"
		}
	}
}

func resolveAll(v reader.Catcher) (string, []string, string, error) {
	lowOrd, structure := "", ""
	highOrd := make([]string, len(v.High), len(v.High))
	switch v.DelimKind {
	case reader.DelimWildcard:
		structure = TextWildCard
	case reader.DelimEvery:
		structure = TextEvery
	case reader.DelimNone:
		structure = TextNone
	case reader.DelimComma:
		structure = commaAddMoreRef(TextComma, v.High)
	case reader.DelimRange:
		structure = TextRange
	default:
		_ = structure
	}

	// if delim os wildcard, every, none, range
	if len(v.High) == 1 {
		if v.Low == v.High[0] {
			ord := NorminalToOrdinal(v.Low)
			lowOrd, highOrd[0] = ord, ord
			return lowOrd, highOrd, structure, nil
		}
		lowOrd = NorminalToOrdinal(v.Low)
		highOrd[0] = NorminalToOrdinal(v.High[0])
		return lowOrd, highOrd, structure, nil
	}

	// when delim is comma
	lowOrd = NorminalToOrdinal(v.Low)
	for i := 0; i < len(v.High); i++ {
		highOrd = append(highOrd, NorminalToOrdinal(v.High[i]))
	}
	return lowOrd, highOrd, structure, nil
}

func commaAddMoreRef(s string, list []int) string {
	if len(list) <= 2 {
		return s
	}

	collab := TextCommaPre

	for i := 2; i < len(list); i++ {
		aft, found := strings.CutPrefix(s, TextCommaPre)
		_ = found
		collab += "%v, " + aft
	}
	return collab
}

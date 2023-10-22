package reader

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	DelimNone = iota
	DelimWildcard
	DelimComma
	DelimRange
	DelimEvery
)

const (
	Minute        = "minute"
	Hour          = "hour"
	Month         = "month"
	DayOfTheMonth = "dayOfTheMonth"
	DayOfTheWeek  = "dayOfTheWeek"
)

var (
	Bounds = []struct {
		low  int
		high int
	}{
		{
			0,
			60,
		},
		{
			0,
			60,
		},
		{
			1,
			31,
		},
		{
			1,
			12,
		},
		{
			1,
			7,
		},
	}
)

var (
	SampleCronFile = `
0 9 * * 6
`
)

type CronRead string

// CronExpression holds the shallow understanding to the cron expression passed in a simple struct.
type CronExpression struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
}

// CronExpressionDecoded holds a deep understanding of the cron expression passed in a deeper struct
type CronExpressionDecoded struct {
	Minute     Catcher
	Hour       Catcher
	DayOfMonth Catcher
	Month      Catcher
	DayOfWeek  Catcher
}

// FlattenToMap helps with easily accessing the values of CronExpressionDecoded in meaning.Explain by keys; helps minimize time complexity on retrieval
func (c *CronExpressionDecoded) FlattenToMap() ([]string, map[string]*Catcher) {
	mapper := map[string]*Catcher{}
	mapper[Minute] = &c.Minute
	mapper[Hour] = &c.Hour
	mapper[DayOfTheMonth] = &c.DayOfMonth
	mapper[Month] = &c.Month
	mapper[DayOfTheWeek] = &c.DayOfWeek
	// Ordered keys so that I can range over the map in an orderly manner
	orderedKeys := []string{Minute, Hour, DayOfTheMonth, Month, DayOfTheWeek}
	return orderedKeys, mapper
}

// Catcher holds a unit of deep cron expression knowledge. It represents the type of token passed in at a time, and the valid bounds for any token at that position.
type Catcher struct {
	Low       int
	High      []int
	DelimKind int
}

// OpenCrontableFile opens the crontab file passed in as argument, casting it into wrapper type CronRead before returning. It errors with os.File errors, and when the file is structurally invalid
func OpenCrontableFile(loc string) (*CronRead, error) {
	file, err := os.ReadFile(loc)
	if err != nil {
		log.Printf("crontab is invalid. \nsample cronfile: %v\n", SampleCronFile)
		return nil, err
	}

	cronTabExpr, _, found := bytes.Cut(file, []byte("\n"))
	if !found {
		bytesLength := len(bytes.Split(cronTabExpr, []byte(" ")))
		if bytesLength != 5 {
			log.Printf("crontab is invalid. the number of crontable arguments is %d \nsample cronfile: %v\n", strings.Count(string(cronTabExpr), " "), SampleCronFile)
			os.Exit(1)
			return nil, fmt.Errorf("crontab is invalid. the number of crontable arguments is %d \nsample cronfile: %s\n", bytesLength, SampleCronFile)
		}
	}
	read := CronRead(string(cronTabExpr))
	return &read, nil
}

// Validate validates a CronRead value. It checks that all the tokens are valid, and/or are within the bounds for their position
func (cr *CronRead) Validate() (bool, error) {
	str := cr.String()
	valErr := 0
	pieces := strings.Split(str, " ")
	for i := 0; i < len(pieces); i++ {
		if !validate(pieces[i]) {
			valErr++
		}
	}

	if valErr > 0 {
		return false, fmt.Errorf("encountered %d validation errors", valErr)
	}
	return true, nil
}

// Decode converts a CronRead into its CronExpressionDecoded, breaking its tokens into their separate units and preserving meaning.
func (cr *CronRead) Decode() *CronExpressionDecoded {
	str := cr.String()
	var catchAll []Catcher
	pieces := strings.Split(str, " ")
	for i := 0; i < len(pieces); i++ {
		var catch Catcher
		var err error
		catch.Low, catch.High, _, catch.DelimKind, err = validateWithFields(pieces[i], Bounds[i].low, Bounds[i].high)
		if err != nil {
			log.Println(fmt.Errorf("w"), err)
		}
		catchAll = append(catchAll, catch)
	}
	var catchAllDecoded CronExpressionDecoded
	catchAllDecoded.Minute = catchAll[0]
	catchAllDecoded.Hour = catchAll[1]
	catchAllDecoded.DayOfMonth = catchAll[2]
	catchAllDecoded.Month = catchAll[3]
	catchAllDecoded.DayOfWeek = catchAll[4]
	return &catchAllDecoded
}

// validate does a quick and shallow validation the string cron string passed into cron table. It returns a bool which signifies if the string is valid or not.
func validate(s string) bool {
	if s == "*" || canBeNumber(s) || containFavoredDelimiter(s) || startSlash(s) {
		return true
	}
	return false
}

// validateWithFields validates that the string argument passed in is a valid cron argument.
// It returns the low, high, boolean (true if string is valid); a false boolean, and error if string isn't
func validateWithFields(s string, low, high int) (int, []int, bool, int, error) {
	if s == "*" {
		return low, []int{high}, true, DelimWildcard, nil
	}
	if canBeNumber(s) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, []int{0}, false, DelimNone, nil
		}
		if i < low || i > high {
			return 0, []int{0}, false, DelimNone, fmt.Errorf("number %d not within acceptable bounds", i)
		}
		return i, []int{i}, true, DelimNone, nil
	}
	containDelim, delim, pieces := containFavoredDelimiterWhich(s)
	if containDelim {
		if len(pieces) < 2 {
			return 0, []int{-1}, false, delim, nil
		}
		var retInt []int
		for i := 0; i < len(pieces); i++ {
			if canBeNumber(pieces[i]) {
				i, err := strconv.Atoi(pieces[i])
				if err != nil {
					return 0, []int{0}, false, delim, nil
				}
				if i < low || i > high {
					return 0, []int{0}, false, delim, fmt.Errorf("number %d not within acceptable bounds", i)
				}
				retInt = append(retInt, i)
			}
		}
		return retInt[0], retInt[1:], true, delim, nil
	}

	isStartSlash, str := startSlashWhich(s)
	if isStartSlash {
		if canBeNumber(str) {
			i, err := strconv.Atoi(s)
			if err != nil {
				return 0, []int{0}, false, delim, nil
			}
			if i < low || i > high {
				return 0, []int{0}, false, delim, fmt.Errorf("number %d not within acceptable bounds", i)
			}
			return i, []int{i}, true, delim, nil
		}
		return 0, []int{0}, false, delim, fmt.Errorf("invald string cannot be number in %v", str)
	}
	return 0, []int{0}, false, delim, fmt.Errorf("invald string passed in %v", s)
}

func startSlash(s string) bool {
	_, aft, found := strings.Cut("/", s)
	if !found {

		if !canBeNumber(aft) {
			return false
		}
	}
	log.Println("does not contain slash")
	return false
}

func startSlashWhich(s string) (bool, string) {
	_, aft, found := strings.Cut("/", s)
	if !found {

		if !canBeNumber(aft) {
			return false, aft
		}
	}
	log.Println("does not contain slash")
	return false, ""
}

func canBeNumber(s string) bool {
	_, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return true
}

func containFavoredDelimiter(s string) bool {
	if strings.Contains(s, "-") {
		log.Println("contains delimiter -")
		byDash := strings.Split(s, "-")
		log.Printf("string split into %v", byDash)
		if len(byDash) != 2 {
			log.Printf("len of string %v is not 2", byDash)
			return false
		} else if len(byDash) == 2 {
			for i := 0; i < len(byDash); i++ {
				if !canBeNumber(byDash[i]) {
					return false
				}
			}
			log.Printf("len of string %v is 2", byDash)
			return true
		}
	}

	if strings.Contains(s, ",") {
		log.Println("contains delimiter ,")
		byComma := strings.Split(s, ",")
		log.Printf("string split into %v", byComma)
		if len(byComma) != 2 {
			log.Printf("len of string %v is not 2", byComma)
			return false
		} else if len(byComma) == 2 {
			log.Printf("len of string %v is 2", byComma)
			return true
		}
	}

	log.Println("does not contain favored delimiter")
	return false
}

func containFavoredDelimiterWhich(s string) (bool, int, []string) {
	if strings.Contains(s, "-") {
		log.Println("contains delimiter -")
		byDash := strings.Split(s, "-")
		log.Printf("string split into %v", byDash)
		if len(byDash) != 2 {
			log.Printf("len of string %v is not 2", byDash)
			return false, DelimRange, byDash
		} else if len(byDash) == 2 {
			for i := 0; i < len(byDash); i++ {
				if !canBeNumber(byDash[i]) {
					return false, DelimRange, byDash
				}
			}
			log.Printf("len of string %v is 2", byDash)
			return true, DelimRange, byDash
		}
	}

	if strings.Contains(s, ",") {
		log.Println("contains delimiter ,")
		byComma := strings.Split(s, ",")
		log.Printf("string split into %v", byComma)
		if len(byComma) < 2 {
			log.Printf("len of string %v is not 2", byComma)
			return false, DelimComma, byComma
		} else if len(byComma) >= 2 {
			log.Printf("len of string %v is 2", byComma)
			return true, DelimComma, byComma
		}
	}

	log.Println("does not contain favored delimiter")
	return false, 0, []string{}
}

func (cr *CronRead) MarshalIntoCronExpression() (*CronExpression, error) {
	str := cr.String()
	pieces := strings.Split(str, " ")
	if len(pieces) < 5 {
		return nil, fmt.Errorf("formatted cron expression has less than 5 arguments: %v", cr)
	}
	return &CronExpression{
		Minute:     pieces[0],
		Hour:       pieces[1],
		DayOfMonth: pieces[2],
		Month:      pieces[3],
		DayOfWeek:  pieces[4],
	}, nil
}

func (cr *CronRead) String() string {
	rattle := fmt.Sprintf("%v", *cr)
	return rattle
}

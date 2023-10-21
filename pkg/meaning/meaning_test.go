package meaning

import (
	"context"
	"fmt"
	"github.com/dark-enstein/crontable/pkg/reader"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
)

type written struct {
	input    []reader.CronExpressionDecoded
	expected []string
}

type CronTab struct {
	suite.Suite
	ctx       context.Context
	log       *log.Logger
	testState *written
	norminal  *map[int]string
}

// NorminalInts holds a map of nominal integers and their expected ordinal format
var NorminalInts = map[int]string{
	13: "13th", 7: "7th", 22: "22nd", 8: "8th", 31: "31st", 11: "11th", 28: "28th", 2: "2nd", 18: "18th", 25: "25th", 6: "6th", 30: "30th", 20: "20th", 14: "14th", 29: "29th", 10: "10th", 1: "1st", 5: "5th", 16: "16th", 26: "26th", 23: "23rd", 4: "4th", 3: "3rd", 19: "19th", 15: "15th", 12: "12th", 21: "21st", 27: "27th", 24: "24th", 9: "9th",
}

var CronDecodedTestInputs = written{
	[]reader.CronExpressionDecoded{
		{
			Minute: reader.Catcher{
				Low:       5,
				High:      []int{9},
				DelimKind: reader.DelimComma,
			},
			Hour: reader.Catcher{
				Low:       3,
				High:      []int{19},
				DelimKind: reader.DelimRange,
			},
			DayOfMonth: reader.Catcher{
				Low:       5,
				High:      []int{5},
				DelimKind: reader.DelimEvery,
			},
			Month: reader.Catcher{
				Low:       5,
				High:      []int{5},
				DelimKind: reader.DelimNone,
			},
			DayOfWeek: reader.Catcher{
				Low:       0,
				High:      []int{0},
				DelimKind: reader.DelimWildcard,
			},
		},
		{
			Minute: reader.Catcher{
				Low:       5,
				High:      []int{9},
				DelimKind: reader.DelimComma,
			},
			Hour: reader.Catcher{
				Low:       3,
				High:      []int{19},
				DelimKind: reader.DelimRange,
			},
			DayOfMonth: reader.Catcher{
				Low:       5,
				High:      []int{5},
				DelimKind: reader.DelimEvery,
			},
			Month: reader.Catcher{
				Low:       5,
				High:      []int{5},
				DelimKind: reader.DelimNone,
			},
			DayOfWeek: reader.Catcher{
				Low:       0,
				High:      []int{0},
				DelimKind: reader.DelimWildcard,
			},
		},
	},
	[]string{
		"On the 5th and 9th minute, between the 3rd and 19th hour, every 5th of the Month, on the 9th Month, every day of the week",
		"On the 5th and 9th minute, between the 3rd and 19th hour, every 5th of the Month, on the 9th Month, every day of the week",
	},
}

func (c *CronTab) SetupTest() {
	c.log = log.New(os.Stdout, "crontable: ", log.LstdFlags)
	log.Println("Starting tests...")
	c.ctx = context.Background()
	c.testState = &CronDecodedTestInputs
	c.norminal = &NorminalInts

	c.log.Println("Tests startup complete...")
}

// TestCronDecoded tests the Write function
func (c *CronTab) TestCronDecoded() {
	log := c.log
	for i := 0; i < len(c.testState.input); i++ {
		log.Printf("Decoding test input %v", c.testState.input[i])
		resStr := Write(&c.testState.input[i])
		c.Assert().Equal(c.testState.expected[i], string(resStr), fmt.Sprintf("expected %v, got %v. input: %v", c.testState.expected[i], resStr, c.testState.input[i]))
	}
}

// TestNorminalToOrdinal tests the NorminalToOrdinal function
func (c *CronTab) TestNorminalToOrdinal() {
	log := c.log
	for k, v := range *c.norminal {
		log.Printf("Testing NorminalToOrdinal conversion for integer: %d", k)
		ord := NorminalToOrdinal(k)
		c.Assert().Equal(v, ord, fmt.Sprintf("expected %v, got %v. input: %v", v, ord, k))
	}
}

func (c *CronTab) TearDownSuite() {
	log := c.log
	log.Println("Commencing test cleanup")
	//err := cleanUpAfterCatTest()
	//s.Require().NoError(err)
	log.Println("All testing complete")
}

func TestUtilTest(t *testing.T) {
	suite.Run(t, new(CronTab))
}

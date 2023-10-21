package reader

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
)

type state struct {
	input    []string
	expected []bool
}

type CronTab struct {
	suite.Suite
	ctx       context.Context
	log       *log.Logger
	testState *map[string]state
}

var CrontableTestInputs = map[string]state{
	"minute": state{
		[]string{
			"0 * * * 6",   // missing
			"0 9 * * SAT", // within bounds
			"0 9 * * SUN", // over bounds
			"3 9 * * 7",   // missing
			"21 9 * * 7",  // within bounds
			"89 9 * * 7",  // over bounds
			"0 9 * * 7",   // missing
			"0 9 * * 7",   // within bounds
			"0 9 * * 7",   // over bounds
		},
		[]bool{
			false, // missing
			true,  // within bounds
			false, // over bounds
			false, // missing
			true,  // within bounds
			false, // over bounds
			false, // missing
			true,  // within bounds
			false, // over bounds
		},
	},
	"hour": state{
		[]string{
			"0 * * * 6",   // missing
			"0 9 * * SAT", // within bounds
			"0 9 * * SUN", // over bounds
			"3 9 * * 7",   // missing
			"21 9 * * 7",  // within bounds
			"89 9 * * 7",  // over bounds
			"0 9 * * 7",   // missing
			"0 9 * * 7",   // within bounds
			"0 9 * * 7",   // over bounds
		},
		[]bool{
			false, // missing
			true,  // within bounds
			false, // over bounds
			false, // missing
			true,  // within bounds
			false, // over bounds
			false, // missing
			true,  // within bounds
			false, // over bounds
		},
	},
	"day_of_month": state{
		[]string{
			"0 * * * 6",   // missing
			"0 9 * * SAT", // within bounds
			"0 9 * * SUN", // over bounds
			"3 9 * * 7",   // missing
			"21 9 * * 7",  // within bounds
			"89 9 * * 7",  // over bounds
			"0 9 * * 7",   // missing
			"0 9 * * 7",   // within bounds
			"0 9 * * 7",   // over bounds
		},
		[]bool{
			false, // missing
			true,  // within bounds
			false, // over bounds
			false, // missing
			true,  // within bounds
			false, // over bounds
			false, // missing
			true,  // within bounds
			false, // over bounds
		},
	},
	"month": state{
		[]string{
			"0 * * * 6",   // missing
			"0 9 * * SAT", // within bounds
			"0 9 * * SUN", // over bounds
			"3 9 * * 7",   // missing
			"21 9 * * 7",  // within bounds
			"89 9 * * 7",  // over bounds
			"0 9 * * 7",   // missing
			"0 9 * * 7",   // within bounds
			"0 9 * * 7",   // over bounds
		},
		[]bool{
			false, // missing
			true,  // within bounds
			false, // over bounds
			false, // missing
			true,  // within bounds
			false, // over bounds
			false, // missing
			true,  // within bounds
			false, // over bounds
		},
	},
	"day_of_week": state{
		[]string{
			"0 * * * 6",   // missing
			"0 9 * * SAT", // within bounds
			"0 9 * * SUN", // over bounds
			"3 9 * * 7",   // missing
			"21 9 * * 7",  // within bounds
			"89 9 * * 7",  // over bounds
			"0 9 * * 7",   // missing
			"0 9 * * 7",   // within bounds
			"0 9 * * 7",   // over bounds
		},
		[]bool{
			false, // missing
			true,  // within bounds
			false, // over bounds
			false, // missing
			true,  // within bounds
			false, // over bounds
			false, // missing
			true,  // within bounds
			false, // over bounds
		},
	},
}

func (c *CronTab) SetupTest() {
	c.log = log.New(os.Stdout, "crontable: ", log.LstdFlags)
	log.Println("Starting tests...")
	c.ctx = context.Background()
	c.testState = &CrontableTestInputs

	c.log.Println("Tests startup complete...")
}

// TestValidation tests that the validations returned by the crontab reader is correct
func (c *CronTab) TestValidation() {
	log := c.log
	for k, v := range *c.testState {
		log.Println("Validating section %v", k)
		for i := 0; i < len(v.input); i++ {
			actual, _ := ValidateExpression(v.input[i])
			log.Println("result from validation:", actual)
			c.Assert().Equal(v.expected[i], actual, fmt.Sprintf("expected %v, got %v. input: %v", v.expected[i], actual, v.input[i]))
		}
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

package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	errorPlanScript = MustGetCurrentFileUseTest("../", "plan_failed.txt")
)

func Test_lineToPlanOutput(t *testing.T) {

	for _, output := range successOutput {

		add, change, destroy := linesToPlanOutput(output)

		assert.Equal(t, add, "0")
		assert.Equal(t, change, "0")
		assert.Equal(t, destroy, "1")
	}
}

func Test_linseToParseErrorOutput(t *testing.T) {
	status, result := LinseToParseLastMesasge(string(errorPlanScript))

	fmt.Println("status : ", status)
	fmt.Println("result : ", result)

}

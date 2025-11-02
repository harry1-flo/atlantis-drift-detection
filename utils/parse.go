package utils

import (
	"regexp"
	"strings"
)

func linesToPlanOutput(tfOutputs string) (add string, change string, destroy string) {

	lines := strings.Split(tfOutputs, "Plan:")
	re := regexp.MustCompile(`(\d+)\s+to add,\s+(\d+)\s+to change,\s+(\d+)\s+to destroy`)
	matches := re.FindStringSubmatch(strings.TrimSpace(lines[1]))

	if len(matches) == 4 {
		return "", "", ""
	}

	add = matches[1]
	change = matches[2]
	destroy = matches[3]

	return add, change, destroy
}

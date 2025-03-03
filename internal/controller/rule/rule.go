package rule

import (
	"strings"
)

func parseRuleParams(ruleName string) []string {
	index := strings.Index(ruleName, ":")
	if index == -1 {
		return []string{}
	}
	keyValue := ruleName[index+1:]
	return strings.Split(keyValue, ",")
}

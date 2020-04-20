package types

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	reDnmString = `[a-z][a-z0-9]{2,}`
	reAmt       = `[[:digit:]]*\.[[:digit:]]+`
	reSpc       = `[[:space:]]*`
	reCoin      = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmt, reSpc, reDnmString))
)

// HasDuplicate checks if the given array contains duplicate elements.
// Return true if there exists duplicate elements, false otherwise
func HasDuplicate(arr []string) bool {
	elementMap := make(map[string]bool)

	for _, elem := range arr {
		if _, ok := elementMap[elem]; ok {
			return true
		}

		elementMap[elem] = true
	}

	return false
}

// ParseCoinParts parses the given string to the amount and denom
func ParseCoinParts(coinStr string) (denom, amount string, err error) {
	coinStr = strings.ToLower(strings.TrimSpace(coinStr))

	matches := reCoin.FindStringSubmatch(coinStr)
	if matches == nil {
		err = fmt.Errorf("invalid coin string: %s", coinStr)
		return
	}

	denom, amount = matches[3], matches[1]
	return
}


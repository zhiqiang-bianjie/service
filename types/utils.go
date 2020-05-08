package types

import (
	"fmt"
	"regexp"
	"strings"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

var (
	reDnmString = `[a-z][a-z0-9]{2,}`
	reAmt       = `[0-9]+(\.[0-9]+)?`
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

// HexBytes wrappers the tendermint HexBytes
type HexBytes tmbytes.HexBytes

func (bz HexBytes) String() string {
	return tmbytes.HexBytes(bz).String()
}

// MarshalYAML returns the YAML representation
func (bz HexBytes) MarshalYAML() (interface{}, error) {
	return bz.String(), nil
}

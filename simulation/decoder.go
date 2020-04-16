package simulation

import (
	"bytes"
	"fmt"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/irismod/service/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding service type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.ServiceDefinitionKey):
		var definition1, definition2 types.ServiceDefinition
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &definition1)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &definition2)
		return fmt.Sprintf("%v\n%v", definition1, definition2)

	default:
		panic(fmt.Sprintf("invalid service key prefix %X", kvA.Key[:1]))
	}
}

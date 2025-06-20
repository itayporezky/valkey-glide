// Copyright Valkey GLIDE Project Contributors - SPDX Identifier: Apache-2.0

package options

import (
	"github.com/itayporezky/valkey-glide/go/v4/constants"
	"github.com/itayporezky/valkey-glide/go/v4/internal/utils"
)

// Aggregate represents the method of aggregating scores from multiple sets
type Aggregate string

const (
	AggregateSum Aggregate = "SUM" // Aggregates by summing the scores of each element across sets
	AggregateMin Aggregate = "MIN" // Aggregates by taking the minimum score of each element across sets
	AggregateMax Aggregate = "MAX" // Aggregates by taking the maximum score of each element across sets
)

// converts the Aggregate to its Valkey API representation
func (a Aggregate) ToArgs() ([]string, error) {
	return []string{constants.AggregateKeyWord, string(a)}, nil
}

// This is a basic interface. Please use one of the following implementations:
// - KeyArray
// - WeightedKeys
type KeysOrWeightedKeys interface {
	ToArgs() ([]string, error)
}

// represents a list of keys of the sorted sets involved in the aggregation operation
type KeyArray struct {
	Keys []string
}

// converts the KeyArray to its Valkey API representation
func (k KeyArray) ToArgs() ([]string, error) {
	args := []string{utils.IntToString(int64(len(k.Keys)))}
	args = append(args, k.Keys...)
	return args, nil
}

type KeyWeightPair struct {
	Key    string
	Weight float64
}

// represents the mapping of sorted set keys to their score weights
type WeightedKeys struct {
	KeyWeightPairs []KeyWeightPair
}

// converts the WeightedKeys to its Valkey API representation
func (w WeightedKeys) ToArgs() ([]string, error) {
	keys := make([]string, 0, len(w.KeyWeightPairs))
	weights := make([]string, 0, len(w.KeyWeightPairs))
	args := make([]string, 0)
	for _, pair := range w.KeyWeightPairs {
		keys = append(keys, pair.Key)
		weights = append(weights, utils.FloatToString(pair.Weight))
	}
	args = append(args, utils.IntToString(int64(len(keys))))
	args = append(args, keys...)
	args = append(args, constants.WeightsKeyword)
	args = append(args, weights...)
	return args, nil
}

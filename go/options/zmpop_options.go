// Copyright Valkey GLIDE Project Contributors - SPDX Identifier: Apache-2.0

package options

import (
	"github.com/itayporezky/valkey-glide/go/v4/constants"
	"github.com/itayporezky/valkey-glide/go/v4/internal/utils"
)

// Optional arguments for `ZMPop` and `BZMPop` in [SortedSetCommands]
type ZMPopOptions struct {
	count      int64
	countIsSet bool
}

func NewZMPopOptions() *ZMPopOptions {
	return &ZMPopOptions{}
}

// Set the count.
func (zmpo *ZMPopOptions) SetCount(count int64) *ZMPopOptions {
	zmpo.count = count
	zmpo.countIsSet = true
	return zmpo
}

func (zmpo *ZMPopOptions) ToArgs() ([]string, error) {
	var args []string

	if zmpo.countIsSet {
		args = append(args, constants.CountKeyword, utils.IntToString(zmpo.count))
	}

	return args, nil
}

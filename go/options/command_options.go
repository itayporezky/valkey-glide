// Copyright Valkey GLIDE Project Contributors - SPDX Identifier: Apache-2.0

package options

import (
	"strconv"

	"github.com/itayporezky/valkey-glide/go/v4/constants"

	"github.com/itayporezky/valkey-glide/go/v4/internal/errors"
	"github.com/itayporezky/valkey-glide/go/v4/internal/utils"
)

// SetOptions represents optional arguments for the [api.StringCommands.SetWithOptions] command.
//
// See [valkey.io]
//
// [valkey.io]: https://valkey.io/commands/set/
type SetOptions struct {
	// If ConditionalSet is not set the value will be set regardless of prior value existence. If value isn't set because of
	// the condition, [api.StringCommands.SetWithOptions] will return a zero-value string ("").
	ConditionalSet constants.ConditionalSet
	// Value to compare when [SetOptions.ConditionalSet] is set to `OnlyIfEquals`.
	ComparisonValue string
	// Set command to return the old value stored at the given key, or a zero-value string ("") if the key did not exist. An
	// error is returned and [api.StringCommands.SetWithOptions] is aborted if the value stored at key is not a string.
	// Equivalent to GET in the valkey API.
	ReturnOldValue bool
	// If not set, no expiry time will be set for the value.
	// Supported ExpiryTypes ("EX", "PX", "EXAT", "PXAT", "KEEPTTL")
	Expiry *Expiry
}

func NewSetOptions() *SetOptions {
	return &SetOptions{}
}

// Sets the condition to [SetOptions.ConditionalSet] for setting the value.
//
// This method overrides any previously set [SetOptions.ConditionalSet] and [SetOptions.ComparisonValue].
//
// Deprecated: Use [SetOptions.SetOnlyIfExists], [SetOptions.SetOnlyIfDoesNotExist], or [SetOptions.SetOnlyIfEquals] instead.
func (setOptions *SetOptions) SetConditionalSet(conditionalSet constants.ConditionalSet) *SetOptions {
	setOptions.ConditionalSet = conditionalSet
	setOptions.ComparisonValue = ""
	return setOptions
}

// Sets the condition to [SetOptions.OnlyIfExists] for setting the value. The key
// will be set if it already exists.
//
// This method overrides any previously set [SetOptions.ConditionalSet] and [SetOptions.ComparisonValue].
func (setOptions *SetOptions) SetOnlyIfExists() *SetOptions {
	setOptions.ConditionalSet = constants.OnlyIfExists
	setOptions.ComparisonValue = ""
	return setOptions
}

// Sets the condition to [SetOptions.OnlyIfDoesNotExist] for setting the value. The key
// will not be set if it already exists.
//
// This method overrides any previously set [SetOptions.ConditionalSet] and [SetOptions.ComparisonValue].
func (setOptions *SetOptions) SetOnlyIfDoesNotExist() *SetOptions {
	setOptions.ConditionalSet = constants.OnlyIfDoesNotExist
	setOptions.ComparisonValue = ""
	return setOptions
}

// Sets the condition to [SetOptions.OnlyIfEquals] for setting the value. The key
// will be set if the provided comparison value matches the existing value.
//
// This method overrides any previously set [SetOptions.ConditionalSet] and [SetOptions.ComparisonValue].
//
// since Valkey 8.1 and above.
func (setOptions *SetOptions) SetOnlyIfEquals(comparisonValue string) *SetOptions {
	setOptions.ConditionalSet = constants.OnlyIfEquals
	setOptions.ComparisonValue = comparisonValue
	return setOptions
}

func (setOptions *SetOptions) SetReturnOldValue(returnOldValue bool) *SetOptions {
	setOptions.ReturnOldValue = returnOldValue
	return setOptions
}

func (setOptions *SetOptions) SetExpiry(expiry *Expiry) *SetOptions {
	setOptions.Expiry = expiry
	return setOptions
}

func (opts *SetOptions) ToArgs() ([]string, error) {
	args := []string{}
	var err error
	if opts.ConditionalSet != "" {
		args = append(args, string(opts.ConditionalSet))
		if opts.ConditionalSet == constants.OnlyIfEquals {
			args = append(args, opts.ComparisonValue)
		}
	}

	if opts.ReturnOldValue {
		args = append(args, constants.ReturnOldValue)
	}

	if opts.Expiry != nil {
		switch opts.Expiry.Type {
		case constants.Seconds, constants.Milliseconds, constants.UnixSeconds, constants.UnixMilliseconds:
			args = append(args, string(opts.Expiry.Type), strconv.FormatUint(opts.Expiry.Count, 10))
		case constants.KeepExisting:
			args = append(args, string(opts.Expiry.Type))
		default:
			err = &errors.RequestError{Msg: "Invalid expiry type"}
		}
	}

	return args, err
}

// GetExOptions represents optional arguments for the [api.StringCommands.GetExWithOptions] command.
//
// See [valkey.io]
//
// [valkey.io]: https://valkey.io/commands/getex/
type GetExOptions struct {
	// If not set, no expiry time will be set for the value.
	// Supported ExpiryTypes ("EX", "PX", "EXAT", "PXAT", "PERSIST")
	Expiry *Expiry
}

func NewGetExOptions() *GetExOptions {
	return &GetExOptions{}
}

func (getExOptions *GetExOptions) SetExpiry(expiry *Expiry) *GetExOptions {
	getExOptions.Expiry = expiry
	return getExOptions
}

func (opts *GetExOptions) ToArgs() ([]string, error) {
	args := []string{}
	var err error

	if opts.Expiry != nil {
		switch opts.Expiry.Type {
		case constants.Seconds, constants.Milliseconds, constants.UnixSeconds, constants.UnixMilliseconds:
			args = append(args, string(opts.Expiry.Type), strconv.FormatUint(opts.Expiry.Count, 10))
		case constants.Persist:
			args = append(args, string(opts.Expiry.Type))
		default:
			err = &errors.RequestError{Msg: "Invalid expiry type"}
		}
	}

	return args, err
}

// Expiry is used to configure the lifetime of a value.
type Expiry struct {
	Type  constants.ExpiryType
	Count uint64
}

func NewExpiry() *Expiry {
	return &Expiry{}
}

func (ex *Expiry) SetType(expiryType constants.ExpiryType) *Expiry {
	ex.Type = expiryType
	return ex
}

func (ex *Expiry) SetCount(count uint64) *Expiry {
	ex.Count = count
	return ex
}

// LPosOptions represents optional arguments for the [api.ListCommands.LPosWithOptions] and
// [api.ListCommands.LPosCountWithOptions] commands.
//
// See [valkey.io]
//
// [valkey.io]: https://valkey.io/commands/lpos/
type LPosOptions struct {
	// Represents if the rank option is set.
	IsRankSet bool
	// The rank of the match to return.
	Rank int64
	// Represents if the max length parameter is set.
	IsMaxLenSet bool
	// The maximum number of comparisons to make between the element and the items in the list.
	MaxLen int64
}

func NewLPosOptions() *LPosOptions {
	return &LPosOptions{}
}

func (lposOptions *LPosOptions) SetRank(rank int64) *LPosOptions {
	lposOptions.IsRankSet = true
	lposOptions.Rank = rank
	return lposOptions
}

func (lposOptions *LPosOptions) SetMaxLen(maxLen int64) *LPosOptions {
	lposOptions.IsMaxLenSet = true
	lposOptions.MaxLen = maxLen
	return lposOptions
}

func (opts *LPosOptions) ToArgs() ([]string, error) {
	args := []string{}
	if opts.IsRankSet {
		args = append(args, constants.RankKeyword, utils.IntToString(opts.Rank))
	}

	if opts.IsMaxLenSet {
		args = append(args, constants.MaxLenKeyword, utils.IntToString(opts.MaxLen))
	}

	return args, nil
}

// Optional arguments to Restore(key string, ttl int64, value string, option RestoreOptions)
//
// Note IDLETIME and FREQ modifiers cannot be set at the same time.
//
// [valkey.io]: https://valkey.io/commands/restore/
type RestoreOptions struct {
	// Subcommand string to replace existing key.
	replace string
	// Subcommand string to represent absolute timestamp (in milliseconds) for TTL.
	absTTL string
	// It represents the idletime/frequency of object.
	eviction Eviction
}

func NewRestoreOptions() *RestoreOptions {
	return &RestoreOptions{}
}

// Custom setter methods to replace existing key.
func (restoreOption *RestoreOptions) SetReplace() *RestoreOptions {
	restoreOption.replace = constants.ReplaceKeyword
	return restoreOption
}

// Custom setter methods to represent absolute timestamp (in milliseconds) for TTL.
func (restoreOption *RestoreOptions) SetABSTTL() *RestoreOptions {
	restoreOption.absTTL = constants.ABSTTLKeyword
	return restoreOption
}

// For eviction purpose, you may use IDLETIME or FREQ modifiers.
type Eviction struct {
	// It represent IDLETIME or FREQ.
	Type constants.EvictionType
	// It represents count(int) of the idletime/frequency of object.
	Count int64
}

// Custom setter methods set the idletime/frequency of object.
func (restoreOption *RestoreOptions) SetEviction(evictionType constants.EvictionType, count int64) *RestoreOptions {
	restoreOption.eviction.Type = evictionType
	restoreOption.eviction.Count = count
	return restoreOption
}

func (opts *RestoreOptions) ToArgs() ([]string, error) {
	args := []string{}
	var err error
	if opts.replace != "" {
		args = append(args, string(opts.replace))
	}
	if opts.absTTL != "" {
		args = append(args, string(opts.absTTL))
	}
	if (opts.eviction != Eviction{}) {
		args = append(args, string(opts.eviction.Type), utils.IntToString(opts.eviction.Count))
	}
	return args, err
}

// Optional arguments for `Info` for standalone client
type InfoOptions struct {
	// A list of [Section] values specifying which sections of information to retrieve.
	// When no parameter is provided, [Section.Default] is assumed.
	// Starting with server version 7.0.0 `INFO` command supports multiple sections.
	Sections []constants.Section
}

// Optional arguments for `Info` for cluster client
type ClusterInfoOptions struct {
	*InfoOptions
	*RouteOption
}

func (opts *InfoOptions) ToArgs() ([]string, error) {
	if opts == nil {
		return []string{}, nil
	}
	args := make([]string, 0, len(opts.Sections))
	for _, section := range opts.Sections {
		args = append(args, string(section))
	}
	return args, nil
}

// Optional arguments to Copy(source string, destination string, option CopyOptions)
//
// [valkey.io]: https://valkey.io/commands/Copy/
type CopyOptions struct {
	// The REPLACE option removes the destination key before copying the value to it.
	replace bool
	// Option allows specifying an alternative logical database index for the destination key
	dbDestination int64
}

func NewCopyOptions() *CopyOptions {
	return &CopyOptions{replace: false}
}

// Custom setter methods to removes the destination key before copying the value to it.
func (restoreOption *CopyOptions) SetReplace() *CopyOptions {
	restoreOption.replace = true
	return restoreOption
}

// Custom setter methods to allows specifying an alternative logical database index for the destination key.
func (copyOption *CopyOptions) SetDBDestination(destinationDB int64) *CopyOptions {
	copyOption.dbDestination = destinationDB
	return copyOption
}

func (opts *CopyOptions) ToArgs() ([]string, error) {
	args := []string{}
	var err error
	if opts.replace {
		args = append(args, string(constants.ReplaceKeyword))
	}
	if opts.dbDestination >= 0 {
		args = append(args, "DB", utils.IntToString(opts.dbDestination))
	}
	return args, err
}

// Optional arguments for `ZPopMin` and `ZPopMax` commands.
type ZPopOptions struct {
	count int64
}

func NewZPopOptions() *ZPopOptions {
	return &ZPopOptions{}
}

// The maximum number of popped elements. If not specified, pops one member.
func (opts *ZPopOptions) SetCount(count int64) *ZPopOptions {
	opts.count = count
	return opts
}

// `ZPopMax/Min` don't use the COUNT keyword, only ZMPop will use .
func (opts *ZPopOptions) ToArgs(withKeyword bool) ([]string, error) {
	if opts.count <= 0 {
		return []string{}, nil
	}
	if withKeyword {
		return []string{"COUNT", strconv.FormatInt(opts.count, 10)}, nil
	}
	return []string{strconv.FormatInt(opts.count, 10)}, nil
}

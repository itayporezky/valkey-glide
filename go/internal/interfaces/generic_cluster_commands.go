// Copyright Valkey GLIDE Project Contributors - SPDX Identifier: Apache-2.0

package interfaces

import (
	"context"

	"github.com/itayporezky/valkey-glide/go/v4/config"
	"github.com/itayporezky/valkey-glide/go/v4/models"
	"github.com/itayporezky/valkey-glide/go/v4/options"
)

// GenericClusterCommands supports commands for the "Generic Commands" group for cluster client.
//
// See [valkey.io] for details.
//
// [valkey.io]: https://valkey.io/commands/#generic
type GenericClusterCommands interface {
	CustomCommand(ctx context.Context, args []string) (models.ClusterValue[any], error)

	CustomCommandWithRoute(ctx context.Context, args []string, route config.Route) (models.ClusterValue[any], error)

	Scan(ctx context.Context, cursor options.ClusterScanCursor) (options.ClusterScanCursor, []string, error)

	ScanWithOptions(
		ctx context.Context,
		cursor options.ClusterScanCursor,
		opts options.ClusterScanOptions,
	) (options.ClusterScanCursor, []string, error)

	RandomKey(ctx context.Context) (models.Result[string], error)

	RandomKeyWithRoute(ctx context.Context, opts options.RouteOption) (models.Result[string], error)
}

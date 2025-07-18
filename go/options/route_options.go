// Copyright Valkey GLIDE Project Contributors - SPDX Identifier: Apache-2.0

package options

import "github.com/itayporezky/valkey-glide/go/v4/config"

// An extension to command option types with Routes
type RouteOption struct {
	// Specifies the routing configuration for the command.
	// The client will route the command to the nodes defined by `route`.
	Route config.Route
}

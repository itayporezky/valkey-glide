// Copyright Valkey GLIDE Project Contributors - SPDX Identifier: Apache-2.0

package options

import "github.com/itayporezky/valkey-glide/go/v4/config"

type DBSizeOptions struct {
	Route config.Route
}

func NewTimeOptionsBuilder() *DBSizeOptions {
	return &DBSizeOptions{}
}

func (dbSizeOptions *DBSizeOptions) SetRoute(route config.Route) *DBSizeOptions {
	dbSizeOptions.Route = route
	return dbSizeOptions
}

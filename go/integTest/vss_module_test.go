// Copyright Valkey GLIDE Project Contributors - SPDX Identifier: Apache-2.0

package integTest

import (
	"context"
	"strings"

	"github.com/itayporezky/valkey-glide/go/v4/constants"

	"github.com/itayporezky/valkey-glide/go/v4/options"
	"github.com/stretchr/testify/assert"
)

func (suite *GlideTestSuite) TestModuleVerifyVssLoaded() {
	client := suite.defaultClusterClient()
	result, err := client.InfoWithOptions(context.Background(),
		options.ClusterInfoOptions{
			InfoOptions: &options.InfoOptions{Sections: []constants.Section{constants.Server}},
			RouteOption: nil,
		},
	)

	assert.Nil(suite.T(), err)
	for _, value := range result.MultiValue() {
		assert.True(suite.T(), strings.Contains(value, "# search_index_stats"))
	}
}

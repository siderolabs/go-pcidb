// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pcidb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/siderolabs/go-pcidb/pkg/pcidb"
)

func TestLookup(t *testing.T) {
	t.Parallel()

	v, ok := pcidb.LookupClass(4)
	assert.True(t, ok)
	assert.Equal(t, "Multimedia controller", v)

	v, ok = pcidb.LookupSubclass(4, 3)
	assert.True(t, ok)
	assert.Equal(t, "Audio device", v)

	v, ok = pcidb.LookupVendor(0x106b)
	assert.True(t, ok)
	assert.Equal(t, "Apple Inc.", v)

	v, ok = pcidb.LookupProduct(0x1002, 0x5960)
	assert.True(t, ok)
	assert.Equal(t, "RV280 [Radeon 9200 PRO / 9250]", v)
}

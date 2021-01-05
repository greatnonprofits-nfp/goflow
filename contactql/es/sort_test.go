package es_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/greatnonprofits-nfp/goflow/contactql/es"
	"github.com/greatnonprofits-nfp/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestElasticSort(t *testing.T) {
	resolver := newMockResolver()

	type testCase struct {
		Description string          `json:"description"`
		SortBy      string          `json:"sort_by"`
		Elastic     json.RawMessage `json:"elastic,omitempty"`
		Error       string          `json:"error,omitempty"`
	}
	tcs := make([]testCase, 0, 20)
	tcJSON, err := ioutil.ReadFile("testdata/to_sort.json")
	require.NoError(t, err)

	err = json.Unmarshal(tcJSON, &tcs)
	require.NoError(t, err)

	for _, tc := range tcs {
		sort, err := es.ToElasticFieldSort(tc.SortBy, resolver)

		if tc.Error != "" {
			assert.EqualError(t, err, tc.Error)
		} else {
			src, _ := sort.Source()
			encoded, _ := jsonx.Marshal(src)
			test.AssertEqualJSON(t, []byte(tc.Elastic), encoded, "field sort mismatch for %s", tc.Description)
		}
	}
}

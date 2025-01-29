package filter_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/planetlabs/go-ogc/filter"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cachedSchema *jsonschema.Schema

func getSchema(t *testing.T) *jsonschema.Schema {
	if cachedSchema != nil {
		return cachedSchema
	}
	schemaData, err := os.ReadFile("testdata/schema/cql2.json")
	require.NoError(t, err)

	schema, err := jsonschema.CompileString("cql2.json", string(schemaData))
	require.NoError(t, err)

	cachedSchema = schema
	return schema
}

type FilterCase struct {
	filter *filter.Filter
	data   string
}

func assertFilterIO(t *testing.T, c *FilterCase) {
	schema := getSchema(t)

	data, err := json.Marshal(c.filter)
	require.NoError(t, err)
	assert.JSONEq(t, c.data, string(data))

	var v any
	require.NoError(t, json.Unmarshal(data, &v))
	if err := schema.Validate(v); err != nil {
		t.Errorf("failed to validate\n%#v", err)
	}

	filter := &filter.Filter{}
	require.NoError(t, json.Unmarshal([]byte(c.data), filter))
	assert.Equal(t, c.filter, filter)

	assert.JSONEq(t, c.data, c.filter.String())
}

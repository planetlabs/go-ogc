package filter_test

import (
	"os"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/require"
)

func getSchema(t *testing.T) *jsonschema.Schema {
	schemaData, err := os.ReadFile("testdata/schema/cql2.json")
	require.NoError(t, err)

	schema, err := jsonschema.CompileString("cql2.json", string(schemaData))
	require.NoError(t, err)

	return schema
}

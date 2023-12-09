package gatewayapi

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/api7/apisix-operator/proto/adminapi"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/yaml"
)

func TestTranslateHTTPRoute(t *testing.T) {
	inputFiles, err := filepath.Glob(filepath.Join("testdata", "httproute*.in.yaml"))
	require.NoError(t, err)

	opts := cmpopts.IgnoreUnexported(adminapi.Route{}, adminapi.Upstream{})
	for _, inputFile := range inputFiles {
		input, err := os.ReadFile(inputFile)
		require.NoError(t, err)

		in := new(gwapiv1.HTTPRoute)
		err = yaml.UnmarshalStrict(input, in)
		require.NoError(t, err)

		outputFilePath := strings.ReplaceAll(inputFile, ".in.yaml", ".out.json")
		output, err := os.ReadFile(outputFilePath)
		require.NoError(t, err)

		want := []*adminapi.Route{}
		err = json.Unmarshal(output, &want)
		require.NoError(t, err)

		got, err := translateHTTPRoute(in)
		require.NoError(t, err)

		require.Empty(t, cmp.Diff(want, got, opts))
	}
}

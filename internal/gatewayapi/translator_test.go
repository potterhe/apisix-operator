package gatewayapi

import (
	"encoding/json"
	"testing"

	"github.com/api7/apisix-operator/proto/adminapi"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/yaml"
)

func TestTranslate(t *testing.T) {
	crd := `
apiVersion: gateway.networking.k8s.io/v1beta1
kind: HTTPRoute
metadata:
  name: example-route
spec:
  parentRefs:
  - name: example-gateway
  hostnames:
  - "example.com"
  rules:
  - backendRefs:
    - name: example-svc
      port: 80
`

	wantStr := `
[{
	"name": "example-route",
	"uris": ["/*"],
	"hosts": ["example.com"]
}]
`

	in := new(gwapiv1.HTTPRoute)
	err := yaml.UnmarshalStrict([]byte(crd), in)
	if err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	got, err := translate(in)
	if err != nil {
		t.Errorf("Translate failed: %v", err)
	}

	want := new([]*adminapi.Route)
	err = json.Unmarshal([]byte(wantStr), want)
	if err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	opts := cmpopts.IgnoreUnexported(adminapi.Route{})
	diff := cmp.Diff(*want, got, opts)
	if diff != "" {
		t.Errorf("Translate failed: %s", diff)
	}

}

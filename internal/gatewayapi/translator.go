package gatewayapi

import (
	"github.com/api7/apisix-operator/proto/adminapi"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func translate(in *gwapiv1.HTTPRoute) (*adminapi.Route, error) {
	r := new(adminapi.Route)
	r.Name = in.Name

	return r, nil
}

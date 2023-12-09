package gatewayapi

import (
	"errors"
	"fmt"

	"github.com/api7/apisix-operator/proto/adminapi"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func translate(in *gwapiv1.HTTPRoute) ([]*adminapi.Route, error) {
	var routes []*adminapi.Route

	var hosts []string
	for _, hostname := range in.Spec.Hostnames {
		hosts = append(hosts, string(hostname))
	}

	for _, rule := range in.Spec.Rules {
		route, err := translateRule(&rule)
		if err != nil {
			return nil, err
		}

		route.Name = in.Name
		route.Hosts = hosts

		routes = append(routes, route)
	}

	return routes, nil
}

func translateRule(rule *gwapiv1.HTTPRouteRule) (*adminapi.Route, error) {
	route := new(adminapi.Route)

	matches := rule.Matches
	if len(matches) == 0 {
		defaultType := gwapiv1.PathMatchPathPrefix
		defaultValue := "/"
		matches = []gwapiv1.HTTPRouteMatch{
			{
				Path: &gwapiv1.HTTPPathMatch{
					Type:  &defaultType,
					Value: &defaultValue,
				},
			},
		}
	}

	for _, match := range matches {
		if match.Path != nil {
			switch *match.Path.Type {
			case gwapiv1.PathMatchExact:
				route.Uris = []string{*match.Path.Value}
			case gwapiv1.PathMatchPathPrefix:
				route.Uris = []string{*match.Path.Value + "*"}
			case gwapiv1.PathMatchRegularExpression:
				return nil, errors.New("unimplement path match type " + string(*match.Path.Type))
			default:
				return nil, errors.New("unknown path match type " + string(*match.Path.Type))
			}
		}
	}

	if len(rule.BackendRefs) > 0 {
		upstream := new(adminapi.Upstream)
		upstream.Nodes = make(map[string]int32)

		for _, ref := range rule.BackendRefs {
			nodeItemKey := fmt.Sprintf("%s:%d", ref.Name, int32(*ref.Port))
			upstream.Nodes[nodeItemKey] = 1
		}

		route.Upstream = upstream
	}

	return route, nil
}

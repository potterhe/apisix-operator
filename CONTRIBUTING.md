## Gateway API

- [Installing Gateway API](https://gateway-api.sigs.k8s.io/guides/#installing-gateway-api)

## operator-sdk

operator-sdk: v1.32.0

```sh
operator-sdk init --domain apisix.apache.org --repo github.com/api7/apisix-operator --plugins=go/v4-alpha
```

### HTTPRouteFilter

```sh
operator-sdk create api --group gateway --version v1alpha1 --kind PluginConfig --resource --controller
```

## protobuf

```sh
protoc --go_out=. --go_opt=paths=source_relative proto/adminapi/*.proto
```
package docs

import _ "embed"

// OpenAPI contains the API contract served by /openapi.yaml.
//
//go:embed openapi.yaml
var OpenAPI []byte

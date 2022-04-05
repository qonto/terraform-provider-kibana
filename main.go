package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	provider "github.com/qonto/terraform-provider-kibana/internal/provider"
	mykibana "github.com/qonto/terraform-provider-kibana/pkg/kibana_api"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website
//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	kibanaClient := mykibana.CreateNewKibanaClient()
	opts := &plugin.ServeOpts{ProviderFunc: provider.New(kibanaClient)}
	plugin.Serve(opts)
}

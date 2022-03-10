package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	provider "github.com/qonto/terraform-provider-kibana/internal/provider"
	mykibana "github.com/qonto/terraform-provider-kibana/pkg/kibana_api"
)

func main() {
	kibanaClient := mykibana.CreateNewKibanaClient()
	opts := &plugin.ServeOpts{ProviderFunc: provider.New(kibanaClient)}
	plugin.Serve(opts)
}

package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	myprovider "github.com/qonto/terraform-provider-kibana/internal/provider"
	mykibana "github.com/qonto/terraform-provider-kibana/pkg/kibana_api"
)

var provider *schema.Provider
var providers map[string]*schema.Provider

func init() {
	k := mykibana.KibanaMockClient{}
	provider = myprovider.New(&k)()
	providers = map[string]*schema.Provider{
		"kibana": provider,
	}
}
func TestProvider(t *testing.T) {
	k := mykibana.KibanaMockClient{}
	if err := myprovider.New(&k)().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	k := mykibana.KibanaMockClient{}
	var _ *schema.Provider = myprovider.New(&k)()
}

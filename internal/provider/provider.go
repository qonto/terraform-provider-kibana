package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mykibana "github.com/qonto/terraform-provider-kibana/pkg/kibana_api"
)

func New(kibanaApi mykibana.KibanaAPI) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"kibana_host": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("KIBANA_HOST", nil),
				},
				"kibana_auth": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("KIBANA_AUTH", nil),
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"kibana_alert_rule": resourceAlertRule(),
			},
		}

		p.ConfigureContextFunc = configure(p, kibanaApi)

		return p
	}
}

func configure(p *schema.Provider, kibanaApi mykibana.KibanaAPI) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		kibana_host := d.Get("kibana_host").(string)
		kibana_auth := d.Get("kibana_auth").(string)
		kibanaApi.SetupClient(kibana_host, kibana_auth)
		return kibanaApi, nil
	}
}

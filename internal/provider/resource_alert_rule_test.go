package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestKibanaAlertRule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: getAlertConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAlertExists("kibana_alert_rule.test"),
				),
			},
		},
	})
}

func getAlertConfig() string {
	return fmt.Sprintf(`
	resource "kibana_alert_rule" "test" {
    consumer     = "siem"
    enabled      = false
    name         = "Test alert"
    notify_when  = "onActiveAlert"
    params       = jsonencode(
        {
            author           = []
            description      = "VPN activity on the same user from different countries"
            exceptionsList   = []
            falsePositives   = []
            filters          = []
            from             = "now-21660s"
            immutable        = false
            index            = [
                "infra-docker-pritunl-*",
            ]
            language         = "kuery"
            license          = ""
            maxSignals       = 100
            meta             = {
                from                = "1m"
                kibana_siem_app_url = "https://dd4b9a16f3264df3a5e87f27be632cc8.eu-central-1.aws.cloud.es.io:9243/app/security"
            }
            outputIndex      = ".siem-signals-default"
            query            = "event.user_name:* and geoip.country_iso_code :fr"
            references       = []
            riskScore        = 47
            riskScoreMapping = []
            ruleId           = "ba443266-0b29-498a-a50e-f0f2f27aa700"
            severity         = "medium"
            severityMapping  = []
            threat           = []
            threshold        = {
                cardinality = [
                    {
                        field = "geoip.country_iso_code.keyword"
                        value = 3
                    },
                ]
                field       = [
                    "event.user_name.keyword",
                ]
                value       = 1
            }
            to               = "now"
            type             = "threshold"
            version          = 3
        }
    )
    rule_type_id = "siem.signals"
    schedule     = {
        "interval" = "5h"
    }
    tags         = ["ok"]
    actions       {
        group  = "default"
        id     = "407ed770-9cf4-47aa-8840-0b5cdb22496e"
        params = jsonencode(
            {
                message = <<-EOT
                    Rule {{context.rule.name}} generated {{state.signals_count}} alerts
                    <{{{context.results_link}}}|Follow on this dashboard>
                EOT
            }
        )
			}
    }
	`)
}

func testCheckAlertExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AlertID set")
		}

		return nil
	}
}

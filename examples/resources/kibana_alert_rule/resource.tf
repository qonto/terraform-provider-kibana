resource "kibana_alert_rule" "example" {
  consumer    = "siem"
  enabled     = false
  name        = "As code - new"
  notify_when = "onActiveAlert"
  params = jsonencode(
    {
      author         = []
      description    = "Some alert rule description"
      exceptionsList = []
      falsePositives = []
      filters        = []
      from           = "now-21660s"
      immutable      = false
      index = [
        "someindex-*",
      ]
      language   = "kuery"
      license    = ""
      maxSignals = 100
      meta = {
        from                = "1m"
        kibana_siem_app_url = "https:/yourkibanainstanceurl"
      }
      outputIndex      = ".siem-signals-default"
      query            = "event.statuscode != 200 and somecondition != true"
      references       = []
      riskScore        = 47
      riskScoreMapping = []
      ruleId           = "theruleid" // Unrelated to the rule object id. https://github.com/elastic/kibana/issues/100667
      severity         = "medium"
      severityMapping  = []
      threat           = []
      threshold = {
        cardinality = [
          {
            field = "somefield.keyword"
            value = 3
          },
        ]
        field = [
          "event.user_name.keyword",
        ]
        value = 1
      }
      to      = "now"
      type    = "threshold"
      version = 3
    }
  )
  rule_type_id = "siem.signals"
  schedule = {
    "interval" = "5h"
  }
  tags = ["ok"]
  actions {
    group = "default"
    id    = "407ed770-9cf4-47aa-8840-0b5cdb22496e" // The Id must refer to an existing Action.
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
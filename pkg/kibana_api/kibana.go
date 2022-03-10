package kibana

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	myhttp "github.com/qonto/terraform-provider-kibana/pkg/httpClient"
)

type KibanaAPI interface {
	SetupClient(kibana_host, kibana_auth string)
	CreateAlertRule(alert Alert) (alertId string, err error)
	DeleteAlertRule(alertId string) error
	UpdateAlertRule(alertId string, alert Alert) error
	ReadAlertRule(alertId string) (Alert, error)
	DisableRule(alertId string) error
	EnableRule(alertId string) error
}

type KibanaClient struct {
	api     myhttp.HttpClientAPI
	headers map[string]string
	host    string
}

type Alert struct {
	Name       string            `json:"name,omitempty"`
	Tags       []string          `json:"tags,omitempty"`
	RuleTypeId string            `json:"rule_type_id,omitempty"`
	Schedule   map[string]string `json:"schedule,omitempty"`
	Throttle   string            `json:"throttle,omitempty"`
	NotifyWhen string            `json:"notify_when,omitempty"`
	Enabled    bool              `json:"enabled,omitempty"`
	Consumer   string            `json:"consumer,omitempty"`
	Params     json.RawMessage   `json:"params,omitempty"`
	Actions    []Action          `json:"actions,omitempty"`
}

type Action struct {
	Id     string          `json:"id"`
	Group  string          `json:"group"`
	Params json.RawMessage `json:"params"`
}
type FindResult struct {
	Alerts []Alert `json:"data"`
}
type ExecutionStatus struct {
	Status string `json:"status"`
}

func CreateNewKibanaClient() KibanaAPI {
	api := myhttp.CreateHTTPClient()
	return &KibanaClient{api: api}
}

func (c *KibanaClient) SetupClient(kibana_host, kibana_auth string) {
	headers := make(map[string]string, 2)
	headers["Authorization"] = fmt.Sprintf("Basic %s", kibana_auth)
	headers["Content-Type"] = "application/json"
	headers["kbn-xsrf"] = "terraform"
	c.headers = headers
	c.host = kibana_host
}

func (c *KibanaClient) CreateAlertRule(alert Alert) (alertId string, err error) {
	var result struct {
		Id string `json:"id"`
	}
	result.Id = ""
	url := fmt.Sprintf("%s/api/alerting/rule", c.host)
	jsonAlert, err := json.Marshal(alert)
	if err != nil {
		return "", err
	}
	r, statusCode, err := c.api.Post(url, c.headers, jsonAlert)
	if err != nil {
		return result.Id, errors.Wrapf(err, "Creating rule failed")
	}
	if statusCode != 200 && statusCode != 204 {
		return "", fmt.Errorf("Received status %d: %s\nRequest body:\n%s", statusCode, string(r), string(jsonAlert))
	}
	err = json.Unmarshal(r, &result)
	return result.Id, nil
}

func (c *KibanaClient) DeleteAlertRule(alertId string) error {
	url := fmt.Sprintf("%s/api/alerting/rule/%s", c.host, alertId)
	r, statusCode, err := c.api.Delete(url, c.headers)
	if err != nil {
		return errors.Wrapf(err, "Deleting rule failed")
	}
	if statusCode != 200 && statusCode != 204 {
		return fmt.Errorf("Received status %d: %s", statusCode, string(r))
	}
	return nil
}

func (c *KibanaClient) UpdateAlertRule(alertId string, alert Alert) error {
	url := fmt.Sprintf("%s/api/alerting/rule/%s", c.host, alertId)
	jsonAlert, err := json.Marshal(alert)
	if err != nil {
		return err
	}
	r, statusCode, err := c.api.Put(url, c.headers, jsonAlert)
	if err != nil {
		return errors.Wrapf(err, "Updating rule failed")
	}
	if statusCode != 200 && statusCode != 204 {
		return fmt.Errorf("Received status %d: %s\nRequest body:\n%s", statusCode, string(r), string(jsonAlert))
	}
	return nil
}

func (c *KibanaClient) ReadAlertRule(alertId string) (Alert, error) {
	var alert Alert
	url := fmt.Sprintf("%s/api/alerting/rule/%s", c.host, alertId)
	r, statusCode, err := c.api.Get(url, c.headers)
	if err != nil {
		return alert, errors.Wrapf(err, "Reading rule failed")
	}
	if statusCode != 200 && statusCode != 204 {
		return alert, fmt.Errorf("Received status %d: %s", statusCode, string(r))
	}
	err = json.Unmarshal(r, &alert)
	return alert, err
}

func (c *KibanaClient) EnableRule(alertId string) error {
	url := fmt.Sprintf("%s/api/alerting/rule/%s/_enable", c.host, alertId)
	r, statusCode, err := c.api.Post(url, c.headers, []byte{})
	if err != nil {
		return errors.Wrapf(err, "Enabling rule failed")
	}
	if statusCode != 200 && statusCode != 204 {
		return fmt.Errorf("Enable rule - Received status %d: %s", statusCode, string(r))
	}
	return nil
}

func (c *KibanaClient) DisableRule(alertId string) error {
	url := fmt.Sprintf("%s/api/alerting/rule/%s/_disable", c.host, alertId)
	r, statusCode, err := c.api.Post(url, c.headers, []byte{})
	if err != nil {
		return errors.Wrapf(err, "Disabling rule failed")
	}
	if statusCode != 200 && statusCode != 204 {
		return fmt.Errorf("Disable rule - Received status %d: %s", statusCode, string(r))
	}
	return nil
}

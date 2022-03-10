package kibana

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

type KibanaMockClient struct {
	CreateAlertShouldFail  bool
	CreateAlertId          string
	DeleteAlertShouldFail  bool
	UpdateAlertShouldFail  bool
	ReadAlertShouldFail    bool
	ReadAlertResult        Alert
	EnableAlertShouldFail  bool
	DisableAlertShouldFail bool
	alerts                 map[string]Alert
}

func (c *KibanaMockClient) CreateAlertRule(alert Alert) (alertId string, err error) {
	if c.CreateAlertShouldFail {
		return "", fmt.Errorf("Creating alert failed")
	}
	rand.Seed(time.Now().UnixNano())
	alertIdBuff := make([]byte, 16)
	rand.Read(alertIdBuff)
	alertId = base64.StdEncoding.EncodeToString(alertIdBuff)
	if c.alerts == nil {
		c.alerts = make(map[string]Alert)
	}
	c.alerts[alertId] = alert
	return alertId, nil
}

func (c *KibanaMockClient) DeleteAlertRule(alertId string) error {
	if c.DeleteAlertShouldFail {
		return fmt.Errorf("Deleting alert failed")
	}
	if c.alerts != nil {
		_, ok := c.alerts[alertId]
		delete(c.alerts, alertId)
		if !ok {
			return fmt.Errorf("Deleting alert failed - unknown id")
		}
	}
	return nil
}

func (c *KibanaMockClient) UpdateAlertRule(alertId string, alert Alert) error {
	if c.alerts == nil {
		c.alerts = make(map[string]Alert)
	}
	if c.UpdateAlertShouldFail {
		return fmt.Errorf("Updating alert failed")
	}
	_, ok := c.alerts[alertId]
	if ok {
		c.alerts[alertId] = alert
	} else {
		return fmt.Errorf("Failed updating alert - unknown alert id")
	}
	return nil
}

func (c *KibanaMockClient) ReadAlertRule(alertId string) (Alert, error) {
	if c.ReadAlertShouldFail {
		return Alert{}, fmt.Errorf("Reading alert failed")
	}
	if c.alerts != nil {
		alert, ok := c.alerts[alertId]
		if !ok {
			return Alert{}, fmt.Errorf("Alert not found")
		}
		return alert, nil
	}
	return Alert{}, fmt.Errorf("Alert not found")
}

func (c *KibanaMockClient) EnableRule(alertId string) error {
	if c.EnableAlertShouldFail || (alertId == "") {
		return fmt.Errorf("Enabling alert failed")
	}
	return nil
}

func (c *KibanaMockClient) DisableRule(alertId string) error {
	if c.DisableAlertShouldFail || (alertId == "") {
		return fmt.Errorf("Disabling alert failed")
	}
	return nil
}

func (c *KibanaMockClient) SetupClient(kibana_host, kibana_auth string) {}

package alerts_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jonboydell/logzio_client/alerts"
	"github.com/stretchr/testify/assert"
)

func TestAlerts_CreateAlert(t *testing.T) {
	underTest, err, teardown := setupAlertsTest()
	defer teardown()

	if assert.NoError(t, err) {
		mux.HandleFunc("/v1/alerts", func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)

			jsonBytes, _ := ioutil.ReadAll(r.Body)
			var target map[string]interface{}
			err = json.Unmarshal(jsonBytes, &target)
			assert.Contains(t, target, "title")
			assert.Contains(t, target, "description")
			assert.Contains(t, target, "query_string")
			assert.Contains(t, target, "severityThresholdTiers")

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, fixture("create_alert.json"))
			w.WriteHeader(http.StatusOK)
		})

		testAlert := alerts.CreateAlertType{
			Title:       "test create alert",
			Description: "this is my description",
			QueryString: "loglevel:ERROR",
			Filter:      "",
			Operation:   alerts.OperatorGreaterThan,
			SeverityThresholdTiers: []alerts.SeverityThresholdType{
				alerts.SeverityThresholdType{
					alerts.SeverityHigh,
					10,
				},
				alerts.SeverityThresholdType{
					alerts.SeverityInfo,
					10,
				},
			},
			SearchTimeFrameMinutes:       0,
			NotificationEmails:           []interface{}{},
			IsEnabled:                    true,
			SuppressNotificationsMinutes: 0,
			ValueAggregationType:         alerts.AggregationTypeCount,
			ValueAggregationField:        nil,
			GroupByAggregationFields:     []interface{}{"my_field"},
			AlertNotificationEndpoints:   []interface{}{},
		}

		alert, err := underTest.CreateAlert(testAlert)
		assert.NoError(t, err)
		assert.Equal(t, int64(1234567), alert.AlertId)
		assert.Equal(t, alerts.SeverityHigh, alert.SeverityThresholdTiers[0].Severity)
		assert.Equal(t, alerts.SeverityInfo, alert.SeverityThresholdTiers[1].Severity)
	}
}

func TestAlerts_CreateAlertAPIFail(t *testing.T) {
	underTest, err, teardown := setupAlertsTest()
	defer teardown()

	if assert.NoError(t, err) {
		mux.HandleFunc("/v1/alerts", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, fixture("create_alert_failed.txt"))
		})
	}

	_, err = underTest.CreateAlert(alerts.CreateAlertType{
		Title:       "test create alert",
		Description: "this is my description",
		QueryString: "loglevel:ERROR",
		Filter:      "",
		Operation:   alerts.OperatorGreaterThan,
		SeverityThresholdTiers: []alerts.SeverityThresholdType{
			alerts.SeverityThresholdType{
				alerts.SeverityHigh,
				10,
			},
		},
		SearchTimeFrameMinutes:       0,
		NotificationEmails:           []interface{}{},
		IsEnabled:                    true,
		SuppressNotificationsMinutes: 0,
		ValueAggregationType:         alerts.AggregationTypeCount,
		ValueAggregationField:        nil,
		GroupByAggregationFields:     []interface{}{"my_field"},
		AlertNotificationEndpoints:   []interface{}{},
	})
	assert.Error(t, err)
}

func TestAlerts_CreateAlertNoTitle(t *testing.T) {
	underTest, err, teardown := setupAlertsTest()
	defer teardown()

	_, err = underTest.CreateAlert(alerts.CreateAlertType{
		Description: "this is my description",
		QueryString: "loglevel:ERROR",
		Filter:      "",
		Operation:   alerts.OperatorGreaterThan,
		SeverityThresholdTiers: []alerts.SeverityThresholdType{
			alerts.SeverityThresholdType{
				alerts.SeverityHigh,
				10,
			},
		},
		SearchTimeFrameMinutes:       0,
		NotificationEmails:           []interface{}{},
		IsEnabled:                    true,
		SuppressNotificationsMinutes: 0,
		ValueAggregationType:         alerts.AggregationTypeCount,
		ValueAggregationField:        nil,
		GroupByAggregationFields:     []interface{}{"my_field"},
		AlertNotificationEndpoints:   []interface{}{},
	})
	assert.Error(t, err)
}

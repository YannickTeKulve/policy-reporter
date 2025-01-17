package loki_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/kyverno/policy-reporter/pkg/fixtures"
	"github.com/kyverno/policy-reporter/pkg/target"
	"github.com/kyverno/policy-reporter/pkg/target/loki"
)

type testClient struct {
	callback   func(req *http.Request)
	statusCode int
}

func (c testClient) Do(req *http.Request) (*http.Response, error) {
	c.callback(req)

	return &http.Response{
		StatusCode: c.statusCode,
	}, nil
}

func Test_LokiTarget(t *testing.T) {
	t.Run("Send Complete Result", func(t *testing.T) {
		callback := func(req *http.Request) {
			if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
				t.Errorf("Unexpected Content-Type: %s", contentType)
			}

			if agend := req.Header.Get("User-Agent"); agend != "Policy-Reporter" {
				t.Errorf("Unexpected Host: %s", agend)
			}

			if url := req.URL.String(); url != "http://localhost:3100/api/prom/push" {
				t.Errorf("Unexpected Host: %s", url)
			}

			expectedLine := fmt.Sprintf("[%s] %s", strings.ToUpper(fixtures.CompleteTargetSendResult.Priority.String()), fixtures.CompleteTargetSendResult.Message)
			labels, line := convertAndValidateBody(req, t)
			if line != expectedLine {
				t.Errorf("Unexpected LineContent: %s", line)
			}
			if !strings.Contains(labels, "policy=\""+fixtures.CompleteTargetSendResult.Policy+"\"") {
				t.Error("Missing Content for Label 'policy'")
			}
			if !strings.Contains(labels, "status=\""+string(fixtures.CompleteTargetSendResult.Result)+"\"") {
				t.Error("Missing Content for Label 'status'")
			}
			if !strings.Contains(labels, "priority=\""+fixtures.CompleteTargetSendResult.Priority.String()+"\"") {
				t.Error("Missing Content for Label 'priority'")
			}
			if !strings.Contains(labels, "source=\"policy-reporter\"") {
				t.Error("Missing Content for Label 'policy-reporter'")
			}
			if !strings.Contains(labels, "rule=\""+fixtures.CompleteTargetSendResult.Rule+"\"") {
				t.Error("Missing Content for Label 'rule'")
			}
			if !strings.Contains(labels, "category=\""+fixtures.CompleteTargetSendResult.Category+"\"") {
				t.Error("Missing Content for Label 'category'")
			}
			if !strings.Contains(labels, "severity=\""+string(fixtures.CompleteTargetSendResult.Severity)+"\"") {
				t.Error("Missing Content for Label 'severity'")
			}
			if !strings.Contains(labels, "custom=\"label\"") {
				t.Error("Missing Content for Label 'severity'")
			}

			res := fixtures.CompleteTargetSendResult.GetResource()
			if !strings.Contains(labels, "kind=\""+res.Kind+"\"") {
				t.Error("Missing Content for Label 'kind'")
			}
			if !strings.Contains(labels, "name=\""+res.Name+"\"") {
				t.Error("Missing Content for Label 'name'")
			}
			if !strings.Contains(labels, "uid=\""+string(res.UID)+"\"") {
				t.Error("Missing Content for Label 'uid'")
			}
			if !strings.Contains(labels, "namespace=\""+res.Namespace+"\"") {
				t.Error("Missing Content for Label 'namespace'")
			}
			if !strings.Contains(labels, "version=\""+fixtures.CompleteTargetSendResult.Properties["version"]+"\"") {
				t.Error("Missing Content for Label 'version'")
			}
		}

		client := loki.NewClient(loki.Options{
			ClientOptions: target.ClientOptions{
				Name: "Loki",
			},
			Host:         "http://localhost:3100/api/prom/push",
			CustomLabels: map[string]string{"custom": "label"},
			HTTPClient:   testClient{callback, 200},
		})
		client.Send(fixtures.CompleteTargetSendResult)
	})

	t.Run("Send Minimal Result", func(t *testing.T) {
		callback := func(req *http.Request) {
			if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
				t.Errorf("Unexpected Content-Type: %s", contentType)
			}

			if agend := req.Header.Get("User-Agent"); agend != "Policy-Reporter" {
				t.Errorf("Unexpected Host: %s", agend)
			}

			if url := req.URL.String(); url != "http://localhost:3100/api/prom/push" {
				t.Errorf("Unexpected Host: %s", url)
			}

			expectedLine := fmt.Sprintf("[%s] %s", strings.ToUpper(fixtures.MinimalTargetSendResult.Priority.String()), fixtures.MinimalTargetSendResult.Message)
			labels, line := convertAndValidateBody(req, t)
			if line != expectedLine {
				t.Errorf("Unexpected LineContent: %s", line)
			}
			if !strings.Contains(labels, "policy=\""+fixtures.MinimalTargetSendResult.Policy+"\"") {
				t.Error("Missing Content for Label 'policy'")
			}
			if !strings.Contains(labels, "status=\""+string(fixtures.MinimalTargetSendResult.Result)+"\"") {
				t.Error("Missing Content for Label 'status'")
			}
			if !strings.Contains(labels, "priority=\""+fixtures.MinimalTargetSendResult.Priority.String()+"\"") {
				t.Error("Missing Content for Label 'priority'")
			}
			if !strings.Contains(labels, "source=\"policy-reporter\"") {
				t.Error("Missing Content for Label 'policy-reporter'")
			}
			if strings.Contains(labels, "rule") {
				t.Error("Unexpected Label 'rule'")
			}
			if strings.Contains(labels, "category") {
				t.Error("Unexpected Label 'category'")
			}
			if strings.Contains(labels, "severity") {
				t.Error("Unexpected 'severity'")
			}
			if strings.Contains(labels, "kind") {
				t.Error("Unexpected Label 'kind'")
			}
			if strings.Contains(labels, "name") {
				t.Error("Unexpected 'name'")
			}
			if strings.Contains(labels, "uid") {
				t.Error("Unexpected 'uid'")
			}
			if strings.Contains(labels, "namespace") {
				t.Error("Unexpected 'namespace'")
			}
		}

		client := loki.NewClient(loki.Options{
			ClientOptions: target.ClientOptions{
				Name: "Loki",
			},
			Host:         "http://localhost:3100/api/prom/push",
			CustomLabels: map[string]string{"custom": "label"},
			HTTPClient:   testClient{callback, 200},
		})
		client.Send(fixtures.MinimalTargetSendResult)
	})
	t.Run("Name", func(t *testing.T) {
		client := loki.NewClient(loki.Options{
			ClientOptions: target.ClientOptions{
				Name: "Loki",
			},
			Host:         "http://localhost:3100/api/prom/push",
			CustomLabels: map[string]string{"custom": "label"},
			HTTPClient:   testClient{},
		})

		if client.Name() != "Loki" {
			t.Errorf("Unexpected Name %s", client.Name())
		}
	})
}

func convertAndValidateBody(req *http.Request, t *testing.T) (string, string) {
	payload := make(map[string]interface{})

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		t.Fatal(err)
	}

	streamsContent, ok := payload["streams"]
	if !ok {
		t.Errorf("Expected payload key 'streams' is missing")
	}

	streams := streamsContent.([]interface{})
	if len(streams) != 1 {
		t.Errorf("Expected one streams entry")
	}

	firstStream := streams[0].(map[string]interface{})
	entriesContent, ok := firstStream["entries"]
	if !ok {
		t.Errorf("Expected stream key 'entries' is missing")
	}
	labels, ok := firstStream["labels"]
	if !ok {
		t.Errorf("Expected stream key 'labels' is missing")
	}

	entryContent := entriesContent.([]interface{})[0]
	entry := entryContent.(map[string]interface{})
	_, ok = entry["ts"]
	if !ok {
		t.Errorf("Expected entry key 'ts' is missing")
	}
	line, ok := entry["line"]
	if !ok {
		t.Errorf("Expected entry key 'line' is missing")
	}

	return labels.(string), line.(string)
}

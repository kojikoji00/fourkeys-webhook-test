package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

var (
	projectName     string
	eventHandlerURL string
)

type CustomEvent struct {
	EventType   string                 `json:"event_type"`
	ID          string                 `json:"id"`
	Metadata    map[string]interface{} `json:"metadata"`
	TimeCreated time.Time              `json:"time_created"`
	Signature   string                 `json:"signature"`
	MsgID       string                 `json:"msg_id"`
	Source      string                 `json:"source"`
}

func main() {

	projectName := os.Getenv("PROJECT_NAME")
	eventHandlerURL := os.Getenv("WEBHOOK")
	secret, err := getSecret(projectName, "custom-ev-handler", "latest")
	if err != nil {
		log.Fatalf("Failed to get secret: %v", err)
	}

	// Define custom event and marshal to JSON
	data := CustomEvent{
		EventType: "custom",
		ID:        "abc123",
		Metadata: map[string]interface{}{
			"id":           "abc123",
			"time_created": time.Now().Format(time.RFC3339Nano),
			"key1":         "value1",
			"key2":         "value2",
		},
		TimeCreated: time.Now(),
		MsgID:       "123",
		Source:      "custom",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal message data: %v", err)
	}

	signature := GenerateSignature(secret, jsonData)

	// Create and send HTTP POST request
	req, err := http.NewRequest("POST", eventHandlerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "custom")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Custom-Event", "custom")
	req.Header.Set("Custom-Signature", signature)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer res.Body.Close()

	log.Printf("Response Status Code: %d", res.StatusCode)

	// respBody, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatalf("Failed to read response body: %v", err)
	// }
	// fmt.Printf("respBody: %s", string(respBody))
}

func GenerateSignature(secret, data []byte) string {
	mac := hmac.New(sha1.New, secret)
	mac.Write(data)
	return fmt.Sprintf("sha1=%x", mac.Sum(nil))
}

func getSecret(projectName, secretName, versionNum string) ([]byte, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret manager client: %v", err)
	}
	defer client.Close()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectName, secretName, versionNum),
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %v", err)
	}

	return result.Payload.Data, nil
}

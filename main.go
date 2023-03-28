package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

var (
	projectName     string
	eventHandlerURL string
)

// type CustomEvent struct {
// 	EventType   string  `json:"event_type"`
// 	ID          string  `json:"id"`
// 	Metadata    string  `json:"metadata"`
// 	TimeCreated string  `json:"time_created"`
// 	Custom-Signature   string  `json:"signature"`
// 	MsgID       string  `json:"msg_id"`
// 	Source      string  `json:"source"`
// 	Payload     Payload `json:"payload"`
// }

func main() {
	projectName = os.Getenv("PROJECT_NAME")
	eventHandlerURL = os.Getenv("WEBHOOK")

	// url := "https://example.com/custom_event"
	data := map[string]interface{}{
		"EventType": "custom",
		"Data": map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}
	// data := map[string]interface{}{
	// 	"EventType": "custom",
	// 	"Data": map[string]string{
	// 		"key1": "value1",
	// 		"key2": "value2",
	// 	},
	// }

	body, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal message data: %v", err)
	}

	req, err := http.NewRequest("POST", eventHandlerURL, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
	}

	// Generate signature
	expectedSignature := "sha1="
	secret, err := getSecret(projectName, "custom-ev-handler", "latest")
	if err != nil {
		log.Printf("Failed to get secret: %v", err)
	}
	hashed := hmac.New(sha1.New, secret)
	hashed.Write(body)
	expectedSignature += hex.EncodeToString(hashed.Sum(nil))

	// Set headers
	req.Header.Set("User-Agent", "custom")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Custom-Event", "custom")
	req.Header.Set("Custom-Signature", expectedSignature)

	// Send request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
	}
	defer res.Body.Close()

	// Handle response
	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
	}
	fmt.Println(string(respBody))
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

// これはHTTPリクエストを加えるとエラーになってしまったコード

// package main

// import (
// 	"bytes"
// 	"context"
// 	"crypto/hmac"
// 	"crypto/sha1"
// 	"encoding/base64"
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"

// 	secretmanager "cloud.google.com/go/secretmanager/apiv1"
// 	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
// )

// type Payload struct {
// 	CreatedAt  string `json:"created_at"`
// 	Repository struct {
// 		Name string `json:"name"`
// 		Url  string `json:"url"`
// 	} `json:"repository"`
// }

// type SampleEvent struct {
// 	EventType   string  `json:"event_type"`
// 	ID          string  `json:"id"`
// 	Metadata    string  `json:"metadata"`
// 	TimeCreated string  `json:"time_created"`
// 	Signature   string  `json:"signature"`
// 	MsgID       string  `json:"msg_id"`
// 	Source      string  `json:"source"`
// 	Payload     Payload `json:"payload"`
// }

// var (
// 	projectName     string
// 	eventHandlerURL string
// )

// func main() {
// 	projectName = os.Getenv("PROJECT_NAME")
// 	eventHandlerURL = os.Getenv("WEBHOOK")

// 	payload := Payload{
// 		CreatedAt: "2022-03-27T09:00:00Z",
// 		Repository: struct {
// 			Name string `json:"name"`
// 			Url  string `json:"url"`
// 		}{
// 			Name: "example-repo",
// 			Url:  "https://github.com/example/repo",
// 		},
// 	}

// 	// イベントデータを作成
// 	data := SampleEvent{
// 		EventType:   "sample",
// 		ID:          "your-id",
// 		Metadata:    "your-metadata",
// 		TimeCreated: "your-time-created",
// 		MsgID:       "your-msg-id",
// 		Source:      "sample",
// 		Payload:     payload,
// 	}

// 	// イベントデータを JSON に変換
// 	msg, err := json.Marshal(data)
// 	if err != nil {
// 		log.Fatalf("Failed to marshal message data: %v", err)
// 	}

// 	// ハッシュ値を生成
// 	expectedSignature := "sha1="
// 	secret, err := getSecret(projectName, "custom-ev-handler", "latest")
// 	if err != nil {
// 		panic(err)
// 	}

// 	hashed := hmac.New(sha1.New, secret)
// 	hashed.Write(msg)
// 	expectedSignature += hex.EncodeToString(hashed.Sum(nil))

// 	// Pub/Sub形式のメッセージを作成
// 	pubsubMessage := map[string]interface{}{
// 		"attributes": map[string]interface{}{
// 			"Sample-Event":     data.EventType,
// 			"Sample-Signature": expectedSignature,
// 			"Content-Type":     "application/json",
// 		},
// 		"data": base64.StdEncoding.EncodeToString(msg),
// 	}

// 	pubsubMessageBytes, err := json.Marshal(pubsubMessage)
// 	if err != nil {
// 		log.Fatalf("Failed to marshal pubsub message data: %v", err)
// 	}

// 	// HTTP リクエストを作成
// 	req, err := http.NewRequest("POST", eventHandlerURL, bytes.NewBuffer(pubsubMessageBytes))
// 	if err != nil {
// 		panic(err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Sample-Event", data.EventType)
// 	req.Header.Set("Sample-Signature", expectedSignature)
// 	req.Header.Set("User-Agent", "sample")

// 	// HTTP リクエストを送信
// 	// client := &http.Client{}
// 	client := http.DefaultClient
// 	res, err := client.Do(req)
// 	if err != nil {
// 		log.Fatalf("Failed to send request: %v", err)
// 	}
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		log.Fatalf("Failed to read response body: %v", err)
// 	}
// 	defer res.Body.Close()
// 	log.Printf("Response Body: %v", string(body))

// 	fmt.Println("Response status code:", res.StatusCode)

// }

// func getSecret(projectName, secretName, versionNum string) ([]byte, error) {
// 	ctx := context.Background()
// 	client, err := secretmanager.NewClient(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer client.Close()

// 	req := &secretmanagerpb.AccessSecretVersionRequest{
// 		Name: "projects/" + projectName + "/secrets/" + secretName + "/versions/" + versionNum,
// 	}

// 	result, err := client.AccessSecretVersion(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result.Payload.Data, nil
// }

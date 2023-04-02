package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	data := CustomEvent{
		EventType:   "custom",
		ID:          "abc123",
		// Metadata:    `{"id":"abc123", "time_created":"` + time.Now().UTC().Format(time.RFC3339Nano) + `","key1":"value1", "key2":"value2"}`,
		Metadata: map[string]interface{}{
			"id":"abc123",
			"time_created": time.Now().UTC().Format(time.RFC3339Nano),
			"key1":"value1",
			"key2":"value2",
		},
		TimeCreated: time.Now().UTC(),
		MsgID:       "123",
		Source:      "custom",
	}
	// data := map[string]interface{}{
	// 	"EventType":   "custom",
	// 	"ID":          "abc123",
	// 	"Metadata":    map[string]interface{}{"id": "abc123", "time_created": time.Now().UTC().Format(time.RFC3339Nano), "key1": "value1", "key2": "value2"},
	// 	"TimeCreated": time.Now().UTC(),
	// 	"MsgID":       "123",
	// 	"Source":      "custom",
	// }

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal message data: %v", err)
	}

	signature := GenerateSignature(secret, jsonData)
	fmt.Printf("Signature: %s\n", signature)

	req, err := http.NewRequest("POST", eventHandlerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "custom")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Custom-Event", "custom")
	req.Header.Set("Custom-Signature", signature)

	fmt.Printf("req: %s", req)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer res.Body.Close()

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	fmt.Printf("respBody: %s", string(respBody))
}

// func main() {
// 	projectName = os.Getenv("PROJECT_NAME")
// 	eventHandlerURL = os.Getenv("WEBHOOK")

// 	data := CustomEvent{
// 		EventType:   "custom",
// 		ID:          "abc123",
// 		Metadata:    `{"id":"abc123", "time_created":"` + time.Now().UTC().Format(time.RFC3339Nano) + `","key1":"value1", "key2":"value2"}`,
// 		TimeCreated: time.Now().UTC(),
// 		MsgID:       "123",
// 		Source:      "custom",
// 		// Signature:   "",
// 	}

// 	// Generate signature
// 	expectedSignature := "sha1="
// 	secret, err := getSecret(projectName, "custom-ev-handler", "latest")
// 	if err != nil {
// 		log.Printf("Failed to get secret: %v", err)
// 	}
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		log.Fatalf("Failed to marshal message data: %v", err)
// 	}

// 	expectedSignature = GenerateSignature(secret, jsonData)
// 	data.Signature = expectedSignature

// 	body, err := json.Marshal(data)
// 	if err != nil {
// 		log.Fatalf("Failed to marshal message data: %v", err)
// 	}

// 	req, err := http.NewRequest("POST", eventHandlerURL, bytes.NewBuffer(body))
// 	if err != nil {
// 		log.Printf("Failed to create request: %v", err)
// 	}

// 	// hashed := hmac.New(sha1.New, secret)
// 	// hashed.Write(jsonData)

// 	fmt.Printf("expectedSignature: %s", expectedSignature)

// 	// hashed := hmac.New(sha1.New, secret)

// 	// // // data.Signature = ""
// 	// // jsonString, _ := json.Marshal(data)
// 	// // log.Printf("JSON string: %s", jsonString)

// 	// // Decode the JSON string into an empty interface
// 	// // hashed = hmac.New(sha1.New, secret)
// 	// hashed.Write(jsonData)
// 	// expectedSignature += hex.EncodeToString(hashed.Sum(nil))

// 	// body, err := json.Marshal(data)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to marshal message data: %v", err)
// 	// }

// 	// req, err := http.NewRequest("POST", eventHandlerURL, bytes.NewBuffer(body))
// 	// if err != nil {
// 	// 	log.Printf("Failed to create request: %v", err)
// 	// }
// 	req.Header.Set("User-Agent", "custom")
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Custom-Event", "custom")
// 	req.Header.Set("Custom-Signature", data.Signature)

// 	fmt.Printf("req: %s", req)

// 	// Set headers

// 	// Send request
// 	client := &http.Client{}
// 	res, err := client.Do(req)
// 	if err != nil {
// 		log.Printf("Failed to send request: %v", err)
// 	}
// 	defer res.Body.Close()

// 	// Handle response
// 	respBody, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		log.Printf("Failed to read response body: %v", err)
// 	}
// 	fmt.Println(string(respBody))
// }

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

// var (
// 	projectName     string
// 	eventHandlerURL string
// )

// type Repository struct {
// 	ID       int    `json:"id"`
// 	Name     string `json:"name"`
// 	FullName string `json:"full_name"`
// }

// type Pusher struct {
// 	Name  string `json:"name"`
// 	Email string `json:"email"`
// }

// type CustomEvent struct {
// 	Ref        string     `json:"ref"`
// 	Repository Repository `json:"repository"`
// 	Pusher     Pusher     `json:"pusher"`
// 	Signature  string     `json:"-"`
// }

// func main() {
// 	projectName = os.Getenv("PROJECT_NAME")
// 	eventHandlerURL = os.Getenv("WEBHOOK")

// 	data := CustomEvent{
// 		Ref: "refs/heads/main",
// 		Repository: Repository{
// 			ID:       1,
// 			Name:     "repo",
// 			FullName: "user/repo",
// 		},
// 		Pusher: Pusher{
// 			Name:  "user",
// 			Email: "user@example.com",
// 		},
// 	}

// 	// Generate signature

// 	secret, err := getSecret(projectName, "custom-ev-handler", "latest")
// 	if err != nil {
// 		log.Fatalf("Failed to get secret: %v", err)
// 	}
// 	hashed := hmac.New(sha1.New, secret)

// 	data.Signature = ""
// 	jsonString, _ := json.Marshal(data)
// 	hashed.Write([]byte(jsonString))
// 	data.Signature = "sha1=" + hex.EncodeToString(hashed.Sum(nil))

// 	log.Printf("data.Signature: %s", data.Signature)

// 	// Encode JSON as Base64
// 	// jsonData, err := json.Marshal(data)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to marshal JSON data: %v", err)
// 	// }

// 	// Encode JSON as Base64
// 	encodedData := base64.StdEncoding.EncodeToString(jsonString)
// 	message := map[string]interface{}{
// 		"data": encodedData,
// 		"attributes": map[string]string{
// 			"headers": "{}",
// 		},
// 		"message_id":   "123",
// 		"publish_time": "2022-09-29T12:34:56.789Z",
// 	}
// 	body, err := json.Marshal(message)
// 	if err != nil {
// 		log.Fatalf("Failed to marshal message data: %v", err)
// 	}

// 	req, err := http.NewRequest("POST", eventHandlerURL, bytes.NewBuffer(body))
// 	if err != nil {
// 		log.Printf("Failed to create request: %v", err)
// 	}

// 	// Hash the encoded payload with the secret
// 	// hash := hmac.New(sha1.New, secret)
// 	// hash.Write([]byte(encodedData))
// 	// signature := hex.EncodeToString(hash.Sum(nil))
// 	// secret, err := getSecret(projectName, "custom-ev-handler", "latest")
// 	// if err != nil {
// 	// 	log.Printf("Failed to get secret: %v", err)
// 	// }
// 	// hashed := hmac.New(sha1.New, secret)

// 	// // Set data.Signature to an empty string for signature calculation
// 	// data.Signature = ""
// 	// jsonString, _ := json.Marshal(data)
// 	// hashed.Write(jsonString)
// 	// data.Signature = "sha1=" + hex.EncodeToString(hashed.Sum(nil))

// 	// log.Printf("data.Signature: %s", data.Signature)

// 	// // Encode JSON as Base64
// 	// encodedData := base64.StdEncoding.EncodeToString(jsonString)
// 	// message := map[string]interface{}{
// 	// 	"message": map[string]interface{}{
// 	// 		"data":         encodedData,
// 	// 		"attributes":   map[string]string{"headers": "{}"},
// 	// 		"message_id":   "123",
// 	// 		"publish_time": "2022-09-29T12:34:56.789Z",
// 	// 	},
// 	// }
// 	// body, err := json.Marshal(message)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to marshal message data: %v", err)
// 	// }

// 	// req, err := http.NewRequest("POST", eventHandlerURL, bytes.NewBuffer(jsonData))
// 	// if err != nil {
// 	// 	log.Printf("Failed to create request: %v", err)
// 	// }

// 	// Set headers
// 	req.Header.Set("User-Agent", "custom")
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Custom-Event", "custom")
// 	req.Header.Set("Custom-Signature", data.Signature)

// 	// Send request
// 	client := &http.Client{}
// 	res, err := client.Do(req)
// 	if err != nil {
// 		log.Printf("Failed to send request: %v", err)
// 	}
// 	defer res.Body.Close()

// 	// Handle response
// 	respBody, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		log.Printf("Failed to read response body: %v", err)
// 	}
// 	fmt.Println(string(respBody))
// }

// func getSecret(projectName, secretName, versionNum string) ([]byte, error) {
// 	ctx := context.Background()
// 	client, err := secretmanager.NewClient(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create secret manager client: %v", err)
// 	}
// 	defer client.Close()

// 	req := &secretmanagerpb.AccessSecretVersionRequest{
// 		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectName, secretName, versionNum),
// 	}

// 	result, err := client.AccessSecretVersion(ctx, req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to access secret version: %v", err)
// 	}

// 	return result.Payload.Data, nil
// }

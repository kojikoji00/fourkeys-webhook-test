// package main

// import (
// 	"bytes"
// 	"context"
// 	"crypto/hmac"
// 	"crypto/sha1"
// 	"encoding/hex"
// 	"log"
// 	"encoding/json"
// 	"encoding/base64"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"os"

// 	secretmanager "cloud.google.com/go/secretmanager/apiv1"
// 	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
// )

// // const secretKey = "your_secret_key"

// var (
// 	// secretKey 			 string
// 	projectID      	 	 string
// 	eventHandlerURL      string
// )

// func main() {
//  	projectID = os.Getenv("PROJECT_ID")
// 	// secretKey = os.Getenv("SECRET")
// 	eventHandlerURL = os.Getenv("WEBHOOK")

// 	payload := map[string]interface{}{"key1": "value1", "key2": "value2"}
// 	payloadBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		log.Fatalf("Failed to marshal payload: %v", err)
// 	}

// 	// イベントデータを作成
// 	data := map[string]interface{}{
// 		"event_type":  "sample",
// 		"id":          "your-id",
// 		"metadata":    "your-metadata",
// 		"time_created": "your-time-created",
// 		"msg_id":       "your-msg-id",
// 		"source":      "sample",
// 		"payload": 		map[string]interface{}{"key1": "value1", "key2": "value2"},
// 	}

// 	// イベントデータを JSON に変換
// 	msg, err := json.Marshal(data)
// 	if err != nil {
// 		log.Fatalf("Failed to marshal message data: %v", err)
// 	}

// 	expectedSignature := "sha1="
// 	secret, err := getSecret(projectID, "custom-ev-handler", "latest")
// 	if err != nil {
// 		panic(err)
// 	}

// 	hashed := hmac.New(sha1.New, secret)
// 	hashed.Write(msg)
// 	expectedSignature += hex.EncodeToString(hashed.Sum(nil))

// 	fmt.Println(expectedSignature)

// 	// Pub/Sub形式のメッセージを作成
// 	pubsubMessage := map[string]interface{}{
// 		"attributes": map[string]interface{}{
// 			// "headers": string(json.RawMessage(`{
// 			// 	"Sample-Event": "` + data["event_type"].(string) + `",
// 			// 	"Sample-Signature": "` + expectedSignature + `",
// 			// 	"Content-Type": "application/json"
// 			// }`)),
// 			// "headers": string(json.RawMessage(`{
// 			"Sample-Event": data["event_type"].(string),
// 			"Sample-Signature": expectedSignature,
// 			"Content-Type": "application/json",
// 			// }`)),
// 		},
// 		"data":       base64.StdEncoding.EncodeToString(payloadBytes),
// 		// "message": map[string]interface{}{
// 		// 	"attributes": map[string]interface{}{
// 		// 		// "headers": string(json.RawMessage(`{
// 		// 		// 	"Sample-Event": "` + data["event_type"].(string) + `",
// 		// 		// 	"Sample-Signature": "` + expectedSignature + `",
// 		// 		// 	"Content-Type": "application/json"
// 		// 		// }`)),
// 		// 		// "headers": string(json.RawMessage(`{
// 		// 		"Sample-Event": data["event_type"].(string),
// 		// 		"Sample-Signature": expectedSignature,
// 		// 		"Content-Type": "application/json",
// 		// 		// }`)),
// 		// 	},
// 		// 	"data":       base64.StdEncoding.EncodeToString(payloadBytes),
// 		// 	"message_id": data["msg_id"].(string),
// 		// },
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
// 	req.Header.Set("Sample-Event", data["event_type"].(string))
// 	req.Header.Set("Sample-Signature", expectedSignature)
// 	req.Header.Set("User-Agent", "sample") // イベントハンドラが受信を許可するソースを設定
// 	fmt.Println(req)
// 	// HTTP クライアントを作成してリクエストを送信
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	// レスポンスを表示
// 	fmt.Printf("Status: %s\n", resp.Status)

// 	respBody, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Response Body:", string(respBody))
// }

// func getSecret(projectID, secretName, versionNum string) ([]byte, error) {
// 	ctx := context.Background()
// 	client, err := secretmanager.NewClient(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer client.Close()

// 	req := &secretmanagerpb.AccessSecretVersionRequest{
// 		Name: "projects/" + projectID + "/secrets/" + secretName + "/versions/" + versionNum,
// 	}

// 	result, err := client.AccessSecretVersion(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result.Payload.Data, nil
// }

package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type Payload struct {
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Repository struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"repository"`
}

type SampleEvent struct {
	EventType   string  `json:"event_type"`
	ID          string  `json:"id"`
	Metadata    string  `json:"metadata"`
	TimeCreated string  `json:"time_created"`
	Signature   string  `json:"signature"`
	MsgID       string  `json:"msg_id"`
	Source      string  `json:"source"`
	Payload     Payload `json:"payload"`
}

var (
	projectID       string
	eventHandlerURL string
)

func main() {
	projectID = os.Getenv("PROJECT_ID")
	eventHandlerURL = os.Getenv("WEBHOOK")

	payload := Payload{
		CreatedAt: "2022-03-27T09:00:00Z",
		Repository: struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		}{
			Name: "example-repo",
			Url:  "https://github.com/example/repo",
		},
	}

	// イベントデータを作成
	data := SampleEvent{
		EventType:   "sample",
		ID:          "your-id",
		Metadata:    "your-metadata",
		TimeCreated: "your-time-created",
		MsgID:       "your-msg-id",
		Source:      "sample",
		Payload:     payload,
	}

	// イベントデータを JSON に変換
	msg, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal message data: %v", err)
	}

	// ハッシュ値を生成
	expectedSignature := "sha1="
	secret, err := getSecret(projectID, "custom-ev-handler", "latest")
	if err != nil {
		panic(err)
	}

	hashed := hmac.New(sha1.New, secret)
	hashed.Write(msg)
	expectedSignature += hex.EncodeToString(hashed.Sum(nil))

	// Pub/Sub形式のメッセージを作成
	pubsubMessage := map[string]interface{}{
		"attributes": map[string]interface{}{
			"Sample-Event":     data.EventType,
			"Sample-Signature": expectedSignature,
			"Content-Type":     "application/json",
		},
		"data": base64.StdEncoding.EncodeToString(msg),
	}

	pubsubMessageBytes, err := json.Marshal(pubsubMessage)
	if err != nil {
		log.Fatalf("Failed to marshal pubsub message data: %v", err)
	}

	// HTTP リクエストを作成
	req, err := http.NewRequest("POST", eventHandlerURL, bytes.NewBuffer(pubsubMessageBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Sample-Event", data.EventType)
	req.Header.Set("Sample-Signature", expectedSignature)
	req.Header.Set("User-Agent", "sample")


}

func getSecret(projectID, secretName, versionNum string) ([]byte, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/" + projectID + "/secrets/" + secretName + "/versions/" + versionNum,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, err
	}

	return result.Payload.Data, nil
}

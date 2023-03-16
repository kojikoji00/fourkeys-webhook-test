package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const secretKey = "your_secret_key"

func main() {
	ctx := context.Background()
	// 送信先のイベントハンドラーのURL
	url := "EVENT_HANDLER_URL"

	// 送信するデータ
	data := map[string]interface{}{
		"event":   "my_custom_event",
		"action":  "update",
		"payload": map[string]interface{}{"key1": "value1", "key2": "value2"},
		"source":  "sample",
	}

	// データをJSONにエンコードする
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// client, err := idtoken.NewClient(ctx, url)
	// if err != nil {
	// 	log.Fatalf("Failed to create new client: %v", err)
	// }

	// HTTPリクエストを作成する
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// 認証トークンを設定する
	// ヘッダーを追加する
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(string(jsonData))))

	// Generate the HMAC-SHA256 signature
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write(jsonData)
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	req.Header.Set("X-Signature", signature)

	// HTTPリクエストを送信する
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	// HTTPレスポンスを読み込む
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// HTTPレスポンスを表示する
	fmt.Println(string(body))
}

// 認証を必要とする
// func main() {
// 	ctx := context.Background()

// 	// 送信先のイベントハンドラーのURL
// 	url := "example_event_handler"

// 	// 送信するデータ
// 	data := map[string]interface{}{
// 		"event":   "my_custom_event",
// 		"action":  "update",
// 		"payload": map[string]interface{}{"key1": "value1", "key2": "value2"},
// 		"source":  "sample",
// 	}

// 	// データをJSONにエンコードする
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// サービスアカウントの秘密鍵ファイルを読み込む
// 	jsonKey, err := ioutil.ReadFile("key.json")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// OAuth2の設定を作成する
// 	config, err := google.ConfigFromJSON(jsonKey, "https://www.googleapis.com/auth/cloud-platform", "https://www.googleapis.com/auth/compute")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 認証コードを取得するためのURLを生成する
// 	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

// 	// 認証コードを取得するための誘導
// 	fmt.Printf("Access the following URL and input the authorization code: %v\n", authURL)
// 	fmt.Print("Authorization Code: ")
// 	var code string
// 	fmt.Scanln(&code)

// 	// 認証コードを用いてアクセストークンを取得する
// 	token, err := config.Exchange(ctx, code)
// 	if err != nil {
// 		panic(err)
// 	}

// 	client := config.Client(ctx, token)

// 	// OAuth2のクライアントを作成する

// 	// HTTPリクエストを作成する
// 	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 認証トークンを設定する
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

// 	// ヘッダーを追加する
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(jsonData)))

// 	// httpClient := &http.Client{}

// 	// HTTPリクエストを送信する
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	// HTTPレスポンスを読み込む
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// HTTPレスポンスを表示する
// 	fmt.Println(string(body))
// }

// // 認証不要
// func main() {
// 	ctx := context.Background()

// 	// 送信先のCloud RunのURL
// 	url := "example_event_handler"

// 	// 送信するデータ
// 	data := map[string]interface{}{
// 		"event":   "my_custom_event",
// 		"action":  "update",
// 		"payload": map[string]interface{}{"key1": "value1", "key2": "value2"},
// 		"source":  "sample",
// 	}

// 	// データをJSONにエンコードする
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Google Cloudの認証を行う
// 	client, err := google.DefaultClient(oauth2.NoContext, "https://www.googleapis.com/auth/cloud-platform")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// HTTPリクエストを作成する
// 	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		panic(err)
// 	}

//     // 認証トークンを設定する
//     _, err = client.Transport.(*oauth2.Transport).RoundTrip(req)
//     if err != nil {
//         panic(err)
//     }

// 	// ヘッダーを追加する
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(jsonData)))

// 	// HTTPリクエストを送信する
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	// HTTPレスポンスを読み込む
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// HTTPレスポンスを表示する
// 	fmt.Println(string(body))
// }

package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "time"
)

const (
    targetURL  = "https://www.target.com"
    webhookURL = "https://discord.com/api/webhooks/<link webhooks>"
    threshold  = 2 * time.Second //time foe delay 
)

type DiscordWebhook struct {
    Content string `json:"content"`
}

func checkWebsite() {
    start := time.Now()
    resp, err := http.Get(targetURL)
    elapsed := time.Since(start)

    if err != nil {
        log.Printf("Error checking website: %v", err)
        sendToDiscord("Error checking website: " + err.Error())
        return
    }
    defer resp.Body.Close()

    if elapsed > threshold {
        //resolve time stamp for server 
        timestamp := time.Now().Format("2006-01-02 15:04:05")
        message := "Website response time: " + elapsed.String() + " at " + timestamp
        log.Println(message)
        sendToDiscord(message)
    }
}

func sendToDiscord(message string) {
    webhook := DiscordWebhook{Content: message}
    payload, err := json.Marshal(webhook)
    if err != nil {
        log.Printf("Error marshalling JSON: %v", err)
        return
    }

    _, err = http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
    if err != nil {
        log.Printf("Error sending to Discord: %v", err)
    }
}

func main() {
    ticker := time.NewTicker(4 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        checkWebsite()
    }
}

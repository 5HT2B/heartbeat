package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/valyala/fasthttp"
)

var (
	// VAPID configuration
	vapidPublicKey  string
	vapidPrivateKey string
	vapidSubject    string

	// Push monitoring configuration
	pushCheckInterval       = 5 * time.Minute
	pushInactivityThreshold = 10 * time.Minute
)

// StoredPushSubscription represents a push subscription stored in Redis
type StoredPushSubscription struct {
	Endpoint   string `json:"endpoint"`
	P256dhKey  string `json:"p256dh"`
	AuthKey    string `json:"auth"`
	DeviceName string `json:"device_name"`
	AuthToken  string `json:"auth_token"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

// PushSubscribeRequest represents the request from the PWA
type PushSubscribeRequest struct {
	Endpoint   string                 `json:"endpoint"`
	Keys       map[string]interface{} `json:"keys"`
	Encoding   string                 `json:"encoding,omitempty"`
	DeviceName string                 `json:"deviceName"`
}

// InitPushConfiguration loads VAPID keys from environment variables
func InitPushConfiguration() {
	vapidPublicKey = os.Getenv("VAPID_PUBLIC_KEY")
	vapidPrivateKey = os.Getenv("VAPID_PRIVATE_KEY")
	vapidSubject = os.Getenv("VAPID_SUBJECT")

	if vapidSubject == "" {
		vapidSubject = "mailto:admin@josie.health"
	}

	// Load push monitoring configuration from environment
	if interval := os.Getenv("PUSH_CHECK_INTERVAL"); interval != "" {
		if d, err := time.ParseDuration(interval); err == nil {
			pushCheckInterval = d
		}
	}

	if threshold := os.Getenv("PUSH_INACTIVITY_THRESHOLD"); threshold != "" {
		if d, err := time.ParseDuration(threshold); err == nil {
			pushInactivityThreshold = d
		}
	}

	LogVAPIDConfiguration()
}

// SavePushSubscription stores a push subscription in Redis
func SavePushSubscription(sub StoredPushSubscription) error {
	key := fmt.Sprintf("push_subscription:%s:%s", sub.DeviceName, sub.AuthToken)
	sub.UpdatedAt = time.Now().Unix()

	if sub.CreatedAt == 0 {
		sub.CreatedAt = sub.UpdatedAt
	}

	_, err := rjh.JSONSet(key, ".", sub)
	if err != nil {
		return fmt.Errorf("error saving push subscription: %v", err)
	}

	// Also add to a set for easy retrieval of all subscriptions
	setKey := fmt.Sprintf("push_subscriptions:%s", sub.AuthToken)
	ctx := context.Background()
	if _, err := rdb.SAdd(ctx, setKey, key).Result(); err != nil {
		log.Printf("- Warning: failed to add subscription to set: %v", err)
	}

	return nil
}

// GetPushSubscriptionsByDevice retrieves all push subscriptions for a device
func GetPushSubscriptionsByDevice(deviceName, authToken string) ([]StoredPushSubscription, error) {
	key := fmt.Sprintf("push_subscription:%s:%s", deviceName, authToken)

	res, err := rjh.JSONGet(key, ".")
	if err != nil {
		return nil, fmt.Errorf("error getting push subscription: %v", err)
	}

	if res == nil {
		return []StoredPushSubscription{}, nil
	}

	var sub StoredPushSubscription
	if err := json.Unmarshal(res.([]byte), &sub); err != nil {
		return nil, fmt.Errorf("error unmarshaling push subscription: %v", err)
	}

	return []StoredPushSubscription{sub}, nil
}

// GetAllPushSubscriptions retrieves all push subscriptions from Redis
func GetAllPushSubscriptions() ([]StoredPushSubscription, error) {
	// Get all push subscription keys
	ctx := context.Background()
	keys, err := rdb.Keys(ctx, "push_subscription:*").Result()
	if err != nil {
		return nil, fmt.Errorf("error getting push subscription keys: %v", err)
	}

	var subscriptions []StoredPushSubscription
	for _, key := range keys {
		res, err := rjh.JSONGet(key, ".")
		if err != nil {
			log.Printf("- Warning: failed to get subscription %s: %v", key, err)
			continue
		}

		if res == nil {
			continue
		}

		var sub StoredPushSubscription
		if err := json.Unmarshal(res.([]byte), &sub); err != nil {
			log.Printf("- Warning: failed to unmarshal subscription %s: %v", key, err)
			continue
		}

		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

// RemovePushSubscription removes a push subscription by endpoint
func RemovePushSubscription(endpoint string) error {
	// Find and remove the subscription
	ctx := context.Background()
	keys, err := rdb.Keys(ctx, "push_subscription:*").Result()
	if err != nil {
		return fmt.Errorf("error getting push subscription keys: %v", err)
	}

	for _, key := range keys {
		res, err := rjh.JSONGet(key, ".")
		if err != nil {
			continue
		}

		if res == nil {
			continue
		}

		var sub StoredPushSubscription
		if err := json.Unmarshal(res.([]byte), &sub); err != nil {
			continue
		}

		if sub.Endpoint == endpoint {
			// Remove from Redis
			ctx := context.Background()
			if _, err := rdb.Del(ctx, key).Result(); err != nil {
				return fmt.Errorf("error removing subscription: %v", err)
			}

			// Remove from set
			setKey := fmt.Sprintf("push_subscriptions:%s", sub.AuthToken)
			if _, err := rdb.SRem(ctx, setKey, key).Result(); err != nil {
				log.Printf("- Warning: failed to remove subscription from set: %v", err)
			}

			log.Printf("- Removed push subscription for %s", sub.DeviceName)
			return nil
		}
	}

	return fmt.Errorf("subscription not found")
}

// HandlePushSubscribe handles the push subscription request from the PWA
func HandlePushSubscribe(ctx *fasthttp.RequestCtx) {
	var req PushSubscribeRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		log.Printf("- Error unmarshalling push subscription: %v", err)
		ctx.SetStatusCode(400)
		ctx.SetContentType(jsonMime)
		ctx.WriteString(`{"error":"Invalid subscription data"}`)
		return
	}

	// Validate required fields
	if req.Endpoint == "" || req.Keys == nil {
		ctx.SetStatusCode(400)
		ctx.SetContentType(jsonMime)
		ctx.WriteString(`{"error":"Missing required fields"}`)
		return
	}

	// Extract keys
	p256dh, ok1 := req.Keys["p256dh"].(string)
	auth, ok2 := req.Keys["auth"].(string)

	if !ok1 || !ok2 {
		ctx.SetStatusCode(400)
		ctx.SetContentType(jsonMime)
		ctx.WriteString(`{"error":"Invalid keys format"}`)
		return
	}

	// Get auth token from headers
	authToken := string(ctx.Request.Header.Peek("Auth"))
	deviceName := string(ctx.Request.Header.Peek("Device"))

	if deviceName == "" {
		deviceName = req.DeviceName
	}
	if deviceName == "" {
		deviceName = "Unknown Device"
	}

	// Store subscription
	sub := StoredPushSubscription{
		Endpoint:   req.Endpoint,
		P256dhKey:  p256dh,
		AuthKey:    auth,
		DeviceName: deviceName,
		AuthToken:  authToken,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	if err := SavePushSubscription(sub); err != nil {
		log.Printf("- Error saving push subscription: %v", err)
		ctx.SetStatusCode(500)
		ctx.SetContentType(jsonMime)
		ctx.WriteString(`{"error":"Failed to save subscription"}`)
		return
	}

	log.Printf("- Push subscription registered for device: %s", deviceName)

	ctx.SetStatusCode(200)
	ctx.SetContentType(jsonMime)
	response := map[string]interface{}{
		"success": true,
		"message": "Push subscription registered successfully",
		"device":  deviceName,
	}

	responseJSON, _ := json.Marshal(response)
	ctx.Write(responseJSON)
}

// HandlePushUnsubscribe handles the push unsubscription request
func HandlePushUnsubscribe(ctx *fasthttp.RequestCtx) {
	var req struct {
		Endpoint   string `json:"endpoint"`
		DeviceName string `json:"deviceName"`
	}

	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		ctx.SetStatusCode(400)
		ctx.SetContentType(jsonMime)
		ctx.WriteString(`{"error":"Invalid request data"}`)
		return
	}

	if err := RemovePushSubscription(req.Endpoint); err != nil {
		log.Printf("- Error removing push subscription: %v", err)
		// Don't return error to client if subscription doesn't exist
		if err.Error() != "subscription not found" {
			ctx.SetStatusCode(500)
			ctx.SetContentType(jsonMime)
			ctx.WriteString(`{"error":"Failed to remove subscription"}`)
			return
		}
	}

	ctx.SetStatusCode(200)
	ctx.SetContentType(jsonMime)
	response := map[string]interface{}{
		"success": true,
		"message": "Push subscription removed successfully",
		"device":  req.DeviceName,
	}

	responseJSON, _ := json.Marshal(response)
	ctx.Write(responseJSON)
}

// HandlePushTest sends a test push notification
func HandlePushTest(ctx *fasthttp.RequestCtx) {
	deviceName := string(ctx.QueryArgs().Peek("device"))
	authToken := string(ctx.Request.Header.Peek("Auth"))

	if deviceName == "" {
		ctx.SetStatusCode(400)
		ctx.SetContentType(jsonMime)
		ctx.WriteString(`{"error":"Device parameter required"}`)
		return
	}

	message := PushMessage{
		Type:       "test",
		DeviceName: deviceName,
		Timestamp:  time.Now().Unix(),
	}

	subscriptions, err := GetPushSubscriptionsByDevice(deviceName, authToken)
	if err != nil || len(subscriptions) == 0 {
		ctx.SetStatusCode(404)
		ctx.SetContentType(jsonMime)
		ctx.WriteString(`{"error":"No subscription found for device"}`)
		return
	}

	sub := subscriptions[0]
	webpushSub := &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			P256dh: sub.P256dhKey,
			Auth:   sub.AuthKey,
		},
	}

	if err := SendPushNotification(webpushSub, message); err != nil {
		ctx.SetStatusCode(500)
		ctx.SetContentType(jsonMime)
		response := map[string]interface{}{
			"error": fmt.Sprintf("Failed to send push notification: %v", err),
		}
		responseJSON, _ := json.Marshal(response)
		ctx.Write(responseJSON)
		return
	}

	ctx.SetStatusCode(200)
	ctx.SetContentType(jsonMime)
	response := map[string]interface{}{
		"success": true,
		"message": "Test push notification sent",
		"device":  deviceName,
	}

	responseJSON, _ := json.Marshal(response)
	ctx.Write(responseJSON)
}
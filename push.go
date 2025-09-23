package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/SherClockHolmes/webpush-go"
)

// PushMessage represents the payload sent to the PWA via push notification
type PushMessage struct {
	Type       string `json:"type"`
	DeviceName string `json:"deviceName,omitempty"`
	Timestamp  int64  `json:"timestamp,omitempty"`
}

// SendPushNotification sends a push notification to a specific subscription
func SendPushNotification(subscription *webpush.Subscription, message PushMessage) error {
	if vapidPrivateKey == "" || vapidPublicKey == "" || vapidSubject == "" {
		return fmt.Errorf("VAPID keys not configured")
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling push message: %v", err)
	}

	options := &webpush.Options{
		Subscriber:      vapidSubject,
		VAPIDPublicKey:  vapidPublicKey,
		VAPIDPrivateKey: vapidPrivateKey,
		TTL:             30, // 30 seconds
	}

	resp, err := webpush.SendNotification(messageBytes, subscription, options)
	if err != nil {
		return fmt.Errorf("error sending push notification: %v", err)
	}
	defer resp.Body.Close()

	// Handle different response status codes
	if resp.StatusCode == 410 {
		// Subscription is no longer valid, remove it
		log.Printf("- Push subscription expired, removing: %s", subscription.Endpoint)
		if err := RemovePushSubscription(subscription.Endpoint); err != nil {
			log.Printf("- Error removing expired subscription: %v", err)
		}
		return fmt.Errorf("subscription expired and removed")
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("push notification failed with status: %d", resp.StatusCode)
	}

	return nil
}

// SendHeartbeatCheckToDevice sends a heartbeat check request to a specific device
func SendHeartbeatCheckToDevice(deviceName, authToken string) (int, error) {
	subscriptions, err := GetPushSubscriptionsByDevice(deviceName, authToken)
	if err != nil {
		return 0, fmt.Errorf("error getting subscriptions for device %s: %v", deviceName, err)
	}

	if len(subscriptions) == 0 {
		return 0, fmt.Errorf("no push subscriptions found for device %s", deviceName)
	}

	message := PushMessage{
		Type:       "heartbeat-request",
		DeviceName: deviceName,
		Timestamp:  getCurrentTimestamp(),
	}

	sentCount := 0
	for _, storedSub := range subscriptions {
		webpushSub := &webpush.Subscription{
			Endpoint: storedSub.Endpoint,
			Keys: webpush.Keys{
				P256dh: storedSub.P256dhKey,
				Auth:   storedSub.AuthKey,
			},
		}

		if err := SendPushNotification(webpushSub, message); err != nil {
			log.Printf("- Error sending push to %s: %v", storedSub.Endpoint, err)
		} else {
			sentCount++
			log.Printf("- Sent heartbeat check to %s (%s)", deviceName, storedSub.Endpoint)
		}
	}

	return sentCount, nil
}

// SendHeartbeatCheckToAllDevices sends heartbeat check requests to all subscribed devices
func SendHeartbeatCheckToAllDevices() (int, error) {
	subscriptions, err := GetAllPushSubscriptions()
	if err != nil {
		return 0, fmt.Errorf("error getting all subscriptions: %v", err)
	}

	if len(subscriptions) == 0 {
		return 0, fmt.Errorf("no push subscriptions found")
	}

	sentCount := 0
	for _, storedSub := range subscriptions {
		message := PushMessage{
			Type:       "heartbeat-request",
			DeviceName: storedSub.DeviceName,
			Timestamp:  getCurrentTimestamp(),
		}

		webpushSub := &webpush.Subscription{
			Endpoint: storedSub.Endpoint,
			Keys: webpush.Keys{
				P256dh: storedSub.P256dhKey,
				Auth:   storedSub.AuthKey,
			},
		}

		if err := SendPushNotification(webpushSub, message); err != nil {
			log.Printf("- Error sending push to %s (%s): %v", storedSub.DeviceName, storedSub.Endpoint, err)
		} else {
			sentCount++
			log.Printf("- Sent heartbeat check to %s (%s)", storedSub.DeviceName, storedSub.Endpoint)
		}
	}

	return sentCount, nil
}

// getCurrentTimestamp returns the current Unix timestamp
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// ValidateVAPIDConfiguration checks if VAPID keys are properly configured
func ValidateVAPIDConfiguration() error {
	if vapidPrivateKey == "" {
		return fmt.Errorf("VAPID_PRIVATE_KEY not set in environment")
	}
	if vapidPublicKey == "" {
		return fmt.Errorf("VAPID_PUBLIC_KEY not set in environment")
	}
	if vapidSubject == "" {
		return fmt.Errorf("VAPID_SUBJECT not set in environment")
	}
	return nil
}

// LogVAPIDConfiguration logs the VAPID configuration status (without exposing private key)
func LogVAPIDConfiguration() {
	if vapidPrivateKey != "" && vapidPublicKey != "" && vapidSubject != "" {
		log.Printf("- VAPID configuration loaded (Subject: %s)", vapidSubject)
	} else {
		log.Printf("- VAPID configuration incomplete - push notifications disabled")
		if vapidPrivateKey == "" {
			log.Printf("  Missing: VAPID_PRIVATE_KEY")
		}
		if vapidPublicKey == "" {
			log.Printf("  Missing: VAPID_PUBLIC_KEY")
		}
		if vapidSubject == "" {
			log.Printf("  Missing: VAPID_SUBJECT")
		}
	}
}

// CleanupExpiredSubscriptions removes subscriptions that are no longer valid
func CleanupExpiredSubscriptions() error {
	subscriptions, err := GetAllPushSubscriptions()
	if err != nil {
		return err
	}

	removedCount := 0
	for _, storedSub := range subscriptions {
		// Try to send a test notification to check if subscription is still valid
		testMessage := PushMessage{
			Type:      "test",
			Timestamp: getCurrentTimestamp(),
		}

		webpushSub := &webpush.Subscription{
			Endpoint: storedSub.Endpoint,
			Keys: webpush.Keys{
				P256dh: storedSub.P256dhKey,
				Auth:   storedSub.AuthKey,
			},
		}

		if err := SendPushNotification(webpushSub, testMessage); err != nil {
			if err.Error() == "subscription expired and removed" {
				removedCount++
			}
		}
	}

	if removedCount > 0 {
		log.Printf("- Cleaned up %d expired push subscriptions", removedCount)
	}

	return nil
}

// DevicePushStatus tracks push notification status for a device
type DevicePushStatus struct {
	DeviceName        string `json:"device_name"`
	LastPushSent      int64  `json:"last_push_sent"`
	PushCount         int64  `json:"push_count"`
	LastHeartbeatTime int64  `json:"last_heartbeat_time"`
}

// StartPushMonitoring starts the background service that monitors device inactivity
// and sends push notifications when devices haven't sent heartbeats for too long
func StartPushMonitoring() {
	if ValidateVAPIDConfiguration() != nil {
		log.Printf("- Push monitoring disabled: VAPID not configured")
		return
	}

	log.Printf("- Starting push monitoring (check every %v, threshold %v)", 
		pushCheckInterval, pushInactivityThreshold)

	ticker := time.NewTicker(pushCheckInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				checkInactiveDevicesAndPush()
			}
		}
	}()
}

// checkInactiveDevicesAndPush checks all devices for inactivity and sends push notifications as needed
func checkInactiveDevicesAndPush() {
	// Get all devices with push subscriptions
	subscriptions, err := GetAllPushSubscriptions()
	if err != nil {
		log.Printf("- Error getting push subscriptions for monitoring: %v", err)
		return
	}

	if len(subscriptions) == 0 {
		return // No devices to monitor
	}

	currentTime := time.Now().Unix()
	devicesChecked := 0
	pushesTriggered := 0

	// Group subscriptions by device to avoid duplicate checks
	deviceSubscriptions := make(map[string][]StoredPushSubscription)
	for _, sub := range subscriptions {
		deviceKey := sub.DeviceName + "|" + sub.AuthToken
		deviceSubscriptions[deviceKey] = append(deviceSubscriptions[deviceKey], sub)
	}

	for deviceKey := range deviceSubscriptions {
		parts := strings.SplitN(deviceKey, "|", 2)
		if len(parts) != 2 {
			continue
		}
		deviceName := parts[0]
		authToken := parts[1]

		// Get the device's last heartbeat time
		lastHeartbeat := getLastHeartbeatForDevice(deviceName)
		if lastHeartbeat == nil {
			continue // No heartbeat data for this device
		}

		// Calculate time since last heartbeat
		timeSinceHeartbeat := time.Duration(currentTime-lastHeartbeat.Timestamp) * time.Second

		// Check if device is inactive beyond threshold
		if timeSinceHeartbeat > pushInactivityThreshold {
			// Check if we've already sent a push recently to avoid spam
			if shouldSendPushNotification(deviceName, authToken, lastHeartbeat.Timestamp) {
				log.Printf("- Device %s inactive for %v (threshold: %v), sending push notification", 
					deviceName, timeSinceHeartbeat.Round(time.Minute), pushInactivityThreshold)

				// Send push notification to wake up the device
				sentCount, err := SendHeartbeatCheckToDevice(deviceName, authToken)
				if err != nil {
					log.Printf("- Error sending push to inactive device %s: %v", deviceName, err)
				} else {
					pushesTriggered += sentCount
					// Record that we sent a push for this device
					recordPushSent(deviceName, authToken, currentTime)
				}
			}
		}

		devicesChecked++
	}

	if pushesTriggered > 0 {
		log.Printf("- Push monitoring: checked %d devices, triggered %d push notifications", 
			devicesChecked, pushesTriggered)
	}
}

// getLastHeartbeatForDevice gets the last heartbeat for a specific device
func getLastHeartbeatForDevice(deviceName string) *HeartbeatBeat {
	// Search through heartbeatDevices for the specific device
	if heartbeatDevices == nil {
		return nil
	}

	for _, device := range *heartbeatDevices {
		if device.DeviceName == deviceName {
			return &device.LastBeat
		}
	}

	return nil
}

// shouldSendPushNotification determines if we should send a push notification to avoid spam
func shouldSendPushNotification(deviceName, authToken string, lastHeartbeatTime int64) bool {
	status := getPushStatus(deviceName, authToken)
	currentTime := time.Now().Unix()

	// Don't send push if we sent one recently (within 15 minutes)
	timeSinceLastPush := time.Duration(currentTime-status.LastPushSent) * time.Second
	if timeSinceLastPush < 15*time.Minute {
		return false
	}

	// Don't send push if heartbeat time hasn't changed (device might be offline)
	if status.LastHeartbeatTime == lastHeartbeatTime && status.PushCount > 0 {
		// Only send another push if it's been a long time since the last one
		return timeSinceLastPush > 4*time.Hour
	}

	return true
}

// getPushStatus retrieves the push status for a device
func getPushStatus(deviceName, authToken string) DevicePushStatus {
	statusKey := fmt.Sprintf("push_status:%s:%s", deviceName, authToken)
	
	res, err := rjh.JSONGet(statusKey, ".")
	if err != nil {
		// Return default status if not found
		return DevicePushStatus{
			DeviceName:        deviceName,
			LastPushSent:      0,
			PushCount:         0,
			LastHeartbeatTime: 0,
		}
	}

	var status DevicePushStatus
	if err = json.Unmarshal(res.([]byte), &status); err != nil {
		log.Printf("- Error unmarshaling push status for %s: %v", deviceName, err)
		return DevicePushStatus{
			DeviceName:        deviceName,
			LastPushSent:      0,
			PushCount:         0,
			LastHeartbeatTime: 0,
		}
	}

	return status
}

// recordPushSent records that a push notification was sent to a device
func recordPushSent(deviceName, authToken string, timestamp int64) {
	status := getPushStatus(deviceName, authToken)
	status.LastPushSent = timestamp
	status.PushCount++
	// Don't update LastHeartbeatTime here - keep the actual last heartbeat time

	statusKey := fmt.Sprintf("push_status:%s:%s", deviceName, authToken)
	if _, err := rjh.JSONSet(statusKey, ".", status); err != nil {
		log.Printf("- Error recording push status for %s: %v", deviceName, err)
	}
}

// GetPushMonitoringStats returns statistics about the push monitoring system
func GetPushMonitoringStats() map[string]interface{} {
	subscriptions, err := GetAllPushSubscriptions()
	if err != nil {
		return map[string]interface{}{
			"error": "Failed to get subscriptions",
		}
	}

	currentTime := time.Now().Unix()
	activeDevices := 0
	inactiveDevices := 0
	totalPushes := int64(0)

	// Group by device
	deviceMap := make(map[string]bool)
	for _, sub := range subscriptions {
		deviceKey := sub.DeviceName + "|" + sub.AuthToken
		if _, exists := deviceMap[deviceKey]; !exists {
			deviceMap[deviceKey] = true

			parts := strings.SplitN(deviceKey, "|", 2)
			if len(parts) == 2 {
				deviceName := parts[0]
				authToken := parts[1]

				lastHeartbeat := getLastHeartbeatForDevice(deviceName)
				if lastHeartbeat != nil {
					timeSinceHeartbeat := time.Duration(currentTime-lastHeartbeat.Timestamp) * time.Second
					if timeSinceHeartbeat > pushInactivityThreshold {
						inactiveDevices++
					} else {
						activeDevices++
					}
				}

				status := getPushStatus(deviceName, authToken)
				totalPushes += status.PushCount
			}
		}
	}

	return map[string]interface{}{
		"total_subscriptions":     len(subscriptions),
		"unique_devices":          len(deviceMap),
		"active_devices":          activeDevices,
		"inactive_devices":        inactiveDevices,
		"total_pushes_sent":       totalPushes,
		"check_interval":          pushCheckInterval.String(),
		"inactivity_threshold":    pushInactivityThreshold.String(),
		"vapid_configured":        ValidateVAPIDConfiguration() == nil,
	}
}
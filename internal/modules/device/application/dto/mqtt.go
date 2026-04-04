package devicedto

type MQTTAuthInput struct {
	ClientID string `json:"clientid"`
	CN       string `json:"cn"`
}

type MQTTACLInput struct {
	ClientID string `json:"clientid"`
	CN       string `json:"cn"`
	Topic    string `json:"topic"`
	Action   string `json:"action"`
}

// identity returns the CN from mTLS if present, falling back to clientid (plain MQTT / dev).
func (i MQTTAuthInput) Identity() string {
	if i.CN != "" {
		return i.CN
	}
	return i.ClientID
}

func (i MQTTACLInput) Identity() string {
	if i.CN != "" {
		return i.CN
	}
	return i.ClientID
}

type MQTTResult struct {
	Result string `json:"result"` // "allow" | "deny"
}

var MQTTAllow = MQTTResult{Result: "allow"}
var MQTTDeny = MQTTResult{Result: "deny"}

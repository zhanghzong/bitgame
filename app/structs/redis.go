package structs

// Redis channel communication struct defined
type RedisChannel struct {
	Type     string      `json:"type"`
	Hostname string      `json:"hostname"`
	Users    []string    `json:"users"`
	Data     interface{} `json:"data"`
}

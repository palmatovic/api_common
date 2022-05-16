package api_common

type ErrorData struct {
	Error Error `json:"error"`
}

type Error struct {
	ErrorCode string `json:"error_code,omitempty"`
	Reason    string `json:"reason,omitempty"`
	Detail    string `json:"detail,omitempty"`
}

type MonitorResponse struct {
	Status bool        `json:"status"`
	Data   MonitorData `json:"data"`
}
type MonitorRequest struct {
	Data MonitorData `json:"data"`
}

type MonitorData struct {
	Monitor Monitor `json:"monitor,omitempty"`
}

type NotificationResponse struct {
	Status bool             `json:"status"`
	Data   NotificationData `json:"data"`
}
type NotificationRequest struct {
	Data NotificationData `json:"data"`
}

type NotificationData struct {
	Notification Notification `json:"notification,omitempty"`
}

type Monitor struct {
	Response   string `json:"response,omitempty"`
	Uuid       string `json:"uuid,omitempty"`
	Source     string `json:"source,omitempty"`
	SourceType string `json:"source_type,omitempty"`
	Success    bool   `json:"success,omitempty"`
	Status     int    `json:"status,omitempty"`
	Endpoint   string `json:"endpoint,omitempty"`
}

type Notification struct {
	NotifyID     string   `json:"notify_id,omitempty"`
	UserID       []string `json:"user_id,omitempty"`
	NotifyTypeID string   `json:"notify_type_id,omitempty"`
	SourceType   string   `json:"source_type,omitempty"`
	Message      string   `json:"message,omitempty"`
	Link         string   `json:"link,omitempty"`
}

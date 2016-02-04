package apnsapi

const (
	DevelopmentServer = "https://api.development.push.apple.com"
	ProductionServer  = "https://api.push.apple.com"
)

type Header struct {
	ApnsID         string
	ApnsExpiration string
	ApnsPriority   string
	ApnsTopic      string
}

type Response struct {
	ApnsID     string
	StatusCode int
}

type ErrorResponse struct {
	Reason    string `json:"reason"`
	Timestamp int    `json:"timestamp"`
}

func (e *ErrorResponse) Error() string {
	return e.Reason
}

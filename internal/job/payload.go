package job

type SendEmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type ResizeImagePayload struct {
	URL   string `json:"url"`
	Width int    `json:"width"`
}

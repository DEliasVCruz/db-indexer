package handlers

type FileUploaded struct {
	Uploaded bool   `json:"uploaded,omitempty"`
	State    string `json:"state"`
}

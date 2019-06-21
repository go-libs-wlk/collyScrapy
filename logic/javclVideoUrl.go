package logic

type VideoUrl struct {
	Success bool `json:"success"`
	Data []DataObject
}

type DataObject struct {
	File string `json:"file"`
	Label string `json:"label"`
	Type string `json:"type"`
}

type VideoDownload struct {
	Href string
	Err int
	VideoNum string
	VideoLabel string
}
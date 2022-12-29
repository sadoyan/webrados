package web

type FileInfo struct {
	Name     string `json:"name"`
	Size     string `json:"size"`
	Pool     string `json:"pool"`
	Segments string `json:"segments"`
}

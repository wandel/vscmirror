package marketplace

type ProductInfo struct {
	Name               string       `json:"name"`
	Version            string       `json:"version"`
	ProductVersion     string       `json:"productVersion"`
	Url                string       `json:"url"`
	Hash               string       `json:"hash"`
	Build              string       `json:"build"`
	Timestamp          int          `json:"timestamp"`
	SHA256Hash         string       `json:"sha256hash"`
	SupportsFastUpdate bool         `json:"supportsFastUpdate"`
	Platform           PlatformInfo `json:"platform"`
}

type PlatformInfo struct {
	OperatingSystem string `json:"os"`
	PrettyName      string `json:"prettyname"`
}

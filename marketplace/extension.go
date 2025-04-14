package marketplace

import "strings"

type Extension struct {
	Publisher        Publisher            `json:"publisher"`
	ExtensionId      string               `json:"extensionId"`
	ExtensionName    string               `json:"extensionName"`
	DisplayName      string               `json:"displayName"`
	Flags            string               `json:"flags"`
	LastUpdated      string               `json:"lastUpdated"`
	PublishedDate    string               `json:"publishedDate"`
	ReleaseDate      string               `json:"releaseDate"`
	ShortDescription string               `json:"shortDescription"`
	Versions         []ExtensionVersion   `json:"versions"`
	Categories       []string             `json:"categories"`
	Tags             []string             `json:"tags"`
	Statistics       []ExtensionStatistic `json:"statistics"`
	DeploymentType   int                  `json:"deploymentType"`
	// This is in the vscodeoffline code, but I haven't seen it from the vscode marketplace yet
	// Recommended      bool                 `json:"recommended"`
}

func (e Extension) GetStatistic(name string) float64 {
	name = strings.ToLower(name)
	for _, stat := range e.Statistics {
		if strings.ToLower(stat.StatisticName) == name {
			return stat.Value
		}
	}

	return 0
}

type Publisher struct {
	PublisherId      string `json:"publisherId"`
	PublisherName    string `json:"publisherName"`
	DisplayName      string `json:"displayName"`
	Flags            string `json:"flags"`
	Domain           string `json:"domain"`
	IsDomainVerified bool   `json:"isDomainVerified"`
}

type ExtensionStatistic struct {
	StatisticName string  `json:"statisticName"`
	Value         float64 `json:"value"`
}

type ExtensionVersion struct {
	Version          string              `json:"version"`
	TargetPlatform   string              `json:"targetPlatform,omitempty"`
	Flags            string              `json:"flags"`
	LastUpdated      string              `json:"lastUpdated"`
	Files            []ExtensionFile     `json:"files"`
	Properties       []ExtensionProperty `json:"properties"`
	AssetURI         string              `json:"assetUri"`
	FallbackAssetURI string              `json:"fallbackAssetUri"`
}

type ExtensionProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ExtensionFile struct {
	AssetType string `json:"assetType"`
	Source    string `json:"source"`
}

package common

import (
	"path"
	"strings"

	"github.com/wandel/vscmirror/marketplace"
)

var PLATFORMS = []string{"win32", "linux", "linux-deb", "linux-rpm", "darwin", "linux-snap", "server-linux", "server-linux-legacy", "cli-alpine"}
var ARCHITECTURES = []string{"", "x64", "ia32"}
var BUILD_TYPES = []string{"", "archive", "user"}
var QUALITY = []string{"stable", "insider"}

type ProductInfoEx struct {
	marketplace.ProductInfo
	Identity         string `json:"identity"`
	Platform         string `json:"platform"`
	Architecture     string `json:"architecture"`
	BuildType        string `json:"buildtype"`
	Quality          string `json:"quality"`
	CheckedForUpdate bool   `json:"checkedForUpdate"`
	UpdateUrl        string `json:"updateUrl"`
}

func (info ProductInfoEx) GetDownloadUrl() string {
	_, filename := path.Split(info.UpdateUrl)
	ext := path.Ext(filename)
	if ext == ".gz" {
		tmp := strings.TrimSuffix(filename, ".gz")
		ext = ext + path.Ext(tmp)
	}

	return path.Join("installers", info.Identity, info.Quality, "vscode-"+info.Name+ext)
}

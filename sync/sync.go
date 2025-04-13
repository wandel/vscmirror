package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"

	marketplace "github.com/wandel/vscmirror/client"
	"github.com/wandel/vscmirror/common"
)

var ARTIFACTS = os.DirFS("D:\\vscmirror")

func DownloadInstallers(ctx context.Context) error {
	var response struct {
		Products []marketplace.ProductInfo `json:"products"`
	}

	url := "https://code.visualstudio.com/sha?build"
	if err := common.DownloadJson(url, &response); err != nil {
		return fmt.Errorf("failed to get list of latest installers: %w", err)
	}

	for _, installer := range response.Products {
		_, name := path.Split(installer.Url)
		slog.Info(name, "platform", installer.Platform.OperatingSystem, "build", installer.Build)
	}

	// https://code.visualstudio.com/sha?build=stable
	// https://code.visualstudio.com/sha?build=insiders
	// if err := os.MkdirAll("installers", 0644); err != nil {
	// 	return fmt.Errorf("failed to create installers directory: %w", err)
	// }

	return nil
}

func DownloadInstaller(ctx context.Context, platform, channel, commit string) error {
	path := filepath.Join("installers", platform, channel)
	if err := os.MkdirAll(path, 0644); err != nil {
		return fmt.Errorf("failed to create '%s': %w", path, err)
	}
	url := fmt.Sprintf("https://update.code.visualstudio.com/api/update/%s/%s/%s", platform, channel, commit)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download installer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := common.WriteFile(fmt.Sprintf("installers/%s-%s.exe", platform, channel), resp.Body, 0644); err != nil {
		return fmt.Errorf("failed to write installer: %w", err)
	}

	return nil
}

func DownloadExtensions(ctx context.Context) error {
	if err := os.MkdirAll("extensions", 0644); err != nil {
		return fmt.Errorf("failed to create installers directory: %w", err)
	}

	client := marketplace.Client{
		HttpClient: http.DefaultClient,
		Version:    "1.99.2",
	}

	var extensions []marketplace.Extension
	for extension := range client.GetAllExtensions(ctx) {
		extensions = append(extensions, extension)
	}

	if data, err := json.Marshal(extensions); err != nil {
		return fmt.Errorf("failed to marshal extensions: %w", err)
	} else if err := os.WriteFile("extensions.json", data, 0644); err != nil {
		return fmt.Errorf("failed to write extensions: %w", err)
	} else {
		slog.Info("downloaded extensions", "count", len(extensions))
	}

	return nil
}

// func DownloadMarketplaceQuery(ctx context.Context) error {
// 	url := "https://marketplace.visualstudio.com/_apis/public/gallery/extensionquery"
// 	return nil
// }

func DownloadMalicious(ctx context.Context) error {
	url := "https://main.vscode-cdn.net/extensions/marketplace.json"
	path := "extensions/marketplace.json"
	if err := common.Download(url, path); err != nil {
		return fmt.Errorf("failed to download marketplace.json")
	}

	return nil
}

// func DownloadRecomendations(ctx context.Context) error {
// 	url := "https://az764295.vo.msecnd.net/extensions/workspaceRecommendations.json.gz"
// 	path := "extensions/workspaceRecommendations.json.gz"
// 	if err := common.Download(url, path); err != nil {
// 		return fmt.Errorf("failed to download workspaceRecommendations.json.gz")
// 	}

// 	return nil
// }

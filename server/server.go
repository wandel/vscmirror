package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/wandel/vscmirror/client"
	"github.com/wandel/vscmirror/common"
)

var ARTIFACTS = os.DirFS("D:\\vscmirror")
var DOMAIN = "https://vscode.cdn.local/"

func NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	// mux.HandleFunc("GET /", IndexHandler)
	// mux.HandleFunc("GET /browse", BrowseHandler)

	// Installer Auto Update
	mux.HandleFunc("GET /api/update/{platform}/{quality}/{commit}", CheckInstallerHandler)
	mux.HandleFunc("GET /{commit}/{platform}/{quality}", DownloadInstallerHandler)

	// Extension Marketplace
	mux.HandleFunc("GET /extensions/marketplace.json", MaliciousHandler)
	mux.HandleFunc("GET /extensions/chat.json", ChatHandler)
	mux.HandleFunc("POST /_apis/public/gallery/extensionquery", GalleryQueryHandler)
	mux.HandleFunc("GET /_apis/public/gallery/vscode/{publisher}/{extension}/latest", GalleryLatestHandler)
	// mux.HandleFunc("GET /_gallery/{publisher}/{extension}/latest", AltGalleryLatestHandler)
	mux.HandleFunc("GET vscode.cdn.local/extensions/", DownloadExtensionHandler)
	// https://vscode.cdn.local/extensions/dongli/python-preview/0.0.4/Microsoft.VisualStudio.Services.Icons.Default?targetPlatform=universal
	// Handles the
	mux.HandleFunc("OPTIONS /", OptionsHandler)

	return mux
}

func OptionsHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("request", "handler", "OptionsHandler", "remote", r.RemoteAddr, "url", r.URL.String())
	w.Header().Set("access-control-allow-origin", "*")
	// w.Header().Set("strict-transport-security", "max-age=31536000; includeSubDomains")
	w.Header().Set("access-control-allow-methods", "OPTIONS,GET,POST,PATCH,PUT,DELETE")
	w.Header().Set("access-control-max-age", "3600")
	w.Header().Set("access-control-allow-headers", "content-type,vscode-sessionid,x-market-client-id,x-market-user-id,authorization")
	w.WriteHeader(http.StatusOK)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("request", "handler", "IndexHandler", "remote", r.RemoteAddr, "url", r.URL.String())
	w.WriteHeader(http.StatusOK)
}

func BrowseHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("request", "handler", "BrowseHandler", "remote", r.RemoteAddr, "url", r.URL.String())
	w.WriteHeader(http.StatusOK)
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("access-control-allow-origin", "*")
	http.NotFound(w, r)
}

func DownloadExtensionHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("request", "handler", "DownloadExtensionHandler", "remote", r.RemoteAddr, "url", r.URL.String())
	w.Header().Set("access-control-allow-origin", "*")
	http.ServeFileFS(w, r, ARTIFACTS, r.URL.Path)
}

func CheckInstallerHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("request", "handler", "CheckInstallerHandler", "remote", r.RemoteAddr, "url", r.URL.String())
	w.Header().Set("access-control-allow-origin", "*")
	platform := r.PathValue("platform")
	quality := r.PathValue("quality")
	commit := r.PathValue("commit")

	slog.Info("installer update check", "platform", platform, "quality", quality, "commit", commit)

	path := path.Join("installers", platform, quality, "latest.json")
	var installer common.ProductInfoEx
	if err := common.LoadJsonFS(ARTIFACTS, path, &installer); err != nil {
		slog.Error("failed to load latest installer metadata", "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// already at the latest version
	if installer.Version == commit {
		slog.Debug("no update found", "platform", platform, "quality", quality, "commit", commit)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	installer.Url = DOMAIN + installer.GetDownloadUrl()
	data, err := json.Marshal(installer)
	if err != nil {
		slog.Error("failed to encode", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		slog.Error("check installer failed to send body", "file", path, "error", err)
	}
}

func DownloadInstallerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("access-control-allow-origin", "*")
	slog.Info("request", "handler", "DownloadInstallerHandler", "remote", r.RemoteAddr, "url", r.URL.String())
	platform := r.PathValue("platform")
	quality := r.PathValue("quality")
	commit := strings.TrimPrefix(r.PathValue("commit"), "commit:")

	slog.Info("installer download", "platform", platform, "quality", quality, "commit", commit)

	filepath := path.Join("installers", platform, quality, commit+".json")
	var installer common.ProductInfoEx
	if err := common.LoadJsonFS(ARTIFACTS, filepath, &installer); err != nil {
		slog.Error("failed to load installer metadata", "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	filepath = installer.GetDownloadUrl()
	slog.Info("sending file", "path", filepath)
	http.ServeFileFS(w, r, ARTIFACTS, filepath)
}

func RecommendationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("access-control-allow-origin", "*")
	slog.Info("request", "handler", "RecommendationHandler", "remote", r.RemoteAddr, "url", r.URL.String())
	http.ServeFileFS(w, r, ARTIFACTS, "recommendations.json")
}

func MaliciousHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("access-control-allow-origin", "*")
	slog.Info("request", "handler", "MaliciousHandler", "remote", r.RemoteAddr, "url", r.URL.String())
	http.ServeFileFS(w, r, ARTIFACTS, "malicious.json")
}

func LoadExtensions(dst *[]client.Extension) error {
	matches, err := fs.Glob(ARTIFACTS, "extensions/*/latest.json")
	if err != nil {
		return fmt.Errorf("failed to glob extensions: %w", err)
	}

	var extensions []client.Extension
	for _, match := range matches {
		var extension client.Extension
		if err := common.LoadJsonFS(ARTIFACTS, match, &extension); err != nil {
			slog.Error("failed to load extension metadata", "path", match, "error", err)
			continue
		}
		extensions = append(extensions, extension)
	}
	*dst = extensions
	return nil
}

func GalleryLatestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("access-control-allow-origin", "*")
	slog.Info("request", "handler", "GalleryLatestHandler", "remote", r.RemoteAddr, "url", r.URL.String())
	identity := r.PathValue("publisher") + "." + r.PathValue("extension")
	filepath := path.Join("extensions", identity, "latest.json")
	f, err := ARTIFACTS.Open(filepath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer f.Close()

	var extension client.Extension
	if err := json.NewDecoder(f).Decode(&extension); err != nil {
		slog.Error("failed to decode json", "path", filepath, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	for j, version := range extension.Versions {
		assetURI := DOMAIN + path.Join("extensions", extension.Publisher.PublisherName, extension.ExtensionName, version.Version)
		for k, file := range version.Files {
			source := strings.Replace(file.Source, version.AssetURI, assetURI, 1)
			extension.Versions[j].Files[k].Source = source
		}
		extension.Versions[j].FallbackAssetURI = version.AssetURI
		extension.Versions[j].AssetURI = assetURI
	}

	data, err := json.Marshal(extension)
	if err != nil {
		slog.Error("failed to marshal query response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write query response", "error", err)
		return
	}

}

func GalleryQueryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("access-control-allow-origin", "*")
	slog.Info("request", "handler", "GalleryQueryHandler", "remote", r.RemoteAddr, "url", r.URL.String())
	var request client.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.Error("failed to parse query request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(request.Filters) == 0 {
		http.Error(w, "no filters specified", http.StatusBadRequest)
		return
	} else if request.Flags == client.QueryFlagNoneDefined {
		http.Error(w, "no flags specified", http.StatusBadRequest)
		return
	}

	for _, filter := range request.Filters {
		if len(filter.Criteria) == 0 {
			http.Error(w, "no criteria specified on filter", http.StatusBadRequest)
			return
		}
	}

	var extensions []client.Extension
	if err := LoadExtensions(&extensions); err != nil {
		slog.Error("failed to load extensions", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filter := request.Filters[0]

	result := []client.Extension{} // initialize so we dont get a null in the json later
	for _, extension := range extensions {
		if filter.Matches(extension) {
			slog.Info("selected extension", "extension", extension.ExtensionName)
			result = append(result, extension)
		}
	}

	slices.SortFunc(result, filter.Compare)

	var reverse bool = filter.SortOrder == client.SortOrderAscending
	switch filter.SortBy {
	case client.SortByNoneOrRelevance:
	case client.SortByLastUpdatedDate:
		reverse = !reverse
	case client.SortByTitle:
		reverse = !reverse
	case client.SortByPublisherName:
	case client.SortByInstallCount:
		reverse = !reverse
	case client.SortByPublishedDate:
		reverse = !reverse
	case client.SortByAverageRating:
		reverse = !reverse
	case client.SortByWeightedRating:
		reverse = !reverse
	default:
	}

	if !reverse {
		slices.Reverse(result)
	}

	// Change Urls to point to us, mimicking vscodeoffline
	for i, extension := range result {
		for j, version := range extension.Versions {
			identity := extension.Publisher.PublisherName + "." + extension.ExtensionName
			uri := DOMAIN + path.Join("extensions", identity, version.Version, version.TargetPlatform)
			for k, file := range version.Files {
				result[i].Versions[j].Files[k].Source = path.Join(uri, file.AssetType)
			}
			result[i].Versions[j].AssetURI = uri
			result[i].Versions[j].FallbackAssetURI = uri
		}
	}

	// build metadta
	categoriesMap := map[string]int{}
	targetsMap := map[string]int{}
	for _, extension := range result {
		for _, category := range extension.Categories {
			name := strings.ToLower(category)
			categoriesMap[name] += 1
		}

		for _, version := range extension.Versions {
			target := strings.ToLower(version.TargetPlatform)
			targetsMap[target] += 1
		}
	}

	var categories []client.MetadataItem
	for name, count := range categoriesMap {
		categories = append(categories, client.MetadataItem{
			Name:  name,
			Count: count,
		})
	}

	var targets []client.MetadataItem
	for name, count := range targetsMap {
		targets = append(targets, client.MetadataItem{
			Name:  name,
			Count: count,
		})
	}

	totalCount := client.MetadataItem{
		Name:  "TotalCount",
		Count: len(result),
	}

	// pageNumber starts at 1, not 0 so we correct it here.
	result = paginate(result, filter.PageNumber-1, filter.PageSize)
	response := client.QueryResponse{
		Extensions: result,
		ResultMetadata: []client.QueryResultMetadata{
			client.QueryResultMetadata{
				MetadataType: "ResultCount",
				MetadataItems: []client.MetadataItem{
					totalCount,
				},
			},
			client.QueryResultMetadata{
				MetadataType:  "Categories",
				MetadataItems: categories,
			},
			client.QueryResultMetadata{
				MetadataType:  "TargetPlatforms",
				MetadataItems: targets,
			},
		},
	}

	wrapper := struct {
		Results []client.QueryResponse `json:"results"`
	}{
		Results: []client.QueryResponse{response},
	}

	data, err := json.Marshal(wrapper)
	if err != nil {
		slog.Error("failed to marshal query response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write query response", "error", err)
		return
	}

	return
}

func paginate[T any](values []T, page, count int) []T {
	if page < 0 {
		page = 0
	}

	// coping marketplace.visualstudio.com pageSize limit
	if count < 1 {
		count = 1
	} else if count > 1000 {
		count = 1000
	}

	start := page * count
	end := start + count

	if start >= len(values) {
		return []T{}
	} else if end >= len(values) {
		return values[start:]
	} else {
		return values[start:end]
	}
}

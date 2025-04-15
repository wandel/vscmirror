package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/urfave/cli/v3"

	"github.com/wandel/vscmirror/marketplace"
	"github.com/wandel/vscmirror/server"
	"github.com/wandel/vscmirror/sync"
)

func main() {
	app := cli.Command{
		Name:  "vsmarket",
		Usage: "A marketplace for Visual Studio Code extensions",
		Commands: []*cli.Command{
			&cli.Command{
				Name:   "download",
				Usage:  "download all extension data",
				Action: DownloadAction,
				Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
					tmp, _ := signal.NotifyContext(ctx, os.Interrupt)
					return tmp, nil
				},
			},
			&cli.Command{
				Name:   "serve",
				Usage:  "serve available extensions",
				Action: ServeAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "address",
						Value: "127.0.0.1:443",
					},
				},
			},
			&cli.Command{
				Name:   "search",
				Usage:  "search the extension marketplace",
				Action: SearchAction,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:  "id",
						Usage: "filter by extension id",
					},
					&cli.StringSliceFlag{
						Name:  "category",
						Usage: "filter by extension category",
					},
					&cli.StringSliceFlag{
						Name:  "name",
						Usage: "filter by extension name",
					},
					&cli.StringSliceFlag{
						Name:  "target",
						Usage: "filter by extension target",
					},
					&cli.BoolFlag{
						Name:  "featured",
						Usage: "filter by extension featured",
					},
					&cli.StringSliceFlag{
						Name:  "tag",
						Usage: "filter by extension tag",
					},
					&cli.StringSliceFlag{
						Name:  "text",
						Usage: "filter by search text",
					},
					&cli.BoolFlag{
						Name:  "exclude",
						Usage: "exclude with flags",
					},
					&cli.BoolFlag{
						Name:  "unpublished",
						Usage: "include unpublished extensions",
					},
					&cli.IntFlag{
						Name:  "page",
						Usage: "which 'page' of results to return",
						Value: 1,
					},
					&cli.IntFlag{
						Name:  "count",
						Usage: "number of entries per 'page'",
						Value: 100,
					},
				},
			},
		},
	}

	ctx := context.Background()
	if err := app.Run(ctx, os.Args); err != nil {
		log.Fatalln(err.Error())
	}
}

func DownloadAction(ctx context.Context, cmd *cli.Command) error {
	if err := sync.DownloadInstallers(ctx); err != nil {
		return fmt.Errorf("failed to download malicious extensions: %w", err)
	}

	// if err := sync.DownloadRecomendations(ctx); err != nil {
	// 	return fmt.Errorf("failed to download malicious extensions: %w", err)
	// }

	if err := sync.DownloadMalicious(ctx); err != nil {
		return fmt.Errorf("failed to download malicious extensions: %w", err)
	}

	if err := sync.DownloadExtensions(ctx); err != nil {
		return fmt.Errorf("failed to download malicious extensions: %w", err)
	}

	return nil
}

func ServeAction(ctx context.Context, cmd *cli.Command) error {
	address := cmd.String("address")
	slog.Info("listening", "address", address)

	var extensions []marketplace.Extension
	if err := server.LoadExtensions(&extensions); err != nil {
		return fmt.Errorf("failed to load extensions: %w", err)
	}

	mux := server.NewServeMux()
	return http.ListenAndServeTLS(address, "visualstudio.com.crt", "visualstudio.com.key", mux)
}

func SearchAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Println("Searching the marketplace...")

	client := marketplace.Client{
		HttpClient: http.DefaultClient,
		Version:    "1.99.1",
	}

	criteria := []marketplace.FilterCriteria{}
	ids := cmd.StringSlice("id")
	for _, id := range ids {
		criteria = append(criteria, marketplace.FilterCriteria{
			FilterType: marketplace.FilterTypeId,
			Value:      id,
		})
	}

	names := cmd.StringSlice("name")
	for _, id := range names {
		criteria = append(criteria, marketplace.FilterCriteria{
			FilterType: marketplace.FilterTypeExtensionName,
			Value:      id,
		})
	}

	texts := cmd.StringSlice("text")
	for _, text := range texts {
		criteria = append(criteria, marketplace.FilterCriteria{
			FilterType: marketplace.FilterTypeSearchText,
			Value:      text,
		})
	}

	criteria = append(criteria, marketplace.FilterCriteria{
		FilterType: marketplace.FilterTypeInstallationTarget,
		Value:      "Microsoft.VisualStudio.Code",
	})

	if !cmd.Bool("unpublished") {
		criteria = append(criteria, marketplace.FilterCriteria{
			FilterType: marketplace.FilterTypeExcludeWithFlags,
			Value:      "4096",
		})
	}

	var flags marketplace.QueryFlag
	flags |= marketplace.QueryFlagIncludeFiles
	flags |= marketplace.QueryFlagIncludeCategoryAndTags
	flags |= marketplace.QueryFlagIncludeVersionProperties
	flags |= marketplace.QueryFlagExcludeNonValidated
	flags |= marketplace.QueryFlagIncludeAssetUri
	flags |= marketplace.QueryFlagIncludeStatistics
	flags |= marketplace.QueryFlagIncludeLatestVersionOnly

	page := int(cmd.Int("page"))
	count := int(cmd.Int("count"))
	if count == 0 {
		count = len(ids) + len(names)
	}

	slog.Info("searching marketplace", "count", count, "page", page, "ids", ids, "names", names)
	request := marketplace.QueryRequest{
		Flags:      flags,
		AssetTypes: []string{},
		Filters: []marketplace.QueryFilter{
			marketplace.QueryFilter{
				Criteria:   criteria,
				SortBy:     marketplace.SortByInstallCount,
				SortOrder:  marketplace.SortOrderDescending,
				PageSize:   count,
				PageNumber: page,
			},
		},
	}

	response, err := client.GenericQuery(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to query marketplace: %w", err)
	}

	total := 0
	for _, metadata := range response.ResultMetadata {
		if metadata.MetadataType != "ResultCount" {
			continue
		}

		for _, item := range metadata.MetadataItems {
			if item.Name == "TotalCount" {
				total = item.Count
			}
		}
	}
	slog.Info("search completed", "results", len(response.Extensions), "total", total)

	f, err := os.Create("search.json")
	if err != nil {
		return fmt.Errorf("failed to create search.json: %w", err)
	}
	defer f.Close()

	archive := struct {
		Request  marketplace.QueryRequest
		Response marketplace.QueryResponse
	}{
		Request:  request,
		Response: response,
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(archive); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}

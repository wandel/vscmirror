package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"net/http"
)

type Client struct {
	HttpClient *http.Client
	Version    string
}

func (c *Client) GenericQuery(ctx context.Context, req QueryRequest) (QueryResponse, error) {
	url := "https://marketplace.visualstudio.com/_apis/public/gallery/extensionquery"

	var response QueryResponse
	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(req); err != nil {
		return response, fmt.Errorf("failed to encode request body to '%s': %w", url, err)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", url, &body)
	if err != nil {
		return response, fmt.Errorf("failed to create request to '%s': %w", url, err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json;api-version=3.0-preview.1")
	request.Header.Set("User-Agent", "VSCode "+c.Version+" (Code)")
	request.Header.Set("x-market-client-id", "VSCode "+c.Version)
	request.Header.Set("x-market-user-id", "469d6fbf-fa8f-48f0-a3f9-3157bce2494b")
	slog.Info("sending request")
	resp, err := c.HttpClient.Do(request)
	if err != nil {
		return response, fmt.Errorf("failed to send request to '%s': %w", url, err)
	}
	defer resp.Body.Close()
	slog.Info("processing response")

	slog.Info("search", "status", resp.StatusCode, "url", url)

	wrapper := struct {
		Results []QueryResponse `json:"results"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return response, fmt.Errorf("failed to decode response from '%s': %w", url, err)
	}

	if len(wrapper.Results) == 0 {
		return response, fmt.Errorf("no results found for '%s'", url)
	}

	return wrapper.Results[0], nil
}

func (c *Client) GetLastestExtensionVersion(ctx context.Context, ids []string) ([]Extension, error) {
	return nil, errors.New("not implemented yet")
}

func (c *Client) GetExtensionsByName(ctx context.Context, names []string) ([]Extension, error) {
	criteria := []FilterCriteria{}
	for _, id := range names {
		criteria = append(criteria, FilterCriteria{
			FilterType: FilterTypeExtensionName,
			Value:      id,
		})
	}

	criteria = append(criteria, FilterCriteria{
		FilterType: FilterTypeExcludeWithFlags,
		Value:      "4096",
	})

	var flags QueryFlag
	flags |= QueryFlagIncludeFiles
	flags |= QueryFlagIncludeCategoryAndTags
	flags |= QueryFlagIncludeVersionProperties
	flags |= QueryFlagExcludeNonValidated
	flags |= QueryFlagIncludeAssetUri
	flags |= QueryFlagIncludeStatistics
	flags |= QueryFlagIncludeLatestVersionOnly

	request := QueryRequest{
		Flags:      flags,
		AssetTypes: []string{},
		Filters: []QueryFilter{
			QueryFilter{
				Criteria:   criteria,
				SortBy:     SortByInstallCount,
				SortOrder:  SortOrderDescending,
				PageSize:   len(names),
				PageNumber: 1,
			},
		},
	}

	response, err := c.GenericQuery(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to query extensions by name: %w", err)
	}

	return response.Extensions, nil
}

func (c *Client) GetExtensionsPaged(ctx context.Context, pageSize int, pageNumber int) ([]Extension, error) {
	criteria := []FilterCriteria{
		FilterCriteria{
			FilterType: FilterTypeInstalltionTarget,
			Value:      "Microsoft.VisualStudio.Code",
		},
	}

	var flags QueryFlag
	flags |= QueryFlagIncludeFiles
	flags |= QueryFlagIncludeCategoryAndTags
	flags |= QueryFlagIncludeVersionProperties
	flags |= QueryFlagExcludeNonValidated
	flags |= QueryFlagIncludeAssetUri
	flags |= QueryFlagIncludeStatistics
	flags |= QueryFlagIncludeLatestVersionOnly

	request := QueryRequest{
		Flags:      flags,
		AssetTypes: []string{},
		Filters: []QueryFilter{
			QueryFilter{
				Criteria:   criteria,
				SortBy:     SortByInstallCount,
				SortOrder:  SortOrderDescending,
				PageSize:   pageSize,
				PageNumber: pageNumber,
			},
		},
	}

	response, err := c.GenericQuery(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to query extensions by text: %w", err)
	}

	return response.Extensions, nil
}

func (c *Client) GetAllExtensions(ctx context.Context) iter.Seq[Extension] {
	return func(yield func(Extension) bool) {
		current := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			extensions, err := c.GetExtensionsPaged(ctx, 1000, current)
			if err != nil {
				slog.Error("failed to get page of extensions", "page", current, "error", err)
			}

			for _, extension := range extensions {
				if !yield(extension) {
					return
				}
			}
		}
	}
}

package client

import (
	"log/slog"
	"strings"
)

type QueryFlag int

const (
	QueryFlagNoneDefined                QueryFlag = 0x0
	QueryFlagIncludeVersions            QueryFlag = 0x1
	QueryFlagIncludeFiles               QueryFlag = 0x2
	QueryFlagIncludeCategoryAndTags     QueryFlag = 0x4
	QueryFlagIncludeSharedAccounts      QueryFlag = 0x8
	QueryFlagIncludeVersionProperties   QueryFlag = 0x10
	QueryFlagExcludeNonValidated        QueryFlag = 0x20
	QueryFlagIncludeInstallationTargets QueryFlag = 0x40
	QueryFlagIncludeAssetUri            QueryFlag = 0x80
	QueryFlagIncludeStatistics          QueryFlag = 0x100
	QueryFlagIncludeLatestVersionOnly   QueryFlag = 0x200
	QueryFlagUnpublished                QueryFlag = 0x1000
)

type FilterType int

const (
	FilterTypeTag                            FilterType = 1
	FilterTypeDisplayName                    FilterType = 2
	FilterTypePrivate                        FilterType = 3
	FilterTypeId                             FilterType = 4
	FilterTypeCategory                       FilterType = 5
	FilterTypeContributionType               FilterType = 6
	FilterTypeName                           FilterType = 7
	FilterTypeInstalltionTarget              FilterType = 8
	FilterTypeFeatured                       FilterType = 9
	FilterTypeSearchText                     FilterType = 10
	FilterTypeFeaturedInCategory             FilterType = 11
	FilterTypeExcludeWithFlags               FilterType = 12
	FilterTypeIncludeWithFlags               FilterType = 13
	FilterTypeLcid                           FilterType = 14
	FilterTypeInstalltionTargetVersion       FilterType = 15
	FilterTypeInstallationTargetVersionRange FilterType = 16
	FilterTypeVsixMetadata                   FilterType = 17
	FilterTypePublisherName                  FilterType = 18
	FilterTypePublisherDisplayName           FilterType = 19
	FilterTypeIncludeWithPublisherFlags      FilterType = 20
	FilterTypeOrganizationSharedWith         FilterType = 21
	FilterTypeProductArchitecture            FilterType = 22
	FilterTypeTargetPlatform                 FilterType = 23
	FilterTypeExtensionName                  FilterType = 24
)

type SortBy int

const (
	SortByNoneOrRelevance SortBy = 0
	SortByLastUpdatedDate SortBy = 1
	SortByTitle           SortBy = 2
	SortByPublisherName   SortBy = 3
	SortByInstallCount    SortBy = 4
	SortByPublishedDate   SortBy = 5
	SortByAverageRating   SortBy = 6
	SortByWeightedRating  SortBy = 12
)

type SortOrder int

const (
	SortOrderDefault    SortOrder = 0
	SortOrderAscending  SortOrder = 1
	SortOrderDescending SortOrder = 2
)

type PublishedExtensionFlags int

const (
	// No flags exist for this extension.
	PublishedExtensionFlagsNone = 0
	// The Disabled flag for an extension means the extension can't be changed and won't be used by consumers. The disabled flag is managed by the service and can't be supplied by the Extension Developers.
	PublishedExtensionFlagsDisabled = 1
	// BuiltIn Extension are available to all Tenants. An explicit registration is not required. This attribute is reserved and can't be supplied by Extension Developers.  BuiltIn extensions are by definition Public. There is no need to set the public flag for extensions marked BuiltIn.
	PublishedExtensionFlagsBuiltIn = 2
	// This was not present in the https://github.com/microsoft/azure-devops-node-api repository
	PublishedExtensionFlagsUnknown = 3
	// This extension has been validated by the service. The extension meets the requirements specified. This attribute is reserved and can't be supplied by the Extension Developers. Validation is a process that ensures that all contributions are well formed. They meet the requirements defined by the contribution type they are extending. Note this attribute will be updated asynchronously as the extension is validated by the developer of the contribution type. There will be restricted access to the extension while this process is performed.
	PublishedExtensionFlagsValidated = 4
	// Trusted extensions are ones that are given special capabilities. These tend to come from Microsoft and can't be published by the general public.  Note: BuiltIn extensions are always trusted.
	PublishedExtensionFlagsTrusted = 8
	// The Paid flag indicates that the commerce can be enabled for this extension. Publisher needs to setup Offer/Pricing plan in Azure. If Paid flag is set and a corresponding Offer is not available, the extension will automatically be marked as Preview. If the publisher intends to make the extension Paid in the future, it is mandatory to set the Preview flag. This is currently available only for VSTS extensions only.
	PublishedExtensionFlagsPaid = 16
	// This extension registration is public, making its visibility open to the public. This means all tenants have the ability to install this extension. Without this flag the extension will be private and will need to be shared with the tenants that can install it.
	PublishedExtensionFlagsPublic = 256
	// This extension has multiple versions active at one time and version discovery should be done using the defined "Version Discovery" protocol to determine the version available to a specific user or tenant.  @TODO: Link to Version Discovery Protocol.
	PublishedExtensionFlagsMultiVersion = 512
	// The system flag is reserved, and cant be used by publishers.
	PublishedExtensionFlagsSystem = 1024
	// The Preview flag indicates that the extension is still under preview (not yet of "release" quality). These extensions may be decorated differently in the gallery and may have different policies applied to them.
	PublishedExtensionFlagsPreview = 2048
	// The Unpublished flag indicates that the extension can't be installed/downloaded. Users who have installed such an extension can continue to use the extension.
	PublishedExtensionFlagsUnpublished = 4096
	// The Trial flag indicates that the extension is in Trial version. The flag is right now being used only with respect to Visual Studio extensions.
	PublishedExtensionFlagsTrial = 8192
	// The Locked flag indicates that extension has been locked from Marketplace. Further updates/acquisitions are not allowed on the extension until this is present. This should be used along with making the extension private/unpublished.
	PublishedExtensionFlagsLocked = 16384
	// This flag is set for extensions we want to hide from Marketplace home and search pages. This will be used to override the exposure of builtIn flags.
	PublishedExtensionFlagsHidden = 32768
)

type PublisherFlags int

const (
	// This should never be returned, it is used to represent a publisher who's flags haven't changed during update calls.
	PublisherFlagsUnChanged = 1073741824
	// No flags exist for this publisher.
	PublisherFlagsNone = 0
	// The Disabled flag for a publisher means the publisher can't be changed and won't be used by consumers, this extends to extensions owned by the publisher as well. The disabled flag is managed by the service and can't be supplied by the Extension Developers.
	PublisherFlagsDisabled = 1
	// A verified publisher is one that Microsoft has done some review of and ensured the publisher meets a set of requirements. The requirements to become a verified publisher are not listed here.  They can be found in public documentation (TBD).
	PublisherFlagsVerified = 2
	// A Certified publisher is one that is Microsoft verified and in addition meets a set of requirements for its published extensions. The requirements to become a certified publisher are not listed here.  They can be found in public documentation (TBD).
	PublisherFlagsCertified = 4
	// This is the set of flags that can't be supplied by the developer and is managed by the service itself.
	PublisherFlagsServiceFlags = 7
)

type PagingDirection int

const (
	PagingDirectionForward  PagingDirection = 1
	PagingDirectionBackward PagingDirection = 2
)

type QueryRequest struct {
	Filters    []QueryFilter `json:"filters"`
	AssetTypes []string      `json:"assetTypes"`
	Flags      QueryFlag     `json:"flags"`
}

type QueryResultMetadata struct {
	MetadataType  string         `json:"metadataType"`
	MetadataItems []MetadataItem `json:"metadataItems"`
}

type MetadataItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type QueryResponse struct {
	Extensions     []Extension           `json:"extensions"`
	PagingToken    string                `json:"pagingToken"`
	ResultMetadata []QueryResultMetadata `json:"resultMetadata"`
}

type QueryFilter struct {
	// The filter values define the set of values in this query. They are applied based on the QueryFilterType.
	Criteria []FilterCriteria `json:"criteria"`
	// The PagingDirection is applied to a paging token if one exists. If not the direction is ignored, and Forward from the start of the resultset is used. Direction should be left out of the request unless a paging token is used to help prevent future issues.
	Direction PagingDirection `json:"direction,omitempty"`
	// The page number requested by the user. If not provided 1 is assumed by default.
	PageNumber int `json:"pageNumber"`
	// The page size defines the number of results the caller wants for this filter. The count can't exceed the overall query size limits.
	PageSize int `json:"pageSize"`
	// The paging token is a distinct type of filter and the other filter fields are ignored. The paging token represents the continuation of a previously executed query. The information about where in the result and what fields are being filtered are embedded in the token.
	PagingToken string `json:"pagingToken,omitempty"`
	// Defines the type of sorting to be applied on the results. The page slice is cut of the sorted results only.
	SortBy SortBy `json:"sortBy"`
	// Defines the order of sorting, 1 for Ascending, 2 for Descending, else default ordering based on the SortBy value
	SortOrder SortOrder `json:"sortOrder"`
}

func (filter *QueryFilter) Matches(extension Extension) bool {
	var include bool
	for _, criteria := range filter.Criteria {
		if criteria.FilterType == FilterTypeExcludeWithFlags {
			if criteria.Matches(extension) {
				return false // exclude this extension
			}
		} else if criteria.Matches(extension) {
			include = true
		}
	}
	return include
}

func (filter QueryFilter) Compare(a, b Extension) int {
	switch filter.SortBy {
	case SortByNoneOrRelevance:
		return CompareStatistic(a, b, "install") // TODO: Implement relevance sorting
	case SortByLastUpdatedDate:
		return strings.Compare(a.LastUpdated, b.LastUpdated)
	case SortByTitle:
		return strings.Compare(a.DisplayName, b.DisplayName)
	case SortByPublisherName:
		return strings.Compare(a.Publisher.PublisherName, b.Publisher.PublisherName)
	case SortByInstallCount:
		return CompareStatistic(a, b, "install")
	case SortByPublishedDate:
		return strings.Compare(a.PublishedDate, b.PublishedDate)
	case SortByAverageRating:
		return CompareStatistic(a, b, "averageRating")
	case SortByWeightedRating:
		return CompareStatistic(a, b, "weightedRating")
	default:
		slog.Error("unknown sortBy", "sortBy", filter.SortBy)
		return 0
	}
}

func CompareStatistic(a, b Extension, name string) int {
	s1 := a.GetStatistic(name)
	s2 := b.GetStatistic(name)
	if s1 == s2 {
		return 0
	} else if s1 < s2 {
		return -1
	}
	return 1
}

type FilterCriteria struct {
	FilterType FilterType `json:"filterType"`
	Value      string     `json:"value"`
}

func (criteria *FilterCriteria) Matches(extension Extension) bool {
	switch criteria.FilterType {
	case FilterTypeId:
		return extension.ExtensionId == criteria.Value
	case FilterTypeCategory:
		return false // TODO: Implement category filtering
	case FilterTypeTag:
		return false // TODO: Implement tag filtering
	case FilterTypeExtensionName:
		tmp := extension.Publisher.PublisherName + "." + extension.ExtensionName
		return strings.EqualFold(tmp, criteria.Value)
	case FilterTypeInstalltionTarget:
		return false // TODO: Implement target filtering
	case FilterTypeFeatured:
		return false // TODO: Implement featured filtering
	case FilterTypeSearchText:
		return MatchSearchTextCriteria(extension, criteria.Value)
	case FilterTypeExcludeWithFlags:
		// this is ususally used with the value "4096" which is "unpublished"
		return false // TODO: Implement exclude with flags filtering
	default:
		return false // unknown filter type
	}
}

func MatchSearchTextCriteria(extension Extension, text string) bool {
	if text == "*" {
		return true
	} else if text == "" {
		return true
	}

	text = strings.ToLower(text)
	if strings.Contains(strings.ToLower(extension.DisplayName), text) {
		return true
	} else if strings.Contains(strings.ToLower(extension.ExtensionName), text) {
		return true
	} else if strings.Contains(strings.ToLower(extension.ShortDescription), text) {
		return true
	}

	return false
}

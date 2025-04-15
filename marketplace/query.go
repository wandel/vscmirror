package marketplace

import (
	"log/slog"
	"strings"
)

type QueryFlag int

const (
	// None is used to retrieve only the basic extension details.
	QueryFlagNone QueryFlag = 0
	// IncludeVersions will return version information for extensions returned
	QueryFlagIncludeVersions QueryFlag = 1
	// IncludeFiles will return information about which files were found within the extension that were stored independent of the manifest. When asking for files, versions will be included as well since files are returned as a property of the versions.  These files can be retrieved using the path to the file without requiring the entire manifest be downloaded.
	QueryFlagIncludeFiles QueryFlag = 2
	// Include the Categories and Tags that were added to the extension definition.
	QueryFlagIncludeCategoryAndTags QueryFlag = 4
	// Include the details about which accounts the extension has been shared with if the extension is a private extension.
	QueryFlagIncludeSharedAccounts QueryFlag = 8
	// Include properties associated with versions of the extension
	QueryFlagIncludeVersionProperties QueryFlag = 16
	// Excluding non-validated extensions will remove any extension versions that either are in the process of being validated or have failed validation.
	QueryFlagExcludeNonValidated QueryFlag = 32
	// Include the set of installation targets the extension has requested.
	QueryFlagIncludeInstallationTargets QueryFlag = 64
	// Include the base uri for assets of this extension
	QueryFlagIncludeAssetUri QueryFlag = 128
	// Include the statistics associated with this extension
	QueryFlagIncludeStatistics QueryFlag = 256
	// When retrieving versions from a query, only include the latest version of the extensions that matched. This is useful when the caller doesn't need all the published versions. It will save a significant size in the returned payload.
	QueryFlagIncludeLatestVersionOnly QueryFlag = 512
	// This flag switches the asset uri to use GetAssetByName instead of CDN When this is used, values of base asset uri and base asset uri fallback are switched When this is used, source of asset files are pointed to Gallery service always even if CDN is available
	QueryFlagUseFallbackAssetUri QueryFlag = 1024
	// This flag is used to get all the metadata values associated with the extension. This is not applicable to VSTS or VSCode extensions and usage is only internal.
	QueryFlagIncludeMetadata QueryFlag = 2048
	// This flag is used to indicate to return very small data for extension required by VS IDE. This flag is only compatible when querying is done by VS IDE
	QueryFlagIncludeMinimalPayloadForVsIde QueryFlag = 4096
	// This flag is used to get Lcid values associated with the extension. This is not applicable to VSTS or VSCode extensions and usage is only internal
	QueryFlagIncludeLcids QueryFlag = 8192
	// Include the details about which organizations the extension has been shared with if the extension is a private extension.
	QueryFlagIncludeSharedOrganizations QueryFlag = 16384
	// Include the details if an extension is in conflict list or not Currently being used for VSCode extensions.
	QueryFlagIncludeNameConflictInfo QueryFlag = 32768
	// AllAttributes is designed to be a mask that defines all sub-elements of the extension should be returned.  NOTE: This is not actually All flags. This is now locked to the set defined since changing this enum would be a breaking change and would change the behavior of anyone using it. Try not to use this value when making calls to the service, instead be explicit about the options required.
	QueryFlagAllAttributes QueryFlag = 16863
)

type FilterType int

const (
	// The values are used as tags. All tags are treated as "OR" conditions with each other. There may be some value put on the number of matched tags from the query.
	FilterTypeTag FilterType = 1
	// The Values are an ExtensionName or fragment that is used to match other extension names.
	FilterTypeDisplayName FilterType = 2
	// The Filter is one or more tokens that define what scope to return private extensions for.
	FilterTypePrivate FilterType = 3
	// Retrieve a set of extensions based on their id's. The values should be the extension id's encoded as strings.
	FilterTypeId FilterType = 4
	// The category is unlike other filters. It is AND'd with the other filters instead of being a separate query.
	FilterTypeCategory FilterType = 5
	// Certain contribution types may be indexed to allow for query by type. User defined types can't be indexed at the moment.
	FilterTypeContributionType FilterType = 6
	// Retrieve an set extension based on the name based identifier. This differs from the internal id (which is being deprecated).
	FilterTypeName FilterType = 7
	// The InstallationTarget for an extension defines the target consumer for the extension. This may be something like VS, VSOnline, or VSCode
	FilterTypeInstallationTarget FilterType = 8
	// Query for featured extensions, no value is allowed when using the query type.
	FilterTypeFeatured FilterType = 9
	// The SearchText provided by the user to search for extensions
	FilterTypeSearchText FilterType = 10
	// Query for extensions that are featured in their own category, The filterValue for this is name of category of extensions.
	FilterTypeFeaturedInCategory FilterType = 11
	// When retrieving extensions from a query, exclude the extensions which are having the given flags. The value specified for this filter should be a string representing the integer values of the flags to be excluded. In case of multiple flags to be specified, a logical OR of the interger values should be given as value for this filter This should be at most one filter of this type. This only acts as a restrictive filter after. In case of having a particular flag in both IncludeWithFlags and ExcludeWithFlags, excludeFlags will remove the included extensions giving empty result for that flag.
	FilterTypeExcludeWithFlags FilterType = 12
	// When retrieving extensions from a query, include the extensions which are having the given flags. The value specified for this filter should be a string representing the integer values of the flags to be included. In case of multiple flags to be specified, a logical OR of the integer values should be given as value for this filter This should be at most one filter of this type. This only acts as a restrictive filter after. In case of having a particular flag in both IncludeWithFlags and ExcludeWithFlags, excludeFlags will remove the included extensions giving empty result for that flag. In case of multiple flags given in IncludeWithFlags in ORed fashion, extensions having any of the given flags will be included.
	FilterTypeIncludeWithFlags FilterType = 13
	// Filter the extensions based on the LCID values applicable. Any extensions which are not having any LCID values will also be filtered. This is currently only supported for VS extensions.
	FilterTypeLcid FilterType = 14
	// Filter to provide the version of the installation target. This filter will be used along with InstallationTarget filter. The value should be a valid version string. Currently supported only if search text is provided.
	FilterTypeInstallationTargetVersion FilterType = 15
	// Filter type for specifying a range of installation target version. The filter will be used along with InstallationTarget filter. The value should be a pair of well formed version values separated by hyphen(-). Currently supported only if search text is provided.
	FilterTypeInstallationTargetVersionRange FilterType = 16
	// Filter type for specifying metadata key and value to be used for filtering.
	FilterTypeVsixMetadata FilterType = 17
	// Filter to get extensions published by a publisher having supplied internal name
	FilterTypePublisherName FilterType = 18
	// Filter to get extensions published by all publishers having supplied display name
	FilterTypePublisherDisplayName FilterType = 19
	// When retrieving extensions from a query, include the extensions which have a publisher having the given flags. The value specified for this filter should be a string representing the integer values of the flags to be included. In case of multiple flags to be specified, a logical OR of the integer values should be given as value for this filter There should be at most one filter of this type. This only acts as a restrictive filter after. In case of multiple flags given in IncludeWithFlags in ORed fashion, extensions having any of the given flags will be included.
	FilterTypeIncludeWithPublisherFlags FilterType = 20
	// Filter to get extensions shared with particular organization
	FilterTypeOrganizationSharedWith FilterType = 21
	// Filter to get VS IDE extensions by Product Architecture
	FilterTypeProductArchitecture FilterType = 22
	// Filter to get VS Code extensions by target platform.
	FilterTypeTargetPlatform FilterType = 23
	// Retrieve an extension based on the extensionName.
	FilterTypeExtensionName FilterType = 24
)

type SortBy int

const (
	// The results will be sorted by relevance in case search query is given, if no search query resutls will be provided as is
	SortByRelevance SortBy = 0
	// The results will be sorted as per Last Updated date of the extensions with recently updated at the top
	SortByLastUpdatedDate SortBy = 1
	// Results will be sorted Alphabetically as per the title of the extension
	SortByTitle SortBy = 2
	// Results will be sorted Alphabetically as per Publisher title
	SortByPublisher SortBy = 3
	// Results will be sorted by Install Count
	SortByInstallCount SortBy = 4
	// The results will be sorted as per Published date of the extensions
	SortByPublishedDate SortBy = 5
	// The results will be sorted as per Average ratings of the extensions
	SortByAverageRating SortBy = 6
	// The results will be sorted as per Trending Daily Score of the extensions
	SortByTrendingDaily SortBy = 7
	// The results will be sorted as per Trending weekly Score of the extensions
	SortByTrendingWeekly SortBy = 8
	// The results will be sorted as per Trending monthly Score of the extensions
	SortByTrendingMonthly SortBy = 9
	// The results will be sorted as per ReleaseDate of the extensions (date on which the extension first went public)
	SortByReleaseDate SortBy = 10
	// The results will be sorted as per Author defined in the VSix/Metadata. If not defined, publisher name is used This is specifically needed by VS IDE, other (new and old) clients are not encouraged to use this
	SortByAuthor SortBy = 11
	// The results will be sorted as per Weighted Rating of the extension.
	SortByWeightedRating SortBy = 12
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
	PublishedExtensionFlagsNone PublishedExtensionFlags = 0
	// The Disabled flag for an extension means the extension can't be changed and won't be used by consumers. The disabled flag is managed by the service and can't be supplied by the Extension Developers.
	PublishedExtensionFlagsDisabled PublishedExtensionFlags = 1
	// BuiltIn Extension are available to all Tenants. An explicit registration is not required. This attribute is reserved and can't be supplied by Extension Developers.  BuiltIn extensions are by definition Public. There is no need to set the public flag for extensions marked BuiltIn.
	PublishedExtensionFlagsBuiltIn PublishedExtensionFlags = 2
	// This was not present in the https://github.com/microsoft/azure-devops-node-api repository
	PublishedExtensionFlagsUnknown PublishedExtensionFlags = 3
	// This extension has been validated by the service. The extension meets the requirements specified. This attribute is reserved and can't be supplied by the Extension Developers. Validation is a process that ensures that all contributions are well formed. They meet the requirements defined by the contribution type they are extending. Note this attribute will be updated asynchronously as the extension is validated by the developer of the contribution type. There will be restricted access to the extension while this process is performed.
	PublishedExtensionFlagsValidated PublishedExtensionFlags = 4
	// Trusted extensions are ones that are given special capabilities. These tend to come from Microsoft and can't be published by the general public.  Note: BuiltIn extensions are always trusted.
	PublishedExtensionFlagsTrusted PublishedExtensionFlags = 8
	// The Paid flag indicates that the commerce can be enabled for this extension. Publisher needs to setup Offer/Pricing plan in Azure. If Paid flag is set and a corresponding Offer is not available, the extension will automatically be marked as Preview. If the publisher intends to make the extension Paid in the future, it is mandatory to set the Preview flag. This is currently available only for VSTS extensions only.
	PublishedExtensionFlagsPaid PublishedExtensionFlags = 16
	// This extension registration is public, making its visibility open to the public. This means all tenants have the ability to install this extension. Without this flag the extension will be private and will need to be shared with the tenants that can install it.
	PublishedExtensionFlagsPublic PublishedExtensionFlags = 256
	// This extension has multiple versions active at one time and version discovery should be done using the defined "Version Discovery" protocol to determine the version available to a specific user or tenant.  @TODO: Link to Version Discovery Protocol.
	PublishedExtensionFlagsMultiVersion PublishedExtensionFlags = 512
	// The system flag is reserved, and cant be used by publishers.
	PublishedExtensionFlagsSystem PublishedExtensionFlags = 1024
	// The Preview flag indicates that the extension is still under preview (not yet of "release" quality). These extensions may be decorated differently in the gallery and may have different policies applied to them.
	PublishedExtensionFlagsPreview PublishedExtensionFlags = 2048
	// The Unpublished flag indicates that the extension can't be installed/downloaded. Users who have installed such an extension can continue to use the extension.
	PublishedExtensionFlagsUnpublished PublishedExtensionFlags = 4096
	// The Trial flag indicates that the extension is in Trial version. The flag is right now being used only with respect to Visual Studio extensions.
	PublishedExtensionFlagsTrial PublishedExtensionFlags = 8192
	// The Locked flag indicates that extension has been locked from Marketplace. Further updates/acquisitions are not allowed on the extension until this is present. This should be used along with making the extension private/unpublished.
	PublishedExtensionFlagsLocked PublishedExtensionFlags = 16384
	// This flag is set for extensions we want to hide from Marketplace home and search pages. This will be used to override the exposure of builtIn flags.
	PublishedExtensionFlagsHidden PublishedExtensionFlags = 32768
)

type PublisherFlags int

const (
	// This should never be returned, it is used to represent a publisher who's flags haven't changed during update calls.
	PublisherFlagsUnChanged PublisherFlags = 1073741824
	// No flags exist for this publisher.
	PublisherFlagsNone PublisherFlags = 0
	// The Disabled flag for a publisher means the publisher can't be changed and won't be used by consumers, this extends to extensions owned by the publisher as well. The disabled flag is managed by the service and can't be supplied by the Extension Developers.
	PublisherFlagsDisabled PublisherFlags = 1
	// A verified publisher is one that Microsoft has done some review of and ensured the publisher meets a set of requirements. The requirements to become a verified publisher are not listed here.  They can be found in public documentation (TBD).
	PublisherFlagsVerified PublisherFlags = 2
	// A Certified publisher is one that is Microsoft verified and in addition meets a set of requirements for its published extensions. The requirements to become a certified publisher are not listed here.  They can be found in public documentation (TBD).
	PublisherFlagsCertified PublisherFlags = 4
	// This is the set of flags that can't be supplied by the developer and is managed by the service itself.
	PublisherFlagsServiceFlags PublisherFlags = 7
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
	for _, criteria := range filter.Criteria {
		if criteria.FilterType == FilterTypeExcludeWithFlags {
			if criteria.Matches(extension) {
				return false // exclude this extension
			}
		} else if criteria.FilterType == FilterTypeIncludeWithFlags {
			if criteria.Matches(extension) {
				return true
			}
		} else if criteria.Matches(extension) {
			return true
		}
	}
	return false
}

func ShouldSkipFirstStageFilters(filter QueryFilter) bool {
	for _, criteria := range filter.Criteria {
		switch criteria.FilterType {
		case FilterTypeExcludeWithFlags:
			continue
		case FilterTypeInstallationTarget:
			continue
		default:
			return false
		}
	}
	return true
}

func (filter *QueryFilter) FilterFirstStage(extensions []Extension) []Extension {
	var results []Extension
	for _, extension := range extensions {
		for _, criteria := range filter.Criteria {
			switch criteria.FilterType {
			case FilterTypeExcludeWithFlags:
				continue
			case FilterTypeInstallationTarget:
				continue
			default:
			}

			if criteria.Matches(extension) {
				results = append(results, extension)
				break
			}
		}
	}
	return results
}

func (filter *QueryFilter) FilterSecondStage(extensions []Extension) []Extension {
	var results []Extension
	for _, extension := range extensions {
		for _, criteria := range filter.Criteria {
			switch criteria.FilterType {
			case FilterTypeExcludeWithFlags:
			case FilterTypeInstallationTarget:
			default:
				continue
			}

			if criteria.Matches(extension) {
				results = append(results, extension)
				break
			}
		}
	}

	return results
}

func (filter QueryFilter) Compare(a, b Extension) int {
	switch filter.SortBy {
	case SortByRelevance:
		return CompareStatistic(a, b, "install") // TODO: Implement relevance sorting
	case SortByLastUpdatedDate:
		return a.LastUpdated.Compare(b.LastUpdated)
	case SortByTitle:
		return strings.Compare(a.DisplayName, b.DisplayName)
	case SortByPublisher:
		return strings.Compare(a.Publisher.PublisherName, b.Publisher.PublisherName)
	case SortByInstallCount:
		return CompareStatistic(a, b, "install")
	case SortByPublishedDate:
		return a.PublishedDate.Compare(b.PublishedDate)
	case SortByAverageRating:
		return CompareStatistic(a, b, "averageRating")
	case SortByReleaseDate:
		return a.ReleaseDate.Compare(b.ReleaseDate)
	case SortByAuthor:
		return strings.Compare(a.Publisher.PublisherName, b.Publisher.PublisherName)
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
	// case FilterTypeTag:
	// 	return false // TODO: Implement tag filtering
	case FilterTypeDisplayName:
		return strings.EqualFold(extension.DisplayName, criteria.Value)
	// case FilterTypePrivate:
	// 	return false
	case FilterTypeId:
		return extension.ExtensionId == criteria.Value
	// case FilterTypeCategory:
	// 	return false // TODO: Implement category filtering
	// case FilterTypeContributionType:
	// 	return false
	case FilterTypeName:
		tmp := extension.Publisher.PublisherName + "." + extension.ExtensionName
		return strings.EqualFold(tmp, criteria.Value)
	case FilterTypeInstallationTarget:
		return true // TODO: Implement target filtering
	// case FilterTypeFeatured:
	// 	return false // TODO: Implement featured filtering
	case FilterTypeSearchText:
		return MatchSearchTextCriteria(extension, criteria.Value)
	// case FilterTypeFeaturedInCategory:
	// 	return false
	case FilterTypeExcludeWithFlags:
		// this is ususally used with the value "4096" which is "unpublished"
		return true // TODO: Implement exclude with flags filtering
	// case FilterTypeIncludeWithFlags:
	// 	return false
	// case FilterTypeLcid:
	// 	return false
	// case FilterTypeInstallationTargetVersion:
	// 	return false
	// case FilterTypeInstallationTargetVersionRange:
	// 	return false
	// case FilterTypeVsixMetadata:
	// 	return false
	case FilterTypePublisherName:
		return strings.EqualFold(extension.Publisher.PublisherName, criteria.Value)
	case FilterTypePublisherDisplayName:
		return strings.EqualFold(extension.Publisher.DisplayName, criteria.Value)
	// case FilterTypeIncludeWithPublisherFlags:
	// 	return false
	// case FilterTypeOrganizationSharedWith:
	// 	return false
	// case FilterTypeProductArchitecture:
	// 	return false
	// case FilterTypeTargetPlatform:
	// 	return false
	case FilterTypeExtensionName:
		return strings.EqualFold(extension.DisplayName, criteria.Value)
	default:
		slog.Warn("unimplemented filter type", "FilterType", criteria.FilterType, "value", criteria.Value)
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

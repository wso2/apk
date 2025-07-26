package semanticcache

const (
	TextCleanRegex            = `^"|"$`                         //TextCleanRegex is the Regex to clean text by removing leading and trailing quotes
	HTTPProtocolType          = "HTTP"                          // HTTPProtocolType is the protocol type for HTTP
	All                       = "*"                             // All is a wildcard for all values
	AnyResponseCode           = ".*"                            // AnyResponseCode is a wildcard for any response code
	DefaultSize               = -1                              // DefaultSize is the default size for cache entries
	DefaultAddAgeHeader       = false                           // DefaultAddAgeHeader indicates whether to add the Age header to responses
	DefaultEnableCacheControl = false                           // DefaultEnableCacheControl indicates whether to enable Cache-Control headers
	RequestEmbeddings         = "requestEmbeddings"             // RequestEmbeddings is the key for request embeddings in the cache
	NoStoreString             = "no-store"                      // NoStoreString is the value for Cache-Control to indicate no caching
	CacheKey                  = "cacheKey"                      // CacheKey is the key used to store cache entries
	DatePattern               = "Mon, 02 Jan 2006 15:04:05 MST" // DatePattern is the format for date headers
	DefaultThreshold          = 80                              // DefaultThreshold is the default threshold for cache entries
	DefaultTimeout            = 5000                            // DefaultTimeout is the default timeout for cache operations in milliseconds
	VectorIndexPrefix         = "apk_semantic_cache_"           // VectorIndexPrefix is the prefix for vector index keys in the cache
	DefaultTTL                = 3600                           // DefaultTTL is the default time-to-live for cache entries in seconds (1 hours, 3600 seconds
	DefaultRequestTimeout     = 30                              // DefaultRequestTimeout is the default timeout for requests in seconds (30 seconds)
)

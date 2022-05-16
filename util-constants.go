package api_common

// API_COMMON_CODES
const API_CODE_COMMON_BAD_REQUEST = "BAD_REQUEST"
const API_CODE_COMMON_UNAUTHORIZED = "UNAUTHORIZED"
const API_CODE_COMMON_INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"

// EXIT_CODES
const EXIT_CODE_MISSING_CONFIG = 10
const EXIT_CODE_CANNOT_INIT_CONFIG = 12

// HTTP_HEADER_REQUEST_ID contains the request id of the
// API call
const HTTP_HEADER_REQUEST_ID = "X-Request-ID"

// CTX_REQUESTID defines the key used when storing the request ID
// in the locals for a specific request
const CTX_REQUESTID = "requestid"

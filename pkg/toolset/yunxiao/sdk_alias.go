package yunxiao

import sdk "github.com/futuretea/yunxiao-mcp-server/pkg/yunxiao"

type Client = sdk.Client
type ClientOption = sdk.ClientOption
type APIError = sdk.APIError
type Response = sdk.Response
type Pagination = sdk.Pagination
type ErrorCategory = sdk.ErrorCategory
type ValidationError = sdk.ValidationError

const (
	AccessTokenHeader     = sdk.AccessTokenHeader
	AccessTokenQueryParam = sdk.AccessTokenQueryParam

	ErrAuth       = sdk.ErrAuth
	ErrPermission = sdk.ErrPermission
	ErrValidation = sdk.ErrValidation
	ErrRateLimit  = sdk.ErrRateLimit
	ErrServer     = sdk.ErrServer
	ErrNetwork    = sdk.ErrNetwork
)

var (
	NewClient                 = sdk.NewClient
	WithInsecureSkipTLSVerify = sdk.WithInsecureSkipTLSVerify
	WithAccessToken           = sdk.WithAccessToken
	AccessTokenFromContext    = sdk.AccessTokenFromContext
	ClassifyError             = sdk.ClassifyError
	WrapError                 = sdk.WrapError
	EncodeRepositoryID        = sdk.EncodeRepositoryID
	PrettyResponseJSON        = sdk.PrettyResponseJSON
	EncodeFilePath            = sdk.EncodeFilePath
	encodePathValue           = sdk.EncodePathValue
)

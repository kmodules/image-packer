// Code generated by smithy-go-codegen DO NOT EDIT.

package ecr

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/smithy-go/middleware"
	smithytime "github.com/aws/smithy-go/time"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	smithywaiter "github.com/aws/smithy-go/waiter"
	"github.com/jmespath/go-jmespath"
	"time"
)

// Retrieves the results of the lifecycle policy preview request for the specified
// repository.
func (c *Client) GetLifecyclePolicyPreview(ctx context.Context, params *GetLifecyclePolicyPreviewInput, optFns ...func(*Options)) (*GetLifecyclePolicyPreviewOutput, error) {
	if params == nil {
		params = &GetLifecyclePolicyPreviewInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "GetLifecyclePolicyPreview", params, optFns, c.addOperationGetLifecyclePolicyPreviewMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*GetLifecyclePolicyPreviewOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type GetLifecyclePolicyPreviewInput struct {

	// The name of the repository.
	//
	// This member is required.
	RepositoryName *string

	// An optional parameter that filters results based on image tag status and all
	// tags, if tagged.
	Filter *types.LifecyclePolicyPreviewFilter

	// The list of imageIDs to be included.
	ImageIds []types.ImageIdentifier

	// The maximum number of repository results returned by
	// GetLifecyclePolicyPreviewRequest in  paginated output. When this parameter is
	// used, GetLifecyclePolicyPreviewRequest only returns  maxResults results in a
	// single page along with a nextToken   response element. The remaining results of
	// the initial request can be seen by sending  another
	// GetLifecyclePolicyPreviewRequest request with the returned nextToken   value.
	// This value can be between 1 and 1000. If this  parameter is not used, then
	// GetLifecyclePolicyPreviewRequest returns up to  100 results and a nextToken
	// value, if  applicable. This option cannot be used when you specify images with
	// imageIds .
	MaxResults *int32

	// The nextToken value returned from a previous paginated
	// GetLifecyclePolicyPreviewRequest request where maxResults was used and the
	// results exceeded the value of that parameter. Pagination continues from the end
	// of the  previous results that returned the nextToken value. This value is  null
	// when there are no more results to return. This option cannot be used when you
	// specify images with imageIds .
	NextToken *string

	// The Amazon Web Services account ID associated with the registry that contains
	// the repository. If you do not specify a registry, the default registry is
	// assumed.
	RegistryId *string

	noSmithyDocumentSerde
}

type GetLifecyclePolicyPreviewOutput struct {

	// The JSON lifecycle policy text.
	LifecyclePolicyText *string

	// The nextToken value to include in a future GetLifecyclePolicyPreview request.
	// When the results of a GetLifecyclePolicyPreview request exceed maxResults , this
	// value can be used to retrieve the next page of results. This value is null when
	// there are no more results to return.
	NextToken *string

	// The results of the lifecycle policy preview request.
	PreviewResults []types.LifecyclePolicyPreviewResult

	// The registry ID associated with the request.
	RegistryId *string

	// The repository name associated with the request.
	RepositoryName *string

	// The status of the lifecycle policy preview request.
	Status types.LifecyclePolicyPreviewStatus

	// The list of images that is returned as a result of the action.
	Summary *types.LifecyclePolicyPreviewSummary

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationGetLifecyclePolicyPreviewMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpGetLifecyclePolicyPreview{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpGetLifecyclePolicyPreview{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "GetLifecyclePolicyPreview"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addOpGetLifecyclePolicyPreviewValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opGetLifecyclePolicyPreview(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

// GetLifecyclePolicyPreviewAPIClient is a client that implements the
// GetLifecyclePolicyPreview operation.
type GetLifecyclePolicyPreviewAPIClient interface {
	GetLifecyclePolicyPreview(context.Context, *GetLifecyclePolicyPreviewInput, ...func(*Options)) (*GetLifecyclePolicyPreviewOutput, error)
}

var _ GetLifecyclePolicyPreviewAPIClient = (*Client)(nil)

// GetLifecyclePolicyPreviewPaginatorOptions is the paginator options for
// GetLifecyclePolicyPreview
type GetLifecyclePolicyPreviewPaginatorOptions struct {
	// The maximum number of repository results returned by
	// GetLifecyclePolicyPreviewRequest in  paginated output. When this parameter is
	// used, GetLifecyclePolicyPreviewRequest only returns  maxResults results in a
	// single page along with a nextToken   response element. The remaining results of
	// the initial request can be seen by sending  another
	// GetLifecyclePolicyPreviewRequest request with the returned nextToken   value.
	// This value can be between 1 and 1000. If this  parameter is not used, then
	// GetLifecyclePolicyPreviewRequest returns up to  100 results and a nextToken
	// value, if  applicable. This option cannot be used when you specify images with
	// imageIds .
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// GetLifecyclePolicyPreviewPaginator is a paginator for GetLifecyclePolicyPreview
type GetLifecyclePolicyPreviewPaginator struct {
	options   GetLifecyclePolicyPreviewPaginatorOptions
	client    GetLifecyclePolicyPreviewAPIClient
	params    *GetLifecyclePolicyPreviewInput
	nextToken *string
	firstPage bool
}

// NewGetLifecyclePolicyPreviewPaginator returns a new
// GetLifecyclePolicyPreviewPaginator
func NewGetLifecyclePolicyPreviewPaginator(client GetLifecyclePolicyPreviewAPIClient, params *GetLifecyclePolicyPreviewInput, optFns ...func(*GetLifecyclePolicyPreviewPaginatorOptions)) *GetLifecyclePolicyPreviewPaginator {
	if params == nil {
		params = &GetLifecyclePolicyPreviewInput{}
	}

	options := GetLifecyclePolicyPreviewPaginatorOptions{}
	if params.MaxResults != nil {
		options.Limit = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &GetLifecyclePolicyPreviewPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *GetLifecyclePolicyPreviewPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next GetLifecyclePolicyPreview page.
func (p *GetLifecyclePolicyPreviewPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*GetLifecyclePolicyPreviewOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxResults = limit

	result, err := p.client.GetLifecyclePolicyPreview(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

// LifecyclePolicyPreviewCompleteWaiterOptions are waiter options for
// LifecyclePolicyPreviewCompleteWaiter
type LifecyclePolicyPreviewCompleteWaiterOptions struct {

	// Set of options to modify how an operation is invoked. These apply to all
	// operations invoked for this client. Use functional options on operation call to
	// modify this list for per operation behavior.
	//
	// Passing options here is functionally equivalent to passing values to this
	// config's ClientOptions field that extend the inner client's APIOptions directly.
	APIOptions []func(*middleware.Stack) error

	// Functional options to be passed to all operations invoked by this client.
	//
	// Function values that modify the inner APIOptions are applied after the waiter
	// config's own APIOptions modifiers.
	ClientOptions []func(*Options)

	// MinDelay is the minimum amount of time to delay between retries. If unset,
	// LifecyclePolicyPreviewCompleteWaiter will use default minimum delay of 5
	// seconds. Note that MinDelay must resolve to a value lesser than or equal to the
	// MaxDelay.
	MinDelay time.Duration

	// MaxDelay is the maximum amount of time to delay between retries. If unset or
	// set to zero, LifecyclePolicyPreviewCompleteWaiter will use default max delay of
	// 120 seconds. Note that MaxDelay must resolve to value greater than or equal to
	// the MinDelay.
	MaxDelay time.Duration

	// LogWaitAttempts is used to enable logging for waiter retry attempts
	LogWaitAttempts bool

	// Retryable is function that can be used to override the service defined
	// waiter-behavior based on operation output, or returned error. This function is
	// used by the waiter to decide if a state is retryable or a terminal state. By
	// default service-modeled logic will populate this option. This option can thus be
	// used to define a custom waiter state with fall-back to service-modeled waiter
	// state mutators.The function returns an error in case of a failure state. In case
	// of retry state, this function returns a bool value of true and nil error, while
	// in case of success it returns a bool value of false and nil error.
	Retryable func(context.Context, *GetLifecyclePolicyPreviewInput, *GetLifecyclePolicyPreviewOutput, error) (bool, error)
}

// LifecyclePolicyPreviewCompleteWaiter defines the waiters for
// LifecyclePolicyPreviewComplete
type LifecyclePolicyPreviewCompleteWaiter struct {
	client GetLifecyclePolicyPreviewAPIClient

	options LifecyclePolicyPreviewCompleteWaiterOptions
}

// NewLifecyclePolicyPreviewCompleteWaiter constructs a
// LifecyclePolicyPreviewCompleteWaiter.
func NewLifecyclePolicyPreviewCompleteWaiter(client GetLifecyclePolicyPreviewAPIClient, optFns ...func(*LifecyclePolicyPreviewCompleteWaiterOptions)) *LifecyclePolicyPreviewCompleteWaiter {
	options := LifecyclePolicyPreviewCompleteWaiterOptions{}
	options.MinDelay = 5 * time.Second
	options.MaxDelay = 120 * time.Second
	options.Retryable = lifecyclePolicyPreviewCompleteStateRetryable

	for _, fn := range optFns {
		fn(&options)
	}
	return &LifecyclePolicyPreviewCompleteWaiter{
		client:  client,
		options: options,
	}
}

// Wait calls the waiter function for LifecyclePolicyPreviewComplete waiter. The
// maxWaitDur is the maximum wait duration the waiter will wait. The maxWaitDur is
// required and must be greater than zero.
func (w *LifecyclePolicyPreviewCompleteWaiter) Wait(ctx context.Context, params *GetLifecyclePolicyPreviewInput, maxWaitDur time.Duration, optFns ...func(*LifecyclePolicyPreviewCompleteWaiterOptions)) error {
	_, err := w.WaitForOutput(ctx, params, maxWaitDur, optFns...)
	return err
}

// WaitForOutput calls the waiter function for LifecyclePolicyPreviewComplete
// waiter and returns the output of the successful operation. The maxWaitDur is the
// maximum wait duration the waiter will wait. The maxWaitDur is required and must
// be greater than zero.
func (w *LifecyclePolicyPreviewCompleteWaiter) WaitForOutput(ctx context.Context, params *GetLifecyclePolicyPreviewInput, maxWaitDur time.Duration, optFns ...func(*LifecyclePolicyPreviewCompleteWaiterOptions)) (*GetLifecyclePolicyPreviewOutput, error) {
	if maxWaitDur <= 0 {
		return nil, fmt.Errorf("maximum wait time for waiter must be greater than zero")
	}

	options := w.options
	for _, fn := range optFns {
		fn(&options)
	}

	if options.MaxDelay <= 0 {
		options.MaxDelay = 120 * time.Second
	}

	if options.MinDelay > options.MaxDelay {
		return nil, fmt.Errorf("minimum waiter delay %v must be lesser than or equal to maximum waiter delay of %v.", options.MinDelay, options.MaxDelay)
	}

	ctx, cancelFn := context.WithTimeout(ctx, maxWaitDur)
	defer cancelFn()

	logger := smithywaiter.Logger{}
	remainingTime := maxWaitDur

	var attempt int64
	for {

		attempt++
		apiOptions := options.APIOptions
		start := time.Now()

		if options.LogWaitAttempts {
			logger.Attempt = attempt
			apiOptions = append([]func(*middleware.Stack) error{}, options.APIOptions...)
			apiOptions = append(apiOptions, logger.AddLogger)
		}

		out, err := w.client.GetLifecyclePolicyPreview(ctx, params, func(o *Options) {
			o.APIOptions = append(o.APIOptions, apiOptions...)
			for _, opt := range options.ClientOptions {
				opt(o)
			}
		})

		retryable, err := options.Retryable(ctx, params, out, err)
		if err != nil {
			return nil, err
		}
		if !retryable {
			return out, nil
		}

		remainingTime -= time.Since(start)
		if remainingTime < options.MinDelay || remainingTime <= 0 {
			break
		}

		// compute exponential backoff between waiter retries
		delay, err := smithywaiter.ComputeDelay(
			attempt, options.MinDelay, options.MaxDelay, remainingTime,
		)
		if err != nil {
			return nil, fmt.Errorf("error computing waiter delay, %w", err)
		}

		remainingTime -= delay
		// sleep for the delay amount before invoking a request
		if err := smithytime.SleepWithContext(ctx, delay); err != nil {
			return nil, fmt.Errorf("request cancelled while waiting, %w", err)
		}
	}
	return nil, fmt.Errorf("exceeded max wait time for LifecyclePolicyPreviewComplete waiter")
}

func lifecyclePolicyPreviewCompleteStateRetryable(ctx context.Context, input *GetLifecyclePolicyPreviewInput, output *GetLifecyclePolicyPreviewOutput, err error) (bool, error) {

	if err == nil {
		pathValue, err := jmespath.Search("status", output)
		if err != nil {
			return false, fmt.Errorf("error evaluating waiter state: %w", err)
		}

		expectedValue := "COMPLETE"
		value, ok := pathValue.(types.LifecyclePolicyPreviewStatus)
		if !ok {
			return false, fmt.Errorf("waiter comparator expected types.LifecyclePolicyPreviewStatus value, got %T", pathValue)
		}

		if string(value) == expectedValue {
			return false, nil
		}
	}

	if err == nil {
		pathValue, err := jmespath.Search("status", output)
		if err != nil {
			return false, fmt.Errorf("error evaluating waiter state: %w", err)
		}

		expectedValue := "FAILED"
		value, ok := pathValue.(types.LifecyclePolicyPreviewStatus)
		if !ok {
			return false, fmt.Errorf("waiter comparator expected types.LifecyclePolicyPreviewStatus value, got %T", pathValue)
		}

		if string(value) == expectedValue {
			return false, fmt.Errorf("waiter state transitioned to Failure")
		}
	}

	return true, nil
}

func newServiceMetadataMiddleware_opGetLifecyclePolicyPreview(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "GetLifecyclePolicyPreview",
	}
}

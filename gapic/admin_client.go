// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go_gapic. DO NOT EDIT.

package gapic

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"time"

	"cloud.google.com/go/longrunning"
	lroauto "cloud.google.com/go/longrunning/autogen"
	rpcpb "github.com/apigee/registry/rpc"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	gtransport "google.golang.org/api/transport/grpc"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var newAdminClientHook clientHook

// AdminCallOptions contains the retry settings for each method of AdminClient.
type AdminCallOptions struct {
	GetStatus       []gax.CallOption
	GetStorage      []gax.CallOption
	MigrateDatabase []gax.CallOption
	ListProjects    []gax.CallOption
	GetProject      []gax.CallOption
	CreateProject   []gax.CallOption
	UpdateProject   []gax.CallOption
	DeleteProject   []gax.CallOption
}

func defaultAdminGRPCClientOptions() []option.ClientOption {
	return []option.ClientOption{
		internaloption.WithDefaultEndpoint("apigeeregistry.googleapis.com:443"),
		internaloption.WithDefaultMTLSEndpoint("apigeeregistry.mtls.googleapis.com:443"),
		internaloption.WithDefaultAudience("https://apigeeregistry.googleapis.com/"),
		internaloption.WithDefaultScopes(DefaultAuthScopes()...),
		internaloption.EnableJwtWithScope(),
		option.WithGRPCDialOption(grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32))),
	}
}

func defaultAdminCallOptions() *AdminCallOptions {
	return &AdminCallOptions{
		GetStatus:       []gax.CallOption{},
		GetStorage:      []gax.CallOption{},
		MigrateDatabase: []gax.CallOption{},
		ListProjects:    []gax.CallOption{},
		GetProject:      []gax.CallOption{},
		CreateProject:   []gax.CallOption{},
		UpdateProject:   []gax.CallOption{},
		DeleteProject:   []gax.CallOption{},
	}
}

// internalAdminClient is an interface that defines the methods available from .
type internalAdminClient interface {
	Close() error
	setGoogleClientInfo(...string)
	Connection() *grpc.ClientConn
	GetStatus(context.Context, *emptypb.Empty, ...gax.CallOption) (*rpcpb.Status, error)
	GetStorage(context.Context, *emptypb.Empty, ...gax.CallOption) (*rpcpb.Storage, error)
	MigrateDatabase(context.Context, *rpcpb.MigrateDatabaseRequest, ...gax.CallOption) (*MigrateDatabaseOperation, error)
	MigrateDatabaseOperation(name string) *MigrateDatabaseOperation
	ListProjects(context.Context, *rpcpb.ListProjectsRequest, ...gax.CallOption) *ProjectIterator
	GetProject(context.Context, *rpcpb.GetProjectRequest, ...gax.CallOption) (*rpcpb.Project, error)
	CreateProject(context.Context, *rpcpb.CreateProjectRequest, ...gax.CallOption) (*rpcpb.Project, error)
	UpdateProject(context.Context, *rpcpb.UpdateProjectRequest, ...gax.CallOption) (*rpcpb.Project, error)
	DeleteProject(context.Context, *rpcpb.DeleteProjectRequest, ...gax.CallOption) error
}

// AdminClient is a client for interacting with .
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
//
// The Admin service supports setup and operation of an API registry.
// It is typically not included in hosted versions of the API.
type AdminClient struct {
	// The internal transport-dependent client.
	internalClient internalAdminClient

	// The call options for this service.
	CallOptions *AdminCallOptions

	// LROClient is used internally to handle long-running operations.
	// It is exposed so that its CallOptions can be modified if required.
	// Users should not Close this client.
	LROClient *lroauto.OperationsClient
}

// Wrapper methods routed to the internal client.

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *AdminClient) Close() error {
	return c.internalClient.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *AdminClient) setGoogleClientInfo(keyval ...string) {
	c.internalClient.setGoogleClientInfo(keyval...)
}

// Connection returns a connection to the API service.
//
// Deprecated: Connections are now pooled so this method does not always
// return the same resource.
func (c *AdminClient) Connection() *grpc.ClientConn {
	return c.internalClient.Connection()
}

// GetStatus getStatus returns the status of the service.
// (– api-linter: core::0131::request-message-name=disabled
// aip.dev/not-precedent (at http://aip.dev/not-precedent): Not in the official API. –)
// (– api-linter: core::0131::method-signature=disabled
// aip.dev/not-precedent (at http://aip.dev/not-precedent): Not in the official API. –)
// (– api-linter: core::0131::http-uri-name=disabled
// aip.dev/not-precedent (at http://aip.dev/not-precedent): Not in the official API. –)
func (c *AdminClient) GetStatus(ctx context.Context, req *emptypb.Empty, opts ...gax.CallOption) (*rpcpb.Status, error) {
	return c.internalClient.GetStatus(ctx, req, opts...)
}

// GetStorage getStorage returns information about the storage used by the service.
// (– api-linter: core::0131::request-message-name=disabled
// aip.dev/not-precedent (at http://aip.dev/not-precedent): Not in the official API. –)
// (– api-linter: core::0131::method-signature=disabled
// aip.dev/not-precedent (at http://aip.dev/not-precedent): Not in the official API. –)
// (– api-linter: core::0131::http-uri-name=disabled
// aip.dev/not-precedent (at http://aip.dev/not-precedent): Not in the official API. –)
func (c *AdminClient) GetStorage(ctx context.Context, req *emptypb.Empty, opts ...gax.CallOption) (*rpcpb.Storage, error) {
	return c.internalClient.GetStorage(ctx, req, opts...)
}

// MigrateDatabase migrateDatabase attempts to migrate the database to the current schema.
func (c *AdminClient) MigrateDatabase(ctx context.Context, req *rpcpb.MigrateDatabaseRequest, opts ...gax.CallOption) (*MigrateDatabaseOperation, error) {
	return c.internalClient.MigrateDatabase(ctx, req, opts...)
}

// MigrateDatabaseOperation returns a new MigrateDatabaseOperation from a given name.
// The name must be that of a previously created MigrateDatabaseOperation, possibly from a different process.
func (c *AdminClient) MigrateDatabaseOperation(name string) *MigrateDatabaseOperation {
	return c.internalClient.MigrateDatabaseOperation(name)
}

// ListProjects listProjects returns matching projects.
// (– api-linter: standard-methods=disabled –)
// (– api-linter: core::0132::method-signature=disabled
// aip.dev/not-precedent (at http://aip.dev/not-precedent): projects are top-level resources. –)
func (c *AdminClient) ListProjects(ctx context.Context, req *rpcpb.ListProjectsRequest, opts ...gax.CallOption) *ProjectIterator {
	return c.internalClient.ListProjects(ctx, req, opts...)
}

// GetProject getProject returns a specified project.
func (c *AdminClient) GetProject(ctx context.Context, req *rpcpb.GetProjectRequest, opts ...gax.CallOption) (*rpcpb.Project, error) {
	return c.internalClient.GetProject(ctx, req, opts...)
}

// CreateProject createProject creates a specified project.
// (– api-linter: standard-methods=disabled –)
// (– api-linter: core::0133::http-uri-parent=disabled
// aip.dev/not-precedent (at http://aip.dev/not-precedent): Project has an implicit parent. –)
// (– api-linter: core::0133::method-signature=disabled
// aip.dev/not-precedent (at http://aip.dev/not-precedent): Project has an implicit parent. –)
func (c *AdminClient) CreateProject(ctx context.Context, req *rpcpb.CreateProjectRequest, opts ...gax.CallOption) (*rpcpb.Project, error) {
	return c.internalClient.CreateProject(ctx, req, opts...)
}

// UpdateProject updateProject can be used to modify a specified project.
func (c *AdminClient) UpdateProject(ctx context.Context, req *rpcpb.UpdateProjectRequest, opts ...gax.CallOption) (*rpcpb.Project, error) {
	return c.internalClient.UpdateProject(ctx, req, opts...)
}

// DeleteProject deleteProject removes a specified project and all of the resources that it
// owns.
func (c *AdminClient) DeleteProject(ctx context.Context, req *rpcpb.DeleteProjectRequest, opts ...gax.CallOption) error {
	return c.internalClient.DeleteProject(ctx, req, opts...)
}

// adminGRPCClient is a client for interacting with  over gRPC transport.
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type adminGRPCClient struct {
	// Connection pool of gRPC connections to the service.
	connPool gtransport.ConnPool

	// flag to opt out of default deadlines via GOOGLE_API_GO_EXPERIMENTAL_DISABLE_DEFAULT_DEADLINE
	disableDeadlines bool

	// Points back to the CallOptions field of the containing AdminClient
	CallOptions **AdminCallOptions

	// The gRPC API client.
	adminClient rpcpb.AdminClient

	// LROClient is used internally to handle long-running operations.
	// It is exposed so that its CallOptions can be modified if required.
	// Users should not Close this client.
	LROClient **lroauto.OperationsClient

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewAdminClient creates a new admin client based on gRPC.
// The returned client must be Closed when it is done being used to clean up its underlying connections.
//
// The Admin service supports setup and operation of an API registry.
// It is typically not included in hosted versions of the API.
func NewAdminClient(ctx context.Context, opts ...option.ClientOption) (*AdminClient, error) {
	clientOpts := defaultAdminGRPCClientOptions()
	if newAdminClientHook != nil {
		hookOpts, err := newAdminClientHook(ctx, clientHookParams{})
		if err != nil {
			return nil, err
		}
		clientOpts = append(clientOpts, hookOpts...)
	}

	disableDeadlines, err := checkDisableDeadlines()
	if err != nil {
		return nil, err
	}

	connPool, err := gtransport.DialPool(ctx, append(clientOpts, opts...)...)
	if err != nil {
		return nil, err
	}
	client := AdminClient{CallOptions: defaultAdminCallOptions()}

	c := &adminGRPCClient{
		connPool:         connPool,
		disableDeadlines: disableDeadlines,
		adminClient:      rpcpb.NewAdminClient(connPool),
		CallOptions:      &client.CallOptions,
	}
	c.setGoogleClientInfo()

	client.internalClient = c

	client.LROClient, err = lroauto.NewOperationsClient(ctx, gtransport.WithConnPool(connPool))
	if err != nil {
		// This error "should not happen", since we are just reusing old connection pool
		// and never actually need to dial.
		// If this does happen, we could leak connp. However, we cannot close conn:
		// If the user invoked the constructor with option.WithGRPCConn,
		// we would close a connection that's still in use.
		// TODO: investigate error conditions.
		return nil, err
	}
	c.LROClient = &client.LROClient
	return &client, nil
}

// Connection returns a connection to the API service.
//
// Deprecated: Connections are now pooled so this method does not always
// return the same resource.
func (c *adminGRPCClient) Connection() *grpc.ClientConn {
	return c.connPool.Conn()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *adminGRPCClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", getVersionClient(), "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *adminGRPCClient) Close() error {
	return c.connPool.Close()
}

func (c *adminGRPCClient) GetStatus(ctx context.Context, req *emptypb.Empty, opts ...gax.CallOption) (*rpcpb.Status, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).GetStatus[0:len((*c.CallOptions).GetStatus):len((*c.CallOptions).GetStatus)], opts...)
	var resp *rpcpb.Status
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.adminClient.GetStatus(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *adminGRPCClient) GetStorage(ctx context.Context, req *emptypb.Empty, opts ...gax.CallOption) (*rpcpb.Storage, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).GetStorage[0:len((*c.CallOptions).GetStorage):len((*c.CallOptions).GetStorage)], opts...)
	var resp *rpcpb.Storage
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.adminClient.GetStorage(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *adminGRPCClient) MigrateDatabase(ctx context.Context, req *rpcpb.MigrateDatabaseRequest, opts ...gax.CallOption) (*MigrateDatabaseOperation, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).MigrateDatabase[0:len((*c.CallOptions).MigrateDatabase):len((*c.CallOptions).MigrateDatabase)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.adminClient.MigrateDatabase(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &MigrateDatabaseOperation{
		lro: longrunning.InternalNewOperation(*c.LROClient, resp),
	}, nil
}

func (c *adminGRPCClient) ListProjects(ctx context.Context, req *rpcpb.ListProjectsRequest, opts ...gax.CallOption) *ProjectIterator {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).ListProjects[0:len((*c.CallOptions).ListProjects):len((*c.CallOptions).ListProjects)], opts...)
	it := &ProjectIterator{}
	req = proto.Clone(req).(*rpcpb.ListProjectsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*rpcpb.Project, string, error) {
		resp := &rpcpb.ListProjectsResponse{}
		if pageToken != "" {
			req.PageToken = pageToken
		}
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else if pageSize != 0 {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.adminClient.ListProjects(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}

		it.Response = resp
		return resp.GetProjects(), resp.GetNextPageToken(), nil
	}
	fetch := func(pageSize int, pageToken string) (string, error) {
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil {
			return "", err
		}
		it.items = append(it.items, items...)
		return nextPageToken, nil
	}

	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	it.pageInfo.MaxSize = int(req.GetPageSize())
	it.pageInfo.Token = req.GetPageToken()

	return it
}

func (c *adminGRPCClient) GetProject(ctx context.Context, req *rpcpb.GetProjectRequest, opts ...gax.CallOption) (*rpcpb.Project, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).GetProject[0:len((*c.CallOptions).GetProject):len((*c.CallOptions).GetProject)], opts...)
	var resp *rpcpb.Project
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.adminClient.GetProject(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *adminGRPCClient) CreateProject(ctx context.Context, req *rpcpb.CreateProjectRequest, opts ...gax.CallOption) (*rpcpb.Project, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append((*c.CallOptions).CreateProject[0:len((*c.CallOptions).CreateProject):len((*c.CallOptions).CreateProject)], opts...)
	var resp *rpcpb.Project
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.adminClient.CreateProject(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *adminGRPCClient) UpdateProject(ctx context.Context, req *rpcpb.UpdateProjectRequest, opts ...gax.CallOption) (*rpcpb.Project, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "project.name", url.QueryEscape(req.GetProject().GetName())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).UpdateProject[0:len((*c.CallOptions).UpdateProject):len((*c.CallOptions).UpdateProject)], opts...)
	var resp *rpcpb.Project
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.adminClient.UpdateProject(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *adminGRPCClient) DeleteProject(ctx context.Context, req *rpcpb.DeleteProjectRequest, opts ...gax.CallOption) error {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).DeleteProject[0:len((*c.CallOptions).DeleteProject):len((*c.CallOptions).DeleteProject)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.adminClient.DeleteProject(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// MigrateDatabaseOperation manages a long-running operation from MigrateDatabase.
type MigrateDatabaseOperation struct {
	lro *longrunning.Operation
}

// MigrateDatabaseOperation returns a new MigrateDatabaseOperation from a given name.
// The name must be that of a previously created MigrateDatabaseOperation, possibly from a different process.
func (c *adminGRPCClient) MigrateDatabaseOperation(name string) *MigrateDatabaseOperation {
	return &MigrateDatabaseOperation{
		lro: longrunning.InternalNewOperation(*c.LROClient, &longrunningpb.Operation{Name: name}),
	}
}

// Wait blocks until the long-running operation is completed, returning the response and any errors encountered.
//
// See documentation of Poll for error-handling information.
func (op *MigrateDatabaseOperation) Wait(ctx context.Context, opts ...gax.CallOption) (*rpcpb.MigrateDatabaseResponse, error) {
	var resp rpcpb.MigrateDatabaseResponse
	if err := op.lro.WaitWithInterval(ctx, &resp, time.Minute, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Poll fetches the latest state of the long-running operation.
//
// Poll also fetches the latest metadata, which can be retrieved by Metadata.
//
// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and
// the operation has completed with failure, the error is returned and op.Done will return true.
// If Poll succeeds and the operation has completed successfully,
// op.Done will return true, and the response of the operation is returned.
// If Poll succeeds and the operation has not completed, the returned response and error are both nil.
func (op *MigrateDatabaseOperation) Poll(ctx context.Context, opts ...gax.CallOption) (*rpcpb.MigrateDatabaseResponse, error) {
	var resp rpcpb.MigrateDatabaseResponse
	if err := op.lro.Poll(ctx, &resp, opts...); err != nil {
		return nil, err
	}
	if !op.Done() {
		return nil, nil
	}
	return &resp, nil
}

// Metadata returns metadata associated with the long-running operation.
// Metadata itself does not contact the server, but Poll does.
// To get the latest metadata, call this method after a successful call to Poll.
// If the metadata is not available, the returned metadata and error are both nil.
func (op *MigrateDatabaseOperation) Metadata() (*rpcpb.MigrateDatabaseMetadata, error) {
	var meta rpcpb.MigrateDatabaseMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}

// Done reports whether the long-running operation has completed.
func (op *MigrateDatabaseOperation) Done() bool {
	return op.lro.Done()
}

// Name returns the name of the long-running operation.
// The name is assigned by the server and is unique within the service from which the operation is created.
func (op *MigrateDatabaseOperation) Name() string {
	return op.lro.Name()
}

// ProjectIterator manages a stream of *rpcpb.Project.
type ProjectIterator struct {
	items    []*rpcpb.Project
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// Response is the raw response for the current page.
	// It must be cast to the RPC response type.
	// Calling Next() or InternalFetch() updates this value.
	Response interface{}

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*rpcpb.Project, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *ProjectIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *ProjectIterator) Next() (*rpcpb.Project, error) {
	var item *rpcpb.Project
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *ProjectIterator) bufLen() int {
	return len(it.items)
}

func (it *ProjectIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}

func (c *AdminClient) GrpcClient() rpcpb.AdminClient {
	return c.internalClient.(*adminGRPCClient).adminClient
}

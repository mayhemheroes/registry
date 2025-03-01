// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/apigee/registry/cmd/registry/core"
	"github.com/apigee/registry/gapic"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/connection/grpctest"
	"github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/registry"
	"github.com/apigee/registry/server/registry/names"
	"github.com/apigee/registry/server/registry/test/seeder"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/proto"
)

// TestMain will set up a local RegistryServer and grpc.Server for all
// tests in this package if APG_REGISTRY_ADDRESS env var is not set
// for the client.
func TestMain(m *testing.M) {
	grpctest.TestMain(m, registry.Config{})
}

const gzipOpenAPIv3 = "application/x.openapi+gzip;version=3.0.0"

var sortActions = cmpopts.SortSlices(func(a, b *Action) bool { return a.Command < b.Command })
var styleguide = &rpc.StyleGuide{
	Id:        "registry-styleguide",
	MimeTypes: []string{gzipOpenAPIv3},
	Guidelines: []*rpc.Guideline{
		{
			Id: "Operation",
			Rules: []*rpc.Rule{
				{
					Id:             "OperationIdValidInURL",
					Linter:         "spectral",
					LinterRulename: "operation-operationId-valid-in-url",
					Severity:       rpc.Rule_WARNING,
				},
			},
			State: rpc.Guideline_ACTIVE,
		},
	},
}

func protoMarshal(m proto.Message) []byte {
	b, _ := proto.Marshal(m)
	return b
}

// Tests for artifacts as resources and specs as dependencies
func TestArtifacts(t *testing.T) {
	tests := []struct {
		desc string
		seed []seeder.RegistryResource
		want []*Action
	}{
		{
			desc: "single spec",
			seed: []seeder.RegistryResource{
				&rpc.ApiSpec{
					Name:     "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
					MimeType: gzipOpenAPIv3,
				},
			},
			want: []*Action{
				{
					Command:           "registry compute lint projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml --linter gnostic",
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/lint-gnostic",
				},
			},
		},
		{
			desc: "multiple specs",
			seed: []seeder.RegistryResource{
				&rpc.ApiSpec{
					Name:     "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
					MimeType: gzipOpenAPIv3,
				},
				&rpc.ApiSpec{
					Name:     "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml",
					MimeType: gzipOpenAPIv3,
				},
				&rpc.ApiSpec{
					Name:     "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml",
					MimeType: gzipOpenAPIv3,
				},
			},
			want: []*Action{
				{
					Command:           "registry compute lint projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml --linter gnostic",
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/lint-gnostic",
				},
				{
					Command:           "registry compute lint projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml --linter gnostic",
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml/artifacts/lint-gnostic",
				},
				{
					Command:           "registry compute lint projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml --linter gnostic",
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml/artifacts/lint-gnostic",
				},
			},
		},
	}

	const projectID = "controller-test"
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			ctx := context.Background()
			registryClient, err := connection.NewRegistryClient(ctx)
			if err != nil {
				t.Fatalf("Failed to create client: %+v", err)
			}
			t.Cleanup(func() { registryClient.Close() })

			adminClient, err := connection.NewAdminClient(ctx)
			if err != nil {
				t.Fatalf("Failed to create client: %+v", err)
			}
			t.Cleanup(func() { adminClient.Close() })

			deleteProject(ctx, adminClient, t, "controller-test")
			t.Cleanup(func() { deleteProject(ctx, adminClient, t, "controller-test") })

			client := seeder.Client{
				RegistryClient: registryClient,
				AdminClient:    adminClient,
			}
			lister := &RegistryLister{RegistryClient: registryClient}

			if err := seeder.SeedRegistry(ctx, client, test.seed...); err != nil {
				t.Fatalf("Setup: failed to seed registry: %s", err)
			}

			manifest := &rpc.Manifest{
				Id: "controller-test",
				GeneratedResources: []*rpc.GeneratedResource{
					{
						Pattern: "apis/-/versions/-/specs/-/artifacts/lint-gnostic",
						Dependencies: []*rpc.Dependency{
							{
								Pattern: "$resource.spec",
								Filter:  "mime_type.contains('openapi')",
							},
						},
						Action: "registry compute lint $resource.spec --linter gnostic",
					},
				},
			}
			actions := ProcessManifest(ctx, lister, projectID, manifest, 10)
			addSpecRevisions(t, ctx, registryClient, test.want)

			if diff := cmp.Diff(test.want, actions, sortActions); diff != "" {
				t.Errorf("ProcessManifest(%+v) returned unexpected diff (-want +got):\n%s", manifest, diff)
			}
		})
	}
}

// Tests for aggregated artifacts at api level and specs as resources
func TestAggregateArtifacts(t *testing.T) {
	tests := []struct {
		desc string
		seed []seeder.RegistryResource
		want []*Action
	}{
		{
			desc: "create artifacts",
			seed: []seeder.RegistryResource{
				// test api 1
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/test-api-1/versions/1.0.0/specs/openapi.yaml",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/test-api-1/versions/1.1.0/specs/openapi.yaml",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/test-api-1/versions/1.0.1/specs/openapi.yaml",
				},
				// test api 2
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/test-api-2/versions/1.0.0/specs/openapi.yaml",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/test-api-2/versions/1.1.0/specs/openapi.yaml",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/test-api-2/versions/1.0.1/specs/openapi.yaml",
				},
			},
			want: []*Action{
				{
					Command:           "registry compute vocabulary projects/controller-test/locations/global/apis/test-api-1",
					GeneratedResource: "projects/controller-test/locations/global/apis/test-api-1/artifacts/vocabulary",
				},
				{
					Command:           "registry compute vocabulary projects/controller-test/locations/global/apis/test-api-2",
					GeneratedResource: "projects/controller-test/locations/global/apis/test-api-2/artifacts/vocabulary",
				},
			},
		},
	}

	const projectID = "controller-test"
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			ctx := context.Background()
			registryClient, err := connection.NewRegistryClient(ctx)
			if err != nil {
				t.Fatalf("Failed to create client: %+v", err)
			}
			t.Cleanup(func() { registryClient.Close() })

			adminClient, err := connection.NewAdminClient(ctx)
			if err != nil {
				t.Fatalf("Failed to create client: %+v", err)
			}
			t.Cleanup(func() { adminClient.Close() })

			deleteProject(ctx, adminClient, t, "controller-test")
			t.Cleanup(func() { deleteProject(ctx, adminClient, t, "controller-test") })

			client := seeder.Client{
				RegistryClient: registryClient,
				AdminClient:    adminClient,
			}
			lister := &RegistryLister{RegistryClient: registryClient}

			if err := seeder.SeedRegistry(ctx, client, test.seed...); err != nil {
				t.Fatalf("Setup: failed to seed registry: %s", err)
			}

			manifest := &rpc.Manifest{

				Id: "controller-test",
				GeneratedResources: []*rpc.GeneratedResource{
					{
						Pattern: "apis/-/artifacts/vocabulary",
						Dependencies: []*rpc.Dependency{
							{
								Pattern: "$resource.api/versions/-/specs/-",
							},
						},
						Action: "registry compute vocabulary $resource.api",
					},
				},
			}
			actions := ProcessManifest(ctx, lister, projectID, manifest, 10)
			addSpecRevisions(t, ctx, registryClient, test.want)

			if diff := cmp.Diff(test.want, actions, sortActions); diff != "" {
				t.Errorf("ProcessManifest(%+v) returned unexpected diff (-want +got):\n%s", manifest, diff)
			}
		})
	}
}

// Tests for derived artifacts with artifacts as dependencies
func TestDerivedArtifacts(t *testing.T) {
	tests := []struct {
		desc string
		seed []seeder.RegistryResource
		want []*Action
	}{
		{
			desc: "create artifacts",
			seed: []seeder.RegistryResource{
				// version 1.0.0
				&rpc.Artifact{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/lint-gnostic",
				},
				&rpc.Artifact{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/complexity",
				},
				// version 1.0.1
				&rpc.Artifact{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml/artifacts/lint-gnostic",
				},
				&rpc.Artifact{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml/artifacts/complexity",
				},
				// version 1.1.0
				&rpc.Artifact{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml/artifacts/lint-gnostic",
				},
				&rpc.Artifact{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml/artifacts/complexity",
				},
			},
			want: []*Action{
				{
					Command: fmt.Sprintf(
						"registry compute summary %s %s",
						"projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/lint-gnostic",
						"projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/complexity"),
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/summary",
				},
				{
					Command: fmt.Sprintf(
						"registry compute summary %s %s",
						"projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml/artifacts/lint-gnostic",
						"projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml/artifacts/complexity"),
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml/artifacts/summary",
				},
				{
					Command: fmt.Sprintf(
						"registry compute summary %s %s",
						"projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml/artifacts/lint-gnostic",
						"projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml/artifacts/complexity"),
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml/artifacts/summary",
				},
			},
		},
		{
			desc: "missing artifacts",
			seed: []seeder.RegistryResource{
				// version 1.0.0
				&rpc.Artifact{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/lint-gnostic",
				},
				// version 1.0.1
				&rpc.Artifact{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml/artifacts/lint-gnostic",
				},
				&rpc.Artifact{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml/artifacts/complexity",
				},
				// version 1.1.0
				&rpc.Artifact{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml/artifacts/complexity",
				},
			},
			want: []*Action{
				{
					Command: fmt.Sprintf(
						"registry compute summary %s %s",
						"projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml/artifacts/lint-gnostic",
						"projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml/artifacts/complexity"),
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml/artifacts/summary",
				},
			},
		},
	}

	const projectID = "controller-test"
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			ctx := context.Background()
			registryClient, err := connection.NewRegistryClient(ctx)
			if err != nil {
				t.Fatalf("Failed to create client: %+v", err)
			}
			t.Cleanup(func() { registryClient.Close() })

			adminClient, err := connection.NewAdminClient(ctx)
			if err != nil {
				t.Fatalf("Failed to create client: %+v", err)
			}
			t.Cleanup(func() { adminClient.Close() })

			deleteProject(ctx, adminClient, t, "controller-test")
			t.Cleanup(func() { deleteProject(ctx, adminClient, t, "controller-test") })

			client := seeder.Client{
				RegistryClient: registryClient,
				AdminClient:    adminClient,
			}
			lister := &RegistryLister{RegistryClient: registryClient}

			if err := seeder.SeedRegistry(ctx, client, test.seed...); err != nil {
				t.Fatalf("Setup: failed to seed registry: %s", err)
			}

			manifest := &rpc.Manifest{
				Id: "controller-test",
				GeneratedResources: []*rpc.GeneratedResource{
					{
						Pattern: "apis/-/versions/-/specs/-/artifacts/summary",
						Dependencies: []*rpc.Dependency{
							{
								Pattern: "$resource.spec/artifacts/lint-gnostic",
							},
							{
								Pattern: "$resource.spec/artifacts/complexity",
							},
						},
						Action: "registry compute summary $resource.spec/artifacts/lint-gnostic $resource.spec/artifacts/complexity",
					},
				},
			}
			actions := ProcessManifest(ctx, lister, projectID, manifest, 10)
			addSpecRevisions(t, ctx, registryClient, test.want)

			if diff := cmp.Diff(test.want, actions, sortActions); diff != "" {
				t.Errorf("ProcessManifest(%+v) returned unexpected diff (-want +got):\n%s", manifest, diff)
			}
		})
	}
}

// Tests for receipt artifacts as generated resource
func TestReceiptArtifacts(t *testing.T) {
	tests := []struct {
		desc string
		seed []seeder.RegistryResource
		want []*Action
	}{
		{
			desc: "create artifacts",
			seed: []seeder.RegistryResource{
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml",
				},
			},
			want: []*Action{
				{
					Command:           "command projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/custom-artifact",
					RequiresReceipt:   true,
				},
				{
					Command:           "command projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml",
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml/artifacts/custom-artifact",
					RequiresReceipt:   true,
				},
				{
					Command:           "command projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml",
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml/artifacts/custom-artifact",
					RequiresReceipt:   true,
				},
			},
		},
	}

	const projectID = "controller-test"
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			ctx := context.Background()
			registryClient, err := connection.NewRegistryClient(ctx)
			if err != nil {
				t.Fatalf("Failed to create client: %+v", err)
			}
			t.Cleanup(func() { registryClient.Close() })

			adminClient, err := connection.NewAdminClient(ctx)
			if err != nil {
				t.Fatalf("Failed to create client: %+v", err)
			}
			t.Cleanup(func() { adminClient.Close() })

			deleteProject(ctx, adminClient, t, "controller-test")
			t.Cleanup(func() { deleteProject(ctx, adminClient, t, "controller-test") })

			client := seeder.Client{
				RegistryClient: registryClient,
				AdminClient:    adminClient,
			}
			lister := &RegistryLister{RegistryClient: registryClient}

			if err := seeder.SeedRegistry(ctx, client, test.seed...); err != nil {
				t.Fatalf("Setup: failed to seed registry: %s", err)
			}

			manifest := &rpc.Manifest{
				Id: "controller-test",
				GeneratedResources: []*rpc.GeneratedResource{
					{
						Pattern: "apis/-/versions/-/specs/-/artifacts/custom-artifact",
						Receipt: true,
						Dependencies: []*rpc.Dependency{
							{
								Pattern: "$resource.spec",
							},
						},
						Action: "command $resource.spec",
					},
				},
			}
			actions := ProcessManifest(ctx, lister, projectID, manifest, 10)
			addSpecRevisions(t, ctx, registryClient, test.want)

			if diff := cmp.Diff(test.want, actions, sortActions); diff != "" {
				t.Errorf("ProcessManifest(%+v) returned unexpected diff (-want +got):\n%s", manifest, diff)
			}
		})
	}
}

// Tests for receipt aggregate artifacts as generated resource
func TestReceiptAggArtifacts(t *testing.T) {
	tests := []struct {
		desc string
		seed []seeder.RegistryResource
		want []*Action
	}{
		{
			desc: "create artifacts",
			seed: []seeder.RegistryResource{
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml",
				},
			},
			want: []*Action{
				{
					Command:           "registry compute search-index projects/controller-test/locations/global/apis/-/versions/-/specs/-",
					GeneratedResource: "projects/controller-test/locations/global/artifacts/search-index",
					RequiresReceipt:   true,
				},
			},
		},
		{
			desc: "updated artifacts",
			seed: []seeder.RegistryResource{
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml",
				},
				&rpc.Artifact{
					Name: "projects/controller-test/locations/global/artifacts/search-index",
				},
				// Add a new spec to make the artifact outdated
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml",
				},
			},
			want: []*Action{
				{
					Command:           "registry compute search-index projects/controller-test/locations/global/apis/-/versions/-/specs/-",
					GeneratedResource: "projects/controller-test/locations/global/artifacts/search-index",
					RequiresReceipt:   true,
				},
			},
		},
	}

	const projectID = "controller-test"
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			ctx := context.Background()
			registryClient, err := connection.NewRegistryClient(ctx)
			if err != nil {
				t.Fatalf("Failed to create client: %+v", err)
			}
			t.Cleanup(func() { registryClient.Close() })

			adminClient, err := connection.NewAdminClient(ctx)
			if err != nil {
				t.Fatalf("Failed to create client: %+v", err)
			}
			t.Cleanup(func() { adminClient.Close() })

			deleteProject(ctx, adminClient, t, "controller-test")
			t.Cleanup(func() { deleteProject(ctx, adminClient, t, "controller-test") })

			client := seeder.Client{
				RegistryClient: registryClient,
				AdminClient:    adminClient,
			}
			lister := &RegistryLister{RegistryClient: registryClient}

			if err := seeder.SeedRegistry(ctx, client, test.seed...); err != nil {
				t.Fatalf("Setup: failed to seed registry: %s", err)
			}

			manifest := &rpc.Manifest{
				Id: "controller-test",
				GeneratedResources: []*rpc.GeneratedResource{
					{
						Pattern: "artifacts/search-index",
						Receipt: true,
						Dependencies: []*rpc.Dependency{
							{
								Pattern: "apis/-/versions/-/specs/-",
							},
						},
						Action: "registry compute search-index projects/controller-test/locations/global/apis/-/versions/-/specs/-",
					},
				},
			}
			actions := ProcessManifest(ctx, lister, projectID, manifest, 10)
			addSpecRevisions(t, ctx, registryClient, test.want)

			if diff := cmp.Diff(test.want, actions, sortActions); diff != "" {
				t.Errorf("ProcessManifest(%+v) returned unexpected diff (-want +got):\n%s", manifest, diff)
			}
		})
	}
}

// Tests for manifest with multiple entity references
func TestMultipleEntitiesArtifacts(t *testing.T) {
	tests := []struct {
		desc string
		seed []seeder.RegistryResource
		want []*Action
	}{
		{
			desc: "create spec artifacts",
			seed: []seeder.RegistryResource{
				&rpc.Artifact{
					Name:     "projects/controller-test/locations/global/artifacts/registry-styleguide",
					MimeType: core.MimeTypeForMessageType("google.cloud.apigeeregistry.v1.style.StyleGuide"),
					Contents: protoMarshal(styleguide),
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml",
				},
			},
			want: []*Action{
				{
					Command:           "registry compute conformance projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/conformance-registry-styleguide",
					RequiresReceipt:   true,
				},
				{
					Command:           "registry compute conformance projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml",
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml/artifacts/conformance-registry-styleguide",
					RequiresReceipt:   true,
				},
				{
					Command:           "registry compute conformance projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml",
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml/artifacts/conformance-registry-styleguide",
					RequiresReceipt:   true,
				},
			},
		},
		{
			desc: "outdated spec artifacts",
			seed: []seeder.RegistryResource{
				&rpc.Artifact{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/conformance-registry-styleguide",
				},
				&rpc.Artifact{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml/artifacts/conformance-registry-styleguide",
				},
				&rpc.Artifact{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml/artifacts/conformance-registry-styleguide",
				},
				//Update styleguide definition to make sure conformance artifacts are outdated
				&rpc.Artifact{
					Name:     "projects/controller-test/locations/global/artifacts/registry-styleguide",
					MimeType: core.MimeTypeForMessageType("google.cloud.apigeeregistry.v1.style.StyleGuide"),
					Contents: protoMarshal(styleguide),
				},
			},
			want: []*Action{
				{
					Command:           "registry compute conformance projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/conformance-registry-styleguide",
					RequiresReceipt:   true,
				},
				{
					Command:           "registry compute conformance projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml",
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml/artifacts/conformance-registry-styleguide",
					RequiresReceipt:   true,
				},
				{
					Command:           "registry compute conformance projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml",
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml/artifacts/conformance-registry-styleguide",
					RequiresReceipt:   true,
				},
			},
		},
		{
			desc: "missing spec dependencies",
			seed: []seeder.RegistryResource{
				&rpc.Artifact{
					Name:     "projects/controller-test/locations/global/artifacts/registry-styleguide",
					MimeType: core.MimeTypeForMessageType("google.cloud.apigeeregistry.v1.style.StyleGuide"),
					Contents: protoMarshal(styleguide),
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
				},
				&rpc.ApiVersion{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml",
				},
			},
			want: []*Action{
				{
					Command:           "registry compute conformance projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml/artifacts/conformance-registry-styleguide",
					RequiresReceipt:   true,
				},
				{
					Command:           "registry compute conformance projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml",
					GeneratedResource: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml/artifacts/conformance-registry-styleguide",
					RequiresReceipt:   true,
				},
			},
		},
	}

	const projectID = "controller-test"
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			ctx := context.Background()
			registryClient, err := connection.NewRegistryClient(ctx)
			if err != nil {
				t.Fatalf("Failed to create client: %+v", err)
			}
			t.Cleanup(func() { registryClient.Close() })

			adminClient, err := connection.NewAdminClient(ctx)
			if err != nil {
				t.Fatalf("Failed to create client: %+v", err)
			}
			t.Cleanup(func() { adminClient.Close() })

			deleteProject(ctx, adminClient, t, "controller-test")
			t.Cleanup(func() { deleteProject(ctx, adminClient, t, "controller-test") })

			client := seeder.Client{
				RegistryClient: registryClient,
				AdminClient:    adminClient,
			}
			lister := &RegistryLister{RegistryClient: registryClient}

			if err := seeder.SeedRegistry(ctx, client, test.seed...); err != nil {
				t.Fatalf("Setup: failed to seed registry: %s", err)
			}

			manifest := &rpc.Manifest{
				Id: "controller-test",
				GeneratedResources: []*rpc.GeneratedResource{
					{
						Pattern: "apis/-/versions/-/specs/-/artifacts/conformance-registry-styleguide",
						Receipt: true,
						Dependencies: []*rpc.Dependency{
							{
								Pattern: "$resource.spec",
							},
							{
								Pattern: "artifacts/registry-styleguide",
							},
						},
						Action: "registry compute conformance $resource.spec",
					},
				},
			}
			actions := ProcessManifest(ctx, lister, projectID, manifest, 10)
			addSpecRevisions(t, ctx, registryClient, test.want)

			if diff := cmp.Diff(test.want, actions, sortActions); diff != "" {
				t.Errorf("ProcessManifest(%+v) returned unexpected diff (-want +got):\n%s", manifest, diff)
			}
		})
	}
}

func addSpecRevisions(t *testing.T, ctx context.Context, registryClient *gapic.RegistryClient, actions []*Action) {
	for _, action := range actions {
		gr := action.GeneratedResource
		a, err := names.ParseArtifact(gr)
		if err != nil {
			t.Fatal("Failed to parse GeneratedResource", err)
		}
		if a.SpecID() == "" {
			return
		}
		sr := names.Spec{
			ProjectID: a.ProjectID(),
			ApiID:     a.ApiID(),
			VersionID: a.VersionID(),
			SpecID:    a.SpecID(),
		}
		if err := core.GetSpec(ctx, registryClient, sr, false, func(s *rpc.ApiSpec) error {
			action.Command = strings.ReplaceAll(action.Command,
				fmt.Sprintf("/%s", a.SpecID()), fmt.Sprintf("/%s@%s", a.SpecID(), s.GetRevisionId()))
			action.GeneratedResource = strings.ReplaceAll(action.GeneratedResource,
				fmt.Sprintf("/%s", a.SpecID()), fmt.Sprintf("/%s@%s", a.SpecID(), s.GetRevisionId()))
			return nil
		}); err != nil {
			t.Fatal("Failed GetSpecRevision", err)
		}
	}
}

func TestMaxActions(t *testing.T) {
	tests := []struct {
		desc       string
		seed       []seeder.RegistryResource
		maxActions int
	}{
		{
			desc: "generated more than maxActions",
			seed: []seeder.RegistryResource{
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml",
				},
			},
			maxActions: 2,
		},
		{
			desc: "generated less than maxActions",
			seed: []seeder.RegistryResource{
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.0/specs/openapi.yaml",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.0.1/specs/openapi.yaml",
				},
				&rpc.ApiSpec{
					Name: "projects/controller-test/locations/global/apis/petstore/versions/1.1.0/specs/openapi.yaml",
				},
			},
			maxActions: 4,
		},
	}

	const projectID = "controller-test"
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			ctx := context.Background()
			registryClient, err := connection.NewRegistryClient(ctx)
			if err != nil {
				t.Fatalf("Failed to create client: %+v", err)
			}
			t.Cleanup(func() { registryClient.Close() })

			adminClient, err := connection.NewAdminClient(ctx)
			if err != nil {
				t.Fatalf("Failed to create client: %+v", err)
			}
			t.Cleanup(func() { adminClient.Close() })

			deleteProject(ctx, adminClient, t, "controller-test")
			t.Cleanup(func() { deleteProject(ctx, adminClient, t, "controller-test") })

			client := seeder.Client{
				RegistryClient: registryClient,
				AdminClient:    adminClient,
			}
			lister := &RegistryLister{RegistryClient: registryClient}

			if err := seeder.SeedRegistry(ctx, client, test.seed...); err != nil {
				t.Fatalf("Setup: failed to seed registry: %s", err)
			}

			manifest := &rpc.Manifest{
				Id: "controller-test",
				GeneratedResources: []*rpc.GeneratedResource{
					{
						Pattern: "apis/-/versions/-/specs/-/artifacts/vocabulary",
						Dependencies: []*rpc.Dependency{
							{
								Pattern: "$resource.spec",
							},
						},
						Action: "registry compute vocabulary $resource.spec",
					},
				},
			}
			actions := ProcessManifest(ctx, lister, projectID, manifest, test.maxActions)
			if len(actions) > test.maxActions {
				t.Errorf("ProcessManifest(%+v) generated unexpected number of actions, wanted <= %d, got %d", manifest, test.maxActions, len(actions))
			}
		})
	}
}

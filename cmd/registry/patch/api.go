// Copyright 2022 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package patch

import (
	"bytes"
	"context"

	"github.com/apigee/registry/cmd/registry/core"
	"github.com/apigee/registry/gapic"
	"github.com/apigee/registry/log"
	"github.com/apigee/registry/pkg/connection"
	"github.com/apigee/registry/pkg/models"
	"github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/registry/names"
	"gopkg.in/yaml.v3"
)

func newApi(ctx context.Context, client *gapic.RegistryClient, message *rpc.Api, nested bool) (*models.Api, error) {
	apiName, err := names.ParseApi(message.Name)
	if err != nil {
		return nil, err
	}
	recommendedVersion, err := relativeVersionName(apiName, message.RecommendedVersion)
	if err != nil {
		return nil, err
	}
	recommendedDeployment, err := relativeDeploymentName(apiName, message.RecommendedDeployment)
	if err != nil {
		return nil, err
	}
	var versions []*models.ApiVersion
	var deployments []*models.ApiDeployment
	var artifacts []*models.Artifact
	if nested {
		versions = make([]*models.ApiVersion, 0)
		if err = core.ListVersions(ctx, client, apiName.Version("-"), "", func(message *rpc.ApiVersion) error {
			var version *models.ApiVersion
			version, err := newApiVersion(ctx, client, message, true)
			if err != nil {
				return err
			}
			// unset these because they can be inferred
			version.ApiVersion = ""
			version.Kind = ""
			version.Metadata.Parent = ""
			versions = append(versions, version)
			return nil
		}); err != nil {
			return nil, err
		}
		deployments = make([]*models.ApiDeployment, 0)
		if err = core.ListDeployments(ctx, client, apiName.Deployment("-"), "", func(message *rpc.ApiDeployment) error {
			var deployment *models.ApiDeployment
			deployment, err = newApiDeployment(ctx, client, message, true)
			if err != nil {
				return err
			}
			// unset these because they can be inferred
			deployment.ApiVersion = ""
			deployment.Kind = ""
			deployment.Metadata.Parent = ""
			deployments = append(deployments, deployment)
			return nil
		}); err != nil {
			return nil, err
		}
		artifacts, err = collectChildArtifacts(ctx, client, apiName.Artifact("-"))
		if err != nil {
			return nil, err
		}
	}

	return &models.Api{
		Header: models.Header{
			ApiVersion: RegistryV1,
			Kind:       "API",
			Metadata: models.Metadata{
				Name:        apiName.ApiID,
				Labels:      message.Labels,
				Annotations: message.Annotations,
			},
		},
		Data: models.ApiData{
			DisplayName:           message.DisplayName,
			Description:           message.Description,
			Availability:          message.Availability,
			RecommendedVersion:    recommendedVersion,
			RecommendedDeployment: recommendedDeployment,
			ApiVersions:           versions,
			ApiDeployments:        deployments,
			Artifacts:             artifacts,
		},
	}, err
}

func collectChildArtifacts(ctx context.Context, client *gapic.RegistryClient, artifactPattern names.Artifact) ([]*models.Artifact, error) {
	artifacts := make([]*models.Artifact, 0)
	if err := core.ListArtifacts(ctx, client, artifactPattern, "", true, func(message *rpc.Artifact) error {
		artifact, err := newArtifact(message)
		if err != nil {
			log.FromContext(ctx).Warnf("Skipping %s: %s", message.Name, err)
			return nil
		}
		if artifact.Kind == "Artifact" { // "Artifact" is the generic artifact type
			log.FromContext(ctx).Warnf("Skipping %s", message.Name)
			return nil
		}
		// unset these because they can be inferred
		artifact.ApiVersion = ""
		artifact.Metadata.Parent = ""
		artifacts = append(artifacts, artifact)
		return nil
	}); err != nil {
		return nil, err
	}
	return artifacts, nil
}

// ExportAPI allows an API to be individually exported as a YAML file.
func ExportAPI(ctx context.Context, client *gapic.RegistryClient, message *rpc.Api, nested bool) ([]byte, *models.Header, error) {
	api, err := newApi(ctx, client, message, nested)
	if err != nil {
		return nil, nil, err
	}
	var b bytes.Buffer
	err = yamlEncoder(&b).Encode(api)
	if err != nil {
		return nil, nil, err
	}
	return b.Bytes(), &api.Header, nil
}

// TODO: These functions assume that their arguments are valid names and fail the export if they aren't.
// They are used to replace the absolute resource names stored by the API with more reusable relative
// resource names that allow the exported files to be applied to arbitrary projects. This more
// concise form is also more convenient for users to specify. But it requires additional validation.
// Either we 1) only allow relative names in these fields in the YAML representation or
// 2) allow both absolute and relative names in these fields and handle them appropriately.
// We should probably support only one output format; for this, relative names seem preferable.

// relativeVersionName returns the version id if the version is within the specified API
func relativeVersionName(apiName names.Api, version string) (string, error) {
	if version == "" {
		return "", nil
	}
	versionName, err := names.ParseVersion(version)
	if err != nil {
		return "", err
	}
	if versionName.Api().String() == apiName.String() {
		return versionName.VersionID, nil
	}
	return version, nil
}

// relativeDeploymentName returns the deployment id if the deployment is within the specified API
func relativeDeploymentName(apiName names.Api, deployment string) (string, error) {
	if deployment == "" {
		return "", nil
	}
	deploymentName, err := names.ParseDeployment(deployment)
	if err != nil {
		return "", err
	}
	if deploymentName.Api().String() == apiName.String() {
		return deploymentName.DeploymentID, nil
	}
	return deployment, nil
}

// TODO: The following functions assume that their arguments are truly an ID,
// but there's there's no validation (here or in the caller) to prevent a full
// resource name from being passed in. Users specifying a full resource name will
// end up with a malformed name being stored.

// optionalVersionName returns a version name if the id is not empty
func optionalVersionName(apiName names.Api, versionID string) string {
	if versionID == "" {
		return ""
	}
	return apiName.Version(versionID).String()
}

// optionalDeploymentName returns a deployment name if the id is not empty
func optionalDeploymentName(apiName names.Api, deploymentID string) string {
	if deploymentID == "" {
		return ""
	}
	return apiName.Deployment(deploymentID).String()
}

func applyApiPatchBytes(ctx context.Context, client connection.RegistryClient, bytes []byte, parent string) error {
	var api models.Api
	err := yaml.Unmarshal(bytes, &api)
	if err != nil {
		return err
	}
	projectName, err := names.ParseProjectWithLocation(parent)
	if err != nil {
		return err
	}
	apiName := projectName.Api(api.Metadata.Name)
	req := &rpc.UpdateApiRequest{
		Api: &rpc.Api{
			Name:                  apiName.String(),
			DisplayName:           api.Data.DisplayName,
			Description:           api.Data.Description,
			Availability:          api.Data.Availability,
			RecommendedVersion:    optionalVersionName(apiName, api.Data.RecommendedVersion),
			RecommendedDeployment: optionalDeploymentName(apiName, api.Data.RecommendedDeployment),
			Labels:                api.Metadata.Labels,
			Annotations:           api.Metadata.Annotations,
		},
		AllowMissing: true,
	}
	_, err = client.UpdateApi(ctx, req)
	if err != nil {
		return err
	}
	for _, versionPatch := range api.Data.ApiVersions {
		err := applyApiVersionPatch(ctx, client, versionPatch, apiName.String())
		if err != nil {
			return err
		}
	}
	for _, deploymentPatch := range api.Data.ApiDeployments {
		err := applyApiDeploymentPatch(ctx, client, deploymentPatch, apiName.String())
		if err != nil {
			return err
		}
	}
	for _, artifactPatch := range api.Data.Artifacts {
		err = applyArtifactPatch(ctx, client, artifactPatch, apiName.String())
		if err != nil {
			return err
		}
	}
	return nil
}

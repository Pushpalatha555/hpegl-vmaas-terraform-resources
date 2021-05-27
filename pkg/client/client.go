// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"fmt"

	"github.com/hpe-hcss/hpegl-provider-lib/pkg/gltform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/hpegl-provider-lib/pkg/client"
	api_client "github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	cmp_client "github.com/hpe-hcss/vmaas-terraform-resources/internal/cmp"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/constants"
)

// keyForGLClientMap is the key in the map[string]interface{} that is passed down by hpegl used to store *Client
// This must be unique, hpegl will error-out if it isn't
const keyForGLClientMap = "vmaasClient"

// Assert that InitialiseClient satisfies the client.Initialisation interface
var _ client.Initialisation = (*InitialiseClient)(nil)

// Client is the client struct that is used by the provider code
type Client struct {
	IAMToken  string
	Location  string
	SpaceName string
	CmpClient *cmp_client.Client
}

// InitialiseClient is imported by hpegl from each service repo
type InitialiseClient struct{}

// NewClient takes an argument of all of the provider.ConfigData, and returns an interface{} and error
// If there is no error interface{} will contain *Client.
// The hpegl provider will put *Client at the value of keyForGLClientMap (returned by ServiceName) in
// the map of clients that it creates and passes down to provider code.  hpegl executes NewClient for each service.
func (i InitialiseClient) NewClient(r *schema.ResourceData) (interface{}, error) {
	token := r.Get("iam_token").(string)
	vmaasProviderSettings, err := client.GetServiceSettingsMap(constants.ServiceName, r)
	if err != nil {
		return nil, nil
	}

	// Read the value supplied in the tf file
	location := vmaasProviderSettings[constants.LOCATION].(string)
	spaceName := vmaasProviderSettings[constants.SPACENAME].(string)

	// Create VMaas Client
	client := new(Client)

	// Token to read from gltform from hpegl-lib
	if token == "" {
		gltoken, err := gltform.GetGLConfig()
		if err != nil {
			return nil, fmt.Errorf("error reading GL token file:  %w", err)
		}
		token = gltoken.Token
	}
	client.IAMToken = token

	// Get the Service Instance using agena-api call by sending space_name amd location
	serviceInstanceID := "SERVICE_INSTANCE_ID"

	// location and space_naem supplied from the terraform tf file
	client.Location = location
	client.SpaceName = spaceName

	cfg := api_client.Configuration{
		Host: constants.ServiceURL,
		DefaultHeader: map[string]string{
			"Authorization": token,
		},
	}
	apiClient := api_client.NewAPIClient(&cfg)
	client.CmpClient = cmp_client.NewClient(apiClient, cfg, serviceInstanceID)

	return client, nil
}

// ServiceName is used to return the value of keyForGLClientMap, for use by hpegl
func (i InitialiseClient) ServiceName() string {
	return keyForGLClientMap
}

// GetClientFromMetaMap is a convenience function used by provider code to extract *Client from the
// meta argument passed-in by terraform
func GetClientFromMetaMap(meta interface{}) (*Client, error) {
	cli := meta.(map[string]interface{})[keyForGLClientMap]
	if cli == nil {
		return nil, fmt.Errorf("client is not initialised, make sure that vmaas block is defined in hpegl stanza")
	}

	return cli.(*Client), nil
}
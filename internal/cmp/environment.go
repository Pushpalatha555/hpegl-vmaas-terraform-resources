// // (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

type environment struct {
	eClient *client.EnvironmentApiService
}

func newEnvironment(eClient *client.EnvironmentApiService) *environment {
	return &environment{
		eClient: eClient,
	}
}

func (c *environment) Read(ctx context.Context, d *utils.Data) error {
	logger.Debug("Get Environments")

	name := d.GetString("name")
	if err := d.Error(); err != nil {
		return err
	}
	resp, err := utils.Retry(func() (interface{}, error) {
		return c.eClient.GetAllEnvironment(ctx, map[string]string{
			nameKey: name,
		})
	})
	if err != nil {
		return err
	}
	environment := resp.(models.GetAllEnvironment)
	if len(environment.Environments) != 1 {
		return fmt.Errorf(errExactMatch, "environments")
	}
	d.SetString("code", environment.Environments[0].Code)
	d.SetID(environment.Environments[0].ID)

	// post check
	return d.Error()
}
// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/models"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/tshihad/tftags"
)

type loadBalancer struct {
	lbClient *client.LoadBalancerAPIService
}

func newLoadBalancer(loadBalancerClient *client.LoadBalancerAPIService) *loadBalancer {
	return &loadBalancer{
		lbClient: loadBalancerClient,
	}
}

func (lb *loadBalancer) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	var lbPoolResp models.GetSpecificLBPoolResp
	if err := tftags.Get(d, &lbPoolResp); err != nil {
		return err
	}
	getlbPoolResp, err := lb.lbClient.GetSpecificLBPool(ctx, lbPoolResp.ID)
	if err != nil {
		return err
	}

	return tftags.Set(d, getlbPoolResp.GetSpecificLBPoolResp)
}

func (lb *loadBalancer) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	createReq := models.CreateLBPool{}
	if err := tftags.Get(d, &createReq.CreateLBPoolReq); err != nil {
		return err
	}

	lbPoolResp, err := lb.lbClient.CreateLBPool(ctx, createReq)
	if err != nil {
		return err
	}
	if !lbPoolResp.Success {
		return fmt.Errorf(successErr, "creating loadBalancer Pool")
	}

	// wait until created
	retry := &utils.CustomRetry{
		RetryDelay:   1,
		InitialDelay: 1,
	}
	_, err = retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.lbClient.GetSpecificLBPool(ctx, lbPoolResp.LBPoolResp.ID)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createReq.CreateLBPoolReq)
}

func (lb *loadBalancer) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	lbPoolID := d.GetID()
	_, err := lb.lbClient.DeleteLBPool(ctx, lbPoolID)
	if err != nil {
		return err
	}

	return nil
}

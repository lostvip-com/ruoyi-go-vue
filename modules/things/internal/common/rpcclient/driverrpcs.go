/*******************************************************************************
 * Copyright 2017 Dell Inc.
 * Copyright (c) 2019 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package rpcclient

import (
	"errors"
	"github.com/lostvip-com/lv_grpc_driver_proto/cloudinstancecallback"
	"github.com/lostvip-com/lv_grpc_driver_proto/devicecallback"
	"github.com/lostvip-com/lv_grpc_driver_proto/productcallback"
	"github.com/lostvip-com/lv_grpc_driver_proto/thingmodel"
	"google.golang.org/grpc"
)

type RpcClient struct {
	address string
	Conn    *grpc.ClientConn
	devicecallback.DeviceCallBackServiceClient
	cloudinstancecallback.CloudInstanceCallBackServiceClient
	productcallback.ProductCallBackServiceClient
	thingmodel.ThingModelDownServiceClient
}

var rpcClient *RpcClient

func GetGRpcClientInstance(address string) (rpcClient *RpcClient, err error) {
	if rpcClient == nil {
		rpcClient, err = NewGRpcClient(address)
	}
	return rpcClient, err
}
func NewGRpcClient(address string) (*RpcClient, error) {
	var (
		err  error
		conn *grpc.ClientConn
	)
	if address == "" {
		return nil, errors.New("required address")
	}
	if conn, err = dialWithLog(address, false, "", ""); err != nil {
		return &RpcClient{}, err
	}
	return &RpcClient{
		address:                            address,
		Conn:                               conn,
		CloudInstanceCallBackServiceClient: cloudinstancecallback.NewCloudInstanceCallBackServiceClient(conn),
		DeviceCallBackServiceClient:        devicecallback.NewDeviceCallBackServiceClient(conn),
		ProductCallBackServiceClient:       productcallback.NewProductCallBackServiceClient(conn),
		ThingModelDownServiceClient:        thingmodel.NewThingModelDownServiceClient(conn),
	}, nil
}
func NewDriverRpcClientTsl(address string, certFile, serverName string) (*RpcClient, error) {
	var (
		err  error
		conn *grpc.ClientConn
	)
	if address == "" {
		return nil, errors.New("required address")
	}
	if conn, err = dialWithLog(address, true, certFile, serverName); err != nil {
		return &RpcClient{}, err
	}
	return &RpcClient{
		address:                            address,
		Conn:                               conn,
		CloudInstanceCallBackServiceClient: cloudinstancecallback.NewCloudInstanceCallBackServiceClient(conn),
		DeviceCallBackServiceClient:        devicecallback.NewDeviceCallBackServiceClient(conn),
		ProductCallBackServiceClient:       productcallback.NewProductCallBackServiceClient(conn),
		ThingModelDownServiceClient:        thingmodel.NewThingModelDownServiceClient(conn),
	}, nil
}

func (d *RpcClient) Close() error {
	return d.Conn.Close()
}

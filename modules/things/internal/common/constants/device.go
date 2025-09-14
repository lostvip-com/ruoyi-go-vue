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
package constants

import (
	"github.com/lostvip-com/lv_grpc_driver_proto/driverdevice"
)

type DeviceStatus string

const (
	DeviceOnline  = "online"
	DeviceOffline = "offline"
)

const (
	DeviceStatusUnKnow   DeviceStatus = "未知"
	DeviceStatusOnline   DeviceStatus = "在线"
	DeviceStatusOffline  DeviceStatus = "离线"
	DeviceStatusUnActive DeviceStatus = "未激活"
	DeviceStatusDisable  DeviceStatus = "禁用"
)

func (i DeviceStatus) TransformToDriverDeviceStatus() driverdevice.DeviceStatus {
	switch i {
	case DeviceStatusUnKnow:
		return driverdevice.DeviceStatus_UnKnowStatus
	case DeviceStatusOnline:
		return driverdevice.DeviceStatus_OnLine
	case DeviceStatusOffline:
		return driverdevice.DeviceStatus_OffLine
	case DeviceStatusUnActive:
		return driverdevice.DeviceStatus_UnActive
	case DeviceStatusDisable:
		return driverdevice.DeviceStatus_Disable
	default:
		return driverdevice.DeviceStatus_UnKnowStatus
	}
}

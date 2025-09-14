/*******************************************************************************
 * Copyright 2023 Winc link Inc.
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
	"fmt"
	"time"
)

const (
	EXPIRE_DEFAULT = 24 * time.Hour
)

// Gen_KEY_DEV_MODEL 用于缓存设备的物模型信息包括最新采集值
func Gen_KEY_DEV_MODEL(deviceId int, propCode string) string {
	KEY_DEV_MODEL := "dev:prop:%d:%s"
	return fmt.Sprintf(KEY_DEV_MODEL, deviceId, propCode)
}

// Gen_KEY_PRD_MODEL 产品的标准物模型缓存信息
func Gen_KEY_PRD_MODEL(productId int, propCode string) string {
	KEY_PRD_MODEL := "prd:prop:%d:%s"
	return fmt.Sprintf(KEY_PRD_MODEL, productId, propCode)
}

func Gen_KEY_DEVICE_ONLINE(deviceId int) string {
	DEVICE_ONLINE := "dev:online:%d"
	return fmt.Sprintf(DEVICE_ONLINE, deviceId)
}

func Gen_Key_DEV_INFO(deviceId int) string {

	KEY_DEV_INFO := "dev:info:%d" //设备信息
	return fmt.Sprintf(KEY_DEV_INFO, deviceId)
}

func Gen_KEY_DRIVER_INFO(driverId int) string {
	KEY_DRIVER_INFO := "driver:info:%d"
	return fmt.Sprintf(KEY_DRIVER_INFO, driverId)
}

func Gen_KEY_PRD_INFO(productId int) string {
	KEY_PRD_INFO := "prd:info:%d" //产品信息
	return fmt.Sprintf(KEY_PRD_INFO, productId)
}

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
package utils

import (
	"os"
	"regexp"
)

func CreateDirIfNotExist(filePath string) error {
	if !FilePathIsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

func FilePathIsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func RemoveFileOrDir(path string) error {
	if FilePathIsExist(path) {
		return os.RemoveAll(path)
	}
	return nil
}

// 检查文件名是有效-是否有/ \等
func CheckFileValid(fileName string) bool {
	pattern := `\\|\/`
	match, _ := regexp.MatchString(pattern, fileName)
	if match {
		return false
	}
	return true
}

func GetPwdDir() string {
	dir, _ := os.Getwd()

	return dir
}

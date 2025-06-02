// Copyright 2022 The imkuqin-zw Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package random 包提供了一些工具函数，用于生成随机字符串。
package random

import (
	"fmt"
	"strings"
)

// 设置字符集
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Int64ToString 将 int64 类型的整数转换为随机字符串
func Int64ToString(number int64) string {
	// 计算需要生成的随机字符串数量
	numChars := len(charset)
	// 初始化结果字符串
	result := make([]byte, 0)
	// 将整数映射到随机字符串
	for number > 0 {
		index := int(number) % numChars
		result = append(result, charset[index])
		number /= int64(numChars)
	}
	// 返回结果字符串
	return string(result)
}

// StringToInt64 将随机字符串还原为 int64 类型的整数
func StringToInt64(randomStr string) (int64, error) {
	numChars := len(charset)
	var result int64

	// 遍历字符串中的每个字符
	for i := len(randomStr) - 1; i >= 0; i-- {
		index := strings.IndexByte(charset, randomStr[i])
		if index == -1 {
			return 0, fmt.Errorf("invalid character '%c' in the string", randomStr[i])
		}
		result *= int64(numChars)
		result += int64(index)
	}

	return result, nil
}

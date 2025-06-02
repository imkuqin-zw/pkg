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

// Package xemail 包提供了一些工具函数，用于处理email.
package xemail

import "regexp"

// IsValidEmail check email if valid
func IsValidEmail(email string) bool {
	// 定义邮箱的正则表达式模式
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	// 编译正则表达式
	re := regexp.MustCompile(emailRegex)
	// 使用正则表达式匹配邮箱
	return re.MatchString(email)
}

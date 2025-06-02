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

// Package xarray provider some array operation
package xarray

// StableReplace stable replace src with dst, return new src and whether src changed
func StableReplace[T any](src []T, dst []T, eq func(T, T) bool) ([]T, bool) {
	needDelIdx := make(map[int]struct{})
	notNeedAddIdx := make(map[int]struct{})
	for idx, s := range src {
		exist := false
		for j, d := range dst {
			if eq(s, d) {
				exist = true
				notNeedAddIdx[j] = struct{}{}
				break
			}
		}
		if !exist {
			needDelIdx[idx] = struct{}{}
		}
	}
	if len(needDelIdx) == 0 && len(notNeedAddIdx) == len(dst) {
		return src, false
	}
	j := 0
	for i, item := range src {
		if _, ok := needDelIdx[i]; ok && j != i {
			src[j] = item
			j++
		}
	}
	src = src[:j]
	for i, item := range dst {
		if _, ok := notNeedAddIdx[i]; ok {
			continue
		}
		src = append(src, item)
	}
	return src, true
}

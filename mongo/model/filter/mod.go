/*
 Copyright 2023 The typed Authors.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package filter

func Eq(key string, val any) *Filter {
	return New().Eq(key, val)
}

func Gt(key string, val any) *Filter {
	return New().Gt(key, val)
}

func Gte(key string, val any) *Filter {
	return New().Gte(key, val)
}

func In(key string, items []any) *Filter {
	return New().In(key, items)
}

func Lt(key string, val any) *Filter {
	return New().Lt(key, val)
}

func Lte(key string, val any) *Filter {
	return New().Lte(key, val)
}

func Nin(key string, items []any) *Filter {
	return New().Nin(key, items)
}

func Ne(key string, val any) *Filter {
	return New().Ne(key, val)
}

/**
 * MIT License
 *
 * Copyright (c) 2022 qianwj
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package filters

import (
	"github.com/qianwj/typed/mongo/operator"
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/bson"
)

func Eq(key string, val any) *Filter {
	return New().Eq(key, val)
}

func Gt(key string, val any) *Filter {
	return New().Gt(key, val)
}

func Gte(key string, val any) *Filter {
	return New().Gte(key, val)
}

func In[T any](key string, items []T) *Filter {
	return New().In(key, util.ToAny(items))
}

func Lt(key string, val any) *Filter {
	return New().Lt(key, val)
}

func Lte(key string, val any) *Filter {
	return New().Lte(key, val)
}

func Nin[T any](key string, items []T) *Filter {
	return New().Nin(key, util.ToAny(items))
}

func Ne(key string, val any) *Filter {
	return New().Ne(key, val)
}

func WithInterval(key string, val *Interval) *Filter {
	return New().WithInterval(key, val)
}

func (f *Filter) Eq(key string, val any) *Filter {
	f.data.Put(key, val)
	return f
}

func (f *Filter) Gt(key string, val any) *Filter {
	f.data.PutAsHash(key, operator.Gt, val)
	return f
}

func (f *Filter) Gte(key string, val any) *Filter {
	f.data.PutAsHash(key, operator.Gte, val)
	return f
}

func (f *Filter) In(key string, items []any) *Filter {
	f.data.Put(key, bson.M{operator.In: items})
	return f
}

func (f *Filter) Lt(key string, val any) *Filter {
	f.data.PutAsHash(key, operator.Lt, val)
	return f
}

func (f *Filter) Lte(key string, val any) *Filter {
	f.data.PutAsHash(key, operator.Lte, val)
	return f
}

func (f *Filter) Nin(key string, items []any) *Filter {
	f.data.Put(key, bson.M{operator.Nin: items})
	return f
}

func (f *Filter) Ne(key string, val any) *Filter {
	f.data.Put(key, bson.M{operator.Ne: val})
	return f
}

func (f *Filter) WithInterval(key string, val *Interval) *Filter {
	f.data.Put(key, val.query())
	return f
}

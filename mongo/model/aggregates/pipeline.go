/** MIT License
 *
 * Copyright (c) 2022-2024 qianwj
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

package aggregates

import (
	"github.com/qianwj/typed/mongo/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Pipeline struct {
	stages mongo.Pipeline
}

func New() *Pipeline {
	return &Pipeline{
		stages: make(mongo.Pipeline, 0),
	}
}

func (p *Pipeline) append(key string, value any) {
	p.stages = append(p.stages, bson.D(bson.E(key, value)).Primitive())
}

func (p *Pipeline) Stages() mongo.Pipeline {
	return p.stages
}

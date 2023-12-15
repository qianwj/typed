// MIT License
//
// Copyright (c) 2022 qianwj
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package options

import (
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	CompressorType          string
	CredentialAuthMechanism string
)

const (
	Snappy      CompressorType          = "snappy"
	ZLIB        CompressorType          = "zlib"
	ZSTD        CompressorType          = "zstd"
	MongoDBCR   CredentialAuthMechanism = "MONGODB-CR"
	Plain       CredentialAuthMechanism = "PLAIN"
	GSSAPI      CredentialAuthMechanism = "GSSAPI"
	MongoDBX509 CredentialAuthMechanism = "MONGODB-X509"
	MongoDBAWS  CredentialAuthMechanism = "MONGODB-AWS"
)

const (
	AuthMechanismServiceNameKey          = "SERVICE_NAME"
	AuthMechanismCanonicalizeHostNameKey = "CANONICALIZE_HOST_NAME"
	AuthMechanismServiceRealmKey         = "SERVICE_REALM"
	AuthMechanismServiceHostKey          = "SERVICE_HOST"
	AuthMechanismAWSSessionTokenKey      = "AWS_SESSION_TOKEN"
)

func (c CompressorType) String() string {
	return string(c)
}

type ServerAPIOptions struct {
	version           *options.ServerAPIVersion
	strict            bool
	deprecationErrors bool
}

func ServerAPI() *ServerAPIOptions {
	return &ServerAPIOptions{}
}

func (s *ServerAPIOptions) Strict() {
	s.strict = true
}

func (s *ServerAPIOptions) DeprecationErrors() {
	s.deprecationErrors = true
}

func (s *ServerAPIOptions) Raw() *options.ServerAPIOptions {
	if s == nil {
		return nil
	}
	version := options.ServerAPIVersion1
	if util.NonNil(s.version) {
		version = util.FromPtr(s.version)
	}
	opts := options.ServerAPI(version)
	if s.strict {
		opts.SetStrict(true)
	}
	if s.deprecationErrors {
		opts.SetDeprecationErrors(s.deprecationErrors)
	}
	return opts
}

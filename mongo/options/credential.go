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

type CredentialOptions struct {
	authMechanism           CredentialAuthMechanism
	authMechanismProperties map[string]string
	authSource              string
	username                string
	password                *string
}

func Credential() *CredentialOptions {
	return &CredentialOptions{}
}

// AuthMechanism the mechanism to use for authentication. Supported values include "SCRAM-SHA-256", "SCRAM-SHA-1",
// "MONGODB-CR", "PLAIN", "GSSAPI", "MONGODB-X509", and "MONGODB-AWS". This can also be set through the "authMechanism"
// URI option. (e.g. "authMechanism=PLAIN"). For more information, see
// https://www.mongodb.com/docs/manual/core/authentication-mechanisms/.
func (c *CredentialOptions) AuthMechanism(mechanism CredentialAuthMechanism) *CredentialOptions {
	c.authMechanism = mechanism
	return c
}

func (c *CredentialOptions) AuthMechanismServiceName(name string) *CredentialOptions {
	c.authMechanismProperties[AuthMechanismServiceNameKey] = name
	return c
}

func (c *CredentialOptions) AuthMechanismCanonicalizeHostName(name string) *CredentialOptions {
	c.authMechanismProperties[AuthMechanismCanonicalizeHostNameKey] = name
	return c
}

func (c *CredentialOptions) AuthMechanismServiceRealm(realm string) *CredentialOptions {
	c.authMechanismProperties[AuthMechanismServiceRealmKey] = realm
	return c
}

func (c *CredentialOptions) AuthMechanismServiceHost(host string) *CredentialOptions {
	c.authMechanismProperties[AuthMechanismServiceHostKey] = host
	return c
}

func (c *CredentialOptions) AuthMechanismAWSSessionToken(token string) *CredentialOptions {
	c.authMechanismProperties[AuthMechanismAWSSessionTokenKey] = token
	return c
}

// AuthSource the name of the database to use for authentication. This defaults to "$external" for MONGODB-X509,
// GSSAPI, and PLAIN and "admin" for all other mechanisms. This can also be set through the "authSource" URI option
// (e.g. "authSource=otherDb").
func (c *CredentialOptions) AuthSource(authSource string) *CredentialOptions {
	c.authSource = authSource
	return c
}

// Username the username for authentication. This can also be set through the URI as a username:password pair before
// the first @ character. For example, a URI for user "user", password "pwd", and host "localhost:27017" would be
// "mongodb://user:pwd@localhost:27017". This is optional for X509 authentication and will be extracted from the
// client certificate if not specified.
func (c *CredentialOptions) Username(name string) *CredentialOptions {
	c.username = name
	return c
}

// Password the password for authentication. This must not be specified for X509 and is optional for GSSAPI
// authentication.
func (c *CredentialOptions) Password(pwd string) *CredentialOptions {
	c.password = util.ToPtr(pwd)
	return c
}

func (c *CredentialOptions) Auth() options.Credential {
	if c == nil {
		return options.Credential{}
	}
	auth := options.Credential{
		AuthMechanism:           string(c.authMechanism),
		AuthMechanismProperties: c.authMechanismProperties,
		AuthSource:              c.authSource,
		Username:                c.username,
	}
	if util.NonNil(c.password) && util.FromPtr(c.password) != "" {
		auth.Password = util.FromPtr(c.password)
		auth.PasswordSet = true
	}
	return auth
}

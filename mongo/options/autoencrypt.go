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
	"crypto/tls"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AutoEncryptionOptions struct {
	internal *options.AutoEncryptionOptions
}

func NewAutoEncryption() *AutoEncryptionOptions {
	return &AutoEncryptionOptions{
		internal: options.AutoEncryption(),
	}
}

// KeyVaultClientOptions specifies options for the client used to communicate with the key vault collection.
//
// If this is set, it is used to create an internal mongo.Client.
// Otherwise, if the target mongo.Client being configured has an unlimited connection pool size (i.e. maxPoolSize=0),
// it is reused to interact with the key vault collection.
// Otherwise, if the target mongo.Client has a limited connection pool size, a separate internal mongo.Client is used
// (and created if necessary). The internal mongo.Client may be shared during automatic encryption (if
// BypassAutomaticEncryption is false). The internal mongo.Client is configured with the same options as the target
// mongo.Client except minPoolSize is set to 0 and AutoEncryptionOptions is omitted.
func (a *AutoEncryptionOptions) KeyVaultClientOptions(opts *options.ClientOptions) *AutoEncryptionOptions {
	a.internal.SetKeyVaultClientOptions(opts)
	return a
}

// KeyVaultNamespace specifies the namespace of the key vault collection. This is required.
func (a *AutoEncryptionOptions) KeyVaultNamespace(ns string) *AutoEncryptionOptions {
	a.internal.SetKeyVaultNamespace(ns)
	return a
}

// KmsProviders specifies options for KMS providers. This is required.
func (a *AutoEncryptionOptions) KmsProviders(providers map[string]map[string]interface{}) *AutoEncryptionOptions {
	a.internal.SetKmsProviders(providers)
	return a
}

// SchemaMap specifies a map from namespace to local schema document. Schemas supplied in the schemaMap only apply
// to configuring automatic encryption for client side encryption. Other validation rules in the JSON schema will not
// be enforced by the driver and will result in an error.
//
// Supplying a schemaMap provides more security than relying on JSON Schemas obtained from the server. It protects
// against a malicious server advertising a false JSON Schema, which could trick the client into sending unencrypted
// data that should be encrypted.
func (a *AutoEncryptionOptions) SchemaMap(schemaMap map[string]interface{}) *AutoEncryptionOptions {
	a.internal.SetSchemaMap(schemaMap)
	return a
}

// BypassAutoEncryption specifies whether or not auto encryption should be done.
//
// If this is unset or false and target mongo.Client being configured has an unlimited connection pool size
// (i.e. maxPoolSize=0), it is reused in the process of auto encryption.
// Otherwise, if the target mongo.Client has a limited connection pool size, a separate internal mongo.Client is used
// (and created if necessary). The internal mongo.Client may be shared for key vault operations (if KeyVaultClient is
// unset). The internal mongo.Client is configured with the same options as the target mongo.Client except minPoolSize
// is set to 0 and AutoEncryptionOptions is omitted.
func (a *AutoEncryptionOptions) BypassAutoEncryption() *AutoEncryptionOptions {
	a.internal.SetBypassAutoEncryption(true)
	return a
}

// ExtraOptions specifies a map of options to configure the mongocryptd process or mongo_crypt shared library.
//
// # Supported Extra Options
//
// "mongocryptdURI" - The mongocryptd URI. Allows setting a custom URI used to communicate with the
// mongocryptd process. The default is "mongodb://localhost:27020", which works with the default
// mongocryptd process spawned by the Client. Must be a string.
//
// "mongocryptdBypassSpawn" - If set to true, the Client will not attempt to spawn a mongocryptd
// process. Must be a bool.
//
// "mongocryptdSpawnPath" - The path used when spawning mongocryptd.
// Defaults to empty string and spawns mongocryptd from system path. Must be a string.
//
// "mongocryptdSpawnArgs" - Command line arguments passed when spawning mongocryptd.
// Defaults to ["--idleShutdownTimeoutSecs=60"]. Must be an array of strings.
//
// "cryptSharedLibRequired" - If set to true, Client creation will return an error if the
// crypt_shared library is not loaded. If unset or set to false, Client creation will not return an
// error if the crypt_shared library is not loaded. The default is unset. Must be a bool.
//
// "cryptSharedLibPath" - The crypt_shared library override path. This must be the path to the
// crypt_shared dynamic library file (for example, a .so, .dll, or .dylib file), not the directory
// that contains it. If the override path is a relative path, it will be resolved relative to the
// working directory of the process. If the override path is a relative path and the first path
// component is the literal string "$ORIGIN", the "$ORIGIN" component will be replaced by the
// absolute path to the directory containing the linked libmongocrypt library. Setting an override
// path disables the default system library search path. If an override path is specified but the
// crypt_shared library cannot be loaded, Client creation will return an error. Must be a string.
func (a *AutoEncryptionOptions) ExtraOptions(extraOpts map[string]interface{}) *AutoEncryptionOptions {
	a.internal.SetExtraOptions(extraOpts)
	return a
}

// TLSConfig specifies tls.Config instances for each KMS provider to use to configure TLS on all connections created
// to the KMS provider.
//
// This should only be used to set custom TLS configurations. By default, the connection will use an empty tls.Config{} with MinVersion set to tls.VersionTLS12.
func (a *AutoEncryptionOptions) TLSConfig(tlsOpts map[string]*tls.Config) *AutoEncryptionOptions {
	a.internal.SetTLSConfig(tlsOpts)
	return a
}

// EncryptedFieldsMap specifies a map from namespace to local EncryptedFieldsMap document.
// EncryptedFieldsMap is used for Queryable Encryption.
func (a *AutoEncryptionOptions) EncryptedFieldsMap(ef map[string]interface{}) *AutoEncryptionOptions {
	a.internal.SetEncryptedFieldsMap(ef)
	return a
}

// BypassQueryAnalysis specifies whether or not query analysis should be used for automatic encryption.
// Use this option when using explicit encryption with Queryable Encryption.
func (a *AutoEncryptionOptions) BypassQueryAnalysis() *AutoEncryptionOptions {
	a.internal.SetBypassQueryAnalysis(true)
	return a
}

func (a *AutoEncryptionOptions) Raw() *options.AutoEncryptionOptions {
	return a.internal
}

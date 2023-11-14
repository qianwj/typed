package client

import (
	"context"
	"crypto/tls"
	"github.com/qianwj/typed/mongo/options"
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"net/http"
	"time"
)

type Builder struct {
	opts                *rawopts.ClientOptions
	defaultDatabaseName string
	pingReadpref        *readpref.ReadPref
}

func NewBuilder() *Builder {
	return &Builder{opts: rawopts.Client()}
}

func (b *Builder) ApplyUri(uri string) *Builder {
	b.opts.ApplyURI(uri)
	return b
}

func (b *Builder) Ping(readpref *readpref.ReadPref) *Builder {
	b.pingReadpref = readpref
	return b
}

// AppName specifies an application name that is sent to the server when creating new connections. It is used by the
// server to log connection and profiling information (e.g. slow query logs). This can also be set through the "appName"
// URI option (e.g "appName=example_application"). The default is empty, meaning no app name will be sent.
func (b *Builder) AppName(name string) *Builder {
	b.opts.SetAppName(name)
	return b
}

// Auth specifies a Credential containing options for configuring authentication. See the options.Credential
// documentation for more information about Credential fields. The default is an empty Credential, meaning no
// authentication will be configured.
func (b *Builder) Auth(auth options.Credential) *Builder {
	b.opts.SetAuth(rawopts.Credential(auth))
	return b
}

// Compressors sets the compressors that can be used when communicating with a server. Valid values are:
//
// 1. "snappy" - requires server version >= 3.4
//
// 2. "zlib" - requires server version >= 3.6
//
// 3. "zstd" - requires server version >= 4.2, and driver version >= 1.2.0 with cgo support enabled or driver
// version >= 1.3.0 without cgo.
//
// If this option is specified, the driver will perform a negotiation with the server to determine a common list of of
// compressors and will use the first one in that list when performing operations. See
// https://www.mongodb.com/docs/manual/reference/program/mongod/#cmdoption-mongod-networkmessagecompressors for more
// information about configuring compression on the server and the server-side defaults.
//
// This can also be set through the "compressors" URI option (e.g. "compressors=zstd,zlib,snappy"). The default is
// an empty slice, meaning no compression will be enabled.
func (b *Builder) Compressors(comps []options.ClientCompressorType) *Builder {
	b.opts.SetCompressors(util.OrderedMap(comps, options.ClientCompressorType.String))
	return b
}

// ConnectTimeout specifies a timeout that is used for creating connections to the server. If a custom Dialer is
// specified through Dialer, this option must not be used. This can be set through ApplyURI with the
// "connectTimeoutMS" (e.g "connectTimeoutMS=30") option. If set to 0, no timeout will be used. The default is 30
// seconds.
func (b *Builder) ConnectTimeout(d time.Duration) *Builder {
	b.opts.SetConnectTimeout(d)
	return b
}

// Dialer specifies a custom ContextDialer to be used to create new connections to the server. The default is a
// net.Dialer with the Timeout field set to ConnectTimeout. See https://golang.org/pkg/net/#Dialer for more information
// about the net.Dialer type.
func (b *Builder) Dialer(d rawopts.ContextDialer) *Builder {
	b.opts.SetDialer(d)
	return b
}

// Direct specifies whether or not a direct connect should be made. If set to true, the driver will only connect to
// the host provided in the URI and will not discover other hosts in the cluster. This can also be set through the
// "directConnection" URI option. This option cannot be set to true if multiple hosts are specified, either through
// ApplyURI or Hosts, or an SRV URI is used.
//
// As of driver version 1.4, the "connect" URI option has been deprecated and replaced with "directConnection". The
// "connect" URI option has two values:
//
// 1. "connect=direct" for direct connections. This corresponds to "directConnection=true".
//
// 2. "connect=automatic" for automatic discovery. This corresponds to "directConnection=false"
//
// If the "connect" and "directConnection" URI options are both specified in the connection string, their values must
// not conflict. Direct connections are not valid if multiple hosts are specified or an SRV URI is used. The default
// value for this option is false. If you use this, this option will be true.
func (b *Builder) Direct() *Builder {
	b.opts.SetDirect(true)
	return b
}

// HeartbeatInterval specifies the amount of time to wait between periodic background server checks. This can also be
// set through the "heartbeatIntervalMS" URI option (e.g. "heartbeatIntervalMS=10000"). The default is 10 seconds.
func (b *Builder) HeartbeatInterval(d time.Duration) *Builder {
	b.opts.SetHeartbeatInterval(d)
	return b
}

// Hosts specifies a list of host names or IP addresses for servers in a cluster. Both IPv4 and IPv6 addresses are
// supported. IPv6 literals must be enclosed in '[]' following RFC-2732 syntax.
//
// Hosts can also be specified as a comma-separated list in a URI. For example, to include "localhost:27017" and
// "localhost:27018", a URI could be "mongodb://localhost:27017,localhost:27018". The default is ["localhost:27017"]
func (b *Builder) Hosts(s []string) *Builder {
	b.opts.SetHosts(s)
	return b
}

// LoadBalanced specifies whether or not the MongoDB deployment is hosted behind a load balancer. This can also be
// set through the "loadBalanced" URI option. The driver will error during Client configuration if this option is set
// to true and one of the following conditions are met:
//
// 1. Multiple hosts are specified, either via the ApplyURI or Hosts methods. This includes the case where an SRV
// URI is used and the SRV record resolves to multiple hostnames.
// 2. A replica set name is specified, either via the URI or the Replica method.
// 3. The options specify whether or not a direct connection should be made, either via the URI or the Direct method.
//
// The default value is false.
func (b *Builder) LoadBalanced(lb bool) *Builder {
	b.opts.SetLoadBalanced(lb)
	return b
}

// LocalThreshold specifies the width of the 'latency window': when choosing between multiple suitable servers for an
// operation, this is the acceptable non-negative delta between shortest and longest average round-trip times. A server
// within the latency window is selected randomly. This can also be set through the "localThresholdMS" URI option (e.g.
// "localThresholdMS=15000"). The default is 15 milliseconds.
func (b *Builder) LocalThreshold(d time.Duration) *Builder {
	b.opts.SetLocalThreshold(d)
	return b
}

// Logger specifies a LoggerBuilder containing options for
// configuring a logger.
func (b *Builder) Logger(builder *LoggerBuilder) *Builder {
	b.opts.SetLoggerOptions(builder.build())
	return b
}

// MaxConnIdleTime specifies the maximum amount of time that a connection will remain idle in a connection pool
// before it is removed from the pool and closed. This can also be set through the "maxIdleTimeMS" URI option (e.g.
// "maxIdleTimeMS=10000"). The default is 0, meaning a connection can remain unused indefinitely.
func (b *Builder) MaxConnIdleTime(d time.Duration) *Builder {
	b.opts.SetMaxConnIdleTime(d)
	return b
}

// MaxPoolSize specifies that maximum number of connections allowed in the driver's connection pool to each server.
// Requests to a server will block if this maximum is reached. This can also be set through the "maxPoolSize" URI option
// (e.g. "maxPoolSize=100"). If this is 0, maximum connection pool size is not limited. The default is 100.
func (b *Builder) MaxPoolSize(u uint64) *Builder {
	b.opts.SetMaxPoolSize(u)
	return b
}

// MinPoolSize specifies the minimum number of connections allowed in the driver's connection pool to each server. If
// this is non-zero, each server's pool will be maintained in the background to ensure that the size does not fall below
// the minimum. This can also be set through the "minPoolSize" URI option (e.g. "minPoolSize=100"). The default is 0.
func (b *Builder) MinPoolSize(u uint64) *Builder {
	b.opts.SetMinPoolSize(u)
	return b
}

// MaxConnecting specifies the maximum number of connections a connection pool may establish simultaneously. This can
// also be set through the "maxConnecting" URI option (e.g. "maxConnecting=2"). If this is 0, the default is used. The
// default is 2. Values greater than 100 are not recommended.
func (b *Builder) MaxConnecting(u uint64) *Builder {
	b.opts.SetMaxConnecting(u)
	return b
}

// PoolMonitor specifies a PoolMonitor to receive connection pool events. See the event.PoolMonitor documentation
// for more information about the structure of the monitor and events that can be received.
func (b *Builder) PoolMonitor(m *event.PoolMonitor) *Builder {
	b.opts.SetPoolMonitor(m)
	return b
}

// Monitor specifies a CommandMonitor to receive command events. See the event.CommandMonitor documentation for more
// information about the structure of the monitor and events that can be received.
func (b *Builder) Monitor(m *event.CommandMonitor) *Builder {
	b.opts.SetMonitor(m)
	return b
}

// ServerMonitor specifies an SDAM monitor used to monitor SDAM events.
func (b *Builder) ServerMonitor(m *event.ServerMonitor) *Builder {
	b.opts.SetServerMonitor(m)
	return b
}

// ReadConcern specifies the read concern to use for read operations. A read concern level can also be set through
// the "readConcernLevel" URI option (e.g. "readConcernLevel=majority"). The default is nil, meaning the server will use
// its configured default.
func (b *Builder) ReadConcern(rc *readconcern.ReadConcern) *Builder {
	b.opts.SetReadConcern(rc)
	return b
}

// ReadPreference specifies the read preference to use for read operations. This can also be set through the
// following URI options:
//
// 1. "readPreference" - Specify the read preference mode (e.g. "readPreference=primary").
//
// 2. "readPreferenceTags": Specify one or more read preference tags
// (e.g. "readPreferenceTags=region:south,datacenter:A").
//
// 3. "maxStalenessSeconds" (or "maxStaleness"): Specify a maximum replication lag for reads from secondaries in a
// replica set (e.g. "maxStalenessSeconds=10").
//
// The default is readpref.Primary(). See https://www.mongodb.com/docs/manual/core/read-preference/#read-preference for
// more information about read preferences.
func (b *Builder) ReadPreference(rp *readpref.ReadPref) *Builder {
	b.opts.SetReadPreference(rp)
	return b
}

// BSONOptions configures optional BSON marshaling and unmarshaling behavior.
func (b *Builder) BSONOptions(opts *options.BSONOptions) *Builder {
	b.opts.SetBSONOptions((*rawopts.BSONOptions)(opts))
	return b
}

// Registry specifies the BSON registry to use for BSON marshalling/unmarshalling operations. The default is
// bson.DefaultRegistry.
func (b *Builder) Registry(registry *bsoncodec.Registry) *Builder {
	b.opts.SetRegistry(registry)
	return b
}

// ReplicaSet specifies the replica set name for the cluster. If specified, the cluster will be treated as a replica
// set and the driver will automatically discover all servers in the set, starting with the nodes specified through
// ApplyURI or Hosts. All nodes in the replica set must have the same replica set name, or they will not be
// considered as part of the set by the Client. This can also be set through the "replica" URI option (e.g.
// "replica=replset"). The default is empty.
func (b *Builder) ReplicaSet(s string) *Builder {
	b.opts.SetReplicaSet(s)
	return b
}

// DisableRetryWrites specifies whether supported write operations should be retried once on certain errors, such as network
// errors.
//
// Supported operations are InsertOne, UpdateOne, ReplaceOne, DeleteOne, FindOneAndDelete, FindOneAndReplace,
// FindOneAndDelete, InsertMany, and BulkWrite. Note that BulkWrite requests must not include UpdateManyModel or
// DeleteManyModel instances to be considered retryable. Unacknowledged writes will not be retried, even if this option
// is set to true.
//
// This option requires server version >= 3.6 and a replica set or sharded cluster and will be ignored for any other
// cluster type. This can also be set through the "retryWrites" URI option (e.g. "retryWrites=true"). The default is
// true.
func (b *Builder) DisableRetryWrites() *Builder {
	b.opts.SetRetryWrites(false)
	return b
}

// DisableRetryReads specifies whether supported read operations should be retried once on certain errors, such as network
// errors.
//
// Supported operations are Find, FindOne, Aggregate without a $out stage, Distinct, CountDocuments,
// EstimatedDocumentCount, Watch (for Client, Database, and Collection), ListCollections, and ListDatabases. Note that
// operations run through RunCommand are not retried.
//
// This option requires server version >= 3.6 and driver version >= 1.1.0. The default is true.
func (b *Builder) DisableRetryReads() *Builder {
	b.opts.SetRetryReads(false)
	return b
}

// ServerSelectionTimeout specifies how long the driver will wait to find an available, suitable server to execute an
// operation. This can also be set through the "serverSelectionTimeoutMS" URI option (e.g.
// "serverSelectionTimeoutMS=30000"). The default value is 30 seconds.
func (b *Builder) ServerSelectionTimeout(d time.Duration) *Builder {
	b.opts.SetServerSelectionTimeout(d)
	return b
}

// SocketTimeout specifies how long the driver will wait for a socket read or write to return before returning a
// network error. This can also be set through the "socketTimeoutMS" URI option (e.g. "socketTimeoutMS=1000"). The
// default value is 0, meaning no timeout is used and socket operations can block indefinitely.
//
// NOTE(benjirewis): SocketTimeout will be deprecated in a future release. The more general Timeout option may be used
// in its place to control the amount of time that a single operation can run before returning an error. ting
// SocketTimeout and Timeout on a single client will result in undefined behavior.
func (b *Builder) SocketTimeout(d time.Duration) *Builder {
	b.opts.SetSocketTimeout(d)
	return b
}

// Timeout specifies the amount of time that a single operation run on this Client can execute before returning an error.
// The deadline of any operation run through the Client will be honored above any Timeout set on the Client; Timeout will only
// be honored if there is no deadline on the operation Context. Timeout can also be set through the "timeoutMS" URI option
// (e.g. "timeoutMS=1000"). The default value is nil, meaning operations do not inherit a timeout from the Client.
//
// If any Timeout is set (even 0) on the Client, the values of MaxTime on operation options, TransactionOptions.MaxCommitTime and
// SessionOptions.DefaultMaxCommitTime will be ignored. ting Timeout and SocketTimeout or WriteConcern.wTimeout will result
// in undefined behavior.
//
// NOTE(benjirewis): Timeout represents unstable, provisional API. The behavior of the driver when a Timeout is specified is
// subject to change.
func (b *Builder) Timeout(d time.Duration) *Builder {
	b.opts.SetTimeout(d)
	return b
}

// TLSConfig specifies a tls.Config instance to use use to configure TLS on all connections created to the cluster.
// This can also be set through the following URI options:
//
// 1. "tls" (or "ssl"): Specify if TLS should be used (e.g. "tls=true").
//
// 2. Either "tlsCertificateKeyFile" (or "sslClientCertificateKeyFile") or a combination of "tlsCertificateFile" and
// "tlsPrivateKeyFile". The "tlsCertificateKeyFile" option specifies a path to the client certificate and private key,
// which must be concatenated into one file. The "tlsCertificateFile" and "tlsPrivateKey" combination specifies separate
// paths to the client certificate and private key, respectively. Note that if "tlsCertificateKeyFile" is used, the
// other two options must not be specified. Only the subject name of the first certificate is honored as the username
// for X509 auth in a file with multiple certs.
//
// 3. "tlsCertificateKeyFilePassword" (or "sslClientCertificateKeyPassword"): Specify the password to decrypt the client
// private key file (e.g. "tlsCertificateKeyFilePassword=password").
//
// 4. "tlsCaFile" (or "sslCertificateAuthorityFile"): Specify the path to a single or bundle of certificate authorities
// to be considered trusted when making a TLS connection (e.g. "tlsCaFile=/path/to/caFile").
//
// 5. "tlsInsecure" (or "sslInsecure"): Specifies whether or not certificates and hostnames received from the server
// should be validated. If true (e.g. "tlsInsecure=true"), the TLS library will accept any certificate presented by the
// server and any host name in that certificate. Note that setting this to true makes TLS susceptible to
// man-in-the-middle attacks and should only be done for testing.
//
// The default is nil, meaning no TLS will be enabled.
func (b *Builder) TLSConfig(cfg *tls.Config) *Builder {
	b.opts.SetTLSConfig(cfg)
	return b
}

// HTTPClient specifies the http.Client to be used for any HTTP requests.
//
// This should only be used to set custom HTTP client configurations. By default, the connection will use an internal.DefaultHTTPClient.
func (b *Builder) HTTPClient(client *http.Client) *Builder {
	b.opts.SetHTTPClient(client)
	return b
}

// WriteConcern specifies the write concern to use to for write operations. This can also be set through the following
// URI options:
//
// 1. "w": Specify the number of nodes in the cluster that must acknowledge write operations before the operation
// returns or "majority" to specify that a majority of the nodes must acknowledge writes. This can either be an integer
// (e.g. "w=10") or the string "majority" (e.g. "w=majority").
//
// 2. "wTimeoutMS": Specify how long write operations should wait for the correct number of nodes to acknowledge the
// operation (e.g. "wTimeoutMS=1000").
//
// 3. "journal": Specifies whether or not write operations should be written to an on-disk journal on the server before
// returning (e.g. "journal=true").
//
// The default is nil, meaning the server will use its configured default.
func (b *Builder) WriteConcern(wc *writeconcern.WriteConcern) *Builder {
	b.opts.SetWriteConcern(wc)
	return b
}

// ZlibLevel specifies the level for the zlib compressor. This option is ignored if zlib is not specified as a
// compressor through ApplyURI or Compressors. Supported values are -1 through 9, inclusive. -1 tells the zlib
// library to use its default, 0 means no compression, 1 means best speed, and 9 means best compression.
// This can also be set through the "zlibCompressionLevel" URI option (e.g. "zlibCompressionLevel=-1"). Defaults to -1.
func (b *Builder) ZlibLevel(level int) *Builder {
	b.opts.SetZlibLevel(level)
	return b
}

// ZstdLevel sets the level for the zstd compressor. This option is ignored if zstd is not specified as a compressor
// through ApplyURI or Compressors. Supported values are 1 through 20, inclusive. 1 means best speed and 20 means
// best compression. This can also be set through the "zstdCompressionLevel" URI option. Defaults to 6.
func (b *Builder) ZstdLevel(level int) *Builder {
	b.opts.SetZstdLevel(level)
	return b
}

// AutoEncryptionOptions specifies an AutoEncryptionOptions instance to automatically encrypt and decrypt commands
// and their results. See the options.AutoEncryptionOptions documentation for more information about the supported
// options.
func (b *Builder) AutoEncryptionOptions(opts *rawopts.AutoEncryptionOptions) *Builder {
	b.opts.SetAutoEncryptionOptions(opts)
	return b
}

// DisableOCSPEndpointCheck specifies whether or not the driver should reach out to OCSP responders to verify the
// certificate status for certificates presented by the server that contain a list of OCSP responders.
//
// If set to true, the driver will verify the status of the certificate using a response stapled by the server, if there
// is one, but will not send an HTTP request to any responders if there is no staple. In this case, the driver will
// continue the connection even though the certificate status is not known.
//
// This can also be set through the tlsDisableOCSPEndpointCheck URI option. Both this URI option and tlsInsecure must
// not be set at the same time and will error if they are. The default value is false.
func (b *Builder) DisableOCSPEndpointCheck() *Builder {
	b.opts.SetDisableOCSPEndpointCheck(true)
	return b
}

// ServerAPIOptions specifies a ServerAPIOptions instance used to configure the API version sent to the server
// when running commands. See the options.ServerAPIOptions documentation for more information about the supported
// options.
func (b *Builder) ServerAPIOptions(opts *rawopts.ServerAPIOptions) *Builder {
	b.opts.SetServerAPIOptions(opts)
	return b
}

// SRVMaxHosts specifies the maximum number of SRV results to randomly select during polling. To limit the number
// of hosts selected in SRV discovery, this function must be called before ApplyURI. This can also be set through
// the "srvMaxHosts" URI option.
func (b *Builder) SRVMaxHosts(srvMaxHosts int) *Builder {
	b.opts.SetSRVMaxHosts(srvMaxHosts)
	return b
}

// SRVServiceName specifies a custom SRV service name to use in SRV polling. To use a custom SRV service name
// in SRV discovery, this function must be called before ApplyURI. This can also be set through the "srvServiceName"
// URI option.
func (b *Builder) SRVServiceName(srvName string) *Builder {
	b.opts.SetSRVServiceName(srvName)
	return b
}

func (b *Builder) Build(ctx context.Context) (*Client, error) {
	uri, err := connstring.ParseAndValidate(b.opts.GetURI())
	if err != nil {
		return nil, err
	}
	b.defaultDatabaseName = uri.Database
	if b.defaultDatabaseName == "" {
		b.defaultDatabaseName = "admin"
	}
	return newClient(ctx, b.pingReadpref, b.defaultDatabaseName, b.opts)
}

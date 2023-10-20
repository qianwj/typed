package transaction

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

type TxSessionBuilder struct {
	cli  *mongo.Client
	opts *options.SessionOptions
	bind func(txc *TxSession) error
}

func NewTxSessionBuilder(cli *mongo.Client) *TxSessionBuilder {
	return &TxSessionBuilder{
		cli:  cli,
		opts: options.Session(),
	}
}

// CausalConsistency sets the value for the CausalConsistency field.
func (tx *TxSessionBuilder) CausalConsistency() *TxSessionBuilder {
	tx.opts.SetCausalConsistency(true)
	return tx
}

// DefaultReadConcern sets the value for the DefaultReadConcern field.
func (tx *TxSessionBuilder) DefaultReadConcern(rc *readconcern.ReadConcern) *TxSessionBuilder {
	tx.opts.SetDefaultReadConcern(rc)
	return tx
}

// DefaultReadPreference sets the value for the DefaultReadPreference field.
func (tx *TxSessionBuilder) DefaultReadPreference(rp *readpref.ReadPref) *TxSessionBuilder {
	tx.opts.SetDefaultReadPreference(rp)
	return tx
}

// DefaultWriteConcern sets the value for the DefaultWriteConcern field.
func (tx *TxSessionBuilder) DefaultWriteConcern(wc *writeconcern.WriteConcern) *TxSessionBuilder {
	tx.opts.SetDefaultWriteConcern(wc)
	return tx
}

// DefaultMaxCommitTime sets the value for the DefaultMaxCommitTime field.
//
// NOTE(benjirewis): DefaultMaxCommitTime will be deprecated in a future release. The more
// general Timeout option may be used in its place to control the amount of time that a
// single operation can run before returning an error. DefaultMaxCommitTime is ignored if
// Timeout is set on the client.
func (tx *TxSessionBuilder) DefaultMaxCommitTime(mct *time.Duration) *TxSessionBuilder {
	tx.opts.SetDefaultMaxCommitTime(mct)
	return tx
}

// Snapshot sets the value for the Snapshot field.
func (tx *TxSessionBuilder) Snapshot() *TxSessionBuilder {
	tx.opts.SetSnapshot(true)
	return tx
}

func (tx *TxSessionBuilder) Bind(fn func(txc *TxSession) error) *TxSessionBuilder {
	tx.bind = fn
	return tx
}

func (tx *TxSessionBuilder) Execute(ctx context.Context) error {
	session, err := tx.cli.StartSession(tx.opts)
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	err = mongo.WithSession(ctx, session, func(stx mongo.SessionContext) error {
		return tx.bind(&TxSession{
			ctx:  stx,
			opts: options.Transaction(),
		})
	})
	return err
}

type TxSession struct {
	ctx  mongo.SessionContext
	opts *options.TransactionOptions
}

// ReadConcern sets the value for the ReadConcern field.
func (t *TxSession) ReadConcern(rc *readconcern.ReadConcern) *TxSession {
	t.opts.SetReadConcern(rc)
	return t
}

// ReadPreference sets the value for the ReadPreference field.
func (t *TxSession) ReadPreference(rp *readpref.ReadPref) *TxSession {
	t.opts.SetReadPreference(rp)
	return t
}

// WriteConcern sets the value for the WriteConcern field.
func (t *TxSession) WriteConcern(wc *writeconcern.WriteConcern) *TxSession {
	t.opts.SetWriteConcern(wc)
	return t
}

// MaxCommitTime sets the value for the MaxCommitTime field.
//
// NOTE(benjirewis): MaxCommitTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can run before
// returning an error. MaxCommitTime is ignored if Timeout is set on the client.
func (t *TxSession) MaxCommitTime(mct *time.Duration) *TxSession {
	t.opts.SetMaxCommitTime(mct)
	return t
}

func (t *TxSession) Start() error {
	return t.ctx.StartTransaction(t.opts)
}

func (t *TxSession) Abort(ctx context.Context) error {
	return t.ctx.AbortTransaction(ctx)
}

func (t *TxSession) Commit(ctx context.Context) error {
	return t.ctx.CommitTransaction(ctx)
}

func (t *TxSession) EndSession(ctx context.Context) {
	t.ctx.EndSession(ctx)
}

func (t *TxSession) Context() mongo.SessionContext {
	return t.ctx
}

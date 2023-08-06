package executor

import (
	"context"
	"errors"
	"github.com/qianwj/typed/mongo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	ErrEmptyIndex = errors.New("mongo: create index is empty")
)

type IndexViewer struct {
	view mongo.IndexView
}

func (i *IndexViewer) Create(idx ...*model.Index) *IndexCreateBuilder {
	return &IndexCreateBuilder{
		idx:  idx,
		view: i.view,
		opts: options.CreateIndexes(),
	}
}

func (i *IndexViewer) DropOne(name string) *IndexDropBuilder {
	return &IndexDropBuilder{
		name: name,
		view: i.view,
		opts: options.DropIndexes(),
	}
}

func (i *IndexViewer) DropAll() *IndexDropBuilder {
	return &IndexDropBuilder{
		all:  true,
		view: i.view,
		opts: options.DropIndexes(),
	}
}

func (i *IndexViewer) List() *IndexListBuilder {
	return &IndexListBuilder{
		view: i.view,
		opts: options.ListIndexes(),
	}
}

type IndexCreateBuilder struct {
	idx  []*model.Index
	view mongo.IndexView
	opts *options.CreateIndexesOptions
}

// MaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (c *IndexCreateBuilder) MaxTime(d time.Duration) *IndexCreateBuilder {
	c.opts.SetMaxTime(d)
	return c
}

// CommitQuorumInt sets the value for the CommitQuorum field as an int32.
func (c *IndexCreateBuilder) CommitQuorumInt(quorum int32) *IndexCreateBuilder {
	c.opts.SetCommitQuorumInt(quorum)
	return c
}

// CommitQuorumString sets the value for the CommitQuorum field as a string.
func (c *IndexCreateBuilder) CommitQuorumString(quorum string) *IndexCreateBuilder {
	c.opts.SetCommitQuorumString(quorum)
	return c
}

// CommitQuorumMajority sets the value for the CommitQuorum to special "majority" value.
func (c *IndexCreateBuilder) CommitQuorumMajority() *IndexCreateBuilder {
	c.opts.SetCommitQuorumString("majority")
	return c
}

// CommitQuorumVotingMembers sets the value for the CommitQuorum to special "votingMembers" value.
func (c *IndexCreateBuilder) CommitQuorumVotingMembers() *IndexCreateBuilder {
	c.opts.SetCommitQuorumString("votingMembers")
	return c
}

func (c *IndexCreateBuilder) Execute(ctx context.Context) ([]string, error) {
	if len(c.idx) > 1 {
		models := make([]mongo.IndexModel, len(c.idx))
		for i, idx := range c.idx {
			models[i] = idx.Marshal()
		}
		return c.view.CreateMany(ctx, models, c.opts)
	} else if len(c.idx) == 0 {
		return nil, ErrEmptyIndex
	}
	name, err := c.view.CreateOne(ctx, c.idx[0].Marshal(), c.opts)
	if err != nil {
		return nil, err
	}
	return []string{name}, nil
}

type IndexDropBuilder struct {
	name string
	all  bool
	view mongo.IndexView
	opts *options.DropIndexesOptions
}

func (d *IndexDropBuilder) MaxTime(duration time.Duration) *IndexDropBuilder {
	d.opts.SetMaxTime(duration)
	return d
}

func (d *IndexDropBuilder) Execute(ctx context.Context) (bson.Raw, error) {
	if !d.all {
		return d.view.DropOne(ctx, d.name, d.opts)
	}
	return d.view.DropAll(ctx, d.opts)
}

type IndexListBuilder struct {
	view mongo.IndexView
	opts *options.ListIndexesOptions
}

// BatchSize sets the value for the BatchSize field.
func (l *IndexListBuilder) BatchSize(i int32) *IndexListBuilder {
	l.opts.SetBatchSize(i)
	return l
}

// SetMaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (l *IndexListBuilder) SetMaxTime(d time.Duration) *IndexListBuilder {
	l.opts.SetMaxTime(d)
	return l
}

func (l *IndexListBuilder) Execute(ctx context.Context) ([]*mongo.IndexModel, error) {
	cursor, err := l.view.List(ctx, l.opts)
	if err != nil {
		return nil, err
	}
	var data []*mongo.IndexModel
	if err := cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (l *IndexListBuilder) ExecuteSpecifications(ctx context.Context) ([]*mongo.IndexSpecification, error) {
	return l.view.ListSpecifications(ctx, l.opts)
}

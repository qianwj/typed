package requests

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type UpdateRequestBuilder struct {
	internal *update.Request
}

func Update() *UpdateRequestBuilder {
	return &UpdateRequestBuilder{internal: update.NewRequest()}
}

// DetectNoop Set to false to disable setting 'result' in the response
// to 'noop' if no change to the document occurred.
func (u *UpdateRequestBuilder) DetectNoop(detectNoop bool) *UpdateRequestBuilder {
	u.internal.DetectNoop = &detectNoop
	return u
}

// Doc A partial update to an existing document.
func (u *UpdateRequestBuilder) Doc(doc json.RawMessage) *UpdateRequestBuilder {
	u.internal.Doc = doc
	return u
}

// DocAsUpsert Set to true to use the contents of 'doc' as the value of 'upsert'
func (u *UpdateRequestBuilder) DocAsUpsert(docAsUpsert bool) *UpdateRequestBuilder {
	u.internal.DocAsUpsert = &docAsUpsert
	return u
}

// Script Script to execute to update the document.
func (u *UpdateRequestBuilder) Script(script types.Script) *UpdateRequestBuilder {
	u.internal.Script = script
	return u
}

// ScriptedUpsert Set to true to execute the script whether or not the document exists.
func (u *UpdateRequestBuilder) ScriptedUpsert(scriptedUpsert bool) *UpdateRequestBuilder {
	u.internal.ScriptedUpsert = &scriptedUpsert
	return u
}

// Source Set to false to disable source retrieval. You can also specify a
// comma-separated
// list of the fields you want to retrieve.
func (u *UpdateRequestBuilder) Source(config types.SourceConfig) *UpdateRequestBuilder {
	u.internal.Source_ = config
	return u
}

// Upsert If the document does not already exist, the contents of 'upsert' are inserted
// as a
// new document. If the document exists, the 'script' is executed.
func (u *UpdateRequestBuilder) Upsert(upsert json.RawMessage) *UpdateRequestBuilder {
	u.internal.Upsert = upsert
	return u
}

func (u *UpdateRequestBuilder) Build() *update.Request {
	return u.internal
}

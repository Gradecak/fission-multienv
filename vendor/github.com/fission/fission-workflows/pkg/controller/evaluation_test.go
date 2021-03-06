package controller

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var dummyRecord = EvalRecord{
	Timestamp: time.Now(),
	Error:     errors.New("other error"),
	RulePath:  []string{},
}

func TestEvalLog_Append(t *testing.T) {
	var log EvalLog
	record1 := dummyRecord
	log.Record(record1)
	record2 := EvalRecord{
		Timestamp: time.Now(),
		Error:     errors.New("stub error"),
		RulePath:  []string{"a", "b", "c"},
	}
	log.Record(record2)

	assert.Equal(t, 2, log.Len())
	last, ok := log.Last()
	assert.True(t, ok)
	assert.EqualValues(t, record2, last)

}

func TestEvalState_Lock(t *testing.T) {
	es := NewEvalState("id", nil)

	<-es.Lock()
	select {
	case <-es.Lock():
		assert.Fail(t, "test was able to lock EvalState twice.")
	default:
		// ok
	}
}

func TestEvalState_Free(t *testing.T) {
	es := NewEvalState("id", nil)

	es.Free() // Idempotent

	select {
	case <-es.Lock():
		// ok
	default:
		assert.Fail(t, "test failed to lock free EvalState twice.")
	}
	es.Free()
	select {
	case <-es.Lock():
		// ok
	default:
		assert.Fail(t, "test failed to lock free EvalState twice.")
	}
}

func TestEvalState_First(t *testing.T) {
	es := NewEvalState("id", nil)
	assert.Equal(t, "id", es.ID())

	// Test non-existent
	r, ok := es.First()
	assert.False(t, ok)
	assert.Equal(t, EvalRecord{}, r)

	// Test existent
	es.Record(dummyRecord)
	r, ok = es.First()
	assert.True(t, ok)
	assert.Equal(t, dummyRecord, r)
}

func TestEvalState_Last(t *testing.T) {
	es := NewEvalState("id", nil)
	assert.Equal(t, "id", es.ID())

	// Test non-existent
	r, ok := es.Last()
	assert.False(t, ok)
	assert.Equal(t, EvalRecord{}, r)

	// Test existent
	es.Record(dummyRecord)
	r, ok = es.Last()
	assert.True(t, ok)
	assert.Equal(t, dummyRecord, r)
}

func TestEvalState_Count(t *testing.T) {
	es := NewEvalState("id", nil)
	assert.Equal(t, "id", es.ID())

	c := es.Len()
	assert.Equal(t, 0, c)

	es.Record(dummyRecord)
	c = es.Len()
	assert.Equal(t, 1, c)
}

func TestEvalState_Logs(t *testing.T) {
	es := NewEvalState("id", nil)
	assert.Equal(t, "id", es.ID())

	// Test non-existent
	el := es.Logs()
	assert.Equal(t, EvalLog{}, el)

	// Test existent
	es.Record(dummyRecord)
	el = es.Logs()
	assert.Equal(t, EvalLog{dummyRecord}, el)
}

func TestEvalCache_GetOrCreate(t *testing.T) {
	ec := EvalStore{}
	id := "foo"
	es, ok := ec.Load(id)
	assert.False(t, ok)
	assert.Empty(t, es)

	es, _ = ec.LoadOrStore(id, nil)
	assert.Equal(t, id, es.ID())

	es, ok = ec.Load(id)
	assert.True(t, ok)
	assert.Equal(t, id, es.ID())
}

func TestEvalCache_Invalidate(t *testing.T) {
	ec := EvalStore{}
	id := "completedId"

	ec.Store(NewEvalState(id, nil))
	es, ok := ec.Load(id)
	assert.True(t, ok)
	assert.Equal(t, id, es.ID())

	ec.Delete(id)
	es, ok = ec.Load(id)
	assert.False(t, ok)
	assert.Empty(t, es)
}

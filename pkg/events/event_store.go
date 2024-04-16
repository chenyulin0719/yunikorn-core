/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package events

import (
	"go.uber.org/zap"

	"github.com/apache/yunikorn-core/pkg/locking"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/metrics"
	"github.com/apache/yunikorn-scheduler-interface/lib/go/si"
)

// The EventStore operates under the following assumptions:
//   - there is a cap for the number of events stored
//   - the CollectEvents() function clears the currently stored events in the EventStore
//
// Assuming the rate of events generated by the scheduler component in a given time period
// is high, calling CollectEvents() periodically should be fine.
type EventStore struct {
	events   []*si.EventRecord
	idx      uint64 // points where to store the next event
	size     uint64
	lastSize uint64
	locking.RWMutex
}

func newEventStore(size uint64) *EventStore {
	return &EventStore{
		events: make([]*si.EventRecord, size),
		size:   size,
	}
}

func (es *EventStore) Store(event *si.EventRecord) {
	es.Lock()
	defer es.Unlock()

	if es.idx == uint64(len(es.events)) {
		metrics.GetEventMetrics().IncEventsNotStored()
		return
	}
	es.events[es.idx] = event
	es.idx++

	metrics.GetEventMetrics().IncEventsStored()
}

func (es *EventStore) CollectEvents() []*si.EventRecord {
	es.Lock()
	defer es.Unlock()

	messages := make([]*si.EventRecord, len(es.events[:es.idx]))
	copy(messages, es.events[:es.idx])

	if es.size != es.lastSize {
		log.Log(log.Events).Info("Resizing event store", zap.Uint64("last", es.lastSize), zap.Uint64("new", es.size))
		es.events = make([]*si.EventRecord, es.size)
	}
	es.idx = 0
	es.lastSize = es.size

	metrics.GetEventMetrics().AddEventsCollected(len(messages))
	return messages
}

func (es *EventStore) CountStoredEvents() uint64 {
	es.RLock()
	defer es.RUnlock()

	return es.idx
}

func (es *EventStore) SetStoreSize(size uint64) {
	es.Lock()
	defer es.Unlock()
	es.size = size
}

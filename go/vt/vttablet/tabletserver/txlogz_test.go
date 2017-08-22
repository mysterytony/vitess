/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tabletserver

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/youtube/vitess/go/sync2"
	"github.com/youtube/vitess/go/vt/callerid"
	"github.com/youtube/vitess/go/vt/vttablet/tabletserver/tabletenv"
)

func testHandler(req *http.Request, t *testing.T) {
	response := httptest.NewRecorder()
	tabletenv.TxLogger.Send("test msg")
	txlogzHandler(response, req)

	if !strings.Contains(response.Body.String(), "Redacted") {
		t.Fatalf("should have been redacted")
	}

	// skip the rest of the test since it is now always redacted
	return

	if !strings.Contains(response.Body.String(), "error") {
		t.Fatalf("should show an error page since transaction log format is invalid.")
	}
	txConn := &TxConnection{
		TransactionID:     123456,
		StartTime:         time.Now(),
		Queries:           []string{"select * from test"},
		Conclusion:        "unknown",
		LogToFile:         sync2.AtomicInt32{},
		EffectiveCallerID: callerid.NewEffectiveCallerID("effective-caller", "component", "subcomponent"),
		ImmediateCallerID: callerid.NewImmediateCallerID("immediate-caller"),
	}
	txConn.EndTime = txConn.StartTime
	response = httptest.NewRecorder()
	tabletenv.TxLogger.Send(txConn)
	txlogzHandler(response, req)
	txConn.EndTime = txConn.StartTime.Add(time.Duration(2) * time.Second)
	response = httptest.NewRecorder()
	tabletenv.TxLogger.Send(txConn)
	txlogzHandler(response, req)
	txConn.EndTime = txConn.StartTime.Add(time.Duration(500) * time.Millisecond)
	response = httptest.NewRecorder()
	tabletenv.TxLogger.Send(txConn)
	txlogzHandler(response, req)

}

func TestTxlogzHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/txlogz?timeout=0&limit=10", nil)
	testHandler(req, t)
}

func TestTxlogzHandlerWithNegativeTimeout(t *testing.T) {
	req, _ := http.NewRequest("GET", "/txlogz?timeout=-1&limit=10", nil)
	testHandler(req, t)
}

func TestTxlogzHandlerWithLargeLimit(t *testing.T) {
	req, _ := http.NewRequest("GET", "/txlogz?timeout=0&limit=10000000", nil)
	testHandler(req, t)
}

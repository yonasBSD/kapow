/*
 * Copyright 2019 Banco Bilbao Vizcaya Argentaria, S.A.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package httperror_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BBVA/kapow/internal/server/httperror"
)

func TestErrorJSONSetsAppJsonContentType(t *testing.T) {
	w := httptest.NewRecorder()

	httperror.ErrorJSON(w, "Not Important Here", 500)

	if v := w.Result().Header.Get("Content-Type"); v != "application/json; charset=utf-8" {
		t.Errorf("Content-Type header mismatch. Expected: %q, got: %q", "application/json; charset=utf-8", v)
	}
}

func TestErrorJSONSetsRequestedStatusCode(t *testing.T) {
	w := httptest.NewRecorder()

	httperror.ErrorJSON(w, "Not Important Here", http.StatusGone)

	if v := w.Result().StatusCode; v != http.StatusGone {
		t.Errorf("Status code mismatch. Expected: %d, got: %d", http.StatusGone, v)
	}
}

func TestErrorJSONSetsBodyCorrectly(t *testing.T) {
	expectedReason := "Something Not Found"
	w := httptest.NewRecorder()

	httperror.ErrorJSON(w, expectedReason, http.StatusNotFound)

	errMsg := httperror.ServerErrMessage{}
	if bodyBytes, err := io.ReadAll(w.Result().Body); err != nil {
		t.Errorf("Unexpected error reading response body: %v", err)
	} else if err := json.Unmarshal(bodyBytes, &errMsg); err != nil {
		t.Errorf("Response body contains invalid JSON entity: %v", err)
	} else if errMsg.Reason != expectedReason {
		t.Errorf("Unexpected reason in response. Expected: %q, got: %q", expectedReason, errMsg.Reason)
	}
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/flow-hydraulics/flow-wallet-api/errors"
	"github.com/flow-hydraulics/flow-wallet-api/templates"
	"github.com/flow-hydraulics/flow-wallet-api/transactions"
	"github.com/gorilla/mux"
)

func (s *Transactions) ListFunc(rw http.ResponseWriter, r *http.Request) {
	var (
		transactionSlice []transactions.Transaction
		err              error
	)

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		limit = 0
	}

	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		offset = 0
	}

	vars := mux.Vars(r)

	if address, ok := vars["address"]; ok {
		// Handle account specific transactions
		// This endpoint is used to handle "raw" transactions for an account
		// so we use transactions.General type here
		transactionSlice, err = s.service.ListForAccount(transactions.General, address, limit, offset)
	} else {
		// Handle all transactions
		transactionSlice, err = s.service.List(limit, offset)
	}

	if err != nil {
		handleError(rw, s.log, err)
		return
	}

	res := make([]transactions.JSONResponse, len(transactionSlice))
	for i, job := range transactionSlice {
		res[i] = job.ToJSONResponse()
	}

	handleJsonResponse(rw, http.StatusOK, res)
}

func (s *Transactions) CreateFunc(rw http.ResponseWriter, r *http.Request) {
	var err error

	if r.Body == nil || r.Body == http.NoBody {
		err = &errors.RequestError{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("empty body"),
		}
		handleError(rw, s.log, err)
		return
	}

	vars := mux.Vars(r)

	var b templates.Raw

	// Try to decode the request body into the struct.
	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		err = &errors.RequestError{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("invalid body"),
		}
		handleError(rw, s.log, err)
		return
	}

	// Decide whether to serve sync or async, default async
	sync := r.FormValue(SyncQueryParameter) != ""
	job, transaction, err := s.service.Create(r.Context(), sync, vars["address"], b, transactions.General)

	if err != nil {
		handleError(rw, s.log, err)
		return
	}

	var res interface{}
	if sync {
		res = transaction.ToJSONResponse()
	} else {
		res = job.ToJSONResponse()
	}

	handleJsonResponse(rw, http.StatusCreated, res)
}

func (s *Transactions) DetailsFunc(rw http.ResponseWriter, r *http.Request) {
	var (
		transaction *transactions.Transaction
		err         error
	)
	vars := mux.Vars(r)

	if address, ok := vars["address"]; ok {
		// Handle account specific transactions
		// This endpoint is used to handle "raw" transactions for an account
		// so we use transactions.General type here
		transaction, err = s.service.DetailsForAccount(r.Context(), transactions.General, address, vars["transactionId"])
	} else {
		// Handle all transactions
		transaction, err = s.service.Details(r.Context(), vars["transactionId"])
	}

	if err != nil {
		handleError(rw, s.log, err)
		return
	}

	res := transaction.ToJSONResponse()

	handleJsonResponse(rw, http.StatusOK, res)
}

func (s *Transactions) ExecuteScriptFunc(rw http.ResponseWriter, r *http.Request) {
	var err error

	if r.Body == nil || r.Body == http.NoBody {
		err = &errors.RequestError{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("empty body"),
		}
		handleError(rw, s.log, err)
		return
	}

	var b templates.Raw

	// Try to decode the request body into the struct.
	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		err = &errors.RequestError{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("invalid body"),
		}
		handleError(rw, s.log, err)
		return
	}

	res, err := s.service.ExecuteScript(r.Context(), b)

	if err != nil {
		handleError(rw, s.log, err)
		return
	}

	handleJsonResponse(rw, http.StatusOK, res)
}

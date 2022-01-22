package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/denniswon/reddio/module/ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/denniswon/reddio/api/model"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func getLatestBlock(service ethereum.Module) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set our response header
		w.Header().Set("Content-Type", "application/json")

		errorMessage := "Error getting latest block"
		block := service.GetLatestBlock()

		if err := json.NewEncoder(w).Encode(block); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func getTxByHash(service ethereum.Module) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set our response header
		w.Header().Set("Content-Type", "application/json")

		errorMessage := "Error getting transaction by hash"

		vars := mux.Vars(r)
		hash := vars["hash"]

		if hash == "" {
			json.NewEncoder(w).Encode(&model.Error{
				Code:    400,
				Message: "Malformed request",
			})
			return
		}
		txHash := common.HexToHash(hash)
		tx := service.GetTxByHash(txHash)

		if tx == nil {
			json.NewEncoder(w).Encode(&model.Error{
				Code:    404,
				Message: "Tx Not Found",
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(tx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

func getAddressBalance(service ethereum.Module) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set our response header
		w.Header().Set("Content-Type", "application/json")

		errorMessage := "Error getting address balance"

		vars := mux.Vars(r)
		address := vars["address"]

		if address == "" {
			json.NewEncoder(w).Encode(&model.Error{
				Code:    400,
				Message: "Malformed request",
			})
			return
		}

		balance, err := service.GetAddressBalance(address)

		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(&model.Error{
				Code:    500,
				Message: "Internal server error",
			})
			return
		}

		toJ := &model.BalanceResponse{
			Address: address,
			Balance: balance,
			Symbol:  "Ether",
			Units:   "Wei",
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func transferEth(service ethereum.Module) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set our response header
		w.Header().Set("Content-Type", "application/json")

		errorMessage := "Error transferring eth"

		decoder := json.NewDecoder(r.Body)
		var t model.TransferEthRequest

		err := decoder.Decode(&t)

		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(&model.Error{
				Code:    400,
				Message: "Malformed request",
			})
			return
		}
		_hash, err := service.TransferEth(t.PrivKey, t.To, t.Amount)

		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(&model.Error{
				Code:    500,
				Message: "Internal server error",
			})
			return
		}

		toJ := &model.HashResponse{
			Hash: _hash,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

//MakeEthereumHandlers make url handlers
func MakeEthereumHandlers(r *mux.Router, n negroni.Negroni, service ethereum.Module) {
	r.Handle("/v1/ethereum/block/latest", n.With(
		negroni.Wrap(getLatestBlock(service)),
	)).Methods("GET", "OPTIONS").Name("getLatestBlock")

	r.Handle("/v1/ethereum/tx/{hash}", n.With(
		negroni.Wrap(getTxByHash(service)),
	)).Methods("GET", "OPTIONS").Name("getTxByHash")

	r.Handle("/v1/ethereum/address/{address}", n.With(
		negroni.Wrap(getAddressBalance(service)),
	)).Methods("GET", "OPTIONS").Name("getAddressBalance")

	r.Handle("/v1/ethereum/transferEth", n.With(
		negroni.Wrap(transferEth(service)),
	)).Methods("POST", "OPTIONS").Name("transferEth")
}

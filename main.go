package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

var bank atomic.Int64
var mtx = sync.Mutex{}

func payHandler(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		msg := "failed to read r.Body" + err.Error()
		fmt.Println(msg)
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("error while writing http response", err)
			return
		}
		return
	}

	paymentAmount, err := strconv.ParseInt(string(requestBody), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		msg := "error while parsing r.Body" + err.Error()
		fmt.Println(msg)
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("error while writing http response", err)
			return
		}
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	if bank.Load()-paymentAmount >= 0 {
		bank.Add(-paymentAmount)
		fmt.Println("payment succeded!")
		fmt.Println("bank: ", bank.Load())

		msg := "payment succeeded, current balance: " + strconv.FormatInt(bank.Load(), 10)

		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("error while writing http response", err)
			return
		}
	}
}

func main() {
	bank.Add(1000)
	http.HandleFunc("/pay", payHandler)

	http.ListenAndServe(":9091", nil)
}

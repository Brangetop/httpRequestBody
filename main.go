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
	mtx.Unlock()
}

func main() {
	bank.Add(1000)
	http.HandleFunc("/pay", payHandler)

	http.ListenAndServe(":9091", nil)
}

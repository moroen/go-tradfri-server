package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	coap "github.com/moroen/go-tradfricoap"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "Welcome!")
}

func GetLights(w http.ResponseWriter, r *http.Request) {
	answer := returnMessageDevices{Action: "getLights"}

	lights, err := coap.GetDevices()
	if err != nil {
		if err == coap.ErrorTimeout {
			answer.Status = "Timeout"
		} else {
			log.Fatal(err.Error())
		}
	} else {
		answer = returnMessageDevices{Action: "getLights", Status: "Ok", Result: lights}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(answer)
}

func GetLight(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Println("Getting light: ", params["id"])

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err.Error())
	}

	if device, err := coap.GetLight(int64(id)); err == nil {
		answer := returnMessageDevice{Action: "getLight", Status: "Ok", Result: device}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(answer)
	} else {
		answer := returnMessageDevice{Action: "getLight"}
		if err == coap.ErrorTimeout {
			answer.Status = "Timeout"
		} else {
			answer.Status = "Unknown error"
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(answer)
	}
}

func SetState(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	state := 0

	answer := returnMessageDevice{Action: "setState"}

	log.Println("SetState: ", params["id"], params["command"])
	if params["command"] == "on" {
		state = 1
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err.Error())
	}

	device, err := coap.SetState(int64(id), state)
	if err != nil {
		if err == coap.ErrorTimeout {
			answer.Status = "Timeout"
		}
	} else {
		answer.Status = "Ok"
		answer.Result = device
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(answer)
}

func SetDimmer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	answer := returnMessageDevice{Action: "setLevel"}

	log.Println("SetDimmer: ", params["id"], params["value"])

	if value, err := strconv.Atoi(params["value"]); err == nil {
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			panic(err.Error())
		}
		device, err := coap.SetLevel(int64(id), value)
		if err != nil {
			answer.Status = "Error: Failed to set level"
		} else {
			answer.Status = "Ok"
			answer.Result = device
		}
	} else {
		log.Println("Failed to set level")
		answer.Status = "Error: Failed to set level"
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(answer)
}

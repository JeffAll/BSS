package handlers

import (
	"bss/go/data"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type ItemHandler struct {
	Data *data.Data
}

func (ih *ItemHandler) HandleQuery(
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Println("HandleQuery")
	items, err := ih.Data.Query()
	if err != nil {
		log.Printf(
			"Error Querying Data\n\t:%s",
			err,
		)
	}
	log.Printf(
		"Items\n\t:%s",
		items,
	)
	toWrite, err := json.Marshal(items)
	if err != nil {
		log.Printf(
			"Error Marshalling Value\n\t:%s",
			err,
		)
	}
	log.Printf(
		"ToWrite\n\t:%s",
		string(toWrite),
	)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(toWrite)
}

func (ih *ItemHandler) HandleUpdate(
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Println("HandleUpdate")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf(
			"Error Reading Request\n\t:%s",
			err,
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var items data.Items
	err = json.Unmarshal(
		body,
		&items,
	)
	if err != nil {
		log.Printf(
			"Error Unmarshaling Data\n\t:%s\n\tbody:%s",
			err,
			string(body),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = ih.Data.AppendArrayAndClear(
		items,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

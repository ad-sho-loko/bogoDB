package db

import (
	"log"
	"net/http"
)

type ApiServer struct {
	db *BogoDb
}

func NewApiServer(db *BogoDb) *ApiServer {
	return &ApiServer{
		db: db,
	}
}

func (a *ApiServer) executeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/execute requested")
	log.Println(r.URL.Query())

	q := r.URL.Query()["query"]
	if len(q) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := a.db.Execute(q[0], r.UserAgent())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
}

func (a *ApiServer) exitHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/exit requested")
	a.db.Terminate()
}

func (a *ApiServer) Host() {
	http.HandleFunc("/execute", a.executeHandler)
	http.HandleFunc("/exit", a.exitHandler)
	log.Fatal(http.ListenAndServe(":32198", nil))
}

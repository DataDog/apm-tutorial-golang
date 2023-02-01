package calendar

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
	//"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func GetDate(w http.ResponseWriter, r *http.Request) {
	//span, _ := tracer.StartSpanFromContext(r.Context(), "GetDate")
	//defer span.Finish()

	val := rand.Intn(365)
	date := time.Now().AddDate(0, 0, val)
	ranDate, marshError := json.Marshal(date.Format("2006-01-02"))
	if marshError != nil {
		log.Fatalln(marshError)
	}
	log.Println(string(ranDate))

	w.Header().Set("Content-Type", "application/json")
	w.Write(ranDate)
}

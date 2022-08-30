package calendar

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func GetDate(w http.ResponseWriter, r *http.Request) {

	sctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(r.Header))
	if err != nil {
		log.Fatalln(err)
	}

	span := tracer.StartSpan("getDate", tracer.ChildOf(sctx))
	defer span.Finish()

	val := rand.Intn(365)
	date := time.Now().AddDate(0, 0, val)
	ranDate := date.Format("2006-01-02")
	log.Println(ranDate)

	calDoLongRunningProcess(span)
	calAnotherProcess(span)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(ranDate))
}

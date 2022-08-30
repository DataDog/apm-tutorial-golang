package calendar

import (
	"log"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

//tracing annotations go here if possible
func calDoLongRunningProcess(span tracer.Span) {

	childSpan := tracer.StartSpan("traceMethod1", tracer.ChildOf(span.Context()))
	childSpan.SetTag(ext.ResourceName, "CalendarHelper.doLongRunningProcess")
	defer childSpan.Finish()

	time.Sleep(300 * time.Millisecond)
	log.Println("Hello from the long running process in Calendar")
	calPrivateMethod1(childSpan)
}

func calAnotherProcess(span tracer.Span) {

	childSpan := tracer.StartSpan("traceMethod2", tracer.ChildOf(span.Context()))
	childSpan.SetTag(ext.ResourceName, "CalendarHelper.anotherProcess")
	defer childSpan.Finish()

	time.Sleep(50 * time.Millisecond)
	log.Println("Hello from the anotherProcess in Calendar")
}

func calPrivateMethod1(span tracer.Span) {
	childSpan := tracer.StartSpan("manualSpan1", tracer.ChildOf(span.Context()))
	childSpan.SetTag(ext.ResourceName, "calPrivateMethod1")
	defer childSpan.Finish()

	time.Sleep(30 * time.Millisecond)
	log.Println("Hello from the custom privateMethod1 in Calendar")
}

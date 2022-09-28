package notes

import (
	"context"
	"log"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func doLongRunningProcess(ctx context.Context) {
	childSpan, ctx := tracer.StartSpanFromContext(ctx, "traceMethod1")
	childSpan.SetTag(ext.ResourceName, "NotesHelper.doLongRunningProcess")
	defer childSpan.Finish()

	time.Sleep(300 * time.Millisecond)
	log.Println("Hello from the long running process in Notes")
	privateMethod1(ctx)
}

func anotherProcess(ctx context.Context) {
	childSpan, _ := tracer.StartSpanFromContext(ctx, "traceMethod2")
	childSpan.SetTag(ext.ResourceName, "NotesHelper.anotherProcess")
	defer childSpan.Finish()

	time.Sleep(50 * time.Millisecond)
	log.Println("Hello from the anotherProcess in Notes")
}

func privateMethod1(ctx context.Context) {
	childSpan, _ := tracer.StartSpanFromContext(ctx, "manualSpan1",
		tracer.SpanType("web"),
		tracer.ServiceName("noteshelper"),
	)
	childSpan.SetTag(ext.ResourceName, "privateMethod1")
	defer childSpan.Finish()

	time.Sleep(30 * time.Millisecond)
	log.Println("Hello from the custom privateMethod1 in Notes")
}

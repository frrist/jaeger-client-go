package jaeger

import opentracing "github.com/opentracing/opentracing-go"

// Observer can be registered with the Tracer to receive notifications about new Spans.
type Observer interface {
	OnStartSpan(operationName string, options opentracing.StartSpanOptions) SpanObserver
}

// SpanObserver is created by the Observer and receives notifications about other Span events.
type SpanObserver interface {
	OnSetOperationName(operationName string)
	OnSetTag(key string, value interface{})
	OnFinish(options opentracing.FinishOptions)
}

// observer is a dispatcher to other observers
type observer struct {
	observers []Observer
}

// spanObserver is a dispatcher to other span observers
type spanObserver struct {
	observers []SpanObserver
}

// noopSpanObserver is used when there are no observers registered on the Tracer
// or none of them returns span observers from OnStartSpan.
var noopSpanObserver = spanObserver{}

func (o *observer) append(observer Observer) {
	o.observers = append(o.observers, observer)
}

func (o observer) OnStartSpan(operationName string, options opentracing.StartSpanOptions) SpanObserver {
	var spanObservers []SpanObserver
	for _, obs := range o.observers {
		spanObs := obs.OnStartSpan(operationName, options)
		if spanObs != nil {
			if spanObservers == nil {
				spanObservers = make([]SpanObserver, 0, len(o.observers))
			}
			spanObservers = append(spanObservers, spanObs)
		}
	}
	if len(spanObservers) == 0 {
		return noopSpanObserver
	}
	return spanObserver{observers: spanObservers}
}

func (o spanObserver) OnSetOperationName(operationName string) {
	for _, obs := range o.observers {
		obs.OnSetOperationName(operationName)
	}
}

func (o spanObserver) OnSetTag(key string, value interface{}) {
	for _, obs := range o.observers {
		obs.OnSetTag(key, value)
	}
}

func (o spanObserver) OnFinish(options opentracing.FinishOptions) {
	for _, obs := range o.observers {
		obs.OnFinish(options)
	}
}
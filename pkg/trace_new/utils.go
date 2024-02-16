/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

// Package trace provides functions and constants for tracing.
package trace_new

import (

	// "github.com/uber/jaeger-client-go"

	"go.opentelemetry.io/otel/exporters/jaeger"

	// jaegercfg "github.com/uber/jaeger-client-go/config"
	// "github.com/uber/jaeger-client-go/log/zap"
	// "github.com/uber/jaeger-lib/metrics"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func NewDefaultTracer(tracerName string) trace.Tracer {
	return otel.Tracer(tracerName)
}

func InitProvider() (*otelsdk.TracerProvider, error) {
	var err error
	tp := &otelsdk.TracerProvider{}
	// tp, err = initOtelTracer(cfg)

	tp, err = initJaegerTracerProvider("http://localhost:14268/api/traces")

	// otel.SetTracerProvider(tp)
	// otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, err
}

func initJaegerTracerProvider(jaegerCollectorEndpoint string) (*otelsdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerCollectorEndpoint)))
	if err != nil {
		return nil, err
	}

	tp := otelsdk.NewTracerProvider(
		// Always be sure to batch in production.
		otelsdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		otelsdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
		)),
	)

	return tp, nil
}

// // NewConstTracer returns an instance of Jaeger Tracer that samples 100% or 0% of traces for test.
// func NewConstTracer(serviceName string, collect bool) (opentracing.Tracer, io.Closer, error) {
// 	if serviceName == "" {
// 		return nil, nil, fmt.Errorf("service name is empty")
// 	}
// 	param := 0.0
// 	if collect {
// 		param = 1.0
// 	}
// 	// Sample configuration for testing. Use constant sampling to sample every trace
// 	// and enable LogSpan to log every span via configured Logger.
// 	cfg := jaegercfg.Configuration{
// 		ServiceName: serviceName,
// 		Sampler: &jaegercfg.SamplerConfig{
// 			Type:  jaeger.SamplerTypeConst,
// 			Param: param,
// 		},
// 		Reporter: &jaegercfg.ReporterConfig{
// 			LogSpans: true,
// 		},
// 	}
// 	return cfg.NewTracer(
// 		jaegercfg.Logger(zap.NewLogger(log.Log(log.OpenTracing).Named(serviceName))),
// 		jaegercfg.Metrics(metrics.NullFactory),
// 	)
// }

// // NewTracerFromEnv returns an instance of Jaeger Tracer that get sampling strategy from env settings.
// func NewTracerFromEnv(serviceName string) (opentracing.Tracer, io.Closer, error) {
// 	cfg, err := jaegercfg.FromEnv()
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	if serviceName != "" {
// 		cfg.ServiceName = serviceName
// 	}
// 	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
// 	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
// 	// frameworks.
// 	// Initialize tracer with a logger and a metrics factory
// 	return cfg.NewTracer(
// 		jaegercfg.Logger(zap.NewLogger(log.Log(log.OpenTracing).Named(cfg.ServiceName))),
// 		jaegercfg.Metrics(metrics.NullFactory),
// 	)
// }

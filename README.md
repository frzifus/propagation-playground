# Play

Run the root service:
```bash
OTEL_TRACES_EXPORTER=none go run cmd/service1/main.go
# output
initialising OpenTelemetry tracer
Do request, correlationID: 7f64af13-0208-486f-a162-108d406e675d, requestID: 0
something went wrong, correlationID: 7f64af13-0208-486f-a162-108d406e675d, requestID: 0, err: Get "http://localhost:8080": EOF
```

Example of outgoing request
```http
GET / HTTP/1.1
Host: localhost:8080
User-Agent: Go-http-client/1.1
Baggage: correlationID=7f64af13-0208-486f-a162-108d406e675d,requestID=0
Traceparent: 00-bacdfe5f74ee2360575648818575566b-de0098a1153fe5c5-01
Accept-Encoding: gzip
```

Using `OTEL_TRACES_EXPORTER` a tracer can be defined; supported values:
  - "none" - "no operation" exporter
  - "console" - Standard output exporter; see [go.opentelemetry.io/otel/exporters/stdout/stdouttrace]
  - "otlp" (default) - OTLP exporter; see [go.opentelemetry.io/otel/exporters/otlp/otlptrace]

Example logging spans to console:
```bash
OTEL_TRACES_EXPORTER=console go run cmd/service1/main.go
# output
initialising OpenTelemetry tracer
Do request, correlationID: 7f64af13-0208-486f-a162-108d406e675d, requestID: 0
something went wrong, correlationID: 7f64af13-0208-486f-a162-108d406e675d, requestID: 0, err: Get "http://localhost:8080": EOF
{"Name":"HTTP GET","SpanContext":{"TraceID":"bacdfe5f74ee2360575648818575566b","SpanID":"de0098a1153fe5c5","TraceFlags":"01","TraceState":"","Remote":false},"Parent":{"TraceID":"bacdfe5f74ee2360575648818575566b","SpanID":"9b96920bd702bd38","TraceFlags":"01","TraceState":"","Remote":false},"SpanKind":3,"StartTime":"2025-01-14T19:32:31.659422787+01:00","EndTime":"2025-01-14T19:32:33.966030653+01:00","Attributes":[{"Key":"http.method","Value":{"Type":"STRING","Value":"GET"}},{"Key":"http.url","Value":{"Type":"STRING","Value":"http://localhost:8080"}},{"Key":"net.peer.name","Value":{"Type":"STRING","Value":"localhost"}},{"Key":"net.peer.port","Value":{"Type":"INT64","Value":8080}}],"Events":[{"Name":"exception","Attributes":[{"Key":"exception.type","Value":{"Type":"STRING","Value":"*errors.errorString"}},{"Key":"exception.message","Value":{"Type":"STRING","Value":"EOF"}}],"DroppedAttributeCount":0,"Time":"2025-01-14T19:32:33.966029017+01:00"}],"Links":null,"Status":{"Code":"Error","Description":"EOF"},"DroppedAttributes":0,"DroppedEvents":0,"DroppedLinks":0,"ChildSpanCount":0,"Resource":[{"Key":"host.name","Value":{"Type":"STRING","Value":"aqua"}},{"Key":"service.name","Value":{"Type":"STRING","Value":"info.APPName"}},{"Key":"service.version","Value":{"Type":"STRING","Value":"info.Version"}}],"InstrumentationScope":{"Name":"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp","Version":"0.58.0","SchemaURL":"","Attributes":null},"InstrumentationLibrary":{"Name":"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp","Version":"0.58.0","SchemaURL":"","Attributes":null}}
{"Name":"origin","SpanContext":{"TraceID":"bacdfe5f74ee2360575648818575566b","SpanID":"9b96920bd702bd38","TraceFlags":"01","TraceState":"","Remote":false},"Parent":{"TraceID":"00000000000000000000000000000000","SpanID":"0000000000000000","TraceFlags":"00","TraceState":"","Remote":false},"SpanKind":1,"StartTime":"2025-01-14T19:32:31.659402917+01:00","EndTime":"2025-01-14T19:32:33.966066381+01:00","Attributes":null,"Events":null,"Links":null,"Status":{"Code":"Unset","Description":""},"DroppedAttributes":0,"DroppedEvents":0,"DroppedLinks":0,"ChildSpanCount":1,"Resource":[{"Key":"host.name","Value":{"Type":"STRING","Value":"aqua"}},{"Key":"service.name","Value":{"Type":"STRING","Value":"info.APPName"}},{"Key":"service.version","Value":{"Type":"STRING","Value":"info.Version"}}],"InstrumentationScope":{"Name":"github.com/frzifus/propagation-playground/cmd/service1","Version":"","SchemaURL":"","Attributes":null},"InstrumentationLibrary":{"Name":"github.com/frzifus/propagation-playground/cmd/service1","Version":"","SchemaURL":"","Attributes":null}}
```

### Root Span

```json
{
  "Name": "origin",
  "SpanContext": {
    "TraceID": "bacdfe5f74ee2360575648818575566b",
    "SpanID": "9b96920bd702bd38",
    "TraceFlags": "01",
    "TraceState": "",
    "Remote": false
  },
  "Parent": {
    "TraceID": "00000000000000000000000000000000",
    "SpanID": "0000000000000000",
    "TraceFlags": "00",
    "TraceState": "",
    "Remote": false
  },
  "SpanKind": 1,
  "StartTime": "2025-01-14T19:32:31.659402917+01:00",
  "EndTime": "2025-01-14T19:32:33.966066381+01:00",
  "Attributes": null,
  "Events": null,
  "Links": null,
  "Status": {
    "Code": "Unset",
    "Description": ""
  },
  "DroppedAttributes": 0,
  "DroppedEvents": 0,
  "DroppedLinks": 0,
  "ChildSpanCount": 1,
  "Resource": [
    {
      "Key": "host.name",
      "Value": {
        "Type": "STRING",
        "Value": "aqua"
      }
    },
    {
      "Key": "service.name",
      "Value": {
        "Type": "STRING",
        "Value": "info.APPName"
      }
    },
    {
      "Key": "service.version",
      "Value": {
        "Type": "STRING",
        "Value": "info.Version"
      }
    }
  ],
  "InstrumentationScope": {
    "Name": "github.com/frzifus/propagation-playground/cmd/service1",
    "Version": "",
    "SchemaURL": "",
    "Attributes": null
  },
  "InstrumentationLibrary": {
    "Name": "github.com/frzifus/propagation-playground/cmd/service1",
    "Version": "",
    "SchemaURL": "",
    "Attributes": null
  }
}
```

### Next span

```json
{
  "Name": "HTTP GET",
  "SpanContext": {
    "TraceID": "bacdfe5f74ee2360575648818575566b",
    "SpanID": "de0098a1153fe5c5",
    "TraceFlags": "01",
    "TraceState": "",
    "Remote": false
  },
  "Parent": {
    "TraceID": "bacdfe5f74ee2360575648818575566b",
    "SpanID": "9b96920bd702bd38",
    "TraceFlags": "01",
    "TraceState": "",
    "Remote": false
  },
  "SpanKind": 3,
  "StartTime": "2025-01-14T19:32:31.659422787+01:00",
  "EndTime": "2025-01-14T19:32:33.966030653+01:00",
  "Attributes": [
    {
      "Key": "http.method",
      "Value": {
        "Type": "STRING",
        "Value": "GET"
      }
    },
    {
      "Key": "http.url",
      "Value": {
        "Type": "STRING",
        "Value": "http://localhost:8080"
      }
    },
    {
      "Key": "net.peer.name",
      "Value": {
        "Type": "STRING",
        "Value": "localhost"
      }
    },
    {
      "Key": "net.peer.port",
      "Value": {
        "Type": "INT64",
        "Value": 8080
      }
    }
  ],
  "Events": [
    {
      "Name": "exception",
      "Attributes": [
        {
          "Key": "exception.type",
          "Value": {
            "Type": "STRING",
            "Value": "*errors.errorString"
          }
        },
        {
          "Key": "exception.message",
          "Value": {
            "Type": "STRING",
            "Value": "EOF"
          }
        }
      ],
      "DroppedAttributeCount": 0,
      "Time": "2025-01-14T19:32:33.966029017+01:00"
    }
  ],
  "Links": null,
  "Status": {
    "Code": "Error",
    "Description": "EOF"
  },
  "DroppedAttributes": 0,
  "DroppedEvents": 0,
  "DroppedLinks": 0,
  "ChildSpanCount": 0,
  "Resource": [
    {
      "Key": "host.name",
      "Value": {
        "Type": "STRING",
        "Value": "aqua"
      }
    },
    {
      "Key": "service.name",
      "Value": {
        "Type": "STRING",
        "Value": "info.APPName"
      }
    },
    {
      "Key": "service.version",
      "Value": {
        "Type": "STRING",
        "Value": "info.Version"
      }
    }
  ],
  "InstrumentationScope": {
    "Name": "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
    "Version": "0.58.0",
    "SchemaURL": "",
    "Attributes": null
  },
  "InstrumentationLibrary": {
    "Name": "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
    "Version": "0.58.0",
    "SchemaURL": "",
    "Attributes": null
  }
}
```


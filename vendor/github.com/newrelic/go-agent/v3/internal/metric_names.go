package internal

const (
	apdexRollup = "Apdex"
	apdexPrefix = "Apdex/"

	webRollup        = "WebTransaction"
	backgroundRollup = "OtherTransaction/all"

	// https://source.datanerd.us/agents/agent-specs/blob/master/Total-Time-Async.md
	totalTimeWeb        = "WebTransactionTotalTime"
	totalTimeBackground = "OtherTransactionTotalTime"

	errorsPrefix = "Errors/"

	// "HttpDispatcher" metric is used for the overview graph, and
	// therefore should only be made for web transactions.
	dispatcherMetric = "HttpDispatcher"

	queueMetric = "WebFrontend/QueueTime"

	webMetricPrefix        = "WebTransaction/Go"
	backgroundMetricPrefix = "OtherTransaction/Go"

	instanceReporting = "Instance/Reporting"

	// https://newrelic.atlassian.net/wiki/display/eng/Custom+Events+in+New+Relic+Agents
	customEventsSeen = "Supportability/Events/Customer/Seen"
	customEventsSent = "Supportability/Events/Customer/Sent"

	// https://source.datanerd.us/agents/agent-specs/blob/master/Transaction-Events-PORTED.md
	txnEventsSeen = "Supportability/AnalyticsEvents/TotalEventsSeen"
	txnEventsSent = "Supportability/AnalyticsEvents/TotalEventsSent"

	// https://source.datanerd.us/agents/agent-specs/blob/master/Error-Events.md
	errorEventsSeen = "Supportability/Events/TransactionError/Seen"
	errorEventsSent = "Supportability/Events/TransactionError/Sent"

	// https://source.datanerd.us/agents/agent-specs/blob/master/Span-Events.md
	spanEventsSeen = "Supportability/SpanEvent/TotalEventsSeen"
	spanEventsSent = "Supportability/SpanEvent/TotalEventsSent"

	supportabilityDropped = "Supportability/MetricsDropped"

	// Runtime/System Metrics
	memoryPhysical       = "Memory/Physical"
	heapObjectsAllocated = "Memory/Heap/AllocatedObjects"
	cpuUserUtilization   = "CPU/User/Utilization"
	cpuSystemUtilization = "CPU/System/Utilization"
	cpuUserTime          = "CPU/User Time"
	cpuSystemTime        = "CPU/System Time"
	runGoroutine         = "Go/Runtime/Goroutines"
	gcPauseFraction      = "GC/System/Pause Fraction"
	gcPauses             = "GC/System/Pauses"

	// Configurable event harvest supportability metrics
	supportReportPeriod     = "Supportability/EventHarvest/ReportPeriod"
	supportTxnEventLimit    = "Supportability/EventHarvest/AnalyticEventData/HarvestLimit"
	supportCustomEventLimit = "Supportability/EventHarvest/CustomEventData/HarvestLimit"
	supportErrorEventLimit  = "Supportability/EventHarvest/ErrorEventData/HarvestLimit"
	supportSpanEventLimit   = "Supportability/EventHarvest/SpanEventData/HarvestLimit"
)

// DistributedTracingSupport is used to track distributed tracing activity for
// supportability.
type DistributedTracingSupport struct {
	// New Relic DT fields
	AcceptPayloadSuccess            bool // AcceptPayload was called successfully
	AcceptPayloadException          bool // AcceptPayload had a generic exception
	AcceptPayloadParseException     bool // AcceptPayload had a parsing exception
	AcceptPayloadCreateBeforeAccept bool // AcceptPayload was ignored because CreatePayload had already been called
	AcceptPayloadIgnoredMultiple    bool // AcceptPayload was ignored because AcceptPayload had already been called
	AcceptPayloadIgnoredVersion     bool // AcceptPayload was ignored because the payload's major version was greater than the agent's
	AcceptPayloadUntrustedAccount   bool // AcceptPayload was ignored because the payload was untrusted
	AcceptPayloadNullPayload        bool // AcceptPayload was ignored because the payload was nil
	CreatePayloadSuccess            bool // CreatePayload was called successfully
	CreatePayloadException          bool // CreatePayload had a generic exception

	// W3C Trace Context fields
	TraceContextAcceptSuccess        bool // The agent successfully accepted inbound traceparent and tracestate headers.
	TraceContextAcceptException      bool // A generic exception occurred unrelated to parsing while accepting either payload.
	TraceContextParentParseException bool // The inbound traceparent header could not be parsed.
	TraceContextStateParseException  bool // The inbound tracestate header could not be parsed.
	TraceContextStateInvalidNrEntry  bool // The inbound tracestate header exists, and was accepted, but the New Relic entry was invalid.
	TraceContextStateNoNrEntry       bool // The traceparent header exists, and was accepted, but the tracestate header did not contain a trusted New Relic entry.
	TraceContextCreateSuccess        bool // The agent successfully created the outbound payloads.
	TraceContextCreateException      bool // A generic exception occurred while creating the outbound payloads.
}

func (dts DistributedTracingSupport) isEmpty() bool {
	return (DistributedTracingSupport{}) == dts
}

func supportMetric(metrics *metricTable, b bool, metricName string) {
	if b {
		metrics.addSingleCount(metricName, forced)
	}
}

func (dts DistributedTracingSupport) createMetrics(ms *metricTable) {
	// Distributed Tracing Supportability Metrics
	supportMetric(ms, dts.AcceptPayloadSuccess, "Supportability/DistributedTrace/AcceptPayload/Success")
	supportMetric(ms, dts.AcceptPayloadException, "Supportability/DistributedTrace/AcceptPayload/Exception")
	supportMetric(ms, dts.AcceptPayloadParseException, "Supportability/DistributedTrace/AcceptPayload/ParseException")
	supportMetric(ms, dts.AcceptPayloadCreateBeforeAccept, "Supportability/DistributedTrace/AcceptPayload/Ignored/CreateBeforeAccept")
	supportMetric(ms, dts.AcceptPayloadIgnoredMultiple, "Supportability/DistributedTrace/AcceptPayload/Ignored/Multiple")
	supportMetric(ms, dts.AcceptPayloadIgnoredVersion, "Supportability/DistributedTrace/AcceptPayload/Ignored/MajorVersion")
	supportMetric(ms, dts.AcceptPayloadUntrustedAccount, "Supportability/DistributedTrace/AcceptPayload/Ignored/UntrustedAccount")
	supportMetric(ms, dts.AcceptPayloadNullPayload, "Supportability/DistributedTrace/AcceptPayload/Ignored/Null")
	supportMetric(ms, dts.CreatePayloadSuccess, "Supportability/DistributedTrace/CreatePayload/Success")
	supportMetric(ms, dts.CreatePayloadException, "Supportability/DistributedTrace/CreatePayload/Exception")

	// W3C Trace Context Supportability Metrics
	supportMetric(ms, dts.TraceContextAcceptSuccess, "Supportability/TraceContext/Accept/Success")
	supportMetric(ms, dts.TraceContextAcceptException, "Supportability/TraceContext/Accept/Exception")
	supportMetric(ms, dts.TraceContextParentParseException, "Supportability/TraceContext/TraceParent/Parse/Exception")
	supportMetric(ms, dts.TraceContextStateParseException, "Supportability/TraceContext/TraceState/Parse/Exception")
	supportMetric(ms, dts.TraceContextCreateSuccess, "Supportability/TraceContext/Create/Success")
	supportMetric(ms, dts.TraceContextCreateException, "Supportability/TraceContext/Create/Exception")
	supportMetric(ms, dts.TraceContextStateInvalidNrEntry, "Supportability/TraceContext/TraceState/InvalidNrEntry")
	supportMetric(ms, dts.TraceContextStateNoNrEntry, "Supportability/TraceContext/TraceState/NoNrEntry")
}

type rollupMetric struct {
	all      string
	allWeb   string
	allOther string
}

func newRollupMetric(s string) rollupMetric {
	return rollupMetric{
		all:      s + "all",
		allWeb:   s + "allWeb",
		allOther: s + "allOther",
	}
}

func (r rollupMetric) webOrOther(isWeb bool) string {
	if isWeb {
		return r.allWeb
	}
	return r.allOther
}

var (
	errorsRollupMetric = newRollupMetric("Errors/")

	// source.datanerd.us/agents/agent-specs/blob/master/APIs/external_segment.md
	// source.datanerd.us/agents/agent-specs/blob/master/APIs/external_cat.md
	// source.datanerd.us/agents/agent-specs/blob/master/Cross-Application-Tracing-PORTED.md
	externalRollupMetric = newRollupMetric("External/")

	// source.datanerd.us/agents/agent-specs/blob/master/Datastore-Metrics-PORTED.md
	datastoreRollupMetric = newRollupMetric("Datastore/")

	datastoreProductMetricsCache = map[string]rollupMetric{
		"Cassandra":     newRollupMetric("Datastore/Cassandra/"),
		"Derby":         newRollupMetric("Datastore/Derby/"),
		"Elasticsearch": newRollupMetric("Datastore/Elasticsearch/"),
		"Firebird":      newRollupMetric("Datastore/Firebird/"),
		"IBMDB2":        newRollupMetric("Datastore/IBMDB2/"),
		"Informix":      newRollupMetric("Datastore/Informix/"),
		"Memcached":     newRollupMetric("Datastore/Memcached/"),
		"MongoDB":       newRollupMetric("Datastore/MongoDB/"),
		"MySQL":         newRollupMetric("Datastore/MySQL/"),
		"MSSQL":         newRollupMetric("Datastore/MSSQL/"),
		"Oracle":        newRollupMetric("Datastore/Oracle/"),
		"Postgres":      newRollupMetric("Datastore/Postgres/"),
		"Redis":         newRollupMetric("Datastore/Redis/"),
		"Solr":          newRollupMetric("Datastore/Solr/"),
		"SQLite":        newRollupMetric("Datastore/SQLite/"),
		"CouchDB":       newRollupMetric("Datastore/CouchDB/"),
		"Riak":          newRollupMetric("Datastore/Riak/"),
		"VoltDB":        newRollupMetric("Datastore/VoltDB/"),
	}
)

func customSegmentMetric(s string) string {
	return "Custom/" + s
}

// customMetric is used to construct custom metrics from the input given to
// Application.RecordCustomMetric.  Note that the "Custom/" prefix helps prevent
// collision with other agent metrics, but does not eliminate the possibility
// since "Custom/" is also used for segments.
func customMetric(customerInput string) string {
	return "Custom/" + customerInput
}

// DatastoreMetricKey contains the fields by which datastore metrics are
// aggregated.
type DatastoreMetricKey struct {
	Product      string
	Collection   string
	Operation    string
	Host         string
	PortPathOrID string
}

type externalMetricKey struct {
	Host                    string
	Library                 string
	Method                  string
	ExternalCrossProcessID  string
	ExternalTransactionName string
}

// MessageMetricKey is the key to use for message segments.
type MessageMetricKey struct {
	Library         string
	DestinationType string
	Consumer        bool
	DestinationName string
	DestinationTemp bool
}

// Name returns the metric name value for this MessageMetricKey to be used for
// scoped and unscoped metrics.
//
// Producers
// MessageBroker/{Library}/{Destination Type}/{Action}/Named/{Destination Name}
// MessageBroker/{Library}/{Destination Type}/{Action}/Temp
//
// Consumers
// OtherTransaction/Message/{Library}/{DestinationType}/Named/{Destination Name}
// OtherTransaction/Message/{Library}/{DestinationType}/Temp
func (key MessageMetricKey) Name() string {
	var destination string
	if key.DestinationTemp {
		destination = "Temp"
	} else if key.DestinationName == "" {
		destination = "Named/Unknown"
	} else {
		destination = "Named/" + key.DestinationName
	}

	if key.Consumer {
		return "Message/" + key.Library +
			"/" + key.DestinationType +
			"/" + destination
	}
	return "MessageBroker/" + key.Library +
		"/" + key.DestinationType +
		"/Produce/" + destination
}

func datastoreScopedMetric(key DatastoreMetricKey) string {
	if "" != key.Collection {
		return datastoreStatementMetric(key)
	}
	return datastoreOperationMetric(key)
}

// Datastore/{datastore}/*
func datastoreProductMetric(key DatastoreMetricKey) rollupMetric {
	d, ok := datastoreProductMetricsCache[key.Product]
	if ok {
		return d
	}
	return newRollupMetric("Datastore/" + key.Product + "/")
}

// Datastore/operation/{datastore}/{operation}
func datastoreOperationMetric(key DatastoreMetricKey) string {
	return "Datastore/operation/" + key.Product +
		"/" + key.Operation
}

// Datastore/statement/{datastore}/{table}/{operation}
func datastoreStatementMetric(key DatastoreMetricKey) string {
	return "Datastore/statement/" + key.Product +
		"/" + key.Collection +
		"/" + key.Operation
}

// Datastore/instance/{datastore}/{host}/{port_path_or_id}
func datastoreInstanceMetric(key DatastoreMetricKey) string {
	return "Datastore/instance/" + key.Product +
		"/" + key.Host +
		"/" + key.PortPathOrID
}

func (key externalMetricKey) scopedMetric() string {
	if "" != key.ExternalCrossProcessID && "" != key.ExternalTransactionName {
		return externalTransactionMetric(key)
	}

	if key.Method == "" {
		// External/{host}/{library}
		return "External/" + key.Host + "/" + key.Library
	}
	// External/{host}/{library}/{method}
	return "External/" + key.Host + "/" + key.Library + "/" + key.Method
}

// External/{host}/all
func externalHostMetric(key externalMetricKey) string {
	return "External/" + key.Host + "/all"
}

// ExternalApp/{host}/{external_id}/all
func externalAppMetric(key externalMetricKey) string {
	return "ExternalApp/" + key.Host +
		"/" + key.ExternalCrossProcessID + "/all"
}

// ExternalTransaction/{host}/{external_id}/{external_txnname}
func externalTransactionMetric(key externalMetricKey) string {
	return "ExternalTransaction/" + key.Host +
		"/" + key.ExternalCrossProcessID +
		"/" + key.ExternalTransactionName
}

func callerFields(c payloadCaller) string {
	return "/" + c.Type +
		"/" + c.Account +
		"/" + c.App +
		"/" + c.TransportType +
		"/"
}

// DurationByCaller/{type}/{account}/{app}/{transport}/*
func durationByCallerMetric(c payloadCaller) rollupMetric {
	return newRollupMetric("DurationByCaller" + callerFields(c))
}

// ErrorsByCaller/{type}/{account}/{app}/{transport}/*
func errorsByCallerMetric(c payloadCaller) rollupMetric {
	return newRollupMetric("ErrorsByCaller" + callerFields(c))
}

// TransportDuration/{type}/{account}/{app}/{transport}/*
func transportDurationMetric(c payloadCaller) rollupMetric {
	return newRollupMetric("TransportDuration" + callerFields(c))
}
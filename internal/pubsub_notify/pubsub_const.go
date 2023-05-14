package pubsub

const TraceID = "trace_id"

type Mode string

const (
	ModeGoChannel    = Mode("gochannel")
	ModeGooglePubsub = Mode("googlepubsub")
)

type Logger string

const (
	LoggerNoop = Logger("noop")
	LoggerZap  = Logger("zap")
)

const (
	EnvDefaultMode     = "INFRA_PUB_SUB_DEFAULT_MODE"
	EnvLogger          = "INFRA_PUB_SUB_LOGGER"
	EnvGoogleProjectID = "INFRA_PUB_SUB_GOOGLE_PROJECT_ID"
	EnvPrefix          = "INFRA_PUB_SUB_PREFIX"
)

package api

const (
	// API Routes
	APIV1BasePath = "/api/v1"

	// System Routes
	HealthPath       = "/health"
	HealthPathSystem = "/health/system"
	SwaggerPath      = "/swagger/*any"

	// Sender Routes
	SenderBasePath      = "/sender"
	StartSchedulerPath  = "/startScheduler"
	StopSchedulerPath   = "/stopScheduler"
	StatusSchedulerPath = "/statusScheduler"

	// Message Routes
	MessagesBasePath = "/messages"
	SentMessagesPath = "/sent"

	// HTTP Headers
	HeaderAccessToken = "x-access-token"
	HeaderAuthKey     = "x-ins-auth-key"
)

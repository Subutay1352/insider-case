package constants

// Database Types
const (
	DBTypePostgres = "postgres"
	DBTypeSQLite   = "sqlite"
)

// Default Database Values
const (
	DefaultDBUser     = "postgres"
	DefaultDBPassword = "postgres"
	DefaultDBName     = "insider_case"
)

// API Routes
const (
	APIV1BasePath = "/api/v1"

	// System Routes
	HealthPath  = "/health"
	SwaggerPath = "/swagger/*any"

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

// Message Status
const (
	MessageStatusQueued     = "queued"
	MessageStatusProcessing = "processing"
	MessageStatusSent       = "sent"
	MessageStatusDelivered  = "delivered"
	MessageStatusFailed     = "failed"
	MessageStatusCancelled  = "cancelled"
)

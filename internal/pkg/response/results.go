package response

type ErrorCode string

const (
	ErrorCodeSchedulerAlreadyRunning  ErrorCode = "SCHEDULER_ALREADY_RUNNING"
	ErrorCodeSchedulerNotRunning      ErrorCode = "SCHEDULER_NOT_RUNNING"
	ErrorCodeSchedulerStartFailed     ErrorCode = "SCHEDULER_START_FAILED"
	ErrorCodeSchedulerStopFailed      ErrorCode = "SCHEDULER_STOP_FAILED"
	ErrorCodeFailedToRetrieveMessages ErrorCode = "FAILED_TO_RETRIEVE_MESSAGES"
	ErrorCodeUnauthorized             ErrorCode = "UNAUTHORIZED"
	ErrorCodeUnauthorizedMissingToken ErrorCode = "UNAUTHORIZED_MISSING_TOKEN"
	ErrorCodeUnauthorizedInvalidToken ErrorCode = "UNAUTHORIZED_INVALID_TOKEN"
	ErrorCodeInternalServerError      ErrorCode = "INTERNAL_SERVER_ERROR"
)

type SuccessCode string

const (
	SuccessCodeSchedulerStarted         SuccessCode = "SCHEDULER_STARTED"
	SuccessCodeSchedulerStopped         SuccessCode = "SCHEDULER_STOPPED"
	SuccessCodeSchedulerStatusRetrieved SuccessCode = "SCHEDULER_STATUS_RETRIEVED"
	SuccessCodeMessagesRetrieved        SuccessCode = "MESSAGES_RETRIEVED"
)

type ErrorResult struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

type SuccessResult struct {
	Code    SuccessCode `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

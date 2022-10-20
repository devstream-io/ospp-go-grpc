package plugin

type LogLevel int32

const (
	LogLevelDebug LogLevel = 0
	LogLevelInfo  LogLevel = 1
	LogLevelWarn  LogLevel = 2
	LogLevelError LogLevel = 3
)

type UnmountReason int32

const (
	UnmountExit      UnmountReason = 0
	UnmountUnUsed    UnmountReason = 1
	UnmountUpgrade   UnmountReason = 2
	UnmountDowngrade UnmountReason = 3
	UnmountPanic     UnmountReason = 4
)

type Status int32

const (
	StatusConnected    Status = 0
	StatusDisconnected Status = 1
)

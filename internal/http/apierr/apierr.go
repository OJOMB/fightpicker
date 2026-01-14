package apierr

import "github.com/OJOMB/fightpicker/pkg/logs"

type APIError struct {
	Status   int
	Code     string
	LogLevel logs.Level
	LogMsg   string
	Public   error
}

func NewAPIError(
	status int,
	code string,
	logLevel logs.Level,
	logMsg string,
	public error,
) APIError {
	return APIError{
		Status:   status,
		Code:     code,
		LogLevel: logLevel,
		LogMsg:   logMsg,
		Public:   public,
	}
}

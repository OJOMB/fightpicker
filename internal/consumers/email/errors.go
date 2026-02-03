package email

import "fmt"

// private initialization errors
var (
	// errNil logger is returned when a nil logger is provided.
	errNilLogger = fmt.Errorf("logger is required to init image consumer")
	// errNilKafkaClient is returned when a nil Kafka client is provided.
	errNilKafkaClient = fmt.Errorf("kafka client is required to init image consumer")
	// errNilUsersService is returned when a nil users service is provided.
	errNilUsersService = fmt.Errorf("users service is required to init image consumer")
	// errNilIDTool is returned when a nil ID tool is provided.
	errNilIDTool = fmt.Errorf("id tool is required to init image consumer")
)

// public runtime errors
var (
	// ErrInvalidRecordKeyFormat is returned when a record key does not match the expected format.
	ErrInvalidRecordKeyFormat = fmt.Errorf("invalid record key format")
)

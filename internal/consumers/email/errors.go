package email

import "fmt"

var (
	ErrNilLogger              = fmt.Errorf("logger is required to init image consumer")
	ErrNilKafkaClient         = fmt.Errorf("kafka client is required to init image consumer")
	ErrNilUsersService        = fmt.Errorf("users service is required to init image consumer")
	ErrInvalidRecordKeyFormat = fmt.Errorf("invalid record key format")
)

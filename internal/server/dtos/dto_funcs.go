package dtos

func NewErrorEnvelope(err error, code, reqID string) ErrorEnvelope {
	if err == nil {
		return ErrorEnvelope{}
	}

	return ErrorEnvelope{
		Error: ErrorObject{
			Code:      code,
			Message:   err.Error(),
			RequestId: reqID,
		},
	}
}

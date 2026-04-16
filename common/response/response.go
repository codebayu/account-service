package response

type Result struct {
	Code       int      `json:"code"`
	StatusCode int      `json:"statusCode"`
	Message    string   `json:"message"`
	Errors     []string `json:"errors,omitempty"`
}

type Response struct {
	Result Result      `json:"result"`
	Data   interface{} `json:"data"`
}

func Success(statusCode int, code int, message string, data interface{}) Response {
	return Response{
		Result: Result{
			Code:       code,
			StatusCode: statusCode,
			Message:    message,
		},
		Data: data,
	}
}

func Error(statusCode int, code int, message string, errors []string) Response {
	return Response{
		Result: Result{
			Code:       code,
			StatusCode: statusCode,
			Message:    message,
			Errors:     errors,
		},
		Data: nil,
	}
}

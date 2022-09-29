package utils

const successStatusText string = "success"
const errorStatusText string = "error"

// SwaggerGenericErrorResponse is store swagger generic error response
// swagger:response genericError
type SwaggerGenericErrorResponse struct {
	// in: body
	Body struct {
		// Required: true
		// Example: Error
		Status string `json:"status"`
		// Required: true
		// Example: 400
		StatusCode int `json:"status_code"`
		// Required: true
		// Example: Invalid value for x
		Error string `json:"error"`
	}
}

// SwaggerGenericSuccessResponse is store swagger generic response
// swagger:response genericResponse
type SwaggerGenericSuccessResponse struct {
	// in: body
	Body struct {
		// Required: true
		// Example: Success
		Status string `json:"status"`
		// Required: true
		// Example: 200
		StatusCode int `json:"status_code"`
		// Required: true
		// Example: ok
		Message string `json:"message"`
	}
}

// GenericErrorResponse is store generic error response related data
type GenericErrorResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

// GenericSuccessResponse is store generic success response related data
type GenericSuccessResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

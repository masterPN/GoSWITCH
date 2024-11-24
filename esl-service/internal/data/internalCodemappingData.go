package data

type InternalCodemappingData struct {
	ID           int `json:"ID"`
	InternalCode int `json:"InternalCode"`
	OperatorCode int `json:"OperatorCode"`
}

type InternalCodemappingDataError struct {
	Error                   string                  `json:"error"`
	InternalCodemappingData InternalCodemappingData `json:"internalCodemappingData"`
}

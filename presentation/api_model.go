package presentation

// APIModel is an interface for all API objects
type APIModel interface {
	ConvertToAPI(interface{}) error
}

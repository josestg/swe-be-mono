package system

// Status is the health status of the application.
// swagger:model system.Status
type Status string //@name system.Status

// Status constants.
const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
)

// HealthRes represents the health status of a service.
// swagger:model system.HealthRes
type HealthRes struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
} //@name system.HealthRes

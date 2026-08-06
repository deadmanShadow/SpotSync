package dto

// CreateZoneRequest is the payload for POST /api/v1/zones.
type CreateZoneRequest struct {
	Name          string  `json:"name" validate:"required,min=2,max=100"`
	Type          string  `json:"type" validate:"required,oneof=general ev_charging covered"`
	TotalCapacity int     `json:"total_capacity" validate:"required,gt=0"`
	PricePerHour  float64 `json:"price_per_hour" validate:"required,gt=0"`
}

// UpdateZoneRequest allows partial updates to a zone. All fields are optional;
// only provided fields are applied.
type UpdateZoneRequest struct {
	Name          string  `json:"name" validate:"omitempty,min=2,max=100"`
	Type          string  `json:"type" validate:"omitempty,oneof=general ev_charging covered"`
	TotalCapacity int     `json:"total_capacity" validate:"omitempty,gt=0"`
	PricePerHour  float64 `json:"price_per_hour" validate:"omitempty,gt=0"`
}
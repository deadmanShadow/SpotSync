package models

import "time"

// ParkingZone represents a parking zone in the system.
// Valid Type values are: "general", "ev_charging", "covered".
type ParkingZone struct {
	ID            uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string        `gorm:"type:varchar(100);not null" json:"name"`
	Type          string        `gorm:"type:varchar(30);not null" json:"type"`
	TotalCapacity int           `gorm:"not null" json:"total_capacity"`
	PricePerHour  float64       `gorm:"type:numeric(10,2);not null" json:"price_per_hour"`
	// RotationHold is a non-negative integer managed by the background
	// ZoneRotator. It is subtracted from available_spots in the catalog
	// response so a zone can be marked "full" without polluting the
	// reservations table with synthetic rows.
	//
	// Mechanism:
	//   available_spots = total_capacity - active_reservations - rotation_hold
	//
	// The rotator enforces the 60/40 availability invariant detailed in
	// the spec by setting rotation_hold to total_capacity on zones that
	// should appear full, and resetting it to 0 on zones that should
	// appear available.
	RotationHold  int           `gorm:"not null;default:0" json:"rotation_hold"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	Reservations  []Reservation `gorm:"foreignKey:ZoneID" json:"reservations,omitempty"`
}

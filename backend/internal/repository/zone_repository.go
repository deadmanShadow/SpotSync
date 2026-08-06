package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"spotsync-backend/internal/dto"
	"spotsync-backend/internal/models"
)

// ZoneRepository defines the persistence contract for ParkingZone records.
type ZoneRepository interface {
	Create(zone *models.ParkingZone) error
	FindByID(id uint) (*models.ParkingZone, error)
	List() ([]models.ParkingZone, error)
	Update(zone *models.ParkingZone) error
	Delete(id uint) error
	CountActiveReservations(zoneID uint) (int64, error)
	FindAllWithAvailability() ([]dto.ZoneResponse, error)
	FindByIDWithAvailability(id uint) (*dto.ZoneResponse, error)
}

type zoneRepository struct {
	db *gorm.DB
}

// NewZoneRepository wires a ZoneRepository backed by the given GORM DB.
func NewZoneRepository(db *gorm.DB) ZoneRepository {
	return &zoneRepository{db: db}
}

// parseTime accepts the various textual representations that the postgres
// driver can return for timestamp columns and parses them into time.Time.
// An empty string yields the zero time and a nil error.
func parseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.999999-07",
		"2006-01-02 15:04:05.999999",
		"2006-01-02 15:04:05-07",
		"2006-01-02 15:04:05",
	}
	var lastErr error
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		} else {
			lastErr = err
		}
	}
	if lastErr == nil {
		return time.Time{}, errors.New("unrecognized time format")
	}
	return time.Time{}, lastErr
}

func (r *zoneRepository) Create(zone *models.ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *zoneRepository) FindByID(id uint) (*models.ParkingZone, error) {
	var zone models.ParkingZone
	if err := r.db.First(&zone, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &zone, nil
}

func (r *zoneRepository) List() ([]models.ParkingZone, error) {
	var zones []models.ParkingZone
	if err := r.db.Order("id ASC").Find(&zones).Error; err != nil {
		return nil, err
	}
	return zones, nil
}

func (r *zoneRepository) Update(zone *models.ParkingZone) error {
	return r.db.Save(zone).Error
}

func (r *zoneRepository) Delete(id uint) error {
	return r.db.Delete(&models.ParkingZone{}, id).Error
}

// CountActiveReservations returns the count of active reservations for the
// given zone. Used to enforce capacity checks.
func (r *zoneRepository) CountActiveReservations(zoneID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Reservation{}).
		Where("zone_id = ? AND status = ?", zoneID, "active").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindAllWithAvailability returns every parking zone with a dynamically
// computed AvailableSpots value: total_capacity - active reservations.
//
// The availability calculation is performed in a single SQL query using a
// correlated subquery against the reservations table, so callers do not need
// to do an N+1 lookup.
func (r *zoneRepository) FindAllWithAvailability() ([]dto.ZoneResponse, error) {
	type row struct {
		ID             uint    `gorm:"column:id"`
		Name           string  `gorm:"column:name"`
		Type           string  `gorm:"column:type"`
		TotalCapacity  int     `gorm:"column:total_capacity"`
		AvailableSpots int     `gorm:"column:available_spots"`
		PricePerHour   float64 `gorm:"column:price_per_hour"`
		CreatedAt      string  `gorm:"column:created_at"`
		UpdatedAt      string  `gorm:"column:updated_at"`
	}

	var rows []row
	err := r.db.
		Table("parking_zones pz").
		Select(`pz.id, pz.name, pz.type, pz.total_capacity,
		        (pz.total_capacity - COALESCE((SELECT COUNT(*) FROM reservations r
		            WHERE r.zone_id = pz.id AND r.status = 'active'), 0)) AS available_spots,
		        pz.price_per_hour, pz.created_at, pz.updated_at`).
		Order("pz.id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]dto.ZoneResponse, 0, len(rows))
	for _, r := range rows {
		z := dto.ZoneResponse{
			ID:             r.ID,
			Name:           r.Name,
			Type:           r.Type,
			TotalCapacity:  r.TotalCapacity,
			AvailableSpots: r.AvailableSpots,
			PricePerHour:   r.PricePerHour,
		}
		if t, err := parseTime(r.CreatedAt); err == nil {
			z.CreatedAt = t
		}
		if t, err := parseTime(r.UpdatedAt); err == nil {
			z.UpdatedAt = t
		}
		out = append(out, z)
	}
	return out, nil
}

// FindByIDWithAvailability returns a single parking zone with its dynamically
// computed AvailableSpots. Returns gorm.ErrRecordNotFound when the zone
// does not exist.
func (r *zoneRepository) FindByIDWithAvailability(id uint) (*dto.ZoneResponse, error) {
	type row struct {
		ID             uint    `gorm:"column:id"`
		Name           string  `gorm:"column:name"`
		Type           string  `gorm:"column:type"`
		TotalCapacity  int     `gorm:"column:total_capacity"`
		AvailableSpots int     `gorm:"column:available_spots"`
		PricePerHour   float64 `gorm:"column:price_per_hour"`
		CreatedAt      string  `gorm:"column:created_at"`
		UpdatedAt      string  `gorm:"column:updated_at"`
	}

	var r0 row
	err := r.db.
		Table("parking_zones pz").
		Select(`pz.id, pz.name, pz.type, pz.total_capacity,
		        (pz.total_capacity - COALESCE((SELECT COUNT(*) FROM reservations r
		            WHERE r.zone_id = pz.id AND r.status = 'active'), 0)) AS available_spots,
		        pz.price_per_hour, pz.created_at, pz.updated_at`).
		Where("pz.id = ?", id).
		Scan(&r0).Error
	if err != nil {
		return nil, err
	}
	// GORM Scan does not return ErrRecordNotFound automatically. Detect
	// "no row" via zero-value ID combined with the requested id.
	if r0.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	z := &dto.ZoneResponse{
		ID:             r0.ID,
		Name:           r0.Name,
		Type:           r0.Type,
		TotalCapacity:  r0.TotalCapacity,
		AvailableSpots: r0.AvailableSpots,
		PricePerHour:   r0.PricePerHour,
	}
	if t, err := parseTime(r0.CreatedAt); err == nil {
		z.CreatedAt = t
	}
	if t, err := parseTime(r0.UpdatedAt); err == nil {
		z.UpdatedAt = t
	}
	return z, nil
}
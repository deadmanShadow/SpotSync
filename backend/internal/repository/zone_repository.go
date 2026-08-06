package repository

import (
	"errors"

	"gorm.io/gorm"

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
}

type zoneRepository struct {
	db *gorm.DB
}

// NewZoneRepository wires a ZoneRepository backed by the given GORM DB.
func NewZoneRepository(db *gorm.DB) ZoneRepository {
	return &zoneRepository{db: db}
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
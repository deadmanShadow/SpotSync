package repository

import (
	"errors"

	"gorm.io/gorm"

	"spotsync-backend/internal/models"
)

// ReservationRepository defines the persistence contract for Reservation records.
type ReservationRepository interface {
	Create(reservation *models.Reservation) error
	FindByID(id uint) (*models.Reservation, error)
	ListByUser(userID uint) ([]models.Reservation, error)
	ListAll() ([]models.Reservation, error)
	Update(reservation *models.Reservation) error
	Delete(id uint) error
}

type reservationRepository struct {
	db *gorm.DB
}

// NewReservationRepository wires a ReservationRepository backed by the given
// GORM DB.
func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(reservation *models.Reservation) error {
	return r.db.Create(reservation).Error
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := r.db.Preload("Zone").Preload("User").First(&reservation, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) ListByUser(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := r.db.Preload("Zone").
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepository) ListAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := r.db.Preload("Zone").Preload("User").
		Order("id DESC").
		Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepository) Update(reservation *models.Reservation) error {
	return r.db.Save(reservation).Error
}

func (r *reservationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Reservation{}, id).Error
}
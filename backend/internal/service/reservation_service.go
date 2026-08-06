package service

import (
	"errors"

	"gorm.io/gorm"

	"spotsync-backend/internal/dto"
	"spotsync-backend/internal/models"
	"spotsync-backend/internal/repository"
	"spotsync-backend/pkg/utils"
)

// Sentinel errors returned by ReservationService.
var (
	ErrReservationNotFound = errors.New("reservation not found")
	ErrZoneFull            = errors.New("zone has reached full capacity")
	ErrUnauthorized        = errors.New("not authorized to perform this action")
	ErrInvalidStatusChange = errors.New("invalid status transition")
)

// ReservationService defines the business-logic contract for reservations.
type ReservationService interface {
	Create(userID uint, req dto.CreateReservationRequest) (*models.Reservation, error)
	GetByID(userID uint, id uint, isAdmin bool) (*models.Reservation, error)
	ListByUser(userID uint) ([]models.Reservation, error)
	ListAll() ([]models.Reservation, error)
	UpdateStatus(userID uint, id uint, isAdmin bool, req dto.UpdateReservationStatusRequest) (*models.Reservation, error)
	Delete(userID uint, id uint, isAdmin bool) error
}

type reservationService struct {
	reservationRepo repository.ReservationRepository
	zoneRepo        repository.ZoneRepository
}

// NewReservationService wires a ReservationService backed by the given
// repositories.
func NewReservationService(reservationRepo repository.ReservationRepository, zoneRepo repository.ZoneRepository) ReservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		zoneRepo:        zoneRepo,
	}
}

func (s *reservationService) Create(userID uint, req dto.CreateReservationRequest) (*models.Reservation, error) {
	// Verify the zone exists and has capacity available.
	zone, err := s.zoneRepo.FindByID(req.ZoneID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}

	active, err := s.zoneRepo.CountActiveReservations(zone.ID)
	if err != nil {
		return nil, err
	}
	if int(active) >= zone.TotalCapacity {
		return nil, ErrZoneFull
	}

	reservation := &models.Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       "active",
	}
	if err := s.reservationRepo.Create(reservation); err != nil {
		return nil, err
	}
	return reservation, nil
}

func (s *reservationService) GetByID(userID uint, id uint, isAdmin bool) (*models.Reservation, error) {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}
	if !isAdmin && reservation.UserID != userID {
		return nil, ErrUnauthorized
	}
	return reservation, nil
}

func (s *reservationService) ListByUser(userID uint) ([]models.Reservation, error) {
	return s.reservationRepo.ListByUser(userID)
}

func (s *reservationService) ListAll() ([]models.Reservation, error) {
	return s.reservationRepo.ListAll()
}

func (s *reservationService) UpdateStatus(userID uint, id uint, isAdmin bool, req dto.UpdateReservationStatusRequest) (*models.Reservation, error) {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}
	if !isAdmin && reservation.UserID != userID {
		return nil, ErrUnauthorized
	}
	// Only allow transitioning from "active" to a terminal state.
	if reservation.Status != "active" {
		return nil, ErrInvalidStatusChange
	}
	reservation.Status = req.Status
	if err := s.reservationRepo.Update(reservation); err != nil {
		return nil, err
	}
	return reservation, nil
}

func (s *reservationService) Delete(userID uint, id uint, isAdmin bool) error {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReservationNotFound
		}
		return err
	}
	if !isAdmin && reservation.UserID != userID {
		return ErrUnauthorized
	}
	return s.reservationRepo.Delete(id)
}

// Ensure utils is referenced to avoid an unused import if helper functions are
// later moved into this file.
var _ = utils.HashPassword
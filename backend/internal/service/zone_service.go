package service

import (
	"errors"

	"gorm.io/gorm"

	"spotsync-backend/internal/dto"
	"spotsync-backend/internal/models"
	"spotsync-backend/internal/repository"
)

// ErrZoneNotFound is returned when a zone cannot be found.
var ErrZoneNotFound = errors.New("zone not found")

// ZoneService defines the business-logic contract for parking zones.
type ZoneService interface {
	Create(req dto.CreateZoneRequest) (*models.ParkingZone, error)
	GetByID(id uint) (*models.ParkingZone, error)
	List() ([]models.ParkingZone, error)
	Update(id uint, req dto.UpdateZoneRequest) (*models.ParkingZone, error)
	Delete(id uint) error
	ListWithAvailability() ([]dto.ZoneResponse, error)
	GetByIDWithAvailability(id uint) (*dto.ZoneResponse, error)
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
}

// NewZoneService wires a ZoneService backed by the given zone repository.
func NewZoneService(zoneRepo repository.ZoneRepository) ZoneService {
	return &zoneService{zoneRepo: zoneRepo}
}

func (s *zoneService) Create(req dto.CreateZoneRequest) (*models.ParkingZone, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}
	if err := s.zoneRepo.Create(zone); err != nil {
		return nil, err
	}
	return zone, nil
}

func (s *zoneService) GetByID(id uint) (*models.ParkingZone, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}
	return zone, nil
}

func (s *zoneService) List() ([]models.ParkingZone, error) {
	return s.zoneRepo.List()
}

func (s *zoneService) Update(id uint, req dto.UpdateZoneRequest) (*models.ParkingZone, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}
	if req.Name != "" {
		zone.Name = req.Name
	}
	if req.Type != "" {
		zone.Type = req.Type
	}
	if req.TotalCapacity > 0 {
		zone.TotalCapacity = req.TotalCapacity
	}
	if req.PricePerHour > 0 {
		zone.PricePerHour = req.PricePerHour
	}
	if err := s.zoneRepo.Update(zone); err != nil {
		return nil, err
	}
	return zone, nil
}

func (s *zoneService) Delete(id uint) error {
	if _, err := s.zoneRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrZoneNotFound
		}
		return err
	}
	return s.zoneRepo.Delete(id)
}

// ListWithAvailability returns every parking zone with its dynamically
// computed AvailableSpots value.
func (s *zoneService) ListWithAvailability() ([]dto.ZoneResponse, error) {
	return s.zoneRepo.FindAllWithAvailability()
}

// GetByIDWithAvailability returns a single parking zone with its dynamically
// computed AvailableSpots value. Returns ErrZoneNotFound if the zone does
// not exist.
func (s *zoneService) GetByIDWithAvailability(id uint) (*dto.ZoneResponse, error) {
	z, err := s.zoneRepo.FindByIDWithAvailability(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}
	return z, nil
}
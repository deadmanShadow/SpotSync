// Package seeder contains the catalog seeder that ensures the
// parking_zones table is populated with 40 diverse zones across the
// three supported types (EV Charging, General, Covered).
//
// The seeder is intentionally idempotent and safe to run on every
// server start: it only inserts rows when the catalog is empty or
// has fewer than 40 zones. Existing zones are never overwritten, so
// admins can keep customizing prices and capacities without losing
// their changes.
//
// Companion to this package is the `ZoneRotator` (in rotator.go)
// which keeps the 60% available / 40% full invariant maintained on
// a 1-hour rotation cadence.
package seeder

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"spotsync-backend/internal/models"
)

// TargetZoneCount is the canonical number of zones the catalog should
// have after seeding. Used by both the seeder (to know when to top up)
// and the rotator (to compute the 60/40 split).
const TargetZoneCount = 40

// SeedZoneSpec is the in-memory specification of a single seeded
// zone. We keep these as named values rather than raw GORM models so
// the catalog is easy to read and edit in one place.
type SeedZoneSpec struct {
	Name          string
	Type          string
	TotalCapacity int
	PricePerHour  float64
}

// FullCatalog is the canonical 40-zone source of truth. Counts:
//
//	EV Charging: 15 zones
//	General:     15 zones
//	Covered:     10 zones
//
// Capacity ranges from 12 to 60 and prices from $2.00 to $9.00 per hour
// to look realistic across the demo. Names are intentionally varied
// (airports, malls, downtown, hospitals, stadiums) so the catalog
// filter and search have meaningful test data.
var FullCatalog = []SeedZoneSpec{
	// ---------- EV Charging (15 zones) ----------
	{Name: "Terminal 1 EV Charging Hub", Type: "ev_charging", TotalCapacity: 24, PricePerHour: 6.50},
	{Name: "Airport North EV Bay", Type: "ev_charging", TotalCapacity: 18, PricePerHour: 7.00},
	{Name: "Mall Central EV Charging", Type: "ev_charging", TotalCapacity: 32, PricePerHour: 5.50},
	{Name: "Convention Center EV Stalls", Type: "ev_charging", TotalCapacity: 20, PricePerHour: 6.00},
	{Name: "Stadium West EV Parking", Type: "ev_charging", TotalCapacity: 40, PricePerHour: 8.00},
	{Name: "Tech Park EV Plaza", Type: "ev_charging", TotalCapacity: 28, PricePerHour: 5.00},
	{Name: "Riverside EV Charge Point", Type: "ev_charging", TotalCapacity: 16, PricePerHour: 5.50},
	{Name: "City Hall EV Priority", Type: "ev_charging", TotalCapacity: 12, PricePerHour: 4.50},
	{Name: "Hospital EV Emergency Bay", Type: "ev_charging", TotalCapacity: 14, PricePerHour: 9.00},
	{Name: "University EV Fast Charge", Type: "ev_charging", TotalCapacity: 22, PricePerHour: 4.00},
	{Name: "Highway Rest EV Station", Type: "ev_charging", TotalCapacity: 30, PricePerHour: 7.50},
	{Name: "Hotel District EV Valet", Type: "ev_charging", TotalCapacity: 18, PricePerHour: 8.50},
	{Name: "Sports Arena EV Lot", Type: "ev_charging", TotalCapacity: 26, PricePerHour: 6.00},
	{Name: "Outlet Mall EV Corner", Type: "ev_charging", TotalCapacity: 20, PricePerHour: 5.00},
	{Name: "Beachfront EV Charging", Type: "ev_charging", TotalCapacity: 16, PricePerHour: 6.50},

	// ---------- General (15 zones) ----------
	{Name: "Terminal 1 General Parking", Type: "general", TotalCapacity: 60, PricePerHour: 4.50},
	{Name: "Terminal 2 Open Lot", Type: "general", TotalCapacity: 48, PricePerHour: 4.00},
	{Name: "Downtown Street Parking", Type: "general", TotalCapacity: 36, PricePerHour: 5.50},
	{Name: "Market Square Lot", Type: "general", TotalCapacity: 28, PricePerHour: 3.50},
	{Name: "Train Station North Lot", Type: "general", TotalCapacity: 45, PricePerHour: 4.00},
	{Name: "Bus Depot Long Stay", Type: "general", TotalCapacity: 52, PricePerHour: 3.00},
	{Name: "Festival Grounds Parking", Type: "general", TotalCapacity: 60, PricePerHour: 5.00},
	{Name: "Riverside Open Lot", Type: "general", TotalCapacity: 40, PricePerHour: 3.50},
	{Name: "Eastside Community Lot", Type: "general", TotalCapacity: 32, PricePerHour: 2.50},
	{Name: "Westgate Plaza Parking", Type: "general", TotalCapacity: 38, PricePerHour: 4.50},
	{Name: "Northgate Shopping Lot", Type: "general", TotalCapacity: 44, PricePerHour: 4.00},
	{Name: "South Pier Parking", Type: "general", TotalCapacity: 50, PricePerHour: 5.50},
	{Name: "Civic Center Open Lot", Type: "general", TotalCapacity: 34, PricePerHour: 3.50},
	{Name: "Park & Walk General", Type: "general", TotalCapacity: 28, PricePerHour: 2.00},
	{Name: "Industrial District Lot", Type: "general", TotalCapacity: 42, PricePerHour: 3.00},

	// ---------- Covered (10 zones) ----------
	{Name: "Garage Level 1 — Downtown", Type: "covered", TotalCapacity: 50, PricePerHour: 6.00},
	{Name: "Garage Level 2 — Downtown", Type: "covered", TotalCapacity: 50, PricePerHour: 6.00},
	{Name: "Underground P-1 — Mall", Type: "covered", TotalCapacity: 36, PricePerHour: 5.50},
	{Name: "Underground P-2 — Mall", Type: "covered", TotalCapacity: 36, PricePerHour: 5.50},
	{Name: "Covered Garage Central", Type: "covered", TotalCapacity: 44, PricePerHour: 5.00},
	{Name: "Airport Covered Deck", Type: "covered", TotalCapacity: 60, PricePerHour: 7.50},
	{Name: "Hospital Covered Garage", Type: "covered", TotalCapacity: 32, PricePerHour: 6.50},
	{Name: "Stadium North Covered", Type: "covered", TotalCapacity: 40, PricePerHour: 7.00},
	{Name: "Convention Center Garage", Type: "covered", TotalCapacity: 28, PricePerHour: 6.50},
	{Name: "Beachside Covered Deck", Type: "covered", TotalCapacity: 22, PricePerHour: 7.50},
}

// SeedIfNeeded inspects the current parking_zones table and tops it up
// to TargetZoneCount zones when it is empty (or has fewer rows than
// the catalog). Existing zones are left untouched so admin edits to
// prices/capacities are preserved.
//
// The function is safe to call concurrently with the HTTP server
// (it only INSERTs and only when the count is below target).
func SeedIfNeeded(db *gorm.DB) (int, error) {
	var existingCount int64
	if err := db.Model(&models.ParkingZone{}).Count(&existingCount).Error; err != nil {
		return 0, fmt.Errorf("seeder: failed to count existing zones: %w", err)
	}

	// Already at or above target — nothing to do.
	if int(existingCount) >= TargetZoneCount {
		log.Printf("[seeder] Catalog already has %d zones (target: %d). Skipping seed.",
			existingCount, TargetZoneCount)
		return 0, nil
	}

	// Build the list of zones we still need to insert. We pull existing
	// names so we never duplicate a row, even across restarts.
	existingNames, err := loadExistingZoneNames(db)
	if err != nil {
		return 0, err
	}

	var toInsert []models.ParkingZone
	for _, spec := range FullCatalog {
		if _, existing := existingNames[spec.Name]; existing {
			continue
		}
		toInsert = append(toInsert, models.ParkingZone{
			Name:          spec.Name,
			Type:          spec.Type,
			TotalCapacity: spec.TotalCapacity,
			PricePerHour:  spec.PricePerHour,
		})
	}

	if len(toInsert) == 0 {
		log.Printf("[seeder] No new zones to insert (existing: %d).", existingCount)
		return 0, nil
	}

	if len(toInsert) == 0 {
		log.Printf("[seeder] No new zones to insert (existing: %d).", existingCount)
		return 0, nil
	}

	if err := db.Create(&toInsert).Error; err != nil {
		return 0, fmt.Errorf("seeder: failed to insert zones: %w", err)
	}

	// Per spec §3.5: regenerate spot_holds for each newly seeded zone
	// so a fresh DB boots into realistic per-spot availability instead
	// of every spot appearing available.
	for i := range toInsert {
		if err := RegenerateSpotHolds(db, toInsert[i].ID); err != nil {
			log.Printf("[seeder] RegenerateSpotHolds(zone %d) failed: %v", toInsert[i].ID, err)
		}
	}

	log.Printf("[seeder] Inserted %d new zones. Catalog now has %d zones (target: %d).",
		len(toInsert), int(existingCount)+len(toInsert), TargetZoneCount)
	return len(toInsert), nil
}

// loadExistingZoneNames returns a set of zone names already present in
// the database so we can avoid duplicate inserts. Used internally by
// SeedIfNeeded.
func loadExistingZoneNames(db *gorm.DB) (map[string]struct{}, error) {
	var rows []struct {
		Name string `gorm:"column:name"`
	}
	if err := db.Model(&models.ParkingZone{}).
		Select("name").
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("seeder: failed to load zone names: %w", err)
	}
	set := make(map[string]struct{}, len(rows))
	for _, r := range rows {
		set[r.Name] = struct{}{}
	}
	return set, nil
}

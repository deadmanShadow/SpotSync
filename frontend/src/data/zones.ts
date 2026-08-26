/**
 * Mock dataset of 40 parking zones used for the frontend demo.
 *
 * The dataset is split across the three supported zone types and is
 * long-lived enough to exercise the catalog's filter / search / scroll
 * behavior end-to-end without touching the backend.
 *
 * Naming convention:
 *   - EV Charging zones: "Terminal N EV Charging Hub" / "Mall EV Bay"
 *   - General zones:     "Terminal N Parking" / "Street Parking"
 *   - Covered zones:     "Garage Level N" / "Underground P-N"
 *
 * Prices are realistic USD-per-hour values to make the catalog feel
 * grounded. Capacities range from 12 to 60 spots per zone.
 *
 * Each zone's `available_spots` is mutated at runtime by the
 * `availabilityRotator` in `src/services/mockZoneData.ts` so the
 * catalog stays fresh on a 1-hour rotation cadence.
 */

import type { ParkingZone, ZoneType } from "../types/zone";

/**
 * Build a stable Unix timestamp for the seed moment. Using `Date.now()`
 * at module load time is fine — the seed itself doesn't need to be
 * deterministic across builds, only the zone IDs and capacities.
 */
const SEED_TIMESTAMP = new Date("2026-01-15T08:00:00Z").toISOString();

interface ZoneSpec {
  name: string;
  type: ZoneType;
  total_capacity: number;
  price_per_hour: number;
}

/**
 * The 40 source-of-truth zone specifications. The numeric `id` and
 * initial `available_spots` are assigned below by `buildSeedZones()`
 * so we don't have to hand-author them.
 */
const ZONE_SPECS: ZoneSpec[] = [
  // ---------- EV Charging (15 zones) ----------
  { name: "Terminal 1 EV Charging Hub",    type: "ev_charging", total_capacity: 24, price_per_hour: 6.50 },
  { name: "Airport North EV Bay",          type: "ev_charging", total_capacity: 18, price_per_hour: 7.00 },
  { name: "Mall Central EV Charging",      type: "ev_charging", total_capacity: 32, price_per_hour: 5.50 },
  { name: "Convention Center EV Stalls",   type: "ev_charging", total_capacity: 20, price_per_hour: 6.00 },
  { name: "Stadium West EV Parking",       type: "ev_charging", total_capacity: 40, price_per_hour: 8.00 },
  { name: "Tech Park EV Plaza",            type: "ev_charging", total_capacity: 28, price_per_hour: 5.00 },
  { name: "Riverside EV Charge Point",     type: "ev_charging", total_capacity: 16, price_per_hour: 5.50 },
  { name: "City Hall EV Priority",         type: "ev_charging", total_capacity: 12, price_per_hour: 4.50 },
  { name: "Hospital EV Emergency Bay",     type: "ev_charging", total_capacity: 14, price_per_hour: 9.00 },
  { name: "University EV Fast Charge",     type: "ev_charging", total_capacity: 22, price_per_hour: 4.00 },
  { name: "Highway Rest EV Station",       type: "ev_charging", total_capacity: 30, price_per_hour: 7.50 },
  { name: "Hotel District EV Valet",       type: "ev_charging", total_capacity: 18, price_per_hour: 8.50 },
  { name: "Sports Arena EV Lot",           type: "ev_charging", total_capacity: 26, price_per_hour: 6.00 },
  { name: "Outlet Mall EV Corner",         type: "ev_charging", total_capacity: 20, price_per_hour: 5.00 },
  { name: "Beachfront EV Charging",        type: "ev_charging", total_capacity: 16, price_per_hour: 6.50 },

  // ---------- General (15 zones) ----------
  { name: "Terminal 1 General Parking",    type: "general",     total_capacity: 60, price_per_hour: 4.50 },
  { name: "Terminal 2 Open Lot",           type: "general",     total_capacity: 48, price_per_hour: 4.00 },
  { name: "Downtown Street Parking",      type: "general",     total_capacity: 36, price_per_hour: 5.50 },
  { name: "Market Square Lot",             type: "general",     total_capacity: 28, price_per_hour: 3.50 },
  { name: "Train Station North Lot",      type: "general",     total_capacity: 45, price_per_hour: 4.00 },
  { name: "Bus Depot Long Stay",           type: "general",     total_capacity: 52, price_per_hour: 3.00 },
  { name: "Festival Grounds Parking",     type: "general",     total_capacity: 60, price_per_hour: 5.00 },
  { name: "Riverside Open Lot",            type: "general",     total_capacity: 40, price_per_hour: 3.50 },
  { name: "Eastside Community Lot",        type: "general",     total_capacity: 32, price_per_hour: 2.50 },
  { name: "Westgate Plaza Parking",        type: "general",     total_capacity: 38, price_per_hour: 4.50 },
  { name: "Northgate Shopping Lot",        type: "general",     total_capacity: 44, price_per_hour: 4.00 },
  { name: "South Pier Parking",            type: "general",     total_capacity: 50, price_per_hour: 5.50 },
  { name: "Civic Center Open Lot",         type: "general",     total_capacity: 34, price_per_hour: 3.50 },
  { name: "Park & Walk General",           type: "general",     total_capacity: 28, price_per_hour: 2.00 },
  { name: "Industrial District Lot",       type: "general",     total_capacity: 42, price_per_hour: 3.00 },

  // ---------- Covered (10 zones) ----------
  { name: "Garage Level 1 — Downtown",     type: "covered",     total_capacity: 50, price_per_hour: 6.00 },
  { name: "Garage Level 2 — Downtown",     type: "covered",     total_capacity: 50, price_per_hour: 6.00 },
  { name: "Underground P-1 — Mall",        type: "covered",     total_capacity: 36, price_per_hour: 5.50 },
  { name: "Underground P-2 — Mall",        type: "covered",     total_capacity: 36, price_per_hour: 5.50 },
  { name: "Covered Garage Central",        type: "covered",     total_capacity: 44, price_per_hour: 5.00 },
  { name: "Airport Covered Deck",          type: "covered",     total_capacity: 60, price_per_hour: 7.50 },
  { name: "Hospital Covered Garage",       type: "covered",     total_capacity: 32, price_per_hour: 6.50 },
  { name: "Stadium North Covered",         type: "covered",     total_capacity: 40, price_per_hour: 7.00 },
  { name: "Convention Center Garage",      type: "covered",     total_capacity: 28, price_per_hour: 6.50 },
  { name: "Beachside Covered Deck",        type: "covered",     total_capacity: 22, price_per_hour: 7.50 },
];

/**
 * Build the seed zones with IDs and an initial "60% available" config:
 *   - 24 zones (60%) start with available_spots between 50% and 100% of capacity
 *   - 16 zones (40%) start full (available_spots = 0)
 *
 * The actual rotation behaviour is owned by `availabilityRotator`; this
 * seed is just the first paint.
 */
export function buildSeedZones(): ParkingZone[] {
  const total = ZONE_SPECS.length;
  const availableCount = Math.round(total * 0.6); // 24 zones available

  return ZONE_SPECS.map((spec, index) => {
    const isAvailable = index < availableCount;
    const available_spots = isAvailable
      ? Math.max(1, Math.round(spec.total_capacity * (0.5 + Math.random() * 0.5)))
      : 0;
    // Build a placeholder spot_holds array that matches total_capacity
    // so the type satisfies ParkingZone. The runtime layer
    // (mockZoneData.ts) regenerates spot_holds on first read so the
    // first paint shows realistic per-spot availability rather than
    // whatever we put here.
    const spot_holds: number[] = new Array(spec.total_capacity).fill(0);
    const zone: ParkingZone = {
      id: index + 1,
      name: spec.name,
      type: spec.type,
      total_capacity: spec.total_capacity,
      available_spots,
      price_per_hour: spec.price_per_hour,
      spot_holds,
      created_at: SEED_TIMESTAMP,
    };
    return zone;
  });
}

/**
 * The complete seed dataset. Exported as a frozen array so callers
 * can't accidentally mutate it — the live, mutable copy lives in
 * `mockZoneData.ts`.
 */
export const SEED_ZONES: readonly ParkingZone[] = Object.freeze(buildSeedZones());

/**
 * Number of zones in the seed dataset. Exported so the rotator can
 * compute the 60% / 40% split without hard-coding 40.
 */
export const TOTAL_ZONE_COUNT: number = SEED_ZONES.length;

/**
 * Domain type contracts for parking zones.
 *
 * A ParkingZone represents a logical grouping of parking spots (e.g. an EV
 * charging bay, a covered garage floor, or an open-air general lot). The
 * `available_spots` field is the live, server-managed counter that powers
 * the dynamic capacity progress bars on the home page.
 */

export type ZoneType = "general" | "ev_charging" | "covered";

/** Parking zone as returned by GET /zones and friends. */
export interface ParkingZone {
  id: number;
  name: string;
  type: ZoneType;
  total_capacity: number;
  available_spots: number;
  price_per_hour: number;
  created_at: string;
}

/** Payload used by POST /zones (admin-only endpoint). */
export interface CreateZonePayload {
  name: string;
  type: ZoneType;
  total_capacity: number;
  price_per_hour: number;
}

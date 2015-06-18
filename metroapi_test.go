package metroapi

import "testing"

func TestGetTimepointDepartures(t *testing.T) {

  departures := GetTimepointDepartures(902, 3, "STVI")
  if len(departures) == 0 {
    t.Errorf("GetTimepointDepartures returned 0 items")
  } else {
    if departures[0].TimeOfDeparture == nil {
      t.Errorf("GeneralDepartures could not parse time")
    }
  }
}

func TestGetProviders(t *testing.T) {

  if len(GetProviders()) == 0 {
    t.Errorf("GetProviders returned 0 items")
  }
}

func TestGetRoutes(t *testing.T) {

  if len(GetRoutes()) == 0 {
    t.Errorf("GetRoutes returned 0 items")
  }
}

func TestGetDirections(t *testing.T) {

  if len(GetDirections(5)) == 0 {
    t.Errorf("GetDirections returned 0 items")
  }
}

func TestGetStops(t *testing.T) {

  if len(GetStops(5, 4)) == 0 {
    t.Errorf("GetStops returned 0 items")
  }
}

func TestGetDepartures(t *testing.T) {

  if len(GetDepartures(11167)) == 0 {
    t.Errorf("GetDepartures returned 0 items")
  }
}

func TestGetVehicleLocations(t *testing.T) {

  locations := GetVehicleLocations(5)
  if len(locations) == 0 {
    t.Errorf("GetVehicleLocations returned 0 items")
  } else {
    if locations[0].LastCheckinTime == nil {
      t.Errorf("VehicleLocations could not parse date")
    }
  }
}

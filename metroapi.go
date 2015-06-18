package metroapi

import (
  "encoding/json"
  "net/http"
  "io/ioutil"
  "strconv"
	"time"
	"regexp"
	"errors"
)

// Metro Transit API data structure types
type GeneralResponse struct {
  Text              string
  Value             string
}

type GeneralDepartures struct {
  Actual            bool
  BlockNumber       int
  DepartureText     string
  Description       string
  Gate              string
  Route             string
  RouteDirection    string
  Terminal          string
  VehicleHeading    int
  VehicleLatitude   float32
  VehicleLongitude  float32
	TimeOfDeparture   *time.Time
}

type Departures struct {
  GeneralDepartures
}

type Directions struct {
  GeneralResponse
}

type Providers struct {
  GeneralResponse
}

type Routes struct {
  Description       string
  ProviderID        string
  Route             string
}

type Stops struct {
  GeneralResponse
}

type TimepointDepartures struct {
  GeneralDepartures
}

type VehicleLocations struct {
  BlockNumber         int
  Direction           int
  Route               string
  Terminal            string
  VehicleLatitude     float32
  VehicleLongitude    float32
  Bearing             int
  Odometer            int
  Speed               int
	LastCheckinTime     *time.Time
}

// Retrieve JSON data from the URL provided and return a byte array of the result.
func retrieveData(url string) []byte {

  // retrieve the JSON from the provided URL
  res, _ := http.Get(url)
  defer res.Body.Close()
  body, _ := ioutil.ReadAll(res.Body)
  return body
}

// convert a JavaScript date representation to a golang Time
func convertJsDate(jsDate string) *time.Time {
	pattern := regexp.MustCompile(`Date\(([0-9]{13})-[0-9]{4}`)
	if matched := pattern.FindStringSubmatch(jsDate); len(matched) > 0 {
		if unixTimestamp, err := strconv.Atoi(matched[1][:10]); err == nil {
			converted := time.Unix(int64(unixTimestamp), 0)
			return &converted
		}
	}
	panic(errors.New("Invalid JavaScript date format string"))
}

// GeneralDepartures @implements json.Unmarshaler
func (g *GeneralDepartures) UnmarshalJSON(b []byte) error {
	
	// Alias for GeneralDepartures that prevents infinite recursion in UnmarshallJSON
	type generalDeparturesAlias GeneralDepartures
	
  var err error = nil
	newGeneralDepartures := generalDeparturesAlias{}
	jsonMap := make(map[string]*json.RawMessage)

	err = json.Unmarshal(b, &newGeneralDepartures)
	err = json.Unmarshal(b, &jsonMap)

	datestr := jsonMap["DepartureTime"]
	newGeneralDepartures.TimeOfDeparture = convertJsDate(string(*datestr))

	*g = GeneralDepartures(newGeneralDepartures)
	return err
}


// VehicleLocations @implements json.Unmarshaler
func (v *VehicleLocations) UnmarshalJSON(b []byte) error {
	
	// Alias for GeneralDepartures that prevents infinite recursion in UnmarshallJSON
	type vehicleLocationsAlias VehicleLocations

	var err error = nil
	newVehicleLocations := vehicleLocationsAlias{}
	jsonMap := make(map[string]*json.RawMessage)

	err = json.Unmarshal(b, &newVehicleLocations)
	err = json.Unmarshal(b, &jsonMap)

	datestr := jsonMap["LocationTime"]
	newVehicleLocations.LastCheckinTime = convertJsDate(string(*datestr))

	*v = VehicleLocations(newVehicleLocations)
	return err
}

// Returns the scheduled departures for a selected route, direction and timepoint stop.
func GetTimepointDepartures(route, direction int, stop string) []TimepointDepartures {

  url := "http://svc.metrotransit.org/NexTrip/" + strconv.Itoa(route) + "/" + strconv.Itoa(direction) + "/" + stop + "?format=json"
  body := retrieveData(url)
  var data []TimepointDepartures
  json.Unmarshal(body, &data)
  return data
}

// Returns a list of area Transit providers. Providers are identified in
// the list of Routes allowing routes to be selected for a single provider.
func GetProviders() []Providers {

  url := "http://svc.metrotransit.org/nextrip/providers?format=json"
  body := retrieveData(url)
  var data []Providers
  json.Unmarshal(body, &data)
  return data
}

// Returns a list of Transit routes that are in service on the current day.
func GetRoutes() []Routes {

  url := "http://svc.metrotransit.org/NexTrip/Routes?format=json"
  body := retrieveData(url)
  var data []Routes
  json.Unmarshal(body, &data)
  return data
}

// Returns the two directions that are valid for a given route. Either North/South or East/West.
// The result includes text/value pair with the direction name and an ID. Directions are
// identified with an ID value. 1 = South, 2 = East, 3 = West, 4 = North.
func GetDirections(route int) []Directions {

  url := "http://svc.metrotransit.org/NexTrip/Directions/" + strconv.Itoa(route) + "?format=json"
  body := retrieveData(url)
  var data []Directions
  json.Unmarshal(body, &data)
  return data
}

// Returns a list of Timepoint stops for the given Route/Direction. The result includes
// text/value pairs with the stop description and a 4 character stop (or node) identifier.
func GetStops(route, direction int) []Stops {

  url := "http://svc.metrotransit.org/NexTrip/Stops/" + strconv.Itoa(route) + "/" + strconv.Itoa(direction) + "?format=json"
  body := retrieveData(url)
  var data []Stops
  json.Unmarshal(body, &data)
  return data
}

// This operation is used to return a list of departures scheduled for any given bus stop.
// A StopID is an integer value identifying any one of the many thousands of bus stops in the metro.
// Stop information can be derived from the GTFS schedule data updated weekly for public use.
// http://datafinder.org/metadata/transit_schedule_google_feed.html
func GetDepartures(stop int) []Departures {

  url := "http://svc.metrotransit.org/NexTrip/" + strconv.Itoa(stop) + "?format=json"
  body := retrieveData(url)
  var data []Departures
  json.Unmarshal(body, &data)
  return data
}

// This operation returns a list of vehicles currently in service that have recently
// (within 5 minutes) reported their locations. A route paramter is used to return
// results for the given route. Use "0" for the route parameter to return a list of
// all vehicles in service.
func GetVehicleLocations(route int) []VehicleLocations {

  url := "http://svc.metrotransit.org/NexTrip/VehicleLocations/" + strconv.Itoa(route) + "?format=json"
  body := retrieveData(url)
  var data []VehicleLocations
  json.Unmarshal(body, &data)
  return data
}

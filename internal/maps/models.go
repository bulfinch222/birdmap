package maps

type Observation struct {
	SpeciesCode string  `json:"speciesCode"`
	ComName 	string	`json:"comName"`
	SciName 	string	`json:"sciName"`
	LocId 		string	`json:"locId"`
	LocName 	string	`json:"locName"`
	ObsDt 		string	`json:"obsDt"`
	Lat 		float64	`json:"lat"`
	Lng 		float64	`json:"lng"`
}

type InatResponce struct {
	Results []InatObservation `json:"results"`
}

type Geojson struct{
	  Coordinates []float64 `json:"coordinates"`
}

type InatObservation struct {
	SpeciesGuess string          `json:"species_guess"`
	QualityGrade string          `json:"quality_grade"`
	TimeObservedAt string        `json:"time_observed_at"` 
	Taxon        Taxon           `json:"taxon"`
	Geojson      Geojson    	 `json:"geojson"`
}

type Taxon struct {
	Name string `json:"name"`
}

type Location struct {
	LocId 		 string				`json:"locId"`
	LocName 	 string				`json:"locName"`
    Lat 		 float64			`json:"lat"`
	Lng 		 float64 			`json:"lng"`
	Observations []Observation 		`json:"observations"`
}

type Ip struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}


var UserLocation *Ip
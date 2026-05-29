package maps

import (
	"birdmap/internal"
	"birdmap/internal/lifelist"
	"birdmap/internal/taxonomy"
	"log"
	"time"
	"sync"
)

	func GetEbirdLocations(key string, s *taxonomy.Service) ([]Location, error){
		locationMap := make(map[string]*Location)
		if UserLocation == nil {
			GetCurrentLocation()
		}
		log.Println(1)
		data, err := getEBirdObservations(key)
		if (err != nil) {
			return nil, err
		}
		log.Println(2)
		var wg sync.WaitGroup
		for ind, observation := range data {
			exist, err := lifelist.CheckIfInList(s, internal.Bird{
				ComName: observation.ComName,
				SpeciesCode: observation.SpeciesCode,
				SciName: observation.SciName,
			})
			if err != nil {
				return nil, err
			}
			if exist {
				continue
			}
			if ind >= 25 {
				break
			}
			wg.Go(func() {GetEbirdLocationsOfSpecies(key, locationMap, observation.SpeciesCode)})
		}
		wg.Wait()
		log.Println(2)
		locations := make([]Location, 0, len(locationMap))

		for _, location := range locationMap {
			locations = append(locations, *location)
		}
		
		return locations, nil

	}
	func GetEbirdLocationsOfSpecies(key string, locationMap map[string]*Location, speciesCode string) error{
		if UserLocation == nil {
			GetCurrentLocation()
		}
		log.Println(key)
		data, err := getEbirdObservationsOfSpecies(key, speciesCode)
		if err != nil {
			return err
		}

		for _, observation := range data {
			if loc, ok := locationMap[observation.LocId]; ok {
				loc.Observations = append(loc.Observations, observation)
			} else {
				locationMap[observation.LocId] = &Location{
					LocId: observation.LocId,
					LocName: observation.LocName,
					Lat: observation.Lat,
					Lng: observation.Lng,
					Observations: []Observation{observation},
				}
			}
		}
		return nil
	}

	func ConvertInatToEbird(service *taxonomy.Service) (locations []Location, _ error){
		observations, err := getInatObservations()

		if err != nil {
			return nil, err
		}
		log.Println(5)
		for _, obs := range observations {
			bird, err := service.SearchBirds(obs.Taxon.Name)

			if err != nil {
				log.Println(err)
				continue
			}
			if len(bird) == 0{
				log.Printf("Проигнорировано наблюдение из Inat: %s\n", obs.Taxon.Name)
				continue
			}

			t, err := time.Parse(time.RFC3339, obs.TimeObservedAt)
			
			if err != nil {
				log.Printf("Ошибка при парсинге времени: %v", err)
				continue
			}

			output := t.Format("2006-01-02 15:04")

			b := bird[0]
			newObs := Observation{
				SpeciesCode: b.SpeciesCode,
				ComName: b.ComName,
				SciName: b.SciName,
				ObsDt: output,
				Lat: obs.Geojson.Coordinates[1],
				Lng: obs.Geojson.Coordinates[0],
			}

			loc := &Location{
				Lat: newObs.Lat,
				Lng: newObs.Lng,
				Observations: []Observation{newObs},
			}

			locations = append(locations, *loc)
		}
		return locations, nil
	}

		func GetLocations(key string, service *taxonomy.Service) ([]Location, error) {
			locations, err := GetEbirdLocations(key, service)
			log.Println(8)
			//log.Println(locations)
			if err != nil {
				return nil, err
			}

			inatLocs, err := ConvertInatToEbird(service)

			if err != nil {
				return nil, err
			}
			locations = append(locations, inatLocs...)
			//log.Println(locations)
			return locations, nil
		}

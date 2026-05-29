package maps

import (
	"net/http"
	"net/url"
	"io"
	"encoding/json"
	"fmt"
	"time"
	"log"
)

func GetCurrentLocation() error {
	resp, err := http.Get("https://ipwhois.app/json/")
	log.Println(6)
	if (err != nil) {
		return fmt.Errorf("Ошибка получения текущей локации: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if (err != nil) {
		return err
	}
	err = json.Unmarshal(body, &UserLocation)
	if (err != nil) {
		return err
	}
	log.Printf("Текущая локация: %+v", UserLocation)
	log.Println(10)
	return nil
}

func getEBirdObservations(key string) ([]Observation, error){
	url := fmt.Sprintf("https://api.ebird.org/v2/data/obs/geo/recent?lat=%v&lng=%v&dist=50&key=%v", UserLocation.Latitude, UserLocation.Longitude, key)
	resp, err := http.Get(url)
	if (err != nil) {
		return nil, fmt.Errorf("Ошибка при запросе к наблюдениям EBird: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if (err != nil){
		return nil, err
	}
	var data []Observation 
	err = json.Unmarshal(body, &data)
	if (err != nil) {
		return nil, err
	}
	return data, nil

}

func getEbirdObservationsOfSpecies(key string,speciesCode string) ([]Observation, error){
	url := fmt.Sprintf("https://api.ebird.org/v2/data/obs/geo/recent/%v?lat=%v&lng=%v&dist=50&key=%v", speciesCode, UserLocation.Latitude, UserLocation.Longitude, key)
	resp, err := http.Get(url)
	if (err != nil) {
		return nil, fmt.Errorf("Ошибка при запросе к наблюдениям EBird конкретных видов: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if (err != nil) {
		return nil, fmt.Errorf("2%w", err)
	}
	var data []Observation
	err = json.Unmarshal(body, &data)
	if (err != nil) {
		return nil, fmt.Errorf("1%w", err)
	}
	return data, nil
}

func getInatObservations() ([]InatObservation, error){
	now := time.Now()
	back := time.Now().AddDate(0, 0, -14)

	d2:=now.Format("2006-01-02")
	d1:=back.Format("2006-01-02")
	
	lat := fmt.Sprintf("%f", UserLocation.Latitude)
	lng := fmt.Sprintf("%f", UserLocation.Longitude)


	params := url.Values{}
	params.Add("radius", "50")
	params.Add("d1", d1)
	params.Add("d2", d2)
	params.Add("lat", lat)
	params.Add("lng", lng)
	params.Add("geoprivacy", "open")
	params.Add("quality_grade", "research")
	params.Add("iconic_taxa[]", "Aves")

	p := params.Encode()
	query := fmt.Sprintf("https://api.inaturalist.org/v1/observations?%v", p)

	resp, err := http.Get(query)

	if err != nil {
		return nil, fmt.Errorf("Ошибка при запросе в INat")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if (err != nil){
		return nil, err
	}

	var data InatResponce
	err = json.Unmarshal(body, &data)
	if (err != nil) {
		return nil, err
	}
	//log.Println(data)
	return data.Results, nil
}
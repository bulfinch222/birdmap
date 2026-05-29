package lifelist

import (
	"birdmap/internal"
	"birdmap/internal/taxonomy"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func CreateLifelist(s *taxonomy.Service) error {
	query := `CREATE TABLE IF NOT EXISTS lifelist(
		id INTEGER PRIMARY KEY,
		speciescode TEXT NOT NULL UNIQUE,
		scientificname TEXT NOT NULL UNIQUE,	
		commonname TEXT NOT NULL

	);`

	_, err := s.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("Не удалось создать базу данных лайфлиста: %w", err)
	}
	return nil
}	

func AddBird(bird internal.Bird, s *taxonomy.Service) error {
	query := `INSERT OR IGNORE INTO lifelist (speciescode, scientificname, commonname)
	VALUES (?, ?, ?)`;
	_, err := s.Db.Exec(query, bird.SpeciesCode, bird.SciName, bird.ComName)
	if err != nil {
		return fmt.Errorf("Ошибка при добавлении птицы в лайлист: %w", err)
	}
	return nil;
}

func GetLifelist(s *taxonomy.Service) (birds []internal.Bird, _ error) {
	query := `SELECT speciescode, scientificname, commonname FROM lifelist`	

	rows, err := s.Db.Query(query)
	
	if err != nil {
		return nil, err
	}
	
	for rows.Next() {
		bird := internal.Bird{}
		err := rows.Scan(&bird.SpeciesCode, &bird.SciName, &bird.ComName)

		if err != nil {
			fmt.Print(err)
			continue
		}

		birds= append(birds, bird)
	}
	return birds, nil
}

func ImportLifelist(s *taxonomy.Service, filepath string) error{
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("Ошибка при открытии файла csv: %w", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	
	if _, err := reader.Read(); err != nil {
		return err
	}
	deleteLifelist(s)
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Ошибка чтения файла csv: %w", err)
		}
		bird := internal.Bird{
			SpeciesCode: record[0],
			ComName: record[3],
			SciName: record[4],
		}
		AddBird(bird, s)
	}
	return nil
}

func CheckIfInList(s *taxonomy.Service, bird internal.Bird)(bool, error){
	query := `SELECT COUNT(*) FROM lifelist WHERE commonname = ?`
	var count int 

	err := s.Db.QueryRow(query, bird.ComName).Scan(&count)
	log.Printf("Количество птиц %v+: %v", bird.ComName, count)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func deleteLifelist(s *taxonomy.Service){
	_, err := s.Db.Exec("DELETE FROM lifelist")
	if err != nil {
		log.Fatal(err)
	}
}
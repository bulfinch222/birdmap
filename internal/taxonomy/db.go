package taxonomy

import (
	"birdmap/internal"
	"context"
	"fmt"
)

func (s *Service) isEmpty() (bool, error) {
	var count int 
	query := `SELECT COUNT(*) FROM taxonomy`
	err := s.Db.QueryRow(query).Scan(&count)
	if err != nil {
		return false, err 
	}
	return count == 0, nil
}

func (s *Service) CreateTaxonomy() error{
	query := `CREATE TABLE IF NOT EXISTS taxonomy(
		id INTEGER PRIMARY KEY,
		speciescode TEXT NOT NULL UNIQUE,
		scientificname TEXT NOT NULL UNIQUE,	
		commonname TEXT NOT NULL
	);`
	_, err := s.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("Не удалось создать базу данных таксономии: %w", err)
	}
	return nil
}

func (s *Service) SaveTaxonomy(ctx context.Context, birds []internal.Bird) error {
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Не удалось начать транзакцию: %v", err)
	}
	defer tx.Rollback()
	
	stmt, err := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO taxonomy (speciescode, scientificname, commonname)
	VALUES (?, ?, ?)
	`)
	if err != nil {
		return err
	} 
	defer stmt.Close()
	for _, bird := range birds {
		_, err := stmt.ExecContext(ctx, bird.SpeciesCode, bird.SciName, bird.ComName)
		if err != nil {
			return fmt.Errorf("Ошибка при вставке птицы %s в таксономию: %w", bird.ComName, err)
		}
	}
	return tx.Commit()
}

func (s *Service) GetBird(speciesCode string) (internal.Bird, error) {
	var bird internal.Bird 
	query := `SELECT speciescode, scientificname, commonname FROM taxonomy WHERE speciescode = ?`
	err := s.Db.QueryRow(query, speciesCode).Scan(&bird.SpeciesCode, &bird.SciName, &bird.ComName)	
	if err != nil {
		return internal.Bird{}, fmt.Errorf("Ошибка получения птицы из таксономии: %w", err)
	}
	return bird, nil
} 

	func (s *Service) SearchBirds(input string) (birds []internal.Bird, e error) {
	query := `SELECT speciescode, scientificname, commonname FROM taxonomy 
	WHERE scientificname LIKE ? OR commonname LIKE ?
	LIMIT 10`

	pattern := "%" + input + "%"
	birds = []internal.Bird{}

	rows, err := s.Db.Query(query, pattern, pattern)
	if (err != nil) {
		return nil, fmt.Errorf("Ошибка при поиске птиц: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		bird := internal.Bird{}
		err := rows.Scan(&bird.SpeciesCode, &bird.SciName, &bird.ComName)
		if (err != nil) {
			fmt.Println(err)
			continue
		}
		birds = append(birds, bird)
	}
	return birds, nil
}


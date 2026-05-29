package taxonomy

import (
	"birdmap/internal"
	"context"
	"database/sql"
)

type Service struct {
	Db *sql.DB
}

func NewService(database *sql.DB) *Service {
	return &Service{
		Db: database, 
	}
}

func (s *Service) SyncTaxonomy(key string, ctx context.Context) error {
	s.CreateTaxonomy()
	empty, err := s.isEmpty()
	
	if err != nil {
		return err 
	}
	if !empty {
		return nil 
	}
	
	var birds []internal.Bird
	birds, err = s.getTaxonomy(key)
	if err != nil {
		return err
	}
	err = s.SaveTaxonomy(ctx, birds)
	if err != nil {
		return err 
	}
	return nil
	
}
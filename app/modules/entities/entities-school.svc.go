package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.SchoolEntity = (*Service)(nil)

func (s *Service) CreateSchool(ctx context.Context, school *ent.School) (*ent.School, error) {
	if _, err := s.db.NewInsert().Model(school).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return school, nil
}

func (s *Service) GetSchoolByID(ctx context.Context, id uuid.UUID) (*ent.School, error) {
	school := new(ent.School)
	if err := s.db.NewSelect().Model(school).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return school, nil
}

func (s *Service) UpdateSchoolByID(ctx context.Context, id uuid.UUID, school *ent.School) (*ent.School, error) {
	updated := new(ent.School)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("name = ?", school.Name).
		Set("logo_url = ?", school.LogoURL).
		Set("theme_color = ?", school.ThemeColor).
		Set("address = ?", school.Address).
		Set("description = ?", school.Description).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) DeleteSchoolByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.School)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListSchools(ctx context.Context) ([]*ent.School, error) {
	var schools []*ent.School
	if err := s.db.NewSelect().Model(&schools).Order("name ASC").Scan(ctx); err != nil {
		return nil, err
	}
	return schools, nil
}

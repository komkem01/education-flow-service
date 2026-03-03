package entities

import (
	"context"
	"errors"
	"strings"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/uptrace/bun"
)

var _ entitiesinf.MemberRoleEntity = (*Service)(nil)

func (s *Service) AddMemberRole(ctx context.Context, memberID uuid.UUID, role ent.MemberRole) error {
	link := &ent.MemberRoleLink{MemberID: memberID, Role: role}
	if _, err := s.db.NewInsert().Model(link).Exec(ctx); err != nil {
		if isMemberRoleDuplicateError(err) {
			return nil
		}
		return err
	}

	return nil
}

func (s *Service) ListMemberRolesByMemberID(ctx context.Context, memberID uuid.UUID) ([]ent.MemberRole, error) {
	var links []*ent.MemberRoleLink
	if err := s.db.NewSelect().
		Model(&links).
		Where("member_id = ?", memberID).
		Order("created_at ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	roles := make([]ent.MemberRole, 0, len(links))
	for _, link := range links {
		roles = append(roles, link.Role)
	}

	return roles, nil
}

func (s *Service) MemberHasAnyRole(ctx context.Context, memberID uuid.UUID, roles []ent.MemberRole) (bool, error) {
	if len(roles) == 0 {
		return false, nil
	}

	return s.db.NewSelect().
		Model((*ent.MemberRoleLink)(nil)).
		Where("member_id = ?", memberID).
		Where("role IN (?)", bun.In(roles)).
		Exists(ctx)
}

func isMemberRoleDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Code != "23505" {
		return false
	}

	constraint := strings.ToLower(pgErr.ConstraintName)
	return strings.Contains(constraint, "uq_member_roles_member_id_role")
}

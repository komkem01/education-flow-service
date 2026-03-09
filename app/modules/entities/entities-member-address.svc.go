package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ entitiesinf.MemberAddressEntity = (*Service)(nil)

func (s *Service) ListMemberAddressesByMemberID(ctx context.Context, memberID uuid.UUID) ([]*ent.MemberAddress, error) {
	var addresses []*ent.MemberAddress
	if err := s.db.NewSelect().
		Model(&addresses).
		Where("member_id = ?", memberID).
		OrderExpr("is_primary DESC, sort_order ASC, created_at ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (s *Service) ReplaceMemberAddressesByMemberID(ctx context.Context, memberID uuid.UUID, addresses []*ent.MemberAddress) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().
			Model((*ent.MemberAddress)(nil)).
			Where("member_id = ?", memberID).
			ForceDelete().
			Exec(ctx); err != nil {
			return err
		}

		for _, item := range addresses {
			item.MemberID = memberID
			if _, err := tx.NewInsert().Model(item).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

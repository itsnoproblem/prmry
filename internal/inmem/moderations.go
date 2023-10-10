package inmem

import (
	"context"

	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/moderation"
)

type ModerationMemoryRepo struct {
	moderations map[string]moderation.Moderation // keyed by moderationID
}

func NewModerationMemoryRepo() *ModerationMemoryRepo {
	return &ModerationMemoryRepo{
		moderations: make(map[string]moderation.Moderation),
	}
}

func (r *ModerationMemoryRepo) Insert(ctx context.Context, mod moderation.Moderation) error {
	if _, ok := r.moderations[mod.ID]; ok {
		return errors.New("moderation already exists")
	}

	r.moderations[mod.ID] = mod

	return nil
}

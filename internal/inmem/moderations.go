package inmem

import (
	"context"

	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/interaction"
)

type ModerationMemoryRepo struct {
	moderations map[string]interaction.Moderation // keyed by moderationID
}

func NewModerationMemoryRepo() *ModerationMemoryRepo {
	return &ModerationMemoryRepo{
		moderations: make(map[string]interaction.Moderation),
	}
}

func (r *ModerationMemoryRepo) Add(ctx context.Context, mod interaction.Moderation) error {
	if _, ok := r.moderations[mod.ID]; ok {
		return errors.New("moderation already exists")
	}

	r.moderations[mod.ID] = mod

	return nil
}

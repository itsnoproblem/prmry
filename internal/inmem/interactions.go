package inmem

import (
	"context"
	"github.com/itsnoproblem/prmry/internal/interaction"
	"github.com/pkg/errors"
)

type InteractionMemoryRepo struct {
	interactions     map[string]interaction.Interaction // keyed by interactionID
	userInteractions map[string][]interaction.Summary   // keyed by userID
}

func NewInteractionMemoryRepo() *InteractionMemoryRepo {
	return &InteractionMemoryRepo{
		interactions:     make(map[string]interaction.Interaction),
		userInteractions: make(map[string][]interaction.Summary),
	}
}

func (r *InteractionMemoryRepo) Add(ctx context.Context, in interaction.Interaction) (id string, err error) {
	if _, ok := r.interactions[in.ID]; ok {
		return "", errors.New("interaction already exists")
	}

	r.interactions[in.ID] = in

	// Generate a summary and add it to the userInteractions map
	summary := interaction.Summary{
		ID:             in.ID,
		FlowID:         in.FlowID,
		FlowName:       in.FlowName,
		Model:          in.Request.Model,
		Prompt:         in.PromptHTML(),
		TokensUsed:     in.Response.Usage.TotalTokens,
		ResponseLength: len(in.ResponseHTML()),
		CreatedAt:      in.CreatedAt,
		UserID:         in.UserID,
	}
	r.userInteractions[in.UserID] = append(r.userInteractions[in.UserID], summary)

	return in.ID, nil
}

func (r *InteractionMemoryRepo) Remove(ctx context.Context, id string) error {
	in, ok := r.interactions[id]
	if !ok {
		return errors.New("interaction not found")
	}

	delete(r.interactions, id)

	// Delete summary from userInteractions
	for i, s := range r.userInteractions[in.UserID] {
		if s.ID == id {
			r.userInteractions[in.UserID] = append(r.userInteractions[in.UserID][:i], r.userInteractions[in.UserID][i+1:]...)
			break
		}
	}

	return nil
}

func (r *InteractionMemoryRepo) SummariesForUser(ctx context.Context, userID string) ([]interaction.Summary, error) {
	summaries, ok := r.userInteractions[userID]
	if !ok || len(summaries) == 0 {
		return nil, errors.New("no summaries found for the given user")
	}

	return summaries, nil
}

func (r *InteractionMemoryRepo) Interaction(ctx context.Context, id string) (interaction.Interaction, error) {
	in, ok := r.interactions[id]
	if !ok {
		return interaction.Interaction{}, errors.New("interaction not found")
	}

	return in, nil
}

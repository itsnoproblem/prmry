package interacting

import (
	"context"
	"fmt"

	"github.com/a-h/templ"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/components/chat"
	"github.com/itsnoproblem/prmry/internal/components/interactions"
	"github.com/itsnoproblem/prmry/internal/interaction"
)

func formatInteractionSummaries(ctx context.Context, response interface{}) (components.Component, error) {
	res, ok := response.([]interaction.Summary)
	if !ok {
		return nil, fmt.Errorf("formatInteractionSummaries: failed to parse response")
	}

	cmp := interactions.NewListView(res)
	cmp.SetUser(auth.UserFromContext(ctx))

	fragment := interactions.InteractionsList(cmp)
	page := interactions.InteractionsListPage(cmp)
	cmp.SetTemplates(page, fragment)

	return &cmp, nil
}

func formatGetInteractionResponse(ctx context.Context, response interface{}) (components.Component, error) {
	ixn, ok := response.(interaction.Interaction)
	if !ok {
		return nil, fmt.Errorf("formatGetInteractionResponse: failed to parse response")
	}

	cmp := chat.NewChatDetailView(ixn)
	cmp.SetUser(auth.UserFromContext(ctx))

	page := interactions.InteractionDetailPage(cmp)
	fragment := interactions.InteractionDetail(cmp)
	cmp.SetTemplates(page, fragment)

	return &cmp, nil
}

func formatCreateInteractionResponse(ctx context.Context, response interface{}) (components.Component, error) {
	ixn, ok := response.(interaction.Interaction)
	if !ok {
		return nil, fmt.Errorf("formatGetInteractionResponse: failed to parse response")
	}

	cmp := chat.ChatResponseView{
		Interaction: chat.NewChatDetailView(ixn),
	}
	cmp.SetUser(auth.UserFromContext(ctx))

	page := templ.NopComponent
	fragment := chat.ChatResponse(cmp)
	cmp.SetTemplates(page, fragment)

	return &cmp, nil
}

func formatChatPromptResponse(ctx context.Context, response interface{}) (components.Component, error) {
	res, ok := response.(chatPromptResponse)
	if !ok {
		return nil, fmt.Errorf("formatChatPromptResponse: failed to parse response")
	}

	cmp := chat.ChatControlsView{
		FlowSelector: chat.NewFlowSelector(res.Flows, res.SelectedFlow),
	}
	cmp.SetUser(auth.UserFromContext(ctx))

	chatFragment := chat.ChatConsole(cmp)
	chatPage := chat.ChatPage(cmp)
	cmp.SetTemplates(chatPage, chatFragment)

	return &cmp, nil
}

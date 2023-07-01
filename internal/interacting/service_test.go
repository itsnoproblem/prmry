package interacting

import (
	"context"
	"reflect"
	"testing"

	gogpt "github.com/sashabaranov/go-gpt3"

	"github.com/itsnoproblem/prmry/internal/interaction"
)

func TestNewService(t *testing.T) {
	type args struct {
		c *gogpt.Client
		r InteractionRepo
		m ModerationRepo
	}
	tests := []struct {
		name string
		args args
		want service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewService(tt.args.c, tt.args.r, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Interaction(t *testing.T) {
	type fields struct {
		gptClient   *gogpt.Client
		history     InteractionRepo
		moderations ModerationRepo
	}
	type args struct {
		ctx           context.Context
		interactionID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interaction.Interaction
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				gptClient:   tt.fields.gptClient,
				history:     tt.fields.history,
				moderations: tt.fields.moderations,
			}
			got, err := s.Interaction(tt.args.ctx, tt.args.interactionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Interaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Interaction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Interactions(t *testing.T) {
	type fields struct {
		gptClient   *gogpt.Client
		history     InteractionRepo
		moderations ModerationRepo
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []interaction.Summary
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				gptClient:   tt.fields.gptClient,
				history:     tt.fields.history,
				moderations: tt.fields.moderations,
			}
			got, err := s.Interactions(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Interactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Interactions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Moderation(t *testing.T) {
	type fields struct {
		gptClient   *gogpt.Client
		history     InteractionRepo
		moderations ModerationRepo
	}
	type args struct {
		ctx           context.Context
		interactionID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interaction.Moderation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				gptClient:   tt.fields.gptClient,
				history:     tt.fields.history,
				moderations: tt.fields.moderations,
			}
			got, err := s.Moderation(tt.args.ctx, tt.args.interactionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Moderation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Moderation() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_ModerationByID(t *testing.T) {
	type fields struct {
		gptClient   *gogpt.Client
		history     InteractionRepo
		moderations ModerationRepo
	}
	type args struct {
		ctx          context.Context
		moderationID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interaction.Moderation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				gptClient:   tt.fields.gptClient,
				history:     tt.fields.history,
				moderations: tt.fields.moderations,
			}
			got, err := s.ModerationByID(tt.args.ctx, tt.args.moderationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ModerationByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModerationByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_RespondToMessage(t *testing.T) {
	type fields struct {
		gptClient   *gogpt.Client
		history     InteractionRepo
		moderations ModerationRepo
	}
	type args struct {
		ctx context.Context
		msg string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				gptClient:   tt.fields.gptClient,
				history:     tt.fields.history,
				moderations: tt.fields.moderations,
			}
			got, err := s.GenerateResponse(tt.args.ctx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

package rgb_test

import (
	"github.com/itsnoproblem/prmry/internal/rgb"
	"testing"
)

func TestPersona_Respond(t *testing.T) {
	chefPersona := newPersonaItalianChef()
	copPersona := newPersonaCop()

	tt := []struct {
		name    string
		persona rgb.Persona
		prompt  string
		expect  string
	}{
		{
			name:    "test Contains rules",
			persona: chefPersona,
			prompt:  "I think rigatoni is the best type of pasta!",
			expect:  "Respond in the voice of a cartoonish italian chef",
		},
		{
			name:    "test multiple rules",
			persona: chefPersona,
			prompt:  "I think spaghetti is the best type of pasta!",
			expect:  "Respond in the voice of chef boyardee",
		},
		{
			name:    "test StartsWith rules",
			persona: copPersona,
			prompt:  "Excuse me mr. officer man",
			expect:  "Respond in the voice of a no-nonsense cop",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			response, err := tc.persona.Respond(tc.prompt)
			if err != nil {
				t.Errorf("%s: %s", tc.name, err)
			}

			if response != tc.expect {
				t.Errorf("%s: got [%s] expected [%s]", tc.name, response, tc.expect)
			}
		})
	}
}

func newPersonaItalianChef() rgb.Persona {
	return rgb.Persona{
		ID:   "123456",
		Name: "Italian Chef",
		Rules: []rgb.Rule{
			{
				RequireAll: false,
				Response:   "Respond in the voice of a cartoonish italian chef",
				Conditions: []rgb.Condition{
					rgb.NewPromptCondition(rgb.ConditionTypeContains, "rigatoni"),
					rgb.NewPromptCondition(rgb.ConditionTypeContains, "penne"),
					rgb.NewPromptCondition(rgb.ConditionTypeContains, "fettuccine"),
					rgb.NewPromptCondition(rgb.ConditionTypeContains, "capellini"),
				},
			},
			{
				RequireAll: true,
				Response:   "Respond in the voice of chef boyardee",
				Conditions: []rgb.Condition{
					rgb.NewPromptCondition(rgb.ConditionTypeContains, "spaghetti"),
				},
			},
		},
	}
}

func newPersonaCop() rgb.Persona {
	return rgb.Persona{
		ID:   "123456",
		Name: "Cop Persona",
		Rules: []rgb.Rule{
			{
				RequireAll: true,
				Response:   "Respond in the voice of a no-nonsense cop",
				Conditions: []rgb.Condition{
					rgb.NewPromptCondition(rgb.ConditionTypeStartsWith, "Excuse me"),
					rgb.NewPromptCondition(rgb.ConditionTypeContains, "Officer"),
				},
			},
			{
				RequireAll: true,
				Response:   "Respond in the voice of an approval-seeking mall cop",
				Conditions: []rgb.Condition{
					rgb.NewPromptCondition(rgb.ConditionTypeStartsWith, "Wait"),
				},
			},
		},
	}
}

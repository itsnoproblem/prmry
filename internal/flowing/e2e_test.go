package flowing_test

import (
    "bytes"
    "context"
    "encoding/json"
    "io"
    "net/http"
    "net/http/httptest"
    "net/url"
    "regexp"
    "testing"

    "github.com/approvals/go-approval-tests"
    "github.com/approvals/go-approval-tests/reporters"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/render"
    "github.com/pkg/errors"

    "github.com/itsnoproblem/prmry/internal/auth"
    "github.com/itsnoproblem/prmry/internal/flow"
    "github.com/itsnoproblem/prmry/internal/flowing"
    "github.com/itsnoproblem/prmry/internal/htmx"
    "github.com/itsnoproblem/prmry/internal/inmem"
)

func init() {
    approvals.UseFolder("testdata")
}

type testHelper struct {
    Renderer    flowing.Renderer
    FlowsRepo   flowing.FlowsRepo
    FlowService flowing.Service
    Router      *chi.Mux
}

func TestFlowing(t *testing.T) {
    reporter := approvals.UseReporter(reporters.NewFileMergeReporter())
    defer reporter.Close()

    helper := newTestHelper()

    helper.Router.Route("/", flowing.RouteHandler(helper.FlowService, helper.Renderer))
    server := httptest.NewServer(helper.Router)
    defer server.Close()

    client := server.Client()
    user := auth.TestModeUser()

    if err := createSomeFlows(context.Background(), user.ID, helper.FlowsRepo); err != nil {
        t.Fatalf(err.Error())
    }

    testCases := []struct {
        Name           string
        Method         string
        Path           string
        FullPage       bool
        WantStatusCode int
        RequestPayload interface{}
    }{
        {
            Name:           "FlowBuilderPage",
            Method:         http.MethodGet,
            Path:           "/flow-builder",
            FullPage:       true,
            WantStatusCode: http.StatusOK,
            RequestPayload: nil,
        },
        {
            Name:           "FlowBuilderFragment",
            Method:         http.MethodGet,
            Path:           "/flow-builder",
            FullPage:       false,
            WantStatusCode: http.StatusOK,
            RequestPayload: nil,
        },
        {
            Name:           "FlowBuilderFragment Add Trigger",
            Method:         http.MethodPost,
            Path:           "/flow-builder/rules",
            FullPage:       false,
            WantStatusCode: http.StatusOK,
            RequestPayload: flowing.FlowBuilderFormRequest{
                ID:         "",
                Name:       "",
                Prompt:     "",
                RequireAll: "false",
            },
        },
        {
            Name:           "FlowBuilderFragment Delete Trigger",
            Method:         http.MethodDelete,
            Path:           "/flow-builder/rules/0",
            FullPage:       false,
            WantStatusCode: http.StatusOK,
            RequestPayload: flowing.FlowBuilderFormRequest{
                ID:         "",
                Name:       "",
                Prompt:     "",
                RequireAll: "false",
            },
        },
        {
            Name:           "FlowBuilderFragment Update Prompt",
            Method:         http.MethodPut,
            Path:           "/flow-builder/prompt",
            FullPage:       false,
            WantStatusCode: http.StatusOK,
            RequestPayload: flowing.FlowBuilderFormRequest{
                ID:         "",
                Name:       "Test Flow",
                Prompt:     "The quick brown %s jumps over the hill",
                RequireAll: "false",
            },
        },
        {
            Name:           "FlowBuilderFragment Select Prompt Replacement",
            Method:         http.MethodPut,
            Path:           "/flow-builder/prompt",
            FullPage:       false,
            WantStatusCode: http.StatusOK,
            RequestPayload: flowing.FlowBuilderFormRequest{
                ID:         "",
                Name:       "Hello Flow",
                Prompt:     "respond to this text with a greeting: %s",
                RequireAll: "false",
                PromptArgs: "input message",
            },
        },
        {
            Name:           "FlowBuilderFragment Select Multiple Prompt Replacements",
            Method:         http.MethodPut,
            Path:           "/flow-builder/prompt",
            FullPage:       false,
            WantStatusCode: http.StatusOK,
            RequestPayload: flowing.FlowBuilderFormRequest{
                ID:              "",
                Name:            "Hello Flow",
                FieldNames:      "input message",
                ConditionTypes:  "does not equal",
                ConditionValues: "X",
                Prompt:          "respond to this text with a greeting: %s",
                PromptArgs: []string{
                    "interaction result from another flow",
                    "input message",
                },
                PromptArgFlows: "123",
                RequireAll:     "false",
            },
        },
        {
            Name:           "FlowBuilderFragment SaveFlow",
            Method:         http.MethodPost,
            Path:           "/flows",
            FullPage:       false,
            WantStatusCode: http.StatusOK,
            RequestPayload: flowing.FlowBuilderFormRequest{
                ID:              "",
                Name:            "Hello Flow",
                FieldNames:      "input message",
                ConditionTypes:  "does not equal",
                ConditionValues: "X",
                Prompt:          "respond to this text with a greeting: %s",
                PromptArgs: []string{
                    "interaction result from another flow",
                    "input message",
                },
                PromptArgFlows: "123",
                RequireAll:     "false",
            },
        },
        {
            Name:           "ListFlowsPage",
            Method:         http.MethodGet,
            Path:           "/flows",
            FullPage:       true,
            WantStatusCode: http.StatusOK,
            RequestPayload: nil,
        },
        {
            Name:           "ListFlowsFragment",
            Method:         http.MethodGet,
            Path:           "/flows",
            FullPage:       false,
            WantStatusCode: http.StatusOK,
            RequestPayload: nil,
        },
        {
            Name:           "EditFlowPage",
            Method:         http.MethodGet,
            Path:           "/flows/123/edit",
            FullPage:       true,
            WantStatusCode: http.StatusOK,
            RequestPayload: nil,
        },
        {
            Name:           "EditFlowFragment",
            Method:         http.MethodGet,
            Path:           "/flows/123/edit",
            FullPage:       false,
            WantStatusCode: http.StatusOK,
            RequestPayload: nil,
        },
    }

    for _, tt := range testCases {

        t.Run(tt.Name, func(t *testing.T) {
            var reqPayload json.RawMessage

            u, err := url.Parse(server.URL + tt.Path)
            if err != nil {
                t.Fatalf("parse url: %s", err)
            }

            if tt.RequestPayload != nil {
                reqPayload, err = json.Marshal(tt.RequestPayload)
                if err != nil {
                    t.Fatalf("encoding request payload: %s", err.Error())
                }
            }
            req, err := http.NewRequest(tt.Method, u.String(), bytes.NewReader(reqPayload))

            if !tt.FullPage {
                if req.Header == nil {
                    req.Header = make(http.Header)
                }
                req.Header.Add(htmx.HeaderHXRequest, "true")
            }

            var resp *http.Response

            for resp == nil || resp.Header.Get(htmx.HeaderHXRedirect) != "" {
                if resp != nil {
                    req.Method = http.MethodGet
                    originalURL := req.URL
                    originalURL.Path = resp.Header.Get(htmx.HeaderHXRedirect)

                    t.Logf("[%s] %s", req.Method, originalURL.String())
                }

                resp, err = client.Do(req)
                if err != nil {
                    t.Fatalf("Get error: %v", err)
                }
            }

            if resp.StatusCode != tt.WantStatusCode {
                t.Fatalf("expected response code %d but received %d", tt.WantStatusCode, resp.StatusCode)
            }

            data, err := io.ReadAll(resp.Body)
            if err != nil {
                t.Errorf("io.ReadAll error: %v", err)
            }

            dateScrubber, err := regexp.Compile("last-changed\">[\t\n\f\r ]?(.*)[\t\n\f\r ]+?<")
            if err != nil {
                t.Fatalf("dateScrubber: %s", err)
            }

            uuidScrubber, err := regexp.Compile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
            if err != nil {
                t.Fatalf("uuidScrubber: %s", err)
            }

            opts := approvals.Options().WithRegexScrubber(dateScrubber, "last-changed\">DATE PLACEHOLDER<")
            opts = opts.WithRegexScrubber(uuidScrubber, "UUID-PLACEHOLDER")
            opts = opts.WithExtension("html")
            approvals.VerifyString(t, string(data[:]), opts)
        })
    }

}

func newTestHelper() testHelper {
    r := chi.NewRouter()
    r.Use(auth.Middleware(auth.Byte32{}, false))
    r.Use(render.SetContentType(render.ContentTypeHTML))
    r.Use(htmx.Middleware)
    r.Use(auth.TestUserMiddleware)

    flowsRepo := inmem.NewFlowsRepo()
    funnelsRepo := inmem.NewFunnelsRepo()

    return testHelper{
        Renderer:    htmx.NewRenderer(),
        FlowsRepo:   flowsRepo,
        FlowService: flowing.NewService(flowsRepo, funnelsRepo, "http://localhost:8080"),
        Router:      r,
    }
}

type requestPayloadNewRule struct {
    Id         string `json:"id"`
    Name       string `json:"name"`
    RequireAll string `json:"requireAll"`
    Prompt     string `json:"prompt"`
}

func createSomeFlows(ctx context.Context, userID string, repo flowing.FlowsRepo) error {
    fakeRules := []flow.Trigger{
        flow.Trigger{
            Field: flow.Field{
                Source: flow.FieldSourceInput,
            },
            Condition: "equals",
            Value:     "test",
        },
    }

    flows := []flow.Flow{
        flow.Flow{
            ID:       "123",
            UserID:   userID,
            Name:     "Test Flow A",
            Triggers: fakeRules,
            Prompt:   "Tell me a story based on this text: %s",
            PromptArgs: []flow.Field{
                {
                    Source: flow.FieldSourceInput,
                },
            },
        },
        flow.Flow{
            ID:       "567",
            UserID:   userID,
            Name:     "Test Flow B",
            Triggers: fakeRules,
        },
        flow.Flow{
            ID:       "998",
            UserID:   userID,
            Name:     "Test Flow C",
            Triggers: nil,
        },
    }

    for _, f := range flows {
        if err := repo.InsertFlow(ctx, f); err != nil {
            return errors.Wrap(err, "createSomeFlows")
        }
    }

    return nil
}

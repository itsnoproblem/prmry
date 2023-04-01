package profiling_test

//
//import (
//	"bytes"
//	approvals "github.com/approvals/go-approval-tests"
//	"github.com/approvals/go-approval-tests/reporters"
//	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/profiling"
//	"github.com/markbates/goth"
//	"github.com/stretchr/testify/require"
//	"io"
//	"testing"
//)
//
//const (
//	testDataDir = "testdata"
//)
//
//func TestRenderHome(t *testing.T) {
//	r := approvals.UseReporter(reporters.NewGoLandReporter())
//	defer r.Close()
//	approvals.UseFolder(testDataDir)
//
//	renderer, err := profiling.NewRenderer()
//	require.NoError(t, err)
//
//	tt := []struct {
//		Description string
//		RenderFn    func(w io.Writer, user *goth.User) error
//	}{
//		{
//			Description: "HTML document of home page",
//			RenderFn:    renderer.RenderHomePage,
//		},
//		{
//			Description: "HTML fragment of home page content",
//			RenderFn:    renderer.RenderHome,
//		},
//	}
//
//	for _, tc := range tt {
//		t.Run(tc.Description, func(t *testing.T) {
//			buf := bytes.Buffer{}
//			err = tc.RenderFn(&buf, &goth.User{})
//			if err != nil {
//				t.Errorf("ERROR in %s: %s", tc.Description, err)
//			}
//
//			approvals.VerifyString(t, buf.String())
//		})
//	}
//}

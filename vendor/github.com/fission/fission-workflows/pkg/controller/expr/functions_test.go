package expr

import (
	"testing"

	"github.com/fission/fission-workflows/pkg/types"
	"github.com/fission/fission-workflows/pkg/types/typedvalues"
	"github.com/fission/fission-workflows/pkg/util"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
)

func makeTestScope() *Scope {
	scope, _ := NewScope(nil, &types.Workflow{
		Metadata: &types.ObjectMetadata{
			Id:        "testWorkflow",
			CreatedAt: ptypes.TimestampNow(),
		},
		Status: &types.WorkflowStatus{
			Status:    types.WorkflowStatus_READY,
			UpdatedAt: ptypes.TimestampNow(),
			Tasks: map[string]*types.TaskStatus{
				"TaskA": {
					FnRef: &types.FnRef{
						Runtime: "fission",
						ID:      "resolvedFissionFunction",
					},
				},
			},
		},
		Spec: &types.WorkflowSpec{
			ApiVersion: "1",
			OutputTask: "TaskA",
			Tasks: map[string]*types.TaskSpec{
				"TaskA": {
					FunctionRef: "fissionFunction",
					Inputs: map[string]*typedvalues.TypedValue{
						types.InputMain: typedvalues.MustWrap("input-default"),
						"otherInput":    typedvalues.MustWrap("input-otherInput"),
					},
				},
			},
		},
	}, &types.WorkflowInvocation{
		Metadata: &types.ObjectMetadata{
			Id:        "testWorkflowInvocation",
			CreatedAt: ptypes.TimestampNow(),
		},
		Spec: &types.WorkflowInvocationSpec{
			WorkflowId: "testWorkflow",
			Inputs: map[string]*typedvalues.TypedValue{
				types.InputMain: typedvalues.MustWrap("body"),
				"headers":       typedvalues.MustWrap("http-headers"),
			},
		},
		Status: &types.WorkflowInvocationStatus{
			Status: types.WorkflowInvocationStatus_IN_PROGRESS,
			Tasks: map[string]*types.TaskInvocation{
				"TaskA": {
					Spec: &types.TaskInvocationSpec{},
					Status: &types.TaskInvocationStatus{
						Output: typedvalues.MustWrap("some output"),
					},
				},
			},
		},
	})
	return scope
}

func TestOutputFn_Apply_OneArgument(t *testing.T) {
	parser := NewJavascriptExpressionParser()

	testScope := makeTestScope()
	result, err := parser.Resolve(testScope, "", mustParseExpr("{ output('TaskA') }"))
	assert.NoError(t, err)

	i := typedvalues.MustUnwrap(result)

	assert.Equal(t, testScope.Tasks["TaskA"].Output, i)
}

func TestOutputFn_Apply_NoArgument(t *testing.T) {
	parser := NewJavascriptExpressionParser()

	testScope := makeTestScope()
	result, err := parser.Resolve(testScope, "TaskA", mustParseExpr("{ output() }"))
	assert.NoError(t, err)

	i := typedvalues.MustUnwrap(result)

	assert.Equal(t, testScope.Tasks["TaskA"].Output, i)
}

func TestInputFn_Apply_NoArgument(t *testing.T) {
	parser := NewJavascriptExpressionParser()

	testScope := makeTestScope()
	result, err := parser.Resolve(testScope, "TaskA", mustParseExpr("{ input() }"))
	assert.NoError(t, err)

	i := typedvalues.MustUnwrap(result)

	assert.Equal(t, "input-default", i)
}

func TestInputFn_Apply_OneArgument(t *testing.T) {
	parser := NewJavascriptExpressionParser()

	testScope := makeTestScope()
	result, err := parser.Resolve(testScope, "", mustParseExpr("{ input('TaskA') }"))
	assert.NoError(t, err)

	i := typedvalues.MustUnwrap(result)

	assert.Equal(t, "input-default", i)
}

func TestInputFn_Apply_TwoArguments(t *testing.T) {
	parser := NewJavascriptExpressionParser()

	testScope := makeTestScope()
	result, err := parser.Resolve(testScope, "", mustParseExpr("{ input('TaskA', 'otherInput') }"))
	assert.NoError(t, err)

	i := typedvalues.MustUnwrap(result)

	assert.Equal(t, "input-otherInput", i)
}

func TestParamFn_Apply_NoArgument(t *testing.T) {
	parser := NewJavascriptExpressionParser()
	testScope := makeTestScope()
	result, err := parser.Resolve(testScope, "", mustParseExpr("{ param() }"))
	assert.NoError(t, err)
	assert.Equal(t, "body", typedvalues.MustUnwrap(result))
}

func TestParamFn_Apply_OneArgument(t *testing.T) {
	parser := NewJavascriptExpressionParser()
	testScope := makeTestScope()
	result, err := parser.Resolve(testScope, "", mustParseExpr("{ param('headers') }"))
	assert.NoError(t, err)
	assert.Equal(t, "http-headers", typedvalues.MustUnwrap(result))
}

func TestUidFn_Apply(t *testing.T) {
	parser := NewJavascriptExpressionParser()
	testScope := makeTestScope()
	result, err := parser.Resolve(testScope, "", mustParseExpr("{ uid() }"))
	assert.NoError(t, err)
	assert.NotEmpty(t, typedvalues.MustUnwrap(result))
}

func TestTaskFn_Apply_OneArgument(t *testing.T) {
	parser := NewJavascriptExpressionParser()

	testScope := makeTestScope()
	result, err := parser.Resolve(testScope, "", mustParseExpr("{ task('TaskA') }"))
	assert.NoError(t, err)
	i := typedvalues.MustUnwrap(result)

	assert.Equal(t, util.MustConvertStructsToMap(testScope.Tasks["TaskA"]), i)
}

func TestTaskFn_Apply_NoArgument(t *testing.T) {
	parser := NewJavascriptExpressionParser()

	testScope := makeTestScope()
	result, err := parser.Resolve(testScope, "TaskA", mustParseExpr("{ task() }"))
	assert.NoError(t, err)

	i := typedvalues.MustUnwrap(result)

	assert.Equal(t, util.MustConvertStructsToMap(testScope.Tasks["TaskA"]), i)
}

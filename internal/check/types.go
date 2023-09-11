package check

import (
	"context"

	"github.com/opdev/operator-certification/internal/operatorsdk"
)

type operatorSdk interface {
	Scorecard(context.Context, string, operatorsdk.OperatorSdkScorecardOptions) (*operatorsdk.OperatorSdkScorecardReport, error)
}

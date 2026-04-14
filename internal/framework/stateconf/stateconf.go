package stateconf

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/samber/lo"
)

const defaultTimeout = 5 * time.Minute

type StateRefreshFunc[S ~string] func(context.Context) (any, S, error)

func States[S ~string](values ...S) []S {
	return values
}

// Generic wrapper of retry.StateChangeConf
type StateChangeConf[S ~string] struct {
	Delay                     time.Duration
	Pending                   []S
	Target                    []S
	Refresh                   StateRefreshFunc[S]
	Timeout                   time.Duration
	MinTimeout                time.Duration
	PollInterval              time.Duration
	NotFoundChecks            int
	ContinuousTargetOccurence int
}

func (c *StateChangeConf[S]) WaitForStateContext(ctx context.Context) (any, error) {
	if c.Timeout == 0 {
		tflog.Debug(ctx, fmt.Sprintf("stateconf: timeout is not defined, setting to default timeout (%v)", defaultTimeout))
		c.Timeout = defaultTimeout
	}
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	conf := &retry.StateChangeConf{
		Delay:                     c.Delay,
		Pending:                   conv(c.Pending),
		Target:                    conv(c.Target),
		Timeout:                   c.Timeout,
		MinTimeout:                c.MinTimeout,
		PollInterval:              c.PollInterval,
		NotFoundChecks:            c.NotFoundChecks,
		ContinuousTargetOccurence: c.ContinuousTargetOccurence,
		Refresh: func() (any, string, error) {
			result, state, err := c.Refresh(ctx)
			return result, string(state), err
		},
	}
	return conf.WaitForStateContext(ctx)
}

func conv[S ~string](ss []S) []string {
	return lo.Map(ss, func(s S, _ int) string {
		return string(s)
	})
}

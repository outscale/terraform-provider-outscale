package logging

import (
	"context"
	"fmt"

	smithylogging "github.com/aws/smithy-go/logging"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type smithyTFLogWrapper struct {
	ctx context.Context
}

var (
	_ smithylogging.Logger        = (*smithyTFLogWrapper)(nil)
	_ smithylogging.ContextLogger = (*smithyTFLogWrapper)(nil)
)

func NewSmithyTFLogWrapper() *smithyTFLogWrapper {
	return &smithyTFLogWrapper{
		ctx: context.Background(),
	}
}

func (t *smithyTFLogWrapper) WithContext(ctx context.Context) smithylogging.Logger {
	return &smithyTFLogWrapper{
		ctx: ctx,
	}
}

func (t *smithyTFLogWrapper) Logf(classification smithylogging.Classification, format string, values ...any) {
	message := fmt.Sprintf(format, values...)
	fields := map[string]any{"message": message}

	if classification == smithylogging.Warn {
		tflog.Warn(t.ctx, "SDK warning", fields)
		return
	}

	tflog.Debug(t.ctx, "SDK log", fields)
}

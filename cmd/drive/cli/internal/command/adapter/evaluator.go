package adapter

import (
	"context"
	"fmt"
	"github.com/khulnasoft/drive/cmd/drive/cli/internal/command/ci"
	"github.com/khulnasoft/drive/drive/image"
	"github.com/khulnasoft/drive/internal/bus"
	"github.com/khulnasoft/drive/internal/bus/event/payload"
	"github.com/khulnasoft/drive/internal/log"
)

type Evaluator interface {
	Evaluate(ctx context.Context, analysis *image.Analysis) ci.Evaluation
}

type evaluationActionObserver struct {
	ci.Evaluator
}

func NewEvaluator(rules []ci.Rule) Evaluator {
	return evaluationActionObserver{
		Evaluator: ci.NewEvaluator(rules),
	}
}

func (c evaluationActionObserver) Evaluate(ctx context.Context, analysis *image.Analysis) ci.Evaluation {
	log.WithFields("image", analysis.Image).Infof("evaluating image")
	mon := bus.StartTask(payload.GenericTask{
		Title: payload.Title{
			Default:      "Evaluating image",
			WhileRunning: "Evaluating image",
			OnSuccess:    "Evaluated image",
		},
		HideOnSuccess:      false,
		HideStageOnSuccess: false,
		ID:                 analysis.Image,
		Context:            fmt.Sprintf("[rules: %d]", len(c.Rules)),
	})
	eval := c.Evaluator.Evaluate(ctx, analysis)
	if eval.Pass {
		mon.SetCompleted()
	} else {
		mon.SetError(fmt.Errorf("failed evaluation"))
	}
	bus.Report(eval.Report)
	return eval
}

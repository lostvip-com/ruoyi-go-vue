package ekuiperclient

import (
	"context"
	"things/internal/common/dtos"
)

type EkuiperClient interface {
	RuleExist(ctx context.Context, ruleId string) (bool, error)
	CreateRule(ctx context.Context, actions []dtos.Actions, ruleId string, sql string) error
	UpdateRule(ctx context.Context, actions []dtos.Actions, ruleId string, sql string) error
	GetRuleStats(ctx context.Context, ruleId string) (map[string]interface{}, error)
	StartRule(ctx context.Context, ruleId string) error
	StopRule(ctx context.Context, ruleId string) error
	RestartRule(ctx context.Context, ruleId string) error
	DeleteRule(ctx context.Context, ruleId string) error
}

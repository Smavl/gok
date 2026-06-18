package config

import (
	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/misc"
	"github.com/smavl/gok/internal/prober/types"
)

func ConfigForMode(mode domain.ProbingMode) (types.ProbeConfig, error) {
	switch mode {
	case domain.Default:
		return DefaultConfig(), nil
	}
	return types.ProbeConfig{}, misc.ErrNoProberMode
}

package checkers

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/ports"
	"github.com/dimiro1/health/redis"
)

type redisChecker struct {
}

func NewRedisChecker() ports.ComponentChecker {
	return &redisChecker{}
}

func (c *redisChecker) Name() string {
	return "redis"
}

func (c *redisChecker) Check() models.Health {
	checker := redis.NewChecker("tcp", ":6379")
	health := checker.Check()
	if health.IsDown() {
		details := "Redis service is down"
		if health.GetInfo("error") != nil {
			details = fmt.Sprintf("%s: %v", details, health.GetInfo("error"))
		}
		return models.Health{
			Status:  constants.StatusDown,
			Details: details,
		}
	}

	return models.Health{
		Status:  constants.StatusUp,
		Details: "Redis service is healthy",
	}
}

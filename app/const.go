package app

import "time"

const (
	AppName = "piston"

	EthUrlKey    = "ethUrl"
	EngineUrlKey = "engineUrl"

	EthUrlDefault    = "http://localhost:8545"
	EngineUrlDefault = "http://localhost:8551"

	PeriodKey     = "period"
	PeriodDefault = 10 * time.Second
	PeriodMinimum = 1 * time.Second

	JWTKey = "jwt"
)

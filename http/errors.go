package http

import "errors"

var (
	// ErrVerificationLevelBounds indicate that you sent a out of bound VerificationLevel
	ErrVerificationLevelBounds = errors.New("VerificationLevel out of bounds, should be between 0 and 3")
	// ErrPruneDaysBounds indicate that you sent an invalid number of days for in PruneMember
	ErrPruneDaysBounds = errors.New("the number of days should be more than or equal to 1")
	// ErrGuildNoIcon indicate that there is no Guild Icon
	ErrGuildNoIcon = errors.New("guild does not have an icon set")
	// ErrGuildNoSplash indicate that there is no Guild Splash
	ErrGuildNoSplash = errors.New("guild does not have a splash set")
)

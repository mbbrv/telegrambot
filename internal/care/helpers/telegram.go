package helpers

import "telegrambot/internal/care/vars"

func GetTimeDivisionInSeconds() int {
	if vars.TimeToSleep%2 == 1 {
		return vars.TimeToSleep/2*60 + 30
	}
	return vars.TimeToSleep / 2 * 60
}

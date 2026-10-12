package legacy

func nox_common_randomInt_415FA0(min, max int) int {
	return GetServer().S().Rand.Logic.IntClamp(min, max)
}

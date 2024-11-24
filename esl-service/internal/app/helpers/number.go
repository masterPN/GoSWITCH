package helpers

func NormalizeDestinationNumber(destination string) string {
	// Begin with 0, means local call
	if len(destination) > 0 && destination[0] == '0' {
		return "66" + destination[1:]
	}
	return destination
}

package lib

func CheckBank(val string) bool {
	banks := []string{"Main", "WE", "TPG"}

	for _, bank := range banks {
		if bank == val {
			return true
		}
	}

	return false
}

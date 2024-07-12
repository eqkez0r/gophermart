package luhn

func CalculateLuhn(number uint64) uint64 {
	checkNumber := checksum(number)

	if checkNumber == 0 {
		return 0
	}
	return 10 - checkNumber
}

func Valid(number uint64) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number uint64) uint64 {
	var luhn uint64

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}

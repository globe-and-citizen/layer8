package utils

func Equal(fst []byte, snd []byte) bool {
	if len(fst) != len(snd) {
		return false
	}

	for i := 0; i < len(fst); i++ {
		if fst[i] != snd[i] {
			return false
		}
	}

	return true
}

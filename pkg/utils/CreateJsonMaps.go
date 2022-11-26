package utils

func CreateJsonMaps(chars string) []string {
	track := 0
	isStringOpened := false
	opened := false
	var objectSlice []byte
	var finishedSlice []string

	for i := 0; i < len(chars); i++ {
		if chars[i] == '"' {
			isStringOpened = !isStringOpened
		}

		if chars[i] == '{' && !isStringOpened {
			track++
			if !opened {
				opened = true
			}
		}
		if chars[i] == '}' && !isStringOpened {
			track--
		}

		if opened {
			objectSlice = append(objectSlice, chars[i])
		}
		if track == 0 && opened {
			opened = false
			finishedSlice = append(finishedSlice, string(objectSlice))
			objectSlice = nil
		}
	}
	return finishedSlice
}

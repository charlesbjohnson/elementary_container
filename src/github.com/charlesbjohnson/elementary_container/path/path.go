package path

func FullExt(path string) string {
	cursor := len(path) - 1

	for i := cursor; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			cursor = i
		}
	}

	return path[cursor:]
}

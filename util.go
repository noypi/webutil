package webutil

func MergeMapStrIf(dst, src map[string]interface{}) {
	for k, v := range src {
		dst[k] = v
	}
}

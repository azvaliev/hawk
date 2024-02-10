package httpfs

import "fmt"

func fileSizeFormat(bytes int64) string {
	const (
		_          = iota // ignore first value by assigning to blank identifier
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
		PB
		EB
		ZB
		YB
	)

	var (
		size   = float64(bytes)
		suffix string
	)

	switch {
	case size >= YB:
		size /= YB
		suffix = "YB"
	case size >= ZB:
		size /= ZB
		suffix = "ZB"
	case size >= EB:
		size /= EB
		suffix = "EB"
	case size >= PB:
		size /= PB
		suffix = "PB"
	case size >= TB:
		size /= TB
		suffix = "TB"
	case size >= GB:
		size /= GB
		suffix = "GB"
	case size >= MB:
		size /= MB
		suffix = "MB"
	case size >= KB:
		size /= KB
		suffix = "KB"
	default:
		suffix = "B"
	}

	return fmt.Sprintf("%.2f%s", size, suffix)
}

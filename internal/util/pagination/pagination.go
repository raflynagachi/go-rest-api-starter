package pagination

// TotalPage count total page for specific limit data per page
func TotalPage(total, limit int64) int64 {
	return (total + int64(limit) - 1) / int64(limit)
}

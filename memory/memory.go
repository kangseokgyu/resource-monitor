package memory

import "fmt"

type Stats struct {
	Total         uint64
	Used          uint64
	UsedPercent   float32
	Free          uint64
	FreePercent   float32
	Cached        uint64
	CachedPercent float32
}

func (s Stats) String() string {
	return fmt.Sprintf("Total: %.2f GB, Used: %.2f GB(%.1f %%), Free: %.2f GB(%.1f %%), Cached: %.2f GB(%.1f %%)",
		float32(s.Total)/1024/1024/1024,
		float32(s.Used)/1024/1024/1024, s.UsedPercent,
		float32(s.Free)/1024/1024/1024, s.FreePercent,
		float32(s.Cached)/1024/1024/1024, s.CachedPercent)
}

func Get() (*Stats, error) {
	return GetStats()
}

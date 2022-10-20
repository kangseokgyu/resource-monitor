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
	return fmt.Sprintf("Total: %d, Used: %d(%.1f %%), Free: %d(%.1f %%), Cached: %d(%.1f %%)", s.Total, s.Used, s.UsedPercent, s.Free, s.FreePercent, s.Cached, s.CachedPercent)
}

func Get() (*Stats, error) {
	return GetStats()
}

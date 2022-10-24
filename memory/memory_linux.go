package memory

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetStats() (*Stats, error) {
	meminfoStats, err := getMeminfoStats()
	if err != nil {
		return nil, err
	}

	return getStats(*&meminfoStats)
}

func getMeminfoStats() (*[]string, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var res []string
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	return &res, nil
}

func getStats(meminfoStatStr *[]string) (*Stats, error) {
	if meminfoStatStr == nil {
		return nil, fmt.Errorf("No input.")
	}

	type MeminfoStat struct {
		MemTotal      uint64
		MemFree       uint64
		MemAvailable  uint64
		MemAvailabled bool
		Buffers       uint64
		Cached        uint64
		SwapCached    uint64
		Active        uint64
		Inactive      uint64
		SwapTotal     uint64
		SwapFree      uint64
		Mapped        uint64
		Shmem         uint64
		Slab          uint64
		PageTables    uint64
		Committed_AS  uint64
		VmallocUsed   uint64
	}
	var meminfoStat MeminfoStat

	type StatFormat struct {
		format string
		value  *uint64
	}
	mapStat := map[string]StatFormat{
		"MemTotal":     StatFormat{"MemTotal: %d kB", &meminfoStat.MemTotal},
		"MemFree":      StatFormat{"MemFree: %d kB", &meminfoStat.MemFree},
		"MemAvailable": StatFormat{"MemAvailable: %d kB", &meminfoStat.MemAvailable},
		"Buffers":      StatFormat{"Buffers: %d kB", &meminfoStat.Buffers},
		"Cached":       StatFormat{"Cached: %d kB", &meminfoStat.Cached},
		"SwapCached":   StatFormat{"SwapCached: %d kB", &meminfoStat.SwapCached},
		"Active":       StatFormat{"Active: %d kB", &meminfoStat.Active},
		"Inactive":     StatFormat{"Inactive: %d kB", &meminfoStat.Inactive},
		"SwapTotal":    StatFormat{"SwapTotal: %d kB", &meminfoStat.SwapTotal},
		"SwapFree":     StatFormat{"SwapFree: %d kB", &meminfoStat.SwapFree},
		"Mapped":       StatFormat{"Mapped: %d kB", &meminfoStat.Mapped},
		"Shmem":        StatFormat{"Shmem: %d kB", &meminfoStat.Shmem},
		"Slab":         StatFormat{"Slab: %d kB", &meminfoStat.Slab},
		"PageTables":   StatFormat{"PageTables: %d kB", &meminfoStat.PageTables},
		"Committed_AS": StatFormat{"Committed_AS: %d kB", &meminfoStat.Committed_AS},
		"VmallocUsed":  StatFormat{"VmallocUsed: %d kB", &meminfoStat.VmallocUsed},
	}

	for _, line := range *meminfoStatStr {
		i := strings.Index(line, ":")
		if i > 0 {
			format, exists := mapStat[line[:i]]
			if exists {
				fmt.Sscanf(line, format.format, format.value)
				*format.value *= 1024

				if line[:i] == "MemAvailable" {
					meminfoStat.MemAvailabled = true
				}
			}
		}
	}

	var stats Stats

	stats.Total = meminfoStat.MemTotal
	stats.Free = meminfoStat.MemFree
	stats.Cached = meminfoStat.Cached

	if meminfoStat.MemAvailabled {
		stats.Used = meminfoStat.MemTotal - meminfoStat.MemAvailable
	} else {
		stats.Used = meminfoStat.MemTotal - meminfoStat.MemFree - meminfoStat.Buffers - meminfoStat.Cached
	}

	stats.UsedPercent = float32(stats.Used) / float32(stats.Total) * 100
	stats.FreePercent = float32(stats.Free) / float32(stats.Total) * 100
	stats.CachedPercent = float32(stats.Cached) / float32(stats.Total) * 100

	return &stats, nil
}

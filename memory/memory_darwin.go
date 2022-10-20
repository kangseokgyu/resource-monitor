package memory

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func GetStats() (*Stats, error) {
	vm_stat, err := getVMStat()
	if err != nil {
		return nil, err
	}

	return getStats(vm_stat)
}

func getVMStat() (*[]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "vm_stat")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err1 := cmd.Start()
	if err1 != nil {
		return nil, err1
	}

	scanner := bufio.NewScanner(out)
	var res []string
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	return &res, nil
}

func getStats(vm_stat_str *[]string) (*Stats, error) {
	var stats Stats

	type VMStat struct {
		PageSize                  uint64
		PagesFree                 uint64
		PagesActive               uint64
		PagesInactive             uint64
		PagesSpeculative          uint64
		PagesThrottled            uint64
		PagesWiredDown            uint64
		PagesPurgeable            uint64
		FileBackedPages           uint64
		PagesOccupiedByCompressor uint64
	}

	var vmStat VMStat

	type StatFormat struct {
		format string
		value  *uint64
	}
	mapStat := map[string]StatFormat{
		"Mach Virtual Memory Statistics": StatFormat{"Mach Virtual Memory Statistics: (page size of %d bytes)", &vmStat.PageSize},
		"Pages free":                     StatFormat{"Pages free: %d.", &vmStat.PagesFree},
		"Pages active":                   StatFormat{"Pages active: %d.", &vmStat.PagesActive},
		"Pages inactive":                 StatFormat{"Pages inactive: %d.", &vmStat.PagesInactive},
		"Pages speculative":              StatFormat{"Pages speculative: %d.", &vmStat.PagesSpeculative},
		"Pages throttled":                StatFormat{"Pages throttled: %d.", &vmStat.PagesThrottled},
		"Pages wired down":               StatFormat{"Pages wired down: %d.", &vmStat.PagesWiredDown},
		"Pages purgeable":                StatFormat{"Pages purgeable: %d.", &vmStat.PagesPurgeable},
		"File-backed pages":              StatFormat{"File-backed pages: %d.", &vmStat.FileBackedPages},
		"Pages occupied by compressor":   StatFormat{"Pages occupied by compressor: %d.", &vmStat.PagesOccupiedByCompressor},
	}

	for _, line := range *vm_stat_str {
		i := strings.Index(line, ":")
		format, exists := mapStat[line[:i]]
		if exists {
			fmt.Sscanf(line, format.format, format.value)
		}
	}

	stats.Cached = (vmStat.PagesPurgeable + vmStat.FileBackedPages) * vmStat.PageSize
	stats.Used = (vmStat.PagesWiredDown+vmStat.PagesOccupiedByCompressor+vmStat.PagesActive+vmStat.PagesInactive+vmStat.PagesSpeculative)*vmStat.PageSize - stats.Cached
	stats.Total = stats.Used + stats.Cached + (vmStat.PagesFree)*vmStat.PageSize
	stats.Free = (vmStat.PagesFree) * vmStat.PageSize

	stats.UsedPercent = float32(stats.Used) / float32(stats.Total) * 100
	stats.FreePercent = float32(stats.Free) / float32(stats.Total) * 100
	stats.CachedPercent = float32(stats.Cached) / float32(stats.Total) * 100

	return &stats, nil
}

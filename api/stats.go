package api

import (
	"fmt"
	"github.com/awilliams/couchdb-utils/util"
	"strings"
)

type Stat struct {
	Description string
	Current     float64
	Sum         float64
	Mean        float64
	Min         float64
	Max         float64
	Stddev      float64
	Section     string
	SubSection  string
}

func (s Stat) TrimmedDescription() string {
	return strings.Title(strings.TrimPrefix(s.Description, "number of "))
}

func (s Stat) PP(printer util.Printer) {
	var header string
	if s.Description == "" {
		header = fmt.Sprintf("[%s:%s]", s.Section, s.SubSection)
	} else {
		header = fmt.Sprintf(`[%s:%s "%s"]`, s.Section, s.SubSection, s.TrimmedDescription())
	}
	printer.Print(header)
	printer.Print(" Current %v", s.Current)
	printer.Print(" Sum %v", s.Sum)
	printer.Print(" Min/Max %v/%v", s.Min, s.Max)
	printer.Print(" Mean %v", s.Mean)
}

type Stats []Stat

func (s Stats) PP(printer util.Printer) {
	for _, stat := range s {
		stat.PP(printer)
	}
}

type statsJson map[string]map[string]Stat

func (s *statsJson) path(sectionA, sectionB string) string {
	base := "_stats"
	if sectionA == "" && sectionB == "" {
		return base
	} else {
		return base + "/" + sectionA + "/" + sectionB
	}
}

func (c Couchdb) GetStats(sectionA, sectionB string) (Stats, error) {
	statsMap := new(statsJson)
	var stats Stats
	err := c.getJson(statsMap, statsMap.path(sectionA, sectionB))
	for section, _stats := range *statsMap {
		for subSection, stat := range _stats {
			stat.Section = section
			stat.SubSection = subSection
			stats = append(stats, stat)
		}
	}
	return stats, err
}

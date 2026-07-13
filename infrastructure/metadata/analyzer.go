package metadata

import (
	"io"

	"go-api/domain/entity"
)

const maxAnalysisBytes = MaxAnalysisBytes

type AnalysisResult struct {
	Signal    entity.Signal     `json:"signal"`
	Extracted ExtractedMetadata `json:"extracted"`
}

type Analyzer struct{}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(r io.Reader) (AnalysisResult, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxAnalysisBytes))
	if err != nil {
		return AnalysisResult{}, err
	}

	extracted := Extract(data)
	signal := Score(extracted)

	return AnalysisResult{
		Signal:    signal,
		Extracted: extracted,
	}, nil
}

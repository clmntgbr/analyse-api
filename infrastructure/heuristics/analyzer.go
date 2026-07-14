package heuristics

import (
	"io"
)

type Analyzer struct{}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(r io.Reader) (AnalysisResult, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxAnalysisBytes))
	if err != nil {
		return AnalysisResult{}, err
	}

	img, format, err := DecodeImage(data)
	if err != nil {
		return AnalysisResult{}, err
	}

	heuristics := HeuristicsResult{
		NoiseScore:       AnalyzeNoise(img),
		CompressionScore: AnalyzeCompression(img, format, data),
		FrequencyScore:   AnalyzeFrequency(img),
		HistogramScore:   AnalyzeHistogram(img),
		Width:            img.Width,
		Height:           img.Height,
		Format:           format,
	}

	return AnalysisResult{
		Signal:     heuristics.ToSignal(),
		Heuristics: heuristics,
	}, nil
}

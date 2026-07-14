package aggregate

import (
	"testing"

	"go-api/domain/entity"
)

func TestComputeExcludesUnknownMetadataScore(t *testing.T) {
	result := Compute([]entity.Signal{
		{Name: "metadata", Score: -1, Confidence: entity.ConfidenceUnknown},
		{Name: "heuristics", Score: 40, Confidence: entity.ConfidenceLow},
		{Name: "ai_model", Score: 90, Confidence: entity.ConfidenceHigh},
	})

	if result.FinalScore < 80 {
		t.Fatalf("expected metadata -1 to be excluded and ai_model to dominate, got %.2f", result.FinalScore)
	}
}

func TestComputeRenormalizesWeightsWhenMetadataUnknown(t *testing.T) {
	result := Compute([]entity.Signal{
		{Name: "metadata", Score: -1, Confidence: entity.ConfidenceUnknown},
		{Name: "heuristics", Score: 20, Confidence: entity.ConfidenceMedium},
		{Name: "ai_model", Score: 80, Confidence: entity.ConfidenceHigh},
	})

	// heuristics: 0.30 * 0.75 = 0.225, ai_model: 0.50 * 1.0 * 1.5 = 0.75
	expected := (20*0.225 + 80*0.75) / (0.225 + 0.75)
	if mathDiff(result.FinalScore, expected) > 0.5 {
		t.Fatalf("expected renormalized score %.2f, got %.2f", expected, result.FinalScore)
	}
}

func TestComputeStrongAIModelNotDilutedByWeakHeuristics(t *testing.T) {
	result := Compute([]entity.Signal{
		{Name: "metadata", Score: 10, Confidence: entity.ConfidenceLow},
		{Name: "heuristics", Score: 20, Confidence: entity.ConfidenceLow},
		{Name: "ai_model", Score: 95, Confidence: entity.ConfidenceHigh},
	})

	if result.FinalScore < 75 {
		t.Fatalf("expected strong ai_model to keep score high, got %.2f", result.FinalScore)
	}
}

func TestComputeVerdictThresholds(t *testing.T) {
	cases := []struct {
		score    int
		expected string
	}{
		{10, VerdictLikelyReal},
		{50, VerdictUncertain},
		{85, VerdictLikelyAI},
	}

	for _, testCase := range cases {
		result := Compute([]entity.Signal{
			{Name: "ai_model", Score: testCase.score, Confidence: entity.ConfidenceHigh},
		})

		if result.Verdict != testCase.expected {
			t.Fatalf("score %d: expected verdict %s, got %s", testCase.score, testCase.expected, result.Verdict)
		}
	}
}

func TestComputeAllUnknownReturnsInconclusive(t *testing.T) {
	result := Compute([]entity.Signal{
		{Name: "metadata", Score: -1, Confidence: entity.ConfidenceUnknown},
	})

	if result.FinalScore != -1 {
		t.Fatalf("expected final score -1, got %.2f", result.FinalScore)
	}
	if result.Confidence != entity.ConfidenceUnknown {
		t.Fatalf("expected unknown confidence, got %s", result.Confidence)
	}
}

func mathDiff(a, b float64) float64 {
	diff := a - b
	if diff < 0 {
		return -diff
	}

	return diff
}

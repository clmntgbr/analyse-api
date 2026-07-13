package metadata

import (
	"fmt"
	"strings"

	"go-api/domain/entity"
)

const signalName = "metadata"

func Score(extracted ExtractedMetadata) entity.Signal {
	if c2paSignal := scoreC2PA(extracted); c2paSignal != nil {
		return *c2paSignal
	}

	if softwareSignal := scoreSoftwareSignature(extracted); softwareSignal != nil {
		return *softwareSignal
	}

	if promptSignal := scorePromptOrSeed(extracted); promptSignal != nil {
		return *promptSignal
	}

	if extracted.HasAnyMetadata() {
		return neutralMetadataSignal(extracted)
	}

	return entity.Signal{
		Name:       signalName,
		Score:      -1,
		Confidence: entity.ConfidenceUnknown,
		Details:    []string{"no metadata found — inconclusive, not evidence of human origin"},
	}
}

func scoreC2PA(extracted ExtractedMetadata) *entity.Signal {
	if !extracted.C2PA {
		return nil
	}

	confirmed := false
	for _, detail := range extracted.C2PADetails {
		lowerDetail := strings.ToLower(detail)
		for _, indicator := range C2PAAIIndicators {
			if strings.Contains(lowerDetail, indicator) {
				confirmed = true
				break
			}
		}
		if strings.HasPrefix(lowerDetail, "generator:") {
			confirmed = true
		}
	}

	if !confirmed {
		return nil
	}

	details := []string{"C2PA / Content Credentials detected"}
	for _, detail := range extracted.C2PADetails {
		details = append(details, fmt.Sprintf("C2PA indicator: %s", detail))
	}

	return &entity.Signal{
		Name:       signalName,
		Score:      97,
		Confidence: entity.ConfidenceHigh,
		Details:    details,
	}
}

func scoreSoftwareSignature(extracted ExtractedMetadata) *entity.Signal {
	candidates := []string{
		extracted.Software,
		extracted.CreatorTool,
		extracted.Artist,
	}

	joined := strings.ToLower(strings.Join(candidates, " "))
	for _, signature := range SoftwareSignatures {
		if !signatureMatches(joined, signature) {
			continue
		}

		score := (signature.ScoreMin + signature.ScoreMax) / 2
		details := []string{fmt.Sprintf("known AI software signature: %s", signature.Label)}

		if extracted.Software != "" {
			details = append(details, fmt.Sprintf("EXIF Software=%s", extracted.Software))
		}
		if extracted.CreatorTool != "" {
			details = append(details, fmt.Sprintf("XMP CreatorTool=%s", extracted.CreatorTool))
		}

		return &entity.Signal{
			Name:       signalName,
			Score:      score,
			Confidence: entity.ConfidenceHigh,
			Details:    details,
		}
	}

	return nil
}

func scorePromptOrSeed(extracted ExtractedMetadata) *entity.Signal {
	hasPrompt := extracted.Prompt != ""
	hasSeed := extracted.Seed != ""

	if !hasPrompt && !hasSeed {
		for key := range extracted.RawTags {
			lowerKey := strings.ToLower(key)
			for _, indicator := range PromptIndicators {
				if strings.Contains(lowerKey, indicator) {
					hasPrompt = true
					break
				}
			}
		}
	}

	if !hasPrompt && !hasSeed {
		return nil
	}

	details := []string{"generation parameters detected in metadata"}
	if extracted.Prompt != "" {
		details = append(details, fmt.Sprintf("prompt=%s", truncate(extracted.Prompt, 120)))
	}
	if extracted.Seed != "" {
		details = append(details, fmt.Sprintf("seed=%s", extracted.Seed))
	}

	return &entity.Signal{
		Name:       signalName,
		Score:      68,
		Confidence: entity.ConfidenceMedium,
		Details:    details,
	}
}

func neutralMetadataSignal(extracted ExtractedMetadata) entity.Signal {
	details := make([]string, 0, 6)
	score := 18

	if matched, label := matchesNeutralSoftware(
		extracted.Software,
		extracted.CreatorTool,
		extracted.Make,
		extracted.Model,
	); matched {
		score = 12
		details = append(details, fmt.Sprintf("recognized neutral source: %s", label))
	}

	if extracted.Software != "" {
		details = append(details, fmt.Sprintf("Software=%s", extracted.Software))
	}
	if extracted.CreatorTool != "" {
		details = append(details, fmt.Sprintf("CreatorTool=%s", extracted.CreatorTool))
	}
	if extracted.Make != "" {
		details = append(details, fmt.Sprintf("Make=%s", extracted.Make))
	}
	if extracted.Model != "" {
		details = append(details, fmt.Sprintf("Model=%s", extracted.Model))
	}
	if extracted.Artist != "" {
		details = append(details, fmt.Sprintf("Artist=%s", extracted.Artist))
	}

	if len(details) == 0 {
		details = append(details, "metadata present but no recognizable tags")
	} else {
		details = append([]string{"neutral metadata — no AI signature matched"}, details...)
	}

	if !extracted.HasCameraTags() {
		details = append(details, "weak signal: no camera Make/Model tags (not sufficient alone)")
	}

	return entity.Signal{
		Name:       signalName,
		Score:      score,
		Confidence: entity.ConfidenceLow,
		Details:    details,
	}
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}

	return value[:max] + "..."
}

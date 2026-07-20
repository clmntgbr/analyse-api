package centrifugo

import (
	"encoding/json"
	"time"

	"go-api/domain/entity"
)

const (
	EventAnalysisStarted   = "analysis_started"
	EventAnalysisCompleted = "analysis_completed"
)

type MediaEvent struct {
	Type       string          `json:"type"`
	AnalysisID string          `json:"analysisId"`
	MediaID    string          `json:"mediaId"`
	UserID     string          `json:"userId"`
	Status     string          `json:"status"`
	FinalScore float64         `json:"finalScore,omitempty"`
	Confidence string          `json:"confidence,omitempty"`
	Verdict    string          `json:"verdict,omitempty"`
	Signals    []SignalPayload `json:"signals,omitempty"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

type SignalPayload struct {
	Name       string   `json:"name"`
	Score      int      `json:"score"`
	Confidence string   `json:"confidence"`
	Details    []string `json:"details"`
}

func NewAnalysisStartedEvent(media *entity.Media) (MediaEvent, error) {
	if media == nil {
		return MediaEvent{}, ErrInvalidMedia
	}

	return MediaEvent{
		Type:       EventAnalysisStarted,
		AnalysisID: media.AnalysisID.String(),
		MediaID:    media.ID.String(),
		UserID:     media.UserID.String(),
		Status:     string(media.Status),
		UpdatedAt:  media.UpdatedAt,
	}, nil
}

func NewAnalysisCompletedEvent(analysis *entity.Analysis, media *entity.Media, signals []*entity.Signal) (MediaEvent, error) {
	if analysis == nil || media == nil {
		return MediaEvent{}, ErrInvalidMedia
	}

	event := MediaEvent{
		Type:       EventAnalysisCompleted,
		AnalysisID: analysis.ID.String(),
		MediaID:    media.ID.String(),
		UserID:     analysis.UserID.String(),
		Status:     string(media.Status),
		FinalScore: analysis.FinalScore,
		Confidence: string(analysis.AnalysisConfidence),
		Verdict:    analysis.Verdict,
		UpdatedAt:  analysis.UpdatedAt,
		Signals:    make([]SignalPayload, 0, len(signals)),
	}

	for _, signal := range signals {
		if signal == nil {
			continue
		}

		event.Signals = append(event.Signals, SignalPayload{
			Name:       signal.Name,
			Score:      signal.Score,
			Confidence: string(signal.Confidence),
			Details:    signal.Details,
		})
	}

	return event, nil
}

func (e MediaEvent) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

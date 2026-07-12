package response

import (
	"time"

	"github.com/admalfrizi/weekly-wrapped-be/internal/model"
)

type CategoryResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Icon     *string `json:"icon,omitempty"`
	ColorHex *string `json:"color_hex,omitempty"`
}

type ActivityResponse struct {
	ID         string            `json:"id"`
	Category   *CategoryResponse `json:"category,omitempty"`
	Value      float64           `json:"value"`
	Note       *string           `json:"note,omitempty"`
	OccurredAt time.Time         `json:"occurred_at"`
	CreatedAt  time.Time         `json:"created_at"`
}

func MapToActivityResponse(activity model.Activity) ActivityResponse {
	res := ActivityResponse{
		ID:         activity.ID.String(),
		Value:      activity.Value,
		Note:       activity.Note,
		OccurredAt: activity.OccurredAt,
		CreatedAt:  activity.CreatedAt,
	}

	if activity.Category != nil {
		res.Category = &CategoryResponse{
			ID:       activity.Category.ID.String(),
			Name:     activity.Category.Name,
		}
	}

	return res
}

func MapToActivityResponseList(activities []model.Activity) []ActivityResponse {
	var res []ActivityResponse
	for _, a := range activities {
		res = append(res, MapToActivityResponse(a))
	}
	
	if res == nil {
		res = make([]ActivityResponse, 0)
	}
	
	return res
}
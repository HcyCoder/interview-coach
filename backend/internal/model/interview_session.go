package model

import (
	"encoding/json"
	"time"
)

// InterviewSession stores the interview session snapshot for a candidate.
type InterviewSession struct {
	ID               uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID           uint            `json:"userId" gorm:"column:user_id;not null;index"`
	ResumeJSON       json.RawMessage `json:"resumeJson" gorm:"column:resume_json;type:json;not null"`
	ConversationJSON json.RawMessage `json:"conversationJson" gorm:"column:conversation_json;type:json;not null"`
	CurrentQuestion  string          `json:"currentQuestion" gorm:"column:current_question;type:text;not null"`
	JDText           string          `json:"jdText" gorm:"column:jd_text;type:longtext;not null"`
	Status           string          `json:"status" gorm:"column:status;type:varchar(32);not null;index"`
	StartedAt        *time.Time      `json:"startedAt,omitempty" gorm:"column:started_at"`
	FinishedAt       *time.Time      `json:"finishedAt,omitempty" gorm:"column:finished_at"`
	CreatedAt        time.Time       `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time       `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName returns the exact MySQL table name used for interview sessions.
//
// Inputs: none.
// Outputs: the table name expected by GORM migrations and queries.
func (InterviewSession) TableName() string {
	return "interview_sessions"
}

package models

type AliceRequest struct {
	Request struct {
		Command           string `json:"command"`
		OriginalUtterance string `json:"original_utterance"`
	} `json:"request"`
	Session struct {
		New       bool   `json:"new"`
		MessageID int    `json:"message_id"`
		SessionID string `json:"session_id"`
		SkillID   string `json:"skill_id"`
		UserID    string `json:"user_id"`
	} `json:"session"`
	Version string `json:"version"`
}

type Button struct {
	Title   string      `json:"title"`
	Payload interface{} `json:"payload,omitempty"`
	URL     string      `json:"url,omitempty"`
	Hide    bool        `json:"hide"`
}

type AliceResponse struct {
	Response struct {
		Text       string   `json:"text"`
		TTs        string   `json:"tts,omitempty"`
		Buttons    []Button `json:"buttons,omitempty"`
		EndSession bool     `json:"end_session"`
	} `json:"response"`
	Session struct {
		SessionID string `json:"session_id"`
		MessageID int    `json:"message_id"`
		UserID    string `json:"user_id"`
	} `json:"session"`
	Version string `json:"version"`
}

type UserState struct {
	AgreedToPA        bool
	AgreedToMakeSkill bool
}

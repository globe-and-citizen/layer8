package dto

type MessageUpdateDTO struct {
	UpdateID int64       `json:"update_id"`
	Message  *MessageDTO `json:"message,omitempty"`
}

type MessageDTO struct {
	MessageID int64       `json:"message_id"`
	Date      int64       `json:"date"`
	Text      string      `json:"text,omitempty"`
	From      *UserDTO    `json:"from,omitempty"`
	Chat      ChatDTO     `json:"chat"`
	Contact   *ContactDTO `json:"contact,omitempty"`
}

type UserDTO struct {
	ID       int64  `json:"id"`
	IsBot    bool   `json:"is_bot"`
	Username string `json:"username,omitempty"`
}

type ChatDTO struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

type ContactDTO struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	UserID      int64  `json:"user_id,omitempty"`
}

type GetUpdatesResponseDTO struct {
	Ok     bool               `json:"ok"`
	Result []MessageUpdateDTO `json:"result"`
}

type SendMessageRequestDTO struct {
	ChatID                   int64       `json:"chat_id"`
	Text                     string      `json:"text"`
	ParseMode                string      `json:"parse_mode,omitempty"`
	ReplyMarkup              interface{} `json:"reply_markup,omitempty"`
	DisableWebPagePreview    bool        `json:"disable_web_page_preview,omitempty"`
	DisableNotification      bool        `json:"disable_notification,omitempty"`
	ProtectContent           bool        `json:"protect_content,omitempty"`
	ReplyToMessageID         int64       `json:"reply_to_message_id,omitempty"`
	AllowSendingWithoutReply bool        `json:"allow_sending_without_reply,omitempty"`
}

type SendMessageResponseDTO struct {
	Ok     bool       `json:"ok"`
	Result MessageDTO `json:"result"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool               `json:"one_time_keyboard,omitempty"`
	InputFieldPlace string             `json:"input_field_placeholder,omitempty"`
	Selective       bool               `json:"selective,omitempty"`
}

type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact,omitempty"`
	RequestLocation bool   `json:"request_location,omitempty"`
}

type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective,omitempty"`
}

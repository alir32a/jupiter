package tg

const (
	MessageTypeCommand     = "bot_command"
	MessageTypeMention     = "mention"
	MessageTypeHashtag     = "hashtag"
	MessageTypeCashtag     = "cashtag"
	MessageTypeEmail       = "email"
	MessageTypePhoneNumber = "phone_number"
	MessageTypeTextMention = "text_mention"
	MessageTypeTextLink    = "text_link"

	ParseModeMarkdown = "MarkdownV2"
	ParseModeHTML     = "HTML"
)

type GetUpdatesResult struct {
	OK      bool     `json:"ok"`
	Updates []Update `json:"result"`
}

type Update struct {
	UpdateID      int            `json:"update_id"`
	Message       Message        `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

type CallbackQuery struct {
	ID           string  `json:"id"`
	From         From    `json:"from"`
	Message      Message `json:"message"`
	ChatInstance string  `json:"chat_instance"`
	Data         string  `json:"data"`
}

type Message struct {
	MessageID int      `json:"message_id"`
	From      From     `json:"from"`
	Chat      Chat     `json:"chat"`
	Date      int      `json:"date"`
	Text      string   `json:"text"`
	Entities  []Entity `json:"entities"`
	Type      string   `json:"-"`
}

type From struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type Entity struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Type   string `json:"type"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	Url          string `json:"url,omitempty"`
	CallbackData string `json:"callback_data,omitempty"`
}

type InlineKeyboard struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type SendMessageRequest struct {
	ChatID              int             `json:"chat_id"`
	Text                string          `json:"text"`
	MessageThreadID     int             `json:"message_thread_id,omitempty"`
	ParseMode           string          `json:"parse_mode,omitempty"`
	DisableNotification bool            `json:"disable_notification,omitempty"`
	ProtectContent      bool            `json:"protect_content,omitempty"`
	ReplyMarkup         *InlineKeyboard `json:"reply_markup,omitempty"`
}

type SendMessageResponse struct {
	OK       bool      `json:"ok"`
	Messages []Message `json:"result"`
}

type DeleteMessageResponse struct {
	OK     bool `json:"ok,omitempty"`
	Result bool `json:"result,omitempty"`
}

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type commandScope struct {
	Type string `json:"type"`
}

type sendBotCommands struct {
	Commands []BotCommand `json:"commands"`
	Scope    commandScope `json:"scope"`
}

type GetBotCommands struct {
	OK       bool         `json:"ok"`
	Commands []BotCommand `json:"result"`
}

type AnswerCallbackQueryRequest struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text"`
	ShowAlert       bool   `json:"show_alert,omitempty"`
}

type AnswerCallbackQueryResponse struct {
	OK     bool           `json:"ok"`
	Result map[string]any `json:"result,omitempty"`
}

func NewInlineKeyboard(buttons ...InlineKeyboardButton) *InlineKeyboard {
	return &InlineKeyboard{
		InlineKeyboard: append([][]InlineKeyboardButton{}, append([]InlineKeyboardButton{}, buttons...))}
}

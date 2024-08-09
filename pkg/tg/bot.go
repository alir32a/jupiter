package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alir32a/jupiter/internal/errorext"
	"github.com/charmbracelet/log"
	"io"
	"net/http"
)

const DefaultTimeoutInSecond = 5

type Bot struct {
	Token          string
	FailureHandler func(err error)
	lastFetchedID  int
	baseUrl        string
	client         *http.Client
}

func NewBot(token string) *Bot {
	return &Bot{
		Token:   token,
		baseUrl: fmt.Sprintf("https://api.telegram.org/bot%s", token),
		client:  &http.Client{},
	}
}

func (b *Bot) Run(handler func([]Update) error) {
	for {
		updates, err := b.GetUpdates()
		if err != nil {
			if b.FailureHandler != nil {
				b.FailureHandler(err)
			}

			continue
		}

		setMessageTypes(updates)

		if len(updates) > 0 {
			b.lastFetchedID = updates[len(updates)-1].UpdateID
		}

		if err := handler(updates); err != nil {
			continue
		}
	}
}

func (b *Bot) GetUpdates() ([]Update, error) {
	resp, err := b.client.Get(getUpdatesUrl(b.baseUrl, b.lastFetchedID+1, DefaultTimeoutInSecond))
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GetUpdatesResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if !result.OK {
		return nil, fmt.Errorf("got none OK status: %s", data)
	}

	return result.Updates, nil
}

func (b *Bot) SendMessage(req SendMessageRequest) ([]Message, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)

	resp, err := b.client.Post(fmt.Sprintf("%s/sendMessage", b.baseUrl), "application/json", reader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SendMessageResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if !result.OK {
		return nil, fmt.Errorf("got none OK status: %s", data)
	}

	return result.Messages, nil
}

func (b *Bot) DeleteMessage(chatID, msgID int) error {
	resp, err := b.client.Post(
		fmt.Sprintf("%s/deleteMessage?chat_id=%d&message_id=%d", b.baseUrl, chatID, msgID),
		"application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var result DeleteMessageResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	if !result.OK || !result.Result {
		return fmt.Errorf("got none OK status: %s", data)
	}

	return nil
}

func (b *Bot) SetCommands(commands ...BotCommand) error {
	data, err := json.Marshal(sendBotCommands{
		Commands: commands,
		Scope:    commandScope{Type: "all_private_chats"},
	})
	if err != nil {
		return err
	}

	reader := bytes.NewReader(data)

	resp, err := b.client.Post(
		fmt.Sprintf("%s/setMyCommands", b.baseUrl),
		"application/json", reader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result sendBotCommands
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	if len(result.Commands) <= 0 {
		return errorext.New("set command result is empty")
	}

	return nil
}

func (b *Bot) GetCommands() ([]BotCommand, error) {
	resp, err := b.client.Get(fmt.Sprintf("%s/getMyCommands", b.baseUrl))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var commands GetBotCommands
	if err := json.Unmarshal(data, &commands); err != nil {
		return nil, err
	}

	if !commands.OK {
		return nil, fmt.Errorf("got none OK status: %s", data)
	}

	return commands.Commands, nil
}

func (b *Bot) AnswerCallbackQuery(callbackQueryID, text string) error {
	data, err := json.Marshal(AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
		Text:            text,
	})
	if err != nil {
		return err
	}

	reader := bytes.NewReader(data)

	resp, err := b.client.Post(
		fmt.Sprintf("%s/answerCallbackQuery", b.baseUrl),
		"application/json", reader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result AnswerCallbackQueryResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	if !result.OK {
		return errorext.New("got none ok status from answerCallbackQuery")
	}

	return nil
}

func setMessageTypes(updates []Update) {
	for i, update := range updates {
		updates[i].Message.Type = parseMessageType(update.Message.Entities)
	}
}

func parseMessageType(entities []Entity) string {
	for _, entity := range entities {
		return entity.Type
	}

	return ""
}

func getMessages(updates []Update) []Message {
	result := make([]Message, 0, len(updates))

	for _, update := range updates {
		result = append(result, update.Message)
	}

	return result
}

func getUpdatesUrl(baseUrl string, offset, timeout int) string {
	return fmt.Sprintf("%s/getUpdates?offset=%d&timeout=%d", baseUrl, offset, timeout)
}

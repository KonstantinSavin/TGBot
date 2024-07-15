package telegram

import (
	"fmt"
)

const MsgHelp = `Я умею сохранять ваши идеи для досуга. Также я умею вам их напоминать.

Если хотите сохранить идею, пришлите мне ссылку(url) на мероприятие.

Если хотите выбрать идею для досуга нажмите "Выбрать ивент" на виртуальной клавиатуре.

Если хотите использовать меня в совместной беседе, дайте мне права администратора беседы`

const (
	MsgUnknownCommand = "Неизвестная команда 🤔"
	MsgNoSavedPages   = "У вас нет подходящих ивентов 🙈"
	MsgSaved          = "Сохранено! 👌"
	MsgAlreadyExists  = "Этот ивент у вас уже сохранён 🤗"
)

func MsgSaveEvent(url, category string, price, timeDuration int) string {
	return fmt.Sprintf("\nurl: %s \nкатегория: %s \nстоимость: %d \nпродолжительность: %d", url, category, price, timeDuration)
}

func MsgShowEvent(url, category string, price, timeDuration int) string {
	return fmt.Sprintf("\nкатегория: %s \nстоимость: %d \nпродолжительность: %d", category, price, timeDuration)
}

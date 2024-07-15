package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"projects/DAB/internal/config"
	"projects/DAB/internal/page"
	pageDB "projects/DAB/internal/page/db/postgresql"
	tg "projects/DAB/internal/telegram"
	"projects/DAB/pkg/logging"
	"projects/DAB/pkg/storage/postgres"
	"strconv"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var logger = logging.GetLogger()
var p page.Page

func main() {
	cfg := config.GetConfig()
	logger.Info("конфиг получен")

	postgreSQLClient, err := postgres.New(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatal("%w", err)
	}

	storage := pageDB.New(postgreSQLClient, logger)
	logger.Info("база данных подключена")

	if err := storage.Init(context.TODO()); err != nil {
		logger.Fatal("не можем проинициализировать бд из-за ошибки: %w", err)
	}

	botToken := mustToken()
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		logger.Fatal("не удалось подключиться к телеграм из-за %w", err)
		os.Exit(1)
	}

	logger.Info("сервис запущен")

	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		logger.Error("не удалось получить update из-за ошибки: %w", err)
	}

	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		logger.Error("не удалось создать бот хендлер из-за ошибки: %w", err)
	}

	defer bh.Stop()
	defer bot.StopLongPolling()

	//		КОМАНДА START

	bh.Handle(func(bot *telego.Bot, update telego.Update) {

		//	ОБРАБОТКА КОМАНДЫ /start

		logger.Infof("получили новую команду '%s' от '%s'", update.Message.Text, update.Message.From.Username)

		chatID := tu.ID(update.Message.Chat.ID)

		keyBoard := tu.Keyboard(
			tu.KeyboardRow(tu.KeyboardButton("Выбрать ивент")),
		)

		message := tu.Message(chatID, "Привет, "+update.Message.From.FirstName+"!\n\n"+tg.MsgHelp).WithReplyMarkup(keyBoard)

		_, err := bot.SendMessage(message)
		if err != nil {
			logger.Error(err)
		}
	}, th.CommandEqual("start"))

	//		КОМАНДА HELP

	bh.Handle(func(bot *telego.Bot, update telego.Update) {

		//	ОБРАБОТКА ВЫЗОВА ПОМОЩИ

		logger.Infof("получили новую команду '%s' от '%s'", update.Message.Text, update.Message.From.Username)

		chatID := tu.ID(update.Message.Chat.ID)

		keyBoard := tu.Keyboard(
			tu.KeyboardRow(tu.KeyboardButton("Выбрать ивент")),
		)

		message := tu.Message(chatID, tg.MsgHelp).WithReplyMarkup(keyBoard)

		_, err = bot.SendMessage(message)
		if err != nil {
			logger.Error(err)
		}
	}, th.CommandEqual("help"))

	//		ЗАПОЛНЕНИЕ CATEGORY

	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {

		logger.Infof("получили категорию ивента '%s' от '%s'", query.Data, query.From.Username)

		category := query.Data
		p.Category = category

		logger.Infof("записали category '%s'", p.Category)

		_, _ = bot.SendMessage(tu.Message(
			tu.ID(query.Message.GetChat().ID),
			"Выберите, ценовую категорию.",
		).WithReplyMarkup(tu.InlineKeyboard(
			tu.InlineKeyboardRow(tu.InlineKeyboardButton("Бесплатно").WithCallbackData("0")),
			tu.InlineKeyboardRow(tu.InlineKeyboardButton("до 1000р").WithCallbackData("1000")),
			tu.InlineKeyboardRow(tu.InlineKeyboardButton("1000 - 3000р").WithCallbackData("3000")),
			tu.InlineKeyboardRow(tu.InlineKeyboardButton("3000 - 6000р").WithCallbackData("6000")),
			tu.InlineKeyboardRow(tu.InlineKeyboardButton("более 6000р").WithCallbackData("10000")),
		),
		))

		_ = bot.AnswerCallbackQuery(tu.CallbackQuery(query.ID).WithText("Категория выбрана"))
	}, th.AnyCallbackQueryWithMessage(), th.Or(th.CallbackDataEqual("sport"), th.CallbackDataEqual("cafe"), th.CallbackDataEqual("fest"), th.CallbackDataEqual("nature")))

	//		ЗАПОЛНЕНИЕ PRICE

	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {

		logger.Infof("получили ценовую категорию ивента '%s' от '%s'", query.Data, query.From.Username)

		price, err := strconv.Atoi(query.Data)
		if err != nil {
			logger.Error(err)
		}
		p.Price = price

		logger.Infof("записали price '%d'", p.Price)

		_, _ = bot.SendMessage(tu.Message(
			tu.ID(query.Message.GetChat().ID),
			"Выберите, временной промежуток.",
		).WithReplyMarkup(tu.InlineKeyboard(
			tu.InlineKeyboardRow(tu.InlineKeyboardButton("до 1ч").WithCallbackData("1")),
			tu.InlineKeyboardRow(tu.InlineKeyboardButton("1 - 3ч").WithCallbackData("3")),
			tu.InlineKeyboardRow(tu.InlineKeyboardButton("3 - 6ч").WithCallbackData("6")),
			tu.InlineKeyboardRow(tu.InlineKeyboardButton("более 6ч").WithCallbackData("12")),
		),
		))

		_ = bot.AnswerCallbackQuery(tu.CallbackQuery(query.ID).WithText("Ценовая категория выбрана"))
	}, th.AnyCallbackQueryWithMessage(), th.Or(th.CallbackDataEqual("0"), th.CallbackDataEqual("1000"), th.CallbackDataEqual("3000"), th.CallbackDataEqual("6000"), th.CallbackDataEqual("10000")))

	//		ЗАПОЛНЕНИЕ TIMEDURATION

	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {

		logger.Infof("получили продолжительность ивента '%s' от '%s'", query.Data, query.From.Username)

		timeDuration, err := strconv.Atoi(query.Data)
		if err != nil {
			logger.Error(err)
		}
		p.TimeDuration = timeDuration

		logger.Infof("записали timeDuration '%d'", p.TimeDuration)

		if p.URL != "" {
			_, _ = bot.SendMessage(tu.Message(
				tu.ID(query.Message.GetChat().ID),
				fmt.Sprintf("Готово! \n"+tg.MsgSaveEvent(p.URL, p.Category, p.Price, p.TimeDuration)),
			).WithReplyMarkup(tu.InlineKeyboard(
				tu.InlineKeyboardRow(tu.InlineKeyboardButton("Сохранить ивент").WithCallbackData("save")),
			),
			))
		} else {
			_, _ = bot.SendMessage(tu.Message(
				tu.ID(query.Message.GetChat().ID),
				fmt.Sprintf("Готово! \n"+tg.MsgShowEvent(p.URL, p.Category, p.Price, p.TimeDuration)),
			).WithReplyMarkup(tu.InlineKeyboard(
				tu.InlineKeyboardRow(tu.InlineKeyboardButton("Показать ивенты").WithCallbackData("show")),
			),
			))
		}

		_ = bot.AnswerCallbackQuery(tu.CallbackQuery(query.ID).WithText("Готово"))
	}, th.AnyCallbackQueryWithMessage(), th.Or(th.CallbackDataEqual("1"), th.CallbackDataEqual("3"), th.CallbackDataEqual("6"), th.CallbackDataEqual("12")))

	//		ОБРАБОТКА СОХРАНЕНИЯ ИВЕНТА

	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {

		logger.Infof("получили запрос на сохранение ивента от '%s'", query.From.Username)

		err = storage.Save(context.TODO(), &p)
		if err != nil {
			logger.Error(err)
		}

		p = page.Page{}

		_, _ = bot.SendMessage(tu.Message(
			tu.ID(query.Message.GetChat().ID),
			tg.MsgSaved,
		).WithReplyMarkup(tu.Keyboard(
			tu.KeyboardRow(tu.KeyboardButton("Выбрать ивент")),
		)))

		_ = bot.AnswerCallbackQuery(tu.CallbackQuery(query.ID).WithText("Готово"))
	}, th.AnyCallbackQueryWithMessage(), th.CallbackDataEqual("save"))

	//		ОБРАБОТКА ПОКАЗА ИВЕНТОВ

	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {

		logger.Infof("получили запрос на показ ивентов от '%s'", query.From.Username)

		pages, err := storage.Show(context.TODO(), &p)
		if err != nil {
			logger.Error(err)
		}

		if len(pages) == 0 {
			_, _ = bot.SendMessage(tu.Message(
				tu.ID(query.Message.GetChat().ID),
				tg.MsgNoSavedPages,
			).WithReplyMarkup(tu.Keyboard(
				tu.KeyboardRow(tu.KeyboardButton("Выбрать ивент")),
			)))
		}

		for _, v := range pages {
			_, _ = bot.SendMessage(tu.Message(
				tu.ID(query.Message.GetChat().ID),
				fmt.Sprintf("%s \nПрислал: %s\nСохранено: %s", v.URL, v.UserName, v.CreateTime.Format("2006.01.02")),
			).WithReplyMarkup(tu.Keyboard(
				tu.KeyboardRow(tu.KeyboardButton("Выбрать ивент")),
			)))
		}

		p = page.Page{}

		_ = bot.AnswerCallbackQuery(tu.CallbackQuery(query.ID).WithText("Готово"))
	}, th.AnyCallbackQueryWithMessage(), th.CallbackDataEqual("show"))

	//		ОБРАБОТКА СООБЩЕНИЙ

	bh.Handle(func(bot *telego.Bot, update telego.Update) {

		logger.Infof("получили новую команду '%s' от '%s", update.Message.Text, update.Message.From.Username)

		chatID := tu.ID(update.Message.Chat.ID)

		if tg.IsURL(update) {

			//		ОБРАБОТКА ДОБАВЛЕНИЯ ИВЕНТА

			p.URL = update.Message.Text
			p.UserName = update.Message.From.Username

			logger.Infof("записали username '%s' и url '%s'", p.UserName, p.URL)

			keyBoard := tu.InlineKeyboard(
				tu.InlineKeyboardRow(tu.InlineKeyboardButton("Спорт").WithCallbackData("sport")),
				tu.InlineKeyboardRow(tu.InlineKeyboardButton("Кафе").WithCallbackData("cafe")),
				tu.InlineKeyboardRow(tu.InlineKeyboardButton("Фестиваль").WithCallbackData("fest")),
				tu.InlineKeyboardRow(tu.InlineKeyboardButton("Природа").WithCallbackData("nature")),
			)

			message := tu.Message(chatID, "Выберите категорию").WithReplyMarkup(keyBoard)
			_, err = bot.SendMessage(message)
			if err != nil {
				logger.Error(err)
			}

		} else if update.Message.Text == "Выбрать ивент" {

			//		ОБРАБОТКА ПОДБОРА ИВЕНТА

			logger.Info("начали поиск ивента")

			keyBoard := tu.InlineKeyboard(
				tu.InlineKeyboardRow(tu.InlineKeyboardButton("Спорт").WithCallbackData("sport")),
				tu.InlineKeyboardRow(tu.InlineKeyboardButton("Кафе").WithCallbackData("cafe")),
				tu.InlineKeyboardRow(tu.InlineKeyboardButton("Фестиваль").WithCallbackData("fest")),
				tu.InlineKeyboardRow(tu.InlineKeyboardButton("Природа").WithCallbackData("nature")),
			)

			message := tu.Message(chatID, "Выберите категорию").WithReplyMarkup(keyBoard)
			_, err = bot.SendMessage(message)
			if err != nil {
				logger.Error(err)
			}
		} else {

			//	ОБРАБОТКА ЛЮБОГО СООБЩЕНИЯ

			keyBoard := tu.Keyboard(
				tu.KeyboardRow(tu.KeyboardButton("Выбрать ивент")),
			)

			message := tu.Message(chatID, "Непонятная команда").WithReplyMarkup(keyBoard)
			_, err = bot.SendMessage(message)
			if err != nil {
				logger.Error(err)
			}
		}

	}, th.AnyMessage())

	bh.Start()
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"токен для доступа к тгБоту",
	)

	flag.Parse()

	if *token == "" {
		logger.Fatal("токен не получен")
	}

	logger.Info("токен получен")
	return *token
}

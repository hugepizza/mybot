package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mybot/gmap"
	"os"
	"strconv"
)

const template = `
Route %s

From %s
[%s]

To %s
[%s]

Check the route here 👇
[%s]

`

func main() {
	var (
		routes []*Route
		err    error
	)
	db, err := sql.Open("sqlite3", "mango.sqlite")
	if err != nil {
		panic("failed to load sqlite3 database")
	}
	defer db.Close()
	if routes, err = QueryRoutes(db); err != nil {
		panic("failed to load routes data")
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		panic("BOT_TOKEN is not set")
	}
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}
	bot.GetMyCommands()
	uf := tgbotapi.NewUpdate(0)
	uf.Timeout = 100

	routesKeyboard := assembleRoutesKeyboard(routes)
	pointsKeyboard := assemblePointRoutesKeyboard(routes)
	updateChan := bot.GetUpdatesChan(uf)

	for update := range updateChan {
		if update.Message != nil && update.Message.IsCommand() {
			reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, "+update.Message.From.FirstName)
			switch update.Message.Command() {
			case "mango":
				reply.ReplyMarkup = &routesKeyboard
				bot.Send(reply)
			case "driedmango":
				reply.ReplyMarkup = &pointsKeyboard
				bot.Send(reply)
			case "Ai":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "love u"))
			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Hello %s, if u have useful information to complete this database, please send it to me, my email is wangleilei950325@gmail.com", update.Message.From.FirstName)))
			}
			continue
		}
		if update.CallbackQuery != nil {
			switch CallbackData(update.CallbackData()).GetCallbackType() {
			case CallbackTypeRouteKey: // return text message
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
				if _, err := bot.Request(callback); err != nil {
					continue
				}
				indexes := CallbackData(update.CallbackData()).GetCallbackIndex()
				if len(indexes) != 1 {
					continue
				}
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, routes[indexes[0]].Format())
				if !routes[indexes[0]].Validated {
					msg.ReplyMarkup = assembleVerifyKeyboard(routes[indexes[0]])
				}
				bot.Send(msg)
			case CallbackTypePointKey: // return a route list keyboard
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
				if _, err := bot.Request(callback); err != nil {
					continue
				}
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Passing Routes")
				indexes := CallbackData(update.CallbackData()).GetCallbackIndex()
				msg.ReplyMarkup = assemblePassingRouteKeyboard(routes, indexes...)
				bot.Send(msg)
			case CallbackTypeVerify:
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
				if _, err := bot.Request(callback); err != nil {
					continue
				}
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Passing Routes")
				indexes := CallbackData(update.CallbackData()).GetCallbackIndex()
				msg.Text = strconv.Itoa(indexes[0])
				bot.Send(msg)
			}
			continue
		}
		reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Use commands to get routes info")
		reply.ReplyToMessageID = update.Message.MessageID
		bot.Send(reply)
	}
}

type Route struct {
	gmap.Route `json:",inline"`

	Index     int    `json:"index"`
	Validated bool   `json:"validated"`
	Note      string `json:"note"`
}

func (r *Route) Format() string {
	if len(r.PassingPoints) < 2 {
		return ""
	}

	var departure, destination = r.PassingPoints[0], r.PassingPoints[len(r.PassingPoints)-1]
	text := fmt.Sprintf(template,
		r.Name, departure.Name, departure.ShareUrl(), destination.Name, destination.ShareUrl(), r.ShareUrl())
	if !r.Validated {
		text += "⚠️Warning: This route is not validated by author yet, plz confirm before your trip\n"
	}
	if r.Note != "" {
		text += fmt.Sprintf("️⚛ Note: %s", r.Note)
	}
	return text
}

func assembleRoutesKeyboard(routes []*Route) tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup()
	for _, route := range routes {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				route.Name, string(NewCallbackData(CallbackTypeRouteKey, route.Index)),
			))
		markup.InlineKeyboard = append(markup.InlineKeyboard, row)
	}
	return markup
}

func assemblePointRoutesKeyboard(routes []*Route) tgbotapi.InlineKeyboardMarkup {
	pointWithRouteIndexes := make(map[string][]int)
	for _, route := range routes {
		for _, point := range route.PassingPoints {
			if _, ok := pointWithRouteIndexes[point.Name]; !ok {
				pointWithRouteIndexes[point.Name] = []int{route.Index}
			} else {
				pointWithRouteIndexes[point.Name] = append(pointWithRouteIndexes[point.Name], route.Index)
			}
		}
	}
	markup := tgbotapi.NewInlineKeyboardMarkup()
	for pointName, routeIndexes := range pointWithRouteIndexes {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				pointName, string(NewCallbackData(CallbackTypePointKey, routeIndexes...)),
			))
		markup.InlineKeyboard = append(markup.InlineKeyboard, row)
	}
	return markup
}

func assemblePassingRouteKeyboard(routes []*Route, routeIndexes ...int) tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup()
	var selectedRoutes []*Route
	for _, index := range routeIndexes {
		for _, route := range routes {
			if route.Index == index {
				selectedRoutes = append(selectedRoutes, route)
			}
		}
	}
	for _, route := range selectedRoutes {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				route.Name, string(NewCallbackData(CallbackTypeRouteKey, route.Index)),
			))
		markup.InlineKeyboard = append(markup.InlineKeyboard, row)
	}
	return markup
}

func assembleVerifyKeyboard(route *Route) tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup()
	row := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"verify "+route.Name, string(NewCallbackData(CallbackTypeVerify, route.Index)),
		))
	markup.InlineKeyboard = append(markup.InlineKeyboard, row)
	return markup
}

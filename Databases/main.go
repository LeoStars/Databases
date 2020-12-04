package main

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)
var keyboardChooseKorz = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Да"),
		tgbotapi.NewKeyboardButton("Нет"),
		tgbotapi.NewKeyboardButton("Очистить корзину"),
	),
)

var keyboardChooseKorz2 = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Да"),
		tgbotapi.NewKeyboardButton("Нет"),
	),
)

var keyboardAdmin = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Добавить товар"),
		tgbotapi.NewKeyboardButton("Найти товар"),
		tgbotapi.NewKeyboardButton("Добавить покупателя"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Найти покупателя"),
		tgbotapi.NewKeyboardButton("Скидка на товар"),
		tgbotapi.NewKeyboardButton("Скидка на отдел"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Корзина"),
		tgbotapi.NewKeyboardButton("Список товаров"),
	),
)



var keyboardUser = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Найти покупателя"),
		tgbotapi.NewKeyboardButton("Найти товар"),
		tgbotapi.NewKeyboardButton("Добавить покупателя"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Корзина"),
		tgbotapi.NewKeyboardButton("Список товаров"),
	),
)

var bot, _ = tgbotapi.NewBotAPI("1439475414:AAGDsm1g5uwMtCtT7TrZxBvJFPHx2dfRcCA")
var u = tgbotapi.NewUpdate(0)

// States:
// 0 - default
// 1 - login
// 2 - password
// 3 - keyboard
//

type User struct {
	state int64
	stateTovar int64
	auth [2]string
	tovar []string
	client []string
	sold []Korzina
	t *[6]string
	price float64
}

type Korzina struct {
	id string
	name string
	otdel string
	marka string
	srok string
	price string

}


type server struct {
	data     *database
}

var s = &server{}

func main() {
	s.setupDB()

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)

	//s0tate := -1
	//stateTovar := -1
	//var tovar = []string{}
	//var auth = [2]string{}
	userMap := make(map[int64]*User)
	admin := 0

	for update := range updates {
		if userMap[update.Message.Chat.ID] == nil {
			userMap[update.Message.Chat.ID] = new(User)
		}
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствуем в ShoppingBot! Пожалуйста, авторизуйтесь")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			bot.Send(msg)


			userMap[update.Message.Chat.ID].price = 0
			userMap[update.Message.Chat.ID].state = 8
			userMap[update.Message.Chat.ID].stateTovar = 0

		case "/close":
			userMap[update.Message.Chat.ID].stateTovar = 0
			userMap[update.Message.Chat.ID].state = 0
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сессия завершена")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			bot.Send(msg)
			userMap[update.Message.Chat.ID].price = 0

			continue
		case "Добавить товар":
			userMap[update.Message.Chat.ID].price = 0
			userMap[update.Message.Chat.ID].state = 0
			userMap[update.Message.Chat.ID].stateTovar = 0
			if userMap[update.Message.Chat.ID].stateTovar == 0 {
				command := "Введите характеристики товара. По характеристике на каждой строке. Характеристики" +
					" в следующем порядке:\n" + "1. Название товара.\n" +
					"2. Цена.\n" + "3. Номер отдела\n" + "4. Марка товара\n" +
					"5. Срок годности товара"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, command)
				bot.Send(msg)
				userMap[update.Message.Chat.ID].stateTovar = 1

				continue
			}
		case "Найти товар":
			userMap[update.Message.Chat.ID].price = 0
			userMap[update.Message.Chat.ID].state = 0
			userMap[update.Message.Chat.ID].stateTovar = 0
			if userMap[update.Message.Chat.ID].stateTovar == 0 {
				command := "Введите ID товара:"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, command)
				bot.Send(msg)
				userMap[update.Message.Chat.ID].stateTovar = 11
				continue
			}
		case "Найти покупателя":
			userMap[update.Message.Chat.ID].price = 0
			userMap[update.Message.Chat.ID].state = 0
			userMap[update.Message.Chat.ID].stateTovar = 0
			if userMap[update.Message.Chat.ID].stateTovar == 0 {
				command := "Введите ID покупателя:"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, command)
				bot.Send(msg)
				userMap[update.Message.Chat.ID].stateTovar = 21
				continue
			}
		case "Список товаров":
			userMap[update.Message.Chat.ID].price = 0
			userMap[update.Message.Chat.ID].state = 0
			userMap[update.Message.Chat.ID].stateTovar = 0
			every := s.data.getAll()
			e := ""
			for j, i := range every {
				e += strconv.Itoa(j) + ". " + "Название товара: " + i.name + "\n" + "Марка товара: " + i.marka + "\n" +
					"Цена: " + i.price + "\n" + "Срок годности: " + i.srok + "\n" +	"Номер отдела: " + i.otdel + "\n" +
					"ID товара: " + i.id + "\n" + "\n"
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, e)
			bot.Send(msg)
			continue
		case "Скидка на товар":
		case "Скидка на отдел":
		case "Добавить покупателя":
			userMap[update.Message.Chat.ID].price = 0
			userMap[update.Message.Chat.ID].state = 0
			userMap[update.Message.Chat.ID].stateTovar = 0
			if userMap[update.Message.Chat.ID].stateTovar == 0 {
				command := "Введите описание покупателя. По описанию на каждой строке. Описание" +
					" в следующем порядке:\n" + "1. ФИО.\n" +
					"2. Номер телефона.\n" + "3. Адрес\n"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, command)
				bot.Send(msg)
				userMap[update.Message.Chat.ID].stateTovar = 31

				continue
			}
		case "Корзина":
			userMap[update.Message.Chat.ID].price = 0
			userMap[update.Message.Chat.ID].state = 0
			userMap[update.Message.Chat.ID].stateTovar = 0
			if len(userMap[update.Message.Chat.ID].sold) == 0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Корзина пуста")
				bot.Send(msg)
				continue
			}
			korz := ""
			for j, i := range userMap[update.Message.Chat.ID].sold {
				korz += strconv.Itoa(j) + ". " + "Название товара: " + i.name + "\n" + "Марка товара: " + i.marka + "\n" +
					"Цена: " + i.price + "\n" + "Срок годности: " + i.srok + "\n" +	"Номер отдела: " + i.otdel + "\n" +
					"ID товара: " + i.id + "\n"
			}
			korz += "Общая цена: " + fmt.Sprintf("%f", userMap[update.Message.Chat.ID].price)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, korz)
			bot.Send(msg)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Продать товар?")
			msg.ReplyMarkup = keyboardChooseKorz
			bot.Send(msg)
			userMap[update.Message.Chat.ID].stateTovar = 101
			continue

		}

		if userMap[update.Message.Chat.ID].stateTovar == 101 {
			ans := update.Message.Text
			if ans == "Да" {
				s.data.remove(userMap[update.Message.Chat.ID].sold)
				userMap[update.Message.Chat.ID].sold = userMap[update.Message.Chat.ID].sold[:0]
				userMap[update.Message.Chat.ID].stateTovar = 0
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Товар успешно продан")
				if admin == 1 {
					msg.ReplyMarkup = keyboardAdmin
				} else {
					msg.ReplyMarkup = keyboardUser
				}
				bot.Send(msg)
				continue
			} else if ans == "Очистить корзину" {
				userMap[update.Message.Chat.ID].sold = userMap[update.Message.Chat.ID].sold[:0]
				userMap[update.Message.Chat.ID].stateTovar = 0
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Корзина успешно очищена")
				if admin == 1 {
					msg.ReplyMarkup = keyboardAdmin
				} else {
					msg.ReplyMarkup = keyboardUser
				}
				bot.Send(msg)
				continue
			} else {
				userMap[update.Message.Chat.ID].stateTovar = 0
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Тогда продолжаем выбор")
				if admin == 1 {
					msg.ReplyMarkup = keyboardAdmin
				} else {
					msg.ReplyMarkup = keyboardUser
				}
				bot.Send(msg)
				continue
			}
		}

		if userMap[update.Message.Chat.ID].stateTovar == 1 {
			userMap[update.Message.Chat.ID].tovar = strings.Split(update.Message.Text, "\n")
			//command := "Отправьте изображение товара:"
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, command)
			//bot.Send(msg)
			s.data.newAddress(userMap[update.Message.Chat.ID].tovar)
			fmt.Println(userMap[update.Message.Chat.ID].tovar)
			userMap[update.Message.Chat.ID].stateTovar = 0
			continue
		} /*else if userMap[update.Message.Chat.ID].stateTovar == 2 {
			photo := update.Message.Photo
			fmt.Println(photo)
			msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "kek.png")
			bot.Send(msg)
			userMap[update.Message.Chat.ID].stateTovar = 0
			continue
		}*/



		if userMap[update.Message.Chat.ID].stateTovar == 11 {
			userMap[update.Message.Chat.ID].t = s.data.getNeeded(update.Message.Text)
			if userMap[update.Message.Chat.ID].t[0] == "ID не существует. Попробуйте снова" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ID не существует. Попробуйте снова")
				bot.Send(msg)
				userMap[update.Message.Chat.ID].stateTovar = 11
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Название товара: " + userMap[update.Message.Chat.ID].t[1] + "\n" + "Марка товара: " + userMap[update.Message.Chat.ID].t[3] + "\n" +
				"Цена: " + userMap[update.Message.Chat.ID].t[5] + "\n" + "Срок годности: " + userMap[update.Message.Chat.ID].t[4] + "\n" +	"Номер отдела: " + userMap[update.Message.Chat.ID].t[2] + "\n" +
				"ID товара: " + userMap[update.Message.Chat.ID].t[0])
			bot.Send(msg)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Добавить товар в корзину?")
			msg.ReplyMarkup = keyboardChooseKorz2
			bot.Send(msg)
			userMap[update.Message.Chat.ID].stateTovar = 12
			continue
		} else if userMap[update.Message.Chat.ID].stateTovar == 12 {
			ans := update.Message.Text
			if ans == "Да" {
				a := Korzina{
					id: userMap[update.Message.Chat.ID].t[0],
					marka: userMap[update.Message.Chat.ID].t[1],
					otdel: userMap[update.Message.Chat.ID].t[2],
					name: userMap[update.Message.Chat.ID].t[3],
					srok: userMap[update.Message.Chat.ID].t[4],
					price: userMap[update.Message.Chat.ID].t[5],
				}
				userMap[update.Message.Chat.ID].sold = append(userMap[update.Message.Chat.ID].sold, a)
				fmt.Println(userMap[update.Message.Chat.ID].sold)
				k, _:=strconv.ParseFloat(userMap[update.Message.Chat.ID].t[5], 64)
				userMap[update.Message.Chat.ID].price += k
				userMap[update.Message.Chat.ID].stateTovar = 0
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Товар добавлен в корзину")
				if admin == 1 {
					msg.ReplyMarkup = keyboardAdmin
				} else {
					msg.ReplyMarkup = keyboardUser
				}
				bot.Send(msg)
				continue
			} else {
				userMap[update.Message.Chat.ID].stateTovar = 0
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Тогда продолжаем выбор")
				if admin == 1 {
					msg.ReplyMarkup = keyboardAdmin
				} else {
					msg.ReplyMarkup = keyboardUser
				}
				bot.Send(msg)
				continue
			}
		}


		if userMap[update.Message.Chat.ID].stateTovar == 21 {
			client := s.data.getNeededUser(update.Message.Text)
			if client[0] == "ID не существует. Попробуйте снова" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ID не существует. Попробуйте снова")
				bot.Send(msg)
				userMap[update.Message.Chat.ID].stateTovar = 21
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, client[0] +
				"\n" + client[1] + "\n" + client[2] + "\n" + client[3])
			bot.Send(msg)
			userMap[update.Message.Chat.ID].stateTovar = 0
			continue
		}

		if userMap[update.Message.Chat.ID].stateTovar == 31 {
			userMap[update.Message.Chat.ID].client = strings.Split(update.Message.Text, "\n")
			s.data.newAddressUser(userMap[update.Message.Chat.ID].client)
			fmt.Println(userMap[update.Message.Chat.ID].client)
			userMap[update.Message.Chat.ID].stateTovar = 0
			continue
		}


		if userMap[update.Message.Chat.ID].state == 8 {
			command := "Введите логин:"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, command)
			bot.Send(msg)
			userMap[update.Message.Chat.ID].state = 1
			continue
		} else if userMap[update.Message.Chat.ID].state == 1 {
			userMap[update.Message.Chat.ID].auth[0] = update.Message.Text
			command := "Введите пароль:"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, command)
			bot.Send(msg)
			userMap[update.Message.Chat.ID].state = 2
			continue
		} else if userMap[update.Message.Chat.ID].state == 2 {
			userMap[update.Message.Chat.ID].auth[1] = update.Message.Text
			resp := login(userMap[update.Message.Chat.ID].auth[0], userMap[update.Message.Chat.ID].auth[1])
			if resp == "admin" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите пункт меню")
				msg.ReplyMarkup = keyboardAdmin
				admin = 1
				bot.Send(msg)
				userMap[update.Message.Chat.ID].state = -1
				continue
			} else if resp == "user" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите пункт меню")
				msg.ReplyMarkup = keyboardUser
				admin = 0
				bot.Send(msg)
				userMap[update.Message.Chat.ID].state = -1
				continue
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Попробуйте ещё раз! Введите логин:")
				bot.Send(msg)
				userMap[update.Message.Chat.ID].state = 1
				continue
			}
		}
	}
}
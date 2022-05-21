package vars

import "time"

var (
	ErrorNoUser  = "Пользователь не найден"
	ErrorDefault = "Ошибка отправки сообщения"

	CareDisabled = "Напонинание об уходе отключено"
	CareEnabled  = "Напоминание об уходе включено"

	WelcomeMessage            = "Добро пожаловать в «Косметолог-Бот»! \n\n❗️Для продолжения работы авторизуйтесь по номеру телефона или имени пользователя Telegram.\n\n <i>Имя пользователя необходимо предварительно сообщить специалисту</i>"
	GreetingsMessage          = "Вы успешно авторизовались!"
	DescriptionMessage        = "Вы можете воспользоваться следующими командами:"
	DailyMessage              = "Здесь Вы можете настроить время оповещения для дневных и вечерних процедур"
	DailySetTimeMessage       = "Пожалуйста, выберите время напоминания для "
	TimeChangedSuccessMessage = "Вы успешно изменили время "
	NoAppointmentMessage      = "К сожалению, у вас нет ближайших записей."

	InlineButtonAppointment = "Запись"
	InlineButtonCare        = "Уход "
	InlineButtonDailyCare   = "Ежедневный уход"
	InlineButtonOk          = "Ок"
	InlineButtonBack        = "Назад"

	InlineButtonMorning = "Утренняя"
	InlineButtonEvening = "Вечерняя"

	InlineButtonInc = "↑"
	InlineButtonDec = "↓"

	KeyboardButtonMobilePhone = "Номер телефона"
	KeyboardButtonUsername    = "Имя пользователя"

	TimeToSleep time.Duration = 5

	MorningGreetings1 = "🌞 Доброе утро, "
	MorningGreetings2 = "Здравствуйте, "
	MorningGreetings3 = "👋 Привет, "
	MorningGreetings4 = "👋 Hello "
	MorningGreetings5 = "Good morning "

	EveningGreetings1 = "🌜Добрый вечер, "
	EveningGreetings2 = "Здравствуйте, "
	EveningGreetings3 = "Good evening "
	EveningGreetings4 = "👋 Привет, "
	EveningGreetings5 = "👋 Hello, "
)

var MorningGreetings = [...]string{
	MorningGreetings1,
	MorningGreetings2,
	MorningGreetings3,
	MorningGreetings4,
	MorningGreetings5,
}

var EveningGreetings = [...]string{
	EveningGreetings1,
	EveningGreetings2,
	EveningGreetings3,
	EveningGreetings4,
	EveningGreetings5,
}

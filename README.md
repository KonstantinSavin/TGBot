# TGBot

Бот разработан для подбора оптимального досуга по нескольким критериям. На данный момент реализованы следуищие критерии: тематика (спорт, ресторан и тд), стоимость, продолжительность. При желании можно относительно легко добавить иные критерии. Идеи для досуга пользователь записывает сам исходя из своих предпочтений.

В коде реализована связь с базой данных Postgresql через Docker compose (для работы с postgres клиентом использован фреймворк github.com/jackc/pgconn), логирование (github.com/sirupsen/logrus), подключение к telegram API (github.com/mymmrac/telego), работу с конфигурацией (github.com/ilyakaznacheev/cleanenv).


Сервис, предоставляющий API по созданию сокращённых ссылок.

Ссылка должна быть:
— Уникальной; на один оригинальный URL должна ссылаться только одна сокращенная ссылка;
— Длиной 10 символов;
— Из символов латинского алфавита в нижнем и верхнем регистре, цифр и символа _ (подчеркивание).

Сервис должен быть написан на Go и принимать следующие запросы по http:
1. Метод Post, который будет сохранять оригинальный URL в базе и возвращать сокращённый.
2. Метод Get, который будет принимать сокращённый URL и возвращать оригинальный.
3. В качестве хранилища ожидается in-memory решение и PostgreSQL. Какое хранилище использовать, указывается параметром при запуске сервиса.

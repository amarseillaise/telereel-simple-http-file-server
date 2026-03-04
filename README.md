# Simple HTTP File Server

API-сервер для управления видео-контентом (загрузка Instagram Reels через yt-dlp).

## Требования

- Go 1.21+
- yt-dlp
- ffmpeg
- cookies.txt (для Instagram)

## Быстрый старт

```bash
cp .env.example .env
# отредактировать .env

go run cmd/server/main.go
```

## Конфигурация

Переменные окружения (`.env` или системные):

| Переменная    | Описание                                  | По умолчанию |
|---------------|-------------------------------------------|--------------|
| SERVER_PORT   | Порт HTTP сервера                         | 8080         |
| CONTENT_DIR   | Директория для хранения видео             | ./content    |
| TLS_CERT_FILE | Путь к TLS-сертификату (опционально)      |              |
| TLS_KEY_FILE  | Путь к TLS-ключу (опционально)            |              |

## API

| Метод  | Путь                              | Описание                  |
|--------|-----------------------------------|---------------------------|
| POST   | /api/reel/{shortcode}             | Скачать рилс              |
| GET    | /api/reel/{shortcode}/video.mp4       | Получить сам рилс         |
| GET    | /api/reel/{shortcode}/description | Получить описание рилса   |
| DELETE | /api/reel/{shortcode}             | Удалить рилс              |
| GET    | /health                           | Health check              |

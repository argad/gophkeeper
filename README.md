# GophKeeper

GophKeeper — это безопасная клиент-серверная система для хранения и управления конфиденциальными данными. Построена на Go, предоставляет REST API с JWT аутентификацией и поддерживает как HTTP, так и HTTPS/TLS коммуникацию.

## Быстрый старт

### Сборка

Соберите сервер и клиент:

```bash
build.bat
```

Это создаст:
- `bin/gophkeeper-server.exe` - Серверное приложение
- `bin/gophkeeper-cli.exe` - CLI клиент

### Запуск с HTTP (Разработка)

**Запустите сервер:**
```bash
bin\gophkeeper-server.exe --config server/config.example.json
```

**Используйте клиент:**
```bash
# Зарегистрировать нового пользователя
bin\gophkeeper-cli.exe register -l myuser -p mypassword

# Войти
bin\gophkeeper-cli.exe login -l myuser -p mypassword

# Сохранить секрет
bin\gophkeeper-cli.exe set -t login -d "username:password" -m "Логин для сайта"

# Получить секреты
bin\gophkeeper-cli.exe get
```

## Выбор storage

### In-Memory хранилище (По умолчанию)

Подходит для разработки и тестирования. Данные не сохраняются.

```bash
bin\gophkeeper-server.exe --storage-type memory
```

### PostgreSQL хранилище



```bash
bin\gophkeeper-server.exe --storage-type postgres --database-dsn "postgresql://[user]:[pass]@localhost:5432/[database]"
```

## Настройка TLS/HTTPS

### Разработка (Самоподписанные сертификаты)

```bash
# Создать приватный ключ и сертификат одной командой
openssl req -x509 -newkey rsa:4096 -keyout certs/server.key -out certs/server.crt -days 365 -nodes -subj "/CN=localhost"

# Запустить сервер c конфигом
bin\gophkeeper-server.exe --config server/config_tls.example.json

# Запустить сервер
bin\gophkeeper-server.exe -enable-tls -tls-cert certs\server.crt -tls-key certs\server.key

# Настроить клиент
set GOPHKEEPER_SERVER_URL=https://localhost:8443
set GOPHKEEPER_INSECURE_TLS=true
```

## CLI команды

### Управление пользователями

```bash
# Регистрация
gophkeeper-cli register -l username -p password

# Вход
gophkeeper-cli login -l username -p password
```

### Управление секретами

```bash
# Сохранить секрет
gophkeeper-cli set -t <тип> -d <данные> -m <метаданные>

# Получить все секреты
gophkeeper-cli get

# Получить конкретный секрет
gophkeeper-cli get -i <id>

# Удалить секрет
gophkeeper-cli delete -i <id>
```

Типы секретов:
- `login` - Логин/Пароль
- `text` - Текстовые данные
- `binary` - Бинарные данные
- `bankcard` - Банковская карта

## Шифрование данных

GophKeeper поддерживает прозрачное шифрование секретных данных при хранении с использованием AES-256-GCM:

```bash
# Сгенерировать безопасный ключ шифрования
openssl rand -base64 32

# Запустить сервер с шифрованием
bin\gophkeeper-server.exe --encryption-key "your-32-byte-encryption-key-here!"

# Или использовать переменную окружения
set ENCRYPTION_KEY=your-32-byte-encryption-key-here!
bin\gophkeeper-server.exe
```

## Структура проекта

```
├── server/                 # Серверное приложение (отдельный Go модуль)
│   ├── cmd/
│   │   ├── gophkeeper-server/  # Точка входа сервера
│   │   └── gencert/            # Генератор сертификатов
│   └── internal/
│       ├── api/                # HTTP обработчики
│       ├── auth/               # Аутентификация и JWT
│       ├── config/             # Управление конфигурацией
│       ├── crypto/             # Шифрование AES-256-GCM
│       ├── models/             # Модели данных сервера
│       ├── storage/            # Слой хранения
│       └── tls/                # TLS утилиты
├── client/                 # CLI клиентское приложение (отдельный Go модуль)
│   ├── cmd/gophkeeper-cli/     # Точка входа клиента
│   └── internal/
│       ├── api/                # API клиент
│       ├── commands/           # CLI команды
│       ├── config/             # Конфигурация клиента
│       └── models/             # Модели данных клиента
├── bin/                    # Скомпилированные бинарники
└── build.bat              # Скрипт сборки
```

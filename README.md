## Chief Machine

### Автовывод с бирж

**Поддерживаемые биржи:**

- Bybit
- Binance
- OKX
- MEXC
- KuCoin

**Модуль включает:**

- Поддержка рандомизации сети вывода
- Поддержка рандомизации токена вывода
- Поддержка рандомизации суммы вывода
- Поддержка рандомизации времени вывода
- Поддержка множества сетей вывода и любых токенов
- Проверка балансов перед выводом для корректного выполнения
- Возможность указания суммы вывода в долларовом эквиваленте
- Ожидание поступления токенов на счёт (не для всех видов токенов)

#### Настройка

1. Посмотреть список допустимых сетей в файле `/config/data/cexs.txt`
2. Заполнить файл `withdraw_addresses.txt` адресами для вывода
3. Указать одну или несколько сетей для вывода (далее в файле user\_config.json)
4. Указать один или несколько токенов для вывода
5. Указать диапазон времени для вывода
6. Указать сумму вывода в долларах
7. При использовании с модулем автобриджа вместо обычных адресов в `withdraw_addresses.txt` указать приватные ключи, а также сеть и токен назначения
   **Важно.** Перед использованием заполните api\_key/secret\_key/password в файле configuration.json

---

### Автобриджер

**Поддерживаемый механизм:**

- Relay

**Модуль включает:**

- Поддержка рандомизации суммы
- Поддержка рандомизации времени
- Поддержка работы с автовыводом

#### Настройка

1. Заполнить файл `withdraw_addresses.txt` приватными ключами
2. Указать сеть назначения (далее в файле user\_config.json)
3. Указать сеть для бриджа
4. Указать диапазон времени
5. Указать диапазон суммы

---

### Автогенератор кошельков

**Поддерживаемые типы:**

- BTC-Segwit
- SVM
- EVM

**Модуль включает:**

- Поддержка задания параметров сохранения кошельков
- Указание количества генерируемых кошельков

#### Настройка

1. Указать тип кошелька (btc/evm/svm) (далее в файле user\_config.json)
2. Указать количество кошельков
3. Указать формат сохранения

---

### Коллектор

**Поддерживаемые сети:**

- EVM

**Модуль включает:**

- Сбор нативных монет и отправка на указанные кошельки

#### Настройка

1. Указать в файле `withdraw_addresses.txt` приватные ключи
2. Указать в конфигурации сети работы "chains"(далее в файле user\_config.json)
3. Указать адрес для отправки средств - destination\_address или прописать ассоциации в destination\_addresses

### configuration.json

**ip\_addresses.** В массив поместить ip, который вы указали в белом списке на бирже.

**cex.** Указать ключи используемой биржи.

**rpc.** Указать RPC сети, если установленная RPC не работает. Если хотите добавить EVM сеть, то укажите ее далее в списке в верхней капитализации как это сделано с остальными сетями.

**threads.** Количество потоков работы программы.

**attention\_gwei.** Максимально допустимый gwei при ончейн работе.

**max\_attention\_time.** Максимальное время ожидания снижения gwei.

**deposit\_waiting\_time.** Время ожидания поступления средств на кошелек при выводе с биржи.

**sleep\_after\_withdraw.** Время паузы после снятия средств с биржи и последующим бриджем.

---

## Установка

- Go (Version 1.24 or newer)
- Git (for cloning the repository)
- Опционально: утилита make&#x20;

### Шаги

1. Клонировать репозиторий:
   ```sh
   git clone https://github.com/ssq0-0/cheif-machine.git
   cd cheif-machine
   go mod download
   go build -o CM ./core/main.go   
   ```
2. Запуск приложения:
   ```sh
   Setup configuration.json && user_config.json
   ./CM
   ```
   Или run main.go:
   ```sh
   cd cheif-machine
   go run ./core/main.go
   ```

## Donations

If this software has been useful to you, consider making a donation:

- **EVM Address:** `0xDe5d4e16C435bC54aDb128039844B01634aE28Ff`
- **SVM Address:** `6xJrAzhGFJ58snkgeVsPpALMkppCHaoYc841REpT5Py`

---

## Community

- **Telegram Channel:** [https://t.me/cheifssq]
- **Telegram Chat:** [https://t.me/chatcheifssq]

---
## Chief Machine

### Auto Withdrawal from Exchanges

**Supported Exchanges:**

- Bybit
- Binance
- OKX
- MEXC
- KuCoin

**Module Includes:**

- Support for randomizing the withdrawal network
- Support for randomizing the withdrawal token
- Support for randomizing the withdrawal amount
- Support for randomizing the withdrawal time
- Support for multiple withdrawal networks and any tokens
- Balance verification before withdrawal for proper execution
- Ability to specify the withdrawal amount in USD equivalent
- Waiting for tokens to be credited to the account (not applicable for all token types)

#### Setup

1. Check the list of allowed networks in the file `/config/data/cexs.txt`
2. Fill the `withdraw_addresses.txt` file with withdrawal addresses
3. Specify one or more networks for withdrawal (in `user_config.json`)
4. Specify one or more tokens for withdrawal
5. Set the time range for withdrawal
6. Specify the withdrawal amount in USD
7. When using the auto-bridge module, instead of regular addresses in `withdraw_addresses.txt`, provide private keys along with the destination network and token.
   **Important:** Before using, fill in `api_key/secret_key/password` in the `configuration.json` file.

---

### Auto Bridger

**Supported Mechanism:**

- Relay

**Module Includes:**

- Support for randomizing the amount
- Support for randomizing the time
- Support for integration with auto-withdrawal

#### Setup

1. Fill the `withdraw_addresses.txt` file with private keys
2. Specify the destination network (in `user_config.json`)
3. Specify the bridge network
4. Set the time range
5. Set the amount range

---

### Auto Wallet Generator

**Supported Types:**

- BTC-Segwit
- SVM
- EVM

**Module Includes:**

- Support for setting wallet saving parameters
- Specifying the number of generated wallets

#### Setup

1. Specify the wallet type (btc/evm/svm) (in `user_config.json`)
2. Specify the number of wallets
3. Specify the saving format

---

### Collector

**Supported Networks:**

- EVM

**Module Includes:**

- Collection of native coins and sending them to specified wallets

#### Setup

1. Specify private keys in the `withdraw_addresses.txt` file
2. Define operational networks in the "chains" section (in `user_config.json`)
3. Specify the destination address - `destination_address` or set associations in `destination_addresses`

### configuration.json

**ip_addresses.** Add the IP address you have whitelisted on the exchange.

**cex.** Specify the API keys of the exchange being used.

**rpc.** Provide the RPC network if the default one is not working. To add an EVM network, specify it in uppercase, following the format of existing networks.

**threads.** Number of working threads in the program.

**attention_gwei.** Maximum acceptable Gwei for on-chain operations.

**max_attention_time.** Maximum waiting time for Gwei reduction.

**deposit_waiting_time.** Waiting time for funds to be credited to the wallet when withdrawing from the exchange.

**sleep_after_withdraw.** Pause time after withdrawing funds from the exchange before the next bridge operation.

---

## Installation

- Go (Version 1.24 or newer)
- Git (for cloning the repository)
- Optional: `make` utility

### Steps

1. Clone the repository:
   ```sh
   git clone https://github.com/ssq0-0/cheif-machine.git
   cd cheif-machine
   go mod download
   go build -o CM ./core/main.go   
   ```
2. Run the application:
   ```sh
   Setup configuration.json && user_config.json
   ./CM
   ```
   Or run `main.go`:
   ```sh
   cd cheif-machine
   go run ./core/main.go
   ```

## Donations

If this software has been useful to you, consider making a donation:

- **EVM Address:** `0xDe5d4e16C435bC54aDb128039844B01634aE28Ff`
- **SVM Address:** `6xJrAzhGFJ58snkgeVsPpALMkppCHaoYc841REpT5Py`

---

## Community

- **Telegram Channel:** [https://t.me/cheifssq]
- **Telegram Chat:** [https://t.me/chatcheifssq]

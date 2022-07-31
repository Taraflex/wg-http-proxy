# wg-proxer

Портабельный WireGuard клиент предоставляюший HTTP прокси сервер 

## Цель

WireGuard (как и иные vpn клиенты) хорош когда требуется завернуть большую часть трафика на всем устройстве, но не так удобен для оборачивания трафика в отдельных программах и/или к отдельным ресурсам. 

Наиболее гибким в этом отношении является HTTP прокси (через PAC скрипты), однако браузеры не умеют даже в Basic авторизацию для HTTP прокси.

Данный проект - попытка заткнуть слабые места обоих решений.
По сети гоняется защищенный WireGuard трафик, а на клиентском устройстве запускается прокси сервер доступный приложениям. 

## Установка

Установка не требуется. 
[Просто скачайте версию под свою систему](https://github.com/Taraflex/wg-proxer/releases/tag/latest) и запускайте хоть с флешки.

## Как использовать

`wg-proxer.exe -p 8087 wg.conf`
- `-p 8087` локальный прокси порт 
- `-s` выключает логгирование
- `-v` включает логгирование запросов
- `wg.conf` ПОСЛЕДНИМ параметром указываем путь к конфигу Wireguard (имя файла может быть любое)

### Пример wg конфига

```ini
[Interface]
PrivateKey = xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=
Address = 10.8.0.5/24
DNS = 8.8.8.8, 8.8.4.4

[Peer]
PublicKey = xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=
PresharedKey = xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=
AllowedIPs = 192.168.0.0/16, 0.0.0.0/5, 213.5.120.3/32
Endpoint = x.x.x.x:2022
PersistentKeepalive = 25
```
**ВНИМАНИЕ для полей разрешающих множественные значения допускается только синтаксис с перечислением значений через зяпятую, но не дубликаты ключей.**
```ini
;;так нельзя
DNS = 8.8.8.8
DNS = 8.8.4.4
```
```ini
;;так можно
DNS = 8.8.8.8, 8.8.4.4
```

## Что не так с настройкой Address

Маски на данный момент просто игнорируются.

Я не понял как прикрутить диапазон по маске. В https://github.com/octeep/wireproxy (что работает на основе той же библиотеки) автор например велел указывать маски /32 /128 и сделал вид что все работает - мы поступим также. 

## Демонизация

- Windows - [проще всего использовать nssm](https://nssm.cc/usage)
- Linux - Используйте systemd 

## Похожие проекты и недостатки

- https://github.com/samhza/wg-tcp-proxy - не умеет читать WireGuard конфиги, мало параметров конфигурации
- https://github.com/shimberger/wg-http-proxy - не умеет читать WireGuard конфиги, мало параметров конфигурации
- https://github.com/zhsj/wghttp - не умеет читать WireGuard конфиги, не умеет менять AllowedIPs, MTU, PersistentKeepalive, ListenPort  (но есть иные плюшки https://github.com/zhsj/wghttp/blob/master/docs/usage.md - если бы стразу нашел этот проект, то не писал бы свой)
- https://github.com/octeep/wireproxy - может все что мне требовалось, но запускает SOCKS5 

## Полезное

- [Обход блокировок РКН](./PKH.md)
- [Калькулятор AllowedIPs](https://www.procustodibus.com/blog/2021/03/wireguard-allowedips-calculator/)
- [Список специальных ip](https://blog.bullspit.co.uk/2016/11/15/public-internet-ipv4-prefixes/) - если добавить в AllowedIPs то решаться проблемы с доступом к локальным ресурсам
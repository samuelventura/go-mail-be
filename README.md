# go-mail-be

Send only mailer with REST API

- curl is the cli
- open relay 
- no policy

## api

```bash
# add domain
curl -X POST http://127.0.0.1:31650/api/domain/domain.tld
# show domain
curl -X GET http://127.0.0.1:31650/api/domain/domain.tld
# delete domain
curl -X DELETE http://127.0.0.1:31650/api/domain/domain.tld
# list domain names
curl -X GET http://127.0.0.1:31650/api/domain
```

### post message

```bash
ssh -L 23
curl -H "Mail-From: i@i.com" -H "Mail-To: u@u.com" -X POST --data 'email body' http://localhost:port/message
```

- Log-Level
- Mail-From
- Mail-To
- Mail-Id
- Mail-Status
- Content-Type
  - text/plain; charset=UTF-8
  - text/html; charset=UTF-8

## linux

```bash
#for go-sqlite
sudo apt install build-essential
```

# go-mail-be

Send only mailer with REST API

- curl is the cli
- open relay 
- no policy

## api

### post message

```bash
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

### show message

```bash
curl -X GET http://localhost:port/message/{id}/body
curl -X GET http://localhost:port/message/{id}/status
```

### add domain

```bash
curl -X POST http://localhost:port/domain/{domain}
```

### remove domain

```bash
curl -X DELETE http://localhost:port/domain/{domain}
```

### show domain

```bash
curl -X GET http://localhost:port/domain/{domain}
```

### list domains

```bash
curl -X GET http://localhost:port/domain/
```

## linux

sudo apt install build-essential

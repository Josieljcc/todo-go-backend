# 📧 Configuração de Notificações

Este guia explica como configurar as notificações por email e Telegram.

## 📋 Pré-requisitos

1. **Email SMTP**: Conta de email com acesso SMTP (Gmail, Outlook, etc.)
2. **Telegram Bot** (opcional): Bot do Telegram para notificações

---

## 📧 Configuração de Email (SMTP)

### Gmail

1. Ative a verificação em duas etapas na sua conta Google
2. Gere uma "Senha de app":
   - Acesse: https://myaccount.google.com/apppasswords
   - Selecione "App" e "Outro (nome personalizado)"
   - Digite "Todo API" e clique em "Gerar"
   - Copie a senha gerada (16 caracteres)

3. Configure no `.env`:
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=seu-email@gmail.com
SMTP_PASSWORD=senha-de-app-gerada
SMTP_FROM=noreply@todoapp.com
```

### Outlook/Hotmail

```env
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587
SMTP_USER=seu-email@outlook.com
SMTP_PASSWORD=sua-senha
SMTP_FROM=seu-email@outlook.com
```

### Outros provedores

Consulte a documentação do seu provedor de email para as configurações SMTP.

---

## 🤖 Configuração do Telegram Bot

### Passo 1: Criar o Bot

1. Abra o Telegram e procure por **@BotFather**
2. Envie o comando `/newbot`
3. Escolha um nome para o bot (ex: "Todo Notifications Bot")
4. Escolha um username (deve terminar em `bot`, ex: `todo_notifications_bot`)
5. **Copie o token** fornecido pelo BotFather (algo como: `123456789:ABCdefGHIjklMNOpqrsTUVwxyz`)

### Passo 2: Configurar no `.env`

```env
TELEGRAM_BOT_TOKEN=seu-token-aqui
```

### Passo 3: Obter o Chat ID do Usuário

**Opção A: Via Bot (Recomendado)**

1. Envie uma mensagem para o seu bot no Telegram
2. Acesse: `https://api.telegram.org/bot<SEU_TOKEN>/getUpdates`
3. Procure por `"chat":{"id":123456789}` no JSON retornado
4. O número `123456789` é o seu Chat ID

**Opção B: Via @userinfobot**

1. Procure por **@userinfobot** no Telegram
2. Inicie uma conversa
3. O bot retornará seu Chat ID

### Passo 4: Configurar Chat ID no Sistema

Use o endpoint da API para configurar:

```bash
PUT /api/v1/users/telegram-chat-id
Authorization: Bearer <seu-token-jwt>

{
  "telegram_chat_id": "123456789"
}
```

---

## ⚙️ Configuração Geral

### Variáveis de Ambiente

Adicione ao seu `.env`:

```env
# Ativar/desativar notificações
NOTIFICATIONS_ENABLED=true

# Intervalo de verificação (formato cron)
# Exemplos:
# "0 * * * *"     - A cada hora
# "0 */6 * * *"   - A cada 6 horas
# "0 9 * * *"     - Diariamente às 9h
# "*/15 * * * *"  - A cada 15 minutos
NOTIFICATION_CHECK_INTERVAL=0 * * * *

# Email SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=seu-email@gmail.com
SMTP_PASSWORD=sua-senha-app
SMTP_FROM=noreply@todoapp.com

# Telegram Bot
TELEGRAM_BOT_TOKEN=seu-token-do-botfather
```

---

## 🔔 Tipos de Notificações

O sistema envia automaticamente:

1. **Due Soon** (1 dia antes): Notificação quando a tarefa vence amanhã
2. **Due Today**: Notificação quando a tarefa vence hoje
3. **Overdue**: Notificação diária para tarefas atrasadas

---

## 👤 Configuração por Usuário

### Ativar/Desativar Notificações

```bash
PUT /api/v1/users/notifications-enabled
Authorization: Bearer <token>

{
  "notifications_enabled": true
}
```

### Configurar Telegram Chat ID

```bash
PUT /api/v1/users/telegram-chat-id
Authorization: Bearer <token>

{
  "telegram_chat_id": "123456789"
}
```

Para remover o Telegram:
```json
{
  "telegram_chat_id": null
}
```

---

## 🧪 Testando

### Teste Manual

1. Crie uma tarefa com `due_date` = hoje
2. Aguarde o próximo ciclo do scheduler (ou ajuste o intervalo)
3. Verifique seu email e Telegram

### Verificar Logs

O scheduler registra no log:
```
Running notification check...
Notification check completed
```

---

## ❓ Troubleshooting

### Email não está sendo enviado

- Verifique as credenciais SMTP
- Para Gmail, use "Senha de app" (não a senha normal)
- Verifique se o firewall não está bloqueando a porta SMTP

### Telegram não está funcionando

- Verifique se o token do bot está correto
- Verifique se o Chat ID está correto
- Envie uma mensagem para o bot antes de configurar o Chat ID
- Verifique os logs do servidor para erros

### Notificações não estão sendo enviadas

- Verifique se `NOTIFICATIONS_ENABLED=true`
- Verifique se o usuário tem `notifications_enabled=true`
- Verifique se a tarefa tem `due_date` configurado
- Verifique se a tarefa não está `completed=true`

---

## 📝 Notas

- Notificações são enviadas apenas uma vez por dia para cada tipo
- Tarefas completadas não recebem notificações
- O scheduler roda em background e não bloqueia a API
- Histórico de notificações é salvo no banco de dados


# 🔍 Troubleshooting de Notificações

Este guia ajuda a diagnosticar problemas com as notificações.

## 📋 Checklist de Verificação

### 1. Verificar Configuração no `.env`

Certifique-se de que todas as variáveis estão configuradas:

```env
# Ativar notificações globalmente
NOTIFICATIONS_ENABLED=true

# Intervalo do scheduler (formato cron)
NOTIFICATION_CHECK_INTERVAL=0 * * * *

# Email SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=seu-email@gmail.com
SMTP_PASSWORD=sua-senha-app
SMTP_FROM=noreply@todoapp.com

# Telegram
TELEGRAM_BOT_TOKEN=seu-token-aqui
```

### 2. Verificar Configuração do Usuário

Use o endpoint de debug para verificar:

```bash
GET /api/v1/notifications/debug
Authorization: Bearer <token>
```

Isso retorna:
- Se `notifications_enabled` está `true`
- Se `email` está configurado
- Se `telegram_chat_id` está configurado
- Lista de tarefas com `due_date`
- Histórico de notificações enviadas

### 3. Verificar Tarefa

A tarefa deve ter:
- ✅ `due_date` configurado (não pode ser `null`)
- ✅ `completed = false`
- ✅ `due_date` = hoje, amanhã, ou no passado

### 4. Verificar Logs do Servidor

Após executar `POST /api/v1/notifications/test`, verifique os logs:

```
Starting notification check at 2024-12-29 10:00:00
Today: 2024-12-29, Tomorrow: 2024-12-30
Found X tasks with due dates
Task Y: due_date=2024-12-29, user_id=1, notifications_enabled=true, email=user@example.com, telegram_chat_id=123456789
Task Y: DUE TODAY
Sending email notification for task Y to user@example.com
Email notification sent successfully for task Y
Sending telegram notification for task Y to chat 123456789
Telegram notification sent successfully for task Y
```

## 🐛 Problemas Comuns

### ❌ "Notifications are disabled"

**Causa**: `NOTIFICATIONS_ENABLED=false` no `.env`

**Solução**: 
```env
NOTIFICATIONS_ENABLED=true
```

### ❌ "Task X: skipping (user notifications disabled)"

**Causa**: O usuário tem `notifications_enabled=false` no banco

**Solução**:
```bash
PUT /api/v1/users/notifications-enabled
{
  "notifications_enabled": true
}
```

### ❌ "Task X: user has no email address"

**Causa**: O usuário não tem email cadastrado

**Solução**: Verifique se o usuário foi criado com email. O email é obrigatório no registro.

### ❌ "Task X: user has no telegram chat ID"

**Causa**: Telegram Chat ID não foi configurado

**Solução**:
```bash
PUT /api/v1/users/telegram-chat-id
{
  "telegram_chat_id": "123456789"
}
```

### ❌ "Failed to send email notification: email service not configured"

**Causa**: Variáveis SMTP não estão configuradas ou estão vazias

**Solução**: Verifique se todas as variáveis SMTP estão no `.env`:
- `SMTP_HOST`
- `SMTP_PORT`
- `SMTP_USER`
- `SMTP_PASSWORD`
- `SMTP_FROM`

### ❌ "Failed to send email notification: failed to send email: ..."

**Causa**: Erro de autenticação ou conexão SMTP

**Soluções**:
1. **Gmail**: Use "Senha de app" (não a senha normal)
   - Ative verificação em duas etapas
   - Gere senha de app em: https://myaccount.google.com/apppasswords

2. **Outlook**: Verifique se a senha está correta

3. **Firewall**: Verifique se a porta SMTP não está bloqueada

4. **Teste manual**:
   ```bash
   telnet smtp.gmail.com 587
   ```

### ❌ "Failed to send telegram notification: telegram bot token not configured"

**Causa**: `TELEGRAM_BOT_TOKEN` não está configurado

**Solução**: Adicione o token do bot no `.env`

### ❌ "Failed to send telegram notification: user telegram chat ID not configured"

**Causa**: Chat ID do usuário não foi configurado

**Solução**: Configure o Chat ID via API:
```bash
PUT /api/v1/users/telegram-chat-id
{
  "telegram_chat_id": "123456789"
}
```

### ❌ "Failed to send telegram notification: telegram API error: ..."

**Causa**: Erro na API do Telegram

**Possíveis causas**:
1. Token do bot inválido
2. Chat ID incorreto
3. Bot não foi iniciado (envie `/start` para o bot primeiro)

**Solução**:
1. Teste o token:
   ```bash
   curl https://api.telegram.org/bot<SEU_TOKEN>/getMe
   ```

2. Teste enviar mensagem manualmente:
   ```bash
   curl -X POST https://api.telegram.org/bot<SEU_TOKEN>/sendMessage \
     -d "chat_id=123456789&text=Teste"
   ```

3. Certifique-se de ter enviado uma mensagem para o bot antes de configurar o Chat ID

### ❌ "Email notification already sent today for task X, skipping"

**Causa**: Notificação já foi enviada hoje (sistema evita duplicatas)

**Solução**: Isso é normal! O sistema só envia uma notificação por tipo por dia. Aguarde até amanhã ou delete o registro no banco:
```sql
DELETE FROM notifications WHERE task_id = X AND DATE(sent_at) = CURDATE();
```

### ❌ "Task X: not due yet (due 2024-12-30)"

**Causa**: A tarefa vence no futuro (mais de 1 dia)

**Solução**: Crie uma tarefa que vence:
- **Hoje**: `due_date` = data de hoje
- **Amanhã**: `due_date` = data de amanhã
- **Atrasada**: `due_date` = data no passado

## 🔧 Endpoints de Debug

### 1. Testar Notificações Manualmente

```bash
POST /api/v1/notifications/test
Authorization: Bearer <token>
```

Isso executa a verificação imediatamente e mostra logs detalhados.

### 2. Ver Informações de Debug

```bash
GET /api/v1/notifications/debug
Authorization: Bearer <token>
```

Retorna:
- Configuração do usuário
- Tarefas com `due_date`
- Histórico de notificações

## 📝 Exemplo de Teste Completo

```bash
# 1. Login
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"seu-usuario","password":"sua-senha"}' \
  | jq -r '.token')

# 2. Verificar configuração
curl -X GET http://localhost:8080/api/v1/notifications/debug \
  -H "Authorization: Bearer $TOKEN"

# 3. Ativar notificações (se necessário)
curl -X PUT http://localhost:8080/api/v1/users/notifications-enabled \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"notifications_enabled":true}'

# 4. Configurar Telegram (se necessário)
curl -X PUT http://localhost:8080/api/v1/users/telegram-chat-id \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"telegram_chat_id":"123456789"}'

# 5. Criar tarefa que vence hoje
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Teste de Notificação",
    "type": "trabalho",
    "due_date": "2024-12-29T23:59:59Z"
  }'

# 6. Testar notificações
curl -X POST http://localhost:8080/api/v1/notifications/test \
  -H "Authorization: Bearer $TOKEN"

# 7. Verificar logs do servidor para detalhes
```

## 🔍 Verificações no Banco de Dados

Se necessário, verifique diretamente no banco:

```sql
-- Ver configuração do usuário
SELECT id, username, email, notifications_enabled, telegram_chat_id 
FROM users WHERE id = 1;

-- Ver tarefas com due_date
SELECT id, title, due_date, completed, user_id 
FROM tasks 
WHERE user_id = 1 AND due_date IS NOT NULL AND completed = false;

-- Ver notificações enviadas
SELECT * FROM notifications 
WHERE user_id = 1 
ORDER BY sent_at DESC 
LIMIT 10;
```

## 💡 Dicas

1. **Sempre verifique os logs** após executar o teste
2. **Use o endpoint de debug** para verificar a configuração
3. **Teste email e Telegram separadamente** para isolar problemas
4. **Para Gmail**, use sempre "Senha de app"
5. **Para Telegram**, envie uma mensagem para o bot antes de configurar o Chat ID


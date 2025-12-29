# 🧪 Como Testar Notificações

Este guia mostra como testar o sistema de notificações de forma rápida e eficiente.

## 🚀 Método Rápido: Endpoint de Teste

### 1. Criar uma Tarefa de Teste

Crie uma tarefa que vence **hoje** ou **amanhã**:

```bash
POST /api/v1/tasks
Authorization: Bearer <seu-token>

{
  "title": "Tarefa de Teste",
  "description": "Testando notificações",
  "type": "trabalho",
  "priority": "alta",
  "due_date": "2024-12-29T23:59:59Z"  # Ajuste para hoje ou amanhã
}
```

### 2. Forçar Verificação de Notificações

Use o endpoint de teste para verificar imediatamente (sem esperar o scheduler):

```bash
POST /api/v1/notifications/test
Authorization: Bearer <seu-token>
```

**Resposta esperada:**
```json
{
  "message": "Notification check completed. Check your email and Telegram.",
  "data": null
}
```

### 3. Verificar Resultados

- ✅ **Email**: Verifique sua caixa de entrada (e spam)
- ✅ **Telegram**: Verifique as mensagens do bot
- ✅ **Logs**: Verifique os logs do servidor

---

## 📋 Método Completo: Teste Manual

### Passo 1: Configurar Telegram Chat ID

```bash
PUT /api/v1/users/telegram-chat-id
Authorization: Bearer <seu-token>

{
  "telegram_chat_id": "123456789"
}
```

### Passo 2: Garantir que Notificações Estão Ativas

```bash
PUT /api/v1/users/notifications-enabled
Authorization: Bearer <seu-token>

{
  "notifications_enabled": true
}
```

### Passo 3: Criar Tarefas de Teste

#### Tarefa que vence hoje:
```json
{
  "title": "Tarefa vence hoje",
  "type": "trabalho",
  "due_date": "2024-12-29T23:59:59Z"  # Ajuste para hoje
}
```

#### Tarefa que vence amanhã:
```json
{
  "title": "Tarefa vence amanhã",
  "type": "trabalho",
  "due_date": "2024-12-30T23:59:59Z"  # Ajuste para amanhã
}
```

#### Tarefa atrasada:
```json
{
  "title": "Tarefa atrasada",
  "type": "trabalho",
  "due_date": "2024-12-28T23:59:59Z"  # Data no passado
}
```

### Passo 4: Executar Verificação

**Opção A: Usar endpoint de teste (recomendado)**
```bash
POST /api/v1/notifications/test
```

**Opção B: Aguardar o scheduler**
- Por padrão, o scheduler roda a cada hora
- Você pode ajustar `NOTIFICATION_CHECK_INTERVAL` no `.env`:
  ```env
  NOTIFICATION_CHECK_INTERVAL=*/5 * * * *  # A cada 5 minutos (para teste)
  ```

### Passo 5: Verificar Notificações Enviadas

Verifique os logs do servidor:
```
Running notification check...
Notification check completed
```

Se houver erros:
```
Failed to send email notification: ...
Failed to send telegram notification: ...
```

---

## 🔍 Verificações de Troubleshooting

### Email não chegou?

1. **Verifique as credenciais SMTP:**
   ```bash
   # Teste SMTP manualmente (opcional)
   telnet smtp.gmail.com 587
   ```

2. **Verifique os logs:**
   - Procure por "Failed to send email notification"
   - Verifique se há erros de autenticação

3. **Para Gmail:**
   - Certifique-se de usar "Senha de app" (não a senha normal)
   - Verifique se a verificação em duas etapas está ativa

### Telegram não chegou?

1. **Verifique o token do bot:**
   ```bash
   # Teste se o bot está funcionando
   curl https://api.telegram.org/bot<SEU_TOKEN>/getMe
   ```

2. **Verifique o Chat ID:**
   - Certifique-se de ter enviado uma mensagem para o bot primeiro
   - Verifique se o Chat ID está correto no banco

3. **Verifique os logs:**
   - Procure por "Failed to send telegram notification"
   - Verifique se há erros da API do Telegram

### Nenhuma notificação foi enviada?

1. **Verifique se a tarefa tem `due_date`:**
   ```sql
   SELECT id, title, due_date, completed FROM tasks WHERE user_id = <seu-id>;
   ```

2. **Verifique se o usuário tem notificações ativadas:**
   ```sql
   SELECT id, username, notifications_enabled FROM users WHERE id = <seu-id>;
   ```

3. **Verifique se já foi notificado hoje:**
   ```sql
   SELECT * FROM notifications 
   WHERE user_id = <seu-id> 
   AND DATE(sent_at) = CURDATE();
   ```

---

## 🎯 Checklist de Teste

- [ ] Configurou SMTP no `.env`
- [ ] Configurou Telegram Bot Token no `.env`
- [ ] Obteve Chat ID do Telegram
- [ ] Configurou Chat ID via API
- [ ] Criou tarefa com `due_date` = hoje
- [ ] Executou `POST /notifications/test`
- [ ] Recebeu email
- [ ] Recebeu mensagem no Telegram
- [ ] Verificou logs do servidor

---

## 💡 Dicas

1. **Para testes rápidos**, ajuste o intervalo do scheduler:
   ```env
   NOTIFICATION_CHECK_INTERVAL=*/1 * * * *  # A cada minuto (apenas para teste!)
   ```

2. **Para testar diferentes tipos**, crie tarefas com diferentes `due_date`:
   - Hoje → `due_today`
   - Amanhã → `due_soon`
   - Ontem → `overdue`

3. **Para evitar spam**, o sistema só envia uma notificação por tipo por dia

4. **Para limpar notificações de teste**, você pode deletar do banco:
   ```sql
   DELETE FROM notifications WHERE task_id IN (SELECT id FROM tasks WHERE title LIKE '%teste%');
   ```

---

## 📝 Exemplo Completo de Teste

```bash
# 1. Login
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"seu-usuario","password":"sua-senha"}' \
  | jq -r '.token')

# 2. Configurar Telegram
curl -X PUT http://localhost:8080/api/v1/users/telegram-chat-id \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"telegram_chat_id":"123456789"}'

# 3. Criar tarefa que vence hoje
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Teste de Notificação",
    "type": "trabalho",
    "due_date": "2024-12-29T23:59:59Z"
  }'

# 4. Forçar verificação
curl -X POST http://localhost:8080/api/v1/notifications/test \
  -H "Authorization: Bearer $TOKEN"

# 5. Verificar email e Telegram!
```

---

Pronto! Agora você pode testar as notificações facilmente! 🎉


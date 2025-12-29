# 📱 Como Obter o Telegram Chat ID

O erro "chat not found" acontece quando o Chat ID está incorreto ou o usuário não iniciou uma conversa com o bot.

## 🔍 Passo a Passo

### 1. Criar o Bot (se ainda não criou)

1. Abra o Telegram e procure por **@BotFather**
2. Envie `/newbot`
3. Escolha um nome e username para o bot
4. **Copie o token** fornecido (algo como: `123456789:ABCdefGHIjklMNOpqrsTUVwxyz`)

### 2. Configurar o Token no `.env`

```env
TELEGRAM_BOT_TOKEN=seu-token-aqui
```

### 3. Obter o Chat ID

**IMPORTANTE**: Você DEVE enviar uma mensagem para o bot ANTES de obter o Chat ID!

#### Opção A: Via @userinfobot (Mais Fácil)

1. Procure por **@userinfobot** no Telegram
2. Inicie uma conversa
3. O bot retornará seu Chat ID (um número como `123456789`)

#### Opção B: Via API do Telegram

1. **Primeiro**: Envie uma mensagem para o SEU bot (qualquer mensagem, como "oi" ou "/start")

2. Depois, acesse no navegador:
   ```
   https://api.telegram.org/bot<SEU_TOKEN>/getUpdates
   ```

3. Procure no JSON retornado por:
   ```json
   {
     "message": {
       "chat": {
         "id": 123456789
       }
     }
   }
   ```

4. O número `123456789` é o seu Chat ID

#### Opção C: Via @getidsbot

1. Procure por **@getidsbot** no Telegram
2. Inicie uma conversa
3. O bot retornará seu Chat ID

### 4. Configurar o Chat ID na API

```bash
PUT /api/v1/users/telegram-chat-id
Authorization: Bearer <seu-token>

{
  "telegram_chat_id": "123456789"
}
```

**IMPORTANTE**: 
- O Chat ID deve ser uma string numérica (ex: `"123456789"`)
- Para grupos, o Chat ID pode ser negativo (ex: `"-123456789"`)
- Você DEVE ter enviado pelo menos uma mensagem para o bot antes de configurar

## ❌ Erros Comuns

### "chat not found"

**Causa**: O usuário não enviou uma mensagem para o bot ainda.

**Solução**:
1. Abra o Telegram
2. Procure pelo seu bot (pelo username que você criou)
3. Envie uma mensagem qualquer (ex: "/start" ou "oi")
4. Depois configure o Chat ID novamente

### "invalid bot token"

**Causa**: O token do bot está incorreto.

**Solução**:
1. Verifique se copiou o token completo do BotFather
2. Verifique se não há espaços extras no `.env`
3. Teste o token:
   ```bash
   curl https://api.telegram.org/bot<SEU_TOKEN>/getMe
   ```

### "bot was blocked by user"

**Causa**: O usuário bloqueou o bot.

**Solução**:
1. Desbloqueie o bot no Telegram
2. Envie uma mensagem para o bot novamente
3. Configure o Chat ID novamente

## 🧪 Testar a Configuração

### 1. Verificar se o bot está funcionando

```bash
curl https://api.telegram.org/bot<SEU_TOKEN>/getMe
```

Deve retornar informações sobre o bot.

### 2. Testar envio manual

```bash
curl -X POST https://api.telegram.org/bot<SEU_TOKEN>/sendMessage \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": "123456789",
    "text": "Teste de mensagem"
  }'
```

Se funcionar, você receberá a mensagem no Telegram.

### 3. Verificar Chat ID no sistema

```bash
GET /api/v1/notifications/debug
Authorization: Bearer <token>
```

Isso mostra o Chat ID configurado no sistema.

## 💡 Dicas

1. **Sempre envie uma mensagem para o bot primeiro** antes de configurar o Chat ID
2. **Para grupos**: O Chat ID será negativo (ex: `-123456789`)
3. **Para canais**: Use o formato `@channelusername` ou o ID numérico
4. **Teste manualmente** antes de usar no sistema
5. **Mantenha o bot desbloqueado** para receber notificações

## 📝 Exemplo Completo

```bash
# 1. Criar bot no @BotFather e copiar token
# 2. Adicionar token no .env
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz

# 3. Enviar mensagem para o bot no Telegram
# (abra o bot e envie "/start" ou qualquer mensagem)

# 4. Obter Chat ID
curl https://api.telegram.org/bot123456789:ABCdefGHIjklMNOpqrsTUVwxyz/getUpdates

# 5. Copiar o chat.id do JSON retornado (ex: 987654321)

# 6. Configurar na API
curl -X PUT http://localhost:8080/api/v1/users/telegram-chat-id \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"telegram_chat_id": "987654321"}'

# 7. Testar
curl -X POST http://localhost:8080/api/v1/notifications/test \
  -H "Authorization: Bearer <token>"
```

## 🔄 Se Ainda Não Funcionar

1. Verifique os logs do servidor para ver o erro completo
2. Teste o envio manual via curl (passo 2 acima)
3. Verifique se o Chat ID está correto no banco de dados
4. Certifique-se de que enviou uma mensagem para o bot recentemente


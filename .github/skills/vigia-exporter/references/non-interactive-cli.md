# CLI não-interativo

## Objetivo

Permitir coleta em cron, CI e Ansible sem prompts stdin.

## Flags padrão (todos os providers)

```bash
vigia-export <provider> \
  --host localhost \
  --port 1433 \
  --user vigia_reader \
  --password-env VIGIA_DB_PASSWORD \
  --alias prod-erp \
  --output /var/lib/vigia/snapshot.json \
  --format v1 \
  --quiet
```

## Variáveis de ambiente

| Variável | Descrição |
|----------|-----------|
| `VIGIA_DB_HOST` | Host do banco |
| `VIGIA_DB_PORT` | Porta |
| `VIGIA_DB_USER` | Usuário |
| `VIGIA_DB_PASSWORD` | Senha (preferir env a flag) |
| `VIGIA_HOST_ALIAS` | Alias da instância |
| `VIGIA_OUTPUT` | Caminho do arquivo |

Precedência: flag > env > default interativo.

## Exemplo cron

```cron
*/15 * * * * /usr/local/bin/vigia-export sqlserver \
  --host sql.prod.internal --port 1433 \
  --user vigia_reader --password-env VIGIA_DB_PASSWORD \
  --alias prod-erp -o /tmp/snapshot.json -q && \
  curl -sf -X POST https://api.getvigia.com/v1/snapshots \
    -H "Authorization: Bearer $VIGIA_API_KEY" \
    -H "Content-Type: application/json" \
    -d @- <<EOF
{"instance_id":"$VIGIA_INSTANCE_ID","payload":$(cat /tmp/snapshot.json)}
EOF
```

## Códigos de saída

| Código | Significado |
|--------|-------------|
| 0 | Sucesso |
| 1 | Erro de conexão |
| 2 | Erro de execução de script |
| 3 | Erro de escrita |
| 4 | Validação de schema falhou |

## Segurança

- Nunca logar senha
- `--password-env` preferido sobre `--password`
- Modo `--quiet` suprime stdout exceto erros em stderr

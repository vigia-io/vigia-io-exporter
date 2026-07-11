# Licença do vigia-io-exporter — resumo

> Este documento é **resumo em português**. Em caso de conflito, prevalece o texto em inglês em [`LICENSE`](../../LICENSE).

## Por que não é MIT?

O Vigia é um **produto comercial**. O exporter é público para **transparência** (você vê o que roda no seu bastion), não para virar base de um concorrente gratuito.

## O que você pode fazer

| Uso | Permitido? |
|-----|------------|
| Ler e auditar o código no GitHub | ✅ |
| Enviar pull requests | ✅ |
| Rodar em produção **com conta Vigia** (API key válida) | ✅ |
| Testar em homologação (até 30 dias ou 3 instâncias) | ✅ |
| Rodar em produção **sem** Vigia | ❌ |
| Vender/hostear um produto concorrente usando este código | ❌ |
| Embutir o exporter em produto de terceiros | ❌ (sem autorização) |

## Modelo mental

```text
Exporter (público, licença restritiva)
        │
        ▼ snapshot.json
Plataforma Vigia (SaaS — onde está o valor)
        │
        ▼
Console, alertas, histórico, billing
```

**Aberto** = código visível.  
**Não open source** = uso em produção amarrado ao ecossistema Vigia.

## Exceções

Empresas com requisitos especiais (air-gap prolongado, OEM, white-label) podem negociar licença alternativa: **legal@getvigia.com**.

## Repositórios e licença

| Repositório | Licença |
|-------------|---------|
| vigia-io-exporter (público) | Vigia Source Available v1.0 |
| vigia-io-platform, console, scripts | Proprietário (privado) |

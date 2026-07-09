# Licença e uso permitido

## Modelo

O exporter é **source-available**, não MIT:

- Código no GitHub = transparência
- Produção = requer conta Vigia (API key) ou licença comercial
- Homologação = até 30 dias ou 3 instâncias

Texto legal: [`LICENSE`](../../../templates/exporter/LICENSE)  
Resumo PT: [`docs/legal/exporter-license.md`](../../../docs/legal/exporter-license.md)  
ADR: [`docs/adr/0008-exporter-source-available-license.md`](../../../docs/adr/0008-exporter-source-available-license.md)

## O que implementar no CLI

```go
const licenseNotice = `
Vigia Exporter — Vigia Source Available License v1.0
Uso em produção requer conta Vigia: https://getvigia.com
`
```

Exibir no `--help` e na primeira execução interativa.

## Upload de snapshot

O fluxo canônico sempre termina em:

```http
POST https://api.getvigia.com/v1/snapshots
Authorization: Bearer vigia_sk_live_...
```

Documentar `VIGIA_API_KEY` como obrigatório para produção — não promover "uso offline permanente" como produto.

## Tier de scripts no catálogo

| Tier | Significado |
|------|-------------|
| `open` | Embutido no exporter; ainda sujeito à licença do binário |
| `premium` | Só via assinatura; distribuído pelo build privado |

`tier: open` **não** significa licença MIT.

# PayStream

Global payroll on Stellar. Pay distributed teams in stablecoins, settle in 47 local currencies, close the books on time.

## Capabilities

| Capability | Status |
|---|---|
| Batched payouts up to 1,000 recipients per run | Live |
| USDC, EURC, and 12 anchored stablecoins | Live |
| Local-currency delivery in 47 countries via SEP-31 | Live |
| Recurring schedules: weekly, bi-weekly, monthly, custom CRON | Live |
| Multi-sig funding wallet (2-of-3, 3-of-5 configurable up to 20-of-20) | Live |
| 1099, W-2, W-8BEN tax form generation | Beta |
| SAML SSO and SCIM provisioning | Beta |
| Webhook events (18 event types) | Live |
| REST API with idempotency keys | Live |
| Audit log streaming to Splunk, Datadog, S3 | Live |
| SOC 2 Type II | In progress (target Q3 2026) |

## Who PayStream is for

PayStream is built for finance and people-ops teams running distributed workforces of 10 to 10,000 people. It replaces:

- **Deel, Remote, Papaya Global** for contractor and EOR payouts where local-currency delivery is the bottleneck
- **Wise Business, Revolut Business** for batched cross-border salary runs
- **Custom scripts** wired together by an engineering team that grew the company faster than the finance stack

If you currently pay your team via 14 different vendor invoices on the 1st of every month, PayStream collapses that into one funding transaction and one approval click.

## Quickstart

### Prerequisites

- Go 1.22+
- PostgreSQL 14+
- A Stellar wallet funded with USDC (mainnet) or testnet USDC for development
- Anchor credentials for at least one corridor

### Run locally

```bash
git clone https://github.com/breedar/paystream.git
cd paystream
make bootstrap     # installs go deps, pnpm deps, runs migrations
make dev           # starts API on :8080, dashboard on :3000
```

### Send your first batch

```bash
curl -X POST http://localhost:8080/v1/payouts/batch \
  -H "Authorization: Bearer $PAYSTREAM_API_KEY" \
  -H "Idempotency-Key: payroll-2026-05-01" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "May 2026 contractors",
    "asset": "USDC",
    "items": [
      {"recipient_id": "rec_abc", "amount": "1500.00"},
      {"recipient_id": "rec_def", "amount": "2200.00"},
      {"recipient_id": "rec_ghi", "amount": "950.00"}
    ]
  }'
```

The response returns a `batch_id`. PayStream signs and submits the underlying Stellar transactions, monitors settlement on Horizon, and emits `payout.completed` webhooks as each recipient confirms receipt.

## API reference

Full reference at `https://docs.paystream.dev`. Summary of the most-used endpoints below.

### Payouts

#### `POST /v1/payouts`
Single-recipient payout. Returns immediately with status `pending`. Settlement webhook fires when Horizon confirms.

#### `POST /v1/payouts/batch`
Batched payout. Up to 1,000 items per request. PayStream chunks across Stellar transactions (max 100 ops each) and reports an aggregate `batch_id`.

#### `GET /v1/payouts/{id}`
Retrieve payout state, transaction hashes, and last-mile delivery confirmation.

#### `POST /v1/payouts/{id}/cancel`
Cancel a payout that has not yet been signed and submitted. After signing, payouts are irreversible.

### Recipients

#### `POST /v1/recipients`
Register a recipient with their preferred payout method: Stellar wallet, bank account (SEP-31 local rail), or mobile money number.

#### `GET /v1/recipients`
List recipients with optional filters by country, payout method, or tag.

#### `POST /v1/recipients/{id}/method`
Update payout method without losing transaction history.

### Schedules

#### `POST /v1/schedules`
Create a recurring schedule. Example: every Friday 9am UTC, run batch `payroll-weekly-contractors`.

#### `POST /v1/schedules/{id}/pause`
Pause a recurring schedule without deleting it.

### Webhooks

#### `POST /v1/webhook_endpoints`
Register a webhook URL. Events are signed with HMAC-SHA256.

#### Event types
`payout.created`, `payout.signed`, `payout.submitted`, `payout.settled`, `payout.delivered`, `payout.failed`, `batch.created`, `batch.completed`, `batch.partial_failure`, `recipient.created`, `recipient.method_updated`, `schedule.triggered`, `schedule.paused`, `wallet.funded`, `wallet.low_balance`, `signer.added`, `signer.removed`, `audit.suspicious_activity`.

## Security and compliance

### Funding wallet architecture

PayStream never holds customer funds. The funding wallet is a Stellar account owned by the customer, configured with multi-sig:

- Default signer set: 2-of-3 (CFO, CEO, PayStream service signer)
- Configurable up to 20-of-20 for enterprise tier
- The PayStream service signer can only co-sign; it cannot unilaterally move funds
- Daily and per-transaction limits enforced at the protocol level via signer weights and pre-authorized transaction time-bounds

### API authentication

#### API keys
Bearer tokens scoped per environment (test, live). Rotate via dashboard or `POST /v1/api_keys`. Old keys remain valid for 24 hours after rotation.

#### Idempotency
All mutating endpoints accept an `Idempotency-Key` header. Replays within 24 hours return the original response. Safe for retry loops.

#### Webhook signing
Every webhook is signed with HMAC-SHA256 using a per-endpoint secret. Sample verification code in Go, Python, Node, and Ruby lives in `examples/webhook_verification/`.

### Compliance posture

#### SOC 2
Type II audit underway with an AICPA-accredited firm. Evidence collection began January 2026. Report target: Q3 2026.

#### Data residency
Customer data (recipient PII, payout amounts, tax documents) is stored in the customer's chosen region: US, EU, or APAC. Stellar transaction data is on-chain and inherently global.

#### Encryption
TLS 1.3 in transit. AES-256 at rest. Recipient bank account numbers and tax IDs are envelope-encrypted with per-customer KMS keys.

#### Audit logs
Every state-changing action emits an audit event. Stream to Splunk, Datadog, S3, or any HTTPS endpoint via the audit forwarder.

## Integrations

### Accounting
QuickBooks Online, Xero, NetSuite, Sage Intacct. Push payouts as bills with project and class tagging.

### HRIS
Rippling, BambooHR, Gusto, Workday, HiBob. Sync employees and contractors as recipients; PayStream picks up payroll runs from the HRIS as scheduled batches.

### Identity
Okta, Google Workspace, Azure AD via SAML 2.0. SCIM 2.0 for automatic user provisioning.

### Stellar anchors
Cowrie, Vibrant, Pendo, ClickPesa, Saldo, MoneyGram Access. Each anchor is a corridor configuration; adding a new anchor is a Helm value plus a credential pair.

### Notifications
Slack, Microsoft Teams, Discord, Email, PagerDuty for `wallet.low_balance` and `audit.suspicious_activity` events.

## Recently shipped

- **2026-04**: SCIM 2.0 provisioning for Okta and Azure AD
- **2026-04**: Audit forwarder support for S3 and HTTPS sinks
- **2026-03**: Mobile-money last-mile in Tanzania and Uganda via ClickPesa
- **2026-03**: Tax form generation beta (1099, W-2, W-8BEN)
- **2026-02**: Idempotency keys on all mutating endpoints
- **2026-01**: EURC support alongside USDC

## In flight

- Native ACH return handling for US recipients
- Mexico corridor via a new sending-anchor partnership
- Kotlin and Swift SDKs for mobile employer apps
- SOC 2 Type II evidence completion

## Stellar primitives used

- Multi-operation transactions (up to 100 ops/tx) for batched payouts
- Multi-sig accounts with weighted signers for the funding wallet
- Horizon `/transactions` and `/payments` ingestion for settlement confirmation
- Memos (`MEMO_HASH`) carrying payout reference IDs that link back to PayStream batch records
- SEP-31 cross-border anchor flow for local-currency delivery
- SEP-10 authentication for the dashboard and API key issuance
- Trustlines managed automatically per asset on the funding wallet

## Contributing

Internal contributions follow the standard fork → branch → PR flow. External PRs are welcome for SDK improvements, anchor configurations, and webhook examples. Run `make test` before pushing; `make lint` runs `golangci-lint` and `eslint` across the monorepo.

The codebase is split:

- `cmd/paystream-api` — Go API server
- `cmd/paystream-worker` — background jobs (signing, Horizon polling, webhook delivery)
- `web/dashboard` — React and TypeScript admin UI
- `sdks/{go,python,node,ruby}` — official client SDKs
- `deploy/` — Helm charts and Terraform modules

## License

PayStream is released under the Business Source License 1.1. The license converts to Apache 2.0 four years after each release. Production use up to $1M ARR is free; above that, contact licensing@paystream.dev.

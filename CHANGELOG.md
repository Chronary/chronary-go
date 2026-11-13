# Changelog

All notable changes to `chronary-go` will be documented in this file starting with the soft-launch release.

## 0.10.0 — Unreleased

- **Breaking:** `CreateCalendarParams.DefaultReminders`, `UpdateCalendarParams.DefaultReminders`, `CreateEventParams.Reminders` and `UpdateEventParams.Reminders` are now `*[]int` (were `[]int`). An empty slice used to be dropped by `omitempty`, so you could not create an event with no reminders or clear a calendar's default reminders: the event or calendar silently kept the default. Now nil omits the field, `chronary.NoReminders()` sends `[]` (no reminders) and `chronary.InheritReminders()` sends `null` (on an update, revert to inheriting the default). Migrate `Reminders: []int{10, 1440}` to `Reminders: chronary.Reminders(10, 1440)`. The `Calendar` and `Event` response fields stay `[]int` (#1991).

## 0.9.0 — 2026-11-13

- A 429 `quota_exceeded` (a monthly plan limit) is no longer retried: the call returns `ErrQuotaExceeded` at once instead of sleeping and asking again. A rate-limit 429 still retries, waiting at most 60 seconds (#1936).

## 0.8.1 — 2026-11-09

- `go.mod` adds `toolchain go1.26.10`, so the SDK's own builds and tests use a patched Go toolchain. The minimum Go version for your module is still 1.24; build with a current Go 1.26.x (or newer) toolchain to pick up the standard-library security fixes (#1642).

## 0.8.0 — 2026-10-19

- Add `BillingService` (`Billing.Get`, `CreateUpgradeRequest`, `GetUpgradeRequest`, `CancelUpgradeRequest`) so an agent can read plan status and ask its human to upgrade (#1129).
- Add `ErrPlanRequired`. `Error` now carries `UpgradeURL`, `DocsURL`, `RequiredPlan`, `Feature`, `ContactURL`, `Metric`, `Limit` and the raw `Body`.
- Add the `org.plan_changed` webhook event type: fires once when the account's plan changes, so agents don't have to poll billing (#1130).
- Add the `calendar_connection.reauthorization_required` webhook event type: fires once when a connected human calendar stops accepting access, so agents learn their availability now fails closed (#1196).
- Add `ConnectionLinks.List` and the setup details on `ConnectionLink` (`RequestingAgentID`, `RequestedCapabilities`, `RequestedPublicationPolicy`, `CreatedAt`) (#1196).
- `ListWebhookDeliveriesParams` gains `EventType` (exact, or a prefix ending in a dot such as `event.`), `CreatedAfter`, `CreatedBefore`, `DeliveryID`, `IsTest` and `Cursor`, and `WebhookDeliveryListResponse` gains `NextCursor` and `HasMore`. Add `WebhookDeliveryDropped` (a Free-plan delivery the monthly limit stopped) and `WebhookDeliveryStats.Dropped` (#1200).
- Add `Webhooks.GetDelivery`: one delivery with its attempt timeline (`AttemptLog`: HTTP status, duration, response body prefix, error code). `WebhookDelivery` gains `LastResponseStatus`, `LastErrorCode` and `LastDurationMs` (#1200).
- Add `Webhooks.SendTest`: sends a signed sample payload once (`IsTest`, `X-Chronary-Test: 1`, quota-exempt). Deliveries gain `IsTest`, and `Webhook` gains `DeactivatedAt` and `DeactivatedReason` (#1200).
- Add `Events.ListAll` with `ListOrgEventsParams` (`CalendarIDs`, `AgentID`, window, `Status`, `Source`, `Q`, `Expand`): events across all of the organization's calendars (`GET /v1/events`). Add `Calendars.RotateICalToken` to regenerate a calendar's subscribe URL, and `ListCalendarsParams.Q` (#1202).
- Add `BookingPages.ListBookings` with `ListBookingsParams` (`Limit`, `Cursor`): the bookings made through a page with the invitee's name, email and notes, newest first, paged by `NextCursor`. Add `BookingPages.RotateSlug` to replace a page's public link, and `ListBookingPagesParams.CalendarID` (#1202).
- Add `Webhooks.ReplayDelivery` (Pro): sends a `failed` or `dropped` delivery again as a new delivery, with an `X-Chronary-Replay-Of` header naming the original. It counts toward the monthly delivery quota. Add `WebhookReplaySent`. `WebhookDelivery` gains `ReplayOfDeliveryID` (#588).
- **Webhook secret rotation (#594).** `VerifySignature`, `VerifySignatureWithOptions` and `ConstructEvent` now accept a delivery when the secret matches either `X-Signature` or `X-Signature-Previous`. During a rotation's grace window Chronary sends both headers, the second signed with the replaced secret, so a receiver can switch secrets without rejecting deliveries. Update before you rotate a secret. Add `Webhooks.RotateSecret` with `RotateWebhookSecretParams.ExpireOldIn` (`WebhookSecretExpireNow`, `1h`, `24h` (default) or `7d`; Pro): returns the new secret once (`RotatedWebhookSecret`). `Webhook` gains `PreviousSecretExpiresAt`.

## 0.7.1 — 2026-07-14

- Docs-only: neutral wording for the `BookingPage` type (no competitor references). `Usage.BookingPages` already present.

## 0.7.0 — 2026-07-14

- Add `BookingPageService` (CRUD) for agent-created public scheduling links, plus a `BookingPages` usage counter (#1036).

## 0.6.0 — 2026-07-12

- Add `ConnectionLinksService` for create/get/cancel human calendar setup requests.
- Add `AllowStale` query support and provider-neutral availability health fields; correct free-slot decoding to `start`/`end`.

## 0.1.2 — 2026-05-18

- Add `CONTRIBUTING.md` to the public mirror documenting that this repo is generated from a private monorepo; PRs are welcome as proof-of-concept but can't be merged directly. No behavioral change.

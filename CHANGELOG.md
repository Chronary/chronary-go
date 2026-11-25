# Changelog

All notable changes to `chronary-go` will be documented in this file starting with the soft-launch release.

## 0.13.0 — unreleased

- **Breaking: webhook verification requires the v2 signature (#2569).** Chronary now signs every delivery twice: `X-Signature` (v1, over `<X-Timestamp>.<body>`, unchanged and deprecated) and the new `X-Signature-V2`, which also covers `X-Delivery-Id`, `X-Chronary-Event-Type`, `X-Chronary-Test` and `X-Chronary-Replay-Of`. Under v1 those headers were unsigned, so anyone holding a captured delivery could re-send it within the 5-minute tolerance under another event type (for example `event.hold_released` as `event.deleted`), with a fresh delivery id that got past deduplication, or with the test flag removed. `VerifySignature`, `VerifySignatureWithOptions` and `ConstructEvent` now check `X-Signature-V2` (or, during a secret rotation, `X-Signature-V2-Previous`) and never fall back to `X-Signature`: a delivery without `X-Signature-V2` is rejected, as is one whose event type, delivery id, test flag or replay id differs from what was signed, an `X-Chronary-Test` other than `1`, or a header value containing CR or LF. All three also require `X-Delivery-Id` and `X-Chronary-Event-Type`. Every Chronary delivery carries the v2 headers, so a receiver needs no change beyond upgrading. Function signatures are unchanged.
- **`WebhookEvent` gains `DeliveryID`, `Test` and `ReplayOf`** (`""` when the delivery is not a replay), set by `ConstructEvent`. They are covered by the signature, so deduplicate on `event.DeliveryID`.
- **Test helpers sign v2.** New `SignedDelivery`, `ComputeSignatureV2(payload, secret, timestamp, d)` and `SignedHeadersFor(payload, secret, d)`. `SignedHeaders` and `SignedHeadersWithEventType` now also set `X-Signature-V2` and a generated `X-Delivery-Id` (`whd_test_<hex>`), and sign `event.created` when no event type is given (they used to leave `X-Chronary-Event-Type` out). `ComputeSignature` still computes the v1 `X-Signature`.
- **List iterators start at `params.Offset` (#2720).** `All` and `Collect` on every list iterator (`Agents.List`, `ListCalendars`, `ListEvents`, `ListICalSubscriptions`, `ICalSubscriptions.List`, `BookingPages.List`, `Calendars.List`, `Events.List`, `Events.ListAll`, `Scheduling.List`, `Webhooks.List`) ignored the caller's `Offset` and started at row 0, so resuming an export processed rows twice. They now start at `Offset`; later pages follow as before. `GetPage(ctx, offset)` still uses its argument. `BookingPages.List` also honours `Limit` (it always asked for 50).
- **A refusal whose `Retry-After` is beyond 60 seconds is not retried (#2558).** The client clamped such a `Retry-After` to 60 s and retried, so an iCal manual-sync cooldown (429, up to 900 s on Free) or the anonymous sign-up daily cap (503, seconds to UTC midnight) blocked the call for about two minutes (or hit the caller's deadline) and then returned the same error. These now return at once; a `Retry-After` of 60 s or less is still honoured and retried. `Error.RetryAfter` is now also set from an HTTP-date `Retry-After` (seconds until then), and a very large value no longer overflows the wait.
- **`All` and `Collect` follow `next_cursor` (#2691).** `Agents.List`, `Calendars.List`, `Agents.ListCalendars` and `Events.ListAll` paged by offset only, so iterating past 10,000 rows will fail once the API enforces its `offset` cap. They now page by the returned `next_cursor` (other lists still page by offset). `ListResponse` and `Page` gain `NextCursor`; `ListAgentsParams`, `ListCalendarsParams` and `ListOrgEventsParams` gain `Cursor`; new `PageIterator.GetPageAfter(ctx, cursor)`. A list that returns the cursor it was just sent ends the iteration with an error instead of looping.

## 0.12.2 — 2026-11-22

- **Path parameters are escaped (#2519).** Every ID a method puts into the URL path now goes through `url.PathEscape` as one segment, so a value such as `../agents/agt_1`, `whk_1#` or `evt_1?expand=all` can no longer turn `Events.DeleteByID` or `Events.GetByID` into a call on another resource or add a query string. `""`, `"."` and `".."` are refused before any request is sent: the method returns an `*Error` with `Type: "validation_error"` and `StatusCode: 0` (`errors.Is(err, chronary.ErrValidation)` holds). Real Chronary IDs are unaffected.

## 0.12.1 — 2026-11-22

- `ProposalResponse` gains `AutoResolved` and `AutoResolveError` (`*AutoResolveError{Code, Message}`), set only on the `Scheduling.Respond` result. When the last response's auto-resolve fails, the API now answers 200 with `auto_resolve_error` instead of an error status, because the response was saved; `Message` says what to do next (usually `Scheduling.Resolve`). Don't respond again (#2540).

## 0.12.0 — 2026-11-17

- `Usage` gains `CalendarConnections` (`UsageCounter`): live Google / Microsoft calendar connections vs. the plan cap (#2397).
- Doc comments name the Developer plan: `AuditLog.List` retention is Free 3 days / Developer 30 / Pro 90, and `Billing.CreateUpgradeRequest` asks for the next plan up (Free to Developer, Developer to Pro) rather than always Pro (#2397).

## 0.11.0 — 2026-11-16

- **Time-zone-aware recurrence (#2276).** `CreateEventParams.Timezone` (`*string`, an IANA zone such as `chronary.String("America/Los_Angeles")`) and `UpdateEventParams.Timezone`: a recurring series then repeats on local wall time across daylight-saving changes, and BYDAY and a date-only UNTIL are read in that zone. `UpdateEventParams.Timezone` is `**string`: nil leaves the zone unchanged, `chronary.SetTimezone("America/Los_Angeles")` sets it and `chronary.ClearTimezone()` sends `null` (back to UTC). `Event` gains `Timezone`, and expanded instances gain `StartLocal`, `UTCOffset` and `DSTAdjusted` (`"gap"` when a DST gap moved the start forward).
- **Series edits keep cancelled occurrences (#2277).** An update that changes `StartTime`, `RecurrenceRule` or `Timezone` on a recurring series carries each cancelled occurrence to the new occurrence on the same local date instead of clearing it. `Event.DroppedExdates` lists any with no occurrence left on that date; re-cancel them with `CancelOccurrence` if needed.

## 0.10.0 — 2026-11-13

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

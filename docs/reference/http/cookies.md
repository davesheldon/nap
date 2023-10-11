---
layout: default
title: Cookies
nav_order: 5
parent: HTTP
grand_parent: Reference
permalink: /reference/http/cookies
---

{: .fs-10 .fw-300 }
# Cookies

{: .fs-6 .fw-300 }
Cookies allow session data to be carried between requests in Nap.

## Response Cookies

Response cookies are stored and sent in later responses in the same routine. Because of this chaining technique, cookies "just work" in Nap for most use cases.

## Sending Cookies Manually

Cookies can be created sent manually during a request. These cookies are not stored unless they're returned in the response.

{: .highlight }
For information on setting cookies on a request, see [File Types -> Requests](/reference/file-types/requests#cookies---http-request-cookies).

## In Queries

Cookies can be accessed via queries, which makes them available in both captures and asserts.

{: .highlight }
For information on cookies in queries, see [Concepts -> Queries](/reference/concepts/queries).
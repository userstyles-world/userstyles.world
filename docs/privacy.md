# GDPR Privacy Policy of UserStyles.world

Last updated June 06, 2021

Table of Content

- [Introduction](#introduction)
- [What data do we collect?](#what-data-do-we-collect)
- [How do we collect your data?](#how-do-we-collect-your-data)
- [How do we use your data?](#how-do-we-use-your-data)
- [How do we store your data?](#how-do-we-store-your-data)
- [What are your data protection rights?](#what-are-your-data-protection-rights)
- [What log data do we collect?](#what-log-data-do-we-collect)
- [Do we use cookies?](#do-we-use-cookies)
- [Do we use any external cookies?](#do-we-use-any-external-cookies)
- [Privacy policies of other websites](#privacy-policies-of-other-websites)
- [Changes to our privacy policy](#changes-to-our-privacy-policy)
- [How to contact us](#how-to-contact-us)


## Introduction

Thank you for choosing to be part of our community at UserStyles.world (“we”,
“us”, or “our”).

We are committed to protecting your personal information and your right to
privacy. If you have any questions or concerns about our policy, or our
practices with regards to your personal information, please contact us at
[feedback@userstyles.world](mailto:feedback@userstyles.world).

When you visit our website userstyles.world (“Site”) and use our services, you
trust us with your personal information. We take your privacy very seriously.

In this privacy notice, we describe our privacy policy. We seek to explain to
you in the clearest way possible what information we store, how we use it and
what rights you have in relation to it. We hope you take some time to read
through it carefully, as it is important.

If there are any terms in this privacy policy that you do not agree with, please
discontinue use of our site and our services.

This privacy policy applies to all information stored through our Site, and/or
any related services. Please read it carefully as it will help you make informed
decisions about sharing your personal information with us.


## What data do we collect?

We collect the following data:

- Username
- Email address
- Visited userstyle pages
- Visited userstyle install pages


## How do we collect your data?

All of the data is directly provided by your activity on the Site. The data is
processed when you:

- Register a new account
- Visit userstyle details page (e.g. [https://userstyles.world/style/1/userstyles-world-tweaks][details])
- Visit userstyle install page (e.g. [https://userstyles.world/api/style/1.user.css][api])

[details]: https://userstyles.world/style/1/userstyles-world-tweaks
[api]: https://userstyles.world/api/style/1.user.css


## How do we use your data?

We use your data in order to show userstyle statistics.


## How do we store your data?

We follow best practices when it comes to handling your data and store
everything on our own server.

The data used for userstyle statistics is stored as a unique hash, and is not
easily reversible without brute-forcing all public IP addresses in IPv4 address
space in combination with our secret key. This gives us decently accurate style
statistics while respecting your privacy. The unique hash is formed like so:

```pseudo
# Formula:
record = IP + " " + StyleID
secret = SecretKey (string converted to bytes)
hashed = HashFunction(record, secret)

# Example:
record = "1.2.3.4 1"
secret = 73 33 63 72 33 37 6b 33 79 (s3cr37k3y)
hashed = HMAC-SHA512(record, secret)

# Result:
97a5abba601dd18829c33507ecf295bf0d2c05db06a3bf8af2c091dee0a8556500886443b59076057ffc5d8ad429d3d1de141e58684740729f3f24c7c435f7bb
```

Try it out online:

- [HMAC hash generator](https://cryptii.com/pipes/hmac)
- [Convert String to Bytes](https://onlinestringtools.com/convert-string-to-bytes)


## What are your data protection rights?

TODO: Add GDPR data protection rights.


## What log data do we collect?

We automatically store certain information when you use the Site, but this
information doesn't reveal your identity. It comes from default NGINX server
logs, which we flush every 24 hours, and includes the following data.

```
# IP        Date                         Visited page                  Browser Agent
1.2.3.4 - - [06/Jun/2021:23:06:13 +0000] "GET / HTTP/1.1" 200 5217 "-" "Mozilla/5.0 ..."
```

This information is necessary for maintaining the security and operation of
our Site, and for our internal analytic and reporting purposes.


## Do we use cookies?

Yes, we use cookies for the purpose of keeping you signed in and authorizing
various actions on the Site that require an account.


## Do we use any external cookies?

No. And we never will.


## Privacy policies of other websites

Our Site contains links to external websites.

This privacy policy only applies to our website, so if you click on a link to
another website, you should read their privacy policy.


## Changes to our privacy policy

The Site is always in development and we will adjust our privacy policy to
reflect new changes when necessary and without any prior notice. The "last
updated" date will always be visible at the top, beneath the page title.


## How to contact us 

If you have any questions about our privacy policy, the data we hold on you, or
you would like to exercise one of your data protection rights, please don't
hesitate to contact us via email:
[feedback@userstyles.world](mailto:feedback@userstyles.world)

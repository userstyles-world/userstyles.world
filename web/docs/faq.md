---
Title: Frequently Asked Questions
---

# Frequently Asked Questions

Last updated May 17, 2023

<!-- markdown-toc start - Don't edit this section. -->
**Table of Contents**

- [General](#general)
    - [What is our monetization strategy?](#what-is-our-monetization-strategy)
    - [Why isn't \[insert here\] implemented?](#why-isnt-insert-here-implemented)
    - [How can I contact the admins?](#how-can-i-contact-the-admins)
- [Userstyles](#userstyles)
    - [Why are `@updateURL` fields overriden?](#why-are-updateurl-fields-overriden)
    - [Why are ratings different in Stylus' search?](#why-are-ratings-different-in-stylus-search)
    - [How do view/install/update statistics work?](#how-do-viewinstallupdate-statistics-work)
    - [How do I remove the `Get Stylus` button?](#how-do-i-remove-the-get-stylus-button)
    - [Why is mirroring source code updates not working?](#why-is-mirroring-source-code-updates-not-working)
    - [Why is there no support for traditional userstyles?](#why-is-there-no-support-for-traditional-userstyles)
    - ["Bad style format" error](#bad-style-format-error)
    - [How does mirroring source code work?](#how-does-mirroring-source-code-work)

<!-- markdown-toc end -->

## General

Non-specific questions.


### What is our monetization strategy?

Since we don't have any means to make money, we rely on donations. Currently, we
fund everything from our own pockets: server (including backup) cost is $6 per
month (or $72 per year) and domain cost is $22 per year (as of April, 2022). You
can support the project on [Open Collective][oc] and help us stay afloat.

[oc]: https://opencollective.com/userstyles


### Why isn't [insert here] implemented?

Well, it could be a variety of reasons. Contact us directly, or open a new issue
over at [the issue tracker][issues] and let's discuss it first.

[issues]: https://github.com/userstyles-world/userstyles.world/issues


### How can I contact the admins?

You could join us over on [Discord], [Matrix], or send us an [email]. If you
contacted us via email and didn't get a reply within a couple of days, please
join us on either chat platform and let's talk about it.

[Discord]: https://userstyles.world/link/discord
[Matrix]: https://userstyles.world/link/matrix
[email]: mailto:feedback@userstyles.world


## Userstyles

Questions regarding userstyles.


### Why are `@updateURL` fields overriden?

It's done in order to avoid the possibility of tracking, as well as broken URLs.


### Why are ratings different in Stylus' search?

We have to accommodate for how Stylus displays them due to compatibility with
USo: their ratings are on a scale of 1 to 3 (bad, okay, and good), meanwhile our
ratings are on a scale of 1 to 5. That results in having to fit them within
those bounds by multiplying the average rating by 3 then dividing it by 5 — or
multiplying by 0.6 — so that Stylus can display them using appropriate icons.


### How do view/install/update statistics work?

As of July 2021, statistics work like so:

- The view counter increases when the user visits userstyle details page
(`/style/:id/:name`).
- The install counter increases when the user visits userstyle install page
(`/api/style/:id.user.css`).
- Update counter is calculated based on install statistics and when it was
  created/updated in database.

That's the reason why some styles have more installs than views. You don't
necessary need to visit a style page to have it installed, e.g. third-party
applications can directly install any style (e.g. Stylus' inline search).


### How do I remove the `Get Stylus` button?

[Stylus extension] removes it automatically from `v1.5.18`, or you can enable an
option provided by [UserStyles.world Tweaks] userstyle.

[Stylus extension]: https://github.com/openstyles/stylus
[UserStyles.world Tweaks]: https://userstyles.world/style/1/userstyles-world-tweaks


### Why is mirroring source code updates not working?

First of all, make sure the checkbox "Mirror source code updates" is enabled. If
it isn't enabled, enable it on the edit page, then save changes.

If you're mirroring source code from a different URL than your userstyle was
originally imported from, make sure that the new URL location is correct. It's
correct if Stylus' install/reinstall page shows up when you visit it.

Last, but certainly not least, make sure that you increase `@version` field in
the UserStyle metadata header. The code is mirrored only if the new version
doesn't match the one stored in our database, otherwise our updater ignores it.


### Why is there no support for traditional userstyles?

Traditional userstyles don't fit in the current workflow because they can't be
self-hosted like UserCSS userstyles. That means installing and/or updating any
userstyle would be a manual process, which defeats the purpose of having USw.

Converting your traditional userstyle is as simple as exporting it in Mozilla
format and using the mandatory UserStyle metadata header that's provided for you
on the "Add userstyle" page. Documentation is on Stylus' [Writing UserCSS page].

[Writing UserCSS page]: https://github.com/openstyles/stylus/wiki/Writing-UserCSS


### "Bad style format" error

Currently, there is an issue with Stylus integration that allows users to add
broken userstyles because `@-moz-document` fields are not included in the
process of adding styles written in "traditional" format. The issue causes
incorrectly formatted styles to apply globally (in other words, on all sites).

In order to prevent addition of more broken userstyles (roughly 25% of them uses
incorrect format), we have decided to not add them unless they pass a specific
criteria: all newly created userstyles need to have `@-moz-document` fields.

To fix your userstyle, do the following:

1. Open it in the editor
1. Click on "Export" button to get code in Mozilla format
1. Copy the source code to your clipboard (we'll create a new style)
1. Click on "Back to manage" button in order to enable "Usercss" format
1. Toggle on "as Usercss" checkbox (☑), then click on "Write new style" button
1. Paste your code below Stylus' default Usercss template
1. Remove the default `@-moz-document domain("example.com") { ... }` block
1. Edit UserStyle header with your info, links, etc.
1. Finally, click on "Publish style" button

The resulting style should look along the lines of:

```css
/* ==UserStyle==
@name           Test style name
@namespace      userstyles.world
@version        1.0.0
==/UserStyle== */

@-moz-document domain("userstyles.world") {
    * {
        color: crimson;
    }
}

@-moz-document url-prefix("https://userstyles.world/docs/") {
    h1, h2, h3 {
        color: tan;
    }
}
```

Please test your style to see whether it works after you publish it to USw. All
broken styles will be removed. If you have issues, feel free to contact us via
any of the links in the page footer or via our feedback email address.

P.S. If your style does apply globally, for the time being you'll have to wrap
it in the following:

```css
@-moz-document regexp(".*") {
    /* Your code goes here. */
}

/* You can also use the following. */
@-moz-document url-prefix("http") {
    /* Your code goes here. */
}
```


### How does mirroring source code work?

Every 4th minute of every 4 hours (00:04 UTC+0, 04:04, and so on) we check for
updates. Userstyle checking runs in batches of 25 (as of August 27, 2021), so it
can take up to a few minutes for your userstyles to be processed. Code will be
updated if `@version` field doesn't match the one in our database.

If your userstyle isn't being updated, read through [troubleshooting
steps](#why-is-mirroring-source-code-updates-not-working) first.

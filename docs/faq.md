# FAQ

Last updated August 14, 2021

<!-- markdown-toc start - Don't edit this section. -->
**Table of Contents**

- [General](#general)
    - [What is monetization strategy?](#what-is-monetization-strategy)
    - [Why isn't \[insert here\] implemented?](#why-isn-t-insert-here-implemented)
    - [How can I contact the admins?](#how-can-i-contact-the-admins)
- [Userstyles](#userstyles)
    - [Why are `@updateURL` fields overriden?](#why-are-updateurl-fields-overriden)
    - [How do view/install/update statistics work?](#how-do-view-install-update-statistics-work)
    - [How do I remove the `Get Stylus` button?](#how-do-i-remove-the-get-stylus-button)
    - [Why is mirroring source code updates not working?](#why-is-mirroring-source-code-updates-not-working)
    - [Why is there no support for traditional userstyles?](#why-is-there-no-support-for-traditional-userstyles)
    - ["Bad style format" error](#bad-style-format-error)

<!-- markdown-toc end -->

## General

Non-specific questions.


### What is monetization strategy?

There is none. USw is funded from our own pockets.

All donations to Vednoc's [LiberaPay] and/or [Ko-fi] go towards paying for the
server costs as well as domain costs. In terms of numbers, the server costs $5
per month (as of July 2021), and domain costs $22 per year.

[LiberaPay]: https://liberapay.com/vednoc
[Ko-fi]: https://ko-fi.com/vednoc


### Why isn't [insert here] implemented?

Well, it could be for various reasons. Contact us directly, or open a new issue
over at [the issue tracker] and lets discuss it.

[the issue tracker]: https://github.com/userstyles-world/userstyles.world/issues/new/choose


### How can I contact the admins?

You could join us over at the [Discord server], however you can also reach over
the email at [feedback@userstyles.world].

[Discord server]: https://discord.gg/WW6vnFsCpB
[feedback@userstyles.world]: mailto:feedback@userstyles.world


## Userstyles

Questions regarding userstyles.


### Why are `@updateURL` fields overriden?

It's done in order to avoid the possibility of tracking, as well as broken URLs.


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
applications can directly install any style.


### How do I remove the `Get Stylus` button?

[Stylus extension] removes it automatically from `v1.5.18`, or you can enable an
option provided by [UserStyles.world Tweaks] userstyle.

[Stylus extension]: https://github.com/openstyles/stylus
[UserStyles.world Tweaks]: https://userstyles.world/style/1/userstyles-world-tweaks


### Why is mirroring source code updates not working?

First of all, make sure the checkbox "Mirror source code updates" is enabled. If
it isn't enabled, enable it then save changes.

If you're mirroring source code from a different URL than your userstyle was
originally imported from, make sure that the new URL location is correct. It's
correct if Stylus' install/reinstall page shows up when you visit it.

Last but certainly not least, make sure that you increase `@version` field in
the UserStyle metadata header. New userstyle will be mirrored if the new version
doesn't match the one in the database.


### Why is there no support for traditional userstyles?

Traditional userstyles don't fit in the current workflow because they can't be
self-hosted like UserCSS userstyles. That means installing and/or updating any
userstyle would be a manual process, which defeats the purpose of USw.

Converting your traditional userstyle is as simple as exporting it in Mozilla
format and using the mandatory UserStyle metadata header that's provided for you
on the "Add userstyle" page. Documentation is on Stylus' [Writing UserCSS page].

[Writing UserCSS page]: https://github.com/openstyles/stylus/wiki/Writing-UserCSS


### "Bad style format" error

Currently, there is an issue with Stylus integration that allows users to add
broken userstyles because `@-moz-document` fields are not included in the
process of adding styles written in "traditional" format. The issue causes
incorrectly formatted styles to apply globally (in other words, on all styles).

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

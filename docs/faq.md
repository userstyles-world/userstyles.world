# FAQ


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

And because of that it can happen, a style have more installations than views,
because you don't necessary need to visit a style to have it being installed,
e.g. third-party applications directly installing the style.


### How do I remove the `Get Stylus` button?

[Stylus extension] removes it automatically from `v1.5.18`, or can enable an
option provided by [UserStyles.world Tweaks] userstyle.

[Stylus extension]: https://github.com/openstyles/stylus
[UserStyles.world Tweaks]: https://userstyles.world/style/1/userstyles-world-tweaks


### Why is there no support for traditional userstyles?

Traditional userstyles don't fit in the current workflow because they can't be
self-hosted like UserCSS userstyles. That means installing and/or updating any
userstyle would be a manual process, which defeats the purpose of USw.

Converting your traditional userstyle is as simple as exporting it in Mozilla
format and using the mandatory UserStyle metadata header that's provided for you
on the "Add userstyle" page. Documentation is on Stylus' [Writing UserCSS page].

[Writing UserCSS page]: https://github.com/openstyles/stylus/wiki/Writing-UserCSS

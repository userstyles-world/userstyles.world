---
title: Documentation
type: page
desc: Work in progress
---

## Contributing

If you'd like to contribute in any way --- be it writing documentation,
submitting new userstyles, reporting and/or fixing bugs, suggesting ideas,
letting userstyle creators know about this place, or really anything else ---
it will be very appreciated.


## Adding a new userstyle

This is a short guide on how to convert your userstyle to UserCSS format, how to
host it on GitHub (or somewhere else), and adding it to USw.


### Setting up Stylus

The initial configuration and migration to UserCSS format.

1. You need to go to Stylus _Manage page_ by clicking on the Stylus icon on your
   toolbar, then on the Manage button. Once you're there, you'll have to enable
   checkbox next to "Write new style" that says "as UserCSS"
2. With that set, you're ready to click on "Write new style" and you'll be
   presented with a default UserCSS template --- do not delete anything, but you
   should replace default values to match your userstyle
4. With that done, you're ready to add your code. There are two ways:
   1. (Fast way) Go to your old userstyle, then press "Export" button below
      "Mozilla Format" heading in the sidebar. Copy the code over to your new
      userstyle and paste it at the very end. Remove the default `@-moz-document
      domain("example.com")` because it's not necessary.
   2. (Slow way) Go to your old userstyle, copy everything and put it inside of
      `@-moz-document` in your new userstyle. Make sure to change the default
      `domain("example.com")` to match your site — see "Applies to" area in your
      old userstyle for reference.
5. (Bonus) See [Writing-UserCSS][] on Stylus' wiki for more information. There
   are lots of ways to customize the userstyle as well as lots of other options

[Writing-UserCSS]: https://github.com/openstyles/stylus/wiki/Writing-UserCSS


### Setting up a repository

This will be for GitHub, but the same principles apply to other git platforms,
with one exception being that you'll need to make a GitHub account to send the
Pull Request, or send it to me via e-mail.

1. Initialize a new repository
2. Add a readme that includes some basic information and a screenshot
   - Adding an [install badge][] to readme would be great!
3. Make a new file and copy your new userstyle over. Make sure to name it
   something along the lines of `example.user.css`, with `.user.css` being
   required otherwise Stylus won't be able to recognize it as a UserCSS
   userstyle.
4. Commit it when you're done
5. Add a license to your project so that others can contribute
6. Test in a separate profile/browser to confirm that your userstyle is working

[install badge]: https://github.com/openstyles/stylus/wiki/Writing-UserCSS#badges


### Sending a Pull Request

The final step before your userstyle can be seen on USw. Have a look at the
templates for other userstyles --- [Dark-WhatsApp][] is the best one — and view
them as "raw" files to avoid pretty formatting on GitHub.

1. Go to [USw repository][] and click on "Fork" button in the top right corner
2. Use a template to add information about your userstyle:
    1. Have a look at userstyles in `content/userstyle` directory for examples
    2. Copy one of them over, or use `archetypes/userstyle.md` as base template
    3. Fill in the fields
3. Follow [Creating a Pull Request][] article
4. That should be it! I'll review it and let you know if everything is alright

[USw repository]: https://github.com/vednoc/userstyles.world
[Creating a Pull Request]: https://docs.github.com/en/free-pro-team@latest/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request#creating-the-pull-request
[Dark-WhatsApp]: https://github.com/vednoc/userstyles.world/blob/main/content/userstyle/dark-whatsapp.md


## Working with Hugo

TODO Add how to contribute Hugo code.

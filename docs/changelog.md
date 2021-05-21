# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [1.1.0](https://github.com/userstyles-world/userstyles.world/compare/v1.0.0...v1.1.0) (2021-05-21)

### Features

* **build:** exclude templates and scss from air ([48b388b](https://github.com/userstyles-world/userstyles.world/commit/48b388bf2e849d68c83c5f242d9a77ed903bd3f8))
* **css:** enable color-scheme meta in Chromium ([9953239](https://github.com/userstyles-world/userstyles.world/commit/995323943105c599b1b986b297af5b78449f6d19)), closes [#59](https://github.com/userstyles-world/userstyles.world/issues/59)
* **css:** improve responsive design for bars ([c496da7](https://github.com/userstyles-world/userstyles.world/commit/c496da750f8c598263983997c5e235d505b387a3))
* **css:** improve styles for card images ([b79ee27](https://github.com/userstyles-world/userstyles.world/commit/b79ee27300816ee22e327ce64869db44942a2372))
* **css:** improve the look of style cards ([2badb7a](https://github.com/userstyles-world/userstyles.world/commit/2badb7ac15142cec533282ee8d240703f2dbb0bf))
* **html:** add confirm page for style deletion ([3046537](https://github.com/userstyles-world/userstyles.world/commit/30465376450d3f14eb308f28c442eea43fa44ce5))
* **html:** add custom checkboxes on sign in page ([b56728d](https://github.com/userstyles-world/userstyles.world/commit/b56728d0eccca0b5b9974cbb32732bcd8a3a5b2a))
* **html:** add new checkboxes on style edit page ([8067200](https://github.com/userstyles-world/userstyles.world/commit/8067200f4a21a3b2ad3eec7b2abfb8e38ceaa8e8))
* **html:** move updated date to card footer ([fd2cc79](https://github.com/userstyles-world/userstyles.world/commit/fd2cc79bd077052f5a61de290bdba407b372c7cb))
* **html:** truncate source code from large styles ([15e7620](https://github.com/userstyles-world/userstyles.world/commit/15e7620669ee917b12194ca3a87d4c8dc2b70209))
* **ts:** add user settings + programmatic color-scheme ([37be14a](https://github.com/userstyles-world/userstyles.world/commit/37be14a12c4f8643f21cb2fa0627a1b7d4ac32ac))
* **ts:** change color-scheme meta ([1254d50](https://github.com/userstyles-world/userstyles.world/commit/1254d504e0be4119222c21f5c81eb31ed4c38995))
* **ts:** add ESLint ([194a320](https://github.com/userstyles-world/userstyles.world/commit/194a320b5acfa2b2e77f247a87de0ea7acb590ae))
* **ts:** add typescript workflow ([12593be](https://github.com/userstyles-world/userstyles.world/commit/12593be806a27ee57f5f0feac15282129fe1fdde))


### Bug Fixes

* **css:** improve auto fill colors across browsers ([5e255c5](https://github.com/userstyles-world/userstyles.world/commit/5e255c5cb23bddf599ba3eefaf40fcad764caaaf))
* **css:** improve max-width for search in nav menu ([c21d246](https://github.com/userstyles-world/userstyles.world/commit/c21d246cc9704a8299bddabdbdac9b4b68e54047))
* **css:** improve position for short screenshots ([9b9f916](https://github.com/userstyles-world/userstyles.world/commit/9b9f916e84c1c3335994c096b2f5c4dc8cd58eeb))
* **css:** truncate author names in style cards ([e98b41f](https://github.com/userstyles-world/userstyles.world/commit/e98b41fb4cab320e9a61b2808f958ca4b654bb1e))
* **html:** hide unset display name fields ([f8cd137](https://github.com/userstyles-world/userstyles.world/commit/f8cd13787409ad9f2d105eea97d2021cb6ec0007))
* **html:** remove blurred preview image from cards ([928467d](https://github.com/userstyles-world/userstyles.world/commit/928467dd1f9eacfeef057f372eadcfd0e5b3b67e))
* **html:** remove useless element ([c8c1ab7](https://github.com/userstyles-world/userstyles.world/commit/c8c1ab7b3c0c12c6e091488ac9e001100bd78593))
* **html:** resolve bad text alignment in buttons ([9e50327](https://github.com/userstyles-world/userstyles.world/commit/9e5032708ff7d422fbc61de92dad3a545c651e1a))
* **models:** add missing methods to MinimalStyle ([dd0d748](https://github.com/userstyles-world/userstyles.world/commit/dd0d748f437e8319d5b43389a86a0d17696095b2))
* **tools:** don't watch for data folder ([2c089e2](https://github.com/userstyles-world/userstyles.world/commit/2c089e2e5c6d31d5eb9dfcb894a5114f683089db))
* **ts:** yoda-compatible ([224c3d4](https://github.com/userstyles-world/userstyles.world/commit/224c3d49c364f6445deaa77cf862ad148305ec4a))
* **ts:** don't error that esbuild isn't used ([85f616c](https://github.com/userstyles-world/userstyles.world/commit/85f616c5c054e3a6146da22d611d6dee0743ea38))

## [1.0.0](https://github.com/userstyles-world/userstyles.world/compare/fef3eb9...v1.0.0) (2021-05-16)

### Features

* add an option to set custom style mirror URL b8ee98c
* add base functionality for display names cd502d6
* add base style statistics to style cards d0b8fb0
* add changelog links and template 91c053d
* add form for user display names fa4fba2
* add HSTS a3cb8e6
* add input validation to display name form dba3d4a
* add legal into server 3981c47
* add legal into server 0335d41
* add more info to account/profile pages f24a5bd
* add privacy policy ea91d75
* add privacy policy 50f181d
* add server-side validation for display names c74226a
* add ToS 0c958bc
* add ToS 2c2f696
* cache up to 2 weeks c5db4c6
* improve accuracy of style stats on home page fa8e8f3
* improve colors for input areas and buttons d50ff19
* improve form validation for display names 8d85f39
* improve styles for page footer template 19cafcf
* improve styles of non 16:9 preview images 040f909
* show stats for weekly updates on home page 232c657
* show weekly update statistics on style page afa236b
* truncate long style names in userstyle cards 43e6d8c
* tweak background colors for dark/light mode dbdc2bd
* use Stylus cyan color for our gopher mascots 89cdd18

### Bug Fixes

* add missing icons used by style cards 99f1384
* allow more characters 38794e0
* allow parallel test 0412562
* change email ae03616
* change email 490ef83
* check proper values for social links c1a312e
* don't count new installs as updated stats 3654cf0
* improve accuracy for statistics on style page 39849bb
* improve contrast ratio (3 -> 8.3) f57b58d
* improve Lighthouse score for accessibility 62bf063
* include logic for blurred preview images 461ec76
* increase width of style stats and share field 7c0c4dc
* provide good search results 2875bf8
* remove required attr from mirror URL field 252c9b4
* show user state on docs route 9bec8b5
* skip showing site stats if user is logged in 9461eb1
* test cases 2790679
* tweak code selectors in Markdown areas f0bbfc6
* use correct package name 713c2da
* fix formatting for API's index endpoint 7ba7d51

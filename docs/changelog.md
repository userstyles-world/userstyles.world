# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

<!-- markdown-toc start - Don't edit this section. -->
**Table of Contents**

- [v1.6.1](#1-6-1-https-github-com-userstyles-world-userstyles-world-compare-v1-6-0-v1-6-1-2021-08-12)
- [v1.6.0](#1-6-0-https-github-com-userstyles-world-userstyles-world-compare-v1-5-0-v1-6-0-2021-08-11)
- [v1.5.0](#1-5-0-https-github-com-userstyles-world-userstyles-world-compare-v1-4-0-v1-5-0-2021-07-18)
- [v1.4.0](#1-4-0-https-github-com-userstyles-world-userstyles-world-compare-v1-3-0-v1-4-0-2021-07-12)
- [v1.3.0](#1-3-0-https-github-com-userstyles-world-userstyles-world-compare-v1-2-1-v1-3-0-2021-06-19)
- [v1.2.1](#1-2-1-https-github-com-userstyles-world-userstyles-world-compare-v1-2-0-v1-2-1-2021-05-30)
- [v1.2.0](#1-2-0-https-github-com-userstyles-world-userstyles-world-compare-v1-1-0-v1-2-0-2021-05-29)
- [v1.1.0](#1-1-0-https-github-com-userstyles-world-userstyles-world-compare-v1-0-0-v1-1-0-2021-05-21)
- [v1.0.0](#1-0-0-https-github-com-userstyles-world-userstyles-world-compare-fef3eb9-v1-0-0-2021-05-16)

<!-- markdown-toc end -->

## [1.6.1](https://github.com/userstyles-world/userstyles.world/compare/v1.6.0...v1.6.1) (2021-08-12)


### Bug Fixes


* **db:** revert duplicate style conditional ([02a8101](https://github.com/userstyles-world/userstyles.world/commit/02a81019643a1d309c2e46be69baddc2b8c42a10))

## [1.6.0](https://github.com/userstyles-world/userstyles.world/compare/v1.5.0...v1.6.0) (2021-08-11)


### Features

* **api:** add in-memory cache for index endpoint ([1242ffd](https://github.com/userstyles-world/userstyles.world/commit/1242ffd9b49f988b3c3085f13bfe5d2d0bc4be86))
* **api:** cache USo-format JSON to disk ([2842895](https://github.com/userstyles-world/userstyles.world/commit/2842895cc01620bf3a441f6dd648b4e89772d295))
* **api:** optimize caching USo-format index ([16ad0a1](https://github.com/userstyles-world/userstyles.world/commit/16ad0a1f6b4f878634bc6d8165ae3e44c66820d9))
* **api:** speed up sorting on Explore page ([344168b](https://github.com/userstyles-world/userstyles.world/commit/344168bb1c0234f400c1253b8fe5a2409b2b87b6))
* **api:** speed up USo-format endpoint ([7735365](https://github.com/userstyles-world/userstyles.world/commit/7735365c9810c77425029bc15ff7cf1e6acd4db4))
* **css:** improve colors used for buttons ([eeefef3](https://github.com/userstyles-world/userstyles.world/commit/eeefef3968b44bf9e7e563703c829bd928e7e977))
* **css:** improve focus states for elements ([4da4081](https://github.com/userstyles-world/userstyles.world/commit/4da4081cb54527da5cd7d00065a8e759736adcdf))
* **dashboard:** add humanized values to stats ([2a14207](https://github.com/userstyles-world/userstyles.world/commit/2a14207aa07d709877cd1a9130d20b91f35d8963))
* **dashboard:** add Style history bar chart ([b66f3bb](https://github.com/userstyles-world/userstyles.world/commit/b66f3bbdaad684235d0c244dce604d1b5cc80148))
* **dashboard:** make dates/numbers more readable ([85f788a](https://github.com/userstyles-world/userstyles.world/commit/85f788a9c770ff2fa1713c85930ba806cf948212))
* **dashboard:** show data for system status ([c3bc5b8](https://github.com/userstyles-world/userstyles.world/commit/c3bc5b8d66b1ec255b9fd80f8add52b7516e82e7))
* **dashboard:** speed up User and Style overviews ([cc026a3](https://github.com/userstyles-world/userstyles.world/commit/cc026a368b23c79f29c562ded4f86f476dfb01e4))
* **db:** speed up sorting styles on Explore page ([24f2bf4](https://github.com/userstyles-world/userstyles.world/commit/24f2bf4b2f9de2d6f0836ae8428ecf878243568b))
* **docs:** add Kind Communications Guidelines ([b1c7230](https://github.com/userstyles-world/userstyles.world/commit/b1c7230d55729376aef0d58a6b83b5ff0356f5b4)), closes [#69](https://github.com/userstyles-world/userstyles.world/issues/69)
* **docs:** add ToC to all on-site documentation ([683a2a8](https://github.com/userstyles-world/userstyles.world/commit/683a2a8d635381f2a421b9b0046dd9dc030a6082))
* **docs:** document mirroring source code updates ([4a94aa0](https://github.com/userstyles-world/userstyles.world/commit/4a94aa0c75f806d17c5266db1020b35a68efd3a5))
* **docs:** fix typo in crypto file. ([4ccdd57](https://github.com/userstyles-world/userstyles.world/commit/4ccdd5715078a121ae408150b826d18638b390a1))
* **docs:** revise wording/typos in crypto file ([71200d1](https://github.com/userstyles-world/userstyles.world/commit/71200d1be22e6c3ebe8bebf1e7f3a5c7f2dc807f))
* **handlers:** Add CSP headers ([fe3da38](https://github.com/userstyles-world/userstyles.world/commit/fe3da38a7ef9811f0b255ddc40331ed5825ce694))
* **handlers:** add more CSP header values ([9d45a55](https://github.com/userstyles-world/userstyles.world/commit/9d45a5530a1d64472dc8ad98bd7fddf845ded10a))
* **handlers:** add route to shorten external URLs ([7b2405c](https://github.com/userstyles-world/userstyles.world/commit/7b2405c7ee58134b5e8f847fa4af2fa2c83fb3d2))
* **handlers:** add security-policy redirect ([073fcaa](https://github.com/userstyles-world/userstyles.world/commit/073fcaa8676415f6b239c66caf5c1ed10212ab66))
* **html:** add release code name to page footer ([a3958e7](https://github.com/userstyles-world/userstyles.world/commit/a3958e733c46c5a1cabce6bd5bf2d8f4528410ca))
* **html:** improve functionality for pagination ([7c9766d](https://github.com/userstyles-world/userstyles.world/commit/7c9766dde3babfb77004f880ebf63d075eefc8c3))
* **html:** make 'Write a review' button public ([ae40992](https://github.com/userstyles-world/userstyles.world/commit/ae40992d8c834ee65cba82c3f855a4c9db27fd54))
* **images:** bump cache to 3 months ([aeb9fe9](https://github.com/userstyles-world/userstyles.world/commit/aeb9fe91345e8c88c70e1c0c6588ced922634662))
* **log:** implement file-based logging ([ff2458c](https://github.com/userstyles-world/userstyles.world/commit/ff2458c213b4d60cd22496aa354934d7b65a146a))
* **oauthP:** allow vendor_data to be passed ([59b2eb8](https://github.com/userstyles-world/userstyles.world/commit/59b2eb8f3afec8b8cf8efc92c7b8d15d1e3b6fe9))
* **oauthP:** use provided name if possible ([983c2ef](https://github.com/userstyles-world/userstyles.world/commit/983c2efc66386a6cae28a2508f59926b70688936))
* **policy:** add security.txt ([b4b2387](https://github.com/userstyles-world/userstyles.world/commit/b4b238714bbb7c1c656e5cf9367e5a53e51d0c7b))
* **robots:** disallow ModLog and API routes ([15053ff](https://github.com/userstyles-world/userstyles.world/commit/15053ff9ff786b2345a360fd4667ba4a80909afb))
* **SEO:** correct canonical for public pages ([4729b35](https://github.com/userstyles-world/userstyles.world/commit/4729b35309de8da521ae5ce07e45855d9ba17817))
* **static:** add our public GPG key ([4d76e2c](https://github.com/userstyles-world/userstyles.world/commit/4d76e2cf1950691f2d0a4de2dd890ac844d85a64))
* **stats:** add redundant checks for style IDs ([30d79cf](https://github.com/userstyles-world/userstyles.world/commit/30d79cfa284d735030f792c8b9eacbf45537f8f8))
* **templates:** add sort parameter to pagination ([1bd76b9](https://github.com/userstyles-world/userstyles.world/commit/1bd76b97bf838a5fb337077e7498b0153dbba4d2))
* **templates:** make pagination more extendable ([dc8193e](https://github.com/userstyles-world/userstyles.world/commit/dc8193ee5989f8a2e1b5152f4e063836d0a0f72d))
* **tools:** add 'fonts' command to run tool ([#76](https://github.com/userstyles-world/userstyles.world/issues/76)) ([56c6dfb](https://github.com/userstyles-world/userstyles.world/commit/56c6dfb78421adc4afd528f6bdbd4865ca400b5d))
* **tools:** download latest release of fonts ([2b2bfb5](https://github.com/userstyles-world/userstyles.world/commit/2b2bfb5d94a5b8fbc412cf2a7b1f04319b66405c))
* **ts:** check if style is installed ([2d3b0e0](https://github.com/userstyles-world/userstyles.world/commit/2d3b0e0281508e44b0d7ca74dbd42e7c08ee6d6c))
* **ui:** add disabled icon state to pagination ([3fcce76](https://github.com/userstyles-world/userstyles.world/commit/3fcce76f9c5356d737749189be9e093a7bc51716))
* **ui:** add links to our chat platforms ([8c71410](https://github.com/userstyles-world/userstyles.world/commit/8c714100fdb799f31faa20fae67edff34a3fe869))
* **ui:** add new info and styles for new footer ([83b2de4](https://github.com/userstyles-world/userstyles.world/commit/83b2de4abfeb6b7bae3fca2d5ab21e7559d3221e))
* **ui:** improve usability for pagination ([877cc91](https://github.com/userstyles-world/userstyles.world/commit/877cc91f634918d7187945a42dd369327f5bb8f3))
* **ui:** rewrite and expand info in page footer ([e99c3f1](https://github.com/userstyles-world/userstyles.world/commit/e99c3f13a8cb180e27bab331bb2b9a2f5077f9ea))
* update site's description ([90e8a08](https://github.com/userstyles-world/userstyles.world/commit/90e8a082a905319f0b67359efa387f9ad98b6a9f))
* **users:** allow 3–32 characters-long usernames ([f0a6d19](https://github.com/userstyles-world/userstyles.world/commit/f0a6d19104e0d21110b0db36a66f5d1b3fcb3a47))


### Bug Fixes

* **accessibility:** improve accessibility ([95180bc](https://github.com/userstyles-world/userstyles.world/commit/95180bcac21795a57c4f3938007d1995ea723e7c))
* **api:** replace wrong logging function ([9b5d629](https://github.com/userstyles-world/userstyles.world/commit/9b5d62944d7019a3d6d55c73f1f7b73c67de263b))
* **config:** update Matrix invite link ([dc15bca](https://github.com/userstyles-world/userstyles.world/commit/dc15bcad2ff899871505b2b5953f2a73c567d981))
* **CSP:** add exception for special case ([f713aad](https://github.com/userstyles-world/userstyles.world/commit/f713aadb96a71b680a9b12d9cd5f37a6959424de))
* **CSP:** fix typo ([f8ca9ef](https://github.com/userstyles-world/userstyles.world/commit/f8ca9ef4bba96ffdbe2017f392583b91f9322d4e))
* **css:** better alignment for flexbox link icons ([665689b](https://github.com/userstyles-world/userstyles.world/commit/665689b757540c71c66c39455f2a70a444189597))
* **css:** improve button cursor and transitions ([62840de](https://github.com/userstyles-world/userstyles.world/commit/62840de8214142bba1eb03ff08b89b83f7760bb3))
* **css:** improve styles for `select` element hack ([1323abd](https://github.com/userstyles-world/userstyles.world/commit/1323abd88a3129e3364c6eece92e3298d1322e5a))
* **css:** improve visuals for new invalid forms ([cd0c651](https://github.com/userstyles-world/userstyles.world/commit/cd0c6511435479f256bebb5280e23d1cbd95744f))
* **css:** remove extra white space from checkboxes ([72a3f7b](https://github.com/userstyles-world/userstyles.world/commit/72a3f7b46183b69b9c9e681693411bbc0f236165))
* **css:** remove invert filter from GitHub icon ([0f12e47](https://github.com/userstyles-world/userstyles.world/commit/0f12e47fff8ddad2f61c11697f2e32d4ca4fa440))
* **css:** resolve bugs on new style via OAuth page ([2e2b181](https://github.com/userstyles-world/userstyles.world/commit/2e2b181da9148d47df206c5bf4dc64eb3638e24b))
* **dashboard:** add parameter links in templates ([2e77bdd](https://github.com/userstyles-world/userstyles.world/commit/2e77bddac6a707270c761ce146029f6df1b815c8))
* **dashboard:** improve visuals for history charts ([5661557](https://github.com/userstyles-world/userstyles.world/commit/56615573a8dda3305b0fdc39703197ec89eef9be))
* **dashboard:** sort slices with proper data ([7e82111](https://github.com/userstyles-world/userstyles.world/commit/7e821114027cbdc2a4b60a3eaf133182f40037bd))
* **db:** add back sorting styles by views/installs ([79ee371](https://github.com/userstyles-world/userstyles.world/commit/79ee37131378fb5b1d19139b94e05d4359af3c1e))
* **db:** add default clause for USo-format index ([8fcb9a1](https://github.com/userstyles-world/userstyles.world/commit/8fcb9a1f664700202e6a0542cd4798e7d6560c78))
* **db:** don't take into account for deleted styles ([a44fa96](https://github.com/userstyles-world/userstyles.world/commit/a44fa968fef3c7129554e46aa976804aa6b90068))
* **docs:** correct parameters for security.md ([27fa97c](https://github.com/userstyles-world/userstyles.world/commit/27fa97c192b63ee0a7e1304da550441446cf9559))
* **dropdown:** match Gitea's sort menu ([#74](https://github.com/userstyles-world/userstyles.world/issues/74)) ([cf43988](https://github.com/userstyles-world/userstyles.world/commit/cf43988e675f8873c060f6289beebdab619229eb))
* **explore:** correct ordering of styles ([f6ac390](https://github.com/userstyles-world/userstyles.world/commit/f6ac39027f1a46b3d7118205213ddb4f4b817e9d))
* **handlers:** correct CSP values ([a856544](https://github.com/userstyles-world/userstyles.world/commit/a85654437a1cd31bdcf5792d2ff6bbc5e153a7c7))
* **html:** add boxes around forms on account page ([9d7bb79](https://github.com/userstyles-world/userstyles.world/commit/9d7bb79755e33ee36bd13761538d25fba86a466f))
* **html:** fix typo on review page ([d1ab4b9](https://github.com/userstyles-world/userstyles.world/commit/d1ab4b9efe1065e0996988bd8a8eebcc8603574d))
* **html:** improve text info on Review page ([bab897a](https://github.com/userstyles-world/userstyles.world/commit/bab897ab51eeee4f6355d63a3f8959f920275b76))
* **html:** remove last new line on description ([334817f](https://github.com/userstyles-world/userstyles.world/commit/334817f2341f0ceebf992a010a69aaeb0fa65053))
* **html:** resolve backdrop images being hidden ([fafd942](https://github.com/userstyles-world/userstyles.world/commit/fafd9420ac1ca6f927280400c223fd89be58a350))
* **oauth:** ensure correct redirect page after login ([e0514d9](https://github.com/userstyles-world/userstyles.world/commit/e0514d93409efdb83e1959eed5aa232ebf6e7a5a))
* **robots:** correct CSP value for robots.txt ([ffcfb0b](https://github.com/userstyles-world/userstyles.world/commit/ffcfb0b01da9724645d2a2ca3daffea5bbe0f5f2))
* **search:** no exception for notExist ([860f8fd](https://github.com/userstyles-world/userstyles.world/commit/860f8fd30220defb90f93491e50cf6e1d21fb55c))
* **stats:** try until all style IDs are returned ([0e6db8c](https://github.com/userstyles-world/userstyles.world/commit/0e6db8cfc534fba1612b7c9295b3db0315e037ab))
* **style:** correct canonical ([5f4738f](https://github.com/userstyles-world/userstyles.world/commit/5f4738f2379fadf70bcbe5b4ff6c0f001d29be01))
* **ts:** correct copy mechanism for span element ([0ca4543](https://github.com/userstyles-world/userstyles.world/commit/0ca4543e58f74c4e85cb2d2ef36bd57593cb2cd7))
* **ts:** correct function ordering ([5842c7e](https://github.com/userstyles-world/userstyles.world/commit/5842c7e27a18563ed283018a3be98e3e854824be))
* **ts:** defer script execution ([668b4d7](https://github.com/userstyles-world/userstyles.world/commit/668b4d7edb8feeb3811990a4763cb87cbb2def49))
* **ts:** initialize log for watch mode ([f7bace6](https://github.com/userstyles-world/userstyles.world/commit/f7bace6f034216046d82ea6601f98c3aa7123499))
* **ts:** remove redirect when on correct page ([1540f39](https://github.com/userstyles-world/userstyles.world/commit/1540f39f42e76b8b4a587b79ad58e540e3916448))

## [1.5.0](https://github.com/userstyles-world/userstyles.world/compare/v1.4.0...v1.5.0) (2021-07-18)


### Features

* **client:** improve UI for pagination component ([d48d2b1](https://github.com/userstyles-world/userstyles.world/commit/d48d2b18fda2298187e821f6ff4c2bfac34d0339))
* **cron:** speed up updating mirrored styles ([d94f27a](https://github.com/userstyles-world/userstyles.world/commit/d94f27ac74c7e29a99a293bdcc187efea2320d78))
* **css:** add highlight to modlog entry ([047a9ee](https://github.com/userstyles-world/userstyles.world/commit/047a9eea84cd65b237c58b3804a24f704e7de255))
* **dashboard:** add total History graph ([14906d2](https://github.com/userstyles-world/userstyles.world/commit/14906d208c19aa00e2a60652963605693f7cca88))
* **dashboard:** add User History charts ([08ac417](https://github.com/userstyles-world/userstyles.world/commit/08ac4175ddcc438e5eb4d235a4f4a98793fc0d71))
* **dashboard:** summarize new daily users/styles ([b4223a7](https://github.com/userstyles-world/userstyles.world/commit/b4223a717c95e88f3b3efdd3e6426dce940f485d))
* **email:** change email when password is resetted ([fd7eea6](https://github.com/userstyles-world/userstyles.world/commit/fd7eea66719dbb4f0f0a0b8131c5348eb3db2207))
* **html:** add anchors to Review entries ([03f531b](https://github.com/userstyles-world/userstyles.world/commit/03f531b0a3ce71b794559b28b45d63a0f9f5924a))
* **html:** improve the look for Review entries ([67c0820](https://github.com/userstyles-world/userstyles.world/commit/67c0820c13905c6235dd1dbf27fae9021758a7b5))
* **html:** show current release tag in the footer ([2b04944](https://github.com/userstyles-world/userstyles.world/commit/2b049443624b59827903aefdeeff32cc42c9496c))
* **mock db:** add style scope for default oauth app ([0ac8ab8](https://github.com/userstyles-world/userstyles.world/commit/0ac8ab84f5779fc65022e70a05a056ce17248b94))
* **models:** add notifications for style reviews ([b943c95](https://github.com/userstyles-world/userstyles.world/commit/b943c95d2ff8493cd88ab8d3b427443e1dd0f284))
* **modlog:** add don't hide behavior ([5d9b56b](https://github.com/userstyles-world/userstyles.world/commit/5d9b56b5c3689c991fc68d18bd96a21a2a1dbfa2))
* **modlog:** add modlog entry for removed styles ([74fa19f](https://github.com/userstyles-world/userstyles.world/commit/74fa19f6823cf6de8a6be7e93eae4c56891d938c))
* **modlog:** add option to change censor behavior ([deb9db1](https://github.com/userstyles-world/userstyles.world/commit/deb9db10c945cfc19d00143df31de5402ca5c3a2))
* **modlog:** allow entries to inside spoilers ([18ab383](https://github.com/userstyles-world/userstyles.world/commit/18ab3831bf7f41962cca74bf53d490edccfe2282))
* **modlog:** send email about mod actions to designated person ([7efcece](https://github.com/userstyles-world/userstyles.world/commit/7efcece6f34804daee82030d11053db6fa07030f))
* **oauthP:** add metadata validation ([7c6d2c5](https://github.com/userstyles-world/userstyles.world/commit/7c6d2c519b4021c409721d909741c0e4bac0c334))
* **reviews:** add core logic for Style reviews ([6b2d97f](https://github.com/userstyles-world/userstyles.world/commit/6b2d97f46fd07344be00d872356f9f897d816ef7))
* **reviews:** add logic to render Review form ([d95ce85](https://github.com/userstyles-world/userstyles.world/commit/d95ce855916a6f66dbba666ba3d9a304675714cb))
* **reviews:** add spam prevention for reviews ([1d4f3d1](https://github.com/userstyles-world/userstyles.world/commit/1d4f3d1b4b73bf14fefb14f59964993c6d5617ea))
* **reviews:** connect database logic to front-end ([15bde83](https://github.com/userstyles-world/userstyles.world/commit/15bde83f818b565ba1f5fcd7eeb03523b6b532d6))
* **reviews:** prevent too long or tall comments ([458791b](https://github.com/userstyles-world/userstyles.world/commit/458791bca6d02edaaf8dae94c0d68d771292b245))
* **server:** increase limiter's max connections ([827f2a7](https://github.com/userstyles-world/userstyles.world/commit/827f2a795a703a79c3294252fd52d777ffa5ea27))
* **styles:** 150k -> 100k chars for large styles ([acaf361](https://github.com/userstyles-world/userstyles.world/commit/acaf361a15ed6e3a31473ddad5076976c980cbd3))
* **styles:** add promotions to notifications ([c5477ad](https://github.com/userstyles-world/userstyles.world/commit/c5477ad9b885affbddbdaf0d7b85519de3c6b48e))
* **styles:** add total History graph in details ([92b0503](https://github.com/userstyles-world/userstyles.world/commit/92b050319899422c125ee4c2dac4ece1aa44f2e2))
* **styles:** implement ban functionality ([12fac37](https://github.com/userstyles-world/userstyles.world/commit/12fac374ca63397de69a3350481f75c0216f9a0c))
* **styles:** implement pagination on Explore page ([17f78a9](https://github.com/userstyles-world/userstyles.world/commit/17f78a90438498cf3d3ff15859fae2d9dbacbe3b))
* **styles:** send email on style promotion ([8a790a5](https://github.com/userstyles-world/userstyles.world/commit/8a790a5d2978d67187575d5bdb4ec6ca0f10ac00))
* **templates:** add global `config` function ([3544371](https://github.com/userstyles-world/userstyles.world/commit/35443713fbe4e6cbacf52b58aaac022e5b9c1f31))


### Bug Fixes

* **alloc:** ensure it's non-zero initialized ([b4e6e1b](https://github.com/userstyles-world/userstyles.world/commit/b4e6e1b5874a748180030821c6c70cc867fddcb7))
* **archive:** return nil instead of empty style ([2eb7196](https://github.com/userstyles-world/userstyles.world/commit/2eb7196e2bf396351c61441543f57db8e7ceccb0))
* **charts:** avoid rendering empty History charts ([775c7cb](https://github.com/userstyles-world/userstyles.world/commit/775c7cbb42c895b1fbbb653c7dd4c26fecd2068f))
* **charts:** remove total updates from the charts ([f3ba4fb](https://github.com/userstyles-world/userstyles.world/commit/f3ba4fbb36eecb9bc9fb42332488fbe250d7e779))
* **core:** account for remainder in pagination ([4f5f02a](https://github.com/userstyles-world/userstyles.world/commit/4f5f02a0479be0d15799e73a1a138393b0e0e3d8))
* **core:** Correct incorrect page numbers ([d834b82](https://github.com/userstyles-world/userstyles.world/commit/d834b82563b421cf7e861d62b304e23875a29623))
* **crypto:** don't use non-zero initalizion ([17068bb](https://github.com/userstyles-world/userstyles.world/commit/17068bb4ca41da9cc4c392b4ee321ef974ddbdb2))
* **css:** add [@media](https://github.com/media) ([0b57dcc](https://github.com/userstyles-world/userstyles.world/commit/0b57dcc87f42e9ab085f4fc010584bc326a87679))
* **css:** resolve some issues with buttons/forms ([231fb3a](https://github.com/userstyles-world/userstyles.world/commit/231fb3a407c8e1365a17d6f9204421952ba3c9bc))
* **css:** switch back to flexbox in Style previews ([968605c](https://github.com/userstyles-world/userstyles.world/commit/968605cff8d1a1886a7af22dd17d97ed8d33ff7b)), closes [#67](https://github.com/userstyles-world/userstyles.world/issues/67)
* **dashboard:** improve item order in User History ([d2815c4](https://github.com/userstyles-world/userstyles.world/commit/d2815c43d3b94f29e9bd7c525ff7d958b131021c))
* **dashboard:** make User History chart readable ([7c8c804](https://github.com/userstyles-world/userstyles.world/commit/7c8c804a40c1c484a4cc1418544c6a247f8468de))
* **email:** globalify email struct ([d03ada8](https://github.com/userstyles-world/userstyles.world/commit/d03ada8c0237aaacb781b43b14983ff98dcf2cc0))
* **errors:** don't expose literal errors ([899f88e](https://github.com/userstyles-world/userstyles.world/commit/899f88ef222cc88f9884fba35af81cc689e7dfb7))
* **gocritic:** fix false positive ([ba8f23c](https://github.com/userstyles-world/userstyles.world/commit/ba8f23ca99210e8dc0f630b294985a785396d760))
* **imports:** use import from defined name ([9c39b79](https://github.com/userstyles-world/userstyles.world/commit/9c39b7998883039ad853088ad4253517ba334e94))
* **log:** use correct log functions ([9feaea3](https://github.com/userstyles-world/userstyles.world/commit/9feaea3df1c10645a1d9166428c3d8935944ec8d))
* **mail:** Fix mail with only 1 part. ([7629d73](https://github.com/userstyles-world/userstyles.world/commit/7629d734b19679d50193473befddfaa3da4802bb))
* **oauth_login:** correct codebergStr ([2b54c05](https://github.com/userstyles-world/userstyles.world/commit/2b54c05cf73bc3cfecfcdce4a5e69b8e1990ad2f))
* **oauth:** broken code ([4d7e52b](https://github.com/userstyles-world/userstyles.world/commit/4d7e52befb19903f83e35850bea95d96fb70c6ed))
* **printf:** use correct type for string ([674018d](https://github.com/userstyles-world/userstyles.world/commit/674018d00dd37b39e1ecab91036db0b6f29307ec))
* **reviews:** add correct data for review entries ([75dba69](https://github.com/userstyles-world/userstyles.world/commit/75dba6916f8b3529f9dece458a13c1afc008755c))
* **reviews:** check if rating is out of range ([5b8fafe](https://github.com/userstyles-world/userstyles.world/commit/5b8fafe3225a3c52285497f8494035b842b3b9b3))
* **reviews:** resolve empty/multi-line comments ([60713ec](https://github.com/userstyles-world/userstyles.world/commit/60713ec3c06d8283858e78c60f47b94f6e67b1e7))
* **reviews:** use correct time function ([fd91a39](https://github.com/userstyles-world/userstyles.world/commit/fd91a39c9dcfc01a91175177e75c635d84e6e7a6))
* **search:** sort by default on relevance ([6f860fa](https://github.com/userstyles-world/userstyles.world/commit/6f860fa3fad60f00734650a8d2d55f6b605f7300))
* **styles:** hide review link to logged-out users ([d2082a3](https://github.com/userstyles-world/userstyles.world/commit/d2082a32dc12dafda61922b87326ad70b6961b94))
* **styles:** resolve a bug with preview image URLs ([e18255a](https://github.com/userstyles-world/userstyles.world/commit/e18255a35efbdc41d660e8c1ac29c8a054d705e9)), closes [#67](https://github.com/userstyles-world/userstyles.world/issues/67)
* **styles:** update to new USo-archive links ([1427178](https://github.com/userstyles-world/userstyles.world/commit/14271785aefa77a456a14bc7eb011050962b40c1))
* **todo:** remove non-possible todo ([9561f2c](https://github.com/userstyles-world/userstyles.world/commit/9561f2c121782e683a795e5f9322e0f6f490d835))
* **typecheck:** don't force types without checks ([f39bebb](https://github.com/userstyles-world/userstyles.world/commit/f39bebb9ad55a33116856294de5dbfc8ae295b6e))

## [1.4.0](https://github.com/userstyles-world/userstyles.world/compare/v1.3.0...v1.4.0) (2021-07-12)


### Features

* **core:** add more sorting options for styles ([d14f171](https://github.com/userstyles-world/userstyles.world/commit/d14f171abd7e40da4aa61c83eb661b1baeddb1db))
* **crypto:** add nonce scrambling ([972e88a](https://github.com/userstyles-world/userstyles.world/commit/972e88afd4b3c8db76919748d29fc55a0c9fa275))
* **crypto:** add nonce scrambling into crypto functions ([63e2b29](https://github.com/userstyles-world/userstyles.world/commit/63e2b298ea4c1562e5b10bc874a38d29119fb98f))
* **css:** improve styles for charts ([bc9820d](https://github.com/userstyles-world/userstyles.world/commit/bc9820d183084d89380190aa37ea96744cea7029))
* **css:** improve styles for content in docs ([5c2145f](https://github.com/userstyles-world/userstyles.world/commit/5c2145f8ea016fe30acdb9fa1823391c82d298d5))
* **css:** tweak and improve styles for graphs ([8084d3a](https://github.com/userstyles-world/userstyles.world/commit/8084d3a187569e75388daf8eed762ee8b398bf78))
* **dashboard:** add functionality to ban users ([8a400b7](https://github.com/userstyles-world/userstyles.world/commit/8a400b70ffc499adb7fdc6baf40d0b7614367155))
* **dashboard:** allow moderators to ban users ([5f8bad4](https://github.com/userstyles-world/userstyles.world/commit/5f8bad44d93fc61cbca6f39e4d59443014276169))
* **dashboard:** show History chart for all styles ([8388dfa](https://github.com/userstyles-world/userstyles.world/commit/8388dfaa762cf16520a214e7bb2ec6295bbd255d))
* **dashboard:** show user's email field to admins ([dd58541](https://github.com/userstyles-world/userstyles.world/commit/dd58541eddc1c075d5da2edb9c0dda01925a0cc3))
* **db:** improve error checking for table actions ([7e61f23](https://github.com/userstyles-world/userstyles.world/commit/7e61f23bb4706e8c3da4a9614f030a05d6c9a646))
* **db:** improve migrate/drop/seed functionality ([7c054b8](https://github.com/userstyles-world/userstyles.world/commit/7c054b8731bc1f549fe13bbb56d672783c3adb0d))
* **docs:** add more content to FAQ page ([8f29262](https://github.com/userstyles-world/userstyles.world/commit/8f2926202fd48633771aad54847b67fba9b4cc2b))
* **docs:** expand our internal documentation ([458efc6](https://github.com/userstyles-world/userstyles.world/commit/458efc6851398aab4f8d20b64bf213820b6c0e7b))
* **html:** add 'default' option to sort menus ([4559fe3](https://github.com/userstyles-world/userstyles.world/commit/4559fe3419c8f67aae7cff30cf42b307d0cfcfa1))
* **html:** add icons for Style action buttons ([9fb1f62](https://github.com/userstyles-world/userstyles.world/commit/9fb1f62c0889c4e5a74ce7891908f882d8e13671))
* **html:** add icons to homepage/save buttons ([adfb0cf](https://github.com/userstyles-world/userstyles.world/commit/adfb0cfe25442a36f76112152e8828a90eca906c))
* **html:** add icons to search/sort buttons ([fa7d515](https://github.com/userstyles-world/userstyles.world/commit/fa7d51568b4ea52805a18c1008dba36215f43721))
* **html:** add legends to History chart ([fca0ff8](https://github.com/userstyles-world/userstyles.world/commit/fca0ff8fb6dc9195b139b0cb93d950c6a3e3cb6c))
* **html:** append git commit hash to JS file ([920ad65](https://github.com/userstyles-world/userstyles.world/commit/920ad659d230769d6aa7d8880b23c582d61184ef))
* **html:** implement meta tags for User profiles ([f1e7070](https://github.com/userstyles-world/userstyles.world/commit/f1e70704645232bf9766e5bba2bc2b79d56e224a))
* **html:** improve meta tags for Style pages ([5e45f38](https://github.com/userstyles-world/userstyles.world/commit/5e45f38cb27199544060e6f7c0779b14d4f66ee8))
* **html:** show counters on dashboard page ([89cfa53](https://github.com/userstyles-world/userstyles.world/commit/89cfa530cafaea12cefa323b2d8b96752e17f8c7))
* **images:** bump image caching to a month ([eab65f9](https://github.com/userstyles-world/userstyles.world/commit/eab65f965ec796d71d2eeae4a9d9b6a655fc49cf))
* **images:** increase quality for screenshots ([55b7182](https://github.com/userstyles-world/userstyles.world/commit/55b71822b8908f6f03efd74c7e53cdaddd75dbcd))
* **login:** bump remember me to a month ([8457887](https://github.com/userstyles-world/userstyles.world/commit/8457887e8b781610673929509cf31b0fe0c99303))
* **modlog:** add reason & log when banned user ([ea0c575](https://github.com/userstyles-world/userstyles.world/commit/ea0c57516f8701b51c6b4baaf7732bf8f64cc20e))
* **modlog:** create modlog prototype ([664f3eb](https://github.com/userstyles-world/userstyles.world/commit/664f3eb9e0ad637a008973397310a111bebdc2dc))
* **styles:** add data for combined statistics ([68bda02](https://github.com/userstyles-world/userstyles.world/commit/68bda021aba7a6da074945fed45f6cbd898fea2b))
* **styles:** add logic for visualizing statistics ([8be681d](https://github.com/userstyles-world/userstyles.world/commit/8be681db6f9f052ad5722e077ec8f6b4f70992e3))
* **styles:** get history data for Style stats ([fd7b739](https://github.com/userstyles-world/userstyles.world/commit/fd7b739a9368c4ec11f51c57ad9420ebf537f516))
* **styles:** show userstyle stats in History area ([5e82d06](https://github.com/userstyles-world/userstyles.world/commit/5e82d06e4c00eed3f8901b8fda870e21a16dd71f))


### Bug Fixes

* **api:** add missing logic for JPEG screenshots ([e5bd33b](https://github.com/userstyles-world/userstyles.world/commit/e5bd33bca656aab5af629a0acc474c7dcbc57469))
* **auth:** resolve cookie issues with Vim Vixen ([b37d8c1](https://github.com/userstyles-world/userstyles.world/commit/b37d8c1a1229e66c99fde585ee9cbdeadd359fe3))
* **core:** trim spaces in reason field for Mod Log ([2ba5457](https://github.com/userstyles-world/userstyles.world/commit/2ba5457fdf99437bc69e2193f62e461b6b9680d0))
* **crypto:** don't panic on incorrect input ([6e665fe](https://github.com/userstyles-world/userstyles.world/commit/6e665fe0ea342ec790bb55ab76301c7a380462b4))
* **css:** hide horizontal overflow on body element ([fb45db9](https://github.com/userstyles-world/userstyles.world/commit/fb45db907c4590d579b72b2927343229c924400b))
* **css:** hide overflow in style cards ([0bc94b5](https://github.com/userstyles-world/userstyles.world/commit/0bc94b5efa0a79c697d268a4fb2796d47970bd9e))
* **css:** improve styles for textareas and select ([65cd64e](https://github.com/userstyles-world/userstyles.world/commit/65cd64e052af4bdd6aeca7bd561cc84e9b8bc3e6))
* **css:** improve styles on Style's link page ([f56268a](https://github.com/userstyles-world/userstyles.world/commit/f56268a40205d28ab3c41703368ca5ab05254826))
* **dashboard:** prevent being able to ban yourself ([bacff96](https://github.com/userstyles-world/userstyles.world/commit/bacff9646f464a12039a821f3ff61a81386766f8))
* **db:** add created date to search cards query ([5bd8ff6](https://github.com/userstyles-world/userstyles.world/commit/5bd8ff69cae805b17d6ac31fa6f5f3a436e18452))
* **err:** add missed letter ([24b7611](https://github.com/userstyles-world/userstyles.world/commit/24b7611254db71c8b2b621afbdb341971885c454))
* **html:** add missing History chart on dashboard ([925687d](https://github.com/userstyles-world/userstyles.world/commit/925687dbcf6cb6cbef4a14097c622de46c5be63d))
* **html:** resolve a bug with alerts ([61ac607](https://github.com/userstyles-world/userstyles.world/commit/61ac6075ac73d6e668b8700f344c6b2bf0b8f528))
* **images:** remove avif ([9e13023](https://github.com/userstyles-world/userstyles.world/commit/9e13023b9855b6c8a8557ad53c681b60da345d6a))
* **json:** fix json tags ([d23c405](https://github.com/userstyles-world/userstyles.world/commit/d23c40587ec5befb09fcc11b8f6ab64cf4d0463e))
* **modlog:** add Modlog to nav menu ([0df8fa6](https://github.com/userstyles-world/userstyles.world/commit/0df8fa65cb3a6f74044d5ebcbcb8d991b50d98e7))
* **modlog:** add User and Title to render ([f55fc2b](https://github.com/userstyles-world/userstyles.world/commit/f55fc2b91ff5769a1dcc95a9c373733058e127be))
* revert this line ([609729d](https://github.com/userstyles-world/userstyles.world/commit/609729d12c9199fe5337575ea59d696b25b1a2e2))
* **styles:** remove History if data doesn't exist ([f63ed61](https://github.com/userstyles-world/userstyles.world/commit/f63ed610ae76481250bd92a0bc3acda8c06ad126))
* **styles:** resolve an issue with SEO logic ([d6cea95](https://github.com/userstyles-world/userstyles.world/commit/d6cea9547beff25944cc7654defa425d19e78c7d)), closes [#62](https://github.com/userstyles-world/userstyles.world/issues/62)
* **ts:** auto-fill proper metadata for homepage ([afc459b](https://github.com/userstyles-world/userstyles.world/commit/afc459b8752b6d8c1ed714afcd0f354fd39cbcf3))
* **ts:** prevent prototype pollution ([a086a49](https://github.com/userstyles-world/userstyles.world/commit/a086a49b36ffa63f096d6598cfb1aab3425d40b4))

## [1.3.0](https://github.com/userstyles-world/userstyles.world/compare/v1.2.1...v1.3.0) (2021-06-19)


### Features

* **dashboard:** add WIP moderation tools ([c2495ce](https://github.com/userstyles-world/userstyles.world/commit/c2495ce21e799366e320b7e26e1ef509d1dd759d))
* **html:** enable more Markdown extensions ([fee7ced](https://github.com/userstyles-world/userstyles.world/commit/fee7ceda263b98303aace9881a0bd53429862ae1))
* **oauth_login:** retrieve email from authorized user ([6f026ba](https://github.com/userstyles-world/userstyles.world/commit/6f026ba15ab99328a254454d1304cd9512f5526b))
* **oauthP:** create new Style based off styleInfo ([31813da](https://github.com/userstyles-world/userstyles.world/commit/31813da30088eb9438299808ac3232f1e2c9fd5d))
* **oauthP:** enable option to pre-fill information ([461ddb0](https://github.com/userstyles-world/userstyles.world/commit/461ddb03c703867d356235c90659a7f89f47df05))
* **search:** sync with database ([e736e84](https://github.com/userstyles-world/userstyles.world/commit/e736e844e6447822940aefe2d72ed296771b0a7d))
* **ts:** add more fields to auto-fill ([7d8c4c1](https://github.com/userstyles-world/userstyles.world/commit/7d8c4c124891f929d4fc191614a2f16a2ab1fa3b))
* **ts:** allow description to be set ([5a5476e](https://github.com/userstyles-world/userstyles.world/commit/5a5476e0a5c5846237bc47cb020ec95ca8cda6d4))
* **ts:** set default meta ([8e1405e](https://github.com/userstyles-world/userstyles.world/commit/8e1405ea4007c658155d3cdd51845ce5096ca986))


### Bug Fixes

* **css:** align thumbnails on the left side ([cef75dc](https://github.com/userstyles-world/userstyles.world/commit/cef75dc2a56ad4b28dc339096ace2f3ccd284fd0))
* **css:** avoid resizing userstyle screenshots ([d9adc11](https://github.com/userstyles-world/userstyles.world/commit/d9adc1154840b689b109ac251d30b1acc70edbb8))
* **db:** default to null instead of time.Time's default ([846fbec](https://github.com/userstyles-world/userstyles.world/commit/846fbec2c7390308c23a010587e403209ce1faad))
* **history:** update queries for stat history ([2bb1110](https://github.com/userstyles-world/userstyles.world/commit/2bb1110eb261980bc711b27a5cdb77b14f705f79))
* **js:** make sure data is an object ([71303e7](https://github.com/userstyles-world/userstyles.world/commit/71303e77b90ca71c081b30402eedfc003dc077eb))
* **oauth_login:** set default role to regular ([2717bf6](https://github.com/userstyles-world/userstyles.world/commit/2717bf6f06e02518a66445d5768ffccbb55d9389))
* **oauthP:** listen to POST request instead of GET ([869b69d](https://github.com/userstyles-world/userstyles.world/commit/869b69d615049c823e74ddc63657ee15c6e7a5c8))
* **stats:** improve accuracy of style statistics ([d9910a1](https://github.com/userstyles-world/userstyles.world/commit/d9910a18c89d940443ce436f136103505247e12e))
* **stats:** include data from previous scheme ([2bf7ce4](https://github.com/userstyles-world/userstyles.world/commit/2bf7ce44ef28c7c9c3983b0b89db14bcb0ed73b2))
* **stats:** update queries on the home page ([17809dc](https://github.com/userstyles-world/userstyles.world/commit/17809dc2824715b06a7e15a8a035065430cd67a5))
* **styles:** update queries for style cards ([2885dde](https://github.com/userstyles-world/userstyles.world/commit/2885dde32887eb4cbb2f288f21c78af9094d7a0a))
* **ts:** compile production version ([c7f58e0](https://github.com/userstyles-world/userstyles.world/commit/c7f58e041f9080b4fb1a6eb159c68b4bbfccc8ca))
* **ts:** safe remove of element ([30c2a39](https://github.com/userstyles-world/userstyles.world/commit/30c2a39ab87f27a28b95ec99d2b65fed97160191))

## [1.2.1](https://github.com/userstyles-world/userstyles.world/compare/v1.2.0...v1.2.1) (2021-05-30)


### Features

* **oauthP:** show all oauth of user ([76b832b](https://github.com/userstyles-world/userstyles.world/commit/76b832b1379969f6dddeda9cd62c26a345dca674))
* **oauthP:** show clientID + clientSecret  ([5491705](https://github.com/userstyles-world/userstyles.world/commit/54917053a6bda3dd8e09d0b3a8ecca63b89b8cd9))


### Bug Fixes

* **core:** sort by created date by default ([9611b53](https://github.com/userstyles-world/userstyles.world/commit/9611b536bfe7cd817157d0468e25327962ab591e))
* **css:** add back hints to various links ([62949f4](https://github.com/userstyles-world/userstyles.world/commit/62949f48bbc5236ee3a584c2e2436288d21bd588))
* **html:** add custom checkbox to register page ([c2d39c3](https://github.com/userstyles-world/userstyles.world/commit/c2d39c336293a3f9f0064b029e9ab356360ed281))
* **oauth:** redirect to oauth after creation ([e080ea6](https://github.com/userstyles-world/userstyles.world/commit/e080ea6684d03b9abe30bdfa06e9f902bd1dd8de))
* **tools:** don't watch static folder ([df9c0b3](https://github.com/userstyles-world/userstyles.world/commit/df9c0b36900829354f55f909caef858e7a76136b))

## [1.2.0](https://github.com/userstyles-world/userstyles.world/compare/v1.1.0...v1.2.0) (2021-05-29)


### Features

* **api:** add /style endpoint ([b05a87d](https://github.com/userstyles-world/userstyles.world/commit/b05a87d29116e5689a427b48139838226746a444))
* **css:** improve reduced motion CSS ([0b4fb9f](https://github.com/userstyles-world/userstyles.world/commit/0b4fb9fe3a509fe6a6467a5884e07a2a735011f4))
* **css:** improve the look for Style share button ([7bc3750](https://github.com/userstyles-world/userstyles.world/commit/7bc3750e02d95d15c42603f63552de96f85b641d))
* **css:** limit resize to vertical on textarea ([5b7170b](https://github.com/userstyles-world/userstyles.world/commit/5b7170bd41b5b0a9bc22429e79be5ba0972d21e1))
* **css:** tweak white space for Markdown headings ([a35aea0](https://github.com/userstyles-world/userstyles.world/commit/a35aea017d25009726da4ca305bc584926f9acbf))
* **html:** add sorting on Search/Explore pages ([f2570ee](https://github.com/userstyles-world/userstyles.world/commit/f2570ee4e35dca50552d24a9aa7374789e4c0c28)), closes [#42](https://github.com/userstyles-world/userstyles.world/issues/42) [#46](https://github.com/userstyles-world/userstyles.world/issues/46)
* **html:** remember state for select elements ([5dd8170](https://github.com/userstyles-world/userstyles.world/commit/5dd8170ec56b207fbe0ec1e0b6a47b3c147f6520))
* **login:** add redirect_uri after login ([5f98efa](https://github.com/userstyles-world/userstyles.world/commit/5f98efa215c22b4ecc5c5b9d364a363cb59e3aad))
* **oauthP:** return Style ID on authorize_style ([2f724da](https://github.com/userstyles-world/userstyles.world/commit/2f724da093fbfd69f6ff389440b80a7b20de7884))
* actually glue it together ([8868ade](https://github.com/userstyles-world/userstyles.world/commit/8868adeb54d7c500d2fe0af2e63af624e93bb6cf))
* Add documentation ([4696462](https://github.com/userstyles-world/userstyles.world/commit/4696462e45b193d4d95f41fefea91eea7893494b))
* add first endpoints ([635ae9e](https://github.com/userstyles-world/userstyles.world/commit/635ae9e9a8392c7baffceb119706593a1c974a80))
* add scope information  ([2b96bbb](https://github.com/userstyles-world/userstyles.world/commit/2b96bbb70483f7b3adea80a50f10e6bcf1dd81df))
* add style ([9c95227](https://github.com/userstyles-world/userstyles.world/commit/9c952277c9166f1123bb9a7794d1ce1e73d64908))
* add/edit oauth ([294714f](https://github.com/userstyles-world/userstyles.world/commit/294714f09a69be62f53a89ce5df12e62a0fa2d36))
* allow access_token retrieval ([0ac7699](https://github.com/userstyles-world/userstyles.world/commit/0ac7699b835ac53c77f879cc14871723bcb1043f))
* delete style ([c3ff6d7](https://github.com/userstyles-world/userstyles.world/commit/c3ff6d71859b3f9df413aa75e58182db73d1f921))
* edit styles programmatically ([621836c](https://github.com/userstyles-world/userstyles.world/commit/621836c751c52a4dd7dcf681511afff22a560824))
* Faster json encoding ([e9ef25f](https://github.com/userstyles-world/userstyles.world/commit/e9ef25f975b3b48bd83d001311d1f0c146de05e2))
* list all styles of user ([38c6212](https://github.com/userstyles-world/userstyles.world/commit/38c621291710b487349ed7a069674a3a9ac91181))
* make authorize protected ([7104481](https://github.com/userstyles-world/userstyles.world/commit/710448160ae2979e888fe4da6a44f229a7592242))
* process authorization ([a324ecf](https://github.com/userstyles-world/userstyles.world/commit/a324ecf74113178f2fe8c66d6da9538ae581ff8c))
* try to glue this linking style ([ed0bcb7](https://github.com/userstyles-world/userstyles.world/commit/ed0bcb72880a32d2a4669600850da1d8757f0a6e))
* use callback helper ([d420e16](https://github.com/userstyles-world/userstyles.world/commit/d420e1695b888ca9afda9674578e397f792f1c51))
* use OJG for decoding ([47f3391](https://github.com/userstyles-world/userstyles.world/commit/47f3391d6054c679593704f9516e5d61b3dff71c))
* use userID ([af72bfd](https://github.com/userstyles-world/userstyles.world/commit/af72bfd30919296a621050cbced5a52c987948b2))
* validate OAuth input ([c104429](https://github.com/userstyles-world/userstyles.world/commit/c104429270d1960bd5bd9b2ae4e679b2ce7b99df))
* validate upon authorize ([7dfdd20](https://github.com/userstyles-world/userstyles.world/commit/7dfdd208bce7d107a09f641aae38373cf4205bd9))
* **html:** add tooltip to style card ([cbc1e87](https://github.com/userstyles-world/userstyles.world/commit/cbc1e87458d57e3bd9baf6a6f1cbf4dc9b07c5f7))


### Bug Fixes

* **conflicts:** patch up conflicts from rebase ([70a0535](https://github.com/userstyles-world/userstyles.world/commit/70a05353d1b3c65cb5f05b75e67ab99f10611788))
* **css:** update selector for card background ([d52ddd4](https://github.com/userstyles-world/userstyles.world/commit/d52ddd48e30fd560df356e0ba0273af3195917b9))
* **db:** resolve failing auto-migration for OAuth ([7ef03a5](https://github.com/userstyles-world/userstyles.world/commit/7ef03a546083bf368c5778ceae71b6fed7b16ee1))
* **html:** add custom checkboxes to Import page ([5ffc920](https://github.com/userstyles-world/userstyles.world/commit/5ffc9202ec55c67ee8ae86a8d53c0dc57f300de3))
* **html:** add custom checkboxes to OAuth page ([b39ae38](https://github.com/userstyles-world/userstyles.world/commit/b39ae3805af12c1b80f29d40b9124e3ed8dc23c5))
* **html:** add missing sign-in icon ([ea2ba87](https://github.com/userstyles-world/userstyles.world/commit/ea2ba876243a27cb7c555fb7b790bd96850dff97))
* **html:** prevent adding styles with empty names ([f7be151](https://github.com/userstyles-world/userstyles.world/commit/f7be151c1ec3965368a84986c4c5d5d447b86bc9))
* **html:** relax validation for search/style names ([b1eb15b](https://github.com/userstyles-world/userstyles.world/commit/b1eb15bdb2d8275b1d94651c82a5a07513fc5afc)), closes [#60](https://github.com/userstyles-world/userstyles.world/issues/60)
* **login:** just include path within after_login ([133de35](https://github.com/userstyles-world/userstyles.world/commit/133de35517753f340df8c9b46666b104fd0ef381))
* **models:** resolve a bug with UpdateOAuth query ([e024042](https://github.com/userstyles-world/userstyles.world/commit/e024042c00ca3b07c300ff94c88082d636a3f2ac))
* **oauth:** resolve a panic from bad validation ([b64229a](https://github.com/userstyles-world/userstyles.world/commit/b64229a6bd03bc9641ae651a0d7609ff27473c83))
* **oauthP:** correct link ([98a55dc](https://github.com/userstyles-world/userstyles.world/commit/98a55dcc838b80647342b3686b6d17d55efceaec))
* **security:** validate style's userID ([dd8d7b1](https://github.com/userstyles-world/userstyles.world/commit/dd8d7b1545aeb0332cc185a43df3f37bc3585894))
* **server:** remove unused imports ([aac127a](https://github.com/userstyles-world/userstyles.world/commit/aac127abee07d72cbc479ae54f9db14ab5b0defb))
* **tools:** don't watch node_modules ([87f719e](https://github.com/userstyles-world/userstyles.world/commit/87f719ed1d4c03a099631970dc39dd8bffa0e871))
* actually renders authorize page ([59b63b8](https://github.com/userstyles-world/userstyles.world/commit/59b63b8cd73e949a8ec11e0b14afa64c05b6eaaf))
* add explanation to weird filter func ([85143cd](https://github.com/userstyles-world/userstyles.world/commit/85143cdb73dc9440855a2fd36776d25e50e954b7))
* correct biography ([f412ef7](https://github.com/userstyles-world/userstyles.world/commit/f412ef7ec5ce7ae0ddc693a82ec4d8e30593c5e0))
* correct table formatting ([116a746](https://github.com/userstyles-world/userstyles.world/commit/116a7461a7b358f63c7ef2fe26c68723deffa261))
* correct token_type ([5115b57](https://github.com/userstyles-world/userstyles.world/commit/5115b5777dca894a4a5cd188c5793d9fef908ec1))
* correct url's ([6aa6393](https://github.com/userstyles-world/userstyles.world/commit/6aa63934ce4abb9447729ac8bb40e3045246ce34))
* handle multiple scopes ([7900436](https://github.com/userstyles-world/userstyles.world/commit/7900436fedf908e02ad1f52d03b1e07e445cded1))
* it's username ([2603e36](https://github.com/userstyles-world/userstyles.world/commit/2603e360d6a0ca631527fdb0b9e2a4b5b6fb9266))
* naming :D ([90acb75](https://github.com/userstyles-world/userstyles.world/commit/90acb750f14c0d9e8010896e945cec349b1a8418))
* pass options correctly (pointer reference) ([6750a9c](https://github.com/userstyles-world/userstyles.world/commit/6750a9cff04d4737a68f9c16a03a8ab695639296))
* use correct go.mod/sum ([75b1426](https://github.com/userstyles-world/userstyles.world/commit/75b14263e7b9a61ab46ce08c9168f2490def8830))
* use Post ([8efc8b4](https://github.com/userstyles-world/userstyles.world/commit/8efc8b4d28f00b53691b85fbfc650d627764e82d))
* **styles:** make url relative to canonical ([93eded2](https://github.com/userstyles-world/userstyles.world/commit/93eded2678fab8c0e4c1c5e0fe716bca3e52820f))

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

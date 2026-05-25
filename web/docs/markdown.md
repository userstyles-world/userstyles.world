---
Title: Markdown Reference
---

# Markdown Reference

You can use [Markdown](https://en.wikipedia.org/wiki/Markdown) syntax in such places like Style Notes and Profile Biography.


## Linking

### Hyperlinks

Wrap link text in `[`square brackets`]` and the link in `(`round brackets`)` to create [hyperlink](/):

`[Explore](https://userstyles.world/explore)`

### Embedded content

Add `!` to image's hyperlink to make it shown as image:

`![](https://raw.githubusercontent.com/openstyles/stylus/67b17335f9da323a359a74a90312df4ab76c425f/images/icon/128.png "Shades of cyan")`

The text in square brackets serves purpuse of alternative text if the image isn't loaded. You can keep it empty or add something, for example: `Stylus logo`.  
The text in quotation marks is for popup. Its optional, you can remove the quotation marks.

<!-- it won't work on dev due to CSP -->

![Screenshot](/preview/7/1.webp "What sup")


## Styling

You can apply different styles to your text.
You can combine different styles for the same text.

### Bold

Wrap text in `**`stars`**` or `_`underscores`_` to make it **bold**.

### Italic

Wrap text in `*`stars`*` or `__`underscores`__` to make it **italic**.

### Strikethrough

Wrap text in `~~`tildas`~~` it ~~crossed out~~.

### Highlight

You could wrap text in `==`tildas`==` to make it <mark>highlighted</mark>, but you have to wrap it in `<mark>`HTML tags`</mark>` instead.


## Code

### Inline

Wrap code in backticks to make `inline code block` out of it.

### Block

Wrap code in three backticks to create a code block:

```
it could
also be
multi line!
```

### Quote

Start line with `>` to make quote out of it:

> Here's quote.


## Header

You can have different sized headers. Start line with `#` of different amounts to have headers accordingly:

`#`, `##`, `###`, `####`, `#####`, `######`.


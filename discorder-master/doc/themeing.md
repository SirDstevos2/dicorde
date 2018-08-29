# Themeing

Themes are located in a themes folder in your discorder configuration dir ( `~/.config/discorder` on Unix, `%APPDATA%/discorder` on Windows).

To pick a theme you can use the theme selection menu or specify the theme in your `discorder.json` file manually. In your themes folder there is also the theme `default`, which is the default theme.

To force a theme:

`-t "path/to/theme.json"`: Forces use of the given theme

### Theme file

An example theme:

```json
{
    "name": "SUPER FUN THEME",
    "author": "jonas747",
    "comment": "SO MUCH FUN OMG",
    "color_mode": 1,
    "window_border": {
        "fg": {"color": "white","bold": true},
        "bg": {"color": "blue"}
    },
    "window_fill": {
        "fg": { "color": "default"},
        "bg": { "color": "black"}
    }
}
```

This works the same way as keybinds by overriding the default theme, so fields not present here will be set to defaults.

The color mode specified the color mode the theme utilize. Note that some terminals don't support certain color modes.


### Colors

#### Modes

0 - OutputNormal
1 - Output256
2 - Output216
3 - OutputGrayscale

In colormode normal you can specify the colors as strings if you want ("red" for 2 for example).

#### Colors (In colormode normal):

0 - Default
1 - Black
2 - Red
3 - Green
4 - Yellow
5 - Blue
6 - Magenta
7 - Cyan
8 - White

#### Attributes

 - bold
 - underline
 - reverse

Combining colors with bold will give a brighter color if you have enabled it in your terminal profile settings.

For colors in other modes see http://misc.flogisoft.com/bash/tip_colors_and_formatting.

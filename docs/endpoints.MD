# Endpoints

This documentation is focused on the behaviour of Userstyles world(USw)'s `/api` endpoints.
URL's are referred as ex. `/user` and thereby will mean the endpoint is located at `https://userstyles.world/api/user/`.

## Basic information

All requests are returned as `application/json` regardless of the `Accept` header.  
All responses are compressed by default, based on the `Accept-Encoding` header.  
All _correct_ GET-response are returned within a `data` property of the returned JSON,  
please note that the example response won't show this.


### Retrieve user's information

**Authorization is required**
```
GET /user
```
Get the user's information.

Example response
```JSON
{
    "username": "Gusted",
    "email": "notmyemail@example.com",
    "biography": "I'm _one_ of the UserStyles World Developers.\nI'm aware, **I'm gay** 🌈.",
    "role": "Moderator",
    "socials": {
        "github": "gusted",
        "gitlab": "gusted",
        "codeberg": "gusted"
    }
}
```

_**Note:** The email field can be empty when a user has registered through an OAuth provider._

### Retrieve a specific user's information
```
GET /user/<id> 
```
or
```
GET /user/<username>
```
Gets a specific user's information, specified with the user's ID/username. The username is capitalization-sensitive.

Example response of id=8/username="Gusted"
```json
{
    "username": "Gusted",
    "biography": "I'm _one_ of the UserStyles World Developers.\nI'm aware, **I'm gay** 🌈.",
    "role": "Moderator",
    "socials": {
        "github": "gusted",
        "gitlab": "gusted",
        "codeberg": "gusted"
    }
}
```

### List styles

**Authorization is required**
```
GET /styles
```

Gets the user's styles and list them with basic information.

Example response
```JSON
[
    {
        "ID": 1,
        "CreatedAt": "2021-04-01T23:44:27.067587104+02:00",
        "UpdatedAt": "2021-04-06T21:54:08.396097219Z",
        "Name": "UserStyles.world Tweaks",
        "Description": "Various tweaks for UserStyles.world that might end up being upstreamed.",
        "Notes": "After you install this userstyle, you can configure it from _configuration menu_ in Stylus' popout menu by clicking on the icon on your toolbar.",
        "Code": "Left out due to length in this example.",
        "License": "MIT",
        "Preview": "https://user-images.githubusercontent.com/18245694/111090605-9cc4ce00-8530-11eb-8716-2c1df40f32e1.png",
        "Homepage": "https://github.com/userstyles-world/tweaks",
        "Archived": false,
        "Featured": true,
        "Category": "userstyles.world",
        "Original": "https://raw.githubusercontent.com/userstyles-world/tweaks/main/tweaks.user.styl",
        "MirrorCode": true,
        "MirrorMeta": false,
        "UserID": 1,
        "Username": "vednoc"
    },
    {
        "ID": 2,
        "CreatedAt": "2021-04-02T17:20:59.830991125+02:00",
        "UpdatedAt": "2021-04-02T17:25:03.430384406+02:00",
        "Name": "Dark-GitHub",
        "Description": "Dark and light, very customizable theme for GitHub.",
        "Notes": "",
        "Code": "Left out due to length in this example.",
        "License": "MIT",
        "Preview": "https://user-images.githubusercontent.com/18245694/109351173-18f4bb80-7879-11eb-998f-3fc31d2abb55.png",
        "Homepage": "https://github.com/vednoc/dark-github",
        "Archived": false,
        "Featured": true,
        "Category": "github.com",
        "Original": "https://raw.githubusercontent.com/vednoc/dark-github/main/github.user.styl",
        "MirrorCode": true,
        "MirrorMeta": false,
        "UserID": 1,
        "Username": "vednoc"
    }
]
```

### Retrieve specific style's information
```
GET /style/<id>
```
Gets a specific style's information, specified with the style's ID.

Example response ID=29
```JSON
{
    "ID": 29,
    "CreatedAt": "2021-04-03T17:49:35.16120309Z",
    "UpdatedAt": "2021-04-18T00:53:52.35042565Z",
    "Name": "Userstyles.World Tweaks!",
    "Description": "Tweaks for Userstyles.World!",
    "Notes": "",
    "Code": "Left out due to length in this example.",
    "License": "",
    "Preview": "https://codeberg.org/Freeplay/UserStyles/raw/branch/main/images/userstyles.world_tweaks-1.png",
    "Homepage": "https://codeberg.org/Freeplay/UserStyles",
    "Archived": false,
    "Featured": true,
    "Category": "Userstyles.world",
    "Original": "https://codeberg.org/Freeplay/UserStyles/raw/branch/main/Userstyles.world_tweaks.user.css",
    "MirrorCode": true,
    "MirrorMeta": false,
    "UserID": 11,
    "Username": "Freeplay"
}
```

### Edit specific style

**Authorization is required**
```
POST /style/<id>
```

Edit a user's style specified with the style's ID.
The only changes that will be changing are the ones that are being sent,  
this means you can send partial information.

Example request
```JSON
{
    "name": "The changed name of the style",
    "description": "A wordy description of the edited style",
    "notes": "Hey, if you install this make sure to remember X",
    "code": "Some actually valid UserCSS code here",
    "preview": "https://i.ytimg.com/vi/-Cv68B-F5B0/maxresdefault.jpg",
    "homepage": "https://github.com/userstyles-world/userstyles.world/",
    "license": "WTFPL",
    "category": "youtube.com"
}
```

Example response
```JSON
Status: 200
{
    "info": "Succesful edited style!"
}
```
or
```JSON
Status: 406
{
    "err": "Homepage validation is wrong."
}
```

### Add new style

**Authorization is required**
```
POST /style/new
```

Add a new style on the user's behalf.

Example request
```JSON
{
    "name": "The name of the style",
    "description": "A wordy description of the new style",
    "notes": "Hey, if you install this make sure to remember X",
    "code": "Some actually valid UserCSS code here",
    "preview": "https://i.ytimg.com/vi/-Cv68B-F5B0/maxresdefault.jpg",
    "homepage": "https://github.com/userstyles-world/userstyles.world/",
    "license": "WTFPL",
    "category": "youtube.com"
}
```

Example response
```JSON
Status: 200
{
    "info": "Succesful added the style! With ID: 69"
}
```
or
```JSON
Status: 406
{
    "err": "Homepage validation is wrong."
}
```

### Delete a style

**Authorization is required**
```
DELETE /style/<id>
```
Deletes an style of the user.

Example response
```JSON
Status: 200
{
    "info": "Succesful removed the style!"
}
```
or
```JSON
Status: 403
{
    "err": "This style doesn't belong to you! ╰༼⇀︿⇀༽つ-]═──"
}
```

### Retrieve strong ETAG of style
```
HEAD /style/<id>.user.css
```
Returns strong ETAG header of specified style.

### Get style's source code
```
GET /style/<id>.user.css
```
Returns raw source code of specified style,  
this is normally used as updateURL and installation URL.

### Get style's preview
```
GET /style/preview/<id>.<webp | jpeg>
```
Gets an optimized webp/jpeg version of the style's preview.

### List all styles
```
GET /index/<format?>
```
List all available styles on USw.
When format(_optional_) is specified it will go through formating:
```JSON
`uso-format`
Special format that the USo-Archive uses as well. 
Meant for stylus
[
    {
        "i": 1,
        "n": "UserStyles.world Tweaks",
        "c": "userstyles.world",
        "u": 1617746048,
        "t": 5,
        "w": 5,
        "an": "vednoc",
        "sn": "https://userstyles.world/api/style/preview/1.webp"
    }
]
```

### Search for style
```
GET /search/<query>
```
Retrieve all styles based on the text query.

Example response query=userstyles-world
```json
[
    {
      "id": "29",
      "username": "Freeplay",
      "name": "Userstyles.World Tweaks!",
      "description": "Tweaks for Userstyles.World!",
      "preview": "https://codeberg.org/Freeplay/UserStyles/raw/branch/main/images/userstyles.world_tweaks-1.png",
      "notes": ""
    },
    {
      "id": "1",
      "username": "vednoc",
      "name": "UserStyles.world Tweaks",
      "description": "Various tweaks for UserStyles.world that might end up being upstreamed.",
      "preview": "https://user-images.githubusercontent.com/18245694/111090605-9cc4ce00-8530-11eb-8716-2c1df40f32e1.png",
      "notes": "After you install this userstyle, you can configure it from _configuration menu_ in Stylus' popout menu by clicking on the icon on your toolbar."
    },
    {
      "id": "36",
      "username": "Freeplay",
      "name": "Cleaner Stylus",
      "description": "Tweaks for Userstyles.World!",
      "preview": "https://codeberg.org/Freeplay/UserStyles/raw/branch/main/images/cleaner-stylus-1.png",
      "notes": ""
    }
]
```

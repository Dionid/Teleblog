# Teleblog

1. ...

# Stack

1. Go
1. Pocketbase
1. Templ
1. Vue
1. Tailwind
1. daisyUI

# Word of caution

This project is NOT about best practices. It's about making product
and do it efficiently. I haven't been working with Vue for a long time,
and this is first time for me to use Pocketbase, Templ.

Don't take this project as a reference for best practices.

# Project structure

1. `cmd/teleblog` - Teleblog platform
1. `infra` - some infrastructure code
1. `libs` - libraries

# Roadmap

## First phase

MG: Make it so content appears, but customization through Pocketbase admin

1. ~~Verify token~~
1. ~~Check user chat id is the same in channels and groups~~ -> They are not...
1. ~~Recover middleware to bot~~
1. ~~Save new posts~~
1. ~~Save new comments~~
1. ~~Update posts on edit~~
1. ~~Update comments on edit~~
1. ~~Pagination~~
1. ~~Post page~~
    1. ~~Add original post and comments links~~
    1. ~~Add comments authors data~~
        1. ~~From channels~~
        1. ~~Check that personal comments work also~~
1. ~~Post Widget~~
    1. ~~Add "Read more" if there is a comments~~
    1. ~~Add "Expand" if text is bigger than 200 symbols~~
    1. ~~Original post link~~
1. ~~Menu (throw template)~~
1. ~~Hero~~
1. ~~Mobile version~~
1. ~~Entities Markup~~
1. ~~Add license~~
1. ~~Rename "saas" to "teleblog"~~
1. Load history
1. Favicon
1. Deploy

## Second phase

MG: Add content improvement features

1. Images
1. Videos
1. Search
1. Likes counter
1. Extract post title
1. Extract tags
1. Pined messages
1. Links preview
1. Sorting
1. Empty page
1. ...

## Third phase

MG: ...

1. Command "Rebind Channels Group"
1. ...

## X phase

1. User Images
1. Admin page
1. Partial reload
1. Theme changer
1. ...
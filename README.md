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

1. `cmd/saas` - Teleblog SaaS platform
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
1. Update posts on edit
1. Update comments on edit
1. Deleted posts and comments
1. Pagination
1. Add comment authors
1. Post Widget
    1. Add "Read more" if there is a comments
    1. Add "Expand" if text is bigger than 200 symbols
    1. ...
1. Menu
    1. Logo
    1. Elements
    1. Links
    1. ...
1. Mobile version
1. Deploy

## Second phase

MG: Add content improvement features

1. Search
1. Likes counter
1. Extract post title
1. Extract tags
1. Pined messages
1. Links preview
1. Images
1. ...

## Third phase

MG: ...

1. Command "Rebind Channels Group"
1. ...

## X phase

1. Admin page
1. Partial reload
1. Theme changer
1. ...
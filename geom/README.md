## Description

`geom` is under construction. I'm trying to create a small [Golang] framework
which will help playing arround with 2D primitives, intersections, collisions,
etc.

Eventually it can become an impulse physics engine. Not [Box2D] of course, but
hopefully still something pleasant to toy around with.

[Box2D]: https://box2d.org/

## Running

`go run github.com/yarcat/playground/geom/demo`

## Techs Used

I use [Ebiten] for the graphics. It's dead simple!

With Go version 1.11+ the dependency should be fetched automagically. Otherwise
`go get github.com/hajimehoshi/ebiten` is your command. Or follow the
[installation instructions] for more details around possible dependencies.

I do the development on my [Chromebook Pixel 2] using [Linux (beta)], [Debian Linux]
and [Windows 10]. So it should work pretty much everywhere.

[Ebiten]: https://ebiten.org/
[installation instructions]: https://ebiten.org/documents/install.html
[Chromebook Pixel 2]: https://www.google.com/search?q=chromebook+pixel+2&rlz=1CATSID_enCH788CH788&sxsrf=ALeKk00XlfhL6ziH6EeE_8UxU_HEnn62yQ:1588155834763&source=lnms&tbm=isch&sa=X&ved=2ahUKEwjt94HDtY3pAhVR-6QKHQouAPIQ_AUoAXoECAwQAw&biw=1600&bih=991
[Linux (beta)]: https://support.google.com/chromebook/answer/9145439?hl=en
[Windows 10]: https://en.wikipedia.org/wiki/Windows_10

# Pages

Each of these represents a page that can be displayed to the user.

Each page must contain at least the following:

```javascript

class PageNAME {

  constructor(alerts) {
    this.alerts = alerts
  }

  setSelected() {
    console.log("PAGE NAME set to selected")
  }

  setIdle() {
    console.log("PAGE NAME set to idle")
  }

  render(contentTag) {
    console.log("Need to use the given content tag to draw data: " + contentTag)
  }
}
```

## constructor

Takes in "alerts" which is how the page can alert the user for errors, info, warning, etc.
Later on, this will most likely take more information


## setSelected

Indicates that the user has chosen to view the page. This is NOT a request to
display the page, rather, its to indicate that very very shortly there will
most likely be a request to the `render` function (below). The reason its not a guarantee
is that a user may delect multiple buttons in the UI within the render loop.

## setIdle

Indicates that the user has selected a different page. Its possible that this happened
between `setSelected` and `render`, but either way, it means that the user has unselected
the page

## render

Takes in `contentTag` which is the hook to which `div` that the page can use
to draw its data. This function should be lightweight, as it is potentially
ticked multiple times per-second, and any hold-up here may cause confusion and
terror to the user

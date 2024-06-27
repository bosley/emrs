# Pages

Each of these represents a page that can be displayed to the user.

Each page must contain at least the following:

```javascript
class PageTerminal {
  constructor(alerts) {
    this.alerts = alerts
    this.selected = false
  }

  setIdle() {
    this.selected = false
  }

  setSelected(contentHook) {
    this.selected = true
    $(contentHook).html("term")
  }
}

```

## constructor

Takes in "alerts" which is how the page can alert the user for errors, info, warning, etc.
Later on, this will most likely take more information


## setSelected

Indicates that the user has selected the UI element for the page. If the item is already selected, we may want to not update, or we may want to. Depends on the contents of the page.

The contentHook passed in is the UI element that will be drawn to

## setIdle

Indicates that the user has selected a different page.

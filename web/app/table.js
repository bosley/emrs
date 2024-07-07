

class Table {
  constructor(titles, rows) {
    this.titles = titles
    this.data = `
        <table>
          <thead>
            <tr>` 

    for (let i = 0; i < this.titles.length; i++) {
        this.data = this.data + `
              <th>` + this.titles[i] + `</th>`
    }

    this.data = this.data + `
            </tr>
            </thead>
            <tbody>`
    for (let row = 0; row < rows.length; row++) {
      this.data = this.data + `
              <tr>`
      for (let col = 0; col < this.titles.length; col++) {
        this.data = this.data + `
              <td>` + rows[row][col] + `</td>`
      }
      this.data = this.data + `
              </tr>`
    }

    this.data = this.data + `
            </body>
            </table>`
  }

  value() {
    return this.data
  }
}

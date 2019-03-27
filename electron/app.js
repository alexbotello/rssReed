const {app, BrowserWindow} = require('electron') // http://electron.atom.io/docs/api
const axios = require("axios")
const url = require('url');
const path = require('path')

let window = null



// Wait until the app is ready
app.once('ready', () => {
  // Create a new window
  window = new BrowserWindow({
    // Set the initial width to 800px
    width: 1400,
    // Set the initial height to 600px
    height: 1000,
    // Don't show the window until it ready, this prevents any white flickering
    show: false,
    webPreferences: {
      // Disable node integration in remote page
      nodeIntegration: false
    }
  })
  window.loadURL(url.format({
    pathname: path.join(__dirname, 'index.html'),
    protocol: 'file:',
    slashes: true
  }))

  window.once('ready-to-show', () => {
    window.show()
  })
})

app.on("window-all-closed", () => {
    axios.get("http://localhost:5000/exit")
        .then(() => {
            app.quit()
        })
        .catch(err => {
            app.quit()
        })
})
const {
  app,
  BrowserWindow,
  remote
} = require('electron')
const { execFile } = require('child_process')
const WebSocket = require("ws")
const shell = require('electron').shell;
const axios = require("axios")
const url = require('url');
const path = require('path')

let window = null

// const server = execFile("./server", (err, stdout, stderr) => {
//   if (err) {
//     console.log(err)
//     return
//   }
//   console.log(stdout)
// })

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
    icon: path.join(__dirname, "ico.png"),
    webPreferences: {
      // Disable node integration in remote page
      nodeIntegration: true
    }
  })

  // setTimeout(() => {
  //     window.socket = new WebSocket("ws://localhost:5000/stream")
  //     window.socket.setMaxListeners(1)
  //     window.loadURL(url.format({
  //       pathname: path.join(__dirname, 'index.html'),
  //       protocol: 'file:',
  //       slashes: true
  //     }))
  // }, 1000)
  window.socket = new WebSocket("ws://localhost:5000/stream")
  window.socket.setMaxListeners(1)
  window.loadURL(url.format({
    pathname: path.join(__dirname, 'index.html'),
    protocol: 'file:',
    slashes: true
  }))


  window.once('ready-to-show', () => {
    window.show()
  })

  function isSafeishURL(url) {
    return url.startsWith('http:') || url.startsWith('https:');
  }

  window.webContents.on("will-navigate", function (event, url) {
    event.preventDefault()
    if (isSafeishURL(url)) {
      shell.openExternal(url);
    }
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
$(function () {
    const { remote } = require("electron")
    const moment = require("moment")
    var socket = remote.getCurrentWindow()["socket"]
    var items = $("#container");
    socket.on("message", (data) => {
        data = JSON.parse(data)
        var rssItem = $("<div>").addClass("rssItem newItem")
        let link = $("<a>").attr("href", data.Link)
        let title = $("<h3>").text(data.Title)
        link.append(title)
        let date = $("<span>").text(moment(data.Date).format("MMMM Do YYYY"))
        rssItem.append(link, date)
        rssItem.prependTo(items)
    })
    // }
});
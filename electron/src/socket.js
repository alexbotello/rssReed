function parseFeedName(data) {
    let str
    if (typeof data === "string" || data instanceof String) {
        str = data
    } else {
        if (data.Name === "") {
            str = data.URL
        } else {
            str = data.Name
        }
    }
    str = str.split("-")[0]
    str = str.replace("www.", "")
    str = str.replace(".com", "")
    str = str.split(".org")[0]
    str = str.replace("https://", "")
    if (str.length < 5) {
        str = str.toUpperCase()
    }
    if (str.includes(":")) {
        // It's a reddit feed
        str = str.split(":")[1].replace(" ", "")
        str = "r/"+str
    }
    return str
}

$(function () {
    const { remote } = require("electron")
    const moment = require("moment")
    var socket = remote.getCurrentWindow()["socket"]
    var items = $("#container");
    var parseTitle;
    var sourceName;
    socket.on("message", (data) => {
        data = JSON.parse(data)
        var rssItem = $("<div>").addClass("rssItem newItem")
        var meta = $("<div>").addClass("meta")
        let link = $("<a>").attr("href", data.Link)
        if (data.Title.length > 60) {
            parseTitle = data.Title.slice(0, 60) + "..."
        } else {
            parseTitle = data.Title
        }
        if (data.Source === null) {
            sourceName = data.Link
        } else {
            sourceName = data.Source
        }

        let title = $("<h3>").text(parseTitle)
        link.append(title)
        let source = $("<p>").text(parseFeedName(sourceName))
        let date = $("<span>").text(moment(data.Date).format("MMMM Do YYYY"))
        meta.append(source, date)
        rssItem.append(link, meta)
        rssItem.prependTo(items)
    })
    // }
});
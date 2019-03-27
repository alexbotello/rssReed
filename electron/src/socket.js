function sleep(ms)
{
    var date = new Date();
    var curDate = null;
    do { curDate = new Date(); }
    while(curDate-date < ms);
}

$(function () {
    var socket = null;
    var items = $(".container");
    var div
    if (!window["WebSocket"]) {
        alert("Error: Your browser does not support web sockets.")
        console.log("Websocket can't connect")
    } else {
        sleep(12000) // Sleep for 15 seconds until initial load is complete
        socket = new WebSocket("ws://localhost:5000/stream");
        socket.onclose = function () {
            alert("Connection has been closed.");
        }
        socket.onmessage = function (e) {
            let data = JSON.parse(e.data)
            div = $("<div>").addClass("rssItem")
            let title = $("<h3>").text(data.Title)
            let img = $("<img>").attr("src", data.Image).addClass("rssImage")
            // let link = data[i].Link $("<>")
            let desc = $("<p>").text(data.Desc)
            let date = $("<span>").text(data.Date)
            div.append(img, title, desc, date)
            div.prependTo(items)
        }
    }
});
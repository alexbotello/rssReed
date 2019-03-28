function sleep(ms)
{
    var date = new Date();
    var curDate = null;
    do { curDate = new Date(); }
    while(curDate-date < ms);
}

function formatDate(date) {
    date = new Date(date)
    var hours = date.getHours();
    var minutes = date.getMinutes();
    var ampm = hours >= 12 ? 'pm' : 'am';
    hours = hours % 12;
    hours = hours ? hours : 12; // the hour '0' should be '12'
    minutes = minutes < 10 ? '0'+minutes : minutes;
    var strTime = hours + ':' + minutes + ' ' + ampm;
    return date.getMonth()+1 + "/" + date.getDate() + "/" + date.getFullYear() + "  " + strTime;
  }

$(function () {
    var socket = null;
    var items = $("#container");
    if (!window["WebSocket"]) {
        alert("Error: Your browser does not support web sockets.")
        console.log("Websocket can't connect")
    } else {
        // sleep(5000) // Sleep for 15 seconds until initial load is complete
        socket = new WebSocket("ws://localhost:5000/stream");
        socket.onclose = function () {
            alert("Connection has been closed.");
        }
        socket.onmessage = function (e) {
            let data = JSON.parse(e.data)
            var rssItem = $("<div>").addClass("rssItem newItem")
            let link = $("<a>").attr("href", data.Link)
            let title = $("<h3>").text(data.Title)
            link.append(title)
            let date = $("<span>").text(formatDate(data.Date))
            rssItem.append(link, date)
            rssItem.prependTo(items)
        }
    }
});
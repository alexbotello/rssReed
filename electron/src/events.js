$(function () {
    let feeds = document.getElementById("feeds")
    var items = $("#container");

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
      };

    function populateItems(data) {
        for (let i = 0; i < data.length; i++) {
            var rssItem = $("<div>").addClass("rssItem newItem")
            let link = $("<a>").attr("href", data[i].Link)
            let title = $("<h3>").text(data[i].Title)
            link.append(title)
            let date = $("<span>").text(formatDate(data[i].Date))
            rssItem.append(link, date)
            rssItem.prependTo(items)
        }
    };

    function populateFeed(data) {
        for (let i = 0; i < data.length; i++) {
            let li = $("<li>").text(data[i].Name).addClass("Feed")
            li.appendTo(feeds)
        }

    }

    fetch("http://localhost:5000/")
        .then(resp => {
            resp.json()
                .then(data => {
                    respFeeds = data.Feed
                    respItems = data.Item
                    populateItems(respItems)
                    populateFeed(respFeeds)
                })
                .catch(err => {
                    console.log(err)
                })
        })
        .catch(err => {
            console.log(err)
        })
    // Remove green border when new item is touched
    let container = document.getElementById("container")
    container.addEventListener("mouseover", (event) => {
        event.target.classList.remove("newItem")
    })
    container.addEventListener("mouseenter", (event) => {
        event.target.classList.remove("newItem")
    })
    container.addEventListener("mouseenter", (event) => {
        event.target.classList.remove("newItem")
    })
});

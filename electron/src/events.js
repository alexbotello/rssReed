$(function () {
    let container = document.getElementById("container")
    let feeds = document.getElementById("feeds")
    let plus = document.getElementById("plus")
    let label = document.getElementById("label")
    let button = document.getElementById("button")
    let input = document.getElementById("input")
    let items = $("#container");

    // Helper functions
    function formatDate(date) {
        date = new Date(date)
        var hours = date.getHours();
        var minutes = date.getMinutes();
        var ampm = hours >= 12 ? 'pm' : 'am';
        hours = hours % 12;
        hours = hours ? hours : 12; // the hour '0' should be '12'
        minutes = minutes < 10 ? '0' + minutes : minutes;
        var strTime = hours + ':' + minutes + ' ' + ampm;
        return date.getMonth() + 1 + "/" + date.getDate() + "/" + date.getFullYear() + "  " + strTime;
    };

    function sleep(ms) {
        var date = new Date();
        var curDate = null;
        do {
            curDate = new Date();
        }
        while (curDate - date < ms);
    }

    function populateItems(data) {
        for (let i = 0; i < data.length; i++) {
            var rssItem = $("<div>").addClass("rssItem")
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
            let name = parseFeedName(data[i])
            let li = $("<li>").text(name.toUpperCase()).addClass("Feed")
            li.appendTo(feeds)
        }
    }

    function parseFeedName(obj) {
        let str
        if (obj.Name === "") {
            str = obj.URL
        } else {
            str = obj.Name
        }
        let name = str.split("-")[0]
        name = name.replace("www.", "")
        name = name.replace(".com", "")
        name = name.split(".org")[0]
        name = name.replace("https://", "")
        return name
    }

    // Populate Feed sidebar and any feed items currenty in the database
    sleep(2000)
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
    container.addEventListener("mouseover", (event) => {
        event.target.classList.remove("newItem")
    })
    container.addEventListener("mouseenter", (event) => {
        event.target.classList.remove("newItem")
    })
    container.addEventListener("mouseleave", (event) => {
        event.target.classList.remove("newItem")
    })

    // Display tooltip label when Add button is hovered over
    plus.addEventListener("mouseover", (event) => {
        label.classList.add("visible")
    })
    plus.addEventListener("mouseleave", (event) => {
        label.classList.remove("visible")
    })

    // Open popup on plus button click
    plus.addEventListener("click", (event) => {
        popup.classList.add("visible")
    })

    // Submit new rss feed to server
    button.addEventListener("click", (event) => {
        popup.classList.remove("visible")
        let data = {"Url": input.value}
        input.value = ""
        fetch("http://localhost:5000/save", {
            method: "post",
            body: JSON.stringify(data)
        })
        .then(resp => {
            console.log(resp)
        })
        .catch(err => {
            console.log(err)
        })
    })

});
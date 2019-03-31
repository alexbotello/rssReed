$(function () {
    const moment = require("moment")

    let container = document.getElementById("container")
    let feeds = document.getElementById("feeds")
    let plus = document.getElementById("plus")
    let label = document.getElementById("label")
    let button = document.getElementById("button")
    let input = document.getElementById("input")
    let close = document.getElementById("popup-close")
    let items = $("#container");

    function populateItems(data) {
        for (let i = 0; i < data.length; i++) {
            var parseTitle;
            var sourceName;
            var rssItem = $("<div>").addClass("rssItem")
            var meta = $("<div>").addClass("meta")
            let link = $("<a>").attr("href", data[i].Link)
            if (data[i].Title.length > 60) {
                parseTitle = data[i].Title.slice(0, 60) + "..."
            } else {
                parseTitle = data[i].Title
            }

            if (data[i].Source === null) {
                sourceName = data[i].Link
            } else {
                sourceName = data[i].Source
            }
            let title = $("<h3>").text(parseTitle)
            link.append(title)
            let source = $("<p>").text(parseFeedName(sourceName))
            let date = $("<span>").text(moment(data[i].Date).format("MMMM Do YYYY"))
            meta.append(source, date)
            rssItem.append(link, meta)
            rssItem.prependTo(items)
        }
    };

    function populateFeed(data) {
        for (let i = 0; i < data.length; i++) {
            let name = parseFeedName(data[i])
            let p = $("<p>").text(name.toUpperCase()).addClass("Feed")
            p.appendTo(feeds)
        }
    }

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
        if (str.includes(":")) {
            // It's a reddit feed
            str = str.split(":")[1].replace(" ", "")
            str = "r/" + str
        }
        console.log(str)
        return str
    }

    function isValid(feed) {
        if (feed.includes("'")) return false
        if (feed.includes('"')) return false
        if (feed.length < 10) return false
        return true
    }

    function requestServerData() {
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
    }

    // Populate Feed sidebar and any feed items currenty in the database
    setTimeout(() => {
        requestServerData()
    }, 5000)

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

    // Remove error outline if refocusing on input
    input.addEventListener("focus", (event) => {
        input.style.outline = "none"
    })

    // Close popup if X button is clicked
    close.addEventListener('click', (event) => {
        popup.classList.remove("visible")
    })

    // Submit new rss feed to server
    button.addEventListener("click", (event) => {
        if (!isValid(input.value)) {
            input.style.outline = "1px solid #8B0000"
            return
        }
        popup.classList.remove("visible")
        let data = {
            "Url": input.value
        }
        input.value = ""
        fetch("http://localhost:5000/save", {
                method: "post",
                body: JSON.stringify(data)
            })
            .then(resp => {
                resp.json()
                    .then(data => console.log(data))
                    //     respFeeds = data.Feed
                    //     respItems = data.Item
                    //     populateItems(respItems)
                    //     populateFeed(respFeeds)
                    //   })
                    .catch(err => console.log(err))
            })
            .catch(err => {
                console.log(err)
            })
    })

});
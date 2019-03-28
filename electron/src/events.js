$(function () {
    fetch("http://localhost:5000/")
        .then(resp => {
            resp.json()
                .then(data => {
                    console.log(data)
                    let feeds = document.getElementById("feeds")
                    for (let i = 0; i < data.length; i++) {
                        console.log(data[i].URL)
                        let li = $("<li>").text(data[i].URL)
                        feeds.append(li)
                    }
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
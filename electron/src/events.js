$(function () {
    // Remove green border when new item is touched
    let container =  document.getElementById("container")
    container.addEventListener("mouseover", (event) => {
        event.target.classList.remove("newItem")
    })
});

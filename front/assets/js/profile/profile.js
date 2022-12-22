window.addEventListener("DOMContentLoaded", () => {
    getCurrentUser.then(() => {
        fetchPostsProfile();
    });
})

function fetchPostsProfile() {
    let userID = localStorage.getItem("userID")
    let headers = {
        "Authorization": localStorage.getItem("token"),
        "Content-Type": "application/json"
    }
    let xs = localStorage.getItem("xsession");
    if (xs) {
        headers["X-Session-Token"]= xs
    }
    axios({
        url: `${HOSTApi}/v1/users/${userID}/posts`,
        method: "GET",
        headers: headers
    }).then(rsp => {
        if (rsp.status == 200) {
            rsp.data.forEach(post => {
                document.getElementById("content").
                    insertAdjacentHTML("afterbegin", buildPost(post));
            });
        }
    })
}


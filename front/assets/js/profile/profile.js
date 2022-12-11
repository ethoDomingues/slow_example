window.addEventListener("DOMCotentLoaded", () => {
    let token = localStorage.getItem("token");
    getCurrentUser(token);
    fetchPostsProfile();
})

function fetchPostsProfile() {
    let userID = localStorage.getItem("userID")
    let token = localStorage.getItem("token")
    axios({
        url: `${HOSTAPI}/v1/users/${userID}/posts`,
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            Authorization: token,
        }
    }).then(rsp => {
        if (rsp.status == 200) {
            rsp.data.forEach(post => {
                document.getElementById("content").
                    insertAdjacentHTML("afterbegin", buildPost(post));
            });
        }
    })
}


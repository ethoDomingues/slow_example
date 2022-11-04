window.addEventListener("load", () =>{
    let token = localStorage.getItem("token");
    getCurrentUser(token);
    fetchPostsProfile();
})

function fetchPostsProfile() {
    let userID = localStorage.getItem("userID")
    let token = localStorage.getItem("token")
    fetch(`http://api.localhost:5000/v1/users/${userID}/posts`, {
        method:"GET",
        headers:{
            "Content-Type":"application/json",
            Authorization: token,
        }
    }).then(resp => {
        if (resp.status == 200 ){
            resp.json().then(data => {
                data.forEach(post => {
                document.getElementById("content").
                    insertAdjacentHTML("afterbegin",buildPost(post)); 
                });
            });
        }
    })
}


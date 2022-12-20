window.addEventListener("DOMContentLoaded", () => {
    getCurrentUser.then(() => {
        fetchPosts.then(() => {
            let form = document.getElementById("newpub-form");
            form.addEventListener("submit", pubNewPost, false);
        })
    }).catch(err => {
        console.log(err);
        logout();
    })
    return;
})

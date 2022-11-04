
window.addEventListener("load",()=>{
    let tkn = localStorage.getItem("token");
    if (tkn) {
        getCurrentUser(tkn);
    } else {
        logout();
    }
    fetchPosts();
    let form = document.getElementById("newpub-form");
    form.addEventListener("submit", pubNewPost, false);
    return;
})

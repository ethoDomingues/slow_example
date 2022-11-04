
function logout() {
    localStorage.clear();
    location.href = "/login";
    return;
}

function getCurrentUser(token) {
    fetch("http://api.localhost:5000/v1/users/whoami",{
        method:"GET",
        headers:{
            "Authorization": token,
        }
    }).then(resp => {
        if (resp.status == 401 || resp.status == 404 ) {
            logout();
        }
        return resp.json();
    }).then(data => {
        localStorage.setItem("token", data.token);
        localStorage.setItem("user", JSON.stringify(data.user));
        localStorage.setItem("userID", data.user.id);
    });
}

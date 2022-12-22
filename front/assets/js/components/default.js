
function logout() {
    localStorage.clear();
    let n = btoa(location.href);
    location.href = `${HOSTAuth}/v1/auth?next=${n}`;
    return;
}

function userLogout() {
    if (confirm("Deseja realmente sair?")) {
        localStorage.clear();
        location.href = `${HOSTAuth}/v1/auth`;
    }
    return;
}

const getCurrentUser = new Promise((resolve, reject) => {
    let u = localStorage.getItem("user");
    if (u != null) {
        let user = JSON.parse(u);
        resolve(user);
    }
    let reloadReq = false;
    let tkn = localStorage.getItem("token")
    if (!tkn) {
        let params = new Proxy(new URLSearchParams(window.location.search), {
            get: (searchParams, prop) => searchParams.get(prop),
        });
        if (params.token) {
            tkn = params.token;
            localStorage.setItem("token", tkn)
            reloadReq = true;
        }
    }
    if (!tkn) { reject("missing token") }
    let headers = {
        "Authorization": tkn,
        "Content-Type": "application/json"
    }
    let xs = localStorage.getItem("xsession");
    if (xs) {
        headers["X-Session-Token"] = xs;
    }
    axios({
        url: `${HOSTApi}/v1/users/whoami`,
        headers: headers
    }).then(rsp => {
        if (rsp.status == 401 || rsp.status == 404) {
            logout();
        }
        if (rsp.data.id) {
            let user = rsp.data;
            localStorage.setItem("userID", user.id);
            localStorage.setItem("user", JSON.stringify(user));
            let xs = rsp.headers["x-session-token"];
            if (xs) {
                localStorage.setItem("xsession",xs);
            }
            if (reloadReq) {
                console.log(location.pathname)
                location.href = location.pathname;
            }
            resolve(user);
        }
        reject("without user")
    }).catch(err => {
        reject(err)
    });
})

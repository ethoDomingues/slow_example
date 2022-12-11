
function logout() {
    localStorage.clear();
    let n = btoa(location.href);
    location.href = `${HOSTAUTH}/v1/auth?next=${n}`;
    return;
}

function userLogout() {
    if (confirm("Deseja realmente sair?")) {
        localStorage.clear();
        location.href = `${HOSTAUTH}/v1/auth`;
    }
    return;
}

const getCurrentUser = new Promise((resolve, reject) => {
    let u = localStorage.getItem("user");
    if (u != null) {
        let user = JSON.parse(u);
        resolve(user);
    }
    let tkn = localStorage.getItem("token");
    let reloadReq = false;
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
    axios({
        url: `${HOSTAPI}/v1/users/whoami`,
        headers: {
            "Authorization": tkn,
            "Content-Type": "application/json"
        }
    }).then(rsp => {
        if (rsp.status == 401 || rsp.status == 404) {
            logout();
        }
        if (rsp.data.id) {
            let user = rsp.data;
            localStorage.setItem("userID", user.id);
            localStorage.setItem("user", JSON.stringify(user));
            localStorage.setItem("xsession", rsp.headers["x-session-token"]);

            if (reloadReq) {
                location.href = location.pathname;
            }
            resolve(user);
        }
        reject("without user")
    }).catch(err => {
        reject(err)
    });
})

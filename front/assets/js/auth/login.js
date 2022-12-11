window.addEventListener("load", () => {
    let userStr = localStorage.getItem("user")
    let tokenStr = localStorage.getItem("token")
    if (userStr || tokenStr) {
        location.href = "/";
    }
    return;
})

var formLogin = document.getElementById("form-login")
var formRegister = document.getElementById("form-register")

formLogin.addEventListener("submit", login, false);
formRegister.addEventListener("submit", register, false);

function login(event) {
    event.preventDefault()
    let user = document.getElementById("input-login-user");
    let keep = document.getElementById("input-login-keep");
    let passwd = document.getElementById("input-login-passwd");

    if (user.value != "" && passwd.value != "") {
        axios({
            url: `${HOSTAUTH}/v1/login`,
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Basic " + btoa(`${user.value}:${passwd.value}`),
            },
            data: { "keep": keep.checked }
        }).then((rsp) => {
            if (rsp.status == 200) {
                let token = rsp.data["token"];
                localStorage.setItem("token", token);
                if (rsp.data["location"] != null) {
                    location.href = rsp.data["location"];
                }
                location.href = `/`;
            }
        })
    }
}

function register(event) {
    event.preventDefault()
    let name = document.getElementById("input-register-name");
    let user = document.getElementById("input-register-user");
    let keep = document.getElementById("input-register-keep");
    let passwd = document.getElementById("input-register-passwd");

    if (user.value != "" && passwd.value != "") {
        axios({
            url: "http://auth.localhost:5000/v1/register",
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Basic " + btoa(`${user.value}:${passwd.value}`),
            },
            data: {
                "keep": keep.checked,
                "name": name.value
            }
        }).then((rsp) => {
            if (rsp.status == 200) {
                let token = rsp.data["token"];
                localStorage.setItem("token", token);
                if (rsp.data["location"] != null) {
                    location.href = rsp.data["location"];
                }
                location.href = `/`;
            }
        });
    }
    return;
}

window.addEventListener("load",()=>{
    let userStr = localStorage.getItem("user")
    let tokenStr = localStorage.getItem("token")
    if (userStr || tokenStr) {
        location.href = "/";
    }
    return;
})

var formLogin = document.getElementById("form-login")
var formRegister = document.getElementById("form-register")

formLogin.addEventListener("submit",login, false);
formRegister.addEventListener("submit",register, false);

function login(event) {
    event.preventDefault()
    let user = document.getElementById("input-login-user");
    let keep = document.getElementById("input-login-keep");
    let passwd = document.getElementById("input-login-passwd");
    
    if(user.value != "" && passwd.value !="") {
        fetch("http://auth.localhost:5000/v1/login", {
            method:"POST",
            headers:{
                "Content-Type":"application/json",
                "Authorization": "Basic "+btoa(`${user.value}:${passwd.value}`),
            },
            body: JSON.stringify({"keep":keep.checked})
        })
        .then((resp) => {
            if (resp.status == 200) {
                resp.json().then((data) => {
                    let token = data["token"];
                    localStorage.setItem("token",token);
                    if( data["location"] != null ){
                        location.href = data["location"];
                    }
                    location.href = `/`;
                });
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
    
    if(user.value != "" && passwd.value !="") {
        fetch("http://auth.localhost:5000/v1/register", {
            method:"POST",
            headers:{
                "Content-Type":"application/json",
                "Authorization": "Basic "+btoa(`${user.value}:${passwd.value}`),
            },
            body: JSON.stringify({
                "keep":keep.checked,
                "name":name.value
            })
        })
        .then((resp) => {
            if (resp.status == 200) {
                resp.json().then((data) => {
                    let token = data["token"];
                    localStorage.setItem("token",token);
                    if( data["location"] != null ){
                        location.href = data["location"];
                    }
                    location.href = `/`;
                });
            }
        });
    }
    return;
}

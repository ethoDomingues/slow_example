
function pubNewComm(postID) {
    let pElement = document.getElementById(postID);
    let text = pElement.getElementsByClassName("comments-form-text")[0];
    let image = pElement.getElementsByClassName("comments-form-image")[0];

    if (text.value || image.files.length>0) {
        let user = localStorage.getItem("userID");
        let token = localStorage.getItem("token");
        let fd = new FormData();

        fd.set("text", text.value);
        if(image.files.length >0){
            fd.set("image", image.files[0]);
        }

        fetch(`http://api.localhost:5000/v1/users/${user}/posts/${postID}/comments`,{
            method:"POST",
            headers:{
                "Authorization": token,
            },
            body:fd
        }).then(resp => {
            if(resp.status == 201 ){
                text.value = "";
                image.type = "text";
                image.type = "file";
                resp.json().then(data => {
                    pElement.getElementsByClassName("comments-content")[0].
                        insertAdjacentHTML("afterbegin",buildComm(data)) 
                });
            }
        });
    } else {
        console.log("nada :(")
    }
}

function buildComm(comm) {
    let image = "";
    if(comm.image) {
        image = `
        <picture class="comments-content-pictures">
                <img src="${comm.image.url}" alt="user-profile">
        </picture>`;
    }
    return `
        <div class="comments" id="${comm.id}" data='${JSON.stringify(comm)}'>
            ${buildPostHeader(comm)}
            <div class="comments-content">
                <div class="comment-content-text">
                    <pre>${comm.text}</pre>
                </div>
                <div>${image}</div>
            </div>
            <div class="comments-footer">
                <div onclick="reactComm('${comm.id}')" class="comments-footer-reacts">Reacts <span>${comm.reacts.length}</span></div>
            </div>
        </div>`;
}

function deleteComm(commID) {
    let user = localStorage.getItem("userID");
    let cElement = document.getElementById(commID)
    let commData = cElement.getAttribute("data");
    let comm = JSON.parse(commData);
    if (comm.owner.id == user){
        let token = localStorage.getItem("token");
        fetch(`http://api.localhost:5000/v1/users/${comm.post.owner.id}/posts/${comm.post.id}/comments/${commID}`,{
            method:"DELETE",
            headers:{
                Authorization:token,
            }
        }).then(resp =>{
            if (resp.status == 200 ||resp.status == 204) {
                cElement.remove();
            }
        })
    }
}

function reactComm(commID) {
    let c = document.getElementById(commID);
    if (c != null) {
        let comm = JSON.parse(c.getAttribute("data"));
        let tkn = localStorage.getItem("token");
        fetch(`http://api.localhost:5000/v1/users/${comm.post.owner.id}/posts/${comm.post.id}/comments/${comm.id}/reacts`,
        {
            method:"PUT",
            headers:{
                Authorization:tkn
            }
        }).
        then(resp => {
            if (resp.status == 200){
                return resp.json();
            }
        }).
        then(data => {
            let reactSpan = c.getElementsByClassName("comments-footer-reacts")[0];
            reactSpan.getElementsByTagName("span")[0].textContent = data.length;
        }).catch(err => {
            alert("Desculpe, algo deu errado... :(")
        } )
    }
}

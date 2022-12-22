
function pubNewComm(postID) {
    let pElement = document.getElementById(postID);
    let text = pElement.getElementsByClassName("comments-form-text")[0];
    let image = pElement.getElementsByClassName("comments-form-image")[0];

    if (text.value || image.files.length > 0) {
        let user = localStorage.getItem("userID");
        let fd = new FormData();
        fd.set("text", text.value);
        if (image.files.length > 0) {
            fd.set("image", image.files[0]);
        }

        let headers = {
            "Authorization": localStorage.getItem("token"),
            "Content-Type": "multipart/form-data"
        }
        let xs = localStorage.getItem("xsession");
        if (xs) {
            headers["X-Session-Token"]= xs
        }
        axios({
            url: `${HOSTApi}/v1/users/${user}/posts/${postID}/comments`,
            method: "POST",
            headers: headers,
            data: fd
        }).then(rsp => {
            if (rsp.status == 201) {
                text.value = "";
                image.type = "text";
                image.type = "file";

                let data = rsp.data;
                let cElement = pElement.getElementsByClassName("comments-content")[0];
                if (cElement.style.display == "none") {
                    cElement.style.display = "flex";
                }
                pElement.getElementsByClassName("comments-content")[0].
                    insertAdjacentHTML("afterbegin", buildComm(data))
            }
        });
    }
}

function deleteComm(commID) {
    let user = localStorage.getItem("userID");
    let cElement = document.getElementById(commID)
    let commData = cElement.getAttribute("data");
    let comm = JSON.parse(commData);
    if (comm.owner.id == user) {
        let headers = {
            "Authorization": localStorage.getItem("token"),
            "Content-Type": "application/json"
        }
        let xs = localStorage.getItem("xsession");
        if (xs) {
            headers["X-Session-Token"]= xs
        }
        axios({
            url: `${HOSTApi}/v1/users/${comm.post.owner.id}/posts/${comm.post.id}/comments/${commID}`,
            method: "DELETE",
            headers: headers,
        }).then(resp => {
            if (resp.status == 200 || resp.status == 204) {
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
        let headers = {
            "Authorization": tkn,
            "Content-Type": "application/json"
        }
        let xs = localStorage.getItem("xsession");
        if (xs) {
            headers["X-Session-Token"]= xs
        }
        axios({
            url: `${HOSTApi}/v1/users/${comm.post.owner.id}/posts/${comm.post.id}/comments/${comm.id}/reacts`,
            method: "PUT",
            headers: headers,
        }).then(rsp => {
            if (rsp.status == 200 || rsp.status == 201) {
                let data = rsp.data;
                let reactSpan = c.getElementsByClassName("comments-footer-reacts")[0];
                reactSpan.getElementsByTagName("span")[0].textContent = data.length;
            }
        }).catch(() => {
            alert("Desculpe, algo deu errado... :(")
        })
    }
}




function buildComm(comm) {
    let image = "";
    if (comm.image) {
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
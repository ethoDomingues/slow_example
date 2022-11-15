

function fetchPosts() {
    let tkn = localStorage.getItem("token")
    if (tkn == null ){
        location.href = "/login";
    }
    fetch("http://api.localhost:5000/v1/posts").
    then((resp) => {
        if (resp.status == 200) {
            return resp.json();
        }}).then((data) => {         
            let cttValue = "";
            let cttDiv = document.getElementById("content");
            data.forEach((p) => {
                cttDiv.insertAdjacentHTML("beforeend",buildPost(p));
            });
        })
}

function deletePost(postID) {
    let user = localStorage.getItem("userID");
    let pElement = document.getElementById(postID)
    let postData = pElement.getAttribute("data");
    let post = JSON.parse(postData);
    if (post.owner.id == user){
        let token = localStorage.getItem("token");
        fetch(`http://api.localhost:5000/v1/users/${post.owner.id}/posts/${postID}`,{
            method:"DELETE",
            headers:{
                Authorization:token,
            }
        }).then(resp =>{
            if (resp.status == 200 || resp.status == 204) {
                Array.from(document.getElementsByClassName(`posts-shared-content-${postID}`)).
                forEach(elem => {
                    elem.innerHTML = `
                    <div>
                        <b>O Conteudo Já Não Está Mais Disponível<b>
                    </div>`;
                })
                pElement.remove();
            }
        })
    }
}

function reactPost(postID) {
    let p = document.getElementById(postID);
    if (p != null) {
        let post = JSON.parse(p.getAttribute("data"));
        let tkn = localStorage.getItem("token");
        fetch(`http://api.localhost:5000/v1/users/${post.owner.id}/posts/${post.id}/reacts`,
        {
            method:"PUT",
            headers:{
                Authorization:tkn
            }
        }).
        then(resp => {
            if (resp.status == 200){
                return resp.json();
            } else if (resp.status == 401) {
                location.href = "/login";
            }
        }).
        then(data => {
            let reactSpan = p.getElementsByClassName("posts-footer-reacts")[0];
            reactSpan.getElementsByTagName("span")[0].textContent = data.length;
        }).catch(err => {
            console.log(err)
            alert("Desculpe, algo deu errado... :(")
        } )
    }
}

function newpubToggle() {
    let newpub = document.getElementById("newpub");
    if (newpub.style.display == "none") {
        newpub.style.display = "flex";
    } else {
        document.getElementById("newpub-text").value = "";
        let img = document.getElementById("newpub-images");
        img.type = "text";
        img.type = "file";
        newpub.style.display = "none";
    }
    return;
}

function openViewer(postID) {
    let viewer = document.getElementById("viewer");
    viewer.innerHTML = "";

    let pElement = document.getElementById(postID);
    let pData = pElement.getAttribute("data");
    let post = JSON.parse(pData);
    viewer.innerHTML = buildPost(post);

    viewer.style.display = "flex";
}

function closeViewer() {
    let viewer = document.getElementById("viewer")
    viewer.style.display = "none";
}

function toggleComments(postID) {
    let pElement = document.getElementById(postID);
    let cElement = pElement.getElementsByClassName("comments-content")[0];

    if(cElement.style.display == "none"){
        cElement.style.display = "flex";
    } else {
        cElement.style.display = "none";
    }
    return;
}

function buildPostHeader(p) {
    let url = "http://localhost:5000/assets/imgs/user.png";
    if (p.owner.profile) {
        url = p.owner.profile.url;
    } 
    let imOwner = p.owner.id == localStorage.getItem("userID");
    let del = "";
    if (p.id.startsWith("comms@")){
        if (!p.deleted && imOwner) {
            del = `<span onclick="deleteComm('${p.id}')">Excluir</span>`;
        }
    } else if (p.id.startsWith("posts@")){
        if (!p.deleted && imOwner) {
            del = `<span onclick="deletePost('${p.id}')">Excluir</span>`;
        }
    }

    return `
    <div class="posts-headers">
        <picture class="posts-headers-picture">
            <img src="${url}" width="50px" alt="user-profile" class="userProfile-${p.owner.id}">
        </picture>
        <div class="posts-headers-data">
            <span>${p.owner.name}</span>
            <span><sub>${moment(p.createdAt).fromNow()}</sub></span>
        </div>
        ${del}
    </div>`
}

function buildPostContent(p) {
    if (p.deleted) {
        return `
        <div class="posts-content">
            <div class="shared-posts-404">
                <b>O Conteudo Já Não Está Mais Disponível</b>
            </div>
        </div>`;
    }
    let content = ''
    if (p.shared != null) {
        content = `
        <div class="posts-content posts-shared-content-${p.id}">
            <div class="posts-content-text">
                <pre>${p.text}</pre>
            </div>
            <div class="shared-${p.shared.id} shared-posts">
                ${buildPostHeader(p.shared)}
                ${buildPostContent(p.shared)}
            </div>
        </div>`;
    } else {
        let images = "";
        p.images.forEach((img)=>{
            images += `
            <picture class="posts-content-pictures">
                <img src="${img.url}" alt="uma imagem">
            </picture>`;
        });
        content = `
            <div class="posts-content-text">
                <pre>${p.text}</pre>
            </div>
            <div class="posts-content-images">
                ${images}
            </div>`;
    }
    return content;
}

function buildPost(p) {
    
    let comms = "";
    p.comments.forEach(comm => { comms += buildComm(comm);});
    if (comms == "") {
        comms = "<span>100 comentarios...</span>";
    }

    return `
<article class="posts" data='${JSON.stringify(p)}' id="${p.id}" >
    ${buildPostHeader(p)}
    <div class="posts-content">
        ${buildPostContent(p)}
    </div>
    <div class="posts-footer">
        <div onclick="reactPost('${p.id}')" class="posts-footer-reacts">Reacts <span>${p.reacts.length}</span></div>
        <div onclick="toggleComments('${p.id}')" class="posts-footer-comm">Comment <span>${p.comments.length}</span></div>
        <div onclick="sharePost('${p.id}')" class="posts-footer-share">Share <span>${p.sharedCount}</span></div>
    </div>
    <div>
        <div>
            <textarea  class="comments-form-text" placeholder="Comentar"></textarea>
            <input type="file" id="comments-form-image-${p.id}" class="comments-form-image" hidden>
            <label for="comments-form-image-${p.id}">+ Image</label>
            <input type="submit" onclick="pubNewComm('${p.id}')" value="Comentar">
        </div>
        <div class="comments-content" style="display:none;">
            ${comms}
        </div>
    </div>
</article>`
}

function pubNewPost(event) {
    event.preventDefault();
    let text = document.getElementById("newpub-text");
    let images = document.getElementById("newpub-images");
    if (text.value || images.files.length>0) {
        let profile = document.getElementById("newpub-profile")

        let user = localStorage.getItem("userID");
        let token = localStorage.getItem("token");
        let fd = new FormData();

        fd.set("text", text.value);
        fd.set("profile", profile.checked);
        Array.from(images.files).forEach((file) => {
            fd.append("images",file)
        });

        fetch(`http://api.localhost:5000/v1/users/${user}/posts`,{
            method:"POST",
            headers:{
                "Authorization": token,
            },
            body:fd
        }).then(resp => {
            if(resp.status == 201 ){
                newpubToggle();
                text.value = "";
                images.type = "text";
                images.type = "file";
                resp.json().then(data => {
                    if(profile.checked){
                        console.log("tamo ae")
                        Array.from(document.getElementsByClassName("userProfile-"+data.owner.id)).
                        forEach(elem => {
                            elem.src = data.owner.profile.url;
                        });
                    }
                    document.getElementById("content").insertAdjacentHTML("afterbegin",buildPost(data)) 
                    profile.checked = false;
                });
            }
        });
    }
    return;
}

function sharePost(postID) {
    let userID = localStorage.getItem("userID")
    let token = localStorage.getItem("token")

    fetch(`http://api.localhost:5000/v1/users/${userID}/posts`, {
        method:"POST",
        headers:{
            "Authorization": token,
            "Content-Type":"application/json",
        },
        body: JSON.stringify({
            "shared":postID,
            text:""
        })
    }).then(resp => {
        if (resp.status == 201) {
            resp.json().then(data => {
                document.getElementById("content").
                    insertAdjacentHTML("afterbegin",buildPost(data));

            let pElement = document.getElementById(data.shared.id);
            pElement.getElementsByClassName("posts-footer-share")[0].
                    getElementsByTagName("span")[0].innerText = data.shared.sharedCount;
            });
        }
    });
    return;
}

function toggleInputProfile() {
    let dElement = document.getElementById("newpub-profile-div");
    let iElement = document.getElementById("newpub-images");
    if(iElement.files.length > 0 ){
        dElement.hidden = false;
    } else {
        dElement.hidden = true;
    }
    return;
}
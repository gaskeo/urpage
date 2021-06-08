function sendRegistration() {
    if (!checkPasswordsMatch()) {
        setError("passwords-not-match")
        return false
    }

    let status;

    let data = new FormData();

    if (document.getElementById("username").value.replace(/\s/g, "") === "" ||
        document.getElementById("email").value.replace(/\s/g, "") === "" ||
        document.getElementById("password").value.replace(/\s/g, "") === "") {

        setError("empty-input")
        return false
    }

    data.append("csrf", document.getElementById("csrf").value)
    data.append("username", document.getElementById("username").value);
    data.append("email", document.getElementById("email").value)
    data.append("password", document.getElementById("password").value)

    fetch("/do/registration", {method: 'post', body: data}).then(function (r) {
        status = r.status

        if (status === 200) {
            r.json().then(function (j) {
                if (j["Err"] === "") {
                    window.location.replace("/login")
                } else {
                    setError(j["Err"])
                }
            })
        } else {
            alert("что-то пошло не так...");
        }
    })

    return false
}

function sendLogin() {
    let status;

    let data = new FormData();

    if (document.getElementById("email").value.replace(/\s/g, "") === "" ||
        document.getElementById("password").value.replace(/\s/g, "") === "") {

        setError("empty-input")
        return false
    }

    data.append("email", document.getElementById("email").value)
    data.append("password", document.getElementById("password").value)

    fetch("/do/login", {method: 'post', body: data}).then(function (r) {
        status = r.status

        if (status === 200) {
            r.json().then(function (j) {
                if (j["Err"] === "") {
                    window.location.replace("/")
                } else {
                    setError(j["Err"])
                }
            })
        } else {
            alert("что-то пошло не так...");
        }
    })

    return false
}

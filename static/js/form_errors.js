
const errors = {
    "wrong-password": "Неправильный пароль",
    "user-not-exist": "Пользователя не существует",
    "passwords-not-match": "Пароли не совпадают",
    "email-exist": "Пользователь с такой почтой уже существует",
    "other-error": "Что-то пошло не так..."
};

function checkPasswordsMatch() {
    let password = document.getElementById("password").value;
    let passwordAgain = document.getElementById("password-again").value;

    return password === passwordAgain;
}

function getErrorP() {
    let errorsP = document.getElementsByClassName("form-error")

    let errorP

    for (let i = 0; i < errorsP.length; i++) {
        if (errorsP[i].style.display !== "none") {
            errorP = errorsP[i]
            return errorP
        }
    }
    return null

}

function setError(err) {
    let errorP = getErrorP();

    if (errorP === null) {
        console.log(1)
        return
    }
    if (err !== null) {
        console.log(2)

        for (let key in errors) {
            if (err === key) {
                errorP.textContent = errors[key];
                errorP.style.display = "block"
                return
            }
        }
    }

}
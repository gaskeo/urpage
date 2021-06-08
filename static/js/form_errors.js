const statusTexts = {
    "wrong-password": "Неправильный пароль",
    "user-not-exist": "Пользователя не существует",
    "passwords-not-match": "Пароли не совпадают",
    "email-exist": "Пользователь с такой почтой уже существует",
    "other-error": "Что-то пошло не так...",
    "empty-input": "Не все поля заполнены",
    "no-csrf": "Соединение потеряно, перезагрузите страницу",
    "ok": "Успешно"
};

const statusTextColors = {
    "wrong-password": "red",
    "user-not-exist": "red",
    "passwords-not-match": "red",
    "email-exist": "red",
    "other-error": "red",
    "empty-input": "red",
    "no-csrf": "red",
    "ok": "black"
}

function checkPasswordsMatch() {
    let password = document.getElementById("password").value;
    let passwordAgain = document.getElementById("password-again").value;

    return password === passwordAgain;
}

function getErrorP() {
    let errorsP = document.getElementsByClassName("form-error")

    let errorP

    for (let i = 0; i < errorsP.length; i++) {
        if (errorsP[i].offsetWidth > 0 && errorsP[i].offsetHeight > 0) {
            errorP = errorsP[i]
            return errorP
        }
    }
    return null

}

function setError(err) {
    let errorP = getErrorP();

    if (errorP === null) {
        return
    }
    if (err !== null) {

        for (let key in statusTexts) {
            if (err === key) {
                errorP.style.color = statusTextColors[key]

                errorP.textContent = statusTexts[key];
                return
            }
        }
    }

}
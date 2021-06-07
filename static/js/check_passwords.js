function checkPasswordsMatch() {
    let password = document.getElementById("password").value;
    let passwordAgain = document.getElementById("password-again").value;

    return password === passwordAgain;
}
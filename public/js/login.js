tokenIsValid(function(isValid) {
    if (!isValid) localStorage.removeItem("camstream-token");
    else window.location.href = window.location.origin + "/stream";
})

function authenticate() {
    console.log("Authenticating...");
    let password = document.getElementById("password").value;
    doHttpGet(`/authenticate?password=${encodeURIComponent(password)}`, function(response, status) {
        console.log(response + " - " + status);
        if (status == 401) {
            document.getElementById('auth-failed').style.display = 'block';
        } else if (status == 202) {
            window.localStorage.setItem("camstream-token", response);
            window.location.href = window.location.origin + "/stream";
        }
    });
}
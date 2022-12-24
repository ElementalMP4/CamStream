tokenIsValid(function(isValid) {
    if (!isValid) {
        localStorage.removeItem("camstream-token");
        window.location.href = window.location.origin + "/login";
    } else document.getElementById("stream").src = `/video?token=${localStorage.getItem("camstream-token")}`;
})
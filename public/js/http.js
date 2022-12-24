function doHttpGet(url, callback) {
    var xhr = new XMLHttpRequest();

    xhr.onreadystatechange = function() {
        if (xhr.readyState == xhr.DONE) {
            callback(xhr.responseText, xhr.status);
        }
    }
    xhr.open("GET", url, true);
    xhr.send();
}

function tokenIsValid(callback) {
    if (localStorage.hasOwnProperty("camstream-token")) {
        var token = localStorage.getItem("camstream-token");
        doHttpGet(`/validate?token=${token}`, function(response, status) {
            let valid = status == 202
            callback(valid);
        });
    }
}
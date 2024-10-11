var data = {
    'path': window.location.pathname,
    'domain': window.location.host,
    'referer': document.referer,
}

var xhr = new XMLHttpRequest();
xhr.open("POST", 'https://map.invictadux.com/api/v2/ping', true);

xhr.setRequestHeader("Content-Type", "application/json;");

xhr.onload = function () {
};

xhr.onerror = function () {
};

xhr.send(JSON.stringify(data));
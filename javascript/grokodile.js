var grokodile = function(){
    var http = new XMLHttpRequest();
    http.withCredentials = true;
    http.open("GET", "localhost:8000", true);
    http.send();
}

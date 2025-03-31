const socket = new WebSocket("ws://"+location.host+"/endpoint/");

socket.onopen = function(e) {
    console.log('Connection Estabilished');
    socket.send('Hello Server!');
};

socket.onmessage = function(event) {
    console.log(`Dados recebidos: ${event.data}`);
};

socket.onclose = function(event) {
    console.log('Conexão fechada');
};

socket.onerror = function(error) {
    console.log('Erro:', error);
};

function sendToServer() {
    fields = getFields()
    socket.send(fields)
}

//Get input fields 
function getFields() {
    let username = document.getElementById("username-input").value
    let message = document.getElementById("message-input").value

    return `{ "username": ${username ? `"`+ username + `"` : null}, "message" : ${message ? `"` + message + `"` : null} }`
}





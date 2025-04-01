const socket = new WebSocket("wss://"+location.host+"/endpoint/");
const chatBox = document.getElementById("chat")

socket.onopen = function(e) {
    console.log('Connection Estabilished');
    socket.send('Hello Server!');
};


socket.onmessage = function(event) {
    //Check if the current message is the chatlog:
    let msg = event.data
    if (msg.startsWith("!msg")){
        console.log("Message: "+event.data)
        msg = msg.replace("!msg", "")
        splitMsg = msg.split(":")
        //HTML element
        const elementString = `<p class="message-body"> <span class="username">${splitMsg[0]}:</span> ${splitMsg[1]}</p>`
        chatBox.innerHTML += elementString
    }


    //console.log(`Dados recebidos: ${event.data}`);
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





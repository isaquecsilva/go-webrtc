import PeerService from './rtcpeer.js'
import WebsocketService from './websocket.js';


function watch() {
    const socket = new WebSocket('/newconn')

    let ps = new PeerService(null, (candidate) => {
        socket.send(JSON.stringify({
            type: "candidate",
            data: candidate,
        }))
    }, (track) => {
        console.log('track received...', track)
        let video = document.querySelector('video');
        video.srcObject = track.streams[0];
    })
    
    let ws = new WebsocketService(socket, ps)

    setTimeout(() => {
        console.log('requesting to create offer')

        ws.socket.send(JSON.stringify({
            type: "createOffer",
            data: null,
        }))

    }, 5000)

}

const button = document.querySelector('button[class=watch]')
button.onclick = () => watch()
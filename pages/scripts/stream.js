import WebsocketService from './websocket.js'
import PeerService from './rtcpeer.js'

async function getStream() {
    const constraints = {
        video: true,
        audio: true,
    }

    let mediaStream = await navigator.mediaDevices.getDisplayMedia(constraints);
    return mediaStream
}

function createStream() {
    getStream()
        .then(stream => {
            document.querySelector('video').srcObject = stream
    
            let socket = new WebSocket('/newconn');

            let ps = new PeerService(stream, (candidate) => {
                socket.send(JSON.stringify({
                    type: 'candidate',
                    data: candidate,
                }))
            }, (track) => {})

            let ws = new WebsocketService(socket, ps)
            setTimeout(() => {
                console.log('requesting to be streamer')

                ws.socket.send(JSON.stringify({
                    type: "streamer",
                    data: null,
                }))
            }, 3000)
        })
        .catch(error => {
            alert(`media stream error: ${error?.message}`)
        })
}

const button = document.querySelector('button[class=stream]')
button.onclick = () => createStream()


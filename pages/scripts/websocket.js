export default class WebsocketService {
    constructor(socket, peerService) {
        this.peerService = peerService;
        this.socket = socket
        this.socket.onopen = this.openEvent
        this.socket.onclose = this.closeEvent
        this.socket.onerror = this.errorEvent;
        this.socket.onmessage = this.messageEvent.bind(this)
    }

    openEvent(event) {
        console.log(event)
    }

    closeEvent(event) {
        console.log(event)
    }

    errorEvent(event) {
        console.log(event)
    }

    messageEvent({ data }) {
        let message = JSON.parse(data)
        console.log(message.type)

        let handler = this[message.type]
        
        if (!handler) {
            console.error(`not found handler: ${message.type}`);
            return;
        }
        
        handler.call(this, message.data)
    }

    async createOffer() {
        console.log('creating offer...')
        let offer = await this.peerService.createOffer()
        console.log(offer)

        this.socket.send(JSON.stringify({
            type: 'offer',
            data: offer,
        }))

        this.socket.send(JSON.stringify({
            type: "createAnswer",
            data: null,
        }))
    }

    async createAnswer() {
        console.log('creating answer')
        let answer = await this.peerService.createAnswer()
        console.log(answer)

        this.socket.send(JSON.stringify({
            type: 'answer',
            data: answer,
        }))
    }

    async description(sdp) {
        await this.peerService.setRemoteDescription(sdp)
    }

    async candidate(iceCandidate) {
        await this.peerService.addCandidate(iceCandidate)
    }

    static async sendCandidate(candidate) {
        this.socket.send(JSON.stringify({
            type: "candidate",
            data: candidate,
        }))
    }
}
export default class PeerService {
    constructor(stream, socketCandidateCb, trackCb) {
        this.peer = new RTCPeerConnection();
        
        if (stream) {
            this.peer.addStream(stream)
        }

        this.onTrack(trackCb)
        this.onCandidate(socketCandidateCb)
    }

    onTrack(cb) {
        this.peer.ontrack = cb
    }

    onCandidate(socketCandidateCb) {
        this.peer.onicecandidate = ({ candidate }) => {
            if (candidate) {
                socketCandidateCb(candidate)
            }
        }
    }

    async createOffer() {
        let offer = await this.peer.createOffer()
        await this.peer.setLocalDescription(offer)
        return offer;
    }

    async createAnswer() {
        let answer = await this.peer.createAnswer();
        await this.peer.setLocalDescription(answer)
        return answer;
    }

    async setRemoteDescription(description) {
        console.log('REMOTE_DESCRIPTION:' + JSON.stringify(description));
        await this.peer.setRemoteDescription(description)
    }

    async addCandidate(iceCandidate) {
        let candidate = new RTCIceCandidate(iceCandidate)
        console.log(candidate);
        await this.peer.addIceCandidate(candidate)
    }
}
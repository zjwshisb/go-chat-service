export default class Client {
    url: string
    conn : WebSocket
    isClose: Number = 1
    constructor(url: string) {
        this.url = url;
        this.conn = new WebSocket(url);
        this.conn.addEventListener('error', (e: Event) =>  {
            this.isClose = 1;
            console.log(e);
        });
        this.conn.addEventListener('open', (e: Event) => {
            console.log(e);
            this.isClose = 0;
        });
        this.conn.addEventListener('message', (e: Event) => {
            console.log(e);
        });
    }
    connect() {
    }
    close() {
        this.isClose = 1;
        this.conn.close();
    }
    send(text: string) : void {
        this.conn.send(text);
    }
}
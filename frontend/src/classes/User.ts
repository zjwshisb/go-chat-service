export default class {
    constructor(id: number) {
        this.id = id;
    }
    id: number;
    messages: [] = [];
    isOnline: number = 0;
    hadNotReadMsg: number = 0;
}
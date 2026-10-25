export default function Base64ToBuffer(enc: string): Uint8Array<ArrayBuffer> {
    const dec = atob(enc);
    const arr = new Uint8Array(dec.length);

    for (let i = 0; i < dec.length; i++) {
        arr[i] = dec.charCodeAt(i);
    }
    
    return arr;
}
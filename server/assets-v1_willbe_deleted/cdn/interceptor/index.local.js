import * as interceptor from "./interceptor__local.json" assert { type: "json" };
await import("../wasm_exec_v1.js");

// Original -- Instantiate from WASM
// (async function(){
//     const go = new Go();
//     try{
//       await WebAssembly.instantiateStreaming(fetch("http://localhost:5000/assets/cdn/interceptor/interceptor__local.wasm"), go.importObject).then((result) => {
//         go.run(result.instance);
//     // await WebAssembly.instantiateStreaming(fetch("interceptor.wasm"), go.importObject).then((result) => {
//     //     go.run(result.instance);
//       })
//     } catch (err){
//       console.log(err)
//     }
//     console.log("IIFE WASM Loader Called")
//   })()


// Daniel's code -- instantiate from JSON
const decode = (encoded) => {
    var str = window.atob(encoded);
    var bytes = new Uint8Array(str.length);
    for (var i = 0; i < str.length; i++) {
        bytes[i] = str.charCodeAt(i);
    }
    return bytes.buffer;
}

(async () => {
    const go = new Go();
    const importObject = go.importObject;
    return await WebAssembly.instantiate(decode(interceptor.default), importObject).then((result) => {
        go.run(result.instance);
        globalThis.layer8 = layer8;
        //globalThis.BACKEND = "localhost:8000"
        console.log("globalThis.BACKEND: ", globalThis.BACKEND)
        layer8.InitEncryptedTunnel(globalThis.BACKEND)
    });
})();

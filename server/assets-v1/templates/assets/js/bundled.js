var xyz = (() => {
  var __defProp = Object.defineProperty;
  var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
  var __getOwnPropNames = Object.getOwnPropertyNames;
  var __hasOwnProp = Object.prototype.hasOwnProperty;
  var __export = (target, all) => {
    for (var name in all)
      __defProp(target, name, { get: all[name], enumerable: true });
  };
  var __copyProps = (to, from, except, desc) => {
    if (from && typeof from === "object" || typeof from === "function") {
      for (let key of __getOwnPropNames(from))
        if (!__hasOwnProp.call(to, key) && key !== except)
          __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
    }
    return to;
  };
  var __toCommonJS = (mod2) => __copyProps(__defProp({}, "__esModule", { value: true }), mod2);

  // index.ts
  var mnemonic_exports = {};
  __export(mnemonic_exports, {
    generateBip39Mnemonic: () => generateBip39Mnemonic,
    getBinarySeed: () => getBinarySeed,
    getPrivateAndPublicKeys: () => getPrivateAndPublicKeys,
    isMnemonicValid: () => isMnemonicValid
  });

  // node_modules/@noble/hashes/esm/_assert.js
  function number(n) {
    if (!Number.isSafeInteger(n) || n < 0)
      throw new Error(`positive integer expected, not ${n}`);
  }
  function isBytes(a) {
    return a instanceof Uint8Array || a != null && typeof a === "object" && a.constructor.name === "Uint8Array";
  }
  function bytes(b, ...lengths) {
    if (!isBytes(b))
      throw new Error("Uint8Array expected");
    if (lengths.length > 0 && !lengths.includes(b.length))
      throw new Error(`Uint8Array expected of length ${lengths}, not of length=${b.length}`);
  }
  function hash(h) {
    if (typeof h !== "function" || typeof h.create !== "function")
      throw new Error("Hash should be wrapped by utils.wrapConstructor");
    number(h.outputLen);
    number(h.blockLen);
  }
  function exists(instance, checkFinished = true) {
    if (instance.destroyed)
      throw new Error("Hash instance has been destroyed");
    if (checkFinished && instance.finished)
      throw new Error("Hash#digest() has already been called");
  }
  function output(out, instance) {
    bytes(out);
    const min = instance.outputLen;
    if (out.length < min) {
      throw new Error(`digestInto() expects output buffer of length at least ${min}`);
    }
  }

  // node_modules/@noble/hashes/esm/crypto.js
  var crypto = typeof globalThis === "object" && "crypto" in globalThis ? globalThis.crypto : void 0;

  // node_modules/@noble/hashes/esm/utils.js
  var createView = (arr) => new DataView(arr.buffer, arr.byteOffset, arr.byteLength);
  var rotr = (word, shift) => word << 32 - shift | word >>> shift;
  var rotl = (word, shift) => word << shift | word >>> 32 - shift >>> 0;
  var isLE = new Uint8Array(new Uint32Array([287454020]).buffer)[0] === 68;
  var hexes = /* @__PURE__ */ Array.from({ length: 256 }, (_, i) => i.toString(16).padStart(2, "0"));
  function bytesToHex(bytes2) {
    bytes(bytes2);
    let hex = "";
    for (let i = 0; i < bytes2.length; i++) {
      hex += hexes[bytes2[i]];
    }
    return hex;
  }
  var asciis = { _0: 48, _9: 57, _A: 65, _F: 70, _a: 97, _f: 102 };
  function asciiToBase16(char) {
    if (char >= asciis._0 && char <= asciis._9)
      return char - asciis._0;
    if (char >= asciis._A && char <= asciis._F)
      return char - (asciis._A - 10);
    if (char >= asciis._a && char <= asciis._f)
      return char - (asciis._a - 10);
    return;
  }
  function hexToBytes(hex) {
    if (typeof hex !== "string")
      throw new Error("hex string expected, got " + typeof hex);
    const hl = hex.length;
    const al = hl / 2;
    if (hl % 2)
      throw new Error("padded hex string expected, got unpadded hex of length " + hl);
    const array = new Uint8Array(al);
    for (let ai = 0, hi = 0; ai < al; ai++, hi += 2) {
      const n1 = asciiToBase16(hex.charCodeAt(hi));
      const n2 = asciiToBase16(hex.charCodeAt(hi + 1));
      if (n1 === void 0 || n2 === void 0) {
        const char = hex[hi] + hex[hi + 1];
        throw new Error('hex string expected, got non-hex character "' + char + '" at index ' + hi);
      }
      array[ai] = n1 * 16 + n2;
    }
    return array;
  }
  function utf8ToBytes(str) {
    if (typeof str !== "string")
      throw new Error(`utf8ToBytes expected string, got ${typeof str}`);
    return new Uint8Array(new TextEncoder().encode(str));
  }
  function toBytes(data) {
    if (typeof data === "string")
      data = utf8ToBytes(data);
    bytes(data);
    return data;
  }
  function concatBytes(...arrays) {
    let sum = 0;
    for (let i = 0; i < arrays.length; i++) {
      const a = arrays[i];
      bytes(a);
      sum += a.length;
    }
    const res = new Uint8Array(sum);
    for (let i = 0, pad = 0; i < arrays.length; i++) {
      const a = arrays[i];
      res.set(a, pad);
      pad += a.length;
    }
    return res;
  }
  var Hash = class {
    // Safe version that clones internal state
    clone() {
      return this._cloneInto();
    }
  };
  var toStr = {}.toString;
  function checkOpts(defaults, opts) {
    if (opts !== void 0 && toStr.call(opts) !== "[object Object]")
      throw new Error("Options should be object or undefined");
    const merged = Object.assign(defaults, opts);
    return merged;
  }
  function wrapConstructor(hashCons) {
    const hashC = (msg) => hashCons().update(toBytes(msg)).digest();
    const tmp = hashCons();
    hashC.outputLen = tmp.outputLen;
    hashC.blockLen = tmp.blockLen;
    hashC.create = () => hashCons();
    return hashC;
  }
  function randomBytes(bytesLength = 32) {
    if (crypto && typeof crypto.getRandomValues === "function") {
      return crypto.getRandomValues(new Uint8Array(bytesLength));
    }
    if (crypto && typeof crypto.randomBytes === "function") {
      return crypto.randomBytes(bytesLength);
    }
    throw new Error("crypto.getRandomValues must be defined");
  }

  // node_modules/@noble/hashes/esm/hmac.js
  var HMAC = class extends Hash {
    constructor(hash2, _key) {
      super();
      this.finished = false;
      this.destroyed = false;
      hash(hash2);
      const key = toBytes(_key);
      this.iHash = hash2.create();
      if (typeof this.iHash.update !== "function")
        throw new Error("Expected instance of class which extends utils.Hash");
      this.blockLen = this.iHash.blockLen;
      this.outputLen = this.iHash.outputLen;
      const blockLen = this.blockLen;
      const pad = new Uint8Array(blockLen);
      pad.set(key.length > blockLen ? hash2.create().update(key).digest() : key);
      for (let i = 0; i < pad.length; i++)
        pad[i] ^= 54;
      this.iHash.update(pad);
      this.oHash = hash2.create();
      for (let i = 0; i < pad.length; i++)
        pad[i] ^= 54 ^ 92;
      this.oHash.update(pad);
      pad.fill(0);
    }
    update(buf) {
      exists(this);
      this.iHash.update(buf);
      return this;
    }
    digestInto(out) {
      exists(this);
      bytes(out, this.outputLen);
      this.finished = true;
      this.iHash.digestInto(out);
      this.oHash.update(out);
      this.oHash.digestInto(out);
      this.destroy();
    }
    digest() {
      const out = new Uint8Array(this.oHash.outputLen);
      this.digestInto(out);
      return out;
    }
    _cloneInto(to) {
      to || (to = Object.create(Object.getPrototypeOf(this), {}));
      const { oHash, iHash, finished, destroyed, blockLen, outputLen } = this;
      to = to;
      to.finished = finished;
      to.destroyed = destroyed;
      to.blockLen = blockLen;
      to.outputLen = outputLen;
      to.oHash = oHash._cloneInto(to.oHash);
      to.iHash = iHash._cloneInto(to.iHash);
      return to;
    }
    destroy() {
      this.destroyed = true;
      this.oHash.destroy();
      this.iHash.destroy();
    }
  };
  var hmac = (hash2, key, message) => new HMAC(hash2, key).update(message).digest();
  hmac.create = (hash2, key) => new HMAC(hash2, key);

  // node_modules/@noble/hashes/esm/pbkdf2.js
  function pbkdf2Init(hash2, _password, _salt, _opts) {
    hash(hash2);
    const opts = checkOpts({ dkLen: 32, asyncTick: 10 }, _opts);
    const { c, dkLen, asyncTick } = opts;
    number(c);
    number(dkLen);
    number(asyncTick);
    if (c < 1)
      throw new Error("PBKDF2: iterations (c) should be >= 1");
    const password = toBytes(_password);
    const salt2 = toBytes(_salt);
    const DK = new Uint8Array(dkLen);
    const PRF = hmac.create(hash2, password);
    const PRFSalt = PRF._cloneInto().update(salt2);
    return { c, dkLen, asyncTick, DK, PRF, PRFSalt };
  }
  function pbkdf2Output(PRF, PRFSalt, DK, prfW, u) {
    PRF.destroy();
    PRFSalt.destroy();
    if (prfW)
      prfW.destroy();
    u.fill(0);
    return DK;
  }
  function pbkdf2(hash2, password, salt2, opts) {
    const { c, dkLen, DK, PRF, PRFSalt } = pbkdf2Init(hash2, password, salt2, opts);
    let prfW;
    const arr = new Uint8Array(4);
    const view = createView(arr);
    const u = new Uint8Array(PRF.outputLen);
    for (let ti = 1, pos = 0; pos < dkLen; ti++, pos += PRF.outputLen) {
      const Ti = DK.subarray(pos, pos + PRF.outputLen);
      view.setInt32(0, ti, false);
      (prfW = PRFSalt._cloneInto(prfW)).update(arr).digestInto(u);
      Ti.set(u.subarray(0, Ti.length));
      for (let ui = 1; ui < c; ui++) {
        PRF._cloneInto(prfW).update(u).digestInto(u);
        for (let i = 0; i < Ti.length; i++)
          Ti[i] ^= u[i];
      }
    }
    return pbkdf2Output(PRF, PRFSalt, DK, prfW, u);
  }

  // node_modules/@noble/hashes/esm/_md.js
  function setBigUint64(view, byteOffset, value, isLE2) {
    if (typeof view.setBigUint64 === "function")
      return view.setBigUint64(byteOffset, value, isLE2);
    const _32n2 = BigInt(32);
    const _u32_max = BigInt(4294967295);
    const wh = Number(value >> _32n2 & _u32_max);
    const wl = Number(value & _u32_max);
    const h = isLE2 ? 4 : 0;
    const l = isLE2 ? 0 : 4;
    view.setUint32(byteOffset + h, wh, isLE2);
    view.setUint32(byteOffset + l, wl, isLE2);
  }
  var Chi = (a, b, c) => a & b ^ ~a & c;
  var Maj = (a, b, c) => a & b ^ a & c ^ b & c;
  var HashMD = class extends Hash {
    constructor(blockLen, outputLen, padOffset, isLE2) {
      super();
      this.blockLen = blockLen;
      this.outputLen = outputLen;
      this.padOffset = padOffset;
      this.isLE = isLE2;
      this.finished = false;
      this.length = 0;
      this.pos = 0;
      this.destroyed = false;
      this.buffer = new Uint8Array(blockLen);
      this.view = createView(this.buffer);
    }
    update(data) {
      exists(this);
      const { view, buffer, blockLen } = this;
      data = toBytes(data);
      const len = data.length;
      for (let pos = 0; pos < len; ) {
        const take = Math.min(blockLen - this.pos, len - pos);
        if (take === blockLen) {
          const dataView = createView(data);
          for (; blockLen <= len - pos; pos += blockLen)
            this.process(dataView, pos);
          continue;
        }
        buffer.set(data.subarray(pos, pos + take), this.pos);
        this.pos += take;
        pos += take;
        if (this.pos === blockLen) {
          this.process(view, 0);
          this.pos = 0;
        }
      }
      this.length += data.length;
      this.roundClean();
      return this;
    }
    digestInto(out) {
      exists(this);
      output(out, this);
      this.finished = true;
      const { buffer, view, blockLen, isLE: isLE2 } = this;
      let { pos } = this;
      buffer[pos++] = 128;
      this.buffer.subarray(pos).fill(0);
      if (this.padOffset > blockLen - pos) {
        this.process(view, 0);
        pos = 0;
      }
      for (let i = pos; i < blockLen; i++)
        buffer[i] = 0;
      setBigUint64(view, blockLen - 8, BigInt(this.length * 8), isLE2);
      this.process(view, 0);
      const oview = createView(out);
      const len = this.outputLen;
      if (len % 4)
        throw new Error("_sha2: outputLen should be aligned to 32bit");
      const outLen = len / 4;
      const state = this.get();
      if (outLen > state.length)
        throw new Error("_sha2: outputLen bigger than state");
      for (let i = 0; i < outLen; i++)
        oview.setUint32(4 * i, state[i], isLE2);
    }
    digest() {
      const { buffer, outputLen } = this;
      this.digestInto(buffer);
      const res = buffer.slice(0, outputLen);
      this.destroy();
      return res;
    }
    _cloneInto(to) {
      to || (to = new this.constructor());
      to.set(...this.get());
      const { blockLen, buffer, length, finished, destroyed, pos } = this;
      to.length = length;
      to.pos = pos;
      to.finished = finished;
      to.destroyed = destroyed;
      if (length % blockLen)
        to.buffer.set(buffer);
      return to;
    }
  };

  // node_modules/@noble/hashes/esm/sha256.js
  var SHA256_K = /* @__PURE__ */ new Uint32Array([
    1116352408,
    1899447441,
    3049323471,
    3921009573,
    961987163,
    1508970993,
    2453635748,
    2870763221,
    3624381080,
    310598401,
    607225278,
    1426881987,
    1925078388,
    2162078206,
    2614888103,
    3248222580,
    3835390401,
    4022224774,
    264347078,
    604807628,
    770255983,
    1249150122,
    1555081692,
    1996064986,
    2554220882,
    2821834349,
    2952996808,
    3210313671,
    3336571891,
    3584528711,
    113926993,
    338241895,
    666307205,
    773529912,
    1294757372,
    1396182291,
    1695183700,
    1986661051,
    2177026350,
    2456956037,
    2730485921,
    2820302411,
    3259730800,
    3345764771,
    3516065817,
    3600352804,
    4094571909,
    275423344,
    430227734,
    506948616,
    659060556,
    883997877,
    958139571,
    1322822218,
    1537002063,
    1747873779,
    1955562222,
    2024104815,
    2227730452,
    2361852424,
    2428436474,
    2756734187,
    3204031479,
    3329325298
  ]);
  var SHA256_IV = /* @__PURE__ */ new Uint32Array([
    1779033703,
    3144134277,
    1013904242,
    2773480762,
    1359893119,
    2600822924,
    528734635,
    1541459225
  ]);
  var SHA256_W = /* @__PURE__ */ new Uint32Array(64);
  var SHA256 = class extends HashMD {
    constructor() {
      super(64, 32, 8, false);
      this.A = SHA256_IV[0] | 0;
      this.B = SHA256_IV[1] | 0;
      this.C = SHA256_IV[2] | 0;
      this.D = SHA256_IV[3] | 0;
      this.E = SHA256_IV[4] | 0;
      this.F = SHA256_IV[5] | 0;
      this.G = SHA256_IV[6] | 0;
      this.H = SHA256_IV[7] | 0;
    }
    get() {
      const { A, B, C, D, E, F, G, H } = this;
      return [A, B, C, D, E, F, G, H];
    }
    // prettier-ignore
    set(A, B, C, D, E, F, G, H) {
      this.A = A | 0;
      this.B = B | 0;
      this.C = C | 0;
      this.D = D | 0;
      this.E = E | 0;
      this.F = F | 0;
      this.G = G | 0;
      this.H = H | 0;
    }
    process(view, offset) {
      for (let i = 0; i < 16; i++, offset += 4)
        SHA256_W[i] = view.getUint32(offset, false);
      for (let i = 16; i < 64; i++) {
        const W15 = SHA256_W[i - 15];
        const W2 = SHA256_W[i - 2];
        const s0 = rotr(W15, 7) ^ rotr(W15, 18) ^ W15 >>> 3;
        const s1 = rotr(W2, 17) ^ rotr(W2, 19) ^ W2 >>> 10;
        SHA256_W[i] = s1 + SHA256_W[i - 7] + s0 + SHA256_W[i - 16] | 0;
      }
      let { A, B, C, D, E, F, G, H } = this;
      for (let i = 0; i < 64; i++) {
        const sigma1 = rotr(E, 6) ^ rotr(E, 11) ^ rotr(E, 25);
        const T1 = H + sigma1 + Chi(E, F, G) + SHA256_K[i] + SHA256_W[i] | 0;
        const sigma0 = rotr(A, 2) ^ rotr(A, 13) ^ rotr(A, 22);
        const T2 = sigma0 + Maj(A, B, C) | 0;
        H = G;
        G = F;
        F = E;
        E = D + T1 | 0;
        D = C;
        C = B;
        B = A;
        A = T1 + T2 | 0;
      }
      A = A + this.A | 0;
      B = B + this.B | 0;
      C = C + this.C | 0;
      D = D + this.D | 0;
      E = E + this.E | 0;
      F = F + this.F | 0;
      G = G + this.G | 0;
      H = H + this.H | 0;
      this.set(A, B, C, D, E, F, G, H);
    }
    roundClean() {
      SHA256_W.fill(0);
    }
    destroy() {
      this.set(0, 0, 0, 0, 0, 0, 0, 0);
      this.buffer.fill(0);
    }
  };
  var sha256 = /* @__PURE__ */ wrapConstructor(() => new SHA256());

  // node_modules/@noble/hashes/esm/_u64.js
  var U32_MASK64 = /* @__PURE__ */ BigInt(2 ** 32 - 1);
  var _32n = /* @__PURE__ */ BigInt(32);
  function fromBig(n, le = false) {
    if (le)
      return { h: Number(n & U32_MASK64), l: Number(n >> _32n & U32_MASK64) };
    return { h: Number(n >> _32n & U32_MASK64) | 0, l: Number(n & U32_MASK64) | 0 };
  }
  function split(lst, le = false) {
    let Ah = new Uint32Array(lst.length);
    let Al = new Uint32Array(lst.length);
    for (let i = 0; i < lst.length; i++) {
      const { h, l } = fromBig(lst[i], le);
      [Ah[i], Al[i]] = [h, l];
    }
    return [Ah, Al];
  }
  var toBig = (h, l) => BigInt(h >>> 0) << _32n | BigInt(l >>> 0);
  var shrSH = (h, _l, s) => h >>> s;
  var shrSL = (h, l, s) => h << 32 - s | l >>> s;
  var rotrSH = (h, l, s) => h >>> s | l << 32 - s;
  var rotrSL = (h, l, s) => h << 32 - s | l >>> s;
  var rotrBH = (h, l, s) => h << 64 - s | l >>> s - 32;
  var rotrBL = (h, l, s) => h >>> s - 32 | l << 64 - s;
  var rotr32H = (_h, l) => l;
  var rotr32L = (h, _l) => h;
  var rotlSH = (h, l, s) => h << s | l >>> 32 - s;
  var rotlSL = (h, l, s) => l << s | h >>> 32 - s;
  var rotlBH = (h, l, s) => l << s - 32 | h >>> 64 - s;
  var rotlBL = (h, l, s) => h << s - 32 | l >>> 64 - s;
  function add(Ah, Al, Bh, Bl) {
    const l = (Al >>> 0) + (Bl >>> 0);
    return { h: Ah + Bh + (l / 2 ** 32 | 0) | 0, l: l | 0 };
  }
  var add3L = (Al, Bl, Cl) => (Al >>> 0) + (Bl >>> 0) + (Cl >>> 0);
  var add3H = (low, Ah, Bh, Ch) => Ah + Bh + Ch + (low / 2 ** 32 | 0) | 0;
  var add4L = (Al, Bl, Cl, Dl) => (Al >>> 0) + (Bl >>> 0) + (Cl >>> 0) + (Dl >>> 0);
  var add4H = (low, Ah, Bh, Ch, Dh) => Ah + Bh + Ch + Dh + (low / 2 ** 32 | 0) | 0;
  var add5L = (Al, Bl, Cl, Dl, El) => (Al >>> 0) + (Bl >>> 0) + (Cl >>> 0) + (Dl >>> 0) + (El >>> 0);
  var add5H = (low, Ah, Bh, Ch, Dh, Eh) => Ah + Bh + Ch + Dh + Eh + (low / 2 ** 32 | 0) | 0;
  var u64 = {
    fromBig,
    split,
    toBig,
    shrSH,
    shrSL,
    rotrSH,
    rotrSL,
    rotrBH,
    rotrBL,
    rotr32H,
    rotr32L,
    rotlSH,
    rotlSL,
    rotlBH,
    rotlBL,
    add,
    add3L,
    add3H,
    add4L,
    add4H,
    add5H,
    add5L
  };
  var u64_default = u64;

  // node_modules/@noble/hashes/esm/sha512.js
  var [SHA512_Kh, SHA512_Kl] = /* @__PURE__ */ (() => u64_default.split([
    "0x428a2f98d728ae22",
    "0x7137449123ef65cd",
    "0xb5c0fbcfec4d3b2f",
    "0xe9b5dba58189dbbc",
    "0x3956c25bf348b538",
    "0x59f111f1b605d019",
    "0x923f82a4af194f9b",
    "0xab1c5ed5da6d8118",
    "0xd807aa98a3030242",
    "0x12835b0145706fbe",
    "0x243185be4ee4b28c",
    "0x550c7dc3d5ffb4e2",
    "0x72be5d74f27b896f",
    "0x80deb1fe3b1696b1",
    "0x9bdc06a725c71235",
    "0xc19bf174cf692694",
    "0xe49b69c19ef14ad2",
    "0xefbe4786384f25e3",
    "0x0fc19dc68b8cd5b5",
    "0x240ca1cc77ac9c65",
    "0x2de92c6f592b0275",
    "0x4a7484aa6ea6e483",
    "0x5cb0a9dcbd41fbd4",
    "0x76f988da831153b5",
    "0x983e5152ee66dfab",
    "0xa831c66d2db43210",
    "0xb00327c898fb213f",
    "0xbf597fc7beef0ee4",
    "0xc6e00bf33da88fc2",
    "0xd5a79147930aa725",
    "0x06ca6351e003826f",
    "0x142929670a0e6e70",
    "0x27b70a8546d22ffc",
    "0x2e1b21385c26c926",
    "0x4d2c6dfc5ac42aed",
    "0x53380d139d95b3df",
    "0x650a73548baf63de",
    "0x766a0abb3c77b2a8",
    "0x81c2c92e47edaee6",
    "0x92722c851482353b",
    "0xa2bfe8a14cf10364",
    "0xa81a664bbc423001",
    "0xc24b8b70d0f89791",
    "0xc76c51a30654be30",
    "0xd192e819d6ef5218",
    "0xd69906245565a910",
    "0xf40e35855771202a",
    "0x106aa07032bbd1b8",
    "0x19a4c116b8d2d0c8",
    "0x1e376c085141ab53",
    "0x2748774cdf8eeb99",
    "0x34b0bcb5e19b48a8",
    "0x391c0cb3c5c95a63",
    "0x4ed8aa4ae3418acb",
    "0x5b9cca4f7763e373",
    "0x682e6ff3d6b2b8a3",
    "0x748f82ee5defb2fc",
    "0x78a5636f43172f60",
    "0x84c87814a1f0ab72",
    "0x8cc702081a6439ec",
    "0x90befffa23631e28",
    "0xa4506cebde82bde9",
    "0xbef9a3f7b2c67915",
    "0xc67178f2e372532b",
    "0xca273eceea26619c",
    "0xd186b8c721c0c207",
    "0xeada7dd6cde0eb1e",
    "0xf57d4f7fee6ed178",
    "0x06f067aa72176fba",
    "0x0a637dc5a2c898a6",
    "0x113f9804bef90dae",
    "0x1b710b35131c471b",
    "0x28db77f523047d84",
    "0x32caab7b40c72493",
    "0x3c9ebe0a15c9bebc",
    "0x431d67c49c100d4c",
    "0x4cc5d4becb3e42b6",
    "0x597f299cfc657e2a",
    "0x5fcb6fab3ad6faec",
    "0x6c44198c4a475817"
  ].map((n) => BigInt(n))))();
  var SHA512_W_H = /* @__PURE__ */ new Uint32Array(80);
  var SHA512_W_L = /* @__PURE__ */ new Uint32Array(80);
  var SHA512 = class extends HashMD {
    constructor() {
      super(128, 64, 16, false);
      this.Ah = 1779033703 | 0;
      this.Al = 4089235720 | 0;
      this.Bh = 3144134277 | 0;
      this.Bl = 2227873595 | 0;
      this.Ch = 1013904242 | 0;
      this.Cl = 4271175723 | 0;
      this.Dh = 2773480762 | 0;
      this.Dl = 1595750129 | 0;
      this.Eh = 1359893119 | 0;
      this.El = 2917565137 | 0;
      this.Fh = 2600822924 | 0;
      this.Fl = 725511199 | 0;
      this.Gh = 528734635 | 0;
      this.Gl = 4215389547 | 0;
      this.Hh = 1541459225 | 0;
      this.Hl = 327033209 | 0;
    }
    // prettier-ignore
    get() {
      const { Ah, Al, Bh, Bl, Ch, Cl, Dh, Dl, Eh, El, Fh, Fl, Gh, Gl, Hh, Hl } = this;
      return [Ah, Al, Bh, Bl, Ch, Cl, Dh, Dl, Eh, El, Fh, Fl, Gh, Gl, Hh, Hl];
    }
    // prettier-ignore
    set(Ah, Al, Bh, Bl, Ch, Cl, Dh, Dl, Eh, El, Fh, Fl, Gh, Gl, Hh, Hl) {
      this.Ah = Ah | 0;
      this.Al = Al | 0;
      this.Bh = Bh | 0;
      this.Bl = Bl | 0;
      this.Ch = Ch | 0;
      this.Cl = Cl | 0;
      this.Dh = Dh | 0;
      this.Dl = Dl | 0;
      this.Eh = Eh | 0;
      this.El = El | 0;
      this.Fh = Fh | 0;
      this.Fl = Fl | 0;
      this.Gh = Gh | 0;
      this.Gl = Gl | 0;
      this.Hh = Hh | 0;
      this.Hl = Hl | 0;
    }
    process(view, offset) {
      for (let i = 0; i < 16; i++, offset += 4) {
        SHA512_W_H[i] = view.getUint32(offset);
        SHA512_W_L[i] = view.getUint32(offset += 4);
      }
      for (let i = 16; i < 80; i++) {
        const W15h = SHA512_W_H[i - 15] | 0;
        const W15l = SHA512_W_L[i - 15] | 0;
        const s0h = u64_default.rotrSH(W15h, W15l, 1) ^ u64_default.rotrSH(W15h, W15l, 8) ^ u64_default.shrSH(W15h, W15l, 7);
        const s0l = u64_default.rotrSL(W15h, W15l, 1) ^ u64_default.rotrSL(W15h, W15l, 8) ^ u64_default.shrSL(W15h, W15l, 7);
        const W2h = SHA512_W_H[i - 2] | 0;
        const W2l = SHA512_W_L[i - 2] | 0;
        const s1h = u64_default.rotrSH(W2h, W2l, 19) ^ u64_default.rotrBH(W2h, W2l, 61) ^ u64_default.shrSH(W2h, W2l, 6);
        const s1l = u64_default.rotrSL(W2h, W2l, 19) ^ u64_default.rotrBL(W2h, W2l, 61) ^ u64_default.shrSL(W2h, W2l, 6);
        const SUMl = u64_default.add4L(s0l, s1l, SHA512_W_L[i - 7], SHA512_W_L[i - 16]);
        const SUMh = u64_default.add4H(SUMl, s0h, s1h, SHA512_W_H[i - 7], SHA512_W_H[i - 16]);
        SHA512_W_H[i] = SUMh | 0;
        SHA512_W_L[i] = SUMl | 0;
      }
      let { Ah, Al, Bh, Bl, Ch, Cl, Dh, Dl, Eh, El, Fh, Fl, Gh, Gl, Hh, Hl } = this;
      for (let i = 0; i < 80; i++) {
        const sigma1h = u64_default.rotrSH(Eh, El, 14) ^ u64_default.rotrSH(Eh, El, 18) ^ u64_default.rotrBH(Eh, El, 41);
        const sigma1l = u64_default.rotrSL(Eh, El, 14) ^ u64_default.rotrSL(Eh, El, 18) ^ u64_default.rotrBL(Eh, El, 41);
        const CHIh = Eh & Fh ^ ~Eh & Gh;
        const CHIl = El & Fl ^ ~El & Gl;
        const T1ll = u64_default.add5L(Hl, sigma1l, CHIl, SHA512_Kl[i], SHA512_W_L[i]);
        const T1h = u64_default.add5H(T1ll, Hh, sigma1h, CHIh, SHA512_Kh[i], SHA512_W_H[i]);
        const T1l = T1ll | 0;
        const sigma0h = u64_default.rotrSH(Ah, Al, 28) ^ u64_default.rotrBH(Ah, Al, 34) ^ u64_default.rotrBH(Ah, Al, 39);
        const sigma0l = u64_default.rotrSL(Ah, Al, 28) ^ u64_default.rotrBL(Ah, Al, 34) ^ u64_default.rotrBL(Ah, Al, 39);
        const MAJh = Ah & Bh ^ Ah & Ch ^ Bh & Ch;
        const MAJl = Al & Bl ^ Al & Cl ^ Bl & Cl;
        Hh = Gh | 0;
        Hl = Gl | 0;
        Gh = Fh | 0;
        Gl = Fl | 0;
        Fh = Eh | 0;
        Fl = El | 0;
        ({ h: Eh, l: El } = u64_default.add(Dh | 0, Dl | 0, T1h | 0, T1l | 0));
        Dh = Ch | 0;
        Dl = Cl | 0;
        Ch = Bh | 0;
        Cl = Bl | 0;
        Bh = Ah | 0;
        Bl = Al | 0;
        const All = u64_default.add3L(T1l, sigma0l, MAJl);
        Ah = u64_default.add3H(All, T1h, sigma0h, MAJh);
        Al = All | 0;
      }
      ({ h: Ah, l: Al } = u64_default.add(this.Ah | 0, this.Al | 0, Ah | 0, Al | 0));
      ({ h: Bh, l: Bl } = u64_default.add(this.Bh | 0, this.Bl | 0, Bh | 0, Bl | 0));
      ({ h: Ch, l: Cl } = u64_default.add(this.Ch | 0, this.Cl | 0, Ch | 0, Cl | 0));
      ({ h: Dh, l: Dl } = u64_default.add(this.Dh | 0, this.Dl | 0, Dh | 0, Dl | 0));
      ({ h: Eh, l: El } = u64_default.add(this.Eh | 0, this.El | 0, Eh | 0, El | 0));
      ({ h: Fh, l: Fl } = u64_default.add(this.Fh | 0, this.Fl | 0, Fh | 0, Fl | 0));
      ({ h: Gh, l: Gl } = u64_default.add(this.Gh | 0, this.Gl | 0, Gh | 0, Gl | 0));
      ({ h: Hh, l: Hl } = u64_default.add(this.Hh | 0, this.Hl | 0, Hh | 0, Hl | 0));
      this.set(Ah, Al, Bh, Bl, Ch, Cl, Dh, Dl, Eh, El, Fh, Fl, Gh, Gl, Hh, Hl);
    }
    roundClean() {
      SHA512_W_H.fill(0);
      SHA512_W_L.fill(0);
    }
    destroy() {
      this.buffer.fill(0);
      this.set(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0);
    }
  };
  var sha512 = /* @__PURE__ */ wrapConstructor(() => new SHA512());

  // node_modules/@scure/base/lib/esm/index.js
  // @__NO_SIDE_EFFECTS__
  function assertNumber(n) {
    if (!Number.isSafeInteger(n))
      throw new Error(`Wrong integer: ${n}`);
  }
  function isBytes2(a) {
    return a instanceof Uint8Array || a != null && typeof a === "object" && a.constructor.name === "Uint8Array";
  }
  // @__NO_SIDE_EFFECTS__
  function chain(...args) {
    const id = (a) => a;
    const wrap = (a, b) => (c) => a(b(c));
    const encode = args.map((x) => x.encode).reduceRight(wrap, id);
    const decode = args.map((x) => x.decode).reduce(wrap, id);
    return { encode, decode };
  }
  // @__NO_SIDE_EFFECTS__
  function alphabet(alphabet2) {
    return {
      encode: (digits) => {
        if (!Array.isArray(digits) || digits.length && typeof digits[0] !== "number")
          throw new Error("alphabet.encode input should be an array of numbers");
        return digits.map((i) => {
          /* @__PURE__ */ assertNumber(i);
          if (i < 0 || i >= alphabet2.length)
            throw new Error(`Digit index outside alphabet: ${i} (alphabet: ${alphabet2.length})`);
          return alphabet2[i];
        });
      },
      decode: (input) => {
        if (!Array.isArray(input) || input.length && typeof input[0] !== "string")
          throw new Error("alphabet.decode input should be array of strings");
        return input.map((letter) => {
          if (typeof letter !== "string")
            throw new Error(`alphabet.decode: not string element=${letter}`);
          const index = alphabet2.indexOf(letter);
          if (index === -1)
            throw new Error(`Unknown letter: "${letter}". Allowed: ${alphabet2}`);
          return index;
        });
      }
    };
  }
  // @__NO_SIDE_EFFECTS__
  function join(separator = "") {
    if (typeof separator !== "string")
      throw new Error("join separator should be string");
    return {
      encode: (from) => {
        if (!Array.isArray(from) || from.length && typeof from[0] !== "string")
          throw new Error("join.encode input should be array of strings");
        for (let i of from)
          if (typeof i !== "string")
            throw new Error(`join.encode: non-string input=${i}`);
        return from.join(separator);
      },
      decode: (to) => {
        if (typeof to !== "string")
          throw new Error("join.decode input should be string");
        return to.split(separator);
      }
    };
  }
  // @__NO_SIDE_EFFECTS__
  function padding(bits, chr = "=") {
    /* @__PURE__ */ assertNumber(bits);
    if (typeof chr !== "string")
      throw new Error("padding chr should be string");
    return {
      encode(data) {
        if (!Array.isArray(data) || data.length && typeof data[0] !== "string")
          throw new Error("padding.encode input should be array of strings");
        for (let i of data)
          if (typeof i !== "string")
            throw new Error(`padding.encode: non-string input=${i}`);
        while (data.length * bits % 8)
          data.push(chr);
        return data;
      },
      decode(input) {
        if (!Array.isArray(input) || input.length && typeof input[0] !== "string")
          throw new Error("padding.encode input should be array of strings");
        for (let i of input)
          if (typeof i !== "string")
            throw new Error(`padding.decode: non-string input=${i}`);
        let end = input.length;
        if (end * bits % 8)
          throw new Error("Invalid padding: string should have whole number of bytes");
        for (; end > 0 && input[end - 1] === chr; end--) {
          if (!((end - 1) * bits % 8))
            throw new Error("Invalid padding: string has too much padding");
        }
        return input.slice(0, end);
      }
    };
  }
  // @__NO_SIDE_EFFECTS__
  function convertRadix(data, from, to) {
    if (from < 2)
      throw new Error(`convertRadix: wrong from=${from}, base cannot be less than 2`);
    if (to < 2)
      throw new Error(`convertRadix: wrong to=${to}, base cannot be less than 2`);
    if (!Array.isArray(data))
      throw new Error("convertRadix: data should be array");
    if (!data.length)
      return [];
    let pos = 0;
    const res = [];
    const digits = Array.from(data);
    digits.forEach((d) => {
      /* @__PURE__ */ assertNumber(d);
      if (d < 0 || d >= from)
        throw new Error(`Wrong integer: ${d}`);
    });
    while (true) {
      let carry = 0;
      let done = true;
      for (let i = pos; i < digits.length; i++) {
        const digit = digits[i];
        const digitBase = from * carry + digit;
        if (!Number.isSafeInteger(digitBase) || from * carry / from !== carry || digitBase - digit !== from * carry) {
          throw new Error("convertRadix: carry overflow");
        }
        carry = digitBase % to;
        const rounded = Math.floor(digitBase / to);
        digits[i] = rounded;
        if (!Number.isSafeInteger(rounded) || rounded * to + carry !== digitBase)
          throw new Error("convertRadix: carry overflow");
        if (!done)
          continue;
        else if (!rounded)
          pos = i;
        else
          done = false;
      }
      res.push(carry);
      if (done)
        break;
    }
    for (let i = 0; i < data.length - 1 && data[i] === 0; i++)
      res.push(0);
    return res.reverse();
  }
  var gcd = /* @__NO_SIDE_EFFECTS__ */ (a, b) => !b ? a : /* @__PURE__ */ gcd(b, a % b);
  var radix2carry = /* @__NO_SIDE_EFFECTS__ */ (from, to) => from + (to - /* @__PURE__ */ gcd(from, to));
  // @__NO_SIDE_EFFECTS__
  function convertRadix2(data, from, to, padding2) {
    if (!Array.isArray(data))
      throw new Error("convertRadix2: data should be array");
    if (from <= 0 || from > 32)
      throw new Error(`convertRadix2: wrong from=${from}`);
    if (to <= 0 || to > 32)
      throw new Error(`convertRadix2: wrong to=${to}`);
    if (/* @__PURE__ */ radix2carry(from, to) > 32) {
      throw new Error(`convertRadix2: carry overflow from=${from} to=${to} carryBits=${/* @__PURE__ */ radix2carry(from, to)}`);
    }
    let carry = 0;
    let pos = 0;
    const mask = 2 ** to - 1;
    const res = [];
    for (const n of data) {
      /* @__PURE__ */ assertNumber(n);
      if (n >= 2 ** from)
        throw new Error(`convertRadix2: invalid data word=${n} from=${from}`);
      carry = carry << from | n;
      if (pos + from > 32)
        throw new Error(`convertRadix2: carry overflow pos=${pos} from=${from}`);
      pos += from;
      for (; pos >= to; pos -= to)
        res.push((carry >> pos - to & mask) >>> 0);
      carry &= 2 ** pos - 1;
    }
    carry = carry << to - pos & mask;
    if (!padding2 && pos >= from)
      throw new Error("Excess padding");
    if (!padding2 && carry)
      throw new Error(`Non-zero padding: ${carry}`);
    if (padding2 && pos > 0)
      res.push(carry >>> 0);
    return res;
  }
  // @__NO_SIDE_EFFECTS__
  function radix(num) {
    /* @__PURE__ */ assertNumber(num);
    return {
      encode: (bytes2) => {
        if (!isBytes2(bytes2))
          throw new Error("radix.encode input should be Uint8Array");
        return /* @__PURE__ */ convertRadix(Array.from(bytes2), 2 ** 8, num);
      },
      decode: (digits) => {
        if (!Array.isArray(digits) || digits.length && typeof digits[0] !== "number")
          throw new Error("radix.decode input should be array of numbers");
        return Uint8Array.from(/* @__PURE__ */ convertRadix(digits, num, 2 ** 8));
      }
    };
  }
  // @__NO_SIDE_EFFECTS__
  function radix2(bits, revPadding = false) {
    /* @__PURE__ */ assertNumber(bits);
    if (bits <= 0 || bits > 32)
      throw new Error("radix2: bits should be in (0..32]");
    if (/* @__PURE__ */ radix2carry(8, bits) > 32 || /* @__PURE__ */ radix2carry(bits, 8) > 32)
      throw new Error("radix2: carry overflow");
    return {
      encode: (bytes2) => {
        if (!isBytes2(bytes2))
          throw new Error("radix2.encode input should be Uint8Array");
        return /* @__PURE__ */ convertRadix2(Array.from(bytes2), 8, bits, !revPadding);
      },
      decode: (digits) => {
        if (!Array.isArray(digits) || digits.length && typeof digits[0] !== "number")
          throw new Error("radix2.decode input should be array of numbers");
        return Uint8Array.from(/* @__PURE__ */ convertRadix2(digits, bits, 8, revPadding));
      }
    };
  }
  // @__NO_SIDE_EFFECTS__
  function checksum(len, fn) {
    /* @__PURE__ */ assertNumber(len);
    if (typeof fn !== "function")
      throw new Error("checksum fn should be function");
    return {
      encode(data) {
        if (!isBytes2(data))
          throw new Error("checksum.encode: input should be Uint8Array");
        const checksum2 = fn(data).slice(0, len);
        const res = new Uint8Array(data.length + len);
        res.set(data);
        res.set(checksum2, data.length);
        return res;
      },
      decode(data) {
        if (!isBytes2(data))
          throw new Error("checksum.decode: input should be Uint8Array");
        const payload = data.slice(0, -len);
        const newChecksum = fn(payload).slice(0, len);
        const oldChecksum = data.slice(-len);
        for (let i = 0; i < len; i++)
          if (newChecksum[i] !== oldChecksum[i])
            throw new Error("Invalid checksum");
        return payload;
      }
    };
  }
  var utils = {
    alphabet,
    chain,
    checksum,
    convertRadix,
    convertRadix2,
    radix,
    radix2,
    join,
    padding
  };
  var genBase58 = (abc) => /* @__PURE__ */ chain(/* @__PURE__ */ radix(58), /* @__PURE__ */ alphabet(abc), /* @__PURE__ */ join(""));
  var base58 = /* @__PURE__ */ genBase58("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz");
  var createBase58check = (sha2562) => /* @__PURE__ */ chain(/* @__PURE__ */ checksum(4, (data) => sha2562(sha2562(data))), base58);

  // node_modules/@scure/bip39/esm/index.js
  var isJapanese = (wordlist2) => wordlist2[0] === "\u3042\u3044\u3053\u304F\u3057\u3093";
  function nfkd(str) {
    if (typeof str !== "string")
      throw new TypeError(`Invalid mnemonic type: ${typeof str}`);
    return str.normalize("NFKD");
  }
  function normalize(str) {
    const norm = nfkd(str);
    const words = norm.split(" ");
    if (![12, 15, 18, 21, 24].includes(words.length))
      throw new Error("Invalid mnemonic");
    return { nfkd: norm, words };
  }
  function assertEntropy(entropy) {
    bytes(entropy, 16, 20, 24, 28, 32);
  }
  function generateMnemonic(wordlist2, strength = 128) {
    number(strength);
    if (strength % 32 !== 0 || strength > 256)
      throw new TypeError("Invalid entropy");
    return entropyToMnemonic(randomBytes(strength / 8), wordlist2);
  }
  var calcChecksum = (entropy) => {
    const bitsLeft = 8 - entropy.length / 4;
    return new Uint8Array([sha256(entropy)[0] >> bitsLeft << bitsLeft]);
  };
  function getCoder(wordlist2) {
    if (!Array.isArray(wordlist2) || wordlist2.length !== 2048 || typeof wordlist2[0] !== "string")
      throw new Error("Wordlist: expected array of 2048 strings");
    wordlist2.forEach((i) => {
      if (typeof i !== "string")
        throw new Error(`Wordlist: non-string element: ${i}`);
    });
    return utils.chain(utils.checksum(1, calcChecksum), utils.radix2(11, true), utils.alphabet(wordlist2));
  }
  function mnemonicToEntropy(mnemonic, wordlist2) {
    const { words } = normalize(mnemonic);
    const entropy = getCoder(wordlist2).decode(words);
    assertEntropy(entropy);
    return entropy;
  }
  function entropyToMnemonic(entropy, wordlist2) {
    assertEntropy(entropy);
    const words = getCoder(wordlist2).encode(entropy);
    return words.join(isJapanese(wordlist2) ? "\u3000" : " ");
  }
  function validateMnemonic(mnemonic, wordlist2) {
    try {
      mnemonicToEntropy(mnemonic, wordlist2);
    } catch (e) {
      return false;
    }
    return true;
  }
  var salt = (passphrase) => nfkd(`mnemonic${passphrase}`);
  function mnemonicToSeedSync(mnemonic, passphrase = "") {
    return pbkdf2(sha512, normalize(mnemonic).nfkd, salt(passphrase), { c: 2048, dkLen: 64 });
  }

  // node_modules/@scure/bip39/esm/wordlists/english.js
  var wordlist = `abandon
ability
able
about
above
absent
absorb
abstract
absurd
abuse
access
accident
account
accuse
achieve
acid
acoustic
acquire
across
act
action
actor
actress
actual
adapt
add
addict
address
adjust
admit
adult
advance
advice
aerobic
affair
afford
afraid
again
age
agent
agree
ahead
aim
air
airport
aisle
alarm
album
alcohol
alert
alien
all
alley
allow
almost
alone
alpha
already
also
alter
always
amateur
amazing
among
amount
amused
analyst
anchor
ancient
anger
angle
angry
animal
ankle
announce
annual
another
answer
antenna
antique
anxiety
any
apart
apology
appear
apple
approve
april
arch
arctic
area
arena
argue
arm
armed
armor
army
around
arrange
arrest
arrive
arrow
art
artefact
artist
artwork
ask
aspect
assault
asset
assist
assume
asthma
athlete
atom
attack
attend
attitude
attract
auction
audit
august
aunt
author
auto
autumn
average
avocado
avoid
awake
aware
away
awesome
awful
awkward
axis
baby
bachelor
bacon
badge
bag
balance
balcony
ball
bamboo
banana
banner
bar
barely
bargain
barrel
base
basic
basket
battle
beach
bean
beauty
because
become
beef
before
begin
behave
behind
believe
below
belt
bench
benefit
best
betray
better
between
beyond
bicycle
bid
bike
bind
biology
bird
birth
bitter
black
blade
blame
blanket
blast
bleak
bless
blind
blood
blossom
blouse
blue
blur
blush
board
boat
body
boil
bomb
bone
bonus
book
boost
border
boring
borrow
boss
bottom
bounce
box
boy
bracket
brain
brand
brass
brave
bread
breeze
brick
bridge
brief
bright
bring
brisk
broccoli
broken
bronze
broom
brother
brown
brush
bubble
buddy
budget
buffalo
build
bulb
bulk
bullet
bundle
bunker
burden
burger
burst
bus
business
busy
butter
buyer
buzz
cabbage
cabin
cable
cactus
cage
cake
call
calm
camera
camp
can
canal
cancel
candy
cannon
canoe
canvas
canyon
capable
capital
captain
car
carbon
card
cargo
carpet
carry
cart
case
cash
casino
castle
casual
cat
catalog
catch
category
cattle
caught
cause
caution
cave
ceiling
celery
cement
census
century
cereal
certain
chair
chalk
champion
change
chaos
chapter
charge
chase
chat
cheap
check
cheese
chef
cherry
chest
chicken
chief
child
chimney
choice
choose
chronic
chuckle
chunk
churn
cigar
cinnamon
circle
citizen
city
civil
claim
clap
clarify
claw
clay
clean
clerk
clever
click
client
cliff
climb
clinic
clip
clock
clog
close
cloth
cloud
clown
club
clump
cluster
clutch
coach
coast
coconut
code
coffee
coil
coin
collect
color
column
combine
come
comfort
comic
common
company
concert
conduct
confirm
congress
connect
consider
control
convince
cook
cool
copper
copy
coral
core
corn
correct
cost
cotton
couch
country
couple
course
cousin
cover
coyote
crack
cradle
craft
cram
crane
crash
crater
crawl
crazy
cream
credit
creek
crew
cricket
crime
crisp
critic
crop
cross
crouch
crowd
crucial
cruel
cruise
crumble
crunch
crush
cry
crystal
cube
culture
cup
cupboard
curious
current
curtain
curve
cushion
custom
cute
cycle
dad
damage
damp
dance
danger
daring
dash
daughter
dawn
day
deal
debate
debris
decade
december
decide
decline
decorate
decrease
deer
defense
define
defy
degree
delay
deliver
demand
demise
denial
dentist
deny
depart
depend
deposit
depth
deputy
derive
describe
desert
design
desk
despair
destroy
detail
detect
develop
device
devote
diagram
dial
diamond
diary
dice
diesel
diet
differ
digital
dignity
dilemma
dinner
dinosaur
direct
dirt
disagree
discover
disease
dish
dismiss
disorder
display
distance
divert
divide
divorce
dizzy
doctor
document
dog
doll
dolphin
domain
donate
donkey
donor
door
dose
double
dove
draft
dragon
drama
drastic
draw
dream
dress
drift
drill
drink
drip
drive
drop
drum
dry
duck
dumb
dune
during
dust
dutch
duty
dwarf
dynamic
eager
eagle
early
earn
earth
easily
east
easy
echo
ecology
economy
edge
edit
educate
effort
egg
eight
either
elbow
elder
electric
elegant
element
elephant
elevator
elite
else
embark
embody
embrace
emerge
emotion
employ
empower
empty
enable
enact
end
endless
endorse
enemy
energy
enforce
engage
engine
enhance
enjoy
enlist
enough
enrich
enroll
ensure
enter
entire
entry
envelope
episode
equal
equip
era
erase
erode
erosion
error
erupt
escape
essay
essence
estate
eternal
ethics
evidence
evil
evoke
evolve
exact
example
excess
exchange
excite
exclude
excuse
execute
exercise
exhaust
exhibit
exile
exist
exit
exotic
expand
expect
expire
explain
expose
express
extend
extra
eye
eyebrow
fabric
face
faculty
fade
faint
faith
fall
false
fame
family
famous
fan
fancy
fantasy
farm
fashion
fat
fatal
father
fatigue
fault
favorite
feature
february
federal
fee
feed
feel
female
fence
festival
fetch
fever
few
fiber
fiction
field
figure
file
film
filter
final
find
fine
finger
finish
fire
firm
first
fiscal
fish
fit
fitness
fix
flag
flame
flash
flat
flavor
flee
flight
flip
float
flock
floor
flower
fluid
flush
fly
foam
focus
fog
foil
fold
follow
food
foot
force
forest
forget
fork
fortune
forum
forward
fossil
foster
found
fox
fragile
frame
frequent
fresh
friend
fringe
frog
front
frost
frown
frozen
fruit
fuel
fun
funny
furnace
fury
future
gadget
gain
galaxy
gallery
game
gap
garage
garbage
garden
garlic
garment
gas
gasp
gate
gather
gauge
gaze
general
genius
genre
gentle
genuine
gesture
ghost
giant
gift
giggle
ginger
giraffe
girl
give
glad
glance
glare
glass
glide
glimpse
globe
gloom
glory
glove
glow
glue
goat
goddess
gold
good
goose
gorilla
gospel
gossip
govern
gown
grab
grace
grain
grant
grape
grass
gravity
great
green
grid
grief
grit
grocery
group
grow
grunt
guard
guess
guide
guilt
guitar
gun
gym
habit
hair
half
hammer
hamster
hand
happy
harbor
hard
harsh
harvest
hat
have
hawk
hazard
head
health
heart
heavy
hedgehog
height
hello
helmet
help
hen
hero
hidden
high
hill
hint
hip
hire
history
hobby
hockey
hold
hole
holiday
hollow
home
honey
hood
hope
horn
horror
horse
hospital
host
hotel
hour
hover
hub
huge
human
humble
humor
hundred
hungry
hunt
hurdle
hurry
hurt
husband
hybrid
ice
icon
idea
identify
idle
ignore
ill
illegal
illness
image
imitate
immense
immune
impact
impose
improve
impulse
inch
include
income
increase
index
indicate
indoor
industry
infant
inflict
inform
inhale
inherit
initial
inject
injury
inmate
inner
innocent
input
inquiry
insane
insect
inside
inspire
install
intact
interest
into
invest
invite
involve
iron
island
isolate
issue
item
ivory
jacket
jaguar
jar
jazz
jealous
jeans
jelly
jewel
job
join
joke
journey
joy
judge
juice
jump
jungle
junior
junk
just
kangaroo
keen
keep
ketchup
key
kick
kid
kidney
kind
kingdom
kiss
kit
kitchen
kite
kitten
kiwi
knee
knife
knock
know
lab
label
labor
ladder
lady
lake
lamp
language
laptop
large
later
latin
laugh
laundry
lava
law
lawn
lawsuit
layer
lazy
leader
leaf
learn
leave
lecture
left
leg
legal
legend
leisure
lemon
lend
length
lens
leopard
lesson
letter
level
liar
liberty
library
license
life
lift
light
like
limb
limit
link
lion
liquid
list
little
live
lizard
load
loan
lobster
local
lock
logic
lonely
long
loop
lottery
loud
lounge
love
loyal
lucky
luggage
lumber
lunar
lunch
luxury
lyrics
machine
mad
magic
magnet
maid
mail
main
major
make
mammal
man
manage
mandate
mango
mansion
manual
maple
marble
march
margin
marine
market
marriage
mask
mass
master
match
material
math
matrix
matter
maximum
maze
meadow
mean
measure
meat
mechanic
medal
media
melody
melt
member
memory
mention
menu
mercy
merge
merit
merry
mesh
message
metal
method
middle
midnight
milk
million
mimic
mind
minimum
minor
minute
miracle
mirror
misery
miss
mistake
mix
mixed
mixture
mobile
model
modify
mom
moment
monitor
monkey
monster
month
moon
moral
more
morning
mosquito
mother
motion
motor
mountain
mouse
move
movie
much
muffin
mule
multiply
muscle
museum
mushroom
music
must
mutual
myself
mystery
myth
naive
name
napkin
narrow
nasty
nation
nature
near
neck
need
negative
neglect
neither
nephew
nerve
nest
net
network
neutral
never
news
next
nice
night
noble
noise
nominee
noodle
normal
north
nose
notable
note
nothing
notice
novel
now
nuclear
number
nurse
nut
oak
obey
object
oblige
obscure
observe
obtain
obvious
occur
ocean
october
odor
off
offer
office
often
oil
okay
old
olive
olympic
omit
once
one
onion
online
only
open
opera
opinion
oppose
option
orange
orbit
orchard
order
ordinary
organ
orient
original
orphan
ostrich
other
outdoor
outer
output
outside
oval
oven
over
own
owner
oxygen
oyster
ozone
pact
paddle
page
pair
palace
palm
panda
panel
panic
panther
paper
parade
parent
park
parrot
party
pass
patch
path
patient
patrol
pattern
pause
pave
payment
peace
peanut
pear
peasant
pelican
pen
penalty
pencil
people
pepper
perfect
permit
person
pet
phone
photo
phrase
physical
piano
picnic
picture
piece
pig
pigeon
pill
pilot
pink
pioneer
pipe
pistol
pitch
pizza
place
planet
plastic
plate
play
please
pledge
pluck
plug
plunge
poem
poet
point
polar
pole
police
pond
pony
pool
popular
portion
position
possible
post
potato
pottery
poverty
powder
power
practice
praise
predict
prefer
prepare
present
pretty
prevent
price
pride
primary
print
priority
prison
private
prize
problem
process
produce
profit
program
project
promote
proof
property
prosper
protect
proud
provide
public
pudding
pull
pulp
pulse
pumpkin
punch
pupil
puppy
purchase
purity
purpose
purse
push
put
puzzle
pyramid
quality
quantum
quarter
question
quick
quit
quiz
quote
rabbit
raccoon
race
rack
radar
radio
rail
rain
raise
rally
ramp
ranch
random
range
rapid
rare
rate
rather
raven
raw
razor
ready
real
reason
rebel
rebuild
recall
receive
recipe
record
recycle
reduce
reflect
reform
refuse
region
regret
regular
reject
relax
release
relief
rely
remain
remember
remind
remove
render
renew
rent
reopen
repair
repeat
replace
report
require
rescue
resemble
resist
resource
response
result
retire
retreat
return
reunion
reveal
review
reward
rhythm
rib
ribbon
rice
rich
ride
ridge
rifle
right
rigid
ring
riot
ripple
risk
ritual
rival
river
road
roast
robot
robust
rocket
romance
roof
rookie
room
rose
rotate
rough
round
route
royal
rubber
rude
rug
rule
run
runway
rural
sad
saddle
sadness
safe
sail
salad
salmon
salon
salt
salute
same
sample
sand
satisfy
satoshi
sauce
sausage
save
say
scale
scan
scare
scatter
scene
scheme
school
science
scissors
scorpion
scout
scrap
screen
script
scrub
sea
search
season
seat
second
secret
section
security
seed
seek
segment
select
sell
seminar
senior
sense
sentence
series
service
session
settle
setup
seven
shadow
shaft
shallow
share
shed
shell
sheriff
shield
shift
shine
ship
shiver
shock
shoe
shoot
shop
short
shoulder
shove
shrimp
shrug
shuffle
shy
sibling
sick
side
siege
sight
sign
silent
silk
silly
silver
similar
simple
since
sing
siren
sister
situate
six
size
skate
sketch
ski
skill
skin
skirt
skull
slab
slam
sleep
slender
slice
slide
slight
slim
slogan
slot
slow
slush
small
smart
smile
smoke
smooth
snack
snake
snap
sniff
snow
soap
soccer
social
sock
soda
soft
solar
soldier
solid
solution
solve
someone
song
soon
sorry
sort
soul
sound
soup
source
south
space
spare
spatial
spawn
speak
special
speed
spell
spend
sphere
spice
spider
spike
spin
spirit
split
spoil
sponsor
spoon
sport
spot
spray
spread
spring
spy
square
squeeze
squirrel
stable
stadium
staff
stage
stairs
stamp
stand
start
state
stay
steak
steel
stem
step
stereo
stick
still
sting
stock
stomach
stone
stool
story
stove
strategy
street
strike
strong
struggle
student
stuff
stumble
style
subject
submit
subway
success
such
sudden
suffer
sugar
suggest
suit
summer
sun
sunny
sunset
super
supply
supreme
sure
surface
surge
surprise
surround
survey
suspect
sustain
swallow
swamp
swap
swarm
swear
sweet
swift
swim
swing
switch
sword
symbol
symptom
syrup
system
table
tackle
tag
tail
talent
talk
tank
tape
target
task
taste
tattoo
taxi
teach
team
tell
ten
tenant
tennis
tent
term
test
text
thank
that
theme
then
theory
there
they
thing
this
thought
three
thrive
throw
thumb
thunder
ticket
tide
tiger
tilt
timber
time
tiny
tip
tired
tissue
title
toast
tobacco
today
toddler
toe
together
toilet
token
tomato
tomorrow
tone
tongue
tonight
tool
tooth
top
topic
topple
torch
tornado
tortoise
toss
total
tourist
toward
tower
town
toy
track
trade
traffic
tragic
train
transfer
trap
trash
travel
tray
treat
tree
trend
trial
tribe
trick
trigger
trim
trip
trophy
trouble
truck
true
truly
trumpet
trust
truth
try
tube
tuition
tumble
tuna
tunnel
turkey
turn
turtle
twelve
twenty
twice
twin
twist
two
type
typical
ugly
umbrella
unable
unaware
uncle
uncover
under
undo
unfair
unfold
unhappy
uniform
unique
unit
universe
unknown
unlock
until
unusual
unveil
update
upgrade
uphold
upon
upper
upset
urban
urge
usage
use
used
useful
useless
usual
utility
vacant
vacuum
vague
valid
valley
valve
van
vanish
vapor
various
vast
vault
vehicle
velvet
vendor
venture
venue
verb
verify
version
very
vessel
veteran
viable
vibrant
vicious
victory
video
view
village
vintage
violin
virtual
virus
visa
visit
visual
vital
vivid
vocal
voice
void
volcano
volume
vote
voyage
wage
wagon
wait
walk
wall
walnut
want
warfare
warm
warrior
wash
wasp
waste
water
wave
way
wealth
weapon
wear
weasel
weather
web
wedding
weekend
weird
welcome
west
wet
whale
what
wheat
wheel
when
where
whip
whisper
wide
width
wife
wild
will
win
window
wine
wing
wink
winner
winter
wire
wisdom
wise
wish
witness
wolf
woman
wonder
wood
wool
word
work
world
worry
worth
wrap
wreck
wrestle
wrist
write
wrong
yard
year
yellow
you
young
youth
zebra
zero
zone
zoo`.split("\n");

  // node_modules/@noble/hashes/esm/ripemd160.js
  var Rho = /* @__PURE__ */ new Uint8Array([7, 4, 13, 1, 10, 6, 15, 3, 12, 0, 9, 5, 2, 14, 11, 8]);
  var Id = /* @__PURE__ */ new Uint8Array(new Array(16).fill(0).map((_, i) => i));
  var Pi = /* @__PURE__ */ Id.map((i) => (9 * i + 5) % 16);
  var idxL = [Id];
  var idxR = [Pi];
  for (let i = 0; i < 4; i++)
    for (let j of [idxL, idxR])
      j.push(j[i].map((k) => Rho[k]));
  var shifts = /* @__PURE__ */ [
    [11, 14, 15, 12, 5, 8, 7, 9, 11, 13, 14, 15, 6, 7, 9, 8],
    [12, 13, 11, 15, 6, 9, 9, 7, 12, 15, 11, 13, 7, 8, 7, 7],
    [13, 15, 14, 11, 7, 7, 6, 8, 13, 14, 13, 12, 5, 5, 6, 9],
    [14, 11, 12, 14, 8, 6, 5, 5, 15, 12, 15, 14, 9, 9, 8, 6],
    [15, 12, 13, 13, 9, 5, 8, 6, 14, 11, 12, 11, 8, 6, 5, 5]
  ].map((i) => new Uint8Array(i));
  var shiftsL = /* @__PURE__ */ idxL.map((idx, i) => idx.map((j) => shifts[i][j]));
  var shiftsR = /* @__PURE__ */ idxR.map((idx, i) => idx.map((j) => shifts[i][j]));
  var Kl = /* @__PURE__ */ new Uint32Array([
    0,
    1518500249,
    1859775393,
    2400959708,
    2840853838
  ]);
  var Kr = /* @__PURE__ */ new Uint32Array([
    1352829926,
    1548603684,
    1836072691,
    2053994217,
    0
  ]);
  function f(group, x, y, z) {
    if (group === 0)
      return x ^ y ^ z;
    else if (group === 1)
      return x & y | ~x & z;
    else if (group === 2)
      return (x | ~y) ^ z;
    else if (group === 3)
      return x & z | y & ~z;
    else
      return x ^ (y | ~z);
  }
  var R_BUF = /* @__PURE__ */ new Uint32Array(16);
  var RIPEMD160 = class extends HashMD {
    constructor() {
      super(64, 20, 8, true);
      this.h0 = 1732584193 | 0;
      this.h1 = 4023233417 | 0;
      this.h2 = 2562383102 | 0;
      this.h3 = 271733878 | 0;
      this.h4 = 3285377520 | 0;
    }
    get() {
      const { h0, h1, h2, h3, h4 } = this;
      return [h0, h1, h2, h3, h4];
    }
    set(h0, h1, h2, h3, h4) {
      this.h0 = h0 | 0;
      this.h1 = h1 | 0;
      this.h2 = h2 | 0;
      this.h3 = h3 | 0;
      this.h4 = h4 | 0;
    }
    process(view, offset) {
      for (let i = 0; i < 16; i++, offset += 4)
        R_BUF[i] = view.getUint32(offset, true);
      let al = this.h0 | 0, ar = al, bl = this.h1 | 0, br = bl, cl = this.h2 | 0, cr = cl, dl = this.h3 | 0, dr = dl, el = this.h4 | 0, er = el;
      for (let group = 0; group < 5; group++) {
        const rGroup = 4 - group;
        const hbl = Kl[group], hbr = Kr[group];
        const rl = idxL[group], rr = idxR[group];
        const sl = shiftsL[group], sr = shiftsR[group];
        for (let i = 0; i < 16; i++) {
          const tl = rotl(al + f(group, bl, cl, dl) + R_BUF[rl[i]] + hbl, sl[i]) + el | 0;
          al = el, el = dl, dl = rotl(cl, 10) | 0, cl = bl, bl = tl;
        }
        for (let i = 0; i < 16; i++) {
          const tr = rotl(ar + f(rGroup, br, cr, dr) + R_BUF[rr[i]] + hbr, sr[i]) + er | 0;
          ar = er, er = dr, dr = rotl(cr, 10) | 0, cr = br, br = tr;
        }
      }
      this.set(this.h1 + cl + dr | 0, this.h2 + dl + er | 0, this.h3 + el + ar | 0, this.h4 + al + br | 0, this.h0 + bl + cr | 0);
    }
    roundClean() {
      R_BUF.fill(0);
    }
    destroy() {
      this.destroyed = true;
      this.buffer.fill(0);
      this.set(0, 0, 0, 0, 0);
    }
  };
  var ripemd160 = /* @__PURE__ */ wrapConstructor(() => new RIPEMD160());

  // node_modules/@noble/curves/esm/abstract/utils.js
  var utils_exports = {};
  __export(utils_exports, {
    aInRange: () => aInRange,
    abool: () => abool,
    abytes: () => abytes,
    bitGet: () => bitGet,
    bitLen: () => bitLen,
    bitMask: () => bitMask,
    bitSet: () => bitSet,
    bytesToHex: () => bytesToHex2,
    bytesToNumberBE: () => bytesToNumberBE,
    bytesToNumberLE: () => bytesToNumberLE,
    concatBytes: () => concatBytes2,
    createHmacDrbg: () => createHmacDrbg,
    ensureBytes: () => ensureBytes,
    equalBytes: () => equalBytes,
    hexToBytes: () => hexToBytes2,
    hexToNumber: () => hexToNumber,
    inRange: () => inRange,
    isBytes: () => isBytes3,
    memoized: () => memoized,
    notImplemented: () => notImplemented,
    numberToBytesBE: () => numberToBytesBE,
    numberToBytesLE: () => numberToBytesLE,
    numberToHexUnpadded: () => numberToHexUnpadded,
    numberToVarBytesBE: () => numberToVarBytesBE,
    utf8ToBytes: () => utf8ToBytes2,
    validateObject: () => validateObject
  });
  var _0n = /* @__PURE__ */ BigInt(0);
  var _1n = /* @__PURE__ */ BigInt(1);
  var _2n = /* @__PURE__ */ BigInt(2);
  function isBytes3(a) {
    return a instanceof Uint8Array || a != null && typeof a === "object" && a.constructor.name === "Uint8Array";
  }
  function abytes(item) {
    if (!isBytes3(item))
      throw new Error("Uint8Array expected");
  }
  function abool(title, value) {
    if (typeof value !== "boolean")
      throw new Error(`${title} must be valid boolean, got "${value}".`);
  }
  var hexes2 = /* @__PURE__ */ Array.from({ length: 256 }, (_, i) => i.toString(16).padStart(2, "0"));
  function bytesToHex2(bytes2) {
    abytes(bytes2);
    let hex = "";
    for (let i = 0; i < bytes2.length; i++) {
      hex += hexes2[bytes2[i]];
    }
    return hex;
  }
  function numberToHexUnpadded(num) {
    const hex = num.toString(16);
    return hex.length & 1 ? `0${hex}` : hex;
  }
  function hexToNumber(hex) {
    if (typeof hex !== "string")
      throw new Error("hex string expected, got " + typeof hex);
    return BigInt(hex === "" ? "0" : `0x${hex}`);
  }
  var asciis2 = { _0: 48, _9: 57, _A: 65, _F: 70, _a: 97, _f: 102 };
  function asciiToBase162(char) {
    if (char >= asciis2._0 && char <= asciis2._9)
      return char - asciis2._0;
    if (char >= asciis2._A && char <= asciis2._F)
      return char - (asciis2._A - 10);
    if (char >= asciis2._a && char <= asciis2._f)
      return char - (asciis2._a - 10);
    return;
  }
  function hexToBytes2(hex) {
    if (typeof hex !== "string")
      throw new Error("hex string expected, got " + typeof hex);
    const hl = hex.length;
    const al = hl / 2;
    if (hl % 2)
      throw new Error("padded hex string expected, got unpadded hex of length " + hl);
    const array = new Uint8Array(al);
    for (let ai = 0, hi = 0; ai < al; ai++, hi += 2) {
      const n1 = asciiToBase162(hex.charCodeAt(hi));
      const n2 = asciiToBase162(hex.charCodeAt(hi + 1));
      if (n1 === void 0 || n2 === void 0) {
        const char = hex[hi] + hex[hi + 1];
        throw new Error('hex string expected, got non-hex character "' + char + '" at index ' + hi);
      }
      array[ai] = n1 * 16 + n2;
    }
    return array;
  }
  function bytesToNumberBE(bytes2) {
    return hexToNumber(bytesToHex2(bytes2));
  }
  function bytesToNumberLE(bytes2) {
    abytes(bytes2);
    return hexToNumber(bytesToHex2(Uint8Array.from(bytes2).reverse()));
  }
  function numberToBytesBE(n, len) {
    return hexToBytes2(n.toString(16).padStart(len * 2, "0"));
  }
  function numberToBytesLE(n, len) {
    return numberToBytesBE(n, len).reverse();
  }
  function numberToVarBytesBE(n) {
    return hexToBytes2(numberToHexUnpadded(n));
  }
  function ensureBytes(title, hex, expectedLength) {
    let res;
    if (typeof hex === "string") {
      try {
        res = hexToBytes2(hex);
      } catch (e) {
        throw new Error(`${title} must be valid hex string, got "${hex}". Cause: ${e}`);
      }
    } else if (isBytes3(hex)) {
      res = Uint8Array.from(hex);
    } else {
      throw new Error(`${title} must be hex string or Uint8Array`);
    }
    const len = res.length;
    if (typeof expectedLength === "number" && len !== expectedLength)
      throw new Error(`${title} expected ${expectedLength} bytes, got ${len}`);
    return res;
  }
  function concatBytes2(...arrays) {
    let sum = 0;
    for (let i = 0; i < arrays.length; i++) {
      const a = arrays[i];
      abytes(a);
      sum += a.length;
    }
    const res = new Uint8Array(sum);
    for (let i = 0, pad = 0; i < arrays.length; i++) {
      const a = arrays[i];
      res.set(a, pad);
      pad += a.length;
    }
    return res;
  }
  function equalBytes(a, b) {
    if (a.length !== b.length)
      return false;
    let diff = 0;
    for (let i = 0; i < a.length; i++)
      diff |= a[i] ^ b[i];
    return diff === 0;
  }
  function utf8ToBytes2(str) {
    if (typeof str !== "string")
      throw new Error(`utf8ToBytes expected string, got ${typeof str}`);
    return new Uint8Array(new TextEncoder().encode(str));
  }
  var isPosBig = (n) => typeof n === "bigint" && _0n <= n;
  function inRange(n, min, max) {
    return isPosBig(n) && isPosBig(min) && isPosBig(max) && min <= n && n < max;
  }
  function aInRange(title, n, min, max) {
    if (!inRange(n, min, max))
      throw new Error(`expected valid ${title}: ${min} <= n < ${max}, got ${typeof n} ${n}`);
  }
  function bitLen(n) {
    let len;
    for (len = 0; n > _0n; n >>= _1n, len += 1)
      ;
    return len;
  }
  function bitGet(n, pos) {
    return n >> BigInt(pos) & _1n;
  }
  function bitSet(n, pos, value) {
    return n | (value ? _1n : _0n) << BigInt(pos);
  }
  var bitMask = (n) => (_2n << BigInt(n - 1)) - _1n;
  var u8n = (data) => new Uint8Array(data);
  var u8fr = (arr) => Uint8Array.from(arr);
  function createHmacDrbg(hashLen, qByteLen, hmacFn) {
    if (typeof hashLen !== "number" || hashLen < 2)
      throw new Error("hashLen must be a number");
    if (typeof qByteLen !== "number" || qByteLen < 2)
      throw new Error("qByteLen must be a number");
    if (typeof hmacFn !== "function")
      throw new Error("hmacFn must be a function");
    let v = u8n(hashLen);
    let k = u8n(hashLen);
    let i = 0;
    const reset = () => {
      v.fill(1);
      k.fill(0);
      i = 0;
    };
    const h = (...b) => hmacFn(k, v, ...b);
    const reseed = (seed = u8n()) => {
      k = h(u8fr([0]), seed);
      v = h();
      if (seed.length === 0)
        return;
      k = h(u8fr([1]), seed);
      v = h();
    };
    const gen = () => {
      if (i++ >= 1e3)
        throw new Error("drbg: tried 1000 values");
      let len = 0;
      const out = [];
      while (len < qByteLen) {
        v = h();
        const sl = v.slice();
        out.push(sl);
        len += v.length;
      }
      return concatBytes2(...out);
    };
    const genUntil = (seed, pred) => {
      reset();
      reseed(seed);
      let res = void 0;
      while (!(res = pred(gen())))
        reseed();
      reset();
      return res;
    };
    return genUntil;
  }
  var validatorFns = {
    bigint: (val) => typeof val === "bigint",
    function: (val) => typeof val === "function",
    boolean: (val) => typeof val === "boolean",
    string: (val) => typeof val === "string",
    stringOrUint8Array: (val) => typeof val === "string" || isBytes3(val),
    isSafeInteger: (val) => Number.isSafeInteger(val),
    array: (val) => Array.isArray(val),
    field: (val, object) => object.Fp.isValid(val),
    hash: (val) => typeof val === "function" && Number.isSafeInteger(val.outputLen)
  };
  function validateObject(object, validators, optValidators = {}) {
    const checkField = (fieldName, type, isOptional) => {
      const checkVal = validatorFns[type];
      if (typeof checkVal !== "function")
        throw new Error(`Invalid validator "${type}", expected function`);
      const val = object[fieldName];
      if (isOptional && val === void 0)
        return;
      if (!checkVal(val, object)) {
        throw new Error(`Invalid param ${String(fieldName)}=${val} (${typeof val}), expected ${type}`);
      }
    };
    for (const [fieldName, type] of Object.entries(validators))
      checkField(fieldName, type, false);
    for (const [fieldName, type] of Object.entries(optValidators))
      checkField(fieldName, type, true);
    return object;
  }
  var notImplemented = () => {
    throw new Error("not implemented");
  };
  function memoized(fn) {
    const map = /* @__PURE__ */ new WeakMap();
    return (arg, ...args) => {
      const val = map.get(arg);
      if (val !== void 0)
        return val;
      const computed = fn(arg, ...args);
      map.set(arg, computed);
      return computed;
    };
  }

  // node_modules/@noble/curves/esm/abstract/modular.js
  var _0n2 = BigInt(0);
  var _1n2 = BigInt(1);
  var _2n2 = BigInt(2);
  var _3n = BigInt(3);
  var _4n = BigInt(4);
  var _5n = BigInt(5);
  var _8n = BigInt(8);
  var _9n = BigInt(9);
  var _16n = BigInt(16);
  function mod(a, b) {
    const result = a % b;
    return result >= _0n2 ? result : b + result;
  }
  function pow(num, power, modulo) {
    if (modulo <= _0n2 || power < _0n2)
      throw new Error("Expected power/modulo > 0");
    if (modulo === _1n2)
      return _0n2;
    let res = _1n2;
    while (power > _0n2) {
      if (power & _1n2)
        res = res * num % modulo;
      num = num * num % modulo;
      power >>= _1n2;
    }
    return res;
  }
  function pow2(x, power, modulo) {
    let res = x;
    while (power-- > _0n2) {
      res *= res;
      res %= modulo;
    }
    return res;
  }
  function invert(number2, modulo) {
    if (number2 === _0n2 || modulo <= _0n2) {
      throw new Error(`invert: expected positive integers, got n=${number2} mod=${modulo}`);
    }
    let a = mod(number2, modulo);
    let b = modulo;
    let x = _0n2, y = _1n2, u = _1n2, v = _0n2;
    while (a !== _0n2) {
      const q = b / a;
      const r = b % a;
      const m = x - u * q;
      const n = y - v * q;
      b = a, a = r, x = u, y = v, u = m, v = n;
    }
    const gcd2 = b;
    if (gcd2 !== _1n2)
      throw new Error("invert: does not exist");
    return mod(x, modulo);
  }
  function tonelliShanks(P) {
    const legendreC = (P - _1n2) / _2n2;
    let Q, S, Z;
    for (Q = P - _1n2, S = 0; Q % _2n2 === _0n2; Q /= _2n2, S++)
      ;
    for (Z = _2n2; Z < P && pow(Z, legendreC, P) !== P - _1n2; Z++)
      ;
    if (S === 1) {
      const p1div4 = (P + _1n2) / _4n;
      return function tonelliFast(Fp2, n) {
        const root = Fp2.pow(n, p1div4);
        if (!Fp2.eql(Fp2.sqr(root), n))
          throw new Error("Cannot find square root");
        return root;
      };
    }
    const Q1div2 = (Q + _1n2) / _2n2;
    return function tonelliSlow(Fp2, n) {
      if (Fp2.pow(n, legendreC) === Fp2.neg(Fp2.ONE))
        throw new Error("Cannot find square root");
      let r = S;
      let g = Fp2.pow(Fp2.mul(Fp2.ONE, Z), Q);
      let x = Fp2.pow(n, Q1div2);
      let b = Fp2.pow(n, Q);
      while (!Fp2.eql(b, Fp2.ONE)) {
        if (Fp2.eql(b, Fp2.ZERO))
          return Fp2.ZERO;
        let m = 1;
        for (let t2 = Fp2.sqr(b); m < r; m++) {
          if (Fp2.eql(t2, Fp2.ONE))
            break;
          t2 = Fp2.sqr(t2);
        }
        const ge = Fp2.pow(g, _1n2 << BigInt(r - m - 1));
        g = Fp2.sqr(ge);
        x = Fp2.mul(x, ge);
        b = Fp2.mul(b, g);
        r = m;
      }
      return x;
    };
  }
  function FpSqrt(P) {
    if (P % _4n === _3n) {
      const p1div4 = (P + _1n2) / _4n;
      return function sqrt3mod4(Fp2, n) {
        const root = Fp2.pow(n, p1div4);
        if (!Fp2.eql(Fp2.sqr(root), n))
          throw new Error("Cannot find square root");
        return root;
      };
    }
    if (P % _8n === _5n) {
      const c1 = (P - _5n) / _8n;
      return function sqrt5mod8(Fp2, n) {
        const n2 = Fp2.mul(n, _2n2);
        const v = Fp2.pow(n2, c1);
        const nv = Fp2.mul(n, v);
        const i = Fp2.mul(Fp2.mul(nv, _2n2), v);
        const root = Fp2.mul(nv, Fp2.sub(i, Fp2.ONE));
        if (!Fp2.eql(Fp2.sqr(root), n))
          throw new Error("Cannot find square root");
        return root;
      };
    }
    if (P % _16n === _9n) {
    }
    return tonelliShanks(P);
  }
  var FIELD_FIELDS = [
    "create",
    "isValid",
    "is0",
    "neg",
    "inv",
    "sqrt",
    "sqr",
    "eql",
    "add",
    "sub",
    "mul",
    "pow",
    "div",
    "addN",
    "subN",
    "mulN",
    "sqrN"
  ];
  function validateField(field) {
    const initial = {
      ORDER: "bigint",
      MASK: "bigint",
      BYTES: "isSafeInteger",
      BITS: "isSafeInteger"
    };
    const opts = FIELD_FIELDS.reduce((map, val) => {
      map[val] = "function";
      return map;
    }, initial);
    return validateObject(field, opts);
  }
  function FpPow(f2, num, power) {
    if (power < _0n2)
      throw new Error("Expected power > 0");
    if (power === _0n2)
      return f2.ONE;
    if (power === _1n2)
      return num;
    let p = f2.ONE;
    let d = num;
    while (power > _0n2) {
      if (power & _1n2)
        p = f2.mul(p, d);
      d = f2.sqr(d);
      power >>= _1n2;
    }
    return p;
  }
  function FpInvertBatch(f2, nums) {
    const tmp = new Array(nums.length);
    const lastMultiplied = nums.reduce((acc, num, i) => {
      if (f2.is0(num))
        return acc;
      tmp[i] = acc;
      return f2.mul(acc, num);
    }, f2.ONE);
    const inverted = f2.inv(lastMultiplied);
    nums.reduceRight((acc, num, i) => {
      if (f2.is0(num))
        return acc;
      tmp[i] = f2.mul(acc, tmp[i]);
      return f2.mul(acc, num);
    }, inverted);
    return tmp;
  }
  function nLength(n, nBitLength) {
    const _nBitLength = nBitLength !== void 0 ? nBitLength : n.toString(2).length;
    const nByteLength = Math.ceil(_nBitLength / 8);
    return { nBitLength: _nBitLength, nByteLength };
  }
  function Field(ORDER, bitLen2, isLE2 = false, redef = {}) {
    if (ORDER <= _0n2)
      throw new Error(`Expected Field ORDER > 0, got ${ORDER}`);
    const { nBitLength: BITS, nByteLength: BYTES } = nLength(ORDER, bitLen2);
    if (BYTES > 2048)
      throw new Error("Field lengths over 2048 bytes are not supported");
    const sqrtP = FpSqrt(ORDER);
    const f2 = Object.freeze({
      ORDER,
      BITS,
      BYTES,
      MASK: bitMask(BITS),
      ZERO: _0n2,
      ONE: _1n2,
      create: (num) => mod(num, ORDER),
      isValid: (num) => {
        if (typeof num !== "bigint")
          throw new Error(`Invalid field element: expected bigint, got ${typeof num}`);
        return _0n2 <= num && num < ORDER;
      },
      is0: (num) => num === _0n2,
      isOdd: (num) => (num & _1n2) === _1n2,
      neg: (num) => mod(-num, ORDER),
      eql: (lhs, rhs) => lhs === rhs,
      sqr: (num) => mod(num * num, ORDER),
      add: (lhs, rhs) => mod(lhs + rhs, ORDER),
      sub: (lhs, rhs) => mod(lhs - rhs, ORDER),
      mul: (lhs, rhs) => mod(lhs * rhs, ORDER),
      pow: (num, power) => FpPow(f2, num, power),
      div: (lhs, rhs) => mod(lhs * invert(rhs, ORDER), ORDER),
      // Same as above, but doesn't normalize
      sqrN: (num) => num * num,
      addN: (lhs, rhs) => lhs + rhs,
      subN: (lhs, rhs) => lhs - rhs,
      mulN: (lhs, rhs) => lhs * rhs,
      inv: (num) => invert(num, ORDER),
      sqrt: redef.sqrt || ((n) => sqrtP(f2, n)),
      invertBatch: (lst) => FpInvertBatch(f2, lst),
      // TODO: do we really need constant cmov?
      // We don't have const-time bigints anyway, so probably will be not very useful
      cmov: (a, b, c) => c ? b : a,
      toBytes: (num) => isLE2 ? numberToBytesLE(num, BYTES) : numberToBytesBE(num, BYTES),
      fromBytes: (bytes2) => {
        if (bytes2.length !== BYTES)
          throw new Error(`Fp.fromBytes: expected ${BYTES}, got ${bytes2.length}`);
        return isLE2 ? bytesToNumberLE(bytes2) : bytesToNumberBE(bytes2);
      }
    });
    return Object.freeze(f2);
  }
  function getFieldBytesLength(fieldOrder) {
    if (typeof fieldOrder !== "bigint")
      throw new Error("field order must be bigint");
    const bitLength = fieldOrder.toString(2).length;
    return Math.ceil(bitLength / 8);
  }
  function getMinHashLength(fieldOrder) {
    const length = getFieldBytesLength(fieldOrder);
    return length + Math.ceil(length / 2);
  }
  function mapHashToField(key, fieldOrder, isLE2 = false) {
    const len = key.length;
    const fieldLen = getFieldBytesLength(fieldOrder);
    const minLen = getMinHashLength(fieldOrder);
    if (len < 16 || len < minLen || len > 1024)
      throw new Error(`expected ${minLen}-1024 bytes of input, got ${len}`);
    const num = isLE2 ? bytesToNumberBE(key) : bytesToNumberLE(key);
    const reduced = mod(num, fieldOrder - _1n2) + _1n2;
    return isLE2 ? numberToBytesLE(reduced, fieldLen) : numberToBytesBE(reduced, fieldLen);
  }

  // node_modules/@noble/curves/esm/abstract/curve.js
  var _0n3 = BigInt(0);
  var _1n3 = BigInt(1);
  var pointPrecomputes = /* @__PURE__ */ new WeakMap();
  var pointWindowSizes = /* @__PURE__ */ new WeakMap();
  function wNAF(c, bits) {
    const constTimeNegate = (condition, item) => {
      const neg = item.negate();
      return condition ? neg : item;
    };
    const validateW = (W) => {
      if (!Number.isSafeInteger(W) || W <= 0 || W > bits)
        throw new Error(`Wrong window size=${W}, should be [1..${bits}]`);
    };
    const opts = (W) => {
      validateW(W);
      const windows = Math.ceil(bits / W) + 1;
      const windowSize = 2 ** (W - 1);
      return { windows, windowSize };
    };
    return {
      constTimeNegate,
      // non-const time multiplication ladder
      unsafeLadder(elm, n) {
        let p = c.ZERO;
        let d = elm;
        while (n > _0n3) {
          if (n & _1n3)
            p = p.add(d);
          d = d.double();
          n >>= _1n3;
        }
        return p;
      },
      /**
       * Creates a wNAF precomputation window. Used for caching.
       * Default window size is set by `utils.precompute()` and is equal to 8.
       * Number of precomputed points depends on the curve size:
       * 2^(𝑊−1) * (Math.ceil(𝑛 / 𝑊) + 1), where:
       * - 𝑊 is the window size
       * - 𝑛 is the bitlength of the curve order.
       * For a 256-bit curve and window size 8, the number of precomputed points is 128 * 33 = 4224.
       * @returns precomputed point tables flattened to a single array
       */
      precomputeWindow(elm, W) {
        const { windows, windowSize } = opts(W);
        const points = [];
        let p = elm;
        let base = p;
        for (let window = 0; window < windows; window++) {
          base = p;
          points.push(base);
          for (let i = 1; i < windowSize; i++) {
            base = base.add(p);
            points.push(base);
          }
          p = base.double();
        }
        return points;
      },
      /**
       * Implements ec multiplication using precomputed tables and w-ary non-adjacent form.
       * @param W window size
       * @param precomputes precomputed tables
       * @param n scalar (we don't check here, but should be less than curve order)
       * @returns real and fake (for const-time) points
       */
      wNAF(W, precomputes, n) {
        const { windows, windowSize } = opts(W);
        let p = c.ZERO;
        let f2 = c.BASE;
        const mask = BigInt(2 ** W - 1);
        const maxNumber = 2 ** W;
        const shiftBy = BigInt(W);
        for (let window = 0; window < windows; window++) {
          const offset = window * windowSize;
          let wbits = Number(n & mask);
          n >>= shiftBy;
          if (wbits > windowSize) {
            wbits -= maxNumber;
            n += _1n3;
          }
          const offset1 = offset;
          const offset2 = offset + Math.abs(wbits) - 1;
          const cond1 = window % 2 !== 0;
          const cond2 = wbits < 0;
          if (wbits === 0) {
            f2 = f2.add(constTimeNegate(cond1, precomputes[offset1]));
          } else {
            p = p.add(constTimeNegate(cond2, precomputes[offset2]));
          }
        }
        return { p, f: f2 };
      },
      wNAFCached(P, n, transform) {
        const W = pointWindowSizes.get(P) || 1;
        let comp = pointPrecomputes.get(P);
        if (!comp) {
          comp = this.precomputeWindow(P, W);
          if (W !== 1)
            pointPrecomputes.set(P, transform(comp));
        }
        return this.wNAF(W, comp, n);
      },
      // We calculate precomputes for elliptic curve point multiplication
      // using windowed method. This specifies window size and
      // stores precomputed values. Usually only base point would be precomputed.
      setWindowSize(P, W) {
        validateW(W);
        pointWindowSizes.set(P, W);
        pointPrecomputes.delete(P);
      }
    };
  }
  function pippenger(c, field, points, scalars) {
    if (!Array.isArray(points) || !Array.isArray(scalars) || scalars.length !== points.length)
      throw new Error("arrays of points and scalars must have equal length");
    scalars.forEach((s, i) => {
      if (!field.isValid(s))
        throw new Error(`wrong scalar at index ${i}`);
    });
    points.forEach((p, i) => {
      if (!(p instanceof c))
        throw new Error(`wrong point at index ${i}`);
    });
    const wbits = bitLen(BigInt(points.length));
    const windowSize = wbits > 12 ? wbits - 3 : wbits > 4 ? wbits - 2 : wbits ? 2 : 1;
    const MASK = (1 << windowSize) - 1;
    const buckets = new Array(MASK + 1).fill(c.ZERO);
    const lastBits = Math.floor((field.BITS - 1) / windowSize) * windowSize;
    let sum = c.ZERO;
    for (let i = lastBits; i >= 0; i -= windowSize) {
      buckets.fill(c.ZERO);
      for (let j = 0; j < scalars.length; j++) {
        const scalar = scalars[j];
        const wbits2 = Number(scalar >> BigInt(i) & BigInt(MASK));
        buckets[wbits2] = buckets[wbits2].add(points[j]);
      }
      let resI = c.ZERO;
      for (let j = buckets.length - 1, sumI = c.ZERO; j > 0; j--) {
        sumI = sumI.add(buckets[j]);
        resI = resI.add(sumI);
      }
      sum = sum.add(resI);
      if (i !== 0)
        for (let j = 0; j < windowSize; j++)
          sum = sum.double();
    }
    return sum;
  }
  function validateBasic(curve) {
    validateField(curve.Fp);
    validateObject(curve, {
      n: "bigint",
      h: "bigint",
      Gx: "field",
      Gy: "field"
    }, {
      nBitLength: "isSafeInteger",
      nByteLength: "isSafeInteger"
    });
    return Object.freeze({
      ...nLength(curve.n, curve.nBitLength),
      ...curve,
      ...{ p: curve.Fp.ORDER }
    });
  }

  // node_modules/@noble/curves/esm/abstract/weierstrass.js
  function validateSigVerOpts(opts) {
    if (opts.lowS !== void 0)
      abool("lowS", opts.lowS);
    if (opts.prehash !== void 0)
      abool("prehash", opts.prehash);
  }
  function validatePointOpts(curve) {
    const opts = validateBasic(curve);
    validateObject(opts, {
      a: "field",
      b: "field"
    }, {
      allowedPrivateKeyLengths: "array",
      wrapPrivateKey: "boolean",
      isTorsionFree: "function",
      clearCofactor: "function",
      allowInfinityPoint: "boolean",
      fromBytes: "function",
      toBytes: "function"
    });
    const { endo, Fp: Fp2, a } = opts;
    if (endo) {
      if (!Fp2.eql(a, Fp2.ZERO)) {
        throw new Error("Endomorphism can only be defined for Koblitz curves that have a=0");
      }
      if (typeof endo !== "object" || typeof endo.beta !== "bigint" || typeof endo.splitScalar !== "function") {
        throw new Error("Expected endomorphism with beta: bigint and splitScalar: function");
      }
    }
    return Object.freeze({ ...opts });
  }
  var { bytesToNumberBE: b2n, hexToBytes: h2b } = utils_exports;
  var DER = {
    // asn.1 DER encoding utils
    Err: class DERErr extends Error {
      constructor(m = "") {
        super(m);
      }
    },
    // Basic building block is TLV (Tag-Length-Value)
    _tlv: {
      encode: (tag, data) => {
        const { Err: E } = DER;
        if (tag < 0 || tag > 256)
          throw new E("tlv.encode: wrong tag");
        if (data.length & 1)
          throw new E("tlv.encode: unpadded data");
        const dataLen = data.length / 2;
        const len = numberToHexUnpadded(dataLen);
        if (len.length / 2 & 128)
          throw new E("tlv.encode: long form length too big");
        const lenLen = dataLen > 127 ? numberToHexUnpadded(len.length / 2 | 128) : "";
        return `${numberToHexUnpadded(tag)}${lenLen}${len}${data}`;
      },
      // v - value, l - left bytes (unparsed)
      decode(tag, data) {
        const { Err: E } = DER;
        let pos = 0;
        if (tag < 0 || tag > 256)
          throw new E("tlv.encode: wrong tag");
        if (data.length < 2 || data[pos++] !== tag)
          throw new E("tlv.decode: wrong tlv");
        const first = data[pos++];
        const isLong = !!(first & 128);
        let length = 0;
        if (!isLong)
          length = first;
        else {
          const lenLen = first & 127;
          if (!lenLen)
            throw new E("tlv.decode(long): indefinite length not supported");
          if (lenLen > 4)
            throw new E("tlv.decode(long): byte length is too big");
          const lengthBytes = data.subarray(pos, pos + lenLen);
          if (lengthBytes.length !== lenLen)
            throw new E("tlv.decode: length bytes not complete");
          if (lengthBytes[0] === 0)
            throw new E("tlv.decode(long): zero leftmost byte");
          for (const b of lengthBytes)
            length = length << 8 | b;
          pos += lenLen;
          if (length < 128)
            throw new E("tlv.decode(long): not minimal encoding");
        }
        const v = data.subarray(pos, pos + length);
        if (v.length !== length)
          throw new E("tlv.decode: wrong value length");
        return { v, l: data.subarray(pos + length) };
      }
    },
    // https://crypto.stackexchange.com/a/57734 Leftmost bit of first byte is 'negative' flag,
    // since we always use positive integers here. It must always be empty:
    // - add zero byte if exists
    // - if next byte doesn't have a flag, leading zero is not allowed (minimal encoding)
    _int: {
      encode(num) {
        const { Err: E } = DER;
        if (num < _0n4)
          throw new E("integer: negative integers are not allowed");
        let hex = numberToHexUnpadded(num);
        if (Number.parseInt(hex[0], 16) & 8)
          hex = "00" + hex;
        if (hex.length & 1)
          throw new E("unexpected assertion");
        return hex;
      },
      decode(data) {
        const { Err: E } = DER;
        if (data[0] & 128)
          throw new E("Invalid signature integer: negative");
        if (data[0] === 0 && !(data[1] & 128))
          throw new E("Invalid signature integer: unnecessary leading zero");
        return b2n(data);
      }
    },
    toSig(hex) {
      const { Err: E, _int: int, _tlv: tlv } = DER;
      const data = typeof hex === "string" ? h2b(hex) : hex;
      abytes(data);
      const { v: seqBytes, l: seqLeftBytes } = tlv.decode(48, data);
      if (seqLeftBytes.length)
        throw new E("Invalid signature: left bytes after parsing");
      const { v: rBytes, l: rLeftBytes } = tlv.decode(2, seqBytes);
      const { v: sBytes, l: sLeftBytes } = tlv.decode(2, rLeftBytes);
      if (sLeftBytes.length)
        throw new E("Invalid signature: left bytes after parsing");
      return { r: int.decode(rBytes), s: int.decode(sBytes) };
    },
    hexFromSig(sig) {
      const { _tlv: tlv, _int: int } = DER;
      const seq = `${tlv.encode(2, int.encode(sig.r))}${tlv.encode(2, int.encode(sig.s))}`;
      return tlv.encode(48, seq);
    }
  };
  var _0n4 = BigInt(0);
  var _1n4 = BigInt(1);
  var _2n3 = BigInt(2);
  var _3n2 = BigInt(3);
  var _4n2 = BigInt(4);
  function weierstrassPoints(opts) {
    const CURVE = validatePointOpts(opts);
    const { Fp: Fp2 } = CURVE;
    const Fn = Field(CURVE.n, CURVE.nBitLength);
    const toBytes2 = CURVE.toBytes || ((_c, point, _isCompressed) => {
      const a = point.toAffine();
      return concatBytes2(Uint8Array.from([4]), Fp2.toBytes(a.x), Fp2.toBytes(a.y));
    });
    const fromBytes = CURVE.fromBytes || ((bytes2) => {
      const tail = bytes2.subarray(1);
      const x = Fp2.fromBytes(tail.subarray(0, Fp2.BYTES));
      const y = Fp2.fromBytes(tail.subarray(Fp2.BYTES, 2 * Fp2.BYTES));
      return { x, y };
    });
    function weierstrassEquation(x) {
      const { a, b } = CURVE;
      const x2 = Fp2.sqr(x);
      const x3 = Fp2.mul(x2, x);
      return Fp2.add(Fp2.add(x3, Fp2.mul(x, a)), b);
    }
    if (!Fp2.eql(Fp2.sqr(CURVE.Gy), weierstrassEquation(CURVE.Gx)))
      throw new Error("bad generator point: equation left != right");
    function isWithinCurveOrder(num) {
      return inRange(num, _1n4, CURVE.n);
    }
    function normPrivateKeyToScalar(key) {
      const { allowedPrivateKeyLengths: lengths, nByteLength, wrapPrivateKey, n: N } = CURVE;
      if (lengths && typeof key !== "bigint") {
        if (isBytes3(key))
          key = bytesToHex2(key);
        if (typeof key !== "string" || !lengths.includes(key.length))
          throw new Error("Invalid key");
        key = key.padStart(nByteLength * 2, "0");
      }
      let num;
      try {
        num = typeof key === "bigint" ? key : bytesToNumberBE(ensureBytes("private key", key, nByteLength));
      } catch (error) {
        throw new Error(`private key must be ${nByteLength} bytes, hex or bigint, not ${typeof key}`);
      }
      if (wrapPrivateKey)
        num = mod(num, N);
      aInRange("private key", num, _1n4, N);
      return num;
    }
    function assertPrjPoint(other) {
      if (!(other instanceof Point3))
        throw new Error("ProjectivePoint expected");
    }
    const toAffineMemo = memoized((p, iz) => {
      const { px: x, py: y, pz: z } = p;
      if (Fp2.eql(z, Fp2.ONE))
        return { x, y };
      const is0 = p.is0();
      if (iz == null)
        iz = is0 ? Fp2.ONE : Fp2.inv(z);
      const ax = Fp2.mul(x, iz);
      const ay = Fp2.mul(y, iz);
      const zz = Fp2.mul(z, iz);
      if (is0)
        return { x: Fp2.ZERO, y: Fp2.ZERO };
      if (!Fp2.eql(zz, Fp2.ONE))
        throw new Error("invZ was invalid");
      return { x: ax, y: ay };
    });
    const assertValidMemo = memoized((p) => {
      if (p.is0()) {
        if (CURVE.allowInfinityPoint && !Fp2.is0(p.py))
          return;
        throw new Error("bad point: ZERO");
      }
      const { x, y } = p.toAffine();
      if (!Fp2.isValid(x) || !Fp2.isValid(y))
        throw new Error("bad point: x or y not FE");
      const left = Fp2.sqr(y);
      const right = weierstrassEquation(x);
      if (!Fp2.eql(left, right))
        throw new Error("bad point: equation left != right");
      if (!p.isTorsionFree())
        throw new Error("bad point: not in prime-order subgroup");
      return true;
    });
    class Point3 {
      constructor(px, py, pz) {
        this.px = px;
        this.py = py;
        this.pz = pz;
        if (px == null || !Fp2.isValid(px))
          throw new Error("x required");
        if (py == null || !Fp2.isValid(py))
          throw new Error("y required");
        if (pz == null || !Fp2.isValid(pz))
          throw new Error("z required");
        Object.freeze(this);
      }
      // Does not validate if the point is on-curve.
      // Use fromHex instead, or call assertValidity() later.
      static fromAffine(p) {
        const { x, y } = p || {};
        if (!p || !Fp2.isValid(x) || !Fp2.isValid(y))
          throw new Error("invalid affine point");
        if (p instanceof Point3)
          throw new Error("projective point not allowed");
        const is0 = (i) => Fp2.eql(i, Fp2.ZERO);
        if (is0(x) && is0(y))
          return Point3.ZERO;
        return new Point3(x, y, Fp2.ONE);
      }
      get x() {
        return this.toAffine().x;
      }
      get y() {
        return this.toAffine().y;
      }
      /**
       * Takes a bunch of Projective Points but executes only one
       * inversion on all of them. Inversion is very slow operation,
       * so this improves performance massively.
       * Optimization: converts a list of projective points to a list of identical points with Z=1.
       */
      static normalizeZ(points) {
        const toInv = Fp2.invertBatch(points.map((p) => p.pz));
        return points.map((p, i) => p.toAffine(toInv[i])).map(Point3.fromAffine);
      }
      /**
       * Converts hash string or Uint8Array to Point.
       * @param hex short/long ECDSA hex
       */
      static fromHex(hex) {
        const P = Point3.fromAffine(fromBytes(ensureBytes("pointHex", hex)));
        P.assertValidity();
        return P;
      }
      // Multiplies generator point by privateKey.
      static fromPrivateKey(privateKey) {
        return Point3.BASE.multiply(normPrivateKeyToScalar(privateKey));
      }
      // Multiscalar Multiplication
      static msm(points, scalars) {
        return pippenger(Point3, Fn, points, scalars);
      }
      // "Private method", don't use it directly
      _setWindowSize(windowSize) {
        wnaf.setWindowSize(this, windowSize);
      }
      // A point on curve is valid if it conforms to equation.
      assertValidity() {
        assertValidMemo(this);
      }
      hasEvenY() {
        const { y } = this.toAffine();
        if (Fp2.isOdd)
          return !Fp2.isOdd(y);
        throw new Error("Field doesn't support isOdd");
      }
      /**
       * Compare one point to another.
       */
      equals(other) {
        assertPrjPoint(other);
        const { px: X1, py: Y1, pz: Z1 } = this;
        const { px: X2, py: Y2, pz: Z2 } = other;
        const U1 = Fp2.eql(Fp2.mul(X1, Z2), Fp2.mul(X2, Z1));
        const U2 = Fp2.eql(Fp2.mul(Y1, Z2), Fp2.mul(Y2, Z1));
        return U1 && U2;
      }
      /**
       * Flips point to one corresponding to (x, -y) in Affine coordinates.
       */
      negate() {
        return new Point3(this.px, Fp2.neg(this.py), this.pz);
      }
      // Renes-Costello-Batina exception-free doubling formula.
      // There is 30% faster Jacobian formula, but it is not complete.
      // https://eprint.iacr.org/2015/1060, algorithm 3
      // Cost: 8M + 3S + 3*a + 2*b3 + 15add.
      double() {
        const { a, b } = CURVE;
        const b3 = Fp2.mul(b, _3n2);
        const { px: X1, py: Y1, pz: Z1 } = this;
        let X3 = Fp2.ZERO, Y3 = Fp2.ZERO, Z3 = Fp2.ZERO;
        let t0 = Fp2.mul(X1, X1);
        let t1 = Fp2.mul(Y1, Y1);
        let t2 = Fp2.mul(Z1, Z1);
        let t3 = Fp2.mul(X1, Y1);
        t3 = Fp2.add(t3, t3);
        Z3 = Fp2.mul(X1, Z1);
        Z3 = Fp2.add(Z3, Z3);
        X3 = Fp2.mul(a, Z3);
        Y3 = Fp2.mul(b3, t2);
        Y3 = Fp2.add(X3, Y3);
        X3 = Fp2.sub(t1, Y3);
        Y3 = Fp2.add(t1, Y3);
        Y3 = Fp2.mul(X3, Y3);
        X3 = Fp2.mul(t3, X3);
        Z3 = Fp2.mul(b3, Z3);
        t2 = Fp2.mul(a, t2);
        t3 = Fp2.sub(t0, t2);
        t3 = Fp2.mul(a, t3);
        t3 = Fp2.add(t3, Z3);
        Z3 = Fp2.add(t0, t0);
        t0 = Fp2.add(Z3, t0);
        t0 = Fp2.add(t0, t2);
        t0 = Fp2.mul(t0, t3);
        Y3 = Fp2.add(Y3, t0);
        t2 = Fp2.mul(Y1, Z1);
        t2 = Fp2.add(t2, t2);
        t0 = Fp2.mul(t2, t3);
        X3 = Fp2.sub(X3, t0);
        Z3 = Fp2.mul(t2, t1);
        Z3 = Fp2.add(Z3, Z3);
        Z3 = Fp2.add(Z3, Z3);
        return new Point3(X3, Y3, Z3);
      }
      // Renes-Costello-Batina exception-free addition formula.
      // There is 30% faster Jacobian formula, but it is not complete.
      // https://eprint.iacr.org/2015/1060, algorithm 1
      // Cost: 12M + 0S + 3*a + 3*b3 + 23add.
      add(other) {
        assertPrjPoint(other);
        const { px: X1, py: Y1, pz: Z1 } = this;
        const { px: X2, py: Y2, pz: Z2 } = other;
        let X3 = Fp2.ZERO, Y3 = Fp2.ZERO, Z3 = Fp2.ZERO;
        const a = CURVE.a;
        const b3 = Fp2.mul(CURVE.b, _3n2);
        let t0 = Fp2.mul(X1, X2);
        let t1 = Fp2.mul(Y1, Y2);
        let t2 = Fp2.mul(Z1, Z2);
        let t3 = Fp2.add(X1, Y1);
        let t4 = Fp2.add(X2, Y2);
        t3 = Fp2.mul(t3, t4);
        t4 = Fp2.add(t0, t1);
        t3 = Fp2.sub(t3, t4);
        t4 = Fp2.add(X1, Z1);
        let t5 = Fp2.add(X2, Z2);
        t4 = Fp2.mul(t4, t5);
        t5 = Fp2.add(t0, t2);
        t4 = Fp2.sub(t4, t5);
        t5 = Fp2.add(Y1, Z1);
        X3 = Fp2.add(Y2, Z2);
        t5 = Fp2.mul(t5, X3);
        X3 = Fp2.add(t1, t2);
        t5 = Fp2.sub(t5, X3);
        Z3 = Fp2.mul(a, t4);
        X3 = Fp2.mul(b3, t2);
        Z3 = Fp2.add(X3, Z3);
        X3 = Fp2.sub(t1, Z3);
        Z3 = Fp2.add(t1, Z3);
        Y3 = Fp2.mul(X3, Z3);
        t1 = Fp2.add(t0, t0);
        t1 = Fp2.add(t1, t0);
        t2 = Fp2.mul(a, t2);
        t4 = Fp2.mul(b3, t4);
        t1 = Fp2.add(t1, t2);
        t2 = Fp2.sub(t0, t2);
        t2 = Fp2.mul(a, t2);
        t4 = Fp2.add(t4, t2);
        t0 = Fp2.mul(t1, t4);
        Y3 = Fp2.add(Y3, t0);
        t0 = Fp2.mul(t5, t4);
        X3 = Fp2.mul(t3, X3);
        X3 = Fp2.sub(X3, t0);
        t0 = Fp2.mul(t3, t1);
        Z3 = Fp2.mul(t5, Z3);
        Z3 = Fp2.add(Z3, t0);
        return new Point3(X3, Y3, Z3);
      }
      subtract(other) {
        return this.add(other.negate());
      }
      is0() {
        return this.equals(Point3.ZERO);
      }
      wNAF(n) {
        return wnaf.wNAFCached(this, n, Point3.normalizeZ);
      }
      /**
       * Non-constant-time multiplication. Uses double-and-add algorithm.
       * It's faster, but should only be used when you don't care about
       * an exposed private key e.g. sig verification, which works over *public* keys.
       */
      multiplyUnsafe(sc) {
        aInRange("scalar", sc, _0n4, CURVE.n);
        const I = Point3.ZERO;
        if (sc === _0n4)
          return I;
        if (sc === _1n4)
          return this;
        const { endo } = CURVE;
        if (!endo)
          return wnaf.unsafeLadder(this, sc);
        let { k1neg, k1, k2neg, k2 } = endo.splitScalar(sc);
        let k1p = I;
        let k2p = I;
        let d = this;
        while (k1 > _0n4 || k2 > _0n4) {
          if (k1 & _1n4)
            k1p = k1p.add(d);
          if (k2 & _1n4)
            k2p = k2p.add(d);
          d = d.double();
          k1 >>= _1n4;
          k2 >>= _1n4;
        }
        if (k1neg)
          k1p = k1p.negate();
        if (k2neg)
          k2p = k2p.negate();
        k2p = new Point3(Fp2.mul(k2p.px, endo.beta), k2p.py, k2p.pz);
        return k1p.add(k2p);
      }
      /**
       * Constant time multiplication.
       * Uses wNAF method. Windowed method may be 10% faster,
       * but takes 2x longer to generate and consumes 2x memory.
       * Uses precomputes when available.
       * Uses endomorphism for Koblitz curves.
       * @param scalar by which the point would be multiplied
       * @returns New point
       */
      multiply(scalar) {
        const { endo, n: N } = CURVE;
        aInRange("scalar", scalar, _1n4, N);
        let point, fake;
        if (endo) {
          const { k1neg, k1, k2neg, k2 } = endo.splitScalar(scalar);
          let { p: k1p, f: f1p } = this.wNAF(k1);
          let { p: k2p, f: f2p } = this.wNAF(k2);
          k1p = wnaf.constTimeNegate(k1neg, k1p);
          k2p = wnaf.constTimeNegate(k2neg, k2p);
          k2p = new Point3(Fp2.mul(k2p.px, endo.beta), k2p.py, k2p.pz);
          point = k1p.add(k2p);
          fake = f1p.add(f2p);
        } else {
          const { p, f: f2 } = this.wNAF(scalar);
          point = p;
          fake = f2;
        }
        return Point3.normalizeZ([point, fake])[0];
      }
      /**
       * Efficiently calculate `aP + bQ`. Unsafe, can expose private key, if used incorrectly.
       * Not using Strauss-Shamir trick: precomputation tables are faster.
       * The trick could be useful if both P and Q are not G (not in our case).
       * @returns non-zero affine point
       */
      multiplyAndAddUnsafe(Q, a, b) {
        const G = Point3.BASE;
        const mul = (P, a2) => a2 === _0n4 || a2 === _1n4 || !P.equals(G) ? P.multiplyUnsafe(a2) : P.multiply(a2);
        const sum = mul(this, a).add(mul(Q, b));
        return sum.is0() ? void 0 : sum;
      }
      // Converts Projective point to affine (x, y) coordinates.
      // Can accept precomputed Z^-1 - for example, from invertBatch.
      // (x, y, z) ∋ (x=x/z, y=y/z)
      toAffine(iz) {
        return toAffineMemo(this, iz);
      }
      isTorsionFree() {
        const { h: cofactor, isTorsionFree } = CURVE;
        if (cofactor === _1n4)
          return true;
        if (isTorsionFree)
          return isTorsionFree(Point3, this);
        throw new Error("isTorsionFree() has not been declared for the elliptic curve");
      }
      clearCofactor() {
        const { h: cofactor, clearCofactor } = CURVE;
        if (cofactor === _1n4)
          return this;
        if (clearCofactor)
          return clearCofactor(Point3, this);
        return this.multiplyUnsafe(CURVE.h);
      }
      toRawBytes(isCompressed = true) {
        abool("isCompressed", isCompressed);
        this.assertValidity();
        return toBytes2(Point3, this, isCompressed);
      }
      toHex(isCompressed = true) {
        abool("isCompressed", isCompressed);
        return bytesToHex2(this.toRawBytes(isCompressed));
      }
    }
    Point3.BASE = new Point3(CURVE.Gx, CURVE.Gy, Fp2.ONE);
    Point3.ZERO = new Point3(Fp2.ZERO, Fp2.ONE, Fp2.ZERO);
    const _bits = CURVE.nBitLength;
    const wnaf = wNAF(Point3, CURVE.endo ? Math.ceil(_bits / 2) : _bits);
    return {
      CURVE,
      ProjectivePoint: Point3,
      normPrivateKeyToScalar,
      weierstrassEquation,
      isWithinCurveOrder
    };
  }
  function validateOpts(curve) {
    const opts = validateBasic(curve);
    validateObject(opts, {
      hash: "hash",
      hmac: "function",
      randomBytes: "function"
    }, {
      bits2int: "function",
      bits2int_modN: "function",
      lowS: "boolean"
    });
    return Object.freeze({ lowS: true, ...opts });
  }
  function weierstrass(curveDef) {
    const CURVE = validateOpts(curveDef);
    const { Fp: Fp2, n: CURVE_ORDER } = CURVE;
    const compressedLen = Fp2.BYTES + 1;
    const uncompressedLen = 2 * Fp2.BYTES + 1;
    function modN(a) {
      return mod(a, CURVE_ORDER);
    }
    function invN(a) {
      return invert(a, CURVE_ORDER);
    }
    const { ProjectivePoint: Point3, normPrivateKeyToScalar, weierstrassEquation, isWithinCurveOrder } = weierstrassPoints({
      ...CURVE,
      toBytes(_c, point, isCompressed) {
        const a = point.toAffine();
        const x = Fp2.toBytes(a.x);
        const cat = concatBytes2;
        abool("isCompressed", isCompressed);
        if (isCompressed) {
          return cat(Uint8Array.from([point.hasEvenY() ? 2 : 3]), x);
        } else {
          return cat(Uint8Array.from([4]), x, Fp2.toBytes(a.y));
        }
      },
      fromBytes(bytes2) {
        const len = bytes2.length;
        const head = bytes2[0];
        const tail = bytes2.subarray(1);
        if (len === compressedLen && (head === 2 || head === 3)) {
          const x = bytesToNumberBE(tail);
          if (!inRange(x, _1n4, Fp2.ORDER))
            throw new Error("Point is not on curve");
          const y2 = weierstrassEquation(x);
          let y;
          try {
            y = Fp2.sqrt(y2);
          } catch (sqrtError) {
            const suffix = sqrtError instanceof Error ? ": " + sqrtError.message : "";
            throw new Error("Point is not on curve" + suffix);
          }
          const isYOdd = (y & _1n4) === _1n4;
          const isHeadOdd = (head & 1) === 1;
          if (isHeadOdd !== isYOdd)
            y = Fp2.neg(y);
          return { x, y };
        } else if (len === uncompressedLen && head === 4) {
          const x = Fp2.fromBytes(tail.subarray(0, Fp2.BYTES));
          const y = Fp2.fromBytes(tail.subarray(Fp2.BYTES, 2 * Fp2.BYTES));
          return { x, y };
        } else {
          throw new Error(`Point of length ${len} was invalid. Expected ${compressedLen} compressed bytes or ${uncompressedLen} uncompressed bytes`);
        }
      }
    });
    const numToNByteStr = (num) => bytesToHex2(numberToBytesBE(num, CURVE.nByteLength));
    function isBiggerThanHalfOrder(number2) {
      const HALF = CURVE_ORDER >> _1n4;
      return number2 > HALF;
    }
    function normalizeS(s) {
      return isBiggerThanHalfOrder(s) ? modN(-s) : s;
    }
    const slcNum = (b, from, to) => bytesToNumberBE(b.slice(from, to));
    class Signature {
      constructor(r, s, recovery) {
        this.r = r;
        this.s = s;
        this.recovery = recovery;
        this.assertValidity();
      }
      // pair (bytes of r, bytes of s)
      static fromCompact(hex) {
        const l = CURVE.nByteLength;
        hex = ensureBytes("compactSignature", hex, l * 2);
        return new Signature(slcNum(hex, 0, l), slcNum(hex, l, 2 * l));
      }
      // DER encoded ECDSA signature
      // https://bitcoin.stackexchange.com/questions/57644/what-are-the-parts-of-a-bitcoin-transaction-input-script
      static fromDER(hex) {
        const { r, s } = DER.toSig(ensureBytes("DER", hex));
        return new Signature(r, s);
      }
      assertValidity() {
        aInRange("r", this.r, _1n4, CURVE_ORDER);
        aInRange("s", this.s, _1n4, CURVE_ORDER);
      }
      addRecoveryBit(recovery) {
        return new Signature(this.r, this.s, recovery);
      }
      recoverPublicKey(msgHash) {
        const { r, s, recovery: rec } = this;
        const h = bits2int_modN(ensureBytes("msgHash", msgHash));
        if (rec == null || ![0, 1, 2, 3].includes(rec))
          throw new Error("recovery id invalid");
        const radj = rec === 2 || rec === 3 ? r + CURVE.n : r;
        if (radj >= Fp2.ORDER)
          throw new Error("recovery id 2 or 3 invalid");
        const prefix = (rec & 1) === 0 ? "02" : "03";
        const R = Point3.fromHex(prefix + numToNByteStr(radj));
        const ir = invN(radj);
        const u1 = modN(-h * ir);
        const u2 = modN(s * ir);
        const Q = Point3.BASE.multiplyAndAddUnsafe(R, u1, u2);
        if (!Q)
          throw new Error("point at infinify");
        Q.assertValidity();
        return Q;
      }
      // Signatures should be low-s, to prevent malleability.
      hasHighS() {
        return isBiggerThanHalfOrder(this.s);
      }
      normalizeS() {
        return this.hasHighS() ? new Signature(this.r, modN(-this.s), this.recovery) : this;
      }
      // DER-encoded
      toDERRawBytes() {
        return hexToBytes2(this.toDERHex());
      }
      toDERHex() {
        return DER.hexFromSig({ r: this.r, s: this.s });
      }
      // padded bytes of r, then padded bytes of s
      toCompactRawBytes() {
        return hexToBytes2(this.toCompactHex());
      }
      toCompactHex() {
        return numToNByteStr(this.r) + numToNByteStr(this.s);
      }
    }
    const utils2 = {
      isValidPrivateKey(privateKey) {
        try {
          normPrivateKeyToScalar(privateKey);
          return true;
        } catch (error) {
          return false;
        }
      },
      normPrivateKeyToScalar,
      /**
       * Produces cryptographically secure private key from random of size
       * (groupLen + ceil(groupLen / 2)) with modulo bias being negligible.
       */
      randomPrivateKey: () => {
        const length = getMinHashLength(CURVE.n);
        return mapHashToField(CURVE.randomBytes(length), CURVE.n);
      },
      /**
       * Creates precompute table for an arbitrary EC point. Makes point "cached".
       * Allows to massively speed-up `point.multiply(scalar)`.
       * @returns cached point
       * @example
       * const fast = utils.precompute(8, ProjectivePoint.fromHex(someonesPubKey));
       * fast.multiply(privKey); // much faster ECDH now
       */
      precompute(windowSize = 8, point = Point3.BASE) {
        point._setWindowSize(windowSize);
        point.multiply(BigInt(3));
        return point;
      }
    };
    function getPublicKey(privateKey, isCompressed = true) {
      return Point3.fromPrivateKey(privateKey).toRawBytes(isCompressed);
    }
    function isProbPub(item) {
      const arr = isBytes3(item);
      const str = typeof item === "string";
      const len = (arr || str) && item.length;
      if (arr)
        return len === compressedLen || len === uncompressedLen;
      if (str)
        return len === 2 * compressedLen || len === 2 * uncompressedLen;
      if (item instanceof Point3)
        return true;
      return false;
    }
    function getSharedSecret(privateA, publicB, isCompressed = true) {
      if (isProbPub(privateA))
        throw new Error("first arg must be private key");
      if (!isProbPub(publicB))
        throw new Error("second arg must be public key");
      const b = Point3.fromHex(publicB);
      return b.multiply(normPrivateKeyToScalar(privateA)).toRawBytes(isCompressed);
    }
    const bits2int = CURVE.bits2int || function(bytes2) {
      const num = bytesToNumberBE(bytes2);
      const delta = bytes2.length * 8 - CURVE.nBitLength;
      return delta > 0 ? num >> BigInt(delta) : num;
    };
    const bits2int_modN = CURVE.bits2int_modN || function(bytes2) {
      return modN(bits2int(bytes2));
    };
    const ORDER_MASK = bitMask(CURVE.nBitLength);
    function int2octets(num) {
      aInRange(`num < 2^${CURVE.nBitLength}`, num, _0n4, ORDER_MASK);
      return numberToBytesBE(num, CURVE.nByteLength);
    }
    function prepSig(msgHash, privateKey, opts = defaultSigOpts) {
      if (["recovered", "canonical"].some((k) => k in opts))
        throw new Error("sign() legacy options not supported");
      const { hash: hash2, randomBytes: randomBytes2 } = CURVE;
      let { lowS, prehash, extraEntropy: ent } = opts;
      if (lowS == null)
        lowS = true;
      msgHash = ensureBytes("msgHash", msgHash);
      validateSigVerOpts(opts);
      if (prehash)
        msgHash = ensureBytes("prehashed msgHash", hash2(msgHash));
      const h1int = bits2int_modN(msgHash);
      const d = normPrivateKeyToScalar(privateKey);
      const seedArgs = [int2octets(d), int2octets(h1int)];
      if (ent != null && ent !== false) {
        const e = ent === true ? randomBytes2(Fp2.BYTES) : ent;
        seedArgs.push(ensureBytes("extraEntropy", e));
      }
      const seed = concatBytes2(...seedArgs);
      const m = h1int;
      function k2sig(kBytes) {
        const k = bits2int(kBytes);
        if (!isWithinCurveOrder(k))
          return;
        const ik = invN(k);
        const q = Point3.BASE.multiply(k).toAffine();
        const r = modN(q.x);
        if (r === _0n4)
          return;
        const s = modN(ik * modN(m + r * d));
        if (s === _0n4)
          return;
        let recovery = (q.x === r ? 0 : 2) | Number(q.y & _1n4);
        let normS = s;
        if (lowS && isBiggerThanHalfOrder(s)) {
          normS = normalizeS(s);
          recovery ^= 1;
        }
        return new Signature(r, normS, recovery);
      }
      return { seed, k2sig };
    }
    const defaultSigOpts = { lowS: CURVE.lowS, prehash: false };
    const defaultVerOpts = { lowS: CURVE.lowS, prehash: false };
    function sign(msgHash, privKey, opts = defaultSigOpts) {
      const { seed, k2sig } = prepSig(msgHash, privKey, opts);
      const C = CURVE;
      const drbg = createHmacDrbg(C.hash.outputLen, C.nByteLength, C.hmac);
      return drbg(seed, k2sig);
    }
    Point3.BASE._setWindowSize(8);
    function verify(signature, msgHash, publicKey, opts = defaultVerOpts) {
      const sg = signature;
      msgHash = ensureBytes("msgHash", msgHash);
      publicKey = ensureBytes("publicKey", publicKey);
      if ("strict" in opts)
        throw new Error("options.strict was renamed to lowS");
      validateSigVerOpts(opts);
      const { lowS, prehash } = opts;
      let _sig = void 0;
      let P;
      try {
        if (typeof sg === "string" || isBytes3(sg)) {
          try {
            _sig = Signature.fromDER(sg);
          } catch (derError) {
            if (!(derError instanceof DER.Err))
              throw derError;
            _sig = Signature.fromCompact(sg);
          }
        } else if (typeof sg === "object" && typeof sg.r === "bigint" && typeof sg.s === "bigint") {
          const { r: r2, s: s2 } = sg;
          _sig = new Signature(r2, s2);
        } else {
          throw new Error("PARSE");
        }
        P = Point3.fromHex(publicKey);
      } catch (error) {
        if (error.message === "PARSE")
          throw new Error(`signature must be Signature instance, Uint8Array or hex string`);
        return false;
      }
      if (lowS && _sig.hasHighS())
        return false;
      if (prehash)
        msgHash = CURVE.hash(msgHash);
      const { r, s } = _sig;
      const h = bits2int_modN(msgHash);
      const is = invN(s);
      const u1 = modN(h * is);
      const u2 = modN(r * is);
      const R = Point3.BASE.multiplyAndAddUnsafe(P, u1, u2)?.toAffine();
      if (!R)
        return false;
      const v = modN(R.x);
      return v === r;
    }
    return {
      CURVE,
      getPublicKey,
      getSharedSecret,
      sign,
      verify,
      ProjectivePoint: Point3,
      Signature,
      utils: utils2
    };
  }

  // node_modules/@noble/curves/esm/_shortw_utils.js
  function getHash(hash2) {
    return {
      hash: hash2,
      hmac: (key, ...msgs) => hmac(hash2, key, concatBytes(...msgs)),
      randomBytes
    };
  }
  function createCurve(curveDef, defHash) {
    const create = (hash2) => weierstrass({ ...curveDef, ...getHash(hash2) });
    return Object.freeze({ ...create(defHash), create });
  }

  // node_modules/@noble/curves/esm/secp256k1.js
  var secp256k1P = BigInt("0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f");
  var secp256k1N = BigInt("0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141");
  var _1n5 = BigInt(1);
  var _2n4 = BigInt(2);
  var divNearest = (a, b) => (a + b / _2n4) / b;
  function sqrtMod(y) {
    const P = secp256k1P;
    const _3n3 = BigInt(3), _6n = BigInt(6), _11n = BigInt(11), _22n = BigInt(22);
    const _23n = BigInt(23), _44n = BigInt(44), _88n = BigInt(88);
    const b2 = y * y * y % P;
    const b3 = b2 * b2 * y % P;
    const b6 = pow2(b3, _3n3, P) * b3 % P;
    const b9 = pow2(b6, _3n3, P) * b3 % P;
    const b11 = pow2(b9, _2n4, P) * b2 % P;
    const b22 = pow2(b11, _11n, P) * b11 % P;
    const b44 = pow2(b22, _22n, P) * b22 % P;
    const b88 = pow2(b44, _44n, P) * b44 % P;
    const b176 = pow2(b88, _88n, P) * b88 % P;
    const b220 = pow2(b176, _44n, P) * b44 % P;
    const b223 = pow2(b220, _3n3, P) * b3 % P;
    const t1 = pow2(b223, _23n, P) * b22 % P;
    const t2 = pow2(t1, _6n, P) * b2 % P;
    const root = pow2(t2, _2n4, P);
    if (!Fp.eql(Fp.sqr(root), y))
      throw new Error("Cannot find square root");
    return root;
  }
  var Fp = Field(secp256k1P, void 0, void 0, { sqrt: sqrtMod });
  var secp256k1 = createCurve({
    a: BigInt(0),
    // equation params: a, b
    b: BigInt(7),
    // Seem to be rigid: bitcointalk.org/index.php?topic=289795.msg3183975#msg3183975
    Fp,
    // Field's prime: 2n**256n - 2n**32n - 2n**9n - 2n**8n - 2n**7n - 2n**6n - 2n**4n - 1n
    n: secp256k1N,
    // Curve order, total count of valid points in the field
    // Base point (x, y) aka generator point
    Gx: BigInt("55066263022277343669578718895168534326250603453777594175500187360389116729240"),
    Gy: BigInt("32670510020758816978083085130507043184471273380659243275938904335757337482424"),
    h: BigInt(1),
    // Cofactor
    lowS: true,
    // Allow only low-S signatures by default in sign() and verify()
    /**
     * secp256k1 belongs to Koblitz curves: it has efficiently computable endomorphism.
     * Endomorphism uses 2x less RAM, speeds up precomputation by 2x and ECDH / key recovery by 20%.
     * For precomputed wNAF it trades off 1/2 init time & 1/3 ram for 20% perf hit.
     * Explanation: https://gist.github.com/paulmillr/eb670806793e84df628a7c434a873066
     */
    endo: {
      beta: BigInt("0x7ae96a2b657c07106e64479eac3434e99cf0497512f58995c1396c28719501ee"),
      splitScalar: (k) => {
        const n = secp256k1N;
        const a1 = BigInt("0x3086d221a7d46bcde86c90e49284eb15");
        const b1 = -_1n5 * BigInt("0xe4437ed6010e88286f547fa90abfe4c3");
        const a2 = BigInt("0x114ca50f7a8e2f3f657c1108d9d44cfd8");
        const b2 = a1;
        const POW_2_128 = BigInt("0x100000000000000000000000000000000");
        const c1 = divNearest(b2 * k, n);
        const c2 = divNearest(-b1 * k, n);
        let k1 = mod(k - c1 * a1 - c2 * a2, n);
        let k2 = mod(-c1 * b1 - c2 * b2, n);
        const k1neg = k1 > POW_2_128;
        const k2neg = k2 > POW_2_128;
        if (k1neg)
          k1 = n - k1;
        if (k2neg)
          k2 = n - k2;
        if (k1 > POW_2_128 || k2 > POW_2_128) {
          throw new Error("splitScalar: Endomorphism failed, k=" + k);
        }
        return { k1neg, k1, k2neg, k2 };
      }
    }
  }, sha256);
  var _0n5 = BigInt(0);
  var Point = secp256k1.ProjectivePoint;

  // node_modules/@scure/bip32/lib/esm/index.js
  var Point2 = secp256k1.ProjectivePoint;
  var base58check = createBase58check(sha256);
  function bytesToNumber(bytes2) {
    return BigInt(`0x${bytesToHex(bytes2)}`);
  }
  function numberToBytes(num) {
    return hexToBytes(num.toString(16).padStart(64, "0"));
  }
  var MASTER_SECRET = utf8ToBytes("Bitcoin seed");
  var BITCOIN_VERSIONS = { private: 76066276, public: 76067358 };
  var HARDENED_OFFSET = 2147483648;
  var hash160 = (data) => ripemd160(sha256(data));
  var fromU32 = (data) => createView(data).getUint32(0, false);
  var toU32 = (n) => {
    if (!Number.isSafeInteger(n) || n < 0 || n > 2 ** 32 - 1) {
      throw new Error(`Invalid number=${n}. Should be from 0 to 2 ** 32 - 1`);
    }
    const buf = new Uint8Array(4);
    createView(buf).setUint32(0, n, false);
    return buf;
  };
  var HDKey = class _HDKey {
    get fingerprint() {
      if (!this.pubHash) {
        throw new Error("No publicKey set!");
      }
      return fromU32(this.pubHash);
    }
    get identifier() {
      return this.pubHash;
    }
    get pubKeyHash() {
      return this.pubHash;
    }
    get privateKey() {
      return this.privKeyBytes || null;
    }
    get publicKey() {
      return this.pubKey || null;
    }
    get privateExtendedKey() {
      const priv = this.privateKey;
      if (!priv) {
        throw new Error("No private key");
      }
      return base58check.encode(this.serialize(this.versions.private, concatBytes(new Uint8Array([0]), priv)));
    }
    get publicExtendedKey() {
      if (!this.pubKey) {
        throw new Error("No public key");
      }
      return base58check.encode(this.serialize(this.versions.public, this.pubKey));
    }
    static fromMasterSeed(seed, versions = BITCOIN_VERSIONS) {
      bytes(seed);
      if (8 * seed.length < 128 || 8 * seed.length > 512) {
        throw new Error(`HDKey: wrong seed length=${seed.length}. Should be between 128 and 512 bits; 256 bits is advised)`);
      }
      const I = hmac(sha512, MASTER_SECRET, seed);
      return new _HDKey({
        versions,
        chainCode: I.slice(32),
        privateKey: I.slice(0, 32)
      });
    }
    static fromExtendedKey(base58key, versions = BITCOIN_VERSIONS) {
      const keyBuffer = base58check.decode(base58key);
      const keyView = createView(keyBuffer);
      const version = keyView.getUint32(0, false);
      const opt = {
        versions,
        depth: keyBuffer[4],
        parentFingerprint: keyView.getUint32(5, false),
        index: keyView.getUint32(9, false),
        chainCode: keyBuffer.slice(13, 45)
      };
      const key = keyBuffer.slice(45);
      const isPriv = key[0] === 0;
      if (version !== versions[isPriv ? "private" : "public"]) {
        throw new Error("Version mismatch");
      }
      if (isPriv) {
        return new _HDKey({ ...opt, privateKey: key.slice(1) });
      } else {
        return new _HDKey({ ...opt, publicKey: key });
      }
    }
    static fromJSON(json) {
      return _HDKey.fromExtendedKey(json.xpriv);
    }
    constructor(opt) {
      this.depth = 0;
      this.index = 0;
      this.chainCode = null;
      this.parentFingerprint = 0;
      if (!opt || typeof opt !== "object") {
        throw new Error("HDKey.constructor must not be called directly");
      }
      this.versions = opt.versions || BITCOIN_VERSIONS;
      this.depth = opt.depth || 0;
      this.chainCode = opt.chainCode || null;
      this.index = opt.index || 0;
      this.parentFingerprint = opt.parentFingerprint || 0;
      if (!this.depth) {
        if (this.parentFingerprint || this.index) {
          throw new Error("HDKey: zero depth with non-zero index/parent fingerprint");
        }
      }
      if (opt.publicKey && opt.privateKey) {
        throw new Error("HDKey: publicKey and privateKey at same time.");
      }
      if (opt.privateKey) {
        if (!secp256k1.utils.isValidPrivateKey(opt.privateKey)) {
          throw new Error("Invalid private key");
        }
        this.privKey = typeof opt.privateKey === "bigint" ? opt.privateKey : bytesToNumber(opt.privateKey);
        this.privKeyBytes = numberToBytes(this.privKey);
        this.pubKey = secp256k1.getPublicKey(opt.privateKey, true);
      } else if (opt.publicKey) {
        this.pubKey = Point2.fromHex(opt.publicKey).toRawBytes(true);
      } else {
        throw new Error("HDKey: no public or private key provided");
      }
      this.pubHash = hash160(this.pubKey);
    }
    derive(path) {
      if (!/^[mM]'?/.test(path)) {
        throw new Error('Path must start with "m" or "M"');
      }
      if (/^[mM]'?$/.test(path)) {
        return this;
      }
      const parts = path.replace(/^[mM]'?\//, "").split("/");
      let child = this;
      for (const c of parts) {
        const m = /^(\d+)('?)$/.exec(c);
        const m1 = m && m[1];
        if (!m || m.length !== 3 || typeof m1 !== "string") {
          throw new Error(`Invalid child index: ${c}`);
        }
        let idx = +m1;
        if (!Number.isSafeInteger(idx) || idx >= HARDENED_OFFSET) {
          throw new Error("Invalid index");
        }
        if (m[2] === "'") {
          idx += HARDENED_OFFSET;
        }
        child = child.deriveChild(idx);
      }
      return child;
    }
    deriveChild(index) {
      if (!this.pubKey || !this.chainCode) {
        throw new Error("No publicKey or chainCode set");
      }
      let data = toU32(index);
      if (index >= HARDENED_OFFSET) {
        const priv = this.privateKey;
        if (!priv) {
          throw new Error("Could not derive hardened child key");
        }
        data = concatBytes(new Uint8Array([0]), priv, data);
      } else {
        data = concatBytes(this.pubKey, data);
      }
      const I = hmac(sha512, this.chainCode, data);
      const childTweak = bytesToNumber(I.slice(0, 32));
      const chainCode = I.slice(32);
      if (!secp256k1.utils.isValidPrivateKey(childTweak)) {
        throw new Error("Tweak bigger than curve order");
      }
      const opt = {
        versions: this.versions,
        chainCode,
        depth: this.depth + 1,
        parentFingerprint: this.fingerprint,
        index
      };
      try {
        if (this.privateKey) {
          const added = mod(this.privKey + childTweak, secp256k1.CURVE.n);
          if (!secp256k1.utils.isValidPrivateKey(added)) {
            throw new Error("The tweak was out of range or the resulted private key is invalid");
          }
          opt.privateKey = added;
        } else {
          const added = Point2.fromHex(this.pubKey).add(Point2.fromPrivateKey(childTweak));
          if (added.equals(Point2.ZERO)) {
            throw new Error("The tweak was equal to negative P, which made the result key invalid");
          }
          opt.publicKey = added.toRawBytes(true);
        }
        return new _HDKey(opt);
      } catch (err) {
        return this.deriveChild(index + 1);
      }
    }
    sign(hash2) {
      if (!this.privateKey) {
        throw new Error("No privateKey set!");
      }
      bytes(hash2, 32);
      return secp256k1.sign(hash2, this.privKey).toCompactRawBytes();
    }
    verify(hash2, signature) {
      bytes(hash2, 32);
      bytes(signature, 64);
      if (!this.publicKey) {
        throw new Error("No publicKey set!");
      }
      let sig;
      try {
        sig = secp256k1.Signature.fromCompact(signature);
      } catch (error) {
        return false;
      }
      return secp256k1.verify(sig, hash2, this.publicKey);
    }
    wipePrivateData() {
      this.privKey = void 0;
      if (this.privKeyBytes) {
        this.privKeyBytes.fill(0);
        this.privKeyBytes = void 0;
      }
      return this;
    }
    toJSON() {
      return {
        xpriv: this.privateExtendedKey,
        xpub: this.publicExtendedKey
      };
    }
    serialize(version, key) {
      if (!this.chainCode) {
        throw new Error("No chainCode set");
      }
      bytes(key, 33);
      return concatBytes(toU32(version), new Uint8Array([this.depth]), toU32(this.parentFingerprint), toU32(this.index), this.chainCode, key);
    }
  };

  // index.ts
  function generateBip39Mnemonic() {
    return generateMnemonic(wordlist, 128);
  }
  function isMnemonicValid(mnemonic) {
    return validateMnemonic(mnemonic, wordlist);
  }
  function getBinarySeed(mnemonic) {
    return mnemonicToSeedSync(mnemonic);
  }
  function getPrivateAndPublicKeys(binarySeed) {
    const hdKey = HDKey.fromMasterSeed(binarySeed);
    return {
      "privateKey": hdKey.privateKey,
      "publicKey": hdKey.publicKey
    };
  }
  return __toCommonJS(mnemonic_exports);
})();
/*! Bundled license information:

@noble/hashes/esm/utils.js:
  (*! noble-hashes - MIT License (c) 2022 Paul Miller (paulmillr.com) *)

@scure/base/lib/esm/index.js:
  (*! scure-base - MIT License (c) 2022 Paul Miller (paulmillr.com) *)

@scure/bip39/esm/index.js:
  (*! scure-bip39 - MIT License (c) 2022 Patricio Palladino, Paul Miller (paulmillr.com) *)

@noble/curves/esm/abstract/utils.js:
  (*! noble-curves - MIT License (c) 2022 Paul Miller (paulmillr.com) *)

@noble/curves/esm/abstract/modular.js:
  (*! noble-curves - MIT License (c) 2022 Paul Miller (paulmillr.com) *)

@noble/curves/esm/abstract/curve.js:
  (*! noble-curves - MIT License (c) 2022 Paul Miller (paulmillr.com) *)

@noble/curves/esm/abstract/weierstrass.js:
  (*! noble-curves - MIT License (c) 2022 Paul Miller (paulmillr.com) *)

@noble/curves/esm/_shortw_utils.js:
  (*! noble-curves - MIT License (c) 2022 Paul Miller (paulmillr.com) *)

@noble/curves/esm/secp256k1.js:
  (*! noble-curves - MIT License (c) 2022 Paul Miller (paulmillr.com) *)

@scure/bip32/lib/esm/index.js:
  (*! scure-bip32 - MIT License (c) 2022 Patricio Palladino, Paul Miller (paulmillr.com) *)
*/

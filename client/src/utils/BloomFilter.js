
function fnv1Hash64(input, seed) {
    // TODO: Fetch this prime from either API or env
    const prime = BigInt(0x100000001B3); // FNV-1 64-bit prime
    let hash = BigInt(seed);

    for (let i = 0; i < input.length; i++) {
        hash ^= BigInt(input.charCodeAt(i));
        hash *= prime;

        // *NOTE: Keep hash within 64-bit range by applying modulo
        hash = hash & BigInt("0xFFFFFFFFFFFFFFFF");
    }

    return hash;
}

class BloomFilter {
    constructor(bitArray, bitArraySize, hashFunctionCount) {
        if (bitArraySize <= 0) throw new Error('Size of the bit array must be greater than 0.');
        if (hashFunctionCount <= 0) throw new Error('Hash function count must be greater than 0.');

        this.bitArray = bitArray;
        this.bitArraySize = bitArraySize;
        this.hashFunctionCount = hashFunctionCount;
    }

    addPassword(value) {
        if (typeof value !== 'string') {
            throw new TypeError('Value must be a string.');
        }

        for (let i = 0; i < this.hashFunctionCount; i++) {
            const hash = this.getHash(value, i);
            this.bitArray[hash] = 1;
        }
    }

    checkPassword(value) {
        if (typeof value !== 'string') {
            throw new TypeError('Value must be a string.');
        }

        for (let i = 0; i < this.hashFunctionCount; i++) {
            const hash = this.getHash(value, i);
            if (this.bitArray[hash] === 0) {
                return false;
            }
        }
        return true;
    }

    getHash(value, seed) {
        const hash = fnv1Hash64(value, seed);
        return Number(hash % BigInt(this.bitArraySize));
    }
}

export default BloomFilter;

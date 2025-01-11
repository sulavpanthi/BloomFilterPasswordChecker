import { x86 } from 'murmurhash3js';

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
        const hash = x86.hash32(value + seed.toString());
        return Math.abs(hash % this.bitArraySize);
    }
}

export default BloomFilter;

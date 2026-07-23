declare module 'murmurhash3js-revisited' {
  export const x86: {
    hash32(key: string, seed?: number): number;
    hash128(key: string, seed?: number): string;
  };
  export const x64: {
    hash128(key: string, seed?: number): string;
  };
}

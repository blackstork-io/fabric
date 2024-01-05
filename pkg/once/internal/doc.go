// 50% faster and 50% smaller version of sync.OnceFunc, sync.OnceValue and sync.OnceValues
// Doesn't keep closures after they have been executed, allowing GC to collect captured variables
package once

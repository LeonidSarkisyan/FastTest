

export function IsNotEmpty(...strings) {
    for (let s of strings) {
        if (s.length === 0) {
            return false
        }
    }

    return true
}

const pad = (n: number) => String(n).padStart(2, "0");

export function toLocalInput(dt: Date) {
    return `${dt.getFullYear()}-${pad(dt.getMonth() + 1)}-${pad(dt.getDate())}T${pad(dt.getHours())}:${pad(dt.getMinutes())}`;
}

export function toLocalYMD(dt: Date) {
    return `${dt.getFullYear()}-${pad(dt.getMonth() + 1)}-${pad(dt.getDate())}`;
}

export function isoToLocalInput(iso: string) {
    return toLocalInput(new Date(iso));
}
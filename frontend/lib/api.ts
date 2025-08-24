export const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? "/api";

async function handle(resOrPromise: Response | Promise<Response>) {
    const res = await resOrPromise; // ← ждём fetch
    if (!res.ok) {
        let text = "";
        try { text = await res.text(); } catch {}
        throw new Error(text || res.statusText);
    }
    if (res.status === 204) return null;
    const ct = res.headers.get("content-type") ?? "";
    if (ct.includes("application/json")) return res.json();
    return res.text();
}

export const api = {
    listReleases: (params: { date?: string; status?: string; duty?: string } = {}) =>
        handle(fetch(`${API_BASE}/releases${qs(params)}`, { cache: "no-store" })),

    getRelease: (id: number) =>
        handle(fetch(`${API_BASE}/releases/${id}`, { cache: "no-store" })),

    createRelease: (payload: any) =>
        handle(fetch(`${API_BASE}/releases`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
        })),

    updateRelease: (id: number, payload: any) =>
        handle(fetch(`${API_BASE}/releases/${id}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
        })),

    deleteRelease: (id: number) =>
        handle(fetch(`${API_BASE}/releases/${id}`, { method: "DELETE" })),

    listComments: (id: number) =>
        handle(fetch(`${API_BASE}/releases/${id}/comments`, { cache: "no-store" })),

    addComment: (id: number, payload: any) =>
        handle(fetch(`${API_BASE}/releases/${id}/comments`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
        })),

    updateComment: (id: number, commentId: number, payload: any) =>
        handle(fetch(`${API_BASE}/releases/${id}/comments/${commentId}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
        })),

    deleteComment: (id: number, commentId: number) =>
        handle(fetch(`${API_BASE}/releases/${id}/comments/${commentId}`, {
            method: "DELETE",
        })),
};

function qs(params: Record<string, any>) {
    const entries = Object.entries(params).filter(([, v]) => v !== undefined && v !== "");
    if (!entries.length) return "";
    const query = new URLSearchParams(entries as any).toString();
    return `?${query}`;
}

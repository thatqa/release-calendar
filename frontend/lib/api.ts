// api (inline mock) — drop in place of your current file

// small latency to mimic network
const delay = <T,>(data: T, ms = 120) => new Promise<T>(r => setTimeout(() => r(data), ms));

// helpers
const pad = (n: number) => (n < 10 ? `0${n}` : `${n}`);
const ymdLocal = (d: Date) => `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
const atTimeLocalISO = (day: Date, hh: number, mm: number) => {
    const d = new Date(day);
    d.setHours(hh, mm, 0, 0);
    return d.toISOString(); // full ISO for UI formatting
};
const addDays = (base: Date, n: number) => {
    const d = new Date(base);
    d.setDate(d.getDate() + n);
    return d;
};

// dates
const today = new Date();
const Dm2 = addDays(today, -2);
const Dm1 = addDays(today, -1);
const D0  = today;
const Dp1 = addDays(today,  1);
const Dp2 = addDays(today,  2);

// types (aligned with your page)
type Link = { id: number; name: string; url: string };
type Release = {
    id: number;
    title: string;
    date: string; // ISO
    status: "planned" | "success" | "failed";
    notes?: string;
    dutyUsers: string[];
    links: Link[];
    createdAt: string;
    updatedAt: string;
};
type Comment = { id: number; author: string; message: string; createdAt: string };
type Markers = Record<string, Array<"planned" | "success" | "failed">>;

// saturated demo data
const releases: Release[] = [
    // two days ago — success + failed
    {
        id: 101,
        title: "Accounts service refactor",
        date: atTimeLocalISO(Dm2, 10, 15),
        status: "success",
        notes:
            "Scope: handler refactor, dependency bump, retry/backoff tuning.\n" +
            "Observability: error rate ↓, P95 latency -8%.",
        dutyUsers: ["alice", "qa-team"],
        links: [
            { id: 1, name: "Jira ACC-1201", url: "https://example.com/j/ACC-1201" },
            { id: 2, name: "Dashboard", url: "https://example.com/grafana/acc" },
        ],
        createdAt: atTimeLocalISO(Dm2, 9, 0),
        updatedAt: atTimeLocalISO(Dm2, 11, 0),
    },
    {
        id: 102,
        title: "Kubernetes upgrade 1.30",
        date: atTimeLocalISO(Dm2, 14, 40),
        status: "failed",
        notes:
            "CNI incompatibility on tainted nodes. Performed rapid rollback (<15m).\n" +
            "Action: validate CNI matrix before next attempt.",
        dutyUsers: ["bob", "sre"],
        links: [
            { id: 1, name: "Runbook", url: "https://example.com/runbooks/k8s-upgrade" },
            { id: 2, name: "Jira OPS-903", url: "https://example.com/j/OPS-903" },
        ],
        createdAt: atTimeLocalISO(Dm2, 13, 50),
        updatedAt: atTimeLocalISO(Dm2, 15, 10),
    },

    // yesterday — success + failed
    {
        id: 201,
        title: "UI revamp batch #2",
        date: atTimeLocalISO(Dm1, 11, 0),
        status: "success",
        notes:
            "Promoted A/B test; accessibility improvements; Lighthouse scores up.\n" +
            "CTR +5.7%; LCP -120ms on 4G.",
        dutyUsers: ["carol"],
        links: [
            { id: 1, name: "Design spec", url: "https://example.com/design/ui-revamp" },
            { id: 2, name: "Jira WEB-221", url: "https://example.com/j/WEB-221" },
        ],
        createdAt: atTimeLocalISO(Dm1, 9, 30),
        updatedAt: atTimeLocalISO(Dm1, 12, 10),
    },
    {
        id: 202,
        title: "Hotfix: rate-limiter",
        date: atTimeLocalISO(Dm1, 16, 20),
        status: "failed",
        notes:
            "Overly strict limits caused 429 spike. Reverted via feature flag.\n" +
            "Plan: segmented rules by endpoint/IP range.",
        dutyUsers: ["dave", "qa-team"],
        links: [{ id: 1, name: "Jira OPS-921", url: "https://example.com/j/OPS-921" }],
        createdAt: atTimeLocalISO(Dm1, 15, 40),
        updatedAt: atTimeLocalISO(Dm1, 16, 45),
    },

    // today — success + failed + planned
    {
        id: 301,
        title: "Search relevance tuning",
        date: atTimeLocalISO(D0, 9, 30),
        status: "success",
        notes: "Adjusted weights, updated synonyms. Search CTR +2%.",
        dutyUsers: ["erin"],
        links: [{ id: 1, name: "Jira SRCH-510", url: "https://example.com/j/SRCH-510" }],
        createdAt: atTimeLocalISO(D0, 8, 50),
        updatedAt: atTimeLocalISO(D0, 9, 45),
    },
    {
        id: 302,
        title: "Android SDK hotfix",
        date: atTimeLocalISO(D0, 13, 0),
        status: "failed",
        notes:
            "Startup crash on ~20% devices post-SDK bump. Rolled back to 2.9.4.\n" +
            "Investigating OEM-specific WebView init.",
        dutyUsers: ["frank", "mobile"],
        links: [{ id: 1, name: "Crash dashboard", url: "https://example.com/grafana/mobile" }],
        createdAt: atTimeLocalISO(D0, 12, 5),
        updatedAt: atTimeLocalISO(D0, 13, 20),
    },
    {
        id: 303,
        title: "Database maintenance",
        date: atTimeLocalISO(D0, 18, 0),
        status: "planned",
        notes: "Index rebuild on user tables. Expected read-only window up to 10 minutes.",
        dutyUsers: ["gina", "dba"],
        links: [{ id: 1, name: "Change request DB-450", url: "https://example.com/j/DB-450" }],
        createdAt: atTimeLocalISO(D0, 10, 0),
        updatedAt: atTimeLocalISO(D0, 10, 0),
    },

    // tomorrow — planned
    {
        id: 401,
        title: "Landing page redesign",
        date: atTimeLocalISO(Dp1, 11, 0),
        status: "planned",
        notes: "New hero section, optimized images, CLS target < 0.05.",
        dutyUsers: ["helen", "web"],
        links: [{ id: 1, name: "Spec", url: "https://example.com/specs/landing" }],
        createdAt: atTimeLocalISO(D0, 17, 0),
        updatedAt: atTimeLocalISO(D0, 17, 0),
    },

    // day after tomorrow — planned
    {
        id: 501,
        title: "Billing service split (Phase 1)",
        date: atTimeLocalISO(Dp2, 15, 30),
        status: "planned",
        notes: "Extract invoicing module, dark launch behind FF, 5% ramp.",
        dutyUsers: ["ivan", "backend"],
        links: [{ id: 1, name: "Architectural plan", url: "https://example.com/arch/billing" }],
        createdAt: atTimeLocalISO(D0, 17, 30),
        updatedAt: atTimeLocalISO(D0, 17, 30),
    },
];

const comments: Record<number, Comment[]> = {
    101: [
        { id: 1, author: "Alice", createdAt: new Date().toISOString(),
            message: "Rollout clean. P95 improved; error budget intact. Dashboards updated." },
        { id: 2, author: "QA", createdAt: new Date().toISOString(),
            message: "Full regression passed on staging; spot checks OK in prod." },
    ],
    102: [
        { id: 1, author: "Bob", createdAt: new Date().toISOString(),
            message: "CNI plugin crash on tainted nodes. Rolled back; re-validate matrix." },
    ],
    201: [
        { id: 1, author: "Carol", createdAt: new Date().toISOString(),
            message: "A/B promoted; monitoring CLS in RUM for the next 24h." },
    ],
    202: [
        { id: 1, author: "Dave", createdAt: new Date().toISOString(),
            message: "429 spike post-hotfix. Reverted via FF; segmented limits next." },
    ],
    301: [
        { id: 1, author: "Erin", createdAt: new Date().toISOString(),
            message: "Long-tail queries improved. Typo tolerance tuning next sprint." },
    ],
    302: [
        { id: 1, author: "Frank", createdAt: new Date().toISOString(),
            message: "Crash reproduced on 2 OEMs. SDK + WebView interaction under review." },
    ],
    303: [
        { id: 1, author: "DBA", createdAt: new Date().toISOString(),
            message: "Backups verified. Maintenance window announced." },
    ],
    401: [{ id: 1, author: "Helen", createdAt: new Date().toISOString(),
        message: "Design approved. Watch CLS and LCP closely post-deploy." }],
    501: [{ id: 1, author: "Ivan", createdAt: new Date().toISOString(),
        message: "Phase 1 will route 5% traffic behind a feature flag." }],
};

// internal helpers
const ymdFromISO = (iso: string) => ymdLocal(new Date(iso));

// API (inline)
export const api = {
    listReleases: (params: { date?: string; status?: string; duty?: string } = {}) => {
        let list = releases.slice();

        if (params.date) {
            list = list.filter(r => ymdFromISO(r.date) === params.date);
        }
        if (params.status) {
            list = list.filter(r => r.status === params.status);
        }
        if (params.duty) {
            const needle = params.duty.toLowerCase();
            list = list.filter(r => r.dutyUsers.some(u => u.toLowerCase().includes(needle)));
        }
        // sort by time asc
        list.sort((a, b) => a.date.localeCompare(b.date));
        return delay(list);
    },

    getRelease: (id: number) => delay(releases.find(r => r.id === id) || null),

    createRelease: () => Promise.reject("Demo only"),
    updateRelease: () => Promise.reject("Demo only"),
    deleteRelease: () => Promise.reject("Demo only"),

    listComments: (id: number) => delay(comments[id] || []),

    addComment: () => Promise.reject("Demo only"),
    updateComment: () => Promise.reject("Demo only"),
    deleteComment: () => Promise.reject("Demo only"),

    // returns Markers: { "YYYY-MM-DD": ["failed","success","planned"] }
    releaseDays: (from: string, to: string) => {
        const fromS = from || "0000-01-01";
        const toS   = to   || "9999-12-31";
        const map: Markers = {};

        for (const r of releases) {
            const d = ymdFromISO(r.date);
            if (d >= fromS && d <= toS) {
                map[d] ||= [];
                if (!map[d].includes(r.status)) map[d].push(r.status);
            }
        }
        // keep arrays in a stable order (failed, success, planned) for consistent dot colors
        const order = { failed: 0, success: 1, planned: 2 } as const;
        for (const k of Object.keys(map)) {
            map[k].sort((a, b) => order[a] - order[b]);
        }
        return delay(map);
    },

    getSummary: (id: number) => {
        const r = releases.find(x => x.id === id);
        if (!r) return delay(null);
        const summary =
            `Status: ${r.status.toUpperCase()} @ ${new Date(r.date).toLocaleString()}\n` +
            `Title: ${r.title}\n` +
            `Duty: ${r.dutyUsers.join(", ")}\n` +
            `Links: ${r.links.map(l => `${l.name}`).join(", ")}\n` +
            `Notes: ${r.notes ?? ""}`;
        return delay({ id: r.id, summary });
    },
};

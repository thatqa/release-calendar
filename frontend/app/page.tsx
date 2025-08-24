"use client";

import { useEffect, useMemo, useState } from "react";
import { DayPicker } from "react-day-picker";
import "react-day-picker/dist/style.css";
import { api } from "../lib/api";
import { Button } from "../components/ui/button";
import { ReleaseForm, ReleasePayload } from "../components/ReleaseForm";
import { Comments } from "../components/Comments";

type Link = { id:number; name:string; url:string };
type Release = {
  id:number; title:string; date:string; status:"planned"|"success"|"failed";
  notes?:string; dutyUsers:string[]; links: Link[]; createdAt:string; updatedAt:string;
};

export default function HomePage() {
  const [selectedDay, setSelectedDay] = useState<Date>(new Date());
  const [releases, setReleases] = useState<Release[]>([]);
  const [active, setActive] = useState<Release|null>(null);
  const [mode, setMode] = useState<"view"|"create"|"edit">("view");
  const [statusFilter, setStatusFilter] = useState<string>("");
  const [dutyFilter, setDutyFilter] = useState<string>("");

  const dayParam = useMemo(()=> selectedDay.toISOString().slice(0,10), [selectedDay]);

  const load = async ()=> {
    const data = await api.listReleases({ date: dayParam, status: statusFilter||undefined, duty: dutyFilter||undefined}) as Release[];
    setReleases(data);
    if (data.length) setActive(data[0]); else setActive(null);
  };
  useEffect(()=>{ load(); }, [dayParam, statusFilter, dutyFilter]);

  const createRelease = async (payload: ReleasePayload)=> {
    await api.createRelease(payload);
    setMode("view"); await load();
  };
  const updateRelease = async (payload: ReleasePayload)=> {
    if (!active) return;
    await api.updateRelease(active.id, payload);
    setMode("view"); await load();
  };
  const deleteRelease = async ()=> {
    if (!active) return;
    if (!confirm("Delete this release?")) return;
    await api.deleteRelease(active.id);
    setActive(null); await load();
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-12 gap-6">
      <div className="md:col-span-5 card">
        <div className="flex items-center justify-between mb-3">
          <div className="text-lg font-semibold">Calendar</div>
          <Button onClick={()=>{ setMode("create"); setActive(null); }}>New release</Button>
        </div>
        <DayPicker
          mode="single"
          selected={selectedDay}
          onSelect={(d)=> d && setSelectedDay(d)}
          weekStartsOn={1}
        />
        <div className="mt-4">
          <div className="text-sm mb-1">Filters</div>
          <div className="grid grid-cols-2 gap-2">
            <select className="select" value={statusFilter} onChange={e=>setStatusFilter(e.target.value)}>
              <option value="">All statuses</option>
              <option value="planned">planned</option>
              <option value="success">success</option>
              <option value="failed">failed</option>
            </select>
            <input className="input" placeholder="Duty user" value={dutyFilter} onChange={e=>setDutyFilter(e.target.value)} />
          </div>
        </div>

        <div className="mt-4">
          <div className="text-sm mb-2">Releases on {dayParam}</div>
          <div className="space-y-2">
            {releases.map(r => (
              <div key={r.id} className={`border rounded-xl p-3 cursor-pointer ${active?.id===r.id?'bg-neutral-50':''}`} onClick={()=>{ setActive(r); setMode("view"); }}>
                <div className="flex items-center justify-between">
                  <div className="font-medium">{r.title}</div>
                  <span className="badge border-neutral-300">{r.status}</span>
                </div>
                <div className="text-sm text-neutral-500">{new Date(r.date).toLocaleString()}</div>
              </div>
            ))}
            {!releases.length && <div className="text-neutral-500">No releases for this date</div>}
          </div>
        </div>
      </div>

      <div className="md:col-span-7">
        {!active && mode!=="create" && (
          <div className="card">
            <div className="text-lg font-semibold">No release selected</div>
            <div className="mt-2 text-neutral-600">Pick a date and create a release.</div>
          </div>
        )}

        {mode==="create" && (
          <div className="card">
            <div className="text-lg font-semibold mb-3">Create release</div>
            <ReleaseForm
              initial={{ date: selectedDay.toISOString().slice(0,16) }}
              onSubmit={createRelease}
              onCancel={()=>setMode("view")}
            />
          </div>
        )}

        {active && mode==="view" && (
          <div className="space-y-4">
            <div className="card">
              <div className="flex items-center justify-between">
                <div className="text-lg font-semibold">{active.title}</div>
                <div className="flex gap-2">
                  <Button onClick={()=>setMode("edit")}>Edit</Button>
                  <Button onClick={deleteRelease}>Delete</Button>
                </div>
              </div>
              <div className="text-sm text-neutral-500">{new Date(active.date).toLocaleString()} · <span className="badge border-neutral-300">{active.status}</span></div>
              {active.notes && <div className="mt-2 whitespace-pre-wrap">{active.notes}</div>}
              <div className="mt-2">
                <div className="font-medium">Duty users</div>
                <div className="flex flex-wrap gap-2 mt-1">
                  {active.dutyUsers?.map(u => <span key={u} className="badge border-neutral-300">{u}</span>)}
                </div>
              </div>
              <div className="mt-2">
                <div className="font-medium">Links</div>
                <ul className="list-disc ml-6 mt-1">
                  {active.links?.map(l => <li key={l.id}><a className="text-blue-600 underline" href={l.url} target="_blank" rel="noreferrer">{l.name}</a></li>)}
                </ul>
              </div>
            </div>
            <Comments releaseId={active.id} />
          </div>
        )}

        {active && mode==="edit" && (
          <div className="card">
            <div className="text-lg font-semibold mb-3">Edit release</div>
            <ReleaseForm
              initial={{
                title: active.title,
                date: active.date.slice(0,16),
                status: active.status,
                notes: active.notes,
                dutyUsers: active.dutyUsers,
                links: active.links,
              }}
              onSubmit={updateRelease}
              onCancel={()=>setMode("view")}
            />
          </div>
        )}
      </div>
    </div>
  );
}

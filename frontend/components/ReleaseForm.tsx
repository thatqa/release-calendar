"use client";

import { useState } from "react";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Textarea } from "./ui/textarea";
import { Select } from "./ui/select";

export type LinkItem = { id?: number; name: string; url: string };
export type ReleasePayload = {
  title: string;
  date: string;
  status: "planned"|"success"|"failed";
  notes?: string;
  dutyUsers: string[];
  links: LinkItem[];
};

export function ReleaseForm({ initial, onSubmit, onCancel } : {
  initial?: Partial<ReleasePayload>;
  onSubmit: (payload: ReleasePayload) => Promise<void>|void;
  onCancel?: () => void;
}) {
  const [title, setTitle] = useState(initial?.title || "");
  const [date, setDate] = useState(initial?.date || new Date().toISOString().slice(0,16));
  const [status, setStatus] = useState<"planned"|"success"|"failed">((initial?.status as any) || "planned");
  const [notes, setNotes] = useState(initial?.notes || "");
  const [duty, setDuty] = useState((initial?.dutyUsers as string[] | undefined) || []);
  const [links, setLinks] = useState<LinkItem[]>(initial?.links || []);
  const [dutyInput, setDutyInput] = useState("");

  const addDuty = () => {
    const v = dutyInput.trim();
    if (!v) return;
    if (!duty.includes(v)) setDuty([...duty, v]);
    setDutyInput("");
  };
  const removeDuty = (name: string) => setDuty(duty.filter(d => d !== name));

  const addLink = () => setLinks([...links, { name: "", url: "" }]);
  const updateLink = (idx: number, patch: Partial<LinkItem>) => {
    const copy = links.slice();
    copy[idx] = { ...copy[idx], ...patch };
    setLinks(copy);
  };
  const removeLink = (idx: number) => setLinks(links.filter((_,i)=> i!==idx));

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="card">
          <div className="text-lg font-semibold mb-3">Main</div>
          <div className="space-y-3">
            <div>
              <div className="text-sm mb-1">Title</div>
              <Input value={title} onChange={e=>setTitle(e.target.value)} placeholder="Release 1.2.3" />
            </div>
            <div>
              <div className="text-sm mb-1">Date & time</div>
              <Input type="datetime-local" value={date} onChange={e=>setDate(e.target.value)} />
            </div>
            <div>
              <div className="text-sm mb-1">Status</div>
              <Select value={status} onChange={e=>setStatus(e.target.value as any)}>
                <option value="planned">planned</option>
                <option value="success">success</option>
                <option value="failed">failed</option>
              </Select>
            </div>
            <div>
              <div className="text-sm mb-1">Notes</div>
              <Textarea rows={5} value={notes} onChange={e=>setNotes(e.target.value)} placeholder="Notes for release..." />
            </div>
          </div>
        </div>

        <div className="card">
          <div className="text-lg font-semibold mb-3">Duty users</div>
          <div className="flex gap-2">
            <Input value={dutyInput} onChange={(e)=>setDutyInput(e.target.value)} placeholder="username" />
            <Button type="button" onClick={addDuty}>Add</Button>
          </div>
          <div className="mt-2 flex flex-wrap gap-2">
            {duty.map(u => (
              <span key={u} className="badge border-neutral-300">
                {u}
                <button className="ml-2 text-neutral-500" onClick={()=>removeDuty(u)}>✕</button>
              </span>
            ))}
          </div>
        </div>
      </div>

      <div className="card">
        <div className="text-lg font-semibold mb-3">Links</div>
        <div className="space-y-2">
          {links.map((l, idx) => (
            <div key={idx} className="grid grid-cols-1 md:grid-cols-12 gap-2 items-center">
              <Input className="md:col-span-3" placeholder="Name" value={l.name} onChange={e=>updateLink(idx,{name:e.target.value})} />
              <Input className="md:col-span-8" placeholder="https://..." value={l.url} onChange={e=>updateLink(idx,{url:e.target.value})} />
              <Button className="md:col-span-1" type="button" onClick={()=>removeLink(idx)}>Del</Button>
            </div>
          ))}
          <Button type="button" onClick={addLink}>Add link</Button>
        </div>
      </div>

      <div className="flex gap-2">
        <Button onClick={()=>onSubmit({
          title, date: new Date(date).toISOString(), status, notes, dutyUsers: duty, links
        })}>Save</Button>
        {onCancel ? <Button onClick={onCancel}>Cancel</Button> : null}
      </div>
    </div>
  );
}

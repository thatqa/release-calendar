"use client";

import { useEffect, useState } from "react";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Textarea } from "./ui/textarea";
import { api } from "../lib/api";

type Comment = { id:number; author:string; message:string; createdAt:string; updatedAt:string };

export function Comments({ releaseId }: { releaseId: number }) {
  const [items, setItems] = useState<Comment[]>([]);
  const [author,setAuthor] = useState("");
  const [message,setMessage] = useState("");

  const load = async ()=> {
    const data = await api.listComments(releaseId) as Comment[];
    setItems(data);
  };
  useEffect(()=>{ load(); }, [releaseId]);

  const add = async ()=> {
    if (!author.trim() || !message.trim()) return;
    await api.addComment(releaseId, {author, message});
    setAuthor(""); setMessage("");
    await load();
  };
  const update = async (c: Comment)=> {
    const na = prompt("Author", c.author) ?? c.author;
    const nm = prompt("Message", c.message) ?? c.message;
    await api.updateComment(releaseId, c.id, { author: na, message: nm });
    await load();
  };
  const del = async (c: Comment)=> {
    if (!confirm("Delete comment?")) return;
    await api.deleteComment(releaseId, c.id);
    await load();
  };

  return (
    <div className="card mt-4">
      <div className="text-lg font-semibold mb-3">Comments</div>
      <div className="space-y-2">
        {items.map(c => (
          <div key={c.id} className="border rounded-xl p-3 flex items-start justify-between">
            <div>
              <div className="text-sm text-neutral-500">{new Date(c.createdAt).toLocaleString()}</div>
              <div className="font-medium">{c.author}</div>
              <div>{c.message}</div>
            </div>
            <div className="flex gap-2">
              <Button onClick={()=>update(c)}>Edit</Button>
              <Button onClick={()=>del(c)}>Delete</Button>
            </div>
          </div>
        ))}
      </div>

      <div className="mt-3 grid grid-cols-1 md:grid-cols-12 gap-2">
        <Input className="md:col-span-3" placeholder="Author" value={author} onChange={e=>setAuthor(e.target.value)} />
        <Textarea className="md:col-span-8" placeholder="Message" value={message} onChange={e=>setMessage(e.target.value)} />
        <Button className="md:col-span-1" onClick={add}>Add</Button>
      </div>
    </div>
  );
}

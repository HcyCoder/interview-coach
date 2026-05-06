'use client';

import * as React from 'react';

import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { cn } from '@/lib/utils';

type Message = {
  id: number;
  role: 'system' | 'assistant' | 'user';
  text: string;
};

const initialMessages: Message[] = [
  {
    id: 1,
    role: 'assistant',
    text: 'Upload a resume and JD to generate your first interview plan.',
  },
  {
    id: 2,
    role: 'system',
    text: 'Backend API, AI service, and database connectivity will appear here once wired.',
  },
];

function MessageBubble({ message }: { message: Message }) {
  const isUser = message.role === 'user';

  return (
    <div
      className={cn(
        'flex max-w-[85%] flex-col gap-1 rounded-2xl border px-4 py-3 text-sm leading-6',
        isUser
          ? 'ml-auto border-primary/40 bg-primary/10 text-foreground'
          : 'border-border bg-card/70 text-foreground',
      )}
    >
      <span className="text-[11px] uppercase tracking-[0.24em] text-muted-foreground">
        {message.role}
      </span>
      <p>{message.text}</p>
    </div>
  );
}

/**
 * ChatShell renders the landing workspace for the interview coach.
 *
 * Inputs: none.
 * Outputs: a responsive chat-oriented workspace with placeholder state.
 */
export function ChatShell() {
  const [messages, setMessages] = React.useState<Message[]>(initialMessages);
  const [draft, setDraft] = React.useState('');
  const [subject, setSubject] = React.useState(
    'Platform Engineer interview prep with a focus on system design and backend fundamentals.',
  );

  function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const content = draft.trim();
    if (!content) {
      return;
    }

    setMessages((current) => [
      ...current,
      {
        id: Date.now(),
        role: 'user',
        text: content,
      },
    ]);
    setDraft('');
  }

  return (
    <div className="min-h-screen bg-[radial-gradient(circle_at_top,_rgba(34,211,238,0.16),_transparent_32%),linear-gradient(180deg,_#020617_0%,_#0f172a_100%)]">
      <main className="mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-6 px-4 py-6 sm:px-6 lg:px-8">
        <header className="flex flex-col gap-4 rounded-3xl border border-border/70 bg-card/70 px-6 py-5 shadow-sm backdrop-blur md:flex-row md:items-center md:justify-between">
          <div className="space-y-1">
            <p className="text-xs uppercase tracking-[0.28em] text-cyan-300">
              Interview Coach
            </p>
            <h1 className="text-2xl font-semibold tracking-tight text-foreground sm:text-3xl">
              AI interview practice workspace
            </h1>
            <p className="max-w-2xl text-sm text-muted-foreground">
              A placeholder front door for resume parsing, interview planning,
              live questioning, and post-session review.
            </p>
          </div>
          <div className="flex items-center gap-2">
            <span className="rounded-full border border-emerald-400/30 bg-emerald-400/10 px-3 py-1 text-xs font-medium text-emerald-300">
              Frontend online
            </span>
            <span className="rounded-full border border-slate-700 bg-slate-900/60 px-3 py-1 text-xs font-medium text-slate-300">
              App Router
            </span>
          </div>
        </header>

        <section className="grid flex-1 gap-6 lg:grid-cols-[minmax(0,1.65fr)_320px]">
          <Card className="flex flex-col overflow-hidden border-border/70 bg-card/75">
            <CardHeader>
              <CardTitle>Conversation</CardTitle>
              <CardDescription>
                Draft interview prompts, review candidate context, and simulate
                the AI interviewer.
              </CardDescription>
            </CardHeader>
            <CardContent className="flex flex-1 flex-col gap-4">
              <div className="flex min-h-[26rem] flex-1 flex-col gap-3 rounded-2xl border border-border/70 bg-background/50 p-4">
                <div className="flex items-center justify-between border-b border-border/60 pb-3">
                  <div>
                    <p className="text-sm font-medium text-foreground">
                      Interview session placeholder
                    </p>
                    <p className="text-xs text-muted-foreground">
                      No backend session attached yet
                    </p>
                  </div>
                  <span className="rounded-full bg-muted px-3 py-1 text-xs text-muted-foreground">
                    Draft
                  </span>
                </div>

                <div className="flex flex-1 flex-col gap-3 overflow-y-auto pr-1">
                  {messages.map((message) => (
                    <MessageBubble key={message.id} message={message} />
                  ))}
                </div>
              </div>

              <form className="space-y-3" onSubmit={handleSubmit}>
                <label className="sr-only" htmlFor="chat-draft">
                  Chat message
                </label>
                <Textarea
                  id="chat-draft"
                  value={draft}
                  onChange={(event) => setDraft(event.target.value)}
                  placeholder="Type a follow-up question or note for the interviewer..."
                />
                <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                  <Input
                    aria-label="Interview subject"
                    value={subject}
                    onChange={(event) => setSubject(event.target.value)}
                    placeholder="Interview subject"
                  />
                  <Button type="submit" className="sm:w-auto">
                    Send message
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>

          <Card className="border-border/70 bg-card/60">
            <CardHeader>
              <CardTitle>Session state</CardTitle>
              <CardDescription>
                Skeleton data for the rest of the product flows.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-5">
              <div className="space-y-2 rounded-2xl border border-border/70 bg-background/40 p-4">
                <p className="text-xs uppercase tracking-[0.24em] text-muted-foreground">
                  Current focus
                </p>
                <p className="text-sm text-foreground">{subject}</p>
              </div>

              <div className="space-y-3">
                {[
                  'Resume upload route placeholder',
                  'Interview plan generation placeholder',
                  'Report generation placeholder',
                ].map((item) => (
                  <div
                    key={item}
                    className="flex items-start gap-3 rounded-2xl border border-border/70 bg-background/40 p-4"
                  >
                    <div className="mt-1 h-2 w-2 rounded-full bg-cyan-300" />
                    <p className="text-sm text-muted-foreground">{item}</p>
                  </div>
                ))}
              </div>

              <div className="rounded-2xl border border-dashed border-border/80 bg-background/30 p-4">
                <p className="text-xs uppercase tracking-[0.24em] text-muted-foreground">
                  Notes
                </p>
                <p className="mt-2 text-sm leading-6 text-muted-foreground">
                  This shell is intentionally minimal: enough to verify the app
                  router, component wiring, and local interaction loop before
                  the backend APIs are connected.
                </p>
              </div>
            </CardContent>
          </Card>
        </section>
      </main>
    </div>
  );
}

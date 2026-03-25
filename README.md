# Single-Page Developer Portfolio

A server-rendered portfolio site built with Go.

## Run locally

```bash
go run ./cmd/web
```

The app listens on `http://localhost:8080` by default.

## Environment variables

- `PORT`: HTTP port for the Go server. Render sets this automatically.
- `SUBMISSIONS_PATH`: where contact form submissions are stored locally. Defaults to `data/submissions.jsonl`.
- `DATABASE_URL`: optional Postgres connection string. If set, submissions are stored in Postgres instead of a local file.

## Deploy to Render

This app is ready for Render as a Go web service.

Suggested settings:

- Environment: `Go`
- Build Command: `go build -o app ./cmd/web`
- Start Command: `./app`

Or deploy with the included `render.yaml`.

## Contact form storage

By default, the contact form writes to `data/submissions.jsonl`.

For production, set `DATABASE_URL` to your Supabase Postgres connection string. If `DATABASE_URL` is present, the app stores submissions in Postgres instead of the local file.

Create this table in Supabase:

```sql
create table if not exists public.contact_submissions (
  id bigint generated always as identity primary key,
  name text not null,
  email text not null,
  message text not null,
  created_at timestamptz not null default now()
);
```

Recommended deployment setup:

- App hosting: Render
- Database: Supabase Postgres
- Render environment variable: `DATABASE_URL`

## Verify

```bash
go test ./...
go build ./...
```

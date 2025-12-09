Purpose
This file gives AI coding agents the essential, discoverable knowledge to work productively on the Ephemeral File Sharer backend.

**Overview**
- **Project Type:** Minimal Go HTTP backend that accepts multipart uploads and serves files from local storage.
- **Entry point:** `main.go` — contains the HTTP routes and the `uploadHandler` implementation.
- **Module:** `github.com/prawns-89/ephemeral-sharer` (see `go.mod`, `go 1.22.2`).

**How to run (dev)**
- **Run locally:** `go run main.go` (server binds to `:8080`).
- **Build binary:** `go build -o ephemeral-sharer` then `./ephemeral-sharer`.
- **Quick check:** Open `http://localhost:8080/` to view the upload form.

**Key files & locations**
- `main.go`: single-file server. Look for the upload form handler (`/`), upload handler (`/upload`) and static file server (`/files/`).
- `uploads/` directory: local storage target for uploaded files. Server writes files directly here using `handler.Filename`.
- `README.md`: lightweight project notes and planned features (auto-deletion, unique links, DB choices).

**Routing & data flow (what to expect)**
- Upload form posts to `/upload` with multipart/form-data where the file field is named `myFile`.
- `uploadHandler` does:
  - `r.ParseMultipartForm(10 << 20)` — upload size capped at about 10MB.
  - Reads file via `r.FormFile("myFile")` and writes to `./uploads/<filename>`.
  - Returns an HTML response with a download link at `/files/<filename>`.
- Static files are served with `http.FileServer(http.Dir("./uploads"))` and `http.StripPrefix("/files/", fs)`.

**Project-specific conventions & gotchas**
- Filenames: current code writes `handler.Filename` directly — agents should note this is unsanitized and platform-dependent; any change must preserve compatibility with the download URL `/files/<filename>`.
- Storage: uses the local filesystem only; there is no database code in the repo yet — planned DB choices appear in `README.md` but are not implemented.
- Logs: the app uses `fmt.Println` for simple stdout logging. Don't assume a logger exists.
- Port & binding: server listens on `:8080`; changing the port is a global behavior change.

**When adding features (practical pointers)**
- Unique download links: implement a mapping layer (in-memory or DB) and keep `/files/` handler backward-compatible. Example: map token -> filename and add a separate handler that redirects to `/files/<filename>`.
- Auto-deletion: add a background goroutine or a cron job process that removes files from `./uploads` and removes any mapping records. Place cleanup code in a new package (e.g., `cleanup`) and call it from `main.go`.
- Resumable uploads (TUS): treat as a new subsystem — keep it separate from the existing `uploadHandler` to avoid regression.

**Debugging & developer workflow notes**
- Reproduce locally with `go run main.go` and use the HTML form or `curl -F "myFile=@path" http://localhost:8080/upload`.
- Check `./uploads` for saved files and use `http://localhost:8080/files/<filename>` to download.
- There are no tests or CI config in the repo; add unit tests around new packages rather than modifying `main.go` directly.

**Security & compatibility reminders (observed, actionable)**
- Filename sanitization is not implemented — when modifying upload storage, validate or normalize `handler.Filename` and consider using generated IDs for on-disk filenames.
- Be careful when refactoring static serving: the current public URL contract is `/files/<filename>` and downstream clients may rely on it.

**What not to change without explicit tests/consent**
- The public upload route `/upload`, the field name `myFile`, and the `/files/` URL shape.
- The default storage path `./uploads` unless you provide a migration path (and tests/downtime plan).

If anything above is unclear or you want more detail (examples of adding a token-based link layer, safe filename helper, or a sample cleanup goroutine), tell me which area to expand and I'll update this file.

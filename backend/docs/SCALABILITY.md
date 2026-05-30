# BuildingMesh | Scalability & High-Load Strategy

To support massive LiDAR and Wi-Fi CSI data streams from multiple organizations, BuildingMesh employs a specialized high-performance architecture.

## 1. Storage Optimization (Offloading)
- **Binary Data Policy:** Massive binary blobs (Point Clouds, CSI matrices) are NEVER stored in the relational database (SQLite/PostgreSQL).
- **S3-First Strategy:** All heavy assets are streamed directly to object storage (S3 or Local Persistent Volume) via the `S3Service`.
- **Relational Integrity:** SQLite only maintains lightweight metadata and pointers (`storage://`) to the actual blobs, ensuring indexing remains fast and the DB size predictable.

## 2. Peak Load Damping (RabbitMQ)
- **Async Ingestion:** Data uploads return `202 Accepted` immediately after the raw file is safely persisted to storage.
- **Worker Pools:** Intensive computations (e.g., decimation, 3DGS rendering, AI analysis) are handled by background workers consuming from RabbitMQ.
- **Traffic Shaping:** This prevents API request timeouts and ensures the system remains responsive under heavy ingestion bursts.

## 3. SQLite Efficiency (CGO-Free)
- **Driver:** Using `modernc.org/sqlite` ensures maximum portability and eliminates CGO overhead.
- **IO Patterns:** By removing binary blobs from the DB, we minimize Disk I/O contention, allowing SQLite to handle high-frequency relational queries effectively.

## 4. Multi-tenant Isolation
- **Resource Limits:** Hard limits on project counts and users prevent single-tenant noise from impacting the global platform performance.
- **Isolated Buffers:** Multipart form processing uses fixed memory buffers (32MB) to prevent OOM errors during concurrent large uploads.

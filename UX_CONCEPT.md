# BuildingMesh | Enterprise Hybrid Architecture Concept

BuildingMesh V1 is designed as a high-performance construction monitoring platform that leverages a hybrid backend architecture for scalability, data sovereignty, and rapid time-to-market.

## 1. The Hybrid Core
The system splits responsibility between two distinct layers:

### A. Proprietary Construction Engine (The "Brain")
- **Tech Stack:** Go 1.23, SQLite (Metadata), Redis (Cache), RabbitMQ (Task Queue).
- **Responsibility:**
  - Processing heavy LiDAR and 3D Gaussian Splatting data.
  - Analyzing Wi-Fi CSI signals from RuView sensors.
  - Managing chronological "Time-Machine" snapshots of wall structures.
  - Generating legal-grade audit reports (PDF).
- **Isolation:** This layer owns all Intellectual Property (IP) related to construction physics and sensor fusion.

### B. Twenty CRM Engine (The "Enterprise Metadata Head")
- **Tech Stack:** NestJS, PostgreSQL, GraphQL.
- **Responsibility:**
  - Multi-tenant Organization management.
  - Role-Based Access Control (RBAC) definitions.
  - Project metadata (addresses, participants, schedules).
  - Workflow automation (email alerts, status transitions).
- **Integration:** Connected via a secure Proxy Client using GraphQL.

## 2. User Journey: "The Git-like Commit"
1. **Field Scan:** An engineer scans a wall using the iPhone LiDAR + RuView sensor.
2. **Local Processing:** The phone processes the point cloud and sends it to our proprietary Go backend.
3. **Commit Creation:** The backend creates an `InspectionCommit`.
4. **Metadata Sync:** The Go backend tells Twenty CRM: *"Project B2 has a new verified layer."*
5. **Office Review:** Management sees the update in the Twenty web dashboard, while the engineer sees the "X-Ray" view in the AR app.

## 3. Data Sovereignty & Licensing
- **Licensing:** Twenty is licensed under AGPL-3.0. By using the **Proxy Pattern**, we keep our proprietary algorithms isolated from Twenty's source code, protecting our IP for future acquisition.
- **Flexibility:** The architecture allows swapping Twenty for Salesforce, SAP, or a custom ERP without touching the 3D/Sensing logic.

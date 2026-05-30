# BuildingMesh | Enterprise AR-Monitoring System

BuildingMesh is a high-fidelity SaaS platform for construction sites, designed to visualize internal wall structures, hidden installations, and historical construction phases.

## 🏗 Hybrid Enterprise Architecture
BuildingMesh uses a unique two-layer architecture to balance deep tech performance with corporate management needs:

1. **Proprietary Sensing Engine (Go 1.23)**: Handles LiDAR processing, 3D Gaussian Splatting, and RuView Wi-Fi CSI sensing. This is our core IP.
2. **Twenty CRM Engine (Headless Metadata)**: Provides Enterprise-grade multi-tenancy, workflow automation, and metadata management via an isolated AGPL-3.0 compliant proxy.

## 🚀 Local Deployment & Access
To launch the system on your local server or Intel N100 edge node:

1. **Quick Start**:
   ```bash
   chmod +x deploy_server.sh
   ./deploy_server.sh
   ```

2. **Access Points**:
   - **Engineering Console (Web UI)**: [http://localhost:8080](http://localhost:8080)
   - **Twenty CRM (Admin Panel)**: [http://localhost:3000](http://localhost:3000)
   - **API Documentation (Swagger)**: [http://localhost:8080/api/docs](http://localhost:8080/api/docs)
   - **System Health**: [http://localhost:8080/api/v1/health/live](http://localhost:8080/api/v1/health/live)

## 🛠 Key Features
- **AR X-Ray**: See "through" finished walls using historical LiDAR/CSI chronological layers.
- **Layer Time-Machine**: Slide through time to see the state of any wall during electrical, insulation, or finishing stages.
- **Git-like Integrity**: Construction phases are "committed" and, once approved by AI/Engineers, become Read-Only to provide a legally significant audit trail.
- **RuView Sensing**: Integrated support for ESP32-S3 sensors for non-destructive physical verification of installations.

## 📦 Deployment Topology
- **OS Target**: Debian 13 (Trixie)
- **Engine**: Docker Compose (PaaS-ready)
- **Infrastructure**: Redis (Cache), RabbitMQ (Peak Damping), SQLite (Metadata), PostgreSQL (CRM Metadata).

---
*Note: This system is architected for maximum IP protection, isolating proprietary sensing algorithms from open-source dependencies (AGPL-3.0).*

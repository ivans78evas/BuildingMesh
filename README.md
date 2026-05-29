# BuildingMesh | Enterprise AR-Monitoring System

BuildingMesh is a high-fidelity SaaS platform for construction sites, designed to visualize internal wall structures, hidden installations, and historical construction phases.

## 🏗 Hybrid Enterprise Architecture
BuildingMesh uses a unique two-layer architecture to balance deep tech performance with corporate management needs:

1. **Proprietary Sensing Engine (Go 1.23)**: Handles LiDAR processing, 3D Gaussian Splatting, and RuView Wi-Fi CSI sensing. This is our core IP.
2. **Twenty CRM Engine (Headless Metadata)**: Provides Enterprise-grade multi-tenancy, workflow automation, and metadata management via an isolated AGPL-3.0 compliant proxy.

## 🚀 Key Features
- **AR X-Ray**: See "through" finished walls using historical LiDAR/CSI chronological layers.
- **Layer Time-Machine**: Slide through time to see the state of any wall during electrical, insulation, or finishing stages.
- **Git-like Inspection**: Construct digital "commits" for hidden works. Approved states become Read-Only to ensure a legal audit trail.
- **RuView Sensing**: Integrated support for ESP32-S3 sensors for non-destructive physical verification of installations.

## 🛠 Tech Stack
- **Backend**: Go (Proprietary Brain), Twenty CRM (Enterprise Framework), Redis, RabbitMQ.
- **Storage**: SQLite (Local Metadata), PostgreSQL (CRM), S3 (Point Clouds/CSI logs).
- **Frontend**: Flutter (Mobile AR/GPU Compute), React (Web Engineering Console).
- **AI**: Google Cloud Vision & Gemini for BIM/Drawing alignment.

## 📦 Deployment
1. **Docker Compose**: Start the entire cluster including Twenty and our backend:
   ```bash
   docker-compose up -d
   ```
2. **Backend**:
   ```bash
   cd backend
   go run cmd/server/main.go
   ```
3. **Engineering Console**: Access the web dashboard at `http://localhost:8080`.

---
*Note: This system is architected for maximum IP protection, isolating proprietary sensing algorithms from open-source dependencies.*

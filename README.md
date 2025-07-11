# Tower Defense Data Processing System

This repository contains the complete source code for the "Tower Defense" real-time data processing simulation. The project is architected as a two-part system: a high-performance Go backend engine and a reactive Vue.js frontend interface. The entire system is orchestrated with Docker Compose for simple, one-command deployment.

---

## 📖 Project Overview

This application is a strategic game that simulates a real-world data processing pipeline.

-   **Backend (Go Engine):** Acts as the game server. It manages unit spawning, game state, weapon systems, scoring, and all core logic. It communicates with the frontend exclusively via WebSockets.
-   **Frontend (Vue UI):** Serves as a rich, real-time command center for a human operator. The user can select units, form strategic groups, and deploy them to the "Defense Tower" for processing, receiving instant visual feedback.

---

## 🛠️ Technology Stack

| Component | Technology                                                              | Purpose                                          |
| :-------- | :---------------------------------------------------------------------- | :----------------------------------------------- |
| **Backend** | [Go](https://go.dev/)                                                   | High-performance, concurrent game logic          |
| **Frontend**  | [Vue 3](https://vuejs.org/) + [TypeScript](https://www.typescriptlang.org/) | Reactive UI & type safety                        |
| **Styling**   | [Tailwind CSS](https://tailwindcss.com/)                                | Utility-first CSS framework for rapid UI development |
| **Real-time** | [WebSockets](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API)      | Low-latency, bidirectional communication       |
| **Orchestration** | [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)    | Containerization and one-command deployment    |

---

## 🚀 Getting Started

### Prerequisites

-   [Docker](https://docs.docker.com/get-docker/)
-   [Docker Compose](https://docs.docker.com/compose/install/)

### Directory Structure

This `README.md` should be in a parent directory that contains the two project folders, like so:

```
parent-directory/
├── README.md
├── docker-compose.yml
├── go-projects/
│   └── tower-defence-engine/
└── vue-projects/
    └── tower-defence-ui/
```

### One-Command Deployment

With Docker running, navigate to the parent directory containing the `docker-compose.yml` file and run a single command:

```bash
docker-compose up --build
```

-   This command will build the Docker images for both the backend and frontend services.
-   It will then start both containers.
-   The backend will be accessible on port `8080`.
-   The frontend will be served on port `5173`.

Once the build is complete and the containers are running, access the application in your web browser at:

**➡️ `http://localhost:5173`**

### Stopping the System

To stop all running containers, press `Ctrl+C` in the terminal where `docker-compose up` is running. To remove the containers completely, run:

```bash
docker-compose down
```

---

## ⚙️ Configuration

The game's balance and difficulty can be fully customized by editing a single file:

-   **File Path:** `go-projects/tower-defence-engine/config.json`

You can adjust parameters such as unit spawn rates, weapon firing speeds, unit hitpoints, weapon processing power, and scoring values. **Changes require a restart of the Docker containers to take effect.**

```bash
# To apply config changes, restart the containers
docker-compose up --build -d
```
*(The `-d` flag runs the containers in detached mode.)*

---

## 🎮 How to Play

1.  **Observe:** Units will automatically spawn on the main battlefield. Each has a timer (TTL).
2.  **Select:** Click on units to select them for an action.
3.  **Deploy:** Use the action buttons to deploy units.
    -   **Deploy Individuals:** Sends selected units to the fast, low-power weapon.
    -   **Create Squad Blueprint:** Saves the composition of selected units for later use.
    -   **Deploy Squad:** Deploys a group based on a saved blueprint, using available units.
4.  **Manage:** Keep an eye on the Defense Corridor. Units must be destroyed before they reach the border. A backlog in the weapon queues can lead to "Border Breaches."
5.  **Score:** Your goal is to achieve the highest possible score by efficiently processing units while minimizing Escapes and Breaches.

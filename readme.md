# Product Agent (Market Watcher)

Product Agent is an AI-powered tool written in Go that monitors app metrics for anomalies and uses Google's Gemini LLM to analyze user reviews and uncover the reasons behind those anomalies.

## Architecture & Workflow

1. **Anomaly Detection**: Periodically fetches anomalies from a Google BigQuery dataset (e.g., unexpected spikes in `CostPerInstall`).
2. **Review Fetching**: Scrapes the latest user reviews from the Apple App Store (for iOS apps) or Google Play Store (for Android apps) based on the anomalous App ID.
3. **Vector Analytics**: Embeds the fetched reviews using Google AI (`gemini-embedding-001`) and stores them in a local ChromaDB instance for intelligent semantic search.
4. **AI Generation**: Uses the Gemini LLM to interpret the context of the reviews against the specific metric anomaly, automatically generating a comprehensive "Agent Report" detailing potential root causes.

## Prerequisites

- **Go 1.21+**
- **Docker**: Required to run the local Chroma vector database instance.
- **Google Cloud / BigQuery**: A Service Account with access to the required BigQuery datasets.
- **Gemini API Key**: Access to Google's Gemini API for LLMs and embeddings.

## Installation

1. **Clone & Setup Environment**
   Create a `.env` file in the root directory using the `BIGQUERY_SERVICE_ACCOUNT`, `BIGQUERY_PROJECT_ID`, and `GEMINI_API_KEY`:

   ```env
   BIGQUERY_PROJECT_ID=your_gcp_project_id
   BIGQUERY_SERVICE_ACCOUNT={"type": "service_account", ...} # Include the full JSON string
   GEMINI_API_KEY=your_gemini_api_key
   ```
   *(Note: The private key string must use properly escaped `\n` characters.)*

2. **Start ChromaDB**
   Run the following Docker command to start a Chroma instance on port 8000:
   ```bash
   docker run -d -p 8000:8000 chromadb/chroma:0.4.24
   ```

3. **Run the Application**
   ```bash
   go run cmd/agent/main.go
   ```

# Credit Risk Scoring Service (Fintech)

A scoring router that receives applicant data, dispatches it to a Champion and Challenger model API, and returns a unified response while logging comparative performance.

The Gateway serves TCP, HTTP and gRPC endpoints that receives loan application data and calls two external credit scoring APIs for A/B testing.

The Credit Risk Scoring Gateway (CSG) is a distributed system designed to route financial loan applications through multiple scoring models (Champion and Challenger) and compare their performance in real-time.

The system aims to:
* Experimentally evaluate new ML models without disrupting production scoring.
* Collect, stream, and persist model performance metrics for analysis.
* Expose dashboards that visualize model accuracy, latency, and decision divergence.

## High-Level Purpose

Financial institutions use credit scoring models to assess the risk of lending to an applicant.

To deploy new models safely, they run Champion–Challenger experiments, sending a portion of real requests to new (“challenger”) models while keeping production (“champion”) models live.

This server simulates the decision routing gateway in this pipeline.

## System Objectives

| Objective                   | Description                                                                            |
| --------------------------- | -------------------------------------------------------------------------------------- |
| **A/B model comparison**    | Route requests probabilistically between Champion and Challenger models.               |
| **Metrics collection**      | Capture latency, score differences, and decision consistency for each scoring request. |
| **Streaming architecture**  | Publish metrics asynchronously to Kafka for scalability and decoupling.                |
| **Persistence layer**       | Store detailed metrics in a relational database for historical analysis.               |
| **Visualization dashboard** | Display near real-time comparisons of model performance metrics in Streamlit.          |

## Architecture

| Component                       | Type                                  | Description                                                                   |
| ------------------------------- | ------------------------------------- | ----------------------------------------------------------------------------- |
| **CreditScore Gateway**         | TCP, HTTP and gRPC Server             | Listens for client requests, orchestrates model calls, gathers metrics, and publishes `MetricEvent` to Kafka asynchronously before returning a response to the client. |
| **Champion API**                | REST Service                          | Stable, trusted scoring model (e.g., logistic regression).                                        |
| **Challenger API**              | REST Service                          | New, experimental scoring model (e.g., gradient boosting).                                        |
| **Kafka Cluster**               | Message broker                        | Decouples real-time event publishing from persistence. Topic: `credit_score_metrics`.             |
| **Metrics Persistence Service** | Kafka consumer for event metrics      | Consumes Kafka messages, validates, and inserts into a database. Provides an API for data access. |
| **Database**                    | PostgreSQL (or SQLite for simplicity) | Stores metrics for each scoring event (Champion/Challenger results, latency, timestamp).          |
| **Streamlit Dashboard**         | Frontend analytics UI                 | Displays real-time and historical metrics to compare Champion vs Challenger performance.          |


## High-Level Architecture

![High-Level Architecture](./docs/high-level-architecture.png)

### Components Descriptions

1. **Credit Score Gateway**
    - Handle TCP, HTTP and gRPC client requests concurrently.
    - Route requests:
        - Champion API (production model) — always called.
        - Challenger API (experimental model) — probabilistically called (e.g., 20% of traffic).
    - Measure and log metrics:
        - Response latency per model.
        - Score differences between Champion and Challenger.
    - Publish MetricEvent JSON messages to Kafka asynchronously via a buffered channel before returning a response to the client.
    - Output:
        - Response with the Champion's decision over TCP, HTTP or gRPC.
        - Kafka message with experiment metadata.
2. **Credit Score Services (Champion and Challenger Model APIs)**
    - Mocked APIs for local development running on different ports (8081 and 8082, for example).
    - Endpoint: `/score`
    - Request:
        ```json
        {
            "applicant_id": "123",
            "income": 60000,
            "loan_amount": 25000,
            "credit_history": 4
        }
        ```
    - Response:
        ```json
        {
            "score": 0.72,
            "decision": "approve"
        }
        ```
3. **Kafka Broker**
    - Decouple synchronous request processing from downstream metric ingestion.
    - Topic: `credit_score_metrics`
    - Schema:
        ```json
        {
            "applicant_id": "12345",
            "champion_score": 0.73,
            "challenger_score": 0.68,
            "latency_champion_ms": 120,
            "latency_challenger_ms": 135,
            "decision_diff": 0.05,
            "timestamp": "2025-10-07T13:35:42Z",
            "status": "OK"
        }

        ```
4. **Metrics Persistence Service**
    - Subscribe to Kafka topic credit_score_metrics.
    - Deserialize messages and insert them into PostgreSQL.
    - Expose REST endpoints for aggregated queries:
        - `/metrics/summary`
        - `/metrics/latency`
        - `/metrics/distribution`
5. **PostgreSQL Database**
    - Store metrics for each scoring event (Champion/Challenger results, latency, timestamp).
    - Schema: `credit_score_metrics`
        | Field | Type | Description |
        |-------|------|-------------|
        | `id` | SERIAL | Primary key |
        | `applicant_id` | TEXT | Request identifier |
        | `champion_score` | FLOAT | Champion model score |
        | `challenger_score` | FLOAT | Challenger model score |
        | `latency_champion_ms` | INT | Champion latency |
        | `latency_challenger_ms` | INT | Challenger latency |
        | `decision_diff` | FLOAT | Absolute score delta |
        | `timestamp` | TIMESTAMP | Event timestamp |
        | `status` | TEXT | “OK”, “ERROR”, etc. |
6. **Streamlit Dashboard**
    - Provide visual analytics to compare model performance.
    - Data Source: REST API or direct SQL queries from model_metrics.
    - Sections:
        | Section                          | Visualization              | Description                                  |
        | -------------------------------- | -------------------------- | -------------------------------------------- |
        | **Summary KPIs**                 | Metric cards (`st.metric`) | Avg latency, avg score diff, total requests. |
        | **Score Distributions**          | Overlayed histograms       | Champion vs Challenger score comparison.     |
        | **Latency Trends**               | Line charts over time      | Average response time evolution per model.   |
        | **Decision Divergence**          | Pie or bar chart           | % of cases where decisions differ.           |
        | **Uplift Analysis** *(optional)* | Bar chart                  | Simulated performance uplift by segment.     |
        | **Filters**                      | Streamlit sidebar          | Filter by date, income, applicant ID, etc.   |
        | **Raw Data Table**               | `st.dataframe`             | Optional debugging and traceability.         |

## Non-Functional Requirements

| Category            | Requirement                                                               |
| ------------------- | ------------------------------------------------------------------------- |
| **Scalability**     | Kafka-based decoupling enables horizontal scaling of producers/consumers. |
| **Reliability**     | At-least-once message delivery via Kafka consumer offsets.                |
| **Latency**         | Gateway API response must remain under 200ms for Champion path.           |
| **Observability**   | Metrics persisted for dashboard visualization and audits.                 |
| **Extensibility**   | Multiple Challengers can be added easily by extending routing logic.      |
| **Fault Tolerance** | Kafka ensures message durability during downstream service outages.       |


## Expected Outcomes

By the end of implementation:
* The Gateway efficiently routes real-time credit scoring traffic.
* Metrics are streamed asynchronously to Kafka and persisted reliably.
* Stakeholders can visually compare model performance in Streamlit.
* The platform supports experimentation and evidence-based model deployment.

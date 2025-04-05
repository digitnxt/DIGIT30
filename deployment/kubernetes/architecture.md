flowchart TB
    subgraph External_Access["External Access"]
        LB[Load Balancer]
        Ingress[Ingress Controller]
    end

    subgraph API_Gateway["API Gateway & Auth"]
        Kong[Kong API Gateway]
        Keycloak[Keycloak Auth]
    end

    subgraph Core_Services["Core Services"]
        MCS[Model Context Service]
        MCP[MCP Service]
        LLM[LLM Server]
        ID[Identity Service]
        LLC[LLM Client]
    end

    subgraph Data_Services["Data Services"]
        PG[PostgreSQL]
        Kafka[Kafka]
        Redis[Redis]
    end

    subgraph Monitoring["Monitoring & Observability"]
        Prom[Prometheus]
        Graf[Grafana]
        Jaeger[Jaeger]
        KExp[Kafka Exporter]
        RExp[Redis Exporter]
    end

    subgraph Infrastructure["Infrastructure"]
        K8s[Kubernetes Cluster]
        GCP[GCP Infrastructure]
    end

    Internet((Internet)) --> LB
    LB --> Ingress
    Ingress --> Kong
    Kong --> Keycloak
    Kong --> MCS
    Kong --> MCP
    Kong --> LLM
    Kong --> ID
    Kong --> LLC
    MCS --> PG
    MCS --> Redis
    MCP --> PG
    MCP --> Kafka
    LLM --> Redis
    ID --> PG
    LLC --> LLM
    PG --> Prom
    Kafka --> KExp
    Redis --> RExp
    KExp --> Prom
    RExp --> Prom
    Prom --> Graf
    Jaeger --> Graf
    K8s --> GCP

    classDef external fill:#f9f,stroke:#333,stroke-width:2px,color:#000
    classDef gateway fill:#bbf,stroke:#333,stroke-width:2px,color:#000
    classDef core fill:#bfb,stroke:#333,stroke-width:2px,color:#000
    classDef data fill:#fbb,stroke:#333,stroke-width:2px,color:#000
    classDef monitoring fill:#ffd,stroke:#333,stroke-width:2px,color:#000
    classDef infra fill:#ddd,stroke:#333,stroke-width:2px,color:#000
    classDef subgraph fill:#fff,stroke:#333,stroke-width:2px,color:#000

    class LB,Ingress external
    class Kong,Keycloak gateway
    class MCS,MCP,LLM,ID,LLC core
    class PG,Kafka,Redis data
    class Prom,Graf,Jaeger,KExp,RExp monitoring
    class K8s,GCP infra
    class External_Access,API_Gateway,Core_Services,Data_Services,Monitoring,Infrastructure subgraph
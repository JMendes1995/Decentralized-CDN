# Decentralized Virtual CDN Architecture Choices

> This document outlines the architecture for a decentralized virtual Content Delivery Network (CDN) with opportunistic offloading. It details the technology stack and the rationale behind each choice, aiming to provide a scalable, efficient, and secure solution capable of handling global traffic with minimized latency and enhanced user experience.

## Table of Contents
- [Running Instructions](#running-instructions)
- [Infrastructure Setup](#infrastructure-setup)
    - [Global Server Deployment](#global-server-deployment)
    - [Database Selection](#database-selection)
        - [Detailed Explanation](#detailed-explanation)
    - [Storage](#storage)
    - [Content Management](#content-management)
        - [Caching Mechanism](#caching-mechanism)
        - [File Systems](#file-systems)
    - [Client-Side Offloading](#client-side-offloading)
        - [Client-Side Caching](#client-side-caching)
        - [Opportunistic Offloading & Reverse Proxy Configuration](#opportunistic-offloading--reverse-proxy-configuration)
    - [Security and Privacy](#security-and-privacy)
    - [Deployment and Monitoring](#deployment-and-monitoring)
    - [Programming and Customization](#programming-and-customization)
- [Team and Evaluation](#team)

## Running Instructions


## Infrastructure Setup

### Global Server Deployment

- **Tool**: Google Cloud Platform (GCP) & NGINX
- **Why**: GCP offers a global infrastructure that enables the deployment of NGINX servers across various regions, minimizing latency by serving content from locations closest to the user. NGINX is chosen for its high performance and ability to handle a large number of simultaneous connections efficiently.
- **Source**: [NGINX](https://nginx.org/en/), [Google Cloud Regions](https://cloud.google.com/about/locations)

## 🛠 Database Selection

- **Tools**: Cassandra, Redis, PostgreSQL
- **Why**: 
  - **Cassandra** for distributed data storage across regions, ensuring high availability and scalability.
  - **Redis** for its in-memory data structure store, used as a database and cache, facilitating quick access to frequently requested data.
  - **PostgreSQL** for reliable transactional data storage, providing robust data integrity features.
- **Source**: 
  - [Cassandra](https://cassandra.apache.org/_/index.html)
  - [Redis](https://redis.io/)
  - [PostgreSQL](https://www.postgresql.org/)

### 📖 Detailed Explanation

#### 🌐 Rationale for Using Cassandra, Redis, and PostgreSQL in a Decentralized Virtual CDN

The decision to employ a combination of Cassandra, Redis, and PostgreSQL is rooted in leveraging the unique strengths of each database system to fulfill distinct requirements of a decentralized virtual CDN with opportunistic offloading. Below is an in-depth explanation for selecting this diverse data management strategy.

#### 🌀 Cassandra

##### Advantages

- **High Availability and Scalability**: Designed for distributed environments, Cassandra ensures high availability and scalability, perfectly aligning with the global distribution needs of a CDN.
- **Write Performance**: Exceptional at handling write-heavy operations, making it ideal for logging and other write-intensive tasks.
- **Geographical Distribution**: Native support for data replication across multiple data centers enhances content delivery with minimal latency.

##### Use Cases

- Storing user activities, content access patterns, and any other write-heavy operations.
- Managing data in a globally distributed CDN infrastructure to serve users with minimal latency.

#### ⚡ Redis

##### Advantages

- **In-Memory Data Store**: Offers extremely fast read and write operations, beneficial for caching frequently accessed data to reduce access times.
- **Data Structures Support**: Accommodates complex data types like lists, sets, and hashes, facilitating advanced caching mechanisms and real-time data processing.
- **Pub/Sub Messaging**: Enables real-time communication between different components of the CDN, aiding in cache invalidation and synchronization.

##### Use Cases

- Caching content, session data, and user preferences to improve response times.
- Utilizing publish/subscribe messaging for system-wide notifications and cache management.

#### 🗃 PostgreSQL

##### Advantages

- **ACID Compliance**: Ensures reliable transaction processing, critical for managing transactional data with integrity.
- **Complex Queries and Features**: Supports advanced database functionalities like joins, views, and stored procedures for managing complex data relationships.
- **Robustness and Extensibility**: Known for its stability and a rich set of features, making it suitable for critical data storage and operations.

##### Use Cases

- Storing and managing transactional data such as user accounts, billing information, and content metadata.
- Handling complex data relationships and integrity requirements within the CDN's management layer.

#### 🌟 Combining Strengths

By integrating Cassandra, Redis, and PostgreSQL, the architecture capitalizes on:

- **Cassandra's** scalability and high availability for global data distribution.
- **Redis's** rapid data access and real-time capabilities for enhancing user experience.
- **PostgreSQL's** transactional integrity and complex data management capabilities for operational reliability.

This strategic multi-database approach ensures the architecture achieves a harmonious balance of performance, reliability, and consistency, crucial for delivering high-quality service to users worldwide.

## Storage

- **Tool**: Ceph
- **Why**: Ceph offers a highly reliable and scalable distributed storage solution, ensuring data redundancy and high availability across the CDN network.
- **Source**: [Ceph](https://ceph.io/ceph-storage/)

## Content Management

### Caching Mechanism

- **Tools**: Memcache, Redis
- **Why**: These tools provide fast access to cached data, reducing load times and decreasing the load on the backend servers. Memcache is simple and efficient for caching non-persistent data, while Redis offers more complex data structures.
- **Source**: [Memcached](https://memcached.org/), [Redis](https://redis.io/)

### File Systems

- **Tools**: ZFS, Btrfs
- **Why**: 
  - **ZFS** and **Btrfs** are chosen for their advanced features such as high reliability, data integrity, and support for high storage capacities. They offer snapshotting and data compression, beneficial for efficiently managing large volumes of content.
- **Source**: [ZFS](https://openzfs.org/), [Btrfs](https://btrfs.wiki.kernel.org/index.php/Main_Page)

## Client-Side Offloading

### Client-Side Caching

- **Tools**: Java, Rust
- **Why**: Lightweight client applications developed in these languages can efficiently manage local caching and network communication for offloading purposes. They were chosen for their performance, security features, and broad support for network programming.
- **Source**: Java, and Rust official documentation

### Opportunistic Offloading & Reverse Proxy Configuration

- **Tool**: NGINX, Custom Client Application
- **Why**: NGINX serves as the reverse proxy, efficiently managing requests between clients and servers. The custom client application, developed in Python, Java, or Rust, utilizes local caching and discovers nearby clients for content offloading, reducing the load on the central servers.
- **Source**: [NGINX Reverse Proxy](https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/)

## Security and Privacy

- **Tools**: TLS, OAuth, JWT
- **Why**: 
  - **TLS** ensures secure data transmission.
  - **OAuth** and **JWT** provide robust authentication and authorization mechanisms.
- **Source**: [TLS](https://www.openssl.org/), [OAuth](https://oauth.net/2/), [JWT](https://jwt.io/)

## Deployment and Monitoring

- **Tools**: GCP Managed Services, Kubernetes
- **Why**: Automated deployment and scaling are facilitated by GCP's managed services and Kubernetes, ensuring the CDN can efficiently handle varying loads. Monitoring tools provided by GCP offer real-time analytics and performance metrics.
- **Source**: [Google Cloud Kubernetes Engine](https://cloud.google.com/kubernetes-engine)

## Programming and Customization

- **Tool**: NGINX (C Programming)
- **Why**: Custom NGINX modules, written in C, allow for specific routing, caching, or security logic customization, integrating deeply with the CDN's functionality.
- **Source**: [Extending NGINX](https://www.nginx.com/resources/wiki/extending/)


# Team

> This is a project for the UC of Cloud Administration with the final score X out of 20
- José Carvalho up202005827@up.pt
- Jorge Mendes up202308811@up.pt


# elections
Elections FEC Data Project

Cash Flow FEC (cashflowfec.com) is a website and micro services application that I built for researching the US Federal Election Commission’s campaign finance 
data, which includes every federal candidate and political committees’ itemized transactions since 1980. Features include datasets for every individual, 
organization, candidate, and committee recorded in the FEC’s bulk data, ranking objects by funds sent or received, full search indexing, and will soon include 
social network graphing and data aggregation. The main components of this system are Envoy Proxy (http proxy & load balancer), AWS EC2 & DynamoDB, BoltDB, 
and gRPC Go and gRPC web.

The project is organized into 5 main directories / microservices:
- frontend code (html/css/js) and Envoy Proxy configuration for gRRC-web
- source: source code for each microservice (admin, server, index). Most packages in this folder are only exposed to the "admin" and "server" packages, which are 
  called by the admin & index/server apps.
- admin_app: Application inteded to run on local machine as the site administrator. Operations for building datasets from raw data, 
  viewing/uploading/deleting processed data. (note - this app is only intended to be ran on MY machine and contains unresolved errors)
- server_app: code for running web server instances and serving static files and some data from in-memory cache
- index: microservice for search index and retrieving datasets from DynamoDB.



This project is still in progress and may contain some errors in code. Currently, the following tasks are in progress:
- Editing in-line documentation for publishing
- Cleaning up code (err messages, print statements for debugging, unused/deprecated code, ect...)
- Redoing unit tests with Go testing package (old testing folder with previous unit & integration tests was removed)
- Final code checks & updates
- Deployment on EC2 using Kubernetes and Envoy as front proxy and service mesh

Once complete, the application will be available at www.cashflowfec.com. More info and documentation 
will be made available once the project is complete. Please email me at danielgarcia95367@gmail.com for any questions regarding the project.

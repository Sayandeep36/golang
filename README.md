# golang

# Environment Varibles

* GrantType= "client_credentials"
* ClientId= "your_sfdc_client_id"
* ClientSecret= "your_sfdc_client_secret"
* OAuthEndpoint="https://your_sfdc_domain.my.salesforce.com"

# Points
* Publish Saleforce change events to NATS for connecting to external systems for pubsub
* Includes handling of replay id to resume from the point the application restarted
* Includes listening to multiple topics in parallel
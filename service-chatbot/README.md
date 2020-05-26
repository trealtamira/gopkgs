# ChatBot

## Google Cloud Run webhook that push Chat messages from Cloud Builder originated Pub/Sub events.

### Overview

The system is based on multiple Google Project: (A) 'devel-services' and another project that we call (B) 'production'.
The project (A) 'devel-services' contains the chatbot service (Cloud Run) and the DNS configuration.
The project (B) 'production' is the project with the business project that contains the build trigger.
The chatbot service has been mapped via CNAME to chatbot.devtools.tre-altamira.com, its original URL is chatbot-t72fne3txq-ew.a.run.app.

The flow of the entire system is as follow:
Project B: Cloud Build run the build triggered by GitHub events
Project B: Cloud Build Push messages on the "cloud-builds" Pub/Sub topic
Project B: a Pub/Sub subscription push the message to the Chatbot webhook 
Project A: the Goolge Run service behind the webhook read the message and call the Google Chat webhook URL
The message is shown on Google Chat room

### Deployment

1. In project B configure a trigger for a build. This trigger will send messages on a Pub/Sub topic ("cloud-builds") inside the project itself.

2. Deploy in project A a Cloud Run docker image able to process Pub/Sub messages.

> `gcloud config configurations activate devel-services`  
> `docker push gcr.io/devel-services/chatbot`  
> `gcloud beta run deploy chatbot --image gcr.io/devel-services/chatbot --platform managed`  

The public endpoint of the service needs to be mapped to an owned domain trough a CNAME:

> [chatbot.{your own domain}. CNAME 300 ghs.googlehosted.com.]

The domain must be registered in the domain section of the project B (https://cloud.google.com/pubsub/docs/push#domain_ownership_validation)

3. Enable the Pub/Sub service account of project B to create JWT auth tokens:

> `gcloud config configurations activate production`  
> `gcloud projects add-iam-policy-binding PROJECT-B-ID --member=serviceAccount:service-PROJECT-B-NUMBER@gcp-sa-pubsub.iam.gserviceaccount.com --role=roles/iam.serviceAccountTokenCreator`  

4. Create in project B a service account that will be the user managing the request to the Pub/Sub subscription endpoint we are going to define:

> `gcloud iam service-accounts create chatbot-publisher --display-name "Publish message trough ChatBot in project devel-services"`  

5. Give to the account created in step 4 the ability to run the service deployed in step 2:

> `gcloud config configurations activate devel-services`  
> `gcloud beta run services add-iam-policy-binding chatbot --member=serviceAccount:chatbot-publisher@PROJECT-B-ID.iam.gserviceaccount.com --role=roles/run.invoker`  

6. Create a subscription in project B that calls the Cloud Run service of step 2 sending the Pub/Sub message:

> `gcloud beta pubsub subscriptions create chatbotSubscription --topic projects/PROJECT-B-ID/topics/cloud-builds --push-endpoint=https://chatbot.<your own domain> --push-auth-service-account=chatbot-publisher@PROJECT-B-ID.iam.gserviceaccount.com`  

Done!
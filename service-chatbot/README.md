# ChatBot

## Google Cloud Run webhook that push Chat messages from Cloud Builder originated Pub/Sub events.

### Overview

The system is based on multiple Google Project: (SER) 'service' and another project that we call (OTH) 'other'.
The project (SER) 'service' contains the chatbot service (Cloud Run) and the DNS configuration.
The project (OTH) 'other' is the project with the business project that contains the build trigger.

The flow of the entire system is as follow:

Project OTH: Cloud Build run the build triggered by GitHub events  
Project OTH: Cloud Build Push messages on the "cloud-builds" Pub/Sub topic  
Project OTH: a Pub/Sub subscription push the message to the Chatbot webhook of project SER  
Project SER: the Goolge Run service behind the webhook read the message and call the Google Chat webhook URL  

The message is shown on Google Chat room


### Deployment

1. In project OTH configure a trigger for a build. This trigger will send messages on a Pub/Sub topic ("cloud-builds") inside the project itself.

2. Deploy in project SER a Cloud Run docker image able to process Pub/Sub messages.

```
> gcloud config configurations activate <SER>
> docker push <docker image in gcr>
> gcloud beta run deploy chatbot --image <docker image in gcr> --platform managed
```

The public endpoint of the service needs to be mapped to an owned domain trough a CNAME:

[chatbot.{your own domain}. CNAME 300 ghs.googlehosted.com.]

The domain must be registered in the "domain verification" section of the project SER (https://console.cloud.google.com/apis/credentials/domainverification)

3. Enable the Pub/Sub service account of project OTH to create JWT auth tokens:

```
> gcloud config configurations activate OTH
> gcloud projects add-iam-policy-binding OTH-ID --member=serviceAccount:service-OTH-NUMBER@gcp-sa-pubsub.iam.gserviceaccount.com --role=roles/iam.serviceAccountTokenCreator
```

4. Create in project OTH a service account that will be the user used in the the Pub/Sub subscription endpoint in project OTH we are going to define:

```
> gcloud config configurations activate OTH
> gcloud iam service-accounts create chatbot-publisher --display-name "Publish message trough ChatBot in project SER"
```

5. Give to the account created in step 4 the ability to run the service deployed in step 2:

```
> gcloud config configurations activate SER
> gcloud beta run services add-iam-policy-binding chatbot --member=serviceAccount:chatbot-publisher@OTH-ID.iam.gserviceaccount.com --role=roles/run.invoker
```

6. Create the cloud-builds topic in project OTH

```
> gcloud pubsub topics create cloud-builds
```

7. Create a subscription in project OTH that calls the Cloud Run service of step 2 sending the Pub/Sub message:

```
> gcloud beta pubsub subscriptions create chatbotSubscription --topic projects/OTH-ID/topics/cloud-builds --push-endpoint=https://chatbot.<your own domain> --push-auth-service-account=chatbot-publisher@OTH-ID.iam.gserviceaccount.com
```

Done!

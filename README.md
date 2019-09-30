# Autopilot Middleware

This is the repository where I save all related with my technical test for Autopilot, the main goal for this test is to implement a Go caching middleware server for Autopilot API.

# Setup

<ul>
<li>For this project we used Redis for cache contact information from Autopilot API, so the first step will be to install a Redis in your
local environment, you can look how to install redis in https://redis.io , in order to configure all related with the redis you can modify the json file inside the directory /config named redis.json</li>
<li>I used go modules as dependency manager, I used the version go1.13. Make sure you clone the projecte under your go workspace following the directory inside src like this /github.com/ramonmacias</li>
  <li>I have a json file located in /config/autopilot_client.json, inside this file you will see a few of fields to configure the timeout for the client that connects to autopilot API and also the base url from autopilot</li>
<li>As soon as you clone the project into you local environment and start the redis server you can go to the directory /cmd/autopilot/ and use go run main.go</li>
<li>After start the API you can start testing the app on http://localhost:8080</li>
</ul>

# API information

This api have the following endpoints:

<ul>
<li>GET http://localhost/contact/id this id can be either the contact_id or email</li>
<li>PUT http://localhost/contact/id this id can be either the contact_id or email</li>
<li>POST http://localhost/contact</li>
</ul>

All the information about fields and formats related with the body requests you can find it on the API documentation from Autopilot https://autopilot.docs.apiary.io/#introduction/getting-help

All this request must have the custom authorization header named **autopilotapikey**


# Code Insisghts

I made all this API following the Clean Architecture (more information https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html), thanks to this I can test the most important part for this project, that is the switch between get contact information from the cache or from the API and also the management of cache invalidation, you can find this test inside the directory named /internal/app/usecase 

Just mention that for a more production ready version, we are going to need to test also the connection and management with redis and also a test end to end.

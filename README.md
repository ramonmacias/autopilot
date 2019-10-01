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


# Code Insights

I made all this API following the Clean Architecture (more information https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html), following this architecture we can see that there are some differents layers. In the most inner layer we can find the entities, on my case you cand find this layer on **internal/app/domain** , you will see there three different packages, the first one is **model** package where we will find only one model named contact and his definition, then in the second package named **repository** we will find and interface that should be implemented for more outter layers (in our case will be a redis implementation), the third package named **service** is where we can find the declaration of the interface for, in our case, the methods that an external API should implement, having this layer we defined the enterprise business rules.

You will find the second layer named usecase on **internal/app/usecase** , inside this package you will find a file named contact.go that has the definition about the contact usecase, also is important to mention that this usecase will need an implementation of contact repository (Redis) and implementation of external API service (Client to connect to Autopilot API), so probably this is the most important part of the project to test because this usecase determine the most important flow, in our case, the switch between get the information from a cache or get the information from the external API and the invalidation of this cache when we update or create a new contact. Using this architecture and the use of interfaces, will give us the chance to test this essential part of the code without the need of have a real connection with an external API or without having a redis running.

The third layer related with clean architecture named interface layer, cand be found on **internal/app/interface/api** , here is the implementation of this layer using handlers, the main goal of this layer is to be able to manage the data and convert it in the best way for the usecases and also for all the external agencys such as databases or webs, in our case the best format for the contact usecase and also the json format for the client that is querying our service.

The most outter layer in clean architecture frameworks and drivers can be found on **internal/app/interface/apiAutopilot** and **internal/app/interface/persistence** here is where we put all the details related more with the infrastructure in our case a redis and the external Autopilot API.

Just mention that for a more production ready version, we are going to need to test also the connection and management with redis, also a set of test cases to test the Autopilot API and finally and end to end test.

# Test

There is always a very important part in all the software projects **TEST**. So here I tried to use one of the most important golang's feature, interfaces, with the use of interfaces we can easy separate the infrastructure testing, such as databases or external connections, from the test of your internal code. So in this case I tried to achieve both testing, in /internal/app/usecase/contact_test.go and internal/app/interface/apiAutopilot/contact/api_test.go I test my own code without using any interface device such as the Autopilot API or the redis, but in internal/app/interface/persistence/redis/contactController_test.go I used a redis test server to test all the management with a specific device, on this case with a redis. So depending on the requirements from your project you are going to test more or less the infrastructure and using golang's interfaces you can test both without problems. In more production ready it will be useful to create and end to end test, testing the handlers using https://golang.org/pkg/net/http/httptest/ 

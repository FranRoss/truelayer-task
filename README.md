# Truelayer challenge

Here some explanations on my thought process for this challenge.

## Intro
I decided to use Golang as the most familiar BE language for me.
The http standard library was enough to do the project with, to note that frameworks like Chi or Gorilla have been considered but not giving much benefit in this context.

## How to run
I setup a docker image with the classical build/run pattern.
In addition there is a docker compose to facilitate. Use it with:
`docker compose up --build`.
Prefer this method as the env is being passed through it.
Note: normally I would keep everything for deployment under a folder called deploy, in this case it has been left in the root because it was only a few fiels and not differentiation by environemnt.

## Design decisions
Folder structure resemble the classical cmd/main for the entry point and everything else in internal (because the api doesn't need to expose anything for import from other parts of the repo, at least not in this context).
You can start from router and handlers first, then go down a level with the service and clients.
The main utility in here is the cache, more on this later.
And then utils, models, interfaces, conversion are the other small accesories to the api's work.

## The cache
The cache is used from all handlers globally to keep some data. There is only one instance of this in the project, that is keeping the pokemon data from the apis.
The decision to introduce was coming from a fair use policy from PokeAPI: "Locally cache resources whenever you request them.". Also it's good practice in general software development.
Using an in-mmeory solution works for the context of this exercise, in a production setup it's more likely to use a Redis istances o whatever else key-value store.
The cache is implemented as a FIFO queue where the newer items push out the oldest when reaching the limit (for the exercise is se at 3 to facilitate demonstration of it's beahviour), while keeping the entries in a map for O(1) complexity in accessing/deleting/updating items.
To note that the cache could have been used for the translation as well. The choice to not do it is because I have imagined that a translation api in a real-world scenario could have been used by a variety of actors in the project, and so it will be effective only if it's created with a big size or the sentences to translate matches between scopes. This is not a given, so I left ot out to showcase that I thought about this scenario. There is not much stopping to use it in this current setup, as for example it would greatly help with the 429 returned from the translation apis.

## The handlers
The two hanlders are structured in a way to reuse the pokemon retrieval, implemented by creating a pokemon service for it.
This makes full use of the pokemon cache as well, since the pokemon service will reuse the pokemon value between the two endpoints call.

## Tests
Unit tests have been added to the main part of the repo: cache, clients, conversion, handlers and service.
Things like utils don't give a lot of value in this exercise, but are easily coverable. Main and router are better suited for anohter kind of test like a SAT.
Note: the tests have been kept in the same folder of the thing they are testing because all folders are limited to 1-2 files. In a real worl scenario, the tests would be in a test folder in the folder they are testing for, so for example "internal/handlers" and "/internal/handlers/tests".
I left out of the testing the cases for json unmarshaling of responses, while important to consider as we are relying to external services, I wanted to keep the testing limited to the scope we are looking for in this exercise and mocked the responses for valid json accordingly.
To run the tests you can do `go test -tags=unit ./...`.

## Improvables
This is a list of improvables that this project could use:
- more coverage with SAT tests.
- cache usage for the translations (depending on the the usage, as explained in the previous section)
- tests in own folder (as explained in the previous section)